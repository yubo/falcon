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
	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/service/config"
)

var (
	cacheApp *Service
	cache    *CacheModule
	err      error
)

func newDp(i int) *tsdb.DataPoint {
	return &tsdb.DataPoint{
		Key: &tsdb.Key{
			Key:     []byte(fmt.Sprintf("host_%d/metric_%d//GAUGE", i, i)),
			ShardId: 0,
		},
		Value: &tsdb.TimeValuePair{
			Timestamp: int64(i) * DEBUG_STEP,
			Value:     float64(i),
		},
	}
}

func test_cache_init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cacheApp = &Service{
		Conf: &config.Service{
			Name: "cacheApp",
		},
	}
	cacheApp.Conf.Configer.Set(falcon.APP_CONF_FILE, map[string]string{
		"shardIds": "0",
	})

	cache = &CacheModule{}
}

func testCache(t *testing.T) {
	//fmt.Println(runtime.Caller(0))
	test_cache_init()
	cache.prestart(cacheApp)
	cache.start(cacheApp)
	defer cache.stop(cacheApp)
	dp := newDp(1)

	// create
	e, err := cacheApp.cache.put(dp)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("cache add success\n")

	// get
	bucket, err := cacheApp.cache.getBucket(0)
	if err != nil {
		t.Error(err)
	}

	e_, err := bucket.getCacheEntry(string(dp.Key.Key))
	if err != nil || e != e_ {
		t.Errorf("bucket.getDpEntry(%s) error", string(dp.Key.Key))
	}
	fmt.Printf("bucket.getDpEntry success\n")

	if len(bucket.entries) != 1 {
		t.Errorf(" size error size %d want 1", len(bucket.entries))
	}

	dp2 := newDp(2)
	cacheApp.cache.put(dp2)
	if len(bucket.entries) != 2 {
		t.Errorf(" size error size %d want 2", len(bucket.entries))
	}

	// unlink
	bucket.delEntry(string(dp.Key.Key))
	if len(bucket.entries) != 1 {
		t.Errorf(" size error size %d want 1", len(bucket.entries))
	}
	fmt.Printf("c.unlink success\n")

	for k, _ := range bucket.entries {
		bucket.delEntry(k)
	}
	fmt.Printf("all c.unlink success\n")
}

func testCacheQueue(t *testing.T) {
	test_cache_init()
	cache.prestart(cacheApp)
	cache.start(cacheApp)
	defer cache.stop(cacheApp)

	dp := newDp(0)
	e, err := cacheApp.cache.put(dp)
	if err != nil {
		t.Errorf("%s:%s", "testCacheQueue", err.Error())
	}

	//fmt.Printf("cacheEtnry filename: %s\n", entry.filename())
	for i := 1; i < 2*CACHE_DATA_SIZE; i++ {
		e.put(newDp(i))
		if len(e.getDps(CACHE_DATA_SIZE).Values) != CACHE_DATA_SIZE {
			t.Errorf("len(data) %d want %d", len(e.getDps(CACHE_DATA_SIZE).Values), CACHE_DATA_SIZE)
		}
	}
	fmt.Printf("e.getDps() success\n")

}
