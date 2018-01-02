/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/models"
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
	var err error
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigerGet("ctrl")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	// config filter
	ret := map[string]interface{}{
		ctrl.C_AUTH_MODULE: conf.Str(ctrl.C_AUTH_MODULE),
		ctrl.C_MASTER_MODE: conf.DefaultBool(ctrl.C_MASTER_MODE, false),
		ctrl.C_DEV_MODE:    conf.DefaultBool(ctrl.C_DEV_MODE, false),
		ctrl.C_MI_MODE:     conf.DefaultBool(ctrl.C_MI_MODE, false),
		ctrl.C_TAG_SCHEMA:  conf.Str(ctrl.C_TAG_SCHEMA),
	}

	c.SendMsg(200, ret)
}
