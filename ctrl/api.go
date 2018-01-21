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
)

type ApiModule struct {
	Dev   bool
	beego beego.BeegoModule
}

func (p *ApiModule) PreStart(ctrl *Ctrl) error {
	return nil
}

func (p *ApiModule) Start(ctrl *Ctrl) (err error) {
	conf := &ctrl.Conf.Ctrl
	dsn := conf.Str(C_DSN)
	gc_time, _ := conf.Int64(C_SESSION_GC_MAX_LIFETIME)
	cookie_time, _ := conf.Int(C_SESSION_COOKIE_LIFETIME)
	http_addr := conf.Str(C_HTTP_ADDR)

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
	beego.BConfig.AppName = ctrl.Conf.Name
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = gc_time
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = cookie_time
	if RunMode&CTL_RUNMODE_DEV != 0 {
		orm.Debug = true
		beego.BConfig.RunMode = "dev"
		beego.BConfig.WebConfig.EnableDocs = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/v1.0/doc"] = "swagger"
	} else {
		beego.BConfig.RunMode = "prod"
	}
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

func (p *ApiModule) Stop(ctrl *Ctrl) error {
	p.beego.Stop()
	return nil
}

func (p *ApiModule) Reload(ctrl *Ctrl) error {
	return nil
	// TODO
	//p.Stop(c)
	//time.Sleep(time.Second)
	//p.PreStart(c)
	//return p.Start(c)
}
