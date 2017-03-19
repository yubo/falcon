/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/gotool/list"
)

// called by rpc
func (p *Backend) createEntry(key string, item *falcon.RrdItem) (*cacheEntry, error) {
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
		createTs:  p.timeNow(),
		host:      item.Host,
		name:      item.Name,
		tags:      item.Tags,
		typ:       item.Type,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       []byte(item.Min)[0],
		max:       []byte(item.Max)[0],
		hashkey:   key,
		time:      make([]int64, CACHE_SIZE),
		value:     make([]float64, CACHE_SIZE),
	}

	cache.Lock()
	cache.hash[key] = e
	cache.Unlock()

	cache.dataq.enqueue(&e.list_data)
	cache.idx0q.enqueue(&e.list_idx)

	if !p.Conf.Migrate.Disabled {
		_, err := os.Stat(e.filename(p))
		if os.IsNotExist(err) {
			e.flag = RRD_F_MISS
		}
	}
	return e, nil
}

func (p *Backend) getItems(key string) (ret []*falcon.RrdItem) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItems()
}

func (p *Backend) getLastItem(key string) (ret *falcon.RrdItem) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}

func (p *Backend) handleItems(items []*falcon.RrdItem) {
	var (
		err error
		e   *cacheEntry
	)

	if items == nil {
		return
	}

	n := len(items)
	if n == 0 {
		return
	}

	glog.V(4).Infof(MODULE_NAME+"recv %d", n)
	statsInc(ST_RPC_SERV_RECV, 1)
	statsInc(ST_RPC_SERV_RECV_ITEM, n)

	for i := 0; i < n; i++ {
		if items[i] == nil {
			continue
		}
		key := items[i].Csum()

		e = p.cache.get(key)
		if e == nil {
			e, err = p.createEntry(key, items[i])
			if err != nil {
				continue
			}
		}

		if DATA_TIMESTAMP_REGULATE {
			items[i].TimeStemp = items[i].TimeStemp -
				items[i].TimeStemp%int64(items[i].Step)
		}

		if items[i].TimeStemp <= e.lastTs || items[i].TimeStemp <= 0 {
			continue
		}

		e.put(items[i])
	}
}

// 非法值: ts=0,value无意义
func (p *Backend) getLast(csum string) *falcon.RRDData {
	nan := &falcon.RRDData{Ts: 0, V: falcon.JsonFloat(0.0)}

	e := p.cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	typ := e.typ
	if typ == falcon.GAUGE {
		if e.dataId == 0 {
			return nan
		}

		idx := uint32(e.dataId-1) & CACHE_SIZE_MASK
		return &falcon.RRDData{
			Ts: int64(e.time[idx]),
			V:  falcon.JsonFloat(e.value[idx]),
		}
	}

	if typ == falcon.COUNTER || typ == falcon.DERIVE {

		if e.dataId < 2 {
			return nan
		}

		data, _ := e._getData(uint32(e.dataId)-2, uint32(e.dataId))

		delta_ts := data[0].Ts - data[1].Ts
		delta_v := data[0].V - data[1].V
		if delta_ts != int64(e.step) || delta_ts <= 0 {
			return nan
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return &falcon.RRDData{Ts: data[0].Ts,
			V: falcon.JsonFloat(float64(delta_v) / float64(delta_ts))}
	}
	return nan
}

func (p *Backend) getLastRaw(csum string) *falcon.RRDData {
	nan := &falcon.RRDData{Ts: 0, V: falcon.JsonFloat(0.0)}
	e := p.cache.get(csum)
	if e == nil {
		return nan
	}

	e.RLock()
	defer e.RUnlock()

	if e.typ == falcon.GAUGE {
		if e.dataId == 0 {
			return nan
		}
		idx := uint32(e.dataId-1) & CACHE_SIZE_MASK
		return &falcon.RRDData{
			Ts: int64(e.time[idx]),
			V:  falcon.JsonFloat(e.value[idx]),
		}
	}
	return nan
}

func list_data_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_data))))
}

func list_idx_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_idx))))
}

func (p *Backend) taskFileRead(key string) ([]byte, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_FILE_READ,
		args:   &falcon.File{Filename: p.ktofname(key)},
		done:   done,
	}

	p.ktoch(key) <- task
	err := <-done
	return task.args.(*falcon.File).Data, err
}

// get local data
func (p *Backend) taskRrdFetch(key string, cf string, start, end int64,
	step int) ([]*falcon.RRDData, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_RRD_FETCH,
		args: &rrdCheckout{
			filename: p.ktofname(key),
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	p.ktoch(key) <- task
	err := <-done
	return task.args.(*rrdCheckout).data, err
}

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
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
