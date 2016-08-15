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

var (
	/* cache */
	appCache    backendCache
	cacheEvent  chan specs.ProcEvent
	cacheConfig BackendOpts
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
	defer C.free(cs)
	return int(C.set_host(p.e, cs))
}

func (p *cacheEntry) setName(s string) int {
	cs := C.CString(s)
	defer C.free(cs)
	return int(C.set_name(p.e, cs))
}

func (p *cacheEntry) setTags(s string) int {
	cs := C.CString(s)
	defer C.free(cs)
	return int(C.set_tags(p.e, cs))
}

func (p *cacheEntry) setTyp(s string) int {
	cs := C.CString(s)
	defer C.free(cs)
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

/*
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
*/

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
	magic          uint32
	block_size     int //segment size(bytes)
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

func (p *backendCache) importBlocks(magic uint32) (int, error) {
	var (
		shm_id int
		size   int
		addr   uintptr
		block  *C.struct_cache_block
		err    error
	)

	for {
		shm_id, size, err = Shmget(p.endkey, p.block_size, C.IPC_CREAT|0600)
		if err != nil {
			return p.endkey - p.startkey, err
		}

		if size != p.block_size {
			glog.Errorf("segment size error, remove "+
				"segment from 0x%x and retry", p.endkey)
			cleanBlocks(p.endkey)
			continue
		}

		addr, err = Shmat(shm_id, 0, 0)
		if err != nil {
			return p.endkey - p.startkey, err
		}

		block = (*C.struct_cache_block)(unsafe.Pointer(addr))

		for uint32(block.magic) != p.magic {
			glog.Infof("set magic head at 0x%x",
				shm_id)
			block.magic = C.uint32_t(p.magic)
			return p.endkey - p.startkey, specs.EFMT
		}

		p.shmaddrs = append(p.shmaddrs, addr)
		p.shmids = append(p.shmids, shm_id)
		list := &p.poolq

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
				p.hash[key] = e
				p.dataq.enqueue(&e.list_data)
				p.idx0q.enqueue(&e.list_idx)
			}
		}
		p.endkey++
	}
	return p.endkey - p.startkey, nil
}

func cleanBlocks(start int) int {
	var (
		shm_id int
		err    error
	)

	for i := start; ; i++ {
		shm_id, _, err = Shmget(i, 0, 0600)
		if err != nil {
			return i - start
		}
		Shmrm(shm_id)
	}
}

func (p *backendCache) init(config BackendOpts) error {

	p.hash = make(map[string]*cacheEntry)
	p.shmaddrs = make([]uintptr, 0)
	p.shmids = make([]int, 0)
	p.dataq.init()
	p.poolq.init()
	p.idx0q.init()
	p.idx1q.init()
	p.idx2q.init()
	p.magic = config.Shm.Magic
	p.startkey = config.Shm.Key
	p.endkey = config.Shm.Key
	p.block_size = config.Shm.Size
	p.cache_entry_nb = (config.Shm.Size - C.sizeof_struct_cache_block) /
		C.sizeof_struct_cache_entry

	n, err := p.importBlocks(config.Shm.Magic)

	if err != specs.EFMT {
		return err
	}
	glog.V(3).Infof("import %d \n", n)
	return nil
}

func (p *backendCache) reset(config BackendOpts) error {
	p.hash = make(map[string]*cacheEntry)
	p.shmaddrs = make([]uintptr, 0)
	p.shmids = make([]int, 0)
	p.dataq.init()
	p.poolq.init()
	p.idx0q.init()
	p.idx1q.init()
	p.idx2q.init()
	p.magic = config.Shm.Magic
	p.startkey = config.Shm.Key
	p.endkey = config.Shm.Key
	p.block_size = config.Shm.Size
	p.cache_entry_nb = (config.Shm.Size - C.sizeof_struct_cache_block) /
		C.sizeof_struct_cache_entry
	cleanBlocks(p.startkey)
	return nil
}

