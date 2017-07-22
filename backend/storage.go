/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
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
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend/config"
	"github.com/yubo/rrdlite"
	"golang.org/x/net/context"
	"stathat.com/c/consistent"
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
	_ = iota
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

type rrdCheckout struct {
	filename string
	cf       falcon.Cf
	start    int64
	end      int64
	step     int
	data     []*falcon.RRDData
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

func (p *StorageModule) rrdCreate(filename string, e *cacheEntry) error {
	var (
		typ, min, max string
	)

	now := time.Unix(p.b.timeNow(), 0)
	start := now.Add(time.Duration(-1) * time.Hour)
	heartbeat := int(e.step * 2)

	switch e.typ {
	case falcon.ItemType_GAUGE:
		typ = "GAUGE"
		min = "U"
		max = "U"
	case falcon.ItemType_DERIVE:
		typ = "DERIVE"
		min = "0"
		max = "U"
	case falcon.ItemType_COUNTER:
		typ = "COUNTER"
		min = "0"
		max = "U"
	default:
		return falcon.EINVAL
	}

	c := rrdlite.NewCreator(filename, start, uint(e.step))
	c.DS("metric", typ, heartbeat, min, max)

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

func rrdUpdate(filename string, dsType string, ds []*falcon.RRDData) error {
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
func (p *StorageModule) ioRrdAdd(e *cacheEntry) (err error) {
	filename := e.filename(p.b)

	// unlikely
	/*
		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			return falcon.ErrExist
		}
	*/

	path := path.Dir(filename)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	statsInc(ST_RRD_CREAT, 1)
	if err = p.rrdCreate(filename, e); err != nil {
		statsInc(ST_RRD_CREAT_ERR, 1)
	}

	return err
}

func (p *StorageModule) ioRrdUpdate(e *cacheEntry) (err error) {
	if e == nil || e.dataId == 0 {
		return falcon.ErrEmpty
	}
	filename := e.filename(p.b)
	ds := e.dequeueAll()

	statsInc(ST_RRD_UPDATE, 1)
	err = rrdUpdate(filename, e.typ.String(), ds)
	if err != nil {
		statsInc(ST_RRD_UPDATE_ERR, 1)

		// unlikely
		_, err = os.Stat(filename)
		if os.IsNotExist(err) {
			path := path.Dir(filename)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				os.MkdirAll(path, os.ModePerm)
			}

			statsInc(ST_RRD_CREAT, 1)
			if err = p.rrdCreate(filename, e); err != nil {
				statsInc(ST_RRD_CREAT_ERR, 1)
			} else {
				// retry
				statsInc(ST_RRD_UPDATE, 1)
				err = rrdUpdate(filename, e.typ.String(), ds)
				if err != nil {
					statsInc(ST_RRD_UPDATE_ERR, 1)
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
func ioRrdFetch(filename string, cf falcon.Cf, start, end int64,
	step int) ([]*falcon.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	statsInc(ST_RRD_FETCH, 1)
	fetchRes, err := rrdlite.Fetch(filename, cf.String(), start_t, end_t, step_t)
	if err != nil {
		statsInc(ST_RRD_FETCH_ERR, 1)
		return []*falcon.RRDData{}, err
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*falcon.RRDData, size)

	start_ts := fetchRes.Start.Unix()
	step_s := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := start_ts + int64(i+1)*int64(step_s)
		d := &falcon.RRDData{
			Ts: ts,
			V:  val,
		}
		ret[i] = d
	}
	if err != nil {
		err = fmt.Errorf("filename:%s %s", filename, err)
	}
	return ret, nil
}

/*
 * get remote data to local disk
 */
func (p *StorageModule) netRrdFetch(client BackendClient, e *cacheEntry,
	addr string) error {
	var (
		flag    uint32
		rrdfile falcon.File
	)

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(p.callTimeout)*time.Millisecond)

	flag = atomic.LoadUint32(&e.flag)

	atomic.StoreUint32(&e.flag, flag|RRD_F_FETCHING)

	_, err := client.GetRrd(ctx, &falcon.GetRrdRequest{Key: []byte(e.hashkey)})

	if err == nil {
		done := make(chan error, 1)
		p.b.ktoch(e.hashkey) <- &ioTask{
			method: IO_TASK_M_FILE_WRITE,
			args: &falcon.File{
				Name: e.filename(p.b),
				Data: rrdfile.Data[:],
			},
			done: done,
		}
		if err = <-done; err == nil {
			flag &= ^RRD_F_MISS
		}
	} else {
		glog.Warning(MODULE_NAME, err)
	}
	flag &= ^RRD_F_FETCHING
	atomic.StoreUint32(&e.flag, flag)
	return err
	/*
	   	for i = 0; i < CONN_RETRY; i++ {
	   		err = netRpcCall(*client, "Storage.GetRrd", e.hashkey, &rrdfile,
	   			time.Duration(p.callTimeout)*time.Millisecond)

	   		if err == nil {
	   			done := make(chan error, 1)
	   			p.b.ktoch(e.hashkey) <- &ioTask{
	   				method: IO_TASK_M_FILE_WRITE,
	   				args: &falcon.File{
	   					Name: []byte(e.filename(p.b)),
	   					Data: rrdfile.Data[:],
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
	   			glog.Warning(MODULE_NAME, err)
	   		}
	   		if err == rpc.ErrShutdown {
	   			p.reconnection(client, addr)
	   		}
	   	}
	   out:
	   	flag &= ^RRD_F_FETCHING
	   	atomic.StoreUint32(&e.flag, flag)
	   	return err
	*/
}

// called by networker
func (p *StorageModule) netQueryData(client BackendClient, addr string,
	args interface{}, resp interface{}) (err error) {

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(p.callTimeout)*time.Millisecond)

	resp, err = client.Get(ctx, args.(*falcon.GetRequest))
	return err
	/*
		for i := 0; i < CONN_RETRY; i++ {
			err = netRpcCall(*client, "Storage.Query", args, resp,
				time.Duration(p.callTimeout)*time.Millisecond)

			if err == nil {
				break
			}
			if err == rpc.ErrShutdown {
				p.reconnection(client, addr)
			}
		}
		return err
	*/
}

func (p *StorageModule) netWorker(idx int, ch chan *netTask, client BackendClient, addr string) {
	var err error
	for {
		select {
		case <-p.ctx.Done():
			return
		case task := <-ch:
			if task.Method == NET_TASK_M_QUERY {
				statsInc(ST_RPC_CLI_QUERY, 1)

				if err = p.netQueryData(client, addr, task.Args,
					task.Reply); err != nil {
					statsInc(ST_RPC_CLI_QUERY_ERR, 1)
				}
			} else if task.Method == NET_TASK_M_FETCH_COMMIT {
				statsInc(ST_RPC_CLI_FETCH, 1)
				if err = p.netRrdFetch(client, task.e, addr); err != nil {
					statsInc(ST_RPC_CLI_FETCH_ERR, 1)
					if os.IsNotExist(err) {
						statsInc(ST_RPC_CLI_FETCH_ERR_NOEXIST, 1)
						atomic.StoreUint32(&task.e.flag, 0)
					}
				}
				//warning:异常情况，也写入本地存储
				task.e.commit(p.b)
			} else {
				err = errors.New("error net task method")
			}
			if task.Done != nil {
				task.Done <- err
			}
		}
	}
}

// ioworker never return
func (p *StorageModule) ioWorker(ch chan *ioTask) {
	var err error
	for {
		select {
		case <-p.ctx.Done():
			for i, _ := range p.b.hdisk {
				close(p.b.storageIoTaskCh[i])
			}
			return
		case task := <-ch:
			if task.method == IO_TASK_M_FILE_READ {
				if args, ok := task.args.(*falcon.File); ok {
					args.Data, err = ioutil.ReadFile(args.Name)
					task.done <- err
				}
			} else if task.method == IO_TASK_M_FILE_WRITE {
				//filename must not exist
				if args, ok := task.args.(*falcon.File); ok {
					_, err := os.Stat(path.Dir(args.Name))
					if os.IsNotExist(err) {
						task.done <- err
					}
					task.done <- ioFileWrite(args.Name,
						args.Data, 0644)
				}
			} else if task.method == IO_TASK_M_RRD_ADD {
				if args, ok := task.args.(*cacheEntry); ok {
					task.done <- p.ioRrdAdd(args)
				}
			} else if task.method == IO_TASK_M_RRD_UPDATE {
				if args, ok := task.args.(*cacheEntry); ok {
					task.done <- p.ioRrdUpdate(args)
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
	p       *falcon.Process
}

/* called by  commitCacheWorker per FLUSH_DISK_STEP
 * or called after stop
 */
func (p *StorageModule) commitCache(whole bool) {
	var lastTs int64
	var percent int

	now := p.b.timeNow()
	expired := now - CACHE_TIME
	nloop := len(p.b.cache.hash) / (CACHE_TIME / FLUSH_DISK_STEP)
	n := 0

	if whole {
		percent = len(p.b.cache.hash) / 100
	}

	for {
		//e := p.cache.dequeue_data()
		l := p.b.cache.dataq.dequeue()
		if l == nil {
			return
		}

		if whole {
			if n%percent == 0 {
				glog.Infof(MODULE_NAME+"%d %d%%", n, n/percent)
			}
		} else {
			p.b.cache.dataq.enqueue(l)
		}

		n++
		e := list_data_entry(l)
		flag := atomic.LoadUint32(&e.flag)

		//write err data to local filename
		if !p.migrate.Disabled && flag&RRD_F_MISS != 0 {
			//PullByKey(key)
			lastTs = e.commitTs
			if lastTs == 0 {
				lastTs = e.createTs
			}
			e.fetchCommit(p.b)
		} else {
			//CommitByKey(key)
			lastTs = e.commitTs
			if lastTs == 0 {
				lastTs = e.createTs
			}
			if err := e.commit(p.b); err != nil {
				if err != falcon.ErrEmpty {
					glog.Warning(MODULE_NAME, err)
				}
			}
		}
		//glog.V(3).Infof("last %d %d/%d", lastTs, n, cache.dataq.size)

		if whole {
			continue
		}

		if lastTs > expired && n > nloop {
			glog.V(4).Infof(MODULE_NAME+"last %d expired %d now %d n %d/%d nloop %d",
				lastTs, expired, now, n, len(p.b.cache.hash), nloop)
			return
		}
	}
}

func (p *StorageModule) commitCacheWorker() {
	var init bool

	debug := p.debug
	ticker := falconTicker(time.Second*FIRST_FLUSH_DISK, debug)

	for {
		select {
		case <-p.ctx.Done():
			p.commitCache(true)
			return
		case <-ticker:
			if !init {
				ticker = falconTicker(time.Second*
					FLUSH_DISK_STEP, debug)
				init = true
			}
			p.commitCache(false)
		}
	}
}

func storageCheckHds(hds []string) error {
	if len(hds) == 0 {
		return falcon.ErrEmpty
	}
	for _, dir := range hds {
		if _, err := os.Stat(dir); err != nil {
			return err
		}
	}
	return nil
}

type StorageModule struct {
	b                     *Backend
	storageMigrateClients map[string][]BackendClient
	connTimeout           int
	callTimeout           int
	payloadsize           int
	workerProcesses       int
	debug                 int
	migrate               config.Migrate
	ctx                   context.Context
	cancel                context.CancelFunc
}

func (p *StorageModule) prestart(b *Backend) error {
	p.b = b
	p.migrate = b.Conf.Migrate
	p.storageMigrateClients = make(map[string][]BackendClient)

	b.storageNetTaskCh = make(map[string]chan *netTask)
	b.storageMigrateConsistent = consistent.New()

	return nil
}

func (p *StorageModule) start(b *Backend) error {
	b.hdisk = strings.Split(b.Conf.Configer.Str(C_HDISK), ",")

	p.workerProcesses, _ = b.Conf.Configer.Int(C_WORKER_PROCESSES)
	p.connTimeout, _ = b.Conf.Configer.Int(C_CONN_TIMEOUT)
	p.callTimeout, _ = b.Conf.Configer.Int(C_CALL_TIMEOUT)
	p.payloadsize, _ = b.Conf.Configer.Int(C_PAYLOADSIZE)
	p.debug = b.Conf.Debug
	p.ctx, p.cancel = context.WithCancel(context.Background())

	err := storageCheckHds(b.hdisk)
	if err != nil {
		glog.Fatalf(MODULE_NAME+" starage start error, bad data dir %v\n", err)
	}

	b.storageIoTaskCh = make([]chan *ioTask, len(b.hdisk))
	for i, _ := range b.hdisk {
		b.storageIoTaskCh[i] = make(chan *ioTask, 16)
	}

	if !b.Conf.Migrate.Disabled {
		b.storageMigrateConsistent.NumberOfReplicas = falcon.REPLICAS

		for node, addr := range b.Conf.Migrate.Upstream {
			b.storageMigrateConsistent.Add(node)
			b.storageNetTaskCh[node] = make(chan *netTask, 16)
			p.storageMigrateClients[node] = make([]BackendClient, p.workerProcesses)

			for i := 0; i < p.workerProcesses; i++ {
				conn, err := grpc.DialContext(p.ctx, addr,
					grpc.WithInsecure(), grpc.WithDialer(falcon.Dialer),
					grpc.WithBlock())
				if err != nil {
					glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n",
						node, addr, err)
				}
				p.storageMigrateClients[node][i] = NewBackendClient(conn)

				go p.netWorker(i, b.storageNetTaskCh[node],
					p.storageMigrateClients[node][i], addr)
			}
		}
	}

	for i, _ := range b.hdisk {
		go p.ioWorker(b.storageIoTaskCh[i])
	}

	go p.commitCacheWorker()

	return nil
}

func (p *StorageModule) stop(backend *Backend) error {
	p.cancel()
	return nil
}

func (p *StorageModule) reload(backend *Backend) error {
	// TODO
	p.stop(backend)
	time.Sleep(time.Second)
	p.prestart(backend)
	return p.start(backend)
}
