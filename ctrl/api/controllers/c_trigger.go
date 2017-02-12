/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Triggers
type TriggerController struct {
	BaseController
}

// @Title CreateTrigger
// @Description create triggers
// @Param	body	body 	models.Trigger	true	"body for trigger content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *TriggerController) CreateTrigger() {
	var trigger models.Trigger
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &trigger)
	trigger.Id = 0

	id, err := me.AddTrigger(&trigger)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetTriggersCnt
// @Description get Triggers number
// @Param   query     query   string  false    "trigger name"
// @Success 200  {total:int} trigger total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TriggerController) GetTriggersCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTriggersCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTriggers
// @Description get all Triggers
// @Param   query     query   string  false    "trigger name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} models.Trigger
// @Failure 403 error string
// @router /search [get]
func (c *TriggerController) GetTriggers() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	triggers, err := me.GetTriggers(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, triggers)
	}
}

// @Title Get
// @Description get trigger by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Trigger
// @Failure 403 error string
// @router /:id [get]
func (c *TriggerController) GetTrigger() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		trigger, err := me.GetTrigger(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, trigger)
		}
	}
}

// @Title UpdateTrigger
// @Description update the trigger
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Trigger	true		"body for trigger content"
// @Success 200 {object} models.Trigger
// @Failure 403 error string
// @router /:id [put]
func (c *TriggerController) UpdateTrigger() {
	var trigger models.Trigger

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &trigger)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateTrigger(id, &trigger); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteTrigger
// @Description delete the trigger
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TriggerController) DeleteTrigger() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteTrigger(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
