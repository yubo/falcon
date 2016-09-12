/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

const (
	test_dir      = "/tmp/falcon"
	b_size        = 100
	work_nb       = 4
	MAX_HD_NUMBER = 2
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

	os.RemoveAll(test_dir)
	test_dirs = make([]string, MAX_HD_NUMBER)

	for i := 0; i < MAX_HD_NUMBER; i++ {
		test_dirs[i] = fmt.Sprintf("%s/hdd%d", test_dir, i)
		os.MkdirAll(test_dirs[i], 0755)
	}

	err := storageCheckHds(test_dirs)
	if err != nil {
		glog.Fatalf("rrdtool.Start error, bad data dir %v\n", err)
	}

	storageIoTaskCh = make([]chan *ioTask, MAX_HD_NUMBER)
	for i := 0; i < MAX_HD_NUMBER; i++ {
		storageIoTaskCh[i] = make(chan *ioTask, 32)
		removeContents(test_dirs[i])
		go ioWorker(storageIoTaskCh[i])
	}

	timeStart(storageConfig, &appTs)
	now = time.Now().Unix()
	start = now - 120
	end = now + 1800
}

func rrdToEntry(item *specs.RrdItem) (*cacheEntry, error) {
	e, _ := appCache.getPoolEntry()
	if err := e.reset(now, item.Host, item.Name, item.Tags, item.Type,
		item.Step, item.Heartbeat, item.Min[0], item.Max[0]); err != nil {
		appCache.putPoolEntry(e)
		return nil, err
	}
	return e, nil
}

func benchmarkAdd(n int, b *testing.B) {
	var err_cnt, cnt uint64
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
			for j := 0; j < m; j++ {
				e, _ := rrdToEntry(item)
				e.setName(fmt.Sprintf("key_%d_%d", i, j))
				e.setHashkey(e.csum())
				if err := e.createRrd(); err != nil {
					if err_cnt < 10 {
						fmt.Println(err)
					}
					atomic.AddUint64(&err_cnt, 1)
				}
				atomic.AddUint64(&cnt, 1)
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("add_err %d\n", err_cnt)
	fmt.Printf("add number: %d\n", cnt)
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
			e, _ := rrdToEntry(item)
			e.put(item)

			for j := 0; j < m; j++ {
				e.setName(fmt.Sprintf("key_%d_%d", i, j))
				e.setHashkey(e.csum())
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
			e, _ := rrdToEntry(item)
			e.put(item)

			for j := 0; j < m; j++ {
				e.setName(fmt.Sprintf("key_%d_%d", i, j))
				e.setHashkey(e.csum())
				if _, err := taskRrdFetch(e.hashkey(), "AVERAGE",
					start, end, int(e.e.step)); err != nil {
					atomic.AddUint64(&err_cnt, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("fetch_err %d\n", err_cnt)
}

func BenchmarkAdd01(b *testing.B) { benchmarkAdd(1, b) }
func BenchmarkAdd02(b *testing.B) { benchmarkAdd(2, b) }
func BenchmarkAdd03(b *testing.B) { benchmarkAdd(3, b) }
func BenchmarkAdd04(b *testing.B) { benchmarkAdd(4, b) }

//func BenchmarkAdd05(b *testing.B) { benchmarkAdd(5, b) }
//func BenchmarkAdd06(b *testing.B) { benchmarkAdd(6, b) }
//func BenchmarkAdd07(b *testing.B) { benchmarkAdd(7, b) }
//func BenchmarkAdd08(b *testing.B) { benchmarkAdd(8, b) }
//func BenchmarkAdd09(b *testing.B) { benchmarkAdd(9, b) }
//func BenchmarkAdd10(b *testing.B) { benchmarkAdd(10, b) }
//func BenchmarkAdd11(b *testing.B) { benchmarkAdd(11, b) }
//func BenchmarkAdd12(b *testing.B) { benchmarkAdd(12, b) }
//
//func BenchmarkFetch01(b *testing.B) { benchmarkFetch(1, b) }
//func BenchmarkFetch02(b *testing.B) { benchmarkFetch(2, b) }
//func BenchmarkFetch03(b *testing.B) { benchmarkFetch(3, b) }
//func BenchmarkFetch04(b *testing.B) { benchmarkFetch(4, b) }
//func BenchmarkFetch05(b *testing.B) { benchmarkFetch(5, b) }
//func BenchmarkFetch06(b *testing.B) { benchmarkFetch(6, b) }
//func BenchmarkFetch07(b *testing.B) { benchmarkFetch(7, b) }
//func BenchmarkFetch08(b *testing.B) { benchmarkFetch(8, b) }
//func BenchmarkFetch09(b *testing.B) { benchmarkFetch(9, b) }
//func BenchmarkFetch10(b *testing.B) { benchmarkFetch(10, b) }
//func BenchmarkFetch11(b *testing.B) { benchmarkFetch(11, b) }
//func BenchmarkFetch12(b *testing.B) { benchmarkFetch(12, b) }
//
//func BenchmarkUpdate01(b *testing.B) { benchmarkUpdate(1, b) }
//func BenchmarkUpdate02(b *testing.B) { benchmarkUpdate(2, b) }
//func BenchmarkUpdate03(b *testing.B) { benchmarkUpdate(3, b) }
//func BenchmarkUpdate04(b *testing.B) { benchmarkUpdate(4, b) }
//func BenchmarkUpdate05(b *testing.B) { benchmarkUpdate(5, b) }
//func BenchmarkUpdate06(b *testing.B) { benchmarkUpdate(6, b) }
//func BenchmarkUpdate07(b *testing.B) { benchmarkUpdate(7, b) }
//func BenchmarkUpdate08(b *testing.B) { benchmarkUpdate(8, b) }
//func BenchmarkUpdate09(b *testing.B) { benchmarkUpdate(9, b) }
//func BenchmarkUpdate10(b *testing.B) { benchmarkUpdate(10, b) }
//func BenchmarkUpdate11(b *testing.B) { benchmarkUpdate(11, b) }
//func BenchmarkUpdate12(b *testing.B) { benchmarkUpdate(12, b) }
/*
[work@lg-hadoop-prc-st31 backend]$ go test -bench=.
4805404 /home/work/yubo/gopath/src/github.com/yubo/falcon/backend/cache_test.go 47 true
c.createEntry success

c.get success
c.unlink success
all c.unlink success
e.getItems() success
PASS
BenchmarkAdd01-24       add number: 9984
   10000           2417394 ns/op
BenchmarkAdd02-24       add number: 9984
   10000           2341798 ns/op
BenchmarkAdd03-24       add number: 9984
   10000           2426143 ns/op
BenchmarkAdd04-24       add number: 9984
   10000           2628815 ns/op
ok      github.com/yubo/falcon/backend  100.488s
*/
