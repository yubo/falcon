/*
 * Copyright 2016 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/yubo/falcon"
)

var (
	cacheModule *cacheModule
	testEntry   *cacheEntry
	rrdItem     *falcon.RrdItem
	err         error
)

func newRrdItem1(i int) *falcon.RrdItem {
	return &falcon.RrdItem{
		Host:      fmt.Sprintf("host_%d", i),
		Name:      fmt.Sprintf("key_%d", i),
		Value:     float64(i),
		TimeStemp: int64(i) * DEBUG_STEP,
		Step:      60,
		Type:      falcon.GAUGE,
		Tags:      "",
		Heartbeat: 120,
		Min:       "U",
		Max:       "U",
	}
}

func test_cache_init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cacheApp = nil
	cacheApp.Conf = &falcon.ConfBackend{
		Name: "cacheApp",
	}
	cacheApp.Conf.Configer.Set(APP_CONF_FILE, map[string]string{
		"hdisks": "/tmp/falcon",
	})
	cacheApp.cacheInit()
}

func TestCache(t *testing.T) {
	//fmt.Println(runtime.Caller(0))
	test_cache_init()
	cacheApp.cacheReset()
	rrdItem = newRrdItem1(1)
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

	rrdItem = newRrdItem1(2)
	cacheApp.createEntry(rrdItem.Csum(), rrdItem)
	if len(cacheApp.cache.hash) != 2 {
		t.Errorf("c.hash size error size %d want 2", len(cacheApp.cache.hash))
	}

	// unlink
	cacheApp.cache.unlink(newRrdItem1(1).Csum())
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

	rrdItem = newRrdItem1(0)
	testEntry, err = cacheApp.createEntry(rrdItem.Csum(), rrdItem)
	if err != nil {
		t.Errorf("%s:%s", "testCacheQueue", err)
	}

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())

	for i := 1; i < 2*CACHE_SIZE; i++ {
		testEntry.put(newRrdItem1(i))
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
