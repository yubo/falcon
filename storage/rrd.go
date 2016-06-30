/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/rpc"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/yubo/rrdlite"

	"github.com/yubo/falcon/specs"
)

// RRA.Point.Size
const (
	RRA1PointCnt   = 720 // 1m一个点存12h
	RRA5PointCnt   = 576 // 5m一个点存2d
	RRA20PointCnt  = 504 // 20m一个点存7d
	RRA180PointCnt = 766 // 3h一个点存3month
	RRA720PointCnt = 730 // 12h一个点存1year
)

func rrdCreate(filename string, e *cacheEntry) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(e.step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", e.dsType, e.heartbeat, e.min, e.max)

	// 设置各种归档策略
	// 1分钟一个点存 12小时
	c.RRA("AVERAGE", 0.5, 1, RRA1PointCnt)

	// 5m一个点存2d
	c.RRA("AVERAGE", 0.5, 5, RRA5PointCnt)
	c.RRA("MAX", 0.5, 5, RRA5PointCnt)
	c.RRA("MIN", 0.5, 5, RRA5PointCnt)

	// 20m一个点存7d
	c.RRA("AVERAGE", 0.5, 20, RRA20PointCnt)
	c.RRA("MAX", 0.5, 20, RRA20PointCnt)
	c.RRA("MIN", 0.5, 20, RRA20PointCnt)

	// 3小时一个点存3个月
	c.RRA("AVERAGE", 0.5, 180, RRA180PointCnt)
	c.RRA("MAX", 0.5, 180, RRA180PointCnt)
	c.RRA("MIN", 0.5, 180, RRA180PointCnt)

	// 12小时一个点存1year
	c.RRA("AVERAGE", 0.5, 720, RRA720PointCnt)
	c.RRA("MAX", 0.5, 720, RRA720PointCnt)
	c.RRA("MIN", 0.5, 720, RRA720PointCnt)

	return c.Create(true)
}

func rrdUpdate(filename string, dsType string, ds []*specs.RRDData) error {
	var i bool
	u := rrdlite.NewUpdater(filename)

	if dsType == "DERIVE" || dsType == "COUNTER" {
		i = true
	}

	for _, data := range ds {
		v := math.Abs(float64(data.Value))
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if i {
			u.Cache(data.Timestamp, int(data.Value))
		} else {
			u.Cache(data.Timestamp, float64(data.Value))
		}
	}

	return u.Update()
}

// flush to disk from cacheEntry
func rrdCommit(e *cacheEntry) (err error) {
	if e == nil || len(e.cache) == 0 {
		return errors.New("empty items")
	}
	filename := e.filename(config().RrdStorage)
	ds := e.dequeueAll()

	err = rrdUpdate(filename, e.dsType, ds)
	if err != nil {

		// unlikely
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			_, err = os.Stat(path.Dir(filename))
			if os.IsNotExist(err) {
				os.MkdirAll(filename, os.ModePerm)
			}

			err = rrdCreate(filename, e)
			if err == nil {
				// retry
				err = rrdUpdate(filename, e.dsType, ds)
			}
		}
	}
	return err
}

func taskReadFile(filename string) ([]byte, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_READ,
		args:   &readfile_t{filename: filename},
		done:   done,
	}

	io_task_chan <- task
	err := <-done
	return task.args.(*readfile_t).data, err
}

/*
 * use  (*cacheEntry)->commit()
 *
func taskFlushFile(filename string, items []*specs.GraphItem) error {
	done := make(chan error, 1)
	io_task_chan <- &io_task_t{
		method: IO_TASK_M_COMMIT,
		args: &flushfile_t{
			filename: filename,
			items:    items,
		},
		done: done,
	}
	stat_inc(ST_DISK_COUNTER, 1)
	return <-done
}
*/

// get local data
func taskCheckout(filename string, cf string, start, end int64, step int) ([]*specs.RRDData, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_CHECKOUT,
		args: &rrdCheckout_t{
			filename: filename,
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	io_task_chan <- task
	err := <-done
	return task.args.(*rrdCheckout_t).data, err
}

func fetch(filename string, cf string, start, end int64, step int) ([]*specs.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	fetchRes, err := rrdlite.Fetch(filename, cf, start_t, end_t, step_t)
	if err != nil {
		return []*specs.RRDData{}, err
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*specs.RRDData, size)

	start_ts := fetchRes.Start.Unix()
	step_s := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := start_ts + int64(i+1)*int64(step_s)
		d := &specs.RRDData{
			Timestamp: ts,
			Value:     specs.JsonFloat(val),
		}
		ret[i] = d
	}

	return ret, nil
}

