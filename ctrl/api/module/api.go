/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package module

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/api/routers"
	"github.com/yubo/falcon/ctrl/config"

	_ "github.com/yubo/falcon/ctrl/api/models/auth"
)

type ApiModule struct {
	b beego.BeegoModule
}

func (p *ApiModule) PreStart(c *config.ConfCtrl) error {
	if err := routers.PreStart(); err != nil {
		return err
	}
	if err := models.PreStart(c); err != nil {
		return err
	}
	return p.b.PreStart()
}

func (p *ApiModule) Start(c *config.ConfCtrl) error {
	return p.b.Start()
}

func (p *ApiModule) Stop(c *config.ConfCtrl) error {
	return p.b.Stop()
	return nil
}

func (p *ApiModule) Reload(old, c *config.ConfCtrl) error {
	p.Stop(c)
	time.Sleep(time.Second)
	p.PreStart(c)
	return p.Start(c)
}

type DevModule struct {
	b beego.BeegoModule
}

func (p *DevModule) PreStart(c *config.ConfCtrl) error {
	if err := routers.PreStart(); err != nil {
		return err
	}
	if err := models.PreStart(c); err != nil {
		return err
	}
	return nil
}

func (p *DevModule) Start(c *config.ConfCtrl) error {
	beego.Run()
	return nil
}

func (p *DevModule) Stop(c *config.ConfCtrl) error {
	return nil
}

func (p *DevModule) Reload(old, c *config.ConfCtrl) error {
	p.Stop(c)
	time.Sleep(time.Second)
	p.PreStart(c)
	return p.Start(c)
}
