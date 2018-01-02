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
	cacheApp *Service
	cache    *ShardModule
	err      error
)

func newItem(i int) *falcon.Item {
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
	cacheApp = &Service{
		Conf: &config.Service{
			Name: "cacheApp",
		},
	}
	cacheApp.Conf.Configer.Set(fconfig.APP_CONF_FILE, map[string]string{
		"shardIds": "0",
	})

	cache = &ShardModule{}
	cache.prestart(cacheApp)
}

func TestCache(t *testing.T) {
	//fmt.Println(runtime.Caller(0))
	test_cache_init()
	cache.prestart(cacheApp)
	item1 := newItem(1)
	item2 := newItem(2)

	// create
	tie, err := cacheApp.shard.put(item1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("cache add success\n")

	// get
	bucket, err := cacheApp.shard.getBucket(0)
	if err != nil {
		t.Error(err)
	}

	ie, err := bucket.getItem(item1.Key())
	if err != nil || tie != ie {
		t.Errorf("bucket.getItem(%s) error", item1.Key())
	}
	fmt.Printf("bucket.getItem success\n")

	if len(bucket.itemMap) != 1 {
		t.Errorf(" size error size %d want 1", len(bucket.itemMap))
	}

	cacheApp.shard.put(item2)
	if len(bucket.itemMap) != 2 {
		t.Errorf(" size error size %d want 2", len(bucket.itemMap))
	}

	// unlink
	bucket.unlink(item1.Key())
	if len(bucket.itemMap) != 1 {
		t.Errorf(" size error size %d want 1", len(bucket.itemMap))
	}
	fmt.Printf("c.unlink success\n")

	for k, _ := range bucket.itemMap {
		bucket.unlink(k)
	}
	fmt.Printf("all c.unlink success\n")
}

func TestCacheQueue(t *testing.T) {
	cache.prestart(cacheApp)

	item := newItem(0)
	ie, err := cacheApp.shard.put(item)
	if err != nil {
		t.Errorf("%s:%s", "testCacheQueue", err)
	}

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())
	for i := 1; i < 2*CACHE_SIZE; i++ {
		ie.put(newItem(i))
		if len(ie.getItems(CACHE_SIZE)) != CACHE_SIZE {
			t.Errorf("len(data) %d want %d", len(ie.getItems(CACHE_SIZE)), CACHE_SIZE)
		}
	}
	fmt.Printf("e.getItems() success\n")

}
