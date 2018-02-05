/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package module

import (
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/models"
)

type ModelsModule struct {
}

func (p *ModelsModule) PreStart(c *ctrl.Ctrl) error {
	return nil
}

func (p *ModelsModule) Start(c *ctrl.Ctrl) error {
	conf := &c.Conf.Ctrl
	models.Db = ctrl.Orm

	if err := models.InitConfig(conf); err != nil {
		return err
	}
	if err := models.InitAuth(conf); err != nil {
		return err
	}
	if err := models.InitCache(conf); err != nil {
		return err
	}

	models.SysOp = &models.Operator{
		O:     ctrl.Orm.Ctrl,
		Token: models.SYS_F_A_TOKEN | models.SYS_F_O_TOKEN | models.SYS_F_A_TOKEN,
	}
	models.SysOp.User, _ = models.GetUser(1, models.SysOp.O)

	models.PutEtcdConfig()

	return nil
}

func (p *ModelsModule) Stop(c *ctrl.Ctrl) error {
	return nil
}

func (p *ModelsModule) Reload(c *ctrl.Ctrl) error {
	// TODO
	return nil
}
