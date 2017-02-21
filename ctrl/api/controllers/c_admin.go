/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about porfile/config/info
type AdminController struct {
	BaseController
}

// @Title Get config
// @Description get tag role user
// @Param	module	path	string	true	"module  number"
// @Success 200 {object} models.ConfigEntry {defualt{}, conf{}, configfile{}}
// @Failure 403 string error
// @router /config/:module [get]
func (c *AdminController) GetConfig() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title update config
// @Description get tag role user
// @Param	module	path	string	true	"module"
// @Param	body	body	map[string]string	true	""
// @Success 200 {string} success
// @Failure 403 string error
// @router /config/:module [put]
func (c *AdminController) UpdateConfig() {
	var conf map[string]string

	module := c.GetString(":module")

	beego.Debug(string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &conf)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if err := op.ConfigSet(module, conf); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}

// @Title Get config
// @Description get tag role user
// @Param	action	path	string	true	"action"
// @Success 200 {string} result
// @Failure 403 string error
// @router /debug/:action [get]
func (c *AdminController) GetDebugAction() {
	var err error
	var obj interface{}
	action := c.GetString(":action")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	switch action {
	case "populate":
		obj, err = op.Populate()
	case "reset_db":
		obj, err = op.ResetDb()
	default:
		err = fmt.Errorf("%s %s", models.ErrUnsupported.Error(), action)
	}

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}
