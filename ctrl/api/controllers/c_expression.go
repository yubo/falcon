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

// Operations about Expressions
type ExpressionController struct {
	BaseController
}

// @Title CreateExpression
// @Description create expressions
// @Param	body	body 	models.Expression	true	"body for expression content"
// @Success 200 {int} models.Expression.Id
// @Failure 403 string error
// @router / [post]
func (c *ExpressionController) CreateExpression() {
	var expression models.Expression
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &expression)
	expression.Id = 0
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
// @Success 200 {int} expression number
// @Failure 403 string error
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
// @Success 200 {object} models.Expression
// @Failure 403 error string
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
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object models.Expression
// @Failure 403 error string
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
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.Expression	true	"body for expression content"
// @Success 200 {object} models.Expression
// @Failure 403 error string
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
// @Success 200  {string} delete success!
// @Failure 403  error string
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
