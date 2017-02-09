/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

/*
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
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *RuleController) CreateRule() {
	var rule models.Rule
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &rule)
	id, err := me.AddRule(&rule)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetRulesCnt
// @Description get Rules number
// @Param   query     query   string  false    "rule name"
// @Success 200  {total:int} rule number
// @Failure 403 string error
// @router /cnt [get]
func (c *RuleController) GetRulesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetRulesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetRules
// @Description get all Rules
// @Param   query     query   string  false    "rule name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} models.Rule
// @Failure 403 error string
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
// @Success 200 {object} models.Rule
// @Failure 403 error string
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
// @Success 200 {object} models.Rule
// @Failure 403 error string
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
*/
