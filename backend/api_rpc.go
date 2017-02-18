/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type connList struct {
	sync.RWMutex
	list *list.List
}

func (l *connList) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}
func (l *connList) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

type Bkd struct {
	backend *Backend
}

func (p *Bkd) GetRrd(key string, rrdfile *falcon.File) (err error) {

	statInc(ST_RPC_SERV_GETRRD, 1)
	e := p.backend.cache.get(key)
	if e != nil {
		e.commit(p.backend)
	}

	rrdfile.Data, err = p.backend.taskFileRead(key)
	if err != nil {
		statInc(ST_RPC_SERV_GETRRD_ERR, 1)
	}

	return
}

func (p *Bkd) Ping(req falcon.Null,
	resp *falcon.RpcResp) error {
	return nil
}

func demoValue(idx int64, i int) float64 {
	return math.Sin(float64(idx+int64(i)) * math.Pi / 40.0)
}

func (p *Bkd) demoStart() {
	items := make([]*falcon.RrdItem, DEBUG_SAMPLE_NB)
	ticker := falconTicker(time.Second*DEBUG_STEP, p.backend.Conf.Params.Debug)
	step := DEBUG_STEP
	j := 0
	for {
		select {
		case _, ok := <-p.backend.running:
			if !ok {
				return
			}

		case <-ticker:
			for i := 0; i < DEBUG_SAMPLE_NB; i++ {
				ts := p.backend.timeNow()
				items[i] = &falcon.RrdItem{
					Host:      "demo",
					Name:      fmt.Sprintf("%d", i),
					Value:     demoValue(ts/int64(step), i),
					TimeStemp: ts,
					Step:      step,
					Type:      falcon.GAUGE,
					Heartbeat: step * 2,
					Min:       "U",
					Max:       "U",
				}
			}
			p.Put(items, nil)
			j++
		}
	}

}

func (p *Bkd) demoStop() {
}

/* "Put" maybe better than "Send" */
func (p *Bkd) Put(items []*falcon.RrdItem,
	resp *falcon.RpcResp) error {
	go p.backend.handleItems(items)
	return nil
}

func (p *Bkd) Send(items []*falcon.RrdItem,
	resp *falcon.RpcResp) error {
	go p.backend.handleItems(items)
	return nil
}

func (p *Backend) queryGetCacheEntry(param *falcon.RrdQuery,
	resp *falcon.RrdResp) (*cacheEntry, error) {
	// form empty response
	resp.Vs = []*falcon.RRDData{}
	resp.Host = param.Host
	resp.Name = param.Name

	e := p.cache.get(param.Csum())
	if e == nil {
		return nil, falcon.ErrNoent
	}

	resp.Type = e.typ()
	resp.Step = int(e.e.step)

	param.Start = param.Start - param.Start%int64(resp.Step)
	param.End = param.End - param.End%int64(resp.Step) + int64(resp.Step)
	if param.End-param.Start-int64(resp.Step) < 1 {
		return nil, falcon.ErrParam
	}
	return e, nil
}

func (p *Backend) queryGetData(param *falcon.RrdQuery, resp *falcon.RrdResp,
	e *cacheEntry) (rrds, caches []*falcon.RRDData, err error) {

	flag := atomic.LoadUint32(&e.flag)
	caches, _ = e._getData(uint32(e.e.commitId), uint32(e.e.dataId))

	if !p.Conf.Migrate.Disabled && flag&RRD_F_MISS != 0 {
		node, _ := p.storageMigrateConsistent.Get(param.Id())
		done := make(chan error, 1)
		res := &falcon.RrdRespCsum{}
		p.storageNetTaskCh[node] <- &netTask{
			Method: NET_TASK_M_QUERY,
			Done:   done,
			Args:   param,
			Reply:  res,
		}
		<-done
		// fetch data from remote
		rrds = res.Values
	} else {
		// read data from local rrd file
		rrds, _ = p.taskRrdFetch(e.hashkey(), param.ConsolFun,
			param.Start, param.End, int(e.e.step))
	}

	// larger than rra1point range, skip merge
	now := p.timeNow()
	if param.Start < now-now%int64(e.e.step)-RRA1PointCnt*int64(e.e.step) {
		resp.Vs = rrds
		return nil, nil, errors.New("skip merge")
	}

	// no cached caches, do not merge
	if len(caches) < 1 {
		resp.Vs = rrds
		return nil, nil, errors.New("no caches")
	}
	return rrds, caches, nil
}

