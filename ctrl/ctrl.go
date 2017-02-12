/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/yubo/falcon/ctrl/models/auth"
	_ "github.com/yubo/falcon/ctrl/routers"
	"github.com/yubo/falcon/specs"
)

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

type Ctrl struct {
	Conf specs.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
	rpcConnects  connList
}

func (p Ctrl) Desc() string {
	return fmt.Sprintf("%s", p.Conf.Params.Name)
}

func (p Ctrl) String() string {
	return fmt.Sprintf("%s (\n%s\n)",
		"params", specs.IndentLines(1, p.Conf.Params.String()))
}

func (p *Ctrl) Init() error {
	glog.V(3).Infof(MODULE_NAME+"%s Init()", p.Conf.Params.Name)

	// rpc
	p.rpcConnects = connList{list: list.New()}

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Params.Name)

	p.status = specs.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	// TODO: move configurations from conf/app.conf to falcon.conf(yyparse)
	dsn := beego.AppConfig.String("mysqldsn")
	maxIdle, _ := beego.AppConfig.Int("mysqlmaxidle")
	maxConn, _ := beego.AppConfig.Int("mysqlmaxconn")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxConn)

	// p.rpcStart()
	p.httpStart()
	p.statStart()
	go beego.Run()
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Params.Name)
	p.status = specs.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	p.httpStop()
	// p.rpcStop()
	return nil
}

func (p *Ctrl) Reload() error {
	glog.V(3).Infof(MODULE_NAME+"%s Reload", p.Conf.Params.Name)
	return nil
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Params.Name, sig)
	return nil
}
