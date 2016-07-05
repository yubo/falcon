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
	"syscall"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appConfig     StorageOpts = defaultOptions
	appConfigfile string
	appProcess    *specs.Process
)

func init() {
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// rpc
	rpcEvent = make(chan specs.ProcEvent)
	rpcConnects = connList{list: list.New()}

	// http
	httpEvent = make(chan specs.ProcEvent)
	httpRoutes()

	// rrdtool/sync_disk/migrate
	rrdSyncEvent = make(chan specs.ProcEvent)
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
		dbStart(appConfig.Dsn, appConfig.DbMaxIdle)
		rrdStart(appConfig, appProcess)
		rpcStart(appConfig, appProcess)
		//indexStart()
		httpStart(appConfig, appProcess)

		appProcess.StartSignal()
	} else {
		glog.Fatal(specs.ErrUnsupported)
	}

}
