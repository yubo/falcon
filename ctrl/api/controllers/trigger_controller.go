/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Triggers
type TriggerController struct {
	BaseController
}

// @Title CreateTrigger
// @Description create triggers
// @Param	body	body 	models.TriggerApiAdd	true	"body for trigger content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *TriggerController) CreateTrigger() {
	var input models.TriggerApiAdd

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if input.ParentId > 0 {
		p, err := op.GetTrigger(input.ParentId)
		if err != nil {
			c.SendMsg(400, fmt.Sprintf("parent id %d not exist",
				input.ParentId))
			return
		} else if p.ParentId > 0 {
			c.SendMsg(400, fmt.Sprintf("trigger %d can not "+
				"be a parent node", input.ParentId))
			return
		}
	} else {
	}

	trigger := models.Trigger{
		ParentId: input.ParentId,
		Priority: input.Priority,
		Name:     input.Name,
		Metric:   input.Metric,
		Tags:     input.Tags,
		Func:     input.Func,
		Op:       input.Op,
		Value:    input.Value,
		Msg:      input.Msg,
	}

	if id, err := op.CreateTrigger(&trigger); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTriggersCnt
// @Description get Triggers number
// @Param   query	query   string  false    "trigger name"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /cnt [get]
func (c *TriggerController) GetTriggersCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTriggersCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTriggers
// @Description get all Triggers
// @Param   query     query   string  false    "trigger name"
// @Param   limit     query   int     false    "limit page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.TriggerUi triggers ui info
// @Failure 400 string error
// @router /search [get]
func (c *TriggerController) GetTriggers() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTriggers(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title UpdateTrigger
// @Description update the trigger
// @Param	id	path 	string			true	"The id you want to update"
// @Param	body	body 	models.TriggerAction	true	"body for trigger content"
// @Success 200 {object} models.Trigger trigger info
// @Failure 400 string error
// @router /:id [put]
func (c *TriggerController) UpdateTrigger() {
}

// @Title DeleteTrigger
// @Description delete the trigger
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:400, msg:string}
// @router /:id [delete]
func (c *TriggerController) DeleteTrigger() {

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	_, err = op.DeleteTrigger(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	c.SendMsg(200, "delete success!")
}