/*
func CommitByKey(key string) {

	md5, dsType, step, err := SplitRrdCacheKey(key)
	if err != nil {
		return
	}
	filename := RrdFileName(Config().RRD.Storage, md5, dsType, step)

	items := store.GraphItems.PopAll(key)
	if len(items) == 0 {
		return
	}
	FlushFile(filename, items)
}

func PullByKey(key string) {
	done := make(chan error)

	item := store.GraphItems.First(key)
	if item == nil {
		return
	}
	node, err := Consistent.Get(item.PrimaryKey())
	if err != nil {
		return
	}
	Net_task_ch[node] <- &Net_task_t{
		Method: NET_TASK_M_FETCH,
		Key:    key,
		Done:   done,
	}
	// net_task slow, shouldn't block syncDisk() or FlushAll()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			log.Printf("get %s from remote err[%s]\n", key, err)
			return
		}
		stat_inc(ST_NET_COUNTER, 1)
		//todo: flushfile after getfile? not yet
	}()
}
*/

// WriteFile writes data to a file named by filename.
// file must not exist
func _writeFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

/* migrate */
func dial(address string, timeout time.Duration) (*rpc.Client, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return rpc.NewClient(conn), err
}

func taskWorker(idx int, ch chan *Net_task_t, client **rpc.Client, addr string) {
	var err error
	for {
		select {
		case task := <-ch:
			if task.Method == NET_TASK_M_SEND {
				if err = send_data(client, task.e, addr); err != nil {
					stat_inc(ST_SEND_S_ERR, 1)
				} else {
					stat_inc(ST_SEND_S_SUCCESS, 1)
				}
			} else if task.Method == NET_TASK_M_QUERY {
				if err = query_data(client, addr, task.Args, task.Reply); err != nil {
					stat_inc(ST_QUERY_S_ERR, 1)
				} else {
					stat_inc(ST_QUERY_S_SUCCESS, 1)
				}
			} else if task.Method == NET_TASK_M_FETCH {
				if err = taskFetchRrd(client, task.e, addr); err != nil {
					if os.IsNotExist(err) {
						//文件不存在时，直接将缓存数据刷入本地
						stat_inc(ST_FETCH_S_ISNOTEXIST, 1)
						//GraphItems.SetFlag(task.Key, 0)
						atomic.StoreUint32(&task.e.flag, 0)
						//CommitByKey(task.Key)
						task.e.commit()
					} else {
						//warning:其他异常情况，缓存数据会堆积
						stat_inc(ST_FETCH_S_ERR, 1)
					}
				} else {
					stat_inc(ST_FETCH_S_SUCCESS, 1)
				}
			} else {
				err = errors.New("error net task method")
			}
			if task.Done != nil {
				task.Done <- err
			}
		}
	}
}

// TODO addr to node
func reconnection(client **rpc.Client, addr string) {
	var err error

	stat_inc(ST_CONN_S_ERR, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = dial(addr, time.Second)
	stat_inc(ST_CONN_S_DIAL, 1)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = dial(addr, time.Second)
		stat_inc(ST_CONN_S_DIAL, 1)
	}
}

func query_data(client **rpc.Client, addr string,
	args interface{}, resp interface{}) error {
	var (
		err error
		i   int
	)

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.Query", args, resp,
			time.Duration(config().CallTimeout)*time.Millisecond)

		if err == nil {
			break
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
	return err
}

func send_data(client **rpc.Client, e *cacheEntry, addr string) error {
	var (
		err  error
		flag uint32
		resp *specs.SimpleRpcResponse
		i    int
	)

	//remote
	cfg := config()
	flag = atomic.LoadUint32(&e.flag)

	e.Lock()
	defer e.Unlock()

	items := e._getItems()
	items_size := len(items)
	if items_size == 0 {
		goto out
	}
	resp = &specs.SimpleRpcResponse{}

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.Send", items, resp,
			time.Duration(cfg.CallTimeout)*time.Millisecond)

		if err == nil {
			e._dequeueAll()
			goto out
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
out:
	flag &= ^GRAPH_F_SENDING
	atomic.StoreUint32(&e.flag, flag)
	return err

}

/*
 * get remote data to local disk
 */

