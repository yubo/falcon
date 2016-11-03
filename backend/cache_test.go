/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/yubo/falcon/specs"
)

var (
	cacheApp  *Backend
	testEntry *cacheEntry
	rrdItem   *specs.RrdItem
	err       error
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	cacheApp = &Backend{
		ShmMagic: 0x80386,
		ShmKey:   0x6020,
		ShmSize:  4096,
		Storage: Storage{
			Type:   "rrd",
			Hdisks: []string{"/tmp/falcon"},
		},

		ts: time.Now().Unix(),
	}
	cacheApp.cacheInit()
	//cacheApp.cacheStart()
}

func TestCache(t *testing.T) {
	//fmt.Println(runtime.Caller(0))
	cacheApp.cacheReset()
	rrdItem = newRrdItem(1)
	key := rrdItem.Csum()

	// create
	testEntry, err = cacheApp.createEntry(key, rrdItem)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("c.createEntry success\n")

	// get
	p := cacheApp.cache.get(rrdItem.Csum())
	if testEntry != p {
		t.Errorf("c.get(%s) error", rrdItem.Csum())
	}
	fmt.Printf("c.get success\n")

	if len(cacheApp.cache.hash) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cacheApp.cache.hash))
	}

	rrdItem = newRrdItem(2)
	cacheApp.createEntry(rrdItem.Csum(), rrdItem)
	if len(cacheApp.cache.hash) != 2 {
		t.Errorf("c.hash size error size %d want 2", len(cacheApp.cache.hash))
	}

	// unlink
	cacheApp.cache.unlink(newRrdItem(1).Csum())
	if len(cacheApp.cache.hash) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cacheApp.cache.hash))
	}
	fmt.Printf("c.unlink success\n")

	for k, _ := range cacheApp.cache.hash {
		cacheApp.cache.unlink(k)
	}
	fmt.Printf("all c.unlink success\n")

}

func TestCacheQueue(t *testing.T) {
	cacheApp.cacheReset()

	rrdItem = newRrdItem(0)
	testEntry, err = cacheApp.createEntry(rrdItem.Csum(), rrdItem)
	if err != nil {
		t.Errorf("%s:%s", "testCacheQueue", err)
	}

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())

	for i := 1; i < 2*CACHE_SIZE; i++ {
		testEntry.put(newRrdItem(i))
		if i < CACHE_SIZE {
			if i != len(testEntry.getItems()) {
				t.Errorf("len(data) %d want %d", len(testEntry.getItems()), i)
			}
		} else {
			if len(testEntry.getItems()) != CACHE_SIZE {
				t.Errorf("len(data) %d want %d", len(testEntry.getItems()), CACHE_SIZE)
			}
		}
	}
	fmt.Printf("e.getItems() success\n")

	testEntry.dequeueAll()
	if testEntry.e.commitId != testEntry.e.dataId {
		t.Errorf("len(cache) %d want %d", int(testEntry.e.dataId-testEntry.e.commitId), 0)
	}
}

func TestCacheShm(t *testing.T) {
	var (
		block_nb int = 3
		entry_nb int
	)

	cacheApp.ShmSize = 268435456
	cacheApp.cacheReset()
	entry_nb = block_nb * cacheApp.cache.cache_entry_nb
	fmt.Printf("entry_nb %d block_nb %d block_entry_nb %d\n",
		entry_nb, block_nb, cacheApp.cache.cache_entry_nb)

	for i := 0; i < entry_nb; i++ {
		rrdItem = newRrdItem(i)
		if testEntry, err = cacheApp.createEntry(rrdItem.Csum(), rrdItem); err != nil {
			t.Errorf("%s:%s", "testCacheShm", err)
		}
	}

	cacheApp.cacheStop()
	cleanBlocks(cacheApp.cache.startkey)

	if err = cacheApp.cacheInit(); err != nil {
		t.Errorf("%s:%s", "testCacheShm", err)
	}
	if err = cacheApp.cacheStart(); err != nil {
		t.Errorf("%s:%s", "testCacheShm", err)
	}
}
