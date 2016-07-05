/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/rpc"
	"os"
	"path"
	"sync/atomic"
	"time"

	"stathat.com/c/consistent"

	"github.com/golang/glog"
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

var (
	rrdConfig            StorageOpts
	rrdSyncEvent         chan specs.ProcEvent
	rrdIoTaskCh          chan *ioTask
	rrdNetTaskCh         map[string]chan *netTask
	rrdMigrateClients    map[string][]*rpc.Client
	rrdMigrateConsistent *consistent.Consistent
)

func rrdCreate(filename string, e *cacheEntry) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(e.step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", e.typ, e.heartbeat, e.min, e.max)

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
		v := math.Abs(float64(data.V))
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if i {
			u.Cache(data.Ts, int(data.V))
		} else {
			u.Cache(data.Ts, float64(data.V))
		}
	}

	return u.Update()
}

// WriteFile writes data to a file named by filename.
// file must not exist
func ioFileWrite(filename string, data []byte, perm os.FileMode) error {
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

// flush to disk from cacheEntry
// call by ioWorker
func ioRrdUpdate(e *cacheEntry) (err error) {
	if e == nil || len(e.cache) == 0 {
		return errors.New("empty items")
	}
	filename := e.filename(rrdConfig.RrdStorage)
	ds := e.dequeueAll()

	err = rrdUpdate(filename, e.typ, ds)
	if err != nil {

		// unlikely
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			path := path.Dir(filename)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				os.MkdirAll(path, os.ModePerm)
			}

			err = rrdCreate(filename, e)
			if err == nil {
				// retry
				err = rrdUpdate(filename, e.typ, ds)
			}
		}
	}
	return err
}

// call by  ioWorker
func ioRrdFetch(filename string, cf string, start, end int64,
	step int) ([]*specs.RRDData, error) {
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
			Ts: ts,
			V:  specs.JsonFloat(val),
		}
		ret[i] = d
	}

	return ret, nil
}

/* migrate */
func dial(address string, timeout int) (*rpc.Client, error) {
	d := net.Dialer{Timeout: time.Millisecond *
		time.Duration(timeout)}
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

// TODO addr to node
func reconnection(client **rpc.Client, addr string) {
	var err error

	statInc(ST_CONN_ERR, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = dial(addr, rrdConfig.Migrate.ConnTimeout)
	statInc(ST_CONN_DIAL, 1)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = dial(addr, rrdConfig.Migrate.ConnTimeout)
		statInc(ST_CONN_DIAL, 1)
	}
}

func taskFileRead(filename string) ([]byte, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_FILE_READ,
		args:   &specs.File{Filename: filename},
		done:   done,
	}

	rrdIoTaskCh <- task
	err := <-done
	return task.args.(*specs.File).Data, err
}