func taskFetchRrd(client **rpc.Client, e *cacheEntry, addr string) error {
	var (
		err     error
		flag    uint32
		i       int
		rrdfile specs.File
	)

	cfg := config()

	flag = atomic.LoadUint32(&e.flag)

	//GraphItems.SetFlag(key, flag|GRAPH_F_FETCHING)
	atomic.StoreUint32(&e.flag, flag|GRAPH_F_FETCHING)

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.GetRrd", e.key, &rrdfile,
			time.Duration(cfg.CallTimeout)*time.Millisecond)

		if err == nil {
			done := make(chan error, 1)
			io_task_chan <- &io_task_t{
				method: IO_TASK_M_WRITE,
				args: &specs.File{
					Filename: e.filename(cfg.RrdStorage),
					Body:     rrdfile.Body[:],
				},
				done: done,
			}
			if err = <-done; err != nil {
				goto out
			} else {
				flag &= ^GRAPH_F_MISS
				goto out
			}
		} else {
			log.Println(err)
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
out:
	flag &= ^GRAPH_F_FETCHING
	//GraphItems.SetFlag(key, flag)
	atomic.StoreUint32(&e.flag, flag)
	return err
}

func rpc_call(client *rpc.Client, method string, args interface{},
	reply interface{}, timeout time.Duration) error {
	done := make(chan *rpc.Call, 1)
	client.Go(method, args, reply, done)
	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[rpc]")
	case call := <-done:
		if call.Error == nil {
			return nil
		} else {
			return call.Error
		}
	}
}

/* called by  syncDiskWorker per FLUSH_DISK_STEP */
func commitCaches(force bool) {
	expired := time.Now().Unix() - CACHE_TIME
	nloop := len(cache.hash) / (CACHE_TIME / FLUSH_DISK_STEP)
	n := 0

	for {
		e := cache.dequeue_data()
		if e == nil {
			return
		}
		n++

		cache.enqueue_data(e)

		flag := atomic.LoadUint32(&e.flag)

		//write err data to local filename
		if force == false && config().Migrate.Enable && flag&GRAPH_F_MISS != 0 {
			//PullByKey(key)
			e.fetch()
		}
		//CommitByKey(key)
		commitTs := e.commitTs
		e.commit()

		if !force && commitTs > expired && n > nloop {
			return
		}
	}
}

func syncDiskWorker() {
	time.Sleep(time.Second * 300)
	ticker := time.NewTicker(time.Millisecond * FLUSH_DISK_STEP).C

	for {
		select {
		case <-ticker:
			commitCaches(false)
		case done := <-sync_exit:
			log.Println("cron recv sigout and exit...")
			done <- nil
			return
		}
	}
}

func ioWorker() {
	var err error
	for {
		select {
		case task := <-io_task_chan:
			if task.method == IO_TASK_M_READ {
				if args, ok := task.args.(*readfile_t); ok {
					args.data, err = ioutil.ReadFile(args.filename)
					task.done <- err
				}
			} else if task.method == IO_TASK_M_WRITE {
				//filename must not exist
				if args, ok := task.args.(*specs.File); ok {
					_, err := os.Stat(path.Dir(args.Filename))
					if os.IsNotExist(err) {
						task.done <- err
					}
					task.done <- _writeFile(args.Filename, args.Body, 0644)
				}
			} else if task.method == IO_TASK_M_COMMIT {
				if args, ok := task.args.(*cacheEntry); ok {
					task.done <- rrdCommit(args)
				}
			} else if task.method == IO_TASK_M_CHECKOUT {
				if args, ok := task.args.(*rrdCheckout_t); ok {
					args.data, err = fetch(args.filename, args.cf, args.start, args.end, args.step)
					task.done <- err
				}
			}
		}
	}
}

func rrdStart() {
	cfg := config()
	_, err := os.Stat(cfg.RrdStorage)
	if os.IsNotExist(err) {
		log.Printf("rrdtool.Start error, bad data dir %s %v\n",
			cfg.RrdStorage, err)
		os.Exit(2)
	}

	if cfg.Migrate.Enable {
		Consistent.NumberOfReplicas = cfg.Migrate.Replicas

		for node, addr := range cfg.Migrate.Cluster {
			Consistent.Add(node)
			Net_task_ch[node] = make(chan *Net_task_t, 16)
			clients[node] = make([]*rpc.Client, cfg.Migrate.Concurrency)

			for i := 0; i < cfg.Migrate.Concurrency; i++ {
				if clients[node][i], err = dial(addr, time.Second); err != nil {
					log.Fatalf("node:%s addr:%s err:%s\n", node, addr, err)
				}
				go taskWorker(i, Net_task_ch[node], &clients[node][i], addr)
			}
		}
	}

	registerExitChans(sync_exit)
	go syncDiskWorker()
	go ioWorker()

	log.Println("rrdStart ok")
}
