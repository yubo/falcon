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
	"github.com/yubo/falcon/specs"
)

var (
	rpcEvent    chan specs.ProcEvent
	rpcConnects connList
	rpcConfig   BackendOpts
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

type Backend int

func (p *Backend) GetRrd(key string, rrdfile *specs.File) (err error) {

	statInc(ST_RPC_SERV_GETRRD, 1)
	e := appCache.get(key)
	if e != nil {
		e.commit()
	}

	rrdfile.Data, err = taskFileRead(key)
	if err != nil {
		statInc(ST_RPC_SERV_GETRRD_ERR, 1)
	}

	return
}

func (p *Backend) Ping(req specs.Null,
	resp *specs.RpcResp) error {
	return nil
}

func demoValue(idx int64, i int) float64 {
	return math.Sin(float64(idx+int64(i)) * math.Pi / 40.0)
}

func (p *Backend) demo() {
	items := make([]*specs.RrdItem, DEBUG_SAMPLE_NB)
	ticker := falconTicker(time.Second*DEBUG_STEP, rpcConfig.Debug)
	step := DEBUG_STEP
	j := 0
	for {
		select {
		case <-ticker:
			for i := 0; i < DEBUG_SAMPLE_NB; i++ {
				ts := timeNow()
				items[i] = &specs.RrdItem{
					Host:      "demo",
					Name:      fmt.Sprintf("%d", i),
					Value:     demoValue(ts/int64(step), i),
					TimeStemp: ts,
					Step:      step,
					Type:      specs.GAUGE,
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

/* "Put" maybe better than "Send" */
func (p *Backend) Put(items []*specs.RrdItem,
	resp *specs.RpcResp) error {
	go handleItems(items)
	return nil
}

func (p *Backend) Send(items []*specs.RrdItem,
	resp *specs.RpcResp) error {
	go handleItems(items)
	return nil
}

func queryGetCacheEntry(param *specs.RrdQuery,
	resp *specs.RrdResp) (*cacheEntry, error) {
	// form empty response
	resp.Vs = []*specs.RRDData{}
	resp.Host = param.Host
	resp.Name = param.Name

	e := appCache.get(param.Csum())
	if e == nil {
		return nil, specs.ErrNoent
	}

	resp.Type = e.typ
	resp.Step = e.step

	param.Start = param.Start - param.Start%int64(e.step)
	param.End = param.End - param.End%int64(e.step) + int64(e.step)
	if param.End-param.Start-int64(e.step) < 1 {
		return nil, specs.ErrParam
	}
	return e, nil
}

func queryGetData(param *specs.RrdQuery, resp *specs.RrdResp,
	e *cacheEntry) (rrds, caches []*specs.RRDData, err error) {

	flag := atomic.LoadUint32(&e.flag)
	caches, _ = e._getData(e.commitId, e.dataId)

	if rpcConfig.Migrate.Enable && flag&RRD_F_MISS != 0 {
		node, _ := storageMigrateConsistent.Get(param.Id())
		done := make(chan error, 1)
		res := &specs.RrdRespCsum{}
		storageNetTaskCh[node] <- &netTask{
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
		rrds, _ = taskRrdFetch(e.hashkey, param.ConsolFun,
			param.Start, param.End, e.step)
	}

	// larger than rra1point range, skip merge
	now := timeNow()
	if param.Start < now-now%int64(e.step)-int64(RRA1PointCnt*e.step) {
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

func queryFmtCache(items []*specs.RRDData, e *cacheEntry,
	start, end int64) (ret []*specs.RRDData) {

	// fmt cached items
	var val specs.JsonFloat

	ts := items[0].Ts
	n := len(items)

	last := items[n-1].Ts
	i := 0
	if e.typ == specs.DERIVE || e.typ == specs.COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Ts &&
				ts == items[i+1].Ts-int64(e.step) {
				val = specs.JsonFloat(items[i+1].V-
					items[i].V) / specs.JsonFloat(e.step)
				if val < 0 {
					val = specs.JsonFloat(math.NaN())
				}
				i++
			} else {
				// missing
				val = specs.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&specs.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.step)
		}
	} else if e.typ == specs.GAUGE {
		for ts <= last {
			if i < n && ts == items[i].Ts {
				val = specs.JsonFloat(items[i].V)
				i++
			} else {
				// missing
				val = specs.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&specs.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.step)
		}
	}
	return ret
}

/*
 * a older than b
 * c = a <- b
 */
func queryMergeData(a, b []*specs.RRDData, start,
	end, step int64) []*specs.RRDData {

	// do merging
	c := make([]*specs.RRDData, 0)
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
			c = append(c, &specs.RRDData{Ts: ts,
				V: specs.JsonFloat(math.NaN())})
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

func queryFmtRet(a []*specs.RRDData,
	start, end, step int64) []*specs.RRDData {

	// fmt result
	n := int((end - start) / step)
	ret := make([]*specs.RRDData, n)
	j := 0
	ts := start
	al := len(a)

	for i := 0; i < n; i++ {
		if j < al && ts == a[j].Ts {
			ret[i] = a[j]
			j++
		} else {
			ret[i] = &specs.RRDData{Ts: ts,
				V: specs.JsonFloat(math.NaN())}
		}
		ts += step
	}
	return ret
}

func (p *Backend) Query(param specs.RrdQuery,
	resp *specs.RrdResp) (err error) {
	var (
		e    *cacheEntry
		rrds []*specs.RRDData
		ret  []*specs.RRDData
	)

	statInc(ST_RPC_SERV_QUERY, 1)

	e, err = queryGetCacheEntry(&param, resp)
	if err != nil {
		return err
	}

	rrds, ret, err = queryGetData(&param, resp, e)
	if err != nil {
		statInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
		return err
	}

	ret = queryFmtCache(ret, e, param.Start, param.End)

	ret = queryMergeData(rrds, ret, param.Start, param.End, int64(e.step))

	resp.Vs = queryFmtRet(ret, param.Start, param.End, int64(e.step))

	statInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
	return nil
}

func handleItems(items []*specs.RrdItem) {
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

	glog.V(4).Infof("recv %d", n)
	statInc(ST_RPC_SERV_RECV, 1)
	statInc(ST_RPC_SERV_RECV_ITEM, n)

	for i := 0; i < n; i++ {
		if items[i] == nil {
			continue
		}
		key := items[i].Csum()

		e = appCache.get(key)
		if e == nil {
			e, err = appCache.createEntry(key, items[i])
			if err != nil {
				continue
			}
		}

		if DATA_TIMESTAMP_REGULATE {
			items[i].TimeStemp = items[i].TimeStemp -
				items[i].TimeStemp%int64(items[i].Step)
		}

		if items[i].TimeStemp <= e.lastTs || items[i].TimeStemp <= 0 {
			continue
		}

		e.put(items[i])
	}
}

// 非法值: ts=0,value无意义
func getLast(csum string) *specs.RRDData {
	nan := &specs.RRDData{Ts: 0, V: specs.JsonFloat(0.0)}

	e := appCache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	if e.typ == specs.GAUGE {
		if e.dataId == 0 {
			return nan
		}

		return e.data[(e.dataId-1)&CACHE_SIZE_MASK]

	}

	if e.typ == specs.COUNTER || e.typ == specs.DERIVE {

		if e.dataId < 2 {
			return nan
		}

		data, _ := e._getData(e.dataId-2, e.dataId)

		delta_ts := data[0].Ts - data[1].Ts
		delta_v := data[0].V - data[1].V
		if delta_ts != int64(e.step) || delta_ts <= 0 {
			return nan
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return &specs.RRDData{Ts: data[0].Ts,
			V: specs.JsonFloat(float64(delta_v) / float64(delta_ts))}
	}
	return nan
}

func getLastRaw(csum string) *specs.RRDData {
	nan := &specs.RRDData{Ts: 0, V: specs.JsonFloat(0.0)}
	e := appCache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	if e.typ == specs.GAUGE {
		if e.dataId == 0 {
			return nan
		}
		return e.data[(e.dataId-1)&CACHE_SIZE_MASK]
	}
	return nan
}

func _rpcStart(config *BackendOpts, listener **net.TCPListener) (err error) {
	var addr *net.TCPAddr

	if !config.Rpc {
		return nil
	}

	addr, err = net.ResolveTCPAddr("tcp", config.RpcAddr)
	if err != nil {
		glog.Fatalf("rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	*listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf("rpc.Start error, listen %s failed, %s",
			config.RpcAddr, err)
	} else {
		glog.Infof("rpc.Start ok, listening on %s", config.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := (*listener).Accept()
			if err != nil {
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
				e := rpcConnects.insert(conn)
				defer rpcConnects.remove(e)
				rpc.ServeConn(conn)
			}()
		}
	}()
	return err
}

func _rpcStop(config *BackendOpts, listener *net.TCPListener) (err error) {
	if listener == nil {
		return specs.ErrNoent
	}

	listener.Close()
	rpcConnects.Lock()
	for e := rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	rpcConnects.Unlock()

	return nil
}

func rpcStart(config BackendOpts, p *specs.Process) {
	var rpcListener *net.TCPListener

	s := new(Backend)
	rpc.Register(s)
	p.RegisterEvent("rpc", rpcEvent)
	rpcConfig = config

	_rpcStart(&rpcConfig, &rpcListener)

	if rpcConfig.Debug > 1 {
		go s.demo()
	}

	go func() {
		select {
		case event := <-rpcEvent:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				_rpcStop(&rpcConfig, rpcListener)
				event.Done <- nil

				return
			} else if event.Method == specs.ROUTINE_EVENT_M_RELOAD {
				_rpcStop(&rpcConfig, rpcListener)

				glog.V(3).Infof("old:\n%s\n new:\n%s",
					rpcConfig, appConfig)
				rpcConfig = appConfig
				_rpcStart(&rpcConfig, &rpcListener)
			}
		}
	}()

}
