/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"runtime"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appConfig     AgentOpts = DefaultOptions
	appConfigFile string
	appUpdateChan chan *[]*specs.MetaData // upstreams
	appProcess    *specs.Process
)

func init() {
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// http
	httpEvent = make(chan specs.ProcEvent)
	httpRoutes()

	// upstreams
	appUpdateChan = make(chan *[]*specs.MetaData, 16)

}

func Handle(arg interface{}) {

	opts := arg.(*specs.CmdOpts)

	//atomic.StoreUint32(&appStatus, specs.APP_STATUS_PENDING)
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
		httpStart(appConfig, appProcess)
		upstreamStart(appConfig, appProcess)
		collectStart(appConfig, appProcess)

		appProcess.StartSignal()
	} else {
		glog.Fatal(specs.ErrUnsupported)
	}
}
