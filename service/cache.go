/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/gotool/list"
)

// cacheEntry {{{
/* used for share memory modle */
type cacheEntry struct {
	sync.RWMutex
	flag      uint32
	list_idx  list.ListHead // point to idx0q/idx1q
	hashkey   string
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
	//e        *C.struct_cache_entry
}

// should === falcon.Item.Key()
func (p *cacheEntry) key() string {
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
func (p *cacheEntry) put(item *falcon.Item) {
	p.Lock()
	defer p.Unlock()
	p.lastTs = item.Timestamp
	id := p.dataId & CACHE_SIZE_MASK
	p.timestamp[id] = item.Timestamp
	p.value[id] = item.Value
	p.dataId += 1
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *cacheEntry) _getData(n int) (ret []*falcon.DataPoint) {
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

func (p *cacheEntry) _getItems(n int) (ret []*falcon.Item) {
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

func (p *cacheEntry) getItems(n int) (ret []*falcon.Item) {
	p.Lock()
	defer p.Unlock()

	return p._getItems(n)
}

/* the last item(dequeue) */
func (p *cacheEntry) getItem() (ret *falcon.Item) {
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

func (p *cacheEntry) String() string {
	return fmt.Sprintf("key %s endpoint %s metric %s "+
		"tags %s type %s\n",
		p.hashkey, p.endpoint, p.metric,
		p.tags, p.typ)
}

// }}}

// cacheq {{{
type cacheq struct {
	sync.RWMutex
	//size int
	head list.ListHead
}

func (p *cacheq) init() {
	//p.size = 0
	p.head.Init()
}

func (p *cacheq) addHead(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.Add(entry)
	//p.size++
}

func (p *cacheq) enqueue(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.AddTail(entry)
	//p.size++
}

func (p *cacheq) dequeue() *list.ListHead {
	p.Lock()
	defer p.Unlock()

	if p.head.Empty() {
		return nil
	}

	entry := p.head.Next
	entry.Del()
	//p.size--
	return entry
}

// }}}

// ServiceCache {{{
type serviceCache struct {
	sync.RWMutex        // hash lock
	dataq        cacheq //for flush rrddate to disk fifo
	idx0q        cacheq //immediate queue
	idx1q        cacheq //lru queue
	idx2q        cacheq //timeout queue
	data         map[string]*cacheEntry
}

func (p *serviceCache) get(key string) *cacheEntry {
	p.RLock()
	defer p.RUnlock()

	if e, ok := p.data[key]; ok {
		return e
	}
	statsInc(ST_CACHE_MISS, 1)
	return nil

}

/*
 * not idxq.size --
 */
func (p *serviceCache) unlink(key string) *cacheEntry {
	p.Lock()
	defer p.Unlock()
	e, ok := p.data[key]
	if !ok {
		return nil
	}

	e.Lock()
	defer e.Unlock()
	delete(p.data, key)

	p.dataq.Lock()
	//p.dataq.size--
	p.dataq.Unlock()
	e.list_idx.Del()

	e.hashkey = ""

	return e
}

// }}}

// cacheModule {{{
type CacheModule struct {
}

func (p *CacheModule) prestart(b *Service) error {
	glog.V(3).Infof(MODULE_NAME + " cache prestart \n")
	cache := &serviceCache{
		data: make(map[string]*cacheEntry),
	}
	cache.dataq.init()
	cache.idx0q.init()
	cache.idx1q.init()
	cache.idx2q.init()
	b.cache = cache
	return nil
}

func (p *CacheModule) start(b *Service) error {
	glog.V(3).Infof(MODULE_NAME + " cache start \n")
	return nil
}

func (p *CacheModule) stop(b *Service) error {
	//p.cache.close()
	return nil
}

func (p *CacheModule) reload(b *Service) error {
	return nil
}

// }}}
