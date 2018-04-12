/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/yubo/falcon/ctrl/api/models"
)

type modelsModule struct {
}

func (p *modelsModule) PreStart(c *Ctrl) error {
	return nil
}

func (p *modelsModule) Start(c *Ctrl) error {
	return models.Init(c.Conf.Configer, c.db,
		c.etcdCli, c.transferCli)
}

func (p *modelsModule) Stop(c *Ctrl) error {
	return nil
}

func (p *modelsModule) Reload(c *Ctrl) error {
	// TODO
	return nil
}
