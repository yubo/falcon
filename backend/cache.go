/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

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
	flag uint32
	// point to dataq/poolq
	list_data list.ListHead
	// point to idx0q/idx1q
	list_idx list.ListHead // no init queue
	//e        *C.struct_cache_entry
	hashkey   string
	idxTs     int64
	commitTs  int64
	createTs  int64
	lastTs    int64
	host      string
	name      string
	tags      string
	typ       string
	step      int
	heartbeat int
	min       byte
	max       byte
	dataId    uint32
	commitId  uint32
	time      []int64
	value     []float64
}

// should === falcon.RrdItem.Id()
func (p *cacheEntry) id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d",
		p.host,
		p.name,
		p.tags,
		p.typ,
		int(p.step))
}

func (p *cacheEntry) csum() string {
	return falcon.Md5sum(p.id())
}

// called by rpc
func (p *cacheEntry) put(item *falcon.RrdItem) {
	p.Lock()
	defer p.Unlock()
	p.lastTs = item.TimeStemp
	idx := p.dataId & CACHE_SIZE_MASK
	p.time[idx] = item.TimeStemp
	p.value[idx] = item.Value
	p.dataId += 1
}

// fetch remote ds
func (p *cacheEntry) fetchCommit(b *Backend) {
	done := make(chan error)

	node, err := b.storageMigrateConsistent.Get(p.hashkey)
	if err != nil {
		return
	}

	b.storageNetTaskCh[node] <- &netTask{
		Method: NET_TASK_M_FETCH_COMMIT,
		e:      p,
		Done:   done,
	}

	// net_task slow, shouldn't block commitCache()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			glog.Warning(MODULE_NAME+"get %s from remote err[%s]\n", p.hashkey, err)
			return
		}
		//todo: flushfile after getfile? not yet
	}()
}

func (p *cacheEntry) createRrd(b *Backend) error {
	done := make(chan error, 1)

	b.ktoch(p.hashkey) <- &ioTask{
		method: IO_TASK_M_RRD_ADD,
		args:   p,
		done:   done,
	}
	err := <-done

	p.commitTs = b.timeNow()

	return err
}

func (p *cacheEntry) commit(b *Backend) error {
	done := make(chan error, 1)

	b.ktoch(p.hashkey) <- &ioTask{
		method: IO_TASK_M_RRD_UPDATE,
		args:   p,
		done:   done,
	}
	err := <-done

	p.commitTs = b.timeNow()

	return err
}

func (p *cacheEntry) filename(b *Backend) string {
	return b.ktofname(p.hashkey)
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *cacheEntry) _getData(l, h uint32) (ret []*falcon.RRDData,
	overrun int) {

	size := h - l
	if size > CACHE_SIZE {
		overrun = int(size - CACHE_SIZE)
		size = CACHE_SIZE
		l = h - CACHE_SIZE
	}

	if size == 0 {
		return
	}

	ret = make([]*falcon.RRDData, size)

	//H := h & CACHE_SIZE_MASK
	L := l & CACHE_SIZE_MASK

	for i := uint32(0); i < size; i++ {
		idx := (L + i) & CACHE_SIZE_MASK
		ret[i] = &falcon.RRDData{
			Ts: int64(p.time[idx]),
			V:  falcon.JsonFloat(p.value[idx]),
		}
	}
	/*
		if H > L {
			copy(ret, p.data[L:H])
		} else {
			copy(ret[:CACHE_SIZE-L], p.data[L:])
			copy(ret[CACHE_SIZE-L:], p.data[:H])
		}
	*/
	return
}

func (p *cacheEntry) _dequeueAll() []*falcon.RRDData {
	ret, over := p._getData(p.commitId, p.dataId)
	p.commitId = p.dataId
	if over > 0 {
		statsInc(ST_CACHE_OVERRUN, over)
	}

	return ret
}

func (p *cacheEntry) dequeueAll() []*falcon.RRDData {
	p.Lock()
	defer p.Unlock()

	return p._dequeueAll()
}

func (p *cacheEntry) _getItems() (ret []*falcon.RrdItem) {

	rrds, _ := p._getData(0, p.dataId)

	for _, v := range rrds {
		ret = append(ret, &falcon.RrdItem{
			Host:      p.host,
			Name:      p.name,
			Tags:      p.tags,
			Value:     float64(v.V),
			TimeStemp: v.Ts,
			Type:      p.typ,
			Step:      p.step,
		})
	}

	return ret
}

func (p *cacheEntry) getItems() (ret []*falcon.RrdItem) {
	p.Lock()
	defer p.Unlock()

	return p._getItems()
}

/* the last item(dequeue) */
func (p *cacheEntry) getItem() (ret *falcon.RrdItem) {
	p.RLock()
	defer p.RUnlock()

	//p.dataId always > 0
	idx := uint32(p.dataId-1) & CACHE_SIZE_MASK
	return &falcon.RrdItem{
		Host:      p.host,
		Name:      p.name,
		Tags:      p.tags,
		Value:     p.value[idx],
		TimeStemp: p.time[idx],
		Type:      p.typ,
		Step:      p.step,
	}
	return
}

func (p *cacheEntry) String() string {
	return fmt.Sprintf("key:%s host:%s name:%s "+
		"tags:%s type:%s step:%d\n",
		p.hashkey, p.host, p.name,
		p.tags, p.typ, p.step)
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

// backendCache {{{
type backendCache struct {
	sync.RWMutex        // hash lock
	dataq        cacheq //for flush rrddate to disk fifo
	poolq        cacheq //free entry lifo
	idx0q        cacheq //immediate queue
	idx1q        cacheq //lru queue
	idx2q        cacheq //timeout queue
	hash         map[string]*cacheEntry
}

func (p *backendCache) get(key string) *cacheEntry {
	p.RLock()
	defer p.RUnlock()

	if e, ok := p.hash[key]; ok {
		return e
	}
	statsInc(ST_CACHE_MISS, 1)
	return nil

}

/*
 * not idxq.size --
 */
func (p *backendCache) unlink(key string) *cacheEntry {
	p.Lock()
	defer p.Unlock()
	e, ok := p.hash[key]
	if !ok {
		return nil
	}

	e.Lock()
	defer e.Unlock()
	delete(p.hash, key)

	p.dataq.Lock()
	e.list_data.Del()
	//p.dataq.size--
	p.dataq.Unlock()
	e.list_idx.Del()

	p.poolq.addHead(&e.list_data)
	e.hashkey = ""

	return e
}

// }}}

// cacheModule {{{
type CacheModule struct {
}

func (p *CacheModule) prestart(b *Backend) error {
	glog.V(3).Infof(MODULE_NAME + " cache prestart \n")
	cache := &backendCache{
		hash: make(map[string]*cacheEntry),
	}
	cache.dataq.init()
	cache.poolq.init()
	cache.idx0q.init()
	cache.idx1q.init()
	cache.idx2q.init()
	b.cache = cache
	return nil
}

func (p *CacheModule) start(b *Backend) error {
	glog.V(3).Infof(MODULE_NAME + " cache start \n")
	return nil
}

func (p *CacheModule) stop(b *Backend) error {
	//p.cache.close()
	return nil
}

func (p *CacheModule) reload(b *Backend) error {
	return nil
}

// }}}
