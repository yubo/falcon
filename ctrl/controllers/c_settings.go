/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"fmt"

	"github.com/yubo/falcon/ctrl/models"
)

// Operations about porfile/config/info
type SetController struct {
	BaseController
}

// @Title Get config
// @Description get tag role user
// @Param	module	path	string	true	"module  number"
// @Success 200 {object} models.ConfigEntry
// @Failure 403 string error
// @router /config/:module [get]
func (c *SetController) GetConfig() {
	var err error

	module := c.GetString(":module")
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	conf, err := me.ConfigGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title update config
// @Description get tag role user
// @Param	module	path	string	true	"module"
// @Success 200 string success
// @Failure 403 string error
// @router /config/:module [put]
func (c *SetController) UpdateConfig() {
	module := c.GetString(":module")

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if err := me.ConfigSet(module, c.Ctx.Input.RequestBody); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}

// @Title Get config
// @Description get tag role user
// @Param	action	path	string	true	"action"
// @Success 200 obj  result
// @Failure 403 string error
// @router /debug/:action [get]
func (c *SetController) GetDebugAction() {
	var err error
	var obj interface{}
	action := c.GetString(":action")

	switch action {
	case "populate":
		obj, err = models.Populate()
	case "reset_db":
		obj, err = models.ResetDb()
	default:
		err = fmt.Errorf("%s %s", models.ErrUnsupported.Error(), action)
	}

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}
