/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yubo/falcon/specs"
)

type cache_t struct {
	sync.RWMutex // hash lock
	dataq        cacheq
	idxq         cacheq
	hash         map[string]*cacheEntry
}

func (c *cache_t) get(key string) *cacheEntry {
	c.RLock()
	defer c.RUnlock()

	if e, ok := c.hash[key]; ok {
		return e
	}
	return nil

}

func (c *cache_t) unlink(key string) *cacheEntry {
	c.Lock()
	defer c.Unlock()
	e, ok := c.hash[key]
	if !ok {
		return nil
	}
	delete(c.hash, key)

	c.dataq.Lock()
	if e.data_prev != nil {
		e.data_prev.data_next = e.data_next
	}
	if e.data_next != nil {
		e.data_next.data_prev = e.data_prev
	}
	e.data_prev = nil
	e.data_next = nil
	c.dataq.Unlock()

	c.idxq.Lock()
	if e.idx_prev != nil {
		e.idx_prev.idx_next = e.idx_next
	}
	if e.idx_next != nil {
		e.idx_next.idx_prev = e.idx_prev
	}
	e.idx_prev = nil
	e.idx_next = nil
	c.idxq.Unlock()

	return e

}

func (c *cache_t) enqueue_data(e *cacheEntry) {
	c.dataq.Lock()
	defer c.dataq.Unlock()

	e.data_next = nil
	x := c.dataq.last
	if x == nil {
		e.data_prev = nil
		c.dataq.first = e
		c.dataq.last = e
		return
	}
	e.data_prev = x
	x.data_next = e
	c.dataq.last = e

}

func (c *cache_t) dequeue_data() *cacheEntry {
	c.dataq.Lock()
	defer c.dataq.Unlock()

	e := c.dataq.first
	if e == nil {
		return nil
	}

	x := e.data_next
	if x == nil {
		c.dataq.first = nil
		c.dataq.last = nil
	} else {
		x.data_prev = nil
		c.dataq.first = x
		e.data_next = nil
	}
	return e
}

func (c *cache_t) enqueue_idx(e *cacheEntry) {
	c.idxq.Lock()
	defer c.idxq.Unlock()

	e.idx_next = nil
	x := c.idxq.last
	if x == nil {
		e.idx_prev = nil
		c.idxq.first = e
		c.idxq.last = e
		return
	}
	e.idx_prev = x
	x.idx_next = e
	c.idxq.last = e

}

func (c *cache_t) dequeue_idx() *cacheEntry {
	c.idxq.Lock()
	defer c.idxq.Unlock()

	e := c.idxq.first
	if e == nil {
		return nil
	}

	x := e.idx_next
	if x == nil {
		c.idxq.first = nil
		c.idxq.last = nil
	} else {
		x.idx_prev = nil
		c.idxq.first = x
		e.idx_next = nil
	}
	return e
}

func (c *cache_t) put(key string, item *specs.GraphItem) (*cacheEntry, error) {

	c.Lock()
	if e, ok := c.hash[key]; ok {
		return e, ErrExist
	}

	e := &cacheEntry{
		key:       key,
		createTs:  time.Now().Unix(),
		endpoint:  item.Endpoint,
		metric:    item.Metric,
		tags:      make(map[string]string, len(item.Tags)),
		dsType:    item.DsType,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       item.Min,
		max:       item.Max,
		cache:     []*specs.RRDData{},
		history:   []*specs.RRDData{},
	}
	c.hash[key] = e
	c.Unlock()

	for k, v := range item.Tags {
		e.tags[k] = v
	}

	c.enqueue_data(e)
	c.enqueue_idx(e)

	if config().Migrate.Enable {
		_, err := os.Stat(e.filename(config().RrdStorage))
		if os.IsNotExist(err) {
			e.flag = GRAPH_F_MISS
		}
	}
	return e, nil
}