func (p *backendCache) close() {
	for _, addr := range p.shmaddrs {
		Shmdt(addr)
	}
	glog.V(3).Infof("shmdt %d", len(p.shmaddrs))
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

func (p *backendCache) allocPool() (err error) {
	var (
		shm_id int
		size   int
		addr   uintptr
		block  *C.struct_cache_block
	)

	shm_id, size, err = Shmget(p.endkey, p.block_size, C.IPC_CREAT|0700)
	if err != nil {
		return err
	}

	if size != p.block_size {
		glog.Errorf("alloc shm size %d want %d, remove and try again",
			size, p.block_size)
		cleanBlocks(p.endkey)
		shm_id, _, err = Shmget(p.endkey, p.block_size, C.IPC_CREAT|0700)
		if err != nil {
			return err
		}
	}

	addr, err = Shmat(shm_id, 0, 0)
	if err != nil {
		return err
	}
	block = (*C.struct_cache_block)(unsafe.Pointer(addr))

	if uint32(block.magic) == p.magic {
		Shmdt(addr)
		return specs.ErrExist
	}

	p.Lock()
	p.shmaddrs = append(p.shmaddrs, addr)
	p.shmids = append(p.shmids, shm_id)
	p.endkey++
	glog.V(3).Infof("alloc shm segment, endkey %d len(shmaddr) %d",
		p.endkey, len(p.shmaddrs))
	p.Unlock()

	block.magic = C.uint32_t(p.magic)
	block.block_size = C.int(p.block_size)
	block.cache_entry_nb = C.int(p.cache_entry_nb)

	list := &p.poolq
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

func (p *backendCache) getPoolEntry() (*cacheEntry, error) {
	var e *list.ListHead
	if e = p.poolq.dequeue(); e == nil {
		if err := p.allocPool(); err != nil {
			return nil, fmt.Errorf("%s(%s)", specs.ENOSPC, err)
		}
		e = p.poolq.dequeue()
	}

	return list_data_entry(e), nil
}

func (p *backendCache) putPoolEntry(e *cacheEntry) error {
	e.e.hashkey[0] = C.char(0)
	p.poolq.enqueue(&e.list_data)
	return nil
}

// called by rpc
func (p *backendCache) createEntry(key string, item *specs.RrdItem) (*cacheEntry, error) {
	var (
		e   *cacheEntry
		ok  bool
		err error
	)

	statInc(ST_CACHE_CREATE, 1)
	if e, ok = p.hash[key]; ok {
		return e, specs.ErrExist
	}

	e, err = p.getPoolEntry()
	if err != nil {
		return nil, err
	}

	e.reset(timeNow(), item.Host, item.Name, item.Tags, item.Type, item.Step,
		item.Heartbeat, item.Min[0], item.Max[0])
	e.setHashkey(key)

	p.Lock()
	p.hash[key] = e
	p.Unlock()

	p.dataq.enqueue(&e.list_data)
	p.idx0q.enqueue(&e.list_idx)

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
func (p *cacheEntry) fetchCommit() {
	done := make(chan error)

	node, err := storageMigrateConsistent.Get(p.hashkey())
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

func (p *cacheEntry) createRrd() error {
	done := make(chan error, 1)

	ktoch(p.hashkey()) <- &ioTask{
		method: IO_TASK_M_RRD_ADD,
		args:   p,
		done:   done,
	}
	err := <-done

	p.e.commitTs = C.int64_t(timeNow())

	return err
}

func (p *cacheEntry) commit() error {
	done := make(chan error, 1)

	ktoch(p.hashkey()) <- &ioTask{
		method: IO_TASK_M_RRD_UPDATE,
		args:   p,
		done:   done,
	}
	err := <-done

	p.e.commitTs = C.int64_t(timeNow())

	return err
}

func (p *cacheEntry) filename() string {
	return ktofname(p.hashkey())
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

func cacheStart(config BackendOpts, p *specs.Process) {

	cacheConfig = config

	p.RegisterEvent("cache", cacheEvent)
	appCache.init(config)

	go func() {
		select {
		case event := <-cacheEvent:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				appCache.close()
				event.Done <- nil
				return
			}
		}
	}()

}
