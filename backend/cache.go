/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

/*
#include "cache.h"
*/
import "C"

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/list"
)

/* used for share memory modle */
type cacheEntry struct {
	sync.RWMutex
	flag uint32
	// point to dataq/poolq
	list_data list.ListHead
	// point to idx0q/idx1q
	list_idx list.ListHead // no init queue
	e        *C.struct_cache_entry
}

func (p *cacheEntry) setHashkey(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.set_hashkey(p.e, cs))
}

func (p *cacheEntry) setHost(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.set_host(p.e, cs))
}

func (p *cacheEntry) setName(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.set_name(p.e, cs))
}

func (p *cacheEntry) setTags(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.set_tags(p.e, cs))
}

func (p *cacheEntry) setTyp(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.set_typ(p.e, cs))
}

func (p *cacheEntry) hashkey() string {
	return C.GoString((*C.char)(unsafe.Pointer(&p.e.hashkey[0])))
}

func (p *cacheEntry) host() string {
	return C.GoString((*C.char)(unsafe.Pointer(&p.e.host[0])))
}

func (p *cacheEntry) name() string {
	return C.GoString((*C.char)(unsafe.Pointer(&p.e.name[0])))
}

func (p *cacheEntry) tags() string {
	return C.GoString((*C.char)(unsafe.Pointer(&p.e.tags[0])))
}

func (p *cacheEntry) typ() string {
	return C.GoString((*C.char)(unsafe.Pointer(&p.e.typ[0])))
}

// should === specs.RrdItem.Id()
func (p *cacheEntry) id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d",
		p.host(),
		p.name(),
		p.tags(),
		p.typ(),
		int(p.e.step))
}

func (p *cacheEntry) csum() string {
	return specs.Md5sum(p.id())
}

type cacheq struct {
	sync.RWMutex
	//size int
	head list.ListHead
}

type backendCache struct {
	sync.RWMutex   // hash lock
	cache_entry_nb int
	startkey       int
	endkey         int    //[startkey, endkey), endkey not used
	dataq          cacheq //for flush rrddate to disk fifo
	poolq          cacheq //free entry lifo
	idx0q          cacheq //immediate queue
	idx1q          cacheq //lru queue
	idx2q          cacheq //timeout queue
	hash           map[string]*cacheEntry
	shmaddrs       []uintptr
	shmids         []int
}

func (p *backendCache) close() {
	for _, addr := range p.shmaddrs {
		Shmdt(addr)
	}
	glog.V(3).Infof(MODULE_NAME+"shmdt %d", len(p.shmaddrs))
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
	e.setHashkey("")

	return e
}

func (p *Backend) importBlocks() (int, error) {
	var (
		shmid int
		size  int
		addr  uintptr
		block *C.struct_cache_block
		err   error
		cache *backendCache
	)
	cache = p.cache
	cache.Lock()
	defer cache.Unlock()

	for {
		shmid, size, err = Shmget(cache.endkey, p.ShmSize, C.IPC_CREAT|0700)
		if err != nil {
			return cache.endkey - cache.startkey, err
		}
		cache.endkey++

		if size != p.ShmSize {
			glog.Errorf(MODULE_NAME+"segment size error, remove "+
				"segment from 0x%x and retry", cache.endkey)
			cleanBlocks(cache.endkey)
			return cache.endkey - cache.startkey, err
		}

		addr, err = Shmat(shmid, 0, 0)
		if err != nil {
			return cache.endkey - cache.startkey, err
		}

		block = (*C.struct_cache_block)(unsafe.Pointer(addr))

		for uint32(block.magic) != p.ShmMagic {
			glog.Infof(MODULE_NAME+"set magic head at 0x%x",
				shmid)
			block.magic = C.uint32_t(p.ShmMagic)
			return cache.endkey - cache.startkey, specs.EFMT
		}

		cache.shmaddrs = append(cache.shmaddrs, addr)
		cache.shmids = append(cache.shmids, shmid)
		list := &cache.poolq

		for i := 0; i < int(block.cache_entry_nb); i++ {
			e := &cacheEntry{
				e: (*C.struct_cache_entry)(unsafe.Pointer(addr +
					uintptr(C.sizeof_struct_cache_block) +
					uintptr(i*C.sizeof_struct_cache_entry))),
			}
			key := e.hashkey()
			if key == "" {
				// add to pool
				list.head.AddTail(&e.list_data)
			} else {
				// restore to cacheEntry
				cache.hash[key] = e
				cache.dataq.enqueue(&e.list_data)
				cache.idx0q.enqueue(&e.list_idx)
			}
		}
	}
	return cache.endkey - cache.startkey, nil
}

func cleanBlocks(start int) (int, error) {
	var (
		shmid int
		err   error
		i     int
	)

	for i = start; ; i++ {
		shmid, _, err = Shmget(i, 0, 0700)
		if err != nil {
			return i - start, nil
		}
		Shmrm(shmid)
	}
	return i - start, err
}

