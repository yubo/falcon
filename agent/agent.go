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

const (
	MODULE_NAME = "\x1B[32m[AGENT]\x1B[0m "
)

var (
	appAgent Agent = DefaultAgent
)

type Agent struct {
	Params    specs.ModuleParams
	Interval  int
	Batch     int
	IfPre     []string
	Upstreams []string
	// runtime
	appUpdateChan chan *[]*specs.MetaData
	httpListener  *net.TCPListener
	httpMux       *http.ServeMux
	running       chan struct{}
}

func (p Agent) Desc() string {
	return fmt.Sprintf("%s", p.Params.Name)
}

func (p Agent) String() string {
	return fmt.Sprintf("%s (\n%s\n)\n"+
		"%-17s %s\n"+
		"%-17s %d\n"+
		"%-17s %d\n"+
		"%-17s %s",
		"params", specs.IndentLines(1, p.Params.String()),
		"ifprefix", strings.Join(p.IfPre, ", "),
		"interval", p.Interval,
		"batch", p.Batch,
		"upstreams", strings.Join(p.Upstreams, ", "))
}

func (p *Agent) Init() error {
	glog.V(3).Infof(MODULE_NAME+"%s Init()", p.Params.Name)
	// core
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	return nil
}

func (p *Agent) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Params.Name)
	p.running = make(chan struct{}, 0)
	p.ctrlStart()
	p.statStart()
	p.upstreamStart()
	p.httpStart()
	p.collectStart()
	return nil
}

func (p *Agent) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Params.Name)
	close(p.running)
	p.collectStop()
	p.httpStop()
	p.upstreamStop()
	p.statStop()
	p.ctrlStop()
	return nil
}

func (p *Agent) Reload() error {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Params.Name)
	return nil
}

func (p *Agent) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Params.Name, sig)
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
