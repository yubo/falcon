/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

const (
	MODULE_NAME = "\x1B[32m[LB]\x1B[0m "
)

type Lb struct {
	Params   specs.ModuleParams
	Batch    int
	Backends []specs.Backend
	// runtime
	status        uint32
	running       chan struct{}
	rpcListener   *net.TCPListener
	httpListener  *net.TCPListener
	httpMux       *http.ServeMux
	bs            []*backend
	appUpdateChan chan *[]*specs.MetaData // upstreams
	rpcConnects   connList
}

func (p Lb) Desc() string {
	return fmt.Sprintf("%s", p.Params.Name)
}

func (p Lb) String() string {
	var s string
	for _, v := range p.Backends {
		s += fmt.Sprintf("\n%s", specs.IndentLines(1, v.String()))
	}
	if s != "" {
		s = fmt.Sprintf("\n%s\n", specs.IndentLines(1, s))
	}
	return fmt.Sprintf("%s (\n%s\n)\n"+
		"%-17s %d\n"+
		"%s (%s)",
		"params", specs.IndentLines(1, p.Params.String()),
		"batch", p.Batch,
		"backends", s)
}

func (p *Lb) Init() error {
	glog.V(3).Infof(MODULE_NAME+"%s Init()", p.Params.Name)

	// rpc
	p.rpcConnects = connList{list: list.New()}

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	p.status = specs.APP_STATUS_INIT
	return nil

}

func (p *Lb) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Params.Name)
	p.status = specs.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)
	p.statStart()
	p.rpcStart()
	p.httpStart()
	p.upstreamStart()
	p.status = specs.APP_STATUS_RUNING
	return nil
}

func (p *Lb) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Params.Name)
	p.status = specs.APP_STATUS_EXIT
	close(p.running)
	p.upstreamStop()
	p.httpStop()
	p.rpcStop()
	p.statStop()
	return nil
}

func (p *Lb) Reload() error {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Params.Name)
	return nil

}

func (p *Lb) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Params.Name, sig)
	return nil
}
