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
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
)

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

type Ctrl struct {
	Conf falcon.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
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

	beego.Debug(fmt.Sprintf(MODULE_NAME+"%s Init()", p.Conf.Name))
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = p.Conf.Dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"

	// connect db
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", p.Conf.Dsn, p.Conf.DbMaxIdle, p.Conf.DbMaxConn)
	orm.RegisterModelWithPrefix("",
		new(models.User), new(models.Host), new(models.Tag),
		new(models.Role), new(models.Token), new(models.Log),
		new(models.Tag_rel), new(models.Tpl_rel), new(models.Team),
		new(models.Template), new(models.Trigger), new(models.Expression),
		new(models.Action), new(models.Strategy))

	p.Conf.Agent.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["agent"])
	p.Conf.Lb.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["lb"])
	p.Conf.Backend.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["backend"])
	p.Conf.Ctrl.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["ctrl"])

	if conf, err := models.GetDbConfig("ctrl"); err == nil {
		p.Conf.Ctrl.Set(falcon.APP_CONF_DB, conf)
	}
	beego.Debug(fmt.Sprintf("%s Api %s", p.Conf.Name, p.Conf.Ctrl))

	// config -> beego config
	c := &p.Conf.Ctrl
	beego.BConfig.AppName = p.Conf.Name
	beego.BConfig.RunMode = c.Str(falcon.C_RUN_MODE)
	beego.BConfig.Listen.HTTPPort, _ = c.Int(falcon.C_HTTP_PORT)
	beego.BConfig.WebConfig.EnableDocs, _ = c.Bool(falcon.C_ENABLE_DOCS)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime, _ = c.Int64(falcon.C_SEESION_GC_MAX_LIFETIME)
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime, _ = c.Int(falcon.C_SESSION_COOKIE_LIFETIME)

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("orm debug on")
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	}

	models.Init(p.Conf)

	return nil

}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Name)

	p.status = falcon.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
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
