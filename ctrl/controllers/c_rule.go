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
	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Rules
type RuleController struct {
	BaseController
}

// @Title CreateRule
// @Description create rules
// @Param	body	body 	models.Rule	true	"body for rule content"
// @Success {code:200, data:int} models.Rule.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *RuleController) CreateRule() {
	var rule models.Rule
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &rule)
	id, err := me.AddRule(&rule)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, id)
	}
}

// @Title GetRulesCnt
// @Description get Rules number
// @Param   query     query   string  false    "rule name"
// @Success {code:200, data:int} rule number
// @Failure {code:int, msg:string}
// @router /cnt [get]
func (c *RuleController) GetRulesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetRulesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title GetRules
// @Description get all Rules
// @Param   query     query   string  false    "rule name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success {code:200, data:object} models.Rule
// @Failure {code:int, msg:string}
// @router /search [get]
func (c *RuleController) GetRules() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	rules, err := me.GetRules(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, rules)
	}
}

// @Title Get
// @Description get rule by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Rule
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *RuleController) GetRule() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		rule, err := me.GetRule(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, rule)
		}
	}
}

// @Title UpdateRule
// @Description update the rule
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Rule	true		"body for rule content"
// @Success {code:200, data:object} models.Rule
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *RuleController) UpdateRule() {
	var rule models.Rule

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &rule)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateRule(id, &rule); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteRule
// @Description delete the rule
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *RuleController) DeleteRule() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteRule(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetRule() {
	var rules []*models.Rule

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryRules(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per, c.SetPaginator(per, total).Offset()).All(&rules)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Template")
	c.Data["Rules"] = rules
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "rule name"}

	c.TplName = "rule/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditRule() {
	var rule *models.Rule
	var me *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	rule, err = me.GetRule(id)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Template")
	c.Data["Rule"] = rule
	c.Data["H1"] = "edit template"
	c.Data["Method"] = "put"
	c.TplName = "rule/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddRule() {

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Template")
	c.Data["Method"] = "post"
	c.Data["H1"] = "add template"
	c.TplName = "rule/edit.tpl"
}