func queryPruneCache(items []*falcon.RRDData, e *cacheEntry,
	start, end int64) (ret []*falcon.RRDData) {

	// prune cached items
	var val falcon.JsonFloat

	ts := items[0].Ts
	n := len(items)

	last := items[n-1].Ts
	i := 0
	typ := e.typ()
	if typ == falcon.DERIVE || typ == falcon.COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Ts &&
				ts == items[i+1].Ts-int64(e.e.step) {
				val = falcon.JsonFloat(items[i+1].V-
					items[i].V) / falcon.JsonFloat(e.e.step)
				if val < 0 {
					val = falcon.JsonFloat(math.NaN())
				}
				i++
			} else {
				// missing
				val = falcon.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&falcon.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.e.step)
		}
	} else if typ == falcon.GAUGE {
		for ts <= last {
			if i < n && ts == items[i].Ts {
				val = falcon.JsonFloat(items[i].V)
				i++
			} else {
				// missing
				val = falcon.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&falcon.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.e.step)
		}
	}
	return ret
}

/*
 * a older than b
 * c = a <- b
 */
func queryMergeData(a, b []*falcon.RRDData, start,
	end, step int64) []*falcon.RRDData {

	// do merging
	c := make([]*falcon.RRDData, 0)
	if len(a) > 0 {
		for _, v := range a {
			if v.Ts >= start &&
				v.Ts <= end {
				//rrdtool返回的数据,时间戳是连续的、不会有跳点的情况
				c = append(c, v)
			}
		}
	}

	bl := len(b)
	if bl > 0 {
		cl := len(c)
		lastTs := b[0].Ts

		// find junction
		i := 0
		for i = cl - 1; i >= 0; i-- {
			if c[i].Ts < b[0].Ts {
				lastTs = c[i].Ts
				break
			}
		}

		// fix missing
		for ts := lastTs + step; ts < b[0].Ts; ts += step {
			c = append(c, &falcon.RRDData{Ts: ts,
				V: falcon.JsonFloat(math.NaN())})
		}

		// merge cached items to result
		i += 1
		for j := 0; j < bl; j++ {
			if i < cl {
				if !math.IsNaN(float64(b[j].V)) {
					c[i] = b[j]
				}
			} else {
				c = append(c, b[j])
			}
			i++
		}
	}
	return c
}

func queryPruneRet(a []*falcon.RRDData,
	start, end, step int64) []*falcon.RRDData {

	// prune result
	n := int((end - start) / step)
	ret := make([]*falcon.RRDData, n)
	j := 0
	ts := start
	al := len(a)

	for i := 0; i < n; i++ {
		if j < al && ts == a[j].Ts {
			ret[i] = a[j]
			j++
		} else {
			ret[i] = &falcon.RRDData{Ts: ts,
				V: falcon.JsonFloat(math.NaN())}
		}
		ts += step
	}
	return ret
}

func (p *Backend) Query(param falcon.RrdQuery,
	resp *falcon.RrdResp) (err error) {
	var (
		e    *cacheEntry
		rrds []*falcon.RRDData
		ret  []*falcon.RRDData
	)

	statInc(ST_RPC_SERV_QUERY, 1)

	e, err = p.queryGetCacheEntry(&param, resp)
	if err != nil {
		return err
	}

	rrds, ret, err = p.queryGetData(&param, resp, e)
	if err != nil {
		statInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
		return err
	}

	ret = queryPruneCache(ret, e, param.Start, param.End)

	ret = queryMergeData(rrds, ret, param.Start, param.End, int64(e.e.step))

	resp.Vs = queryPruneRet(ret, param.Start, param.End, int64(e.e.step))

	statInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
	return nil
}