// get local data
func taskRrdFetch(filename string, cf string, start, end int64,
	step int) ([]*specs.RRDData, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_RRD_FETCH,
		args: &rrdCheckout{
			filename: filename,
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	rrdIoTaskCh <- task
	err := <-done
	return task.args.(*rrdCheckout).data, err
}

func netRpcCall(client *rpc.Client, method string, args interface{},
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

/*
 * get remote data to local disk
 */
func netRrdFetch(client **rpc.Client, e *cacheEntry, addr string) error {
	var (
		err     error
		flag    uint32
		i       int
		rrdfile specs.File
	)

	flag = atomic.LoadUint32(&e.flag)

	atomic.StoreUint32(&e.flag, flag|RRD_F_FETCHING)

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Storage.GetRrd", e.key, &rrdfile,
			time.Duration(rrdConfig.Migrate.CallTimeout)*time.Millisecond)

		if err == nil {
			done := make(chan error, 1)
			rrdIoTaskCh <- &ioTask{
				method: IO_TASK_M_FILE_WRITE,
				args: &specs.File{
					Filename: e.filename(rrdConfig.RrdStorage),
					Data:     rrdfile.Data[:],
				},
				done: done,
			}
			if err = <-done; err != nil {
				goto out
			} else {
				flag &= ^RRD_F_MISS
				goto out
			}
		} else {
			glog.Warning(err)
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
out:
	flag &= ^RRD_F_FETCHING
	atomic.StoreUint32(&e.flag, flag)
	return err
}

/* push cacheEntry data
 * by
 * call remote storage.send  */
func netSendData(client **rpc.Client, e *cacheEntry, addr string) error {
	var (
		err  error
		flag uint32
		resp *specs.RpcResp
		i    int
	)

	//remote
	flag = atomic.LoadUint32(&e.flag)

	e.Lock()
	defer e.Unlock()

	items := e._getItems()
	items_size := len(items)
	if items_size == 0 {
		goto out
	}
	resp = &specs.RpcResp{}

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Storage.Send", items, resp,
			time.Duration(rrdConfig.Migrate.CallTimeout)*time.Millisecond)

		if err == nil {
			e._dequeueAll()
			goto out
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
out:
	flag &= ^RRd_F_SENDING
	atomic.StoreUint32(&e.flag, flag)
	return err

}

// called by networker
func netQueryData(client **rpc.Client, addr string,
	args interface{}, resp interface{}) error {
	var (
		err error
		i   int
	)

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Storage.Query", args, resp,
			time.Duration(rrdConfig.Migrate.CallTimeout)*time.Millisecond)

		if err == nil {
			break
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
	return err
}

func netWorker(idx int, ch chan *netTask, client **rpc.Client, addr string) {
	var err error
	for {
		select {
		case task := <-ch:
			if task.Method == NET_TASK_M_SEND {
				if err = netSendData(client, task.e, addr); err != nil {
					statInc(ST_SEND_ERR, 1)
				} else {
					statInc(ST_SEND_SUCCESS, 1)
				}
			} else if task.Method == NET_TASK_M_QUERY {
				if err = netQueryData(client, addr, task.Args,
					task.Reply); err != nil {
					statInc(ST_QUERY_ERR, 1)
				} else {
					statInc(ST_QUERY_SUCCESS, 1)
				}
			} else if task.Method == NET_TASK_M_FETCH {
				if err = netRrdFetch(client, task.e, addr); err != nil {
					if os.IsNotExist(err) {
						//文件不存在时，直接将缓存数据刷入本地
						statInc(ST_FETCH_ISNOTEXIST, 1)
						atomic.StoreUint32(&task.e.flag, 0)
						task.e.commit()
					} else {
						//warning:其他异常情况，缓存数据会堆积
						statInc(ST_FETCH_ERR, 1)
					}
				} else {
					statInc(ST_FETCH_SUCCESS, 1)
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

func ioWorker() {
	var err error
	for {
		select {
		case task := <-rrdIoTaskCh:
			if task.method == IO_TASK_M_FILE_READ {
				if args, ok := task.args.(*specs.File); ok {
					args.Data, err = ioutil.ReadFile(args.Filename)
					task.done <- err
				}
			} else if task.method == IO_TASK_M_FILE_WRITE {
				//filename must not exist
				if args, ok := task.args.(*specs.File); ok {
					_, err := os.Stat(path.Dir(args.Filename))
					if os.IsNotExist(err) {
						task.done <- err
					}
					task.done <- ioFileWrite(args.Filename,
						args.Data, 0644)
				}
			} else if task.method == IO_TASK_M_RRD_UPDATE {
				if args, ok := task.args.(*cacheEntry); ok {
					task.done <- ioRrdUpdate(args)
				}
			} else if task.method == IO_TASK_M_RRD_FETCH {
				if args, ok := task.args.(*rrdCheckout); ok {
					args.data, err = ioRrdFetch(args.filename,
						args.cf, args.start, args.end, args.step)
					task.done <- err
				}
			}
		}
	}
}

type commitCacheArg struct {
	migrate bool
	p       *specs.Process
}

/* called by  commitCacheWorker per FLUSH_DISK_STEP */
func commitCache(_arg interface{}) {
	arg := _arg.(*commitCacheArg)

	expired := time.Now().Unix() - CACHE_TIME
	nloop := len(cache.hash) / (CACHE_TIME / FLUSH_DISK_STEP)
	n := 0

	for {
		e := cache.dequeue_data()
		if e == nil {
			return
		}

		if arg.p.Status() != specs.APP_STATUS_EXIT {
			cache.enqueue_data(e)
		} else {
			if n&0x04ff == 0 {
				glog.Infof("%d", n)
			}
		}

		n++
		glog.V(3).Infof("%d", n)

		flag := atomic.LoadUint32(&e.flag)

		//write err data to local filename
		if arg.migrate && flag&RRD_F_MISS != 0 {
			//PullByKey(key)
			e.fetch()
		}
		//CommitByKey(key)
		commitTs := e.commitTs
		e.commit()

		if !arg.migrate && commitTs > expired && n > nloop {
			return
		}
	}
}

func commitCacheWorker(p *specs.Process) {
	var init bool
	var arg commitCacheArg
	arg.p = p

	ticker := time.NewTicker(time.Second * FIRST_FLUSH_DISK).C

	if rrdConfig.Migrate.Enable {
		arg.migrate = true
	}

	for {
		select {
		case <-ticker:
			if !init {
				ticker = time.NewTicker(time.Second *
					FLUSH_DISK_STEP).C
				init = true
			}
			commitCache(&arg)
		case e := <-rrdSyncEvent:
			if e.Method == specs.ROUTINE_EVENT_M_EXIT {
				e.Done <- nil
				return
			}
		}
	}
}

func rrdStart(config StorageOpts, p *specs.Process) {
	_, err := os.Stat(config.RrdStorage)
	if os.IsNotExist(err) {
		glog.Fatalf("rrdtool.Start error, bad data dir %s %v\n",
			config.RrdStorage, err)
	}
	rrdConfig = config

	if rrdConfig.Migrate.Enable {
		rrdMigrateConsistent.NumberOfReplicas = rrdConfig.Migrate.Replicas

		for node, addr := range rrdConfig.Migrate.Upstream {
			rrdMigrateConsistent.Add(node)
			rrdNetTaskCh[node] = make(chan *netTask, 16)
			rrdMigrateClients[node] = make([]*rpc.Client,
				rrdConfig.Migrate.Concurrency)

			for i := 0; i < rrdConfig.Migrate.Concurrency; i++ {
				rrdMigrateClients[node][i], err = dial(addr,
					rrdConfig.Migrate.ConnTimeout)
				if err != nil {
					glog.Fatalf("node:%s addr:%s err:%s\n",
						node, addr, err)
				}
				go netWorker(i, rrdNetTaskCh[node],
					&rrdMigrateClients[node][i], addr)
			}
		}
	}

	go ioWorker()

	p.RegisterEvent("rrdSync", rrdSyncEvent)
	p.RegisterPostHook("rrdSync", commitCache,
		&commitCacheArg{migrate: false, p: p})
	go commitCacheWorker(p)

	glog.Info("rrdStart ok")
}
