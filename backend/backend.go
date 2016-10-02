/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

/*
#include "cache.h"
*/
import "C"

import (
	"container/list"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

const (
	CACHE_TIME              = 1800 //s
	FIRST_FLUSH_DISK        = 1    //s
	FLUSH_DISK_STEP         = 1    //s
	DEFAULT_HISTORY_SIZE    = 3
	CONN_RETRY              = 2
	CACHE_SIZE              = C.CACHE_SIZE     // must pow(2,n)
	CACHE_SIZE_MASK         = C.CACHE_SIZE - 1 //
	DATA_TIMESTAMP_REGULATE = true
	INDEX_QPS               = 100
	INDEX_UPDATE_CYCLE_TIME = 86400
	INDEX_TIMEOUT           = 86400
	INDEX_TRASH_LOOPTIME    = 600
	INDEX_MAX_OPEN_CONNS    = 4
	DEBUG_MULTIPLES         = 20    // demo 时间倍数
	DEBUG_STEP              = 60    //
	DEBUG_SAMPLE_NB         = 18000 //单周期生成样本数量
	DEBUG_STAT_MODULE       = ST_M_CACHE | ST_M_INDEX
	DEBUG_STAT_STEP         = 60
)

type Migrate struct {
	Disabled    bool
	Concurrency int
	Replicas    int
	CallTimeout int
	ConnTimeout int
	Upstreams   map[string]string
}

func (o Migrate) String() string {
	var upstream string
	indent := strings.Repeat(" ", specs.IndentSize)

	for k, v := range o.Upstreams {
		upstream += fmt.Sprintf("%s%-10s = %s\n", indent, k, v)
	}

	return fmt.Sprintf("%-12s %v\n%-12s %d\n"+
		"%-12s %d\n%-12s %d\n"+
		"%-12s %d\n%s (\n%s\n)",
		"disable", o.Disabled, "concurrency", o.Concurrency,
		"replicas", o.Replicas, "callTimeout", o.CallTimeout,
		"conntimeout", o.ConnTimeout, "upstream", strings.TrimRight(upstream, "\n"))
}

type Storage struct {
	Type   string
	Hdisks []string
}

func (o Storage) String() string {

	return fmt.Sprintf("%-12s %s\n%-12s %s",
		"type", o.Type,
		"hdisks", strings.Join(o.Hdisks, ", "))
}

type Shm struct {
	Magic uint32
	Key   int
	Size  int
}

func (o Shm) String() string {

	return fmt.Sprintf("%-14s 0x%x\n%-14s 0x%x\n%-14s %v",
		"magic_code", o.Magic,
		"key_start_id", o.Key,
		"segment_size", o.Size)

}

type Backend struct {
	// config
	Debug           int
	Disabled        bool
	Name            string
	Http            bool
	HttpAddr        string
	Rpc             bool
	RpcAddr         string
	Idx             bool
	IdxInterval     int
	IdxFullInterval int
	Dsn             string
	DbMaxIdle       int
	Migrate         Migrate
	Storage         Storage
	Shm             Shm
	// runtime
	status                   uint32
	running                  chan struct{}
	ts                       int64
	statTicker               chan time.Time
	timeTicker               chan time.Time
	commitTicker             chan time.Time
	rpcListener              *net.TCPListener
	rpcConnects              connList
	rpcBkd                   *Bkd
	httpListener             *net.TCPListener
	httpMux                  *http.ServeMux
	storageIoTaskCh          []chan *ioTask
	storageNetTaskCh         map[string]chan *netTask
	storageMigrateClients    map[string][]*rpc.Client
	storageMigrateConsistent *consistent.Consistent
	cache                    *backendCache
	indexDb                  *sql.DB
	indexUpdateCh            chan *cacheEntry
}

func (p Backend) Desc() string {
	if p.Disabled {
		return fmt.Sprintf("%s(Disabled)", p.Name)
	} else {
		return fmt.Sprintf("%s", p.Name)
	}

}
func (p Backend) String() string {
	http := p.HttpAddr
	rpc := p.RpcAddr

	if !p.Http {
		http += "(disabled)"
	}
	if !p.Rpc {
		rpc += "(disabled)"
	}
	return fmt.Sprintf("%-17s %d\n%-17s %v\n"+
		"%-17s %s\n%-17s %s\n"+
		"%-17s %v\n%-17s %d\n"+
		"%-17s %d\n%-17s %s\n"+
		"%-17s %d\n%s (\n%s\n)\n%s "+
		"(\n%s\n)\n%s (\n%s\n)",
		"debug", p.Debug,
		"disabled", p.Disabled,
		"http", http,
		"rpc", rpc,
		"idx", p.Idx,
		"idx_interval", p.IdxInterval,
		"idx_full_interval", p.IdxFullInterval,
		"dsn", p.Dsn,
		"dbmaxidle", p.DbMaxIdle,
		"migrate", specs.IndentLines(1, p.Migrate.String()),
		"storage", specs.IndentLines(1, p.Storage.String()),
		"shm", specs.IndentLines(1, p.Shm.String()))
}

func (p *Backend) Init() error {
	glog.V(3).Infof("%s Init()", p.Name)
	// core

	//cache
	p.rpcConnects = connList{list: list.New()}
	p.cacheInit()

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	// rrdtool/sync_disk/migrate
	p.storageNetTaskCh = make(map[string]chan *netTask)
	p.storageMigrateClients = make(map[string][]*rpc.Client)
	p.storageMigrateConsistent = consistent.New()

	// store
	size := CACHE_TIME / FLUSH_DISK_STEP
	if size < 0 {
		glog.Fatalf("store.init, bad size %d\n", size)
	}
	p.status = specs.APP_STATUS_INIT
	return nil

}

func (p *Backend) Start() error {
	glog.V(3).Infof("%s Start()", p.Name)
	p.status = specs.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)
	p.timeStart()
	p.rrdStart()
	p.rpcStart()
	p.indexStart()
	p.httpStart()
	p.statStart()
	p.cacheStart()
	p.status = specs.APP_STATUS_RUNING
	return nil
}