func (p *Backend) handleItems(items []*falcon.RrdItem) {
	var (
		err error
		e   *cacheEntry
	)

	if items == nil {
		return
	}

	n := len(items)
	if n == 0 {
		return
	}

	glog.V(4).Infof(MODULE_NAME+"recv %d", n)
	statInc(ST_RPC_SERV_RECV, 1)
	statInc(ST_RPC_SERV_RECV_ITEM, n)

	for i := 0; i < n; i++ {
		if items[i] == nil {
			continue
		}
		key := items[i].Csum()

		e = p.cache.get(key)
		if e == nil {
			e, err = p.createEntry(key, items[i])
			if err != nil {
				continue
			}
		}

		if DATA_TIMESTAMP_REGULATE {
			items[i].TimeStemp = items[i].TimeStemp -
				items[i].TimeStemp%int64(items[i].Step)
		}

		if items[i].TimeStemp <= int64(e.e.lastTs) || items[i].TimeStemp <= 0 {
			continue
		}

		e.put(items[i])
	}
}

// 非法值: ts=0,value无意义
func (p *Backend) getLast(csum string) *falcon.RRDData {
	nan := &falcon.RRDData{Ts: 0, V: falcon.JsonFloat(0.0)}

	e := p.cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	typ := e.typ()
	if typ == falcon.GAUGE {
		if e.e.dataId == 0 {
			return nan
		}

		idx := uint32(e.e.dataId-1) & CACHE_SIZE_MASK
		return &falcon.RRDData{
			Ts: int64(e.e.time[idx]),
			V:  falcon.JsonFloat(e.e.value[idx]),
		}
	}

	if typ == falcon.COUNTER || typ == falcon.DERIVE {

		if e.e.dataId < 2 {
			return nan
		}

		data, _ := e._getData(uint32(e.e.dataId)-2, uint32(e.e.dataId))

		delta_ts := data[0].Ts - data[1].Ts
		delta_v := data[0].V - data[1].V
		if delta_ts != int64(e.e.step) || delta_ts <= 0 {
			return nan
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return &falcon.RRDData{Ts: data[0].Ts,
			V: falcon.JsonFloat(float64(delta_v) / float64(delta_ts))}
	}
	return nan
}

func (p *Backend) getLastRaw(csum string) *falcon.RRDData {
	nan := &falcon.RRDData{Ts: 0, V: falcon.JsonFloat(0.0)}
	e := p.cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	if e.typ() == falcon.GAUGE {
		if e.e.dataId == 0 {
			return nan
		}
		idx := uint32(e.e.dataId-1) & CACHE_SIZE_MASK
		return &falcon.RRDData{
			Ts: int64(e.e.time[idx]),
			V:  falcon.JsonFloat(e.e.value[idx]),
		}
	}
	return nan
}

func (p *Backend) rpcStart() (err error) {
	var addr *net.TCPAddr

	if !p.Conf.Params.Rpc {
		return nil
	}

	addr, err = net.ResolveTCPAddr("tcp", p.Conf.Params.RpcAddr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	p.rpcBkd = &Bkd{
		backend: p,
	}
	rpcServer := rpc.NewServer()
	rpcServer.Register(p.rpcBkd)

	p.rpcListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, listen %s failed, %s",
			p.Conf.Params.RpcAddr, err)
	} else {
		glog.Infof(MODULE_NAME+"%s rpcStart ok, listening on %s", p.Conf.Params.Name, p.Conf.Params.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
				if p.status == falcon.APP_STATUS_EXIT {
					return
				}
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			tempDelay = 0
			go func() {
				e := p.rpcConnects.insert(conn)
				defer p.rpcConnects.remove(e)
				rpcServer.ServeConn(conn)
			}()
		}
	}()

	if p.Conf.Params.Debug > 1 {
		go p.rpcBkd.demoStart()
	}

	return err
}

func (p *Backend) rpcStop() (err error) {
	if p.rpcListener == nil {
		return falcon.ErrNoent
	}

	if p.Conf.Params.Debug > 1 {
		p.rpcBkd.demoStop()
	}

	p.rpcListener.Close()
	p.rpcConnects.Lock()
	for e := p.rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	p.rpcConnects.Unlock()

	return nil
}
