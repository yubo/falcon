/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	/* cache */
	cache storageCache
)

type storageCache struct {
	sync.RWMutex // hash lock
	dataq        cacheq
	idxq         cacheq
	hash         map[string]*cacheEntry
}

func (c *storageCache) get(key string) *cacheEntry {
	c.RLock()
	defer c.RUnlock()

	if e, ok := c.hash[key]; ok {
		return e
	}
	return nil

}

func (c *storageCache) unlink(key string) *cacheEntry {
	c.Lock()
	defer c.Unlock()
	e, ok := c.hash[key]
	if !ok {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	delete(c.hash, key)

	c.dataq.Lock()
	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
	e.prev = nil
	e.next = nil
	c.dataq.Unlock()

	c.idxq.Lock()
	if e.idxPrev != nil {
		e.idxPrev.idxNext = e.idxNext
	}
	if e.idxNext != nil {
		e.idxNext.idxPrev = e.idxPrev
	}
	e.idxPrev = nil
	e.idxNext = nil
	c.idxq.Unlock()

	return e

}

func (c *storageCache) enqueue_data(e *cacheEntry) {
	c.dataq.Lock()
	defer c.dataq.Unlock()

	e.Lock()
	defer e.Unlock()

	e.next = nil
	x := c.dataq.last
	if x == nil {
		e.prev = nil
		c.dataq.first = e
		c.dataq.last = e
		return
	}
	e.prev = x
	x.next = e
	c.dataq.last = e

}

func (c *storageCache) dequeue_data() *cacheEntry {
	c.dataq.Lock()
	defer c.dataq.Unlock()

	e := c.dataq.first
	if e == nil {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	x := e.next
	if x == nil {
		c.dataq.first = nil
		c.dataq.last = nil
	} else {
		x.prev = nil
		c.dataq.first = x
		e.next = nil
	}
	return e
}

func (c *storageCache) enqueue_idx(e *cacheEntry) {
	c.idxq.Lock()
	defer c.idxq.Unlock()

	e.Lock()
	defer e.Unlock()

	e.idxNext = nil
	x := c.idxq.last
	if x == nil {
		e.idxPrev = nil
		c.idxq.first = e
		c.idxq.last = e
		return
	}
	e.idxPrev = x
	x.idxNext = e
	c.idxq.last = e

}

func (c *storageCache) dequeue_idx() *cacheEntry {
	c.idxq.Lock()
	defer c.idxq.Unlock()

	e := c.idxq.first
	if e == nil {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	x := e.idxNext
	if x == nil {
		c.idxq.first = nil
		c.idxq.last = nil
	} else {
		x.idxPrev = nil
		c.idxq.first = x
		e.idxNext = nil
	}
	return e
}

// called by rpc
func (c *storageCache) put(key string, item *specs.RrdItem) (*cacheEntry, error) {

	c.Lock()
	if e, ok := c.hash[key]; ok {
		return e, specs.ErrExist
	}

	e := &cacheEntry{
		key:       key,
		createTs:  time.Now().Unix(),
		host:      item.Host,
		k:         item.K,
		tags:      item.Tags,
		typ:       item.Type,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       item.Min,
		max:       item.Max,
		cache:     []*specs.RRDData{},
		history:   []*specs.RRDData{},
	}
	c.hash[key] = e
	c.Unlock()

	c.enqueue_data(e)
	c.enqueue_idx(e)

	if rpcConfig.Migrate.Enable {
		_, err := os.Stat(e.filename(rpcConfig.RrdStorage))
		if os.IsNotExist(err) {
			e.flag = RRD_F_MISS
		}
	}
	return e, nil
}

// called by rpc
func (c *cacheEntry) put(item *specs.RrdItem) {
	c.Lock()
	defer c.Unlock()
	c.ts = item.Ts
	c.cache = append(c.cache, &specs.RRDData{
		Ts: item.Ts,
		V:  specs.JsonFloat(item.V),
	})
}

// fetch remote ds
func (e *cacheEntry) fetch() {
	done := make(chan error)

	node, err := rrdMigrateConsistent.Get(e.key)
	if err != nil {
		return
	}

	rrdNetTaskCh[node] <- &netTask{
		Method: NET_TASK_M_FETCH,
		e:      e,
		Done:   done,
	}

	// net_task slow, shouldn't block commitCache()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			glog.Warning("get %s from remote err[%s]\n", e.key, err)
			return
		}
		statInc(ST_NET_TASK_CNT, 1)
		//todo: flushfile after getfile? not yet
	}()
}

func (e *cacheEntry) commit() error {
	done := make(chan error, 1)

	rrdIoTaskCh <- &ioTask{
		method: IO_TASK_M_RRD_UPDATE,
		args:   e,
		done:   done,
	}
	err := <-done

	statInc(ST_DISK_TASK_CNT, 1)
	e.commitTs = time.Now().Unix()

	return err
}

func (e *cacheEntry) filename(dir string) string {
	return fmt.Sprintf("%s/%s/%s.rrd", dir, e.key[0:2], e.key)
}

func (e *cacheEntry) _dequeueAll() []*specs.RRDData {
	ret := e.cache
	e.cache = []*specs.RRDData{}
	return ret
}

func (e *cacheEntry) dequeueAll() []*specs.RRDData {
	e.Lock()
	defer e.Unlock()

	return e._dequeueAll()
}

func (e *cacheEntry) _getItems() (ret []*specs.RrdItem) {

	for _, v := range e.cache {
		ret = append(ret, &specs.RrdItem{
			Host: e.host,
			K:    e.k,
			Tags: e.tags,
			V:    float64(v.V),
			Ts:   v.Ts,
			Type: e.typ,
			Step: e.step,
		})
	}

	return ret
}

func (e *cacheEntry) getItems() (ret []*specs.RrdItem) {
	e.Lock()
	defer e.Unlock()

	return e._getItems()
}

func (e *cacheEntry) getItemsAll() (ret []*specs.RrdItem) {
	e.RLock()
	defer e.RUnlock()

	for _, v := range e.history {
		ret = append(ret, &specs.RrdItem{
			Host: e.host,
			K:    e.k,
			Tags: e.tags,
			V:    float64(v.V),
			Ts:   v.Ts,
			Type: e.typ,
			Step: e.step,
		})
	}

	for _, v := range e.cache {
		ret = append(ret, &specs.RrdItem{
			Host: e.host,
			K:    e.k,
			Tags: e.tags,
			V:    float64(v.V),
			Ts:   v.Ts,
			Type: e.typ,
			Step: e.step,
		})
	}
	return
}

/* the last item(dequeue) */
func (e *cacheEntry) getItem() (ret *specs.RrdItem) {
	e.RLock()
	defer e.RUnlock()

	if len(e.cache) > 0 {
		v := e.cache[len(e.cache)-1]
		return &specs.RrdItem{
			Host: e.host,
			K:    e.k,
			Tags: e.tags,
			V:    float64(v.V),
			Ts:   v.Ts,
			Type: e.typ,
			Step: e.step,
		}
	}

	if len(e.history) > 0 {
		v := e.history[len(e.history)-1]
		return &specs.RrdItem{
			Host: e.host,
			K:    e.k,
			Tags: e.tags,
			V:    float64(v.V),
			Ts:   v.Ts,
			Type: e.typ,
			Step: e.step,
		}
	}
	return
}

func getAllItems(key string) (ret []*specs.RrdItem) {
	e := cache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func getLastItem(key string) (ret *specs.RrdItem) {
	e := cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}