func (p *Backend) Stop() error {
	glog.V(3).Infof("%s Stop()", p.Name)
	p.status = specs.APP_STATUS_EXIT
	close(p.running)
	p.cacheStop()
	p.statStop()
	p.httpStop()
	p.indexStop()
	p.rpcStop()
	p.rrdStop()
	p.timeStop()
	return nil
}

func (p *Backend) Reload() error {
	glog.V(3).Infof("%s Reload()", p.Name)
	return nil
}

func (p *Backend) Signal(sig os.Signal) error {
	glog.V(3).Infof("%s signal %v", p.Name, sig)
	return nil
}

func (p *Backend) timeStart() {
	start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}

			case <-ticker:
				now := time.Now().Unix()
				if p.Debug > 1 {
					atomic.StoreInt64(&p.ts,
						start+(now-start)*DEBUG_MULTIPLES)
				} else {
					atomic.StoreInt64(&p.ts, now)
				}
			}
		}
	}()
}

func (p *Backend) timeStop() {
}

func (p *Backend) timeNow() int64 {
	return atomic.LoadInt64(&p.ts)
}

/*

func Handle(arg interface{}) {

	//atomic.StoreUint32(&appStatus, specs.APP_STATUS_PENDING)
	opts := arg.(*specs.CmdOpts)
	parse(&appConfig, opts.ConfigFile)
	appProcess = specs.NewProcess(appConfig.PidFile)

	cmd := "start"
	if len(opts.Args) > 0 {
		cmd = opts.Args[0]
	}

	if cmd == "stop" {

		if err := appProcess.Kill(syscall.SIGTERM); err != nil {
			glog.Fatal(err)
		}
	} else if cmd == "reload" {
		if err := appProcess.Kill(syscall.SIGUSR1); err != nil {
			glog.Fatal(err)
		}
	} else if cmd == "start" {
		if err := appProcess.Check(); err != nil {
			glog.Fatal(err)
		}
		if err := appProcess.Save(); err != nil {
			glog.Fatal(err)
		}
		rrdStart(appConfig, appProcess)
		rpcStart(appConfig, appProcess)
		indexStart(appConfig, appProcess)
		httpStart(appConfig, appProcess)
		statStart(appConfig, appProcess)
		cacheStart(appConfig, appProcess)

		appProcess.StartSignal()
	} else {
		glog.Fatal(specs.ErrUnsupported)
	}

}
*/