func (p *Backend) cacheReset() error {
	cache := p.cache
	cache.hash = make(map[string]*cacheEntry)
	cache.shmaddrs = make([]uintptr, 0)
	cache.shmids = make([]int, 0)
	cache.dataq.init()
	cache.poolq.init()
	cache.idx0q.init()
	cache.idx1q.init()
	cache.idx2q.init()
	cache.startkey = p.ShmKey
	cache.endkey = p.ShmKey
	cache.cache_entry_nb = (p.ShmSize - C.sizeof_struct_cache_block) /
		C.sizeof_struct_cache_entry
	n, err := cleanBlocks(cache.startkey)
	if err != nil {
		return err
	}
	glog.V(3).Infof(MODULE_NAME+"cleanblocks %d", n)
	return nil
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

func (p *Backend) allocPool() (err error) {
	var (
		shmid int
		size  int
		addr  uintptr
		block *C.struct_cache_block
		cache *backendCache
	)
	cache = p.cache
	cache.Lock()
	defer cache.Unlock()

	shmid, size, err = Shmget(cache.endkey, p.ShmSize, C.IPC_CREAT|0700)
	if err != nil {
		glog.Errorf(MODULE_NAME+"shmget(%d, %d, ipc_create ) err %s",
			cache.endkey, p.ShmSize, err)
		return err
	}
	cache.endkey++

	if size != p.ShmSize {
		glog.Errorf(MODULE_NAME+"alloc shm size %d want %d, remove and try again",
			size, p.ShmSize)
		cleanBlocks(cache.endkey)
		shmid, _, err = Shmget(cache.endkey, p.ShmSize, C.IPC_CREAT|0700)
		if err != nil {
			glog.Errorf(MODULE_NAME+"shmget(%d, %d, ipc_create ) err %s",
				cache.endkey, p.ShmSize, err)
			return err
		}
	}

	addr, err = Shmat(shmid, 0, 0)
	if err != nil {
		glog.Errorf(MODULE_NAME+"shmat(%d, 0, 0 ) err %s", shmid, err)
		return err
	}
	block = (*C.struct_cache_block)(unsafe.Pointer(addr))

	if uint32(block.magic) == p.ShmMagic {
		Shmdt(addr)
		glog.Errorf(MODULE_NAME+"magic head already set at shmid %d", shmid)
		return specs.ErrExist
	}
	glog.V(5).Infof(MODULE_NAME+"alloc shm segment, endkey 0x%08x shmid %d len(shmaddr) %d",
		cache.endkey, shmid, len(cache.shmaddrs))

	cache.shmaddrs = append(cache.shmaddrs, addr)
	cache.shmids = append(cache.shmids, shmid)

	block.magic = C.uint32_t(p.ShmMagic)
	block.block_size = C.int(p.ShmSize)
	block.cache_entry_nb = C.int(cache.cache_entry_nb)

	list := &cache.poolq
	list.Lock()
	defer list.Unlock()
	for i := 0; i < int(block.cache_entry_nb); i++ {
		e := &cacheEntry{
			e: (*C.struct_cache_entry)(unsafe.Pointer(addr +
				uintptr(C.sizeof_struct_cache_block) +
				uintptr(i*C.sizeof_struct_cache_entry))),
		}
		list.head.AddTail(&e.list_data)
	}
	//list.size += int(block.cache_entry_nb)

	//fmt.Printf("poolq.size %d \n", list.size)
	return nil
}

func (p *Backend) getPoolEntry() (*cacheEntry, error) {
	var e *list.ListHead
	if e = p.cache.poolq.dequeue(); e == nil {
		if err := p.allocPool(); err != nil {
			return nil, fmt.Errorf("%s(%s)", specs.EALLOC, err)
		}
		e = p.cache.poolq.dequeue()
	}

	return list_data_entry(e), nil
}

func (p *Backend) putPoolEntry(e *cacheEntry) error {
	e.e.hashkey[0] = C.char(0)
	p.cache.poolq.enqueue(&e.list_data)
	return nil
}

// called by rpc
func (p *Backend) createEntry(key string, item *specs.RrdItem) (*cacheEntry, error) {
	var (
		e     *cacheEntry
		ok    bool
		err   error
		cache *backendCache
	)

	cache = p.cache

	statInc(ST_CACHE_CREATE, 1)
	if e, ok = cache.hash[key]; ok {
		return e, specs.ErrExist
	}

	e, err = p.getPoolEntry()
	if err != nil {
		return nil, err
	}

	e.reset(p.timeNow(), item.Host, item.Name, item.Tags, item.Type, item.Step,
		item.Heartbeat, item.Min[0], item.Max[0])
	e.setHashkey(key)

	cache.Lock()
	cache.hash[key] = e
	cache.Unlock()

	cache.dataq.enqueue(&e.list_data)
	cache.idx0q.enqueue(&e.list_idx)

	if !p.Migrate.Disabled {
		_, err := os.Stat(e.filename(p))
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
	p.e.lastTs = C.int64_t(item.TimeStemp)
	idx := p.e.dataId & CACHE_SIZE_MASK
	p.e.time[idx] = C.int64_t(item.TimeStemp)
	p.e.value[idx] = C.double(item.Value)
	p.e.dataId += 1
}

func (p *cacheEntry) reset(createTs int64, _host, _name, _tags, _typ string,
	step, heartbeat int, min, max byte) error {
	p.Lock()
	defer p.Unlock()

	host := C.CString(_host)
	name := C.CString(_name)
	tags := C.CString(_tags)
	typ := C.CString(_typ)

	ret := int(C.cache_entry_reset(p.e, C.int64_t(createTs),
		host, name, tags, typ, C.int(step),
		C.int(heartbeat), C.char(min), C.char(max)))

	C.free(unsafe.Pointer(host))
	C.free(unsafe.Pointer(name))
	C.free(unsafe.Pointer(tags))
	C.free(unsafe.Pointer(typ))

	if ret != 0 {
		return specs.EINVAL
	}
	return nil
}

// fetch remote ds
func (p *cacheEntry) fetchCommit(backend *Backend) {
	done := make(chan error)

	node, err := backend.storageMigrateConsistent.Get(p.hashkey())
	if err != nil {
		return
	}

	backend.storageNetTaskCh[node] <- &netTask{
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

func (p *cacheEntry) createRrd(backend *Backend) error {
	done := make(chan error, 1)

	backend.ktoch(p.hashkey()) <- &ioTask{
		method: IO_TASK_M_RRD_ADD,
		args:   p,
		done:   done,
	}
	err := <-done

	p.e.commitTs = C.int64_t(backend.timeNow())

	return err
}

func (p *cacheEntry) commit(backend *Backend) error {
	done := make(chan error, 1)

	backend.ktoch(p.hashkey()) <- &ioTask{
		method: IO_TASK_M_RRD_UPDATE,
		args:   p,
		done:   done,
	}
	err := <-done

	p.e.commitTs = C.int64_t(backend.timeNow())

	return err
}

func (p *cacheEntry) filename(backend *Backend) string {
	return backend.ktofname(p.hashkey())
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

	//H := h & CACHE_SIZE_MASK
	L := l & CACHE_SIZE_MASK

	for i := uint32(0); i < size; i++ {
		idx := (L + i) & CACHE_SIZE_MASK
		ret[i] = &specs.RRDData{
			Ts: int64(p.e.time[idx]),
			V:  specs.JsonFloat(p.e.value[idx]),
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

func (p *cacheEntry) _dequeueAll() []*specs.RRDData {
	ret, over := p._getData(uint32(p.e.commitId), uint32(p.e.dataId))
	p.e.commitId = p.e.dataId
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

	rrds, _ := p._getData(0, uint32(p.e.dataId))

	for _, v := range rrds {
		ret = append(ret, &specs.RrdItem{
			Host:      p.host(),
			Name:      p.name(),
			Tags:      p.tags(),
			Value:     float64(v.V),
			TimeStemp: v.Ts,
			Type:      p.typ(),
			Step:      int(p.e.step),
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
	idx := uint32(p.e.dataId-1) & CACHE_SIZE_MASK
	return &specs.RrdItem{
		Host:      p.host(),
		Name:      p.name(),
		Tags:      p.tags(),
		Value:     float64(p.e.value[idx]),
		TimeStemp: int64(p.e.time[idx]),
		Type:      p.typ(),
		Step:      int(p.e.step),
	}
	return
}

func (p *cacheEntry) String() string {
	return fmt.Sprintf("key:%s host:%s name:%s "+
		"tags:%s type:%s step:%d\n",
		p.hashkey(), p.host(), p.name(),
		p.tags(), p.typ(), int(p.e.step))
}

func (p *Backend) getItems(key string) (ret []*specs.RrdItem) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func (p *Backend) getLastItem(key string) (ret *specs.RrdItem) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}

func (p *Backend) cacheInit() error {
	cache := &backendCache{
		hash:     make(map[string]*cacheEntry),
		shmaddrs: make([]uintptr, 0),
		shmids:   make([]int, 0),
		startkey: p.ShmKey,
		endkey:   p.ShmKey,
		cache_entry_nb: (p.ShmSize -
			C.sizeof_struct_cache_block) /
			C.sizeof_struct_cache_entry,
	}
	cache.dataq.init()
	cache.poolq.init()
	cache.idx0q.init()
	cache.idx1q.init()
	cache.idx2q.init()
	p.cache = cache
	return nil
}

func (p *Backend) cacheStart() error {
	//p.cache.start(p)

	n, err := p.importBlocks()

	if err != specs.EFMT {
		return err
	}
	glog.V(3).Infof(MODULE_NAME+"import %d \n", n)
	return nil
}

func (p *Backend) cacheStop() {
	p.cache.close()
}
