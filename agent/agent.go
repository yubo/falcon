/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	appAgent Agent = DefaultAgent
)

type Lb struct {
	Batch       int
	ConnTimeout int
	CallTimeout int
	Upstreams   []string
}

func (p Lb) String() string {
	return fmt.Sprintf("%-14s %d\n"+
		"%-14s %d\n%-14s %d\n"+
		"%-14s %s ",
		"batch", p.Batch,
		"conntimeout", p.ConnTimeout, "callTimeout", p.CallTimeout,
		"upstreams", strings.Join(p.Upstreams, ", "))
}

type Agent struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Rpc      bool
	RpcAddr  string
	Http     bool
	HttpAddr string
	IfPre    []string
	Interval int
	Lb  Lb
	// runtime
	appUpdateChan chan *[]*specs.MetaData // upstreams
	httpListener  *net.TCPListener
	httpMux       *http.ServeMux
	running       chan struct{}
}

func (p Agent) Desc() string {
	if p.Disabled {
		return fmt.Sprintf("%s(Disabled)", p.Name)
	} else {
		return fmt.Sprintf("%s", p.Name)
	}

}

func (p Agent) String() string {
	http := p.HttpAddr
	rpc := p.RpcAddr

	if !p.Http {
		http += "(disabled)"
	}
	if !p.Rpc {
		rpc += "(disabled)"
	}

	return fmt.Sprintf("%-17s %s\n%-17s %s\n"+
		"%-17s %d\n%-17s %v\n"+
		"%-17s %s\n%-17s %s\n"+
		"%-17s %s\n%-17s %d\n"+
		"%s (\n%s\n)",
		"Name", p.Name, "Host", p.Host,
		"debug", p.Debug, "disabled", p.Disabled,
		"http", http, "rpc", rpc,
		"ifprefix", strings.Join(p.IfPre, ", "),
		"interval", p.Interval,
		"Lb", specs.IndentLines(1, p.Lb.String()))
}

func (p *Agent) Init() error {
	glog.V(3).Infof("%s Init()", p.Name)
	// core
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	return nil
}

func (p *Agent) Start() error {
	glog.V(3).Infof("%s Start()", p.Name)
	p.running = make(chan struct{}, 0)
	p.statStart()
	p.upstreamStart()
	p.httpStart()
	p.collectStart()
	return nil
}

func (p *Agent) Stop() error {
	glog.V(3).Infof("%s Stop()", p.Name)
	close(p.running)
	p.collectStop()
	p.httpStop()
	p.upstreamStop()
	p.statStop()
	return nil
}

func (p *Agent) Reload() error {
	glog.V(3).Infof("%s Reload()", p.Name)
	return nil
}

func (p *Agent) Signal(sig os.Signal) error {
	glog.V(3).Infof("%s signal %v", p.Name, sig)
	return nil
}

/*
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
		statStart(appConfig, appProcess)

		appProcess.StartSignal()
	} else {
		glog.Fatal(specs.ErrUnsupported)
	}
}
*/
