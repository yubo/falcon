/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	fconfig "github.com/yubo/falcon/config"
)

const (
	testDir       = "/tmp"
	benchSize     = 20
	MAX_HD_NUMBER = 1
)

var (
	storageApp *Backend
	storage    *StorageModule
	testDirs   []string
	wg         sync.WaitGroup
	lock       sync.RWMutex
	now        int64
	es         [MAX_HD_NUMBER][benchSize]*cacheEntry
)

func newRrdItem2(i int, ts int64) *falcon.RrdItem {
	return &falcon.RrdItem{
		Host:      fmt.Sprintf("host_%d", i),
		Name:      fmt.Sprintf("key_%d", i),
		Value:     float64(i),
		TimeStemp: ts,
		Step:      60,
		Type:      falcon.GAUGE,
		Tags:      "",
		Heartbeat: 120,
		Min:       "U",
		Max:       "U",
	}
}

func testStoragetInit() (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()

	storageApp = &Backend{}
	storage = &StorageModule{}
	cache = &CacheModule{}

	storageApp.Conf.Configer.Set(fconfig.APP_CONF_FILE, map[string]string{
		"hdisks": "/tmp/falcon",
	})

	testDirs = make([]string, MAX_HD_NUMBER)

	for i := 0; i < MAX_HD_NUMBER; i++ {
		testDirs[i] = fmt.Sprintf("%s/hdd%d/test", testDir, i+1)
		os.RemoveAll(testDirs[i])
		os.MkdirAll(testDirs[i], 0755)
	}

	err = storageCheckHds(testDirs)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rrdtool.Start error, bad data dir %v\n", err)
	}

	storageApp.storageIoTaskCh = make([]chan *ioTask, MAX_HD_NUMBER)
	for i := 0; i < MAX_HD_NUMBER; i++ {
		storageApp.storageIoTaskCh[i] = make(chan *ioTask, 320)
		go storage.ioWorker(storageApp.storageIoTaskCh[i])
	}

	cache.prestart(storageApp)
	storage.prestart(storageApp)
	cache.start(storageApp)
	storage.start(storageApp)

	now = time.Now().Unix()
	storageApp.ts = now

	for i := 0; i < MAX_HD_NUMBER; i++ {
		item := newRrdItem2(i, now)
		for j := 0; j < benchSize; j++ {
			es[i][j], err = storageApp.rrdToEntry(item)
			if err != nil {
				glog.Infof("benchmarkAdd %s\n", err)
				return err
			}
			es[i][j].name = fmt.Sprintf("key_%d_%d", i, j)
			es[i][j].hashkey = es[i][j].csum()
		}
	}
	return nil
}

func testStorageDone() (err error) {
	storage.stop(storageApp)
	cache.stop(storageApp)
	return nil
}

func (p *Backend) rrdToEntry(item *falcon.RrdItem) (*cacheEntry, error) {

	return &cacheEntry{
		createTs:  p.timeNow(),
		host:      item.Host,
		name:      item.Name,
		tags:      item.Tags,
		typ:       item.Type,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       []byte(item.Min)[0],
		max:       []byte(item.Max)[0],
	}, nil
}

func testAdd(n int, t *testing.T) {
	var ts int64

	hds := &storageApp.hdisk
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", testDirs[i], n)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	m := benchSize
	N := m * n

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := time.Now().UnixNano()

			for j := 0; j < m; j++ {
				if err := es[i][j].createRrd(storageApp); err != nil {
					glog.V(4).Infof("%s", err)
				}
			}
			stop := time.Now().UnixNano()
			atomic.AddInt64(&ts, stop-start)
		}(i)
	}
	wg.Wait()
	glog.Infof("add %d %d %d ns/op", n, N, ts/int64(m*n*n))
}

func testUpdate(n int, t *testing.T) {
	var ts int64

	hds := &storageApp.hdisk
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", testDirs[i], n)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	m := benchSize
	N := m * n

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := time.Now().UnixNano()

			for j := 0; j < m; j++ {
				item := newRrdItem2(i+1, now+60)
				es[i][j].put(item)
				if err := es[i][j].commit(storageApp); err != nil {
					glog.V(4).Infof("%s", err)
				}
			}
			stop := time.Now().UnixNano()
			atomic.AddInt64(&ts, stop-start)
		}(i)
	}
	wg.Wait()
	glog.Infof("update %d %d %d ns/op", n, N, ts/int64(m*n*n))
}

func testFetch(n int, t *testing.T) {
	var ts int64

	hds := &storageApp.hdisk
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", testDirs[i], n)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	m := benchSize
	N := m * n

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := time.Now().UnixNano()

			for j := 0; j < m; j++ {
				if _, err := storageApp.taskRrdFetch(es[i][j].hashkey, "AVERAGE",
					now-60, now+600, es[i][j].step); err != nil {
					glog.V(4).Infof("%s", err)
				}
			}
			stop := time.Now().UnixNano()
			atomic.AddInt64(&ts, stop-start)
		}(i)
	}
	wg.Wait()
	glog.Infof("fetch %d %d %d ns/op", n, N, ts/int64(m*n*n))
}

func TestAll(t *testing.T) {
	testStoragetInit()
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testAdd(i+1, t)
	}
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testUpdate(i+1, t)
	}
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testFetch(i+1, t)
	}
	statsRrd()
	testStorageDone()
}
