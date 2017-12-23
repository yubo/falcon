/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/gotool/list"
)

// called by rpc
func (p *Backend) createEntry(key string, item *falcon.Item) (*cacheEntry, error) {
	var (
		e     *cacheEntry
		ok    bool
		cache *backendCache
	)

	cache = p.cache

	statsInc(ST_CACHE_CREATE, 1)
	if e, ok = cache.hash[key]; ok {
		return e, falcon.ErrExist
	}

	e = &cacheEntry{
		createTs: p.timeNow(),
		endpoint: item.Endpoint,
		metric:   item.Metric,
		tags:     item.Tags,
		typ:      item.Type,
		hashkey:  key,
		time:     make([]int64, CACHE_SIZE),
		value:    make([]float64, CACHE_SIZE),
		//heartbeat: item.Heartbeat,
		//min:       []byte(item.Min)[0],
		//max:       []byte(item.Max)[0],
	}

	cache.Lock()
	cache.hash[key] = e
	cache.Unlock()

	cache.dataq.enqueue(&e.list_data)
	cache.idx0q.enqueue(&e.list_idx)

	return e, nil
}

func (p *Backend) getItems(key string) (ret []*falcon.Item) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func (p *Backend) getLastItem(key string) (ret *falcon.Item) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}

func (p *Backend) handleItems(items []*falcon.Item) (total, errors int) {
	var (
		err error
		e   *cacheEntry
	)

	total = len(items)
	if total == 0 {
		return
	}

	glog.V(4).Infof(MODULE_NAME+"recv %d", total)
	statsInc(ST_RPC_SERV_RECV, 1)
	statsInc(ST_RPC_SERV_RECV_ITEM, total)

	for i := 0; i < total; i++ {
		if items[i] == nil {
			errors++
			continue
		}
		key := items[i].Csum()

		e = p.cache.get(key)
		if e == nil {
			e, err = p.createEntry(key, items[i])
			if err != nil {
				errors++
				continue
			}
		}

		/*
			if DATA_TIMESTAMP_REGULATE {
				items[i].Ts = items[i].Ts -
					items[i].Ts%int64(items[i].Step)
			}

			if items[i].Ts <= e.lastTs || items[i].Ts <= 0 {
				errors++
				continue
			}
		*/

		e.put(items[i])
	}
	return
}

func list_data_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_data))))
}

func list_idx_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_idx))))
}

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
/*
func (p *Backend) ktofname(key string) string {
	csum, _ := strconv.ParseUint(key[0:2], 16, 64)
	return fmt.Sprintf("%s/%s/%s.rrd",
		p.hdisk[int(csum)%len(p.hdisk)],
		key[0:2], key)
}

func (p *Backend) ktoch(key string) chan *ioTask {
	csum, _ := strconv.ParseUint(key[0:2], 16, 64)
	return p.storageIoTaskCh[int(csum)%len(p.hdisk)]
}
*/
