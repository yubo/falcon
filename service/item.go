/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sync"

	"github.com/yubo/falcon"
	"github.com/yubo/gotool/list"
)

type itemEntry struct { // item_t
	sync.RWMutex
	list      list.ListHead // point to newQueue or lruQueue
	shardId   int32
	flag      uint32
	idxTs     int64
	commitTs  int64
	createTs  int64
	lastTs    int64
	endpoint  []byte
	metric    []byte
	tags      []byte
	typ       falcon.ItemType
	dataId    uint32
	timestamp []int64
	value     []float64
	//tsdb hook
}

/* used for share memory modle */

// should === falcon.Item.Key()
func (p *itemEntry) key() string {
	return fmt.Sprintf("%s/%s/%s/%d",
		p.endpoint,
		p.metric,
		p.tags,
		p.typ)
}

/*
func (p *cacheEntry) csum() string {
	return falcon.Md5sum(p.key())
}
*/

// called by rpc
func (p *itemEntry) put(item *falcon.Item) error {
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

// TODO
func (p *itemEntry) getDps(begin, end int64) ([]*falcon.DataPoint, error) {
	return nil, nil
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *itemEntry) _getData(n int) (ret []*falcon.DataPoint) {
	if n == 0 {
		return
	}
	/*
		if n > p.dataId {
			n = p.dataId
		}
	*/

	if n > CACHE_SIZE {
		n = CACHE_SIZE
	}

	ret = make([]*falcon.DataPoint, n)

	//H := h & CACHE_SIZE_MASK
	offset := (p.dataId - uint32(n)) & CACHE_SIZE_MASK

	for i := 0; i < n; i++ {
		id := (uint32(i) + offset) & CACHE_SIZE_MASK
		ret[i] = &falcon.DataPoint{
			Timestamp: p.timestamp[id],
			Value:     p.value[id],
		}
	}
	return
}

func (p *itemEntry) _getItems(n int) (ret []*falcon.Item) {
	data := p._getData(n)

	for _, v := range data {
		ret = append(ret, &falcon.Item{
			Endpoint:  p.endpoint,
			Metric:    p.metric,
			Tags:      p.tags,
			Value:     v.Value,
			Timestamp: v.Timestamp,
			Type:      p.typ,
		})
	}

	return ret
}

func (p *itemEntry) getItems(n int) (ret []*falcon.Item) {
	p.Lock()
	defer p.Unlock()

	return p._getItems(n)
}

/* the last item(dequeue) */
func (p *itemEntry) getItem() (ret *falcon.Item) {
	p.RLock()
	defer p.RUnlock()

	//p.dataId always > 0
	id := uint32(p.dataId-1) & CACHE_SIZE_MASK
	return &falcon.Item{
		Endpoint:  p.endpoint,
		Metric:    p.metric,
		Tags:      p.tags,
		Value:     p.value[id],
		Timestamp: p.timestamp[id],
		Type:      p.typ,
	}
	return
}

func (p *itemEntry) String() string {
	return fmt.Sprintf("endpoint %s metric %s "+
		"tags %s type %s\n",
		p.endpoint, p.metric,
		p.tags, p.typ)
}
