/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"container/list"
	"runtime"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appConfig     HandoffOpts = defaultOptions
	appConfigFile string
	appUpdateChan chan *[]*specs.MetaData // upstreams
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

	// upstreams
	appUpdateChan = make(chan *[]*specs.MetaData, 16)

}

func Handle(arg interface{}) {

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
		rpcStart(appConfig, appProcess)
		httpStart(appConfig, appProcess)
		upstreamStart(appConfig, appProcess)
		statStart(appConfig, appProcess)

		appProcess.StartSignal()
	} else {
		glog.Fatal(specs.ErrUnsupported)
	}
}
