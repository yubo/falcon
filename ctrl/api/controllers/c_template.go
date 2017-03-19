/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Templates
type TemplateController struct {
	BaseController
}

// @Title CreateTemplate
// @Description create templates
// @Param	body	body 	models.TemplateAction	true	"body for template content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *TemplateController) CreateTemplate() {
	var ta models.TemplateAction
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &ta)

	id, err := op.AddAction(&ta.Action)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	ta.Template.ActionId = id
	id, err = op.AddTemplate(&ta.Template)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTemplatesCnt
// @Description get Templates number
// @Param   query     query   string  false    "template name"
// @Success 200 {object} models.Total total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TemplateController) GetTemplatesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTemplatesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTemplates
// @Description get all Templates
// @Param   query     query   string  false    "template name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.TemplateUi templates ui info
// @Failure 403 string error
// @router /search [get]
func (c *TemplateController) GetTemplates() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTemplates(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Get
// @Description get template by id
// @Param	id	path 	int	true		"template id"
// @Param	clone	query 	bool	false		"clone tid to new one"
// @Success 200 {object} models.TemplateAction template and action info
// @Failure 403 string error
// @router /:id [get]
func (c *TemplateController) GetTemplate() {
	var (
		o   *models.TemplateAction
		err error
	)
	id, _ := c.GetInt64(":id", 0)
	clone, _ := c.GetBool("clone", false)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if clone {
		o, err = op.CloneTemplate(id)
	} else {
		o, err = op.GetTemplate(id)
	}
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, o)
	}
}

// @Title UpdateTemplate
// @Description update the template
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Template	true		"body for template content"
// @Success 200 {object} models.Template template info
// @Failure 403 string error
// @router /:id [put]
func (c *TemplateController) UpdateTemplate() {
	var ta models.TemplateAction

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &ta)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if t, err := op.UpdateTemplate(id, &ta); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, t)
	}
}

// @Title DeleteTemplate
// @Description delete the template
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TemplateController) DeleteTemplate() {

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteTemplate(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	c.SendMsg(200, "delete success!")
}
