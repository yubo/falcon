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
	"strings"

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
	MODULE_NAME = "\x1B[32m[CTRL]\x1B[0m "
)

type lb struct {
	name     string
	host     string
	addr     string
	last     int64
	nb_agent int
}

type backend struct {
	name string
	host string
	addr string
	last int64
}

type Ctrl struct {
	Params     specs.ModuleParams
	ConfigFile string
	Migrate    specs.Migrate
	Backends   []specs.Backend
	Lbs        []string

	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
	rpcConnects  connList
}

/*
func list_lb_entry(l *list.ListHead) *lb {
	return (*lb)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*lb)(nil)).list))))
}

func list_backend_entry(l *list.ListHead) *backend {
	return (*backend)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*backend)(nil)).list))))
}
*/

func (p Ctrl) Desc() string {
	return fmt.Sprintf("%s", p.Params.Name)
}

func (p Ctrl) String() string {
	var s string
	for _, v := range p.Backends {
		s += fmt.Sprintf("%s\n", specs.IndentLines(1, v.String()))
	}
	if s != "" {
		s = fmt.Sprintf("\n%s\n", specs.IndentLines(1, s))
	}
	return fmt.Sprintf("%s (\n%s\n)\n"+
		"%s (\n%s\n)\n"+
		"%s (%s)\n"+
		"%-17s %s\n",
		"params", specs.IndentLines(1, p.Params.String()),
		"migrate", specs.IndentLines(1, p.Migrate.String()),
		"backends", s,
		"lbs", strings.Join(p.Lbs, ", "))
}

func (p *Ctrl) Init() error {
	glog.V(3).Infof(MODULE_NAME+"%s Init()", p.Params.Name)

	// rpc
	p.rpcConnects = connList{list: list.New()}

	// http
	p.httpMux = http.NewServeMux()
	p.httpRoutes()

	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Params.Name)

	p.status = specs.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	dsn := beego.AppConfig.String("mysqldsn")
	maxIdle, _ := beego.AppConfig.Int("mysqlmaxidle")
	maxConn, _ := beego.AppConfig.Int("mysqlmaxconn")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxConn)

	p.rpcStart()
	p.httpStart()
	p.statStart()
	go beego.Run()
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Params.Name)
	p.status = specs.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	p.httpStop()
	p.rpcStop()
	return nil
}

func (p *Ctrl) Reload() error {
	glog.V(3).Infof(MODULE_NAME+"%s Reload", p.Params.Name)
	return nil
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Params.Name, sig)
	return nil
}
