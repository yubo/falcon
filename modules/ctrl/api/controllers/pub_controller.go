/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"github.com/yubo/falcon/modules/ctrl/api/models"
	"github.com/yubo/falcon/modules/ctrl/stats"
)

// Operations about porfile/config/info
type PubController struct {
	BaseController
}

// @Title Get config
// @Description get ctrl modules config
// @Success 200  map[string]interface{} ctrl server config
// @Failure 400 string error
// @router /config/ctrl [get]
func (c *PubController) GetConfig() {
	c.SendMsg(200, models.GetConfig())
}

// @Title Get stats
// @Description get ctrl modules config
// @Success 200  map[string]interface{} ctrl server config
// @Failure 400 string error
// @router /stats [get]
func (c *PubController) GetStats() {
	c.SendMsg(200, stats.Gets())
}
