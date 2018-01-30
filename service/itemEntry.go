/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sync"

	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/service/expr"
	"github.com/yubo/gotool/list"
)

var (
	nullTs = &tsdb.TimeValuePair{}
)

type dpEntry struct {
	sync.RWMutex
	expr.ExprItem
	list   list.ListHead // point to newQueue or idxQueue trashQueue
	list_p list.ListHead // point to putQueue
	key    *tsdb.Key
	values []*tsdb.TimeValuePair

	flag     uint32
	idxTs    int64
	commitTs int64
	createTs int64
	lastTs   int64
	dataId   uint32
	endpoint string
	metric   string
	tags     string
	typ      string
}

func dpEntryNew(dp *tsdb.DataPoint) (*dpEntry, error) {
	endpoint, metric, tags, typ, err := keyAttr(dp.Key)
	if err != nil {
		return nil, err
	}

	e := &dpEntry{
		key:      dp.Key,
		values:   make([]*tsdb.TimeValuePair, CACHE_SIZE),
		createTs: timer.now(),
		endpoint: endpoint,
		metric:   metric,
		tags:     tags,
		typ:      typ,
		dataId:   CACHE_SIZE,
	}
	for i := 0; i < len(e.values); i++ {
		e.values[i] = nullTs
	}
	return e, nil

}

// called by rpc
func (p *dpEntry) put(dp *tsdb.DataPoint) error {
	p.Lock()
	defer p.Unlock()
	p.lastTs = dp.Value.Timestamp
	p.values[p.dataId&CACHE_SIZE_MASK] = dp.Value
	p.dataId += 1

	// HOOK TSDB
	return nil
}

//TODO for expr.Item interface{}
func (p *dpEntry) Get(isNum bool, num, shift_time_ int) (ret []float64) {
	p.RLock()
	defer p.RUnlock()

	var i uint32

	now := timer.now()
	shift_time := int64(shift_time_)
	id := p.dataId - 1

	if isNum {
		for i = 0; i < CACHE_SIZE; i++ {
			if now-p.values[(id-i)&CACHE_SIZE_MASK].Timestamp >= shift_time {
				break
			}
		}
		for j := 0; i < CACHE_SIZE && j < num; i++ {
			v := p.values[(id-i)&CACHE_SIZE_MASK]
			if v.Timestamp == 0 {
				break
			}
			ret = append(ret, v.Value)
			j++
		}
		return
	}

	// isSec
	sec := int64(num) + int64(shift_time)
	for i = 0; i < CACHE_SIZE; i++ {
		if now-p.values[(id-i)&CACHE_SIZE_MASK].Timestamp >= shift_time {
			break
		}
	}
	for ; i < CACHE_SIZE; i++ {
		v := p.values[(id-i)&CACHE_SIZE_MASK]
		if now-v.Timestamp >= sec || v.Timestamp == 0 {
			break
		}
		ret = append(ret, v.Value)
	}
	return ret
}

func (p *dpEntry) Nodata(isNum bool, args []float64, get expr.GetHandle) float64 {
	if timer.now()-p.lastTs <= int64(args[0]) {
		return 1
	}
	return 0
}

// TODO
func (p *dpEntry) getValues(begin, end int64) []*tsdb.TimeValuePair {
	p.Lock()
	defer p.Unlock()
	return p._getValues(CACHE_SIZE)
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *dpEntry) _getValues(n int) (ret []*tsdb.TimeValuePair) {
	var num uint32
	ret = make([]*tsdb.TimeValuePair, n)

	if n == 0 {
		return
	}

	if n > CACHE_SIZE {
		num = CACHE_SIZE
	} else {
		num = uint32(n)
	}

	//H := h & CACHE_SIZE_MASK
	begin := (p.dataId - num)

	for i := uint32(0); i < num; i++ {
		ret[i] = p.values[(i+begin)&CACHE_SIZE_MASK]
	}
	return
}

func (p *dpEntry) _getDps(n int) *tsdb.DataPoints {
	return &tsdb.DataPoints{Key: p.key, Values: p._getValues(n)}
}

func (p *dpEntry) getDps(n int) (ret *tsdb.DataPoints) {
	p.Lock()
	defer p.Unlock()

	return p._getDps(n)
}

/* the last dp(dequeue) */
func (p *dpEntry) getDp() (ret *tsdb.DataPoint) {
	p.RLock()
	defer p.RUnlock()

	//p.dataId always > 0
	return &tsdb.DataPoint{
		Key:   p.key,
		Value: p.values[(p.dataId-1)&CACHE_SIZE_MASK],
	}
}

func (p *dpEntry) String() string {
	return fmt.Sprintf("%s\n", p.key)
}
