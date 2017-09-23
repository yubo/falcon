/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package api

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/config"

	_ "github.com/yubo/falcon/ctrl/api/models/auth"
	_ "github.com/yubo/falcon/ctrl/api/routers"
)

type ApiModule struct {
	b beego.BeegoModule
}

func (p *ApiModule) PreStart(c *config.ConfCtrl) error {
	if err := models.InitModels(c); err != nil {
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
