/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"container/list"
	"net/rpc"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appConfig     BackendOpts = defaultOptions
	appConfigfile string
	appProcess    *specs.Process
	appTs         int64
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
	storageSyncEvent = make(chan specs.ProcEvent)
	storageMigrateConsistent = consistent.New()
	storageNetTaskCh = make(map[string]chan *netTask)
	storageMigrateClients = make(map[string][]*rpc.Client)

	// store
	size := CACHE_TIME / FLUSH_DISK_STEP
	if size < 0 {
		glog.Fatalf("store.init, bad size %d\n", size)
	}
}

func timeStart(config BackendOpts, ts *int64) {
	start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-ticker:
				now := time.Now().Unix()
				if config.Debug > 1 {
					atomic.StoreInt64(ts,
						start+(now-start)*DEBUG_MULTIPLES)
				} else {
					atomic.StoreInt64(ts, now)
				}
			}
		}
	}()
}

func timeNow() int64 {
	return atomic.LoadInt64(&appTs)
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
		timeStart(appConfig, &appTs)
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