func (c *cacheEntry) put(item *specs.GraphItem) {
	c.Lock()
	defer c.Unlock()
	c.putTs = item.Timestamp
	c.cache = append(c.cache, &specs.RRDData{
		Timestamp: item.Timestamp,
		Value:     specs.JsonFloat(item.Value),
	})
}

// fetch remote ds
func (e *cacheEntry) fetch() {
	done := make(chan error)

	node, err := Consistent.Get(e.key)
	if err != nil {
		return
	}

	Net_task_ch[node] <- &Net_task_t{
		Method: NET_TASK_M_FETCH,
		e:      e,
		Done:   done,
	}

	// net_task slow, shouldn't block syncDisk() or FlushAll()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			log.Printf("get %s from remote err[%s]\n", e.key, err)
			return
		}
		stat_inc(ST_NET_COUNTER, 1)
		//todo: flushfile after getfile? not yet
	}()
}

func (e *cacheEntry) commit() error {
	done := make(chan error, 1)

	io_task_chan <- &io_task_t{
		method: IO_TASK_M_COMMIT,
		args:   e,
		done:   done,
	}
	err := <-done

	stat_inc(ST_DISK_COUNTER, 1)
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

func (e *cacheEntry) _getItems() (ret []*specs.GraphItem) {
	tags := make(map[string]string, len(e.tags))
	for k, v := range e.tags {
		tags[k] = v
	}

	for _, v := range e.cache {
		ret = append(ret, &specs.GraphItem{
			Endpoint:  e.endpoint,
			Metric:    e.metric,
			Tags:      tags,
			Value:     float64(v.Value),
			Timestamp: v.Timestamp,
			DsType:    e.dsType,
			Step:      e.step,
		})
	}

	return ret
}

func (e *cacheEntry) getItems() (ret []*specs.GraphItem) {
	e.Lock()
	defer e.Unlock()

	return e._getItems()
}

func (e *cacheEntry) getItemsAll() (ret []*specs.GraphItem) {
	e.RLock()
	defer e.RUnlock()

	tags := make(map[string]string, len(e.tags))
	for k, v := range e.tags {
		tags[k] = v
	}

	for _, v := range e.history {
		ret = append(ret, &specs.GraphItem{
			Endpoint:  e.endpoint,
			Metric:    e.metric,
			Tags:      tags,
			Value:     float64(v.Value),
			Timestamp: v.Timestamp,
			DsType:    e.dsType,
			Step:      e.step,
		})
	}

	for _, v := range e.cache {
		ret = append(ret, &specs.GraphItem{
			Endpoint:  e.endpoint,
			Metric:    e.metric,
			Tags:      tags,
			Value:     float64(v.Value),
			Timestamp: v.Timestamp,
			DsType:    e.dsType,
			Step:      e.step,
		})
	}
	return
}

/* the last item(dequeue) */
func (e *cacheEntry) getItem() (ret *specs.GraphItem) {
	e.RLock()
	defer e.RUnlock()

	tags := make(map[string]string, len(e.tags))
	for k, v := range e.tags {
		tags[k] = v
	}

	if len(e.cache) > 0 {
		v := e.cache[len(e.cache)-1]
		return &specs.GraphItem{
			Endpoint:  e.endpoint,
			Metric:    e.metric,
			Tags:      tags,
			Value:     float64(v.Value),
			Timestamp: v.Timestamp,
			DsType:    e.dsType,
			Step:      e.step,
		}
	}

	if len(e.history) > 0 {
		v := e.history[len(e.history)-1]
		return &specs.GraphItem{
			Endpoint:  e.endpoint,
			Metric:    e.metric,
			Tags:      tags,
			Value:     float64(v.Value),
			Timestamp: v.Timestamp,
			DsType:    e.dsType,
			Step:      e.step,
		}
	}
	return
}

func getAllItems(key string) (ret []*specs.GraphItem) {
	e := cache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func getLastItem(key string) (ret *specs.GraphItem) {
	e := cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}
