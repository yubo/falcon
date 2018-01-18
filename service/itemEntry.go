/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sync"

	"github.com/yubo/falcon/service/expr"
	"github.com/yubo/gotool/list"
)

type itemEntry struct { // item_t
	sync.RWMutex
	expr.ExprItem
	list      list.ListHead // point to newQueue or idxQueue trashQueue
	list_p    list.ListHead // point to putQueue
	shardId   int32
	key       string
	flag      uint32
	idxTs     int64
	commitTs  int64
	createTs  int64
	lastTs    int64
	dataId    uint32
	timestamp []int64
	value     []float64
	endpoint  string
	metric    string
	tags      string
	typ       string

	//tsdb hook
}

func itemEntryNew(item *Item) (*itemEntry, error) {
	endpoint, metric, tags, typ, err := item.Attr()
	if err != nil {
		return nil, err
	}

	return &itemEntry{
		createTs:  timer.now(),
		endpoint:  endpoint,
		metric:    metric,
		tags:      tags,
		typ:       typ,
		key:       string(item.Key),
		shardId:   item.ShardId,
		dataId:    CACHE_SIZE,
		timestamp: make([]int64, CACHE_SIZE),
		value:     make([]float64, CACHE_SIZE),
	}, nil

}

// called by rpc
func (p *itemEntry) put(item *Item) error {
	p.Lock()
	defer p.Unlock()
	p.lastTs = item.Timestamp
	id := p.dataId & CACHE_SIZE_MASK
	p.timestamp[id] = item.Timestamp
	p.value[id] = item.Value
	p.dataId += 1

	// HOOK TSDB
	return nil
}

//TODO for expr.Item interface{}
func (p *itemEntry) Get(isNum bool, num, shift_time_ int) (ret []float64) {
	p.RLock()
	defer p.RUnlock()

	var i uint32

	now := timer.now()
	shift_time := int64(shift_time_)
	id := p.dataId - 1

	if isNum {
		for i = 0; i < CACHE_SIZE; i++ {
			if now-p.timestamp[(id-i)&CACHE_SIZE_MASK] >= shift_time {
				break
			}
		}
		for j := 0; i < CACHE_SIZE && j < num; i++ {
			k := (id - i) & CACHE_SIZE_MASK
			if p.timestamp[k] == 0 {
				break
			}
			ret = append(ret, p.value[k])
			j++
		}
		return
	}

	// isSec
	sec := int64(num) + int64(shift_time)
	for i = 0; i < CACHE_SIZE; i++ {
		if now-p.timestamp[(id-i)&CACHE_SIZE_MASK] >= shift_time {
			break
		}
	}
	for ; i < CACHE_SIZE; i++ {
		k := (id - i) & CACHE_SIZE_MASK
		if now-p.timestamp[k] >= sec || p.timestamp[k] == 0 {
			break
		}
		ret = append(ret, p.value[k])
	}
	return ret
}

func (p *itemEntry) Nodata(isNum bool, args []float64, get expr.GetHandle) float64 {
	if timer.now()-p.lastTs <= int64(args[0]) {
		return 1
	}
	return 0
}

// TODO
func (p *itemEntry) getDps(begin, end int64) ([]*DataPoint, error) {
	return nil, nil
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *itemEntry) _getData(n int) (ret []*DataPoint) {
	if n == 0 {
		return
	}

	if n > CACHE_SIZE {
		n = CACHE_SIZE
	}

	ret = make([]*DataPoint, n)

	//H := h & CACHE_SIZE_MASK
	begin := (p.dataId - uint32(n))

	for i := 0; i < n; i++ {
		id := (uint32(i) + begin) & CACHE_SIZE_MASK
		ret[i] = &DataPoint{
			Timestamp: p.timestamp[id],
			Value:     p.value[id],
		}
	}
	return
}

func (p *itemEntry) _getItems(n int) (ret []*Item) {
	data := p._getData(n)

	for _, v := range data {
		ret = append(ret, &Item{
			Key:       []byte(p.key),
			Value:     v.Value,
			Timestamp: v.Timestamp,
		})
	}

	return ret
}

func (p *itemEntry) getItems(n int) (ret []*Item) {
	p.Lock()
	defer p.Unlock()

	return p._getItems(n)
}

/* the last item(dequeue) */
func (p *itemEntry) getItem() (ret *Item) {
	p.RLock()
	defer p.RUnlock()

	//p.dataId always > 0
	id := uint32(p.dataId-1) & CACHE_SIZE_MASK
	return &Item{
		Key:       []byte(p.key),
		Value:     p.value[id],
		Timestamp: p.timestamp[id],
	}
	return
}

func (p *itemEntry) String() string {
	return fmt.Sprintf("%s\n", p.key)
}
