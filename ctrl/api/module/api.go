/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package module

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/routers"
)

type ApiModule struct {
	Dev   bool
	beego beego.BeegoModule
}

func (p *ApiModule) PreStart(c *ctrl.Ctrl) error {
	conf := &c.Conf.Ctrl

	if conf.DefaultBool(ctrl.C_BEEGODEV_MODE, false) {
		orm.Debug = true
		beego.BConfig.RunMode = beego.DEV
		beego.BConfig.WebConfig.EnableDocs = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/v1.0/doc"] = "swagger"
	} else {
		beego.BConfig.RunMode = beego.PROD
	}
	routers.Init()
	return nil
}

func (p *ApiModule) Start(c *ctrl.Ctrl) (err error) {
	conf := &c.Conf.Ctrl
	dsn := conf.Str(ctrl.C_DSN)
	gc_time, _ := conf.Int64(ctrl.C_SESSION_GC_MAX_LIFETIME)
	cookie_time, _ := conf.Int(ctrl.C_SESSION_COOKIE_LIFETIME)
	http_addr := conf.Str(ctrl.C_HTTP_ADDR)

	// ctrl beggo config
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "falcon"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"
	beego.BConfig.AppName = c.Conf.Name
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = gc_time
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = cookie_time
	if addr := strings.Split(http_addr, ":"); len(addr) == 2 {
		beego.BConfig.Listen.HTTPAddr = addr[0]
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[1])
	} else if len(addr) == 1 {
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[0])
	}

	if err := p.beego.Start(p.Dev); err != nil {
		return err
	}

	return nil
}

func (p *ApiModule) Stop(c *ctrl.Ctrl) error {
	p.beego.Stop()
	return nil
}

func (p *ApiModule) Reload(ctrl *ctrl.Ctrl) error {
	return nil
	// TODO
	//p.Stop(c)
	//time.Sleep(time.Second)
	//p.PreStart(c)
	//return p.Start(c)
}
