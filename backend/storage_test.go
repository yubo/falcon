/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

const (
	test_dir      = "/tmp/test"
	b_size        = 10000
	work_nb       = 10
	MAX_HD_NUMBER = 12
)

var (
	test_dirs       []string
	wg              sync.WaitGroup
	lock            sync.RWMutex
	start, end, now int64
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	storageConfig = BackendOpts{
		Storage: StorageOpts{
			Type: "rrdlite",
		},
	}

	test_dirs = make([]string, MAX_HD_NUMBER)

	for i := 0; i < MAX_HD_NUMBER; i++ {
		test_dirs[i] = fmt.Sprintf("%s/hd%02d", test_dir, i)
	}

	err := storageCheckHds(test_dirs)
	if err != nil {
		glog.Fatalf("rrdtool.Start error, bad data dir %v\n", err)
	}

	storageIoTaskCh = make([]chan *ioTask, MAX_HD_NUMBER)
	for i := 0; i < MAX_HD_NUMBER; i++ {
		storageIoTaskCh[i] = make(chan *ioTask, 16)
		removeContents(test_dirs[i])
		go ioWorker(storageIoTaskCh[i])
	}

	timeStart(storageConfig, &appTs)
	now = time.Now().Unix()
	start = now - 120
	end = now + 1800
}

func rrdToEntry(item *specs.RrdItem) *cacheEntry {
	return &cacheEntry{
		createTs:  now,
		host:      item.Host,
		name:      item.Name,
		tags:      item.Tags,
		typ:       item.Type,
		step:      item.Step,
		heartbeat: item.Heartbeat,
		min:       item.Min,
		max:       item.Max,
		dataId:    0,
		commitId:  0,
	}

}

func benchmarkAdd(n int, b *testing.B) {
	var err_cnt uint64
	b.StopTimer()
	hds := &storageConfig.Storage.Hdisks
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", test_dirs[i], n)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	b.N = b_size
	m := b.N / work_nb
	b.StartTimer()

	for i := 0; i < work_nb; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			item := newRrdItem(i)
			e := rrdToEntry(item)
			for j := 0; j < m; j++ {
				e.name = fmt.Sprintf("key_%d_%d", i, j)
				e.hashkey = e.csum()
				if err := e.createRrd(); err != nil {
					if err_cnt < 10 {
						fmt.Println(err)
					}
					atomic.AddUint64(&err_cnt, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("add_err %d\n", err_cnt)
}

func benchmarkUpdate(n int, b *testing.B) {
	var err_cnt uint64
	b.StopTimer()
	hds := &storageConfig.Storage.Hdisks
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", test_dirs[i], i)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	b.N = b_size
	m := b.N / work_nb
	b.StartTimer()

	for i := 0; i < work_nb; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			item := newRrdItem(i)
			e := rrdToEntry(item)
			e.put(item)

			for j := 0; j < m; j++ {
				e.name = fmt.Sprintf("key_%d_%d", i, j)
				e.hashkey = e.csum()
				if err := e.commit(); err != nil {
					atomic.AddUint64(&err_cnt, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("update_err %d\n", err_cnt)
}

func benchmarkFetch(n int, b *testing.B) {
	var err_cnt uint64
	b.StopTimer()
	hds := &storageConfig.Storage.Hdisks
	*hds = make([]string, n)
	for i := 0; i < n; i++ {
		(*hds)[i] = fmt.Sprintf("%s/%d", test_dirs[i], i)
		os.MkdirAll((*hds)[i], os.ModePerm)
	}

	b.N = b_size
	m := b.N / work_nb

	b.StartTimer()

	for i := 0; i < work_nb; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			item := newRrdItem(i)
			e := rrdToEntry(item)
			e.put(item)

			for j := 0; j < m; j++ {
				e.name = fmt.Sprintf("key_%d_%d", i, j)
				e.hashkey = e.csum()
				if _, err := taskRrdFetch(e.hashkey, "AVERAGE",
					start, end, e.step); err != nil {
					atomic.AddUint64(&err_cnt, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("fetch_err %d\n", err_cnt)
}

func BenchmarkAdd1(b *testing.B)     { benchmarkAdd(1, b) }
func BenchmarkUpdate1(b *testing.B)  { benchmarkUpdate(1, b) }
func BenchmarkFetch1(b *testing.B)   { benchmarkFetch(1, b) }
func BenchmarkAdd2(b *testing.B)     { benchmarkAdd(2, b) }
func BenchmarkUpdate2(b *testing.B)  { benchmarkUpdate(2, b) }
func BenchmarkFetch2(b *testing.B)   { benchmarkFetch(2, b) }
func BenchmarkAdd3(b *testing.B)     { benchmarkAdd(3, b) }
func BenchmarkUpdate3(b *testing.B)  { benchmarkUpdate(3, b) }
func BenchmarkFetch3(b *testing.B)   { benchmarkFetch(3, b) }
func BenchmarkAdd4(b *testing.B)     { benchmarkAdd(4, b) }
func BenchmarkUpdate4(b *testing.B)  { benchmarkUpdate(4, b) }
func BenchmarkFetch4(b *testing.B)   { benchmarkFetch(4, b) }
func BenchmarkAdd5(b *testing.B)     { benchmarkAdd(5, b) }
func BenchmarkUpdate5(b *testing.B)  { benchmarkUpdate(5, b) }
func BenchmarkFetch5(b *testing.B)   { benchmarkFetch(5, b) }
func BenchmarkAdd6(b *testing.B)     { benchmarkAdd(6, b) }
func BenchmarkUpdate6(b *testing.B)  { benchmarkUpdate(6, b) }
func BenchmarkFetch6(b *testing.B)   { benchmarkFetch(6, b) }
func BenchmarkAdd7(b *testing.B)     { benchmarkAdd(7, b) }
func BenchmarkUpdate7(b *testing.B)  { benchmarkUpdate(7, b) }
func BenchmarkFetch7(b *testing.B)   { benchmarkFetch(7, b) }
func BenchmarkAdd8(b *testing.B)     { benchmarkAdd(8, b) }
func BenchmarkUpdate8(b *testing.B)  { benchmarkUpdate(8, b) }
func BenchmarkFetch8(b *testing.B)   { benchmarkFetch(8, b) }
func BenchmarkAdd9(b *testing.B)     { benchmarkAdd(9, b) }
func BenchmarkUpdate9(b *testing.B)  { benchmarkUpdate(9, b) }
func BenchmarkFetch9(b *testing.B)   { benchmarkFetch(9, b) }
func BenchmarkAdd10(b *testing.B)    { benchmarkAdd(10, b) }
func BenchmarkUpdate10(b *testing.B) { benchmarkUpdate(10, b) }
func BenchmarkFetch10(b *testing.B)  { benchmarkFetch(10, b) }
func BenchmarkAdd11(b *testing.B)    { benchmarkAdd(11, b) }
func BenchmarkUpdate11(b *testing.B) { benchmarkUpdate(11, b) }
func BenchmarkFetch11(b *testing.B)  { benchmarkFetch(11, b) }
func BenchmarkAdd12(b *testing.B)    { benchmarkAdd(12, b) }
func BenchmarkUpdate12(b *testing.B) { benchmarkUpdate(12, b) }
func BenchmarkFetch12(b *testing.B)  { benchmarkFetch(12, b) }

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
