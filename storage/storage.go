/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"container/list"
	"net/rpc"
	"runtime"
	"sync/atomic"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appConfig     StorageOpts = defaultOptions
	appEvents     []*specs.RoutineEvent
	appStatus     uint32
	appConfigfile string
)

func init() {
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// rpc
	rpcEvent = &specs.RoutineEvent{Name: "rpc",
		E: make(chan specs.REvent)}
	rpcConnects = connList{list: list.New()}

	// http
	httpEvent = &specs.RoutineEvent{Name: "http",
		E: make(chan specs.REvent)}
	httpRoutes()

	// rrdtool/sync_disk/migrate
	rrdSyncEvent = &specs.RoutineEvent{Name: "sync",
		E: make(chan specs.REvent)}
	rrdIoTaskCh = make(chan *ioTask, 16)
	rrdMigrateConsistent = consistent.New()
	rrdNetTaskCh = make(map[string]chan *netTask)
	rrdMigrateClients = make(map[string][]*rpc.Client)

	// store
	size := CACHE_TIME / FLUSH_DISK_STEP
	if size < 0 {
		glog.Fatalf("store.init, bad size %d\n", size)
	}

	// cache
	cache.hash = make(map[string]*cacheEntry)
}

func Handle(arg interface{}) {

	atomic.StoreUint32(&appStatus, specs.APP_STATUS_PENDING)
	parse(&appConfig, arg.(*specs.CmdOpts).ConfigFile)

	dbStart(appConfig.Dsn, appConfig.DbMaxIdle)
	rrdStart(appConfig)
	rpcStart(appConfig)
	//indexStart()
	httpStart(appConfig)

	atomic.StoreUint32(&appStatus, specs.APP_STATUS_RUNING)
	startSignal()
}
