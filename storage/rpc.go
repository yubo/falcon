/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

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
	"github.com/yubo/falcon/specs"
)

var (
	rpcEvent    *specs.RoutineEvent
	rpcConnects connList
	rpcConfig   StorageOpts
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

type Storage int

func (this *Storage) GetRrd(key string, rrdfile *specs.File) (err error) {

	e := cache.get(key)
	if e != nil {
		e.commit()
	}

	rrdfile.Data, err = taskFileRead(key2filename(rpcConfig.RrdStorage, key))
	return
}

func (this *Storage) Ping(req specs.Null,
	resp *specs.RpcResp) error {
	return nil
}

/* "Put" maybe better than "Send" */
func (this *Storage) Put(items []*specs.RrdItem,
	resp *specs.RpcResp) error {
	go handleItems(items)
	return nil
}

func (this *Storage) Send(items []*specs.RrdItem,
	resp *specs.RpcResp) error {
	go handleItems(items)
	return nil
}

func queryCheckParam(param *specs.RrdQuery,
	resp *specs.RrdResp) (*cacheEntry, error) {
	// form empty response
	resp.Values = []*specs.RRDData{}
	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter

	e := cache.get(param.Csum())
	if e == nil {
		return nil, specs.ErrNoent
	}

	resp.DsType = e.dsType
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

	filename := e.filename(rpcConfig.RrdStorage)

	flag := atomic.LoadUint32(&e.flag)
	caches = make([]*specs.RRDData, len(e.cache))
	copy(caches, e.cache)

	if rpcConfig.Migrate.Enable && flag&RRD_F_MISS != 0 {
		node, _ := rrdMigrateConsistent.Get(param.Id())
		done := make(chan error, 1)
		res := &specs.RrdRespCsum{}
		rrdNetTaskCh[node] <- &netTask{
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
		rrds, _ = taskRrdFetch(filename, param.ConsolFun,
			param.Start, param.End, e.step)
	}

	// larger than rra1point range, skip merge
	now := time.Now().Unix()
	if param.Start < now-now%int64(e.step)-int64(RRA1PointCnt*e.step) {
		resp.Values = rrds
		return nil, nil, errors.New("skip merge")
	}

	// no cached caches, do not merge
	if len(caches) < 1 {
		resp.Values = rrds
		return nil, nil, errors.New("no caches")
	}
	return rrds, caches, nil
}

func queryFmtCache(items []*specs.RRDData, e *cacheEntry,
	start, end int64) (ret []*specs.RRDData) {

	// fmt cached items
	var val specs.JsonFloat

	ts := items[0].Timestamp
	n := len(items)

	last := items[n-1].Timestamp
	i := 0
	if e.dsType == specs.DERIVE || e.dsType == specs.COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Timestamp &&
				ts == items[i+1].Timestamp-int64(e.step) {
				val = specs.JsonFloat(items[i+1].Value-
					items[i].Value) / specs.JsonFloat(e.step)
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
					&specs.RRDData{Timestamp: ts, Value: val})
			}
			ts = ts + int64(e.step)
		}
	} else if e.dsType == specs.GAUGE {
		for ts <= last {
			if i < n && ts == items[i].Timestamp {
				val = specs.JsonFloat(items[i].Value)
				i++
			} else {
				// missing
				val = specs.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&specs.RRDData{Timestamp: ts, Value: val})
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
			if v.Timestamp >= start &&
				v.Timestamp <= end {
				//rrdtool返回的数据,时间戳是连续的、不会有跳点的情况
				c = append(c, v)
			}
		}
	}

	bl := len(b)
	if bl > 0 {
		cl := len(c)
		lastTs := b[0].Timestamp

		// find junction
		i := 0
		for i = cl - 1; i >= 0; i-- {
			if c[i].Timestamp < b[0].Timestamp {
				lastTs = c[i].Timestamp
				break
			}
		}

		// fix missing
		for ts := lastTs + step; ts < b[0].Timestamp; ts += step {
			c = append(c, &specs.RRDData{Timestamp: ts,
				Value: specs.JsonFloat(math.NaN())})
		}

		// merge cached items to result
		i += 1
		for j := 0; j < bl; j++ {
			if i < cl {
				if !math.IsNaN(float64(b[j].Value)) {
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
		if j < al && ts == a[j].Timestamp {
			ret[i] = a[j]
			j++
		} else {
			ret[i] = &specs.RRDData{Timestamp: ts,
				Value: specs.JsonFloat(math.NaN())}
		}
		ts += step
	}
	return ret
}

func (this *Storage) Query(param specs.RrdQuery,
	resp *specs.RrdResp) (err error) {
	var (
		e    *cacheEntry
		rrds []*specs.RRDData
		ret  []*specs.RRDData
	)

	statInc(ST_STORAGE_QUERY_CNT, 1)

	e, err = queryCheckParam(&param, resp)
	if err != nil {
		return err
	}

	rrds, ret, err = queryGetData(&param, resp, e)
	if err != nil {
		statInc(ST_STORAGE_QUERY_ITEM_CNT, len(resp.Values))
		return err
	}

	ret = queryFmtCache(ret, e, param.Start, param.End)

	ret = queryMergeData(rrds, ret, param.Start, param.End, int64(e.step))

	resp.Values = queryFmtRet(ret, param.Start, param.End, int64(e.step))

	statInc(ST_STORAGE_QUERY_ITEM_CNT, len(resp.Values))
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

	count := len(items)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		if items[i] == nil {
			continue
		}
		key := items[i].Csum()

		statInc(ST_STORAGE_RPC_RECV_CNT, 1)

		e = cache.get(key)
		if e == nil {
			e, err = cache.put(key, items[i])
			if err == nil {
				continue
			}
		}

		if items[i].Timestamp <= e.lastTs {
			continue
		}
		e.put(items[i])

		// To Index
		//ReceiveItem(items[i], checksum)

		// To History
		//AddItem(checksum, items[i])
	}
}

// 非法值: ts=0,value无意义
func getLast(csum string) *specs.RRDData {
	nan := &specs.RRDData{Timestamp: 0, Value: specs.JsonFloat(0.0)}

	e := cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	cl := len(e.cache)
	hl := len(e.history)

	if e.dsType == specs.GAUGE {
		if cl+hl < 1 {
			return nan
		}
		if cl > 0 {
			return e.cache[cl-1]
		} else {
			return e.history[hl-1]
		}
	}

	if e.dsType == specs.COUNTER || e.dsType == specs.DERIVE {
		var f0, f1 *specs.RRDData

		if cl+hl < 2 {
			return nan
		}

		if cl > 1 {
			f0 = e.cache[cl-1]
			f1 = e.cache[cl-2]
		} else if cl > 0 {
			f0 = e.cache[cl-1]
			f1 = e.history[hl-1]
		} else {
			f0 = e.history[cl-1]
			f1 = e.history[hl-2]
		}
		delta_ts := f0.Timestamp - f1.Timestamp
		delta_v := f0.Value - f1.Value
		if delta_ts != int64(e.step) || delta_ts <= 0 {
			return nan
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return &specs.RRDData{Timestamp: f0.Timestamp,
			Value: specs.JsonFloat(float64(delta_v) / float64(delta_ts))}
	}
	return nan
}

func getLastRaw(csum string) *specs.RRDData {
	nan := &specs.RRDData{Timestamp: 0, Value: specs.JsonFloat(0.0)}
	e := cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	cl := len(e.cache)
	hl := len(e.history)

	if e.dsType == specs.GAUGE {
		if cl+hl < 1 {
			return nan
		}
		if cl > 0 {
			return e.cache[cl-1]
		} else {
			return e.history[hl-1]
		}
	}
	return nan
}

func _rpcStart(config *StorageOpts, listener **net.TCPListener) (err error) {
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

func _rpcStop(config *StorageOpts, listener *net.TCPListener) (err error) {
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

func rpcStart(config StorageOpts) {
	var rpcListener *net.TCPListener

	rpc.Register(new(Storage))
	registerEventChan(rpcEvent)
	rpcConfig = config

	_rpcStart(&rpcConfig, &rpcListener)

	go func() {
		select {
		case event := <-rpcEvent.E:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				_rpcStop(&rpcConfig, rpcListener)

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
