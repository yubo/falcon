/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"container/list"
	"errors"
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
	b *Backend
}

func (p *Bkd) GetRrd(key string, rrdfile *falcon.File) (err error) {

	statsInc(ST_RPC_SERV_GETRRD, 1)
	e := p.b.cache.get(key)
	if e != nil {
		e.commit(p.b)
	}

	rrdfile.Data, err = p.b.taskFileRead(key)
	if err != nil {
		statsInc(ST_RPC_SERV_GETRRD_ERR, 1)
	}

	return
}

func (p *Bkd) Ping(req falcon.Null,
	resp *falcon.RpcResp) error {
	return nil
}

/* "Put" maybe better than "Send" */
func (p *Bkd) Put(items []*falcon.RrdItem,
	resp *falcon.RpcResp) error {
	go p.b.handleItems(items)
	return nil
}

func (p *Bkd) Send(items []*falcon.RrdItem,
	resp *falcon.RpcResp) error {
	go p.b.handleItems(items)
	return nil
}

type rpcModule struct {
	running     chan struct{}
	b           *Backend
	rpcListener *net.TCPListener
	rpcConnects connList
}

/*
func demoValue(idx int64, i int) float64 {
	return math.Sin(float64(idx+int64(i)) * math.Pi / 40.0)
}

func (p *Bkd) demoStart() {
	items := make([]*falcon.RrdItem, DEBUG_SAMPLE_NB)
	ticker := falconTicker(time.Second*DEBUG_STEP, p.b.Conf.Debug)
	step := DEBUG_STEP
	j := 0
	for {
		select {
		case _, ok := <-p.b.running:
			if !ok {
				return
			}

		case <-ticker:
			for i := 0; i < DEBUG_SAMPLE_NB; i++ {
				ts := p.b.timeNow()
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
*/

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

	resp.Type = e.typ
	resp.Step = int(e.step)

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
	caches, _ = e._getData(e.commitId, e.dataId)

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
		rrds, _ = p.taskRrdFetch(e.hashkey, param.ConsolFun,
			param.Start, param.End, e.step)
	}

	// larger than rra1point range, skip merge
	now := p.timeNow()
	if param.Start < now-now%int64(e.step)-RRA1PointCnt*int64(e.step) {
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
	if e.typ == falcon.DERIVE || e.typ == falcon.COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Ts &&
				ts == items[i+1].Ts-int64(e.step) {
				val = falcon.JsonFloat(items[i+1].V-
					items[i].V) / falcon.JsonFloat(e.step)
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
			ts = ts + int64(e.step)
		}
	} else if e.typ == falcon.GAUGE {
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
			ts = ts + int64(e.step)
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

	statsInc(ST_RPC_SERV_QUERY, 1)

	e, err = p.queryGetCacheEntry(&param, resp)
	if err != nil {
		return err
	}

	rrds, ret, err = p.queryGetData(&param, resp, e)
	if err != nil {
		statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
		return err
	}

	ret = queryPruneCache(ret, e, param.Start, param.End)

	ret = queryMergeData(rrds, ret, param.Start, param.End, int64(e.step))

	resp.Vs = queryPruneRet(ret, param.Start, param.End, int64(e.step))

	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
	return nil
}

func (p *rpcModule) prestart(b *Backend) error {
	p.running = make(chan struct{}, 0)
	p.b = b
	return nil
}

func (p *rpcModule) start(b *Backend) error {
	enable, _ := b.Conf.Configer.Bool(falcon.C_HTTP_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "rpc not enabled")
		return nil
	}

	addr, err := net.ResolveTCPAddr("tcp", b.Conf.Configer.Str(falcon.C_RPC_ADDR))
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	rpcServer := rpc.NewServer()
	rpcServer.Register(&Bkd{b: b})

	p.rpcListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, listen %s failed, %s",
			addr, err)
	} else {
		glog.Infof(MODULE_NAME+"%s rpcStart ok, listening on %s", b.Conf.Name, addr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
				// will return when rpcListener closed
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
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
			}
			tempDelay = 0
			go func() {
				e := p.rpcConnects.insert(conn)
				defer p.rpcConnects.remove(e)
				rpcServer.ServeConn(conn)
			}()
		}
	}()

	return err
}

func (p *rpcModule) stop(b *Backend) (err error) {
	if p.rpcListener == nil {
		return falcon.ErrNoent
	}

	p.rpcListener.Close()
	p.rpcConnects.Lock()
	for e := p.rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	p.rpcConnects.Unlock()
	p.rpcConnects = connList{list: list.New()}

	return nil
}

func (p *rpcModule) reload(b *Backend) error {
	// TODO
	return nil
}
