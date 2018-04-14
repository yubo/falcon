/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/ctrl/api/models"
)

// Operations about porfile/config/info
type AdminController struct {
	BaseController
}

// @Title Get config
// @Description list etcd Description
// @Success 200 map[string]string map[string]string etcd map string
// @Failure 400 string error
// @router /config/list/etcd [get]
func (c *AdminController) GetEtcdMap() {
	c.SendMsg(200, models.EtcdMap)
}

// @Title Get config
// @Description list module config api Description
// @Success 200 map[string]string map[string]string etcd map string
// @Failure 400 string error
// @router /config/list/module [get]
func (c *AdminController) GetModuleMap() {
	c.SendMsg(200, models.ModuleMap)
}

// @Title Get online endpoint
// @Description get online endpoint
// @Param	module	path	string	true	"module name"
// @Success 200 map[string]string map[string]string online endpoint list
// @Failure 400 string error
// @router /online/:module [get]
func (c *AdminController) GetOnline() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.OnlineGet(module)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title Get config
// @Description get module config
// @Param	module	path	string	true	"module name"
// @Success 200 [3]map[string]string {defualt{}, conf{}, configfile{}}
// @Failure 400 string error
// @router /config/:module [get]
func (c *AdminController) GetConfig() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigGet(module)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title update config
// @Description get tag role user
// @Param	module	path	string	true	"module"
// @Param	body	body	models.Obj true	"map[string]string"
// @Success 200 {string} success
// @Failure 400 string error
// @router /config/:module [put]
func (c *AdminController) UpdateConfig() {
	var conf map[string]string

	module := c.GetString(":module")

	c.SendMsg(400, core.EACCES.Error())
	return

	beego.Debug(string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &conf)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if err := op.ConfigSet(module, conf); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}

// @Title Get config
// @Description get tag role user
// @Param	action	path	string	true	"action"
// @Success 200 {string} result
// @Failure 400 string error
// @router /debug/:action [get]
func (c *AdminController) GetDebugAction() {
	var err error
	var obj interface{}
	action := c.GetString(":action")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	switch action {
	case "populate":
		obj, err = op.ResetDb(true)
	case "reset_db":
		obj, err = op.ResetDb(false)
	default:
		err = fmt.Errorf("%s %s", core.ErrUnsupported.Error(), action)
	}

	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title Get expansion status
// @Description graph expansion config
// @Param	module	path	string	true	"module name"
// @Success 200 {object} models.ExpansionStatus expansion status
// @Failure 400 string error
// @router /expansion/:module [get]
func (c *AdminController) GetExpansion() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	obj, err := op.ExpansionGet(module)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title set expansion (add endpoint)
// @Description expansion set
// @Param	module	path	string	true	"module name"
// @Param	body	body 	models.ExpansionStatus	true		"body for endpoints content"
// @Success 200 string success
// @Failure 400 string error
// @router /expansion/:module [put]
func (c *AdminController) SetExpansion() {
	var err error
	var status models.ExpansionStatus
	var cur_status *models.ExpansionStatus

	module := c.GetString(":module")
	json.Unmarshal(c.Ctx.Input.RequestBody, &status)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	cur_status, err = op.ExpansionGet(module)
	if err != nil {
		goto out
	}

	if !cur_status.Migrating && status.Migrating {
		err = op.ExpansionBegin(module, status.NewEndpoint)
		goto out
	}

	if cur_status.Migrating && !status.Migrating {
		err = op.ExpansionFinish(module)
		goto out
	}

	err = core.ErrParam
out:
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}
