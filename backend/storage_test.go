/*
 * Copyright 2016 yubo. All rights reserved.
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
)

const (
	testDir       = "/tmp"
	benchSize     = 20
	MAX_HD_NUMBER = 1
)

var (
	storageApp *Backend
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

func test_storage_init() (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()

	storageApp = new(Backend)
	storageApp.Conf = falcon.ConfBackend{
		ShmMagic: 0x80386,
		ShmKey:   0x6020,
		ShmSize:  1 << 28, // 256m
		Storage: falcon.Storage{
			Type: "rrdlite",
		},
	}

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
		go storageApp.ioWorker(storageApp.storageIoTaskCh[i])
	}

	storageApp.cacheInit()
	if err := storageApp.cacheReset(); err != nil {
		fmt.Println(err)
	}

	storageApp.timeStart()
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
			es[i][j].setName(fmt.Sprintf("key_%d_%d", i, j))
			es[i][j].setHashkey(es[i][j].csum())
		}
	}
	return nil
}

func (p *Backend) rrdToEntry(item *falcon.RrdItem) (*cacheEntry, error) {
	e, err := p.getPoolEntry()
	if err != nil {
		glog.V(4).Infoln(err)
		return nil, err
	}
	if err = e.reset(now, item.Host, item.Name, item.Tags, item.Type,
		item.Step, item.Heartbeat, item.Min[0], item.Max[0]); err != nil {
		p.putPoolEntry(e)
		glog.V(4).Infoln(err)
		return nil, err
	}
	return e, nil
}

func testAdd(n int, t *testing.T) {
	var ts int64

	if storageApp == nil {
		test_storage_init()
	}

	hds := &storageApp.Conf.Storage.Hdisks
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

	if storageApp == nil {
		test_storage_init()
	}

	hds := &storageApp.Conf.Storage.Hdisks
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

	if storageApp == nil {
		test_storage_init()
	}
	hds := &storageApp.Conf.Storage.Hdisks
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
				if _, err := storageApp.taskRrdFetch(es[i][j].hashkey(), "AVERAGE",
					now-60, now+600, int(es[i][j].e.step)); err != nil {
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
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testAdd(i+1, t)
	}
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testUpdate(i+1, t)
	}
	for i := 0; i < MAX_HD_NUMBER; i++ {
		testFetch(i+1, t)
	}
	statRrd()
}
