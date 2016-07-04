/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"container/list"
	"runtime"
	"sync/atomic"

	"github.com/yubo/falcon/specs"
)

var (
	appConfig     HandoffOpts = defaultOptions
	appEvents     []*specs.RoutineEvent
	appStatus     uint32
	appConfigFile string
	appUpdateChan chan *[]*specs.MetaData // upstreams
)

func init() {
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// rpc
	rpcEvent = &specs.RoutineEvent{Name: "rpc", E: make(chan specs.REvent)}
	rpcConnects = connList{list: list.New()}

	// http
	httpEvent = &specs.RoutineEvent{Name: "http", E: make(chan specs.REvent)}
	//httpRoutes()

	// upstreams
	appUpdateChan = make(chan *[]*specs.MetaData, 16)

}

func Handle(arg interface{}) {

	atomic.StoreUint32(&appStatus, specs.APP_STATUS_PENDING)
	parse(&appConfig, arg.(*specs.CmdOpts).ConfigFile)

	rpcStart(appConfig)
	httpStart(appConfig)
	upstreamStart(appConfig)

	atomic.StoreUint32(&appStatus, specs.APP_STATUS_RUNING)
	startSignal()
}
