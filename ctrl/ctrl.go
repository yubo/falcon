/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	initHooks = make([]hookfunc, 0)
	Config    *falcon.ConfCtrl
)

type hookfunc func(conf *falcon.ConfCtrl) error

type Ctrl struct {
	Conf falcon.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
}

func RegisterInit(fn hookfunc) {
	initHooks = append(initHooks, fn)
}

func (p Ctrl) Desc() string {
	return fmt.Sprintf("%s", p.Conf.Name)
}

func (p Ctrl) String() string {
	return fmt.Sprintf("%s (\n%s\n)",
		"conf", falcon.IndentLines(1, p.Conf.String()))
}

// ugly hack
// should called by main package
func (p *Ctrl) Init() error {
	Config = &p.Conf

	for _, fn := range initHooks {
		if err := fn(Config); err != nil {
			panic(err)
		}
	}
	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Name)

	p.status = falcon.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	// TODO: move configurations from conf/app.conf to falcon.conf(yyparse)

	// p.rpcStart()
	// p.httpStart()
	p.statStart()
	go beego.Run()
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	// p.httpStop()
	// p.rpcStop()
	return nil
}

func (p *Ctrl) Reload() error {
	glog.V(3).Infof(MODULE_NAME+"%s Reload", p.Conf.Name)
	return nil
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return nil
}
