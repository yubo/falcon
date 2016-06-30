/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"errors"
	"math"
	"sync/atomic"
	"time"

	"github.com/yubo/falcon/specs"
)

type Graph int

func (this *Graph) GetRrd(filename, key string, rrdfile *specs.File) (err error) {

	e := cache.get(key)
	if e != nil {
		e.commit()
	}

	rrdfile.Body, err = taskReadFile(key2filename(config().RrdStorage, key))
	return
}

func (this *Graph) Ping(req specs.NullRpcRequest,
	resp *specs.SimpleRpcResponse) error {
	return nil
}

/* "Put" maybe better than "Send" */
func (this *Graph) Put(items []*specs.GraphItem,
	resp *specs.SimpleRpcResponse) error {
	go handleItems(items)
	return nil
}

func (this *Graph) Send(items []*specs.GraphItem,
	resp *specs.SimpleRpcResponse) error {
	go handleItems(items)
	return nil
}

func query_check_param(param *specs.GraphQueryParam, resp *specs.GraphQueryResponse) (*cacheEntry, error) {
	// form empty response
	resp.Values = []*specs.RRDData{}
	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter

	e := cache.get(param.Key())
	if e == nil {
		return nil, ErrNoent
	}

	resp.DsType = e.dsType
	resp.Step = e.step

	param.Start = param.Start - param.Start%int64(e.step)
	param.End = param.End - param.End%int64(e.step) + int64(e.step)
	if param.End-param.Start-int64(e.step) < 1 {
		return nil, ErrParam
	}
	return e, nil
}

func query_get_data(param *specs.GraphQueryParam,
	resp *specs.GraphQueryResponse, e *cacheEntry) (rrds, caches []*specs.RRDData, err error) {

	filename := e.filename(config().RrdStorage)

	flag := atomic.LoadUint32(&e.flag)
	caches = make([]*specs.RRDData, len(e.cache))
	copy(caches, e.cache)

	if config().Migrate.Enable && flag&GRAPH_F_MISS != 0 {
		node, _ := Consistent.Get(param.Endpoint + "/" + param.Counter)
		done := make(chan error, 1)
		res := &specs.GraphAccurateQueryResponse{}
		Net_task_ch[node] <- &Net_task_t{
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
		rrds, _ = taskCheckout(filename, param.ConsolFun, param.Start, param.End, e.step)
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

func query_fmt_cache(items []*specs.RRDData, e *cacheEntry, start, end int64) (ret []*specs.RRDData) {

	// fmt cached items
	var val specs.JsonFloat

	ts := items[0].Timestamp
	n := len(items)

	last := items[n-1].Timestamp
	i := 0
	if e.dsType == DERIVE || e.dsType == COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Timestamp &&
				ts == items[i+1].Timestamp-int64(e.step) {
				val = specs.JsonFloat(items[i+1].Value-items[i].Value) / specs.JsonFloat(e.step)
				if val < 0 {
					val = specs.JsonFloat(math.NaN())
				}
				i++
			} else {
				// missing
				val = specs.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret, &specs.RRDData{Timestamp: ts, Value: val})
			}
			ts = ts + int64(e.step)
		}
	} else if e.dsType == GAUGE {
		for ts <= last {
			if i < n && ts == items[i].Timestamp {
				val = specs.JsonFloat(items[i].Value)
				i++
			} else {
				// missing
				val = specs.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				ret = append(ret, &specs.RRDData{Timestamp: ts, Value: val})
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
func query_merge_data(a, b []*specs.RRDData, start, end, step int64) []*specs.RRDData {

	// do merging
	c := make([]*specs.RRDData, 0)
	if len(a) > 0 {
		for _, v := range a {
			if v.Timestamp >= start && v.Timestamp <= end {
				c = append(c, v) //rrdtool返回的数据,时间戳是连续的、不会有跳点的情况
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
			c = append(c, &specs.RRDData{Timestamp: ts, Value: specs.JsonFloat(math.NaN())})
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

func query_fmt_ret(a []*specs.RRDData, start, end, step int64) []*specs.RRDData {

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
			ret[i] = &specs.RRDData{Timestamp: ts, Value: specs.JsonFloat(math.NaN())}
		}
		ts += step
	}
	return ret
}

func (this *Graph) Query(param specs.GraphQueryParam, resp *specs.GraphQueryResponse) (err error) {
	var (
		e    *cacheEntry
		rrds []*specs.RRDData
		ret  []*specs.RRDData
	)

	stat_inc(ST_GraphQueryCnt, 1)

	e, err = query_check_param(&param, resp)
	if err != nil {
		return err
	}

	rrds, ret, err = query_get_data(&param, resp, e)
	if err != nil {
		stat_inc(ST_GraphQueryItemCnt, len(resp.Values))
		return err
	}

	ret = query_fmt_cache(ret, e, param.Start, param.End)

	ret = query_merge_data(rrds, ret, param.Start, param.End, int64(e.step))

	resp.Values = query_fmt_ret(ret, param.Start, param.End, int64(e.step))

	stat_inc(ST_GraphQueryItemCnt, len(resp.Values))
	return nil
}

/*
func (this *Graph) Info(param specs.GraphInfoParam, resp *specs.GraphInfoResp) error {
	// statistics
	stat_inc(ST_GraphInfoCnt, 1)

	csum := param.Key()

	e := cache.get(csum)
	if e == nil {
		return nil
	}

	resp.ConsolFun = e.dsType
	resp.Step = e.step
	resp.Filename = csum2filename(config().RrdStorage, csum, e.dsType, e.step)

	return nil
}

func (this *Graph) Last(param specs.GraphLastParam, resp *specs.GraphLastResp) error {
	// statistics
	stat_inc(ST_GraphLastCnt, 1)

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = getLast(param.Checksum())

	return nil
}

func (this *Graph) LastRaw(param specs.GraphLastParam, resp *specs.GraphLastResp) error {
	// statistics
	stat_inc(ST_GraphLastRawCnt, 1)

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = getLastRaw(param.Checksum())

	return nil
}
*/
func handleItems(items []*specs.GraphItem) {
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
		key := items[i].Key()

		stat_inc(ST_GraphRpcRecvCnt, 1)

		// To Graph
		e = cache.get(key)
		if e == nil {
			e, err = cache.put(key, items[i])
			if err == nil {
				continue
			}
		}

		if items[i].Timestamp <= e.putTs {
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

	if e.dsType == GAUGE {
		if cl+hl < 1 {
			return nan
		}
		if cl > 0 {
			return e.cache[cl-1]
		} else {
			return e.history[hl-1]
		}
	}

	if e.dsType == COUNTER || e.dsType == DERIVE {
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

	if e.dsType == GAUGE {
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
