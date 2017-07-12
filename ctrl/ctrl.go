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
	"reflect"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/config"
	"github.com/yubo/falcon/ctrl/parse"
	"github.com/yubo/falcon/utils"
)

type hookfunc func(conf *config.ConfCtrl) error

type Ctrl struct {
	Conf    *config.ConfCtrl
	oldConf *config.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
}

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	prestartHooks = make([]hookfunc, 0)
	reloadHooks   = make([]hookfunc, 0)
	Configure     *config.ConfCtrl
	EtcdCli       *utils.EtcdCli
)

func RegisterPrestart(fn hookfunc) {
	prestartHooks = append(prestartHooks, fn)
}

func RegisterReload(fn hookfunc) {
	reloadHooks = append(reloadHooks, fn)
}

func (p *Ctrl) New(conf interface{}) falcon.Module {
	return &Ctrl{Conf: conf.(*config.ConfCtrl)}
}

func (p *Ctrl) Name() string {
	return fmt.Sprintf("%s", p.Conf.Name)
}

func (p *Ctrl) Parse(text []byte, filename string, lino int, debug bool) interface{} {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfCtrl)
	return p.Conf
}

func (p *Ctrl) String() string {
	return p.Conf.String()
}

// ugly hack
// should called by main package
func (p *Ctrl) Prestart() error {
	glog.V(3).Infof(MODULE_NAME + "Prestart() entering")
	Configure = p.Conf

	EtcdCli = utils.NewEtcdCli(Configure.Ctrl)

	EtcdCli.Prestart()
	for _, fn := range prestartHooks {
		glog.V(3).Infof(MODULE_NAME+"%s() called\n", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}
	glog.V(3).Infof(MODULE_NAME + "Prestart() leaving")
	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME + "Start() entering")

	p.status = falcon.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	EtcdCli.Start()
	// p.rpcStart()
	// p.httpStart()
	p.statStart()
	go beego.Run()
	glog.V(3).Infof(MODULE_NAME + "Start() leaving")
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	EtcdCli.Stop()
	// p.httpStop()
	// p.rpcStop()
	return nil
}

// TODO: reload is not yet implemented
func (p *Ctrl) Reload(conf interface{}) error {
	p.Conf = conf.(*config.ConfCtrl)
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	Configure = p.Conf

	EtcdCli.Reload(Configure.Ctrl)
	for _, fn := range prestartHooks {
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}

	return nil
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return nil
}
