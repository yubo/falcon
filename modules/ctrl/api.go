/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon/modules/ctrl/api/routers"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yubo/falcon/modules/ctrl/api/models/auth"
	_ "github.com/yubo/falcon/modules/ctrl/api/models/session"
)

type apiModule struct {
	Dev   bool
	beego beego.BeegoModule
}

func (p *apiModule) PreStart(c *Ctrl) error {

	if c.Conf.BeegoDevMode {
		orm.Debug = true
		beego.BConfig.RunMode = beego.DEV
		beego.BConfig.WebConfig.EnableDocs = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/v1/doc"] = "swagger"
	} else {
		beego.BConfig.RunMode = beego.PROD
	}
	routers.Init(c.Conf.MiMode)
	return nil
}

func (p *apiModule) Start(c *Ctrl) (err error) {
	conf := c.Conf

	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "falcon"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = conf.Dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"
	beego.BConfig.AppName = "falcon"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = int64(conf.SessionGcMaxLifetime)
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = conf.SessionCookieLifetime
	if addr := strings.Split(conf.HttpAddr, ":"); len(addr) == 2 {
		beego.BConfig.Listen.HTTPAddr = addr[0]
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[1])
	} else if len(addr) == 1 {
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[0])
	}

	glog.V(3).Infof("port %s dev %v", conf.HttpAddr, conf.BeeMode)
	if err := p.beego.Start(conf.BeeMode); err != nil {
		return err
	}

	return nil
}

func (p *apiModule) Stop(c *Ctrl) error {
	p.beego.Stop()
	return nil
}

func (p *apiModule) Reload(ctrl *Ctrl) error {
	return nil
	// TODO
	//p.Stop(c)
	//time.Sleep(time.Second)
	//p.PreStart(c)
	//return p.Start(c)
}
