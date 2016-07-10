/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/rpc"
	"os"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	"stathat.com/c/consistent"

	"github.com/golang/glog"
	"github.com/yubo/rrdlite"

	"github.com/yubo/falcon/specs"
)

const (
	_ = iota
	IO_TASK_M_FILE_READ
	IO_TASK_M_FILE_WRITE
	IO_TASK_M_RRD_ADD
	IO_TASK_M_RRD_UPDATE
	IO_TASK_M_RRD_FETCH
)

const (
	_               = iota
	NET_TASK_M_SEND // no used
	NET_TASK_M_QUERY
	NET_TASK_M_FETCH_COMMIT
)

const (
	RRD_F_MISS uint32 = 1 << iota
	RRD_F_ERR
	RRd_F_SENDING
	RRD_F_FETCHING
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
	storageConfig            BackendOpts
	storageSyncEvent         chan specs.ProcEvent
	storageIoTaskCh          []chan *ioTask
	storageNetTaskCh         map[string]chan *netTask
	storageMigrateClients    map[string][]*rpc.Client
	storageMigrateConsistent *consistent.Consistent
)

type rrdCheckout struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*specs.RRDData
}

type ioTask struct {
	method int
	args   interface{}
	done   chan error
}

type netTask struct {
	Method int
	e      *cacheEntry
	Done   chan error
	Args   interface{}
	Reply  interface{}
}

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
func ktofname(key string) string {
	csum, _ := strconv.ParseUint(key[0:2], 16, 64)
	return fmt.Sprintf("%s/%s/%s.rrd",
		storageConfig.Storage.Hdisks[int(csum)%len(storageConfig.Storage.Hdisks)],
		key[0:2], key)
}

func ktoch(key string) chan *ioTask {
	csum, _ := strconv.ParseUint(key[0:2], 16, 64)
	return storageIoTaskCh[int(csum)%len(storageConfig.Storage.Hdisks)]
}

func rrdCreate(filename string, e *cacheEntry) error {
	now := time.Unix(timeNow(), 0)
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
	if err != nil {
		err = fmt.Errorf("filename:%s %s", filename, err)
	}
	return err
}

// flush to disk from cacheEntry
// call by ioWorker
func ioRrdAdd(e *cacheEntry) (err error) {
	filename := e.filename()

	// unlikely
	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		return specs.ErrExist
	}

	path := path.Dir(filename)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	statInc(ST_RRD_CREAT, 1)
	if err = rrdCreate(filename, e); err != nil {
		statInc(ST_RRD_CREAT_ERR, 1)
	}

	return err
}

func ioRrdUpdate(e *cacheEntry) (err error) {
	if e == nil || e.dataId == 0 {
		return specs.ErrEmpty
	}
	filename := e.filename()
	ds := e.dequeueAll()

	statInc(ST_RRD_UPDATE, 1)
	err = rrdUpdate(filename, e.typ, ds)
	if err != nil {
		statInc(ST_RRD_UPDATE_ERR, 1)

		// unlikely
		_, err = os.Stat(filename)
		if os.IsNotExist(err) {
			path := path.Dir(filename)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				os.MkdirAll(path, os.ModePerm)
			}

			statInc(ST_RRD_CREAT, 1)
			if err = rrdCreate(filename, e); err != nil {
				statInc(ST_RRD_CREAT_ERR, 1)
			} else {
				// retry
				statInc(ST_RRD_UPDATE, 1)
				err = rrdUpdate(filename, e.typ, ds)
				if err != nil {
					statInc(ST_RRD_UPDATE_ERR, 1)
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("filename:%s %s", filename, err)
	}
	return err
}

// call by  ioWorker
func ioRrdFetch(filename string, cf string, start, end int64,
	step int) ([]*specs.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	statInc(ST_RRD_FETCH, 1)
	fetchRes, err := rrdlite.Fetch(filename, cf, start_t, end_t, step_t)
	if err != nil {
		statInc(ST_RRD_FETCH_ERR, 1)
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
	if err != nil {
		err = fmt.Errorf("filename:%s %s", filename, err)
	}
	return ret, nil
}

/* migrate */
func dial(address string, timeout int) (*rpc.Client, error) {
	statInc(ST_RPC_CLI_DIAL, 1)
	d := net.Dialer{Timeout: time.Millisecond *
		time.Duration(timeout)}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		statInc(ST_RPC_CLI_DIAL_ERR, 1)
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			statInc(ST_RPC_CLI_DIAL_ERR, 1)
			conn.Close()
			return nil, err
		}
	}
	return rpc.NewClient(conn), err
}

// TODO addr to node
func reconnection(client **rpc.Client, addr string) {
	var err error

	statInc(ST_RPC_CLI_RECONNECT, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = dial(addr, storageConfig.Migrate.ConnTimeout)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = dial(addr, storageConfig.Migrate.ConnTimeout)
	}
}

func taskFileRead(key string) ([]byte, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_FILE_READ,
		args:   &specs.File{Filename: ktofname(key)},
		done:   done,
	}

	ktoch(key) <- task
	err := <-done
	return task.args.(*specs.File).Data, err
}

