/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yubo/falcon/specs"
)

var (
	cache   backendCache
	entry   *cacheEntry
	rrdItem *specs.RrdItem
	err     error
)

func newRrdItem(i int) *specs.RrdItem {
	return &specs.RrdItem{
		Host:      fmt.Sprintf("host_%d", i),
		Name:      fmt.Sprintf("key_%d", i),
		Value:     float64(i),
		TimeStemp: int64(i) * DEBUG_STEP,
		Step:      60,
		Type:      specs.GAUGE,
		Tags:      "",
		Heartbeat: 120,
		Min:       "U",
		Max:       "U",
	}

}

func init() {
	atomic.StoreInt64(&appTs, time.Now().Unix())
	cache.init()
}

func TestBackendCache(t *testing.T) {
	fmt.Println(runtime.Caller(0))
	rrdItem = newRrdItem(1)
	key := rrdItem.Csum()

	// create
	entry, err = cache.createEntry(key, rrdItem)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("c.createEntry success\n")

	// get
	p := cache.get(rrdItem.Csum())
	if entry != p {
		t.Errorf("c.get(%s) error", rrdItem.Csum())
	}
	fmt.Printf("c.get success\n")

	if len(cache.hash) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cache.hash))
	}

	rrdItem = newRrdItem(2)
	cache.createEntry(rrdItem.Csum(), rrdItem)
	if len(cache.hash) != 2 {
		t.Errorf("c.hash size error size %d want 2", len(cache.hash))
	}

	// unlink
	cache.unlink(newRrdItem(1).Csum())
	if len(cache.hash) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cache.hash))
	}
	fmt.Printf("c.unlink success\n")

	for k, _ := range cache.hash {
		cache.unlink(k)
	}
	fmt.Printf("all c.unlink success\n")

}

func TestCacheQueue(t *testing.T) {
	rrdItem = newRrdItem(0)
	entry, _ = cache.createEntry(rrdItem.Csum(), rrdItem)

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())

	for i := 1; i < 2*CACHE_SIZE; i++ {
		entry.put(newRrdItem(i))
		if i < CACHE_SIZE {
			if i != len(entry.getItems()) {
				t.Errorf("len(data) %d want %d", len(entry.getItems()), i)
			}
		} else {
			if len(entry.getItems()) != CACHE_SIZE {
				t.Errorf("len(data) %d want %d", len(entry.getItems()), CACHE_SIZE)
			}
		}
	}
	fmt.Printf("e.getItems() success\n")

	entry.dequeueAll()
	if entry.commitId != entry.dataId {
		t.Errorf("len(cache) %d want %d", entry.dataId-entry.commitId, 0)
	}
}
