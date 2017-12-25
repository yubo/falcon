/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/yubo/falcon"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/service/config"
)

var (
	cacheApp  *Service
	cache     *CacheModule
	testEntry *cacheEntry
	item      *falcon.Item
	err       error
)

func newItem1(i int) *falcon.Item {
	return &falcon.Item{
		Endpoint:  []byte(fmt.Sprintf("host_%d", i)),
		Metric:    []byte(fmt.Sprintf("key_%d", i)),
		Value:     float64(i),
		Timestamp: int64(i) * DEBUG_STEP,
		Type:      falcon.ItemType_GAUGE,
		Tags:      []byte{},
	}
}

func test_cache_init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cacheApp = &Service{}
	cache = &CacheModule{}
	cacheApp.Conf = &config.Service{
		Name: "cacheApp",
	}
	cacheApp.Conf.Configer.Set(fconfig.APP_CONF_FILE, map[string]string{
		"hdisks": "/tmp/falcon",
	})
	cache.prestart(cacheApp)
}

func TestCache(t *testing.T) {
	//fmt.Println(runtime.Caller(0))
	test_cache_init()
	cache.prestart(cacheApp)
	item = newItem1(1)
	key := item.Csum()

	// create
	testEntry, err = cacheApp.createEntry(key, item)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("c.createEntry success\n")

	// get
	p := cacheApp.cache.get(item.Csum())
	if testEntry != p {
		t.Errorf("c.get(%s) error", item.Csum())
	}
	fmt.Printf("c.get success\n")

	if len(cacheApp.cache.data) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cacheApp.cache.data))
	}

	item = newItem1(2)
	cacheApp.createEntry(item.Csum(), item)
	if len(cacheApp.cache.data) != 2 {
		t.Errorf("c.hash size error size %d want 2", len(cacheApp.cache.data))
	}

	// unlink
	cacheApp.cache.unlink(newItem1(1).Csum())
	if len(cacheApp.cache.data) != 1 {
		t.Errorf("c.hash size error size %d want 1", len(cacheApp.cache.data))
	}
	fmt.Printf("c.unlink success\n")

	for k, _ := range cacheApp.cache.data {
		cacheApp.cache.unlink(k)
	}
	fmt.Printf("all c.unlink success\n")

}

func TestCacheQueue(t *testing.T) {
	cache.prestart(cacheApp)

	item = newItem1(0)
	testEntry, err = cacheApp.createEntry(item.Csum(), item)
	if err != nil {
		t.Errorf("%s:%s", "testCacheQueue", err)
	}

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())
	for i := 1; i < 2*CACHE_SIZE; i++ {
		testEntry.put(newItem1(i))
		if len(testEntry.getItems(CACHE_SIZE)) != CACHE_SIZE {
			t.Errorf("len(data) %d want %d", len(testEntry.getItems(CACHE_SIZE)), CACHE_SIZE)
		}
	}
	fmt.Printf("e.getItems() success\n")

}
