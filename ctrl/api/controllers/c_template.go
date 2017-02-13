/*
 * Copyright 2016 yubo. All rights reserved.
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
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *TemplateController) CreateTemplate() {
	var ta models.TemplateAction
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &ta)

	id, err := me.AddAction(&ta.Action)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	ta.Template.ActionId = id
	id, err = me.AddTemplate(&ta.Template)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetTemplatesCnt
// @Description get Templates number
// @Param   query     query   string  false    "template name"
// @Success 200  {total:int} template number
// @Failure 403 string error
// @router /cnt [get]
func (c *TemplateController) GetTemplatesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTemplatesCnt(query)
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
// @Success 200 [object] []models.TemplateUi
// @Failure 403 string error
// @router /search [get]
func (c *TemplateController) GetTemplates() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	ret, err := me.GetTemplates(query, per, offset)
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
// @Success 200 {object} models.TemplateAction
// @Failure 403 error string
// @router /:id [get]
func (c *TemplateController) GetTemplate() {
	var (
		o   *models.TemplateAction
		err error
	)
	id, _ := c.GetInt64(":id", 0)
	clone, _ := c.GetBool("clone", false)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if clone {
		o, err = me.CloneTemplate(id)
	} else {
		o, err = me.GetTemplate(id)
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
// @Success 200 {object} models.Template
// @Failure 403 error string
// @router /:id [put]
func (c *TemplateController) UpdateTemplate() {
	var ta models.TemplateAction

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &ta)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if t, err := me.UpdateTemplate(id, &ta); err != nil {
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

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteTemplate(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	c.SendMsg(200, "delete success!")
}
