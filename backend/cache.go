/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/list"
)

var (
	/* cache */
	appCache backendCache
)

type cacheEntry struct {
	sync.RWMutex
	flag      uint32
	list_data list.ListHead
	list_idx  list.ListHead // no init queue
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
	min       string
	max       string
	dataId    uint32
	commitId  uint32
	data      [CACHE_SIZE]*specs.RRDData
}

// should === specs.RrdItem.Id()
func (p *cacheEntry) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.host, p.name, p.tags, p.typ, p.step)
}

type cacheq struct {
	sync.RWMutex
	size int
	head list.ListHead
}

type backendCache struct {
	sync.RWMutex // hash lock
	dataq        cacheq
	idx0q        cacheq //index queue
	idx1q        cacheq //timeout queue
	hash         map[string]*cacheEntry
}

func (p *backendCache) init() {
	p.hash = make(map[string]*cacheEntry)
	p.dataq.init()
	p.idx0q.init()
	p.idx1q.init()
}

func (p *backendCache) get(key string) *cacheEntry {
	p.RLock()
	defer p.RUnlock()

	if e, ok := p.hash[key]; ok {
		return e
	}
	statInc(ST_CACHE_MISS, 1)
	return nil

}

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
	p.dataq.size--
	p.dataq.Unlock()

	if e.idxTs == 0 {
		p.idx0q.Lock()
		defer p.idx0q.Unlock()
		p.idx0q.size--
	} else {
		p.idx1q.Lock()
		defer p.idx1q.Unlock()
		p.idx1q.size--
	}
	e.list_idx.Del()

	return e
}

func list_data_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_data))))
}

func list_idx_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_idx))))
}

func (p *cacheq) init() {
	p.size = 0
	p.head.Init()
}

func (p *cacheq) addHead(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.Add(entry)
	p.size++
}

func (p *cacheq) enqueue(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.AddTail(entry)
	p.size++
}

func (p *cacheq) dequeue() *list.ListHead {
	p.Lock()
	defer p.Unlock()

	if p.size == 0 {
		return nil
	}

	entry := p.head.Next
	entry.Del()
	p.size--
	return entry
}

// called by rpc
func (p *backendCache) createEntry(key string, item *specs.RrdItem) (*cacheEntry, error) {

	statInc(ST_CACHE_CREATE, 1)
	p.Lock()
	if e, ok := p.hash[key]; ok {
		return e, specs.ErrExist
	}

	e := &cacheEntry{
		hashkey:   key,
		createTs:  timeNow(),
		host:      item.Host,
		name:      item.Name,
		tags:      item.Tags,
		typ:       item.Type,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       item.Min,
		max:       item.Max,
		dataId:    0,
		commitId:  0,
	}
	p.hash[key] = e
	p.Unlock()

	p.dataq.enqueue(&e.list_data)
	p.idx0q.addHead(&e.list_idx)

	if rpcConfig.Migrate.Enable {
		_, err := os.Stat(e.filename())
		if os.IsNotExist(err) {
			e.flag = RRD_F_MISS
		}
	}
	return e, nil
}

// called by rpc
func (p *cacheEntry) put(item *specs.RrdItem) {
	p.Lock()
	defer p.Unlock()
	p.lastTs = item.TimeStemp
	p.data[p.dataId&CACHE_SIZE_MASK] = &specs.RRDData{
		Ts: item.TimeStemp,
		V:  specs.JsonFloat(item.Value),
	}
	p.dataId += 1
}

// fetch remote ds
func (p *cacheEntry) fetchCommit() {
	done := make(chan error)

	node, err := storageMigrateConsistent.Get(p.hashkey)
	if err != nil {
		return
	}

	storageNetTaskCh[node] <- &netTask{
		Method: NET_TASK_M_FETCH_COMMIT,
		e:      p,
		Done:   done,
	}

	// net_task slow, shouldn't block commitCache()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			glog.Warning("get %s from remote err[%s]\n", p.hashkey, err)
			return
		}
		//todo: flushfile after getfile? not yet
	}()
}

func (p *cacheEntry) commit() error {
	done := make(chan error, 1)

	ktoch(p.hashkey) <- &ioTask{
		method: IO_TASK_M_RRD_UPDATE,
		args:   p,
		done:   done,
	}
	err := <-done

	p.commitTs = timeNow()

	return err
}

func (p *cacheEntry) filename() string {
	return ktofname(p.hashkey)
}

// return [l, h)
// h - l <= CACHE_SIZE
func (p *cacheEntry) _getData(l, h uint32) (ret []*specs.RRDData,
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

	ret = make([]*specs.RRDData, size)

	H := h & CACHE_SIZE_MASK
	L := l & CACHE_SIZE_MASK

	if H > L {
		copy(ret, p.data[L:H])
	} else {
		copy(ret[:CACHE_SIZE-L], p.data[L:])
		copy(ret[CACHE_SIZE-L:], p.data[:H])
	}
	return
}

func (p *cacheEntry) _dequeueAll() []*specs.RRDData {
	ret, over := p._getData(p.commitId, p.dataId)
	p.commitId = p.dataId
	if over > 0 {
		statInc(ST_CACHE_OVERRUN, over)
	}

	return ret
}

func (p *cacheEntry) dequeueAll() []*specs.RRDData {
	p.Lock()
	defer p.Unlock()

	return p._dequeueAll()
}

func (p *cacheEntry) _getItems() (ret []*specs.RrdItem) {

	rrds, _ := p._getData(0, p.dataId)

	for _, v := range rrds {
		ret = append(ret, &specs.RrdItem{
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

func (p *cacheEntry) getItems() (ret []*specs.RrdItem) {
	p.Lock()
	defer p.Unlock()

	return p._getItems()
}

/* the last item(dequeue) */
func (p *cacheEntry) getItem() (ret *specs.RrdItem) {
	p.RLock()
	defer p.RUnlock()

	//p.dataId always > 0
	v := p.data[(p.dataId-1)&CACHE_SIZE_MASK]
	return &specs.RrdItem{
		Host:      p.host,
		Name:      p.name,
		Tags:      p.tags,
		Value:     float64(v.V),
		TimeStemp: v.Ts,
		Type:      p.typ,
		Step:      p.step,
	}
	return
}

func getItems(key string) (ret []*specs.RrdItem) {
	e := appCache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func getLastItem(key string) (ret *specs.RrdItem) {
	e := appCache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}
