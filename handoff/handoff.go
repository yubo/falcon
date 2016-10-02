/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"container/list"
	"fmt"
	"net"
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

type Backend struct {
	Disabled    bool
	Name        string
	Type        string
	Sched       string
	Batch       int
	ConnTimeout int
	CallTimeout int
	Upstreams   map[string]string
}

func (p Backend) String() string {
	var (
		ret  string
		name string
	)

	if p.Disabled {
		name = fmt.Sprintf("%s %s(Disable)", p.Type, p.Name)
	} else {
		name = fmt.Sprintf("%s %s", p.Type, p.Name)
	}

	for k, v := range p.Upstreams {
		ret += fmt.Sprintf("\n%-14s %s", k, v)
	}
	ret = fmt.Sprintf("%-14s %s\n"+
		"%-14s %d\n"+
		"%-14s %d\n"+
		"%-14s %d\n"+
		"%s (%s\n)",
		"sched", p.Sched,
		"batch", p.Batch,
		"conntimeout", p.ConnTimeout,
		"callTimeout", p.CallTimeout,
		"upstreams", specs.IndentLines(1, ret))
	return fmt.Sprintf("%s (\n%s\n)", name, specs.IndentLines(1, ret))
}

type Handoff struct {
	Debug       int
	Disabled    bool
	Name        string
	Http        bool
	HttpAddr    string
	Rpc         bool
	RpcAddr     string
	Replicas    int
	Concurrency int
	Backends    []Backend

	// runtime
	status        uint32
	running       chan struct{}
	rpcListener   *net.TCPListener
	httpListener  *net.TCPListener
	bs            []*backend
	appUpdateChan chan *[]*specs.MetaData // upstreams
	rpcConnects   connList
}

func (p Handoff) Desc() string {
	if p.Disabled {
		return fmt.Sprintf("%s(Disabled)", p.Name)
	} else {
		return fmt.Sprintf("%s", p.Name)
	}

}

func (p Handoff) String() string {
	http := p.HttpAddr
	rpc := p.RpcAddr

	if !p.Http {
		http += "(disabled)"
	}
	if !p.Rpc {
		rpc += "(disabled)"
	}

	ret := fmt.Sprintf("%-17s %d\n%-17s %v\n"+
		"%-17s %s\n%-17s %s\n"+
		"%-17s %d\n%-17s %d\nbackends (",
		"debug", p.Debug, "disabled", p.Disabled,
		"http", http, "rpc", rpc,
		"replicas", p.Replicas, "concurrency", p.Concurrency)
	for _, v := range p.Backends {
		ret += fmt.Sprintf("\n%s", specs.IndentLines(1, v.String()))
	}
	return ret + fmt.Sprintf("\n)")
}

func (p *Handoff) Init() error {
	glog.V(3).Infof("%s Init()", p.Name)

	// rpc
	p.rpcConnects = connList{list: list.New()}

	p.status = specs.APP_STATUS_INIT
	return nil

}

func (p *Handoff) Start() error {
	glog.V(3).Infof("%s Start()", p.Name)
	p.status = specs.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)
	p.statStart()
	p.rpcStart()
	p.httpStart()
	p.upstreamStart()
	p.status = specs.APP_STATUS_RUNING
	return nil
}

func (p *Handoff) Stop() error {
	glog.V(3).Infof("%s Stop()", p.Name)
	p.status = specs.APP_STATUS_EXIT
	close(p.running)
	p.upstreamStop()
	p.httpStop()
	p.rpcStop()
	p.statStop()
	return nil
}

func (p *Handoff) Reload() error {
	glog.V(3).Infof("%s Reload()", p.Name)
	return nil

}

func (p *Handoff) Signal(sig os.Signal) error {
	glog.V(3).Infof("%s signal %v", p.Name, sig)
	return nil
}

/*
func init() {
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// rpc
	rpcEvent = make(chan specs.ProcEvent)
	rpcConnects = connList{list: list.New()}

	// http
	httpEvent = make(chan specs.ProcEvent)
	httpRoutes()

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
*/
