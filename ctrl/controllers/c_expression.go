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

// Operations about Expressions
type ExpressionController struct {
	BaseController
}

// @Title CreateExpression
// @Description create expressions
// @Param	body	body 	models.Expression	true	"body for expression content"
// @Success {code:200, data:int} models.Expression.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *ExpressionController) CreateExpression() {
	var expression models.Expression
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &expression)
	id, err := me.AddExpression(&expression)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, id)
	}
}

// @Title GetExpressionsCnt
// @Description get Expressions number
// @Param   query     query   string  false    "expression name"
// @Success {code:200, data:int} expression number
// @Failure {code:int, msg:string}
// @router /cnt [get]
func (c *ExpressionController) GetExpressionsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetExpressionsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title GetExpressions
// @Description get all Expressions
// @Param   query     query   string  false    "expression name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success {code:200, data:object} models.Expression
// @Failure {code:int, msg:string}
// @router /search [get]
func (c *ExpressionController) GetExpressions() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	expressions, err := me.GetExpressions(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, expressions)
	}
}

// @Title Get
// @Description get expression by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Expression
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *ExpressionController) GetExpression() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		expression, err := me.GetExpression(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, expression)
		}
	}
}

// @Title UpdateExpression
// @Description update the expression
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Expression	true		"body for expression content"
// @Success {code:200, data:object} models.Expression
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *ExpressionController) UpdateExpression() {
	var expression models.Expression

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &expression)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateExpression(id, &expression); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteExpression
// @Description delete the expression
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *ExpressionController) DeleteExpression() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteExpression(id)
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
func (c *MainController) GetExpression() {
	var expressions []*models.Expression

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryExpressions(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&expressions)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Expression")
	c.Data["Expressions"] = expressions
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "expression name"}

	c.TplName = "expression/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditExpression() {
	var expression *models.Expression
	var me *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	expression, err = me.GetExpression(id)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Expression")
	c.Data["Expression"] = expression
	c.Data["H1"] = "edit expression"
	c.Data["Method"] = "put"
	c.TplName = "expression/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddExpression() {

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Expression")
	c.Data["Method"] = "post"
	c.Data["H1"] = "add expression"
	c.TplName = "expression/edit.tpl"
}