// get local data
func taskRrdFetch(key string, cf string, start, end int64,
	step int) ([]*specs.RRDData, error) {
	done := make(chan error, 1)
	task := &ioTask{
		method: IO_TASK_M_RRD_FETCH,
		args: &rrdCheckout{
			filename: ktofname(key),
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	ktoch(key) <- task
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
		err = netRpcCall(*client, "Storage.GetRrd", e.hashkey, &rrdfile,
			time.Duration(storageConfig.Migrate.CallTimeout)*time.Millisecond)

		if err == nil {
			done := make(chan error, 1)
			ktoch(e.hashkey) <- &ioTask{
				method: IO_TASK_M_FILE_WRITE,
				args: &specs.File{
					Filename: e.filename(),
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
			time.Duration(storageConfig.Migrate.CallTimeout)*time.Millisecond)

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
	args interface{}, resp interface{}) (err error) {

	for i := 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Storage.Query", args, resp,
			time.Duration(storageConfig.Migrate.CallTimeout)*time.Millisecond)

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
				statInc(ST_RPC_CLI_SEND, 1)
				if err = netSendData(client, task.e, addr); err != nil {
					statInc(ST_RPC_CLI_SEND_ERR, 1)
				}
			} else if task.Method == NET_TASK_M_QUERY {
				statInc(ST_RPC_CLI_QUERY, 1)
				if err = netQueryData(client, addr, task.Args,
					task.Reply); err != nil {
					statInc(ST_RPC_CLI_QUERY_ERR, 1)
				}
			} else if task.Method == NET_TASK_M_FETCH_COMMIT {
				statInc(ST_RPC_CLI_FETCH, 1)
				if err = netRrdFetch(client, task.e, addr); err != nil {
					statInc(ST_RPC_CLI_FETCH_ERR, 1)
					if os.IsNotExist(err) {
						statInc(ST_RPC_CLI_FETCH_ERR_NOEXIST, 1)
						atomic.StoreUint32(&task.e.flag, 0)
					}
				}
				//warning:异常情况，也写入本地存储
				task.e.commit()
			} else {
				err = errors.New("error net task method")
			}
			if task.Done != nil {
				task.Done <- err
			}
		}
	}
}

func ioWorker(ch chan *ioTask) {
	var err error
	for {
		select {
		case task := <-ch:
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
			} else if task.method == IO_TASK_M_RRD_ADD {
				if args, ok := task.args.(*cacheEntry); ok {
					task.done <- ioRrdAdd(args)
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
	var lastTs int64
	var percent int
	var exit bool

	arg := _arg.(*commitCacheArg)

	now := timeNow()
	expired := now - CACHE_TIME
	nloop := appCache.dataq.size / (CACHE_TIME / FLUSH_DISK_STEP)
	n := 0

	if arg.p.Status() == specs.APP_STATUS_EXIT {
		percent = appCache.dataq.size / 100
		exit = true
	}

	for {
		//e := appCache.dequeue_data()
		l := appCache.dataq.dequeue()
		if l == nil {
			return
		}

		if exit {
			if n%percent == 0 {
				glog.Infof("%d %d%%", n, n/percent)
			}
		} else {
			appCache.dataq.enqueue(l)
		}

		n++
		e := list_data_entry(l)
		flag := atomic.LoadUint32(&e.flag)

		//write err data to local filename
		if arg.migrate && flag&RRD_F_MISS != 0 {
			//PullByKey(key)
			lastTs = e.commitTs
			if lastTs == 0 {
				lastTs = e.createTs
			}
			e.fetchCommit()
		} else {
			//CommitByKey(key)
			lastTs = e.commitTs
			if lastTs == 0 {
				lastTs = e.createTs
			}
			if err := e.commit(); err != nil {
				if err != specs.ErrEmpty {
					glog.Warning(err)
				}
			}
		}
		//glog.V(3).Infof("last %d %d/%d", lastTs, n, cache.dataq.size)

		if exit {
			continue
		}

		if lastTs > expired && n > nloop {
			glog.V(4).Infof("last %d expired %d now %d n %d/%d nloop %d",
				lastTs, expired, now, n, appCache.dataq.size, nloop)
			return
		}
	}
}

func commitCacheWorker(p *specs.Process) {
	var init bool
	var arg commitCacheArg
	arg.p = p

	ticker := falconTicker(time.Second*FIRST_FLUSH_DISK, storageConfig.Debug)

	if storageConfig.Migrate.Enable {
		arg.migrate = true
	}

	for {
		select {
		case <-ticker:
			if !init {
				ticker = falconTicker(time.Second*
					FLUSH_DISK_STEP, storageConfig.Debug)
				init = true
			}
			commitCache(&arg)
		case e := <-storageSyncEvent:
			if e.Method == specs.ROUTINE_EVENT_M_EXIT {
				e.Done <- nil
				return
			}
		}
	}
}

func storageCheckHds(hds []string) error {
	if len(hds) == 0 {
		return specs.ErrEmpty
	}
	for _, dir := range hds {
		if _, err := os.Stat(dir); err != nil {
			return err
		}
	}
	return nil
}

func rrdStart(config BackendOpts, p *specs.Process) {

	storageConfig = config

	err := storageCheckHds(storageConfig.Storage.Hdisks)
	if err != nil {
		glog.Fatalf("rrdtool.Start error, bad data dir %v\n", err)
	}

	hds := &storageConfig.Storage.Hdisks
	storageIoTaskCh = make([]chan *ioTask, len(*hds))
	for i, _ := range *hds {
		storageIoTaskCh[i] = make(chan *ioTask, 16)
	}

	if storageConfig.Migrate.Enable {
		storageMigrateConsistent.NumberOfReplicas = storageConfig.Migrate.Replicas

		for node, addr := range storageConfig.Migrate.Upstream {
			storageMigrateConsistent.Add(node)
			storageNetTaskCh[node] = make(chan *netTask, 16)
			storageMigrateClients[node] = make([]*rpc.Client,
				storageConfig.Migrate.Concurrency)

			for i := 0; i < storageConfig.Migrate.Concurrency; i++ {
				storageMigrateClients[node][i], err = dial(addr,
					storageConfig.Migrate.ConnTimeout)
				if err != nil {
					glog.Fatalf("node:%s addr:%s err:%s\n",
						node, addr, err)
				}
				go netWorker(i, storageNetTaskCh[node],
					&storageMigrateClients[node][i], addr)
			}
		}
	}

	for i, _ := range *hds {
		go ioWorker(storageIoTaskCh[i])
	}

	p.RegisterEvent("rrdSync", storageSyncEvent)
	p.RegisterPostHook("rrdSync", commitCache,
		&commitCacheArg{migrate: false, p: p})
	go commitCacheWorker(p)

	glog.Info("rrdStart ok")
}
