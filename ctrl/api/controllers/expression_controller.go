/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
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
// @Param	body	body 	models.ExpressionActionApiPut	true	"body for expression content"
// @Success 200 {object} models.Id models.Expression.Id
// @Failure 400 string error
// @router / [post]
func (c *ExpressionController) CreateExpression() {
	var ea models.ExpressionActionApiPut
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)
	beego.Debug("pause", ea.Expression.Pause)

	id, err := op.AddAction(&ea.Action)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	ea.Expression.ActionId = id
	ea.Expression.CreateUserId = op.User.Id
	id, err = op.AddExpression(&ea.Expression)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetExpressionsCnt
// @Description get Expressions number
// @Param   query	query   string	false    "expression name"
// @Param   mine	query   bool	false    "only show mine expressions, default true"
// @Success 200 {object} models.Total expression number
// @Failure 400 string error
// @router /cnt [get]
func (c *ExpressionController) GetExpressionsCnt() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	cnt, err := op.GetExpressionsCnt(query, user_id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetExpressions
// @Description get all Expressions
// @Param	query	query   string  false    "expression name"
// @Param	mine	query   bool	false    "only show mine expressions, default true"
// @Param	limit	query   int     false    "limit page number"
// @Param	offset	query   int     false    "offset  number"
// @Success 200 {object} []models.ExpressionApiGet expressionuis
// @Failure 400 string error
// @router /search [get]
func (c *ExpressionController) GetExpressions() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	ret, err := op.GetExpressions(query, user_id, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Get
// @Description get expression by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object} models.ExpressionAction expression and action info
// @Failure 400 string error
// @router /:id [get]
func (c *ExpressionController) GetExpressionAction() {

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if ret, err := op.GetExpressionAction(id); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title UpdateExpressionAction
// @Description update the expression
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.ExpressionActionApiPut	true	"body for expression content"
// @Success 200 {object} models.Expression expression
// @Failure 400 string error
// @router /:id [put]
func (c *ExpressionController) UpdateExpressionAction() {
	var ea models.ExpressionActionApiPut

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateExpressionAction(id, &ea); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title pause Expression
// @Description pause the expression
// @Param	id	query 	string	true	"The id you want to update"
// @Param	pause	query 	int	true	"1: pause, 0: not pause"
// @Success 200 null success
// @Failure 400 string error
// @router /pause [put]
func (c *ExpressionController) PauseExpression() {
	var pause int

	id, err := c.GetInt64("id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	pause, err = c.GetInt("pause")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if _, err := op.PauseExpression(id, pause); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, nil)
	}
}

// @Title DeleteExpression by expression
// @Description delete the expression by expression
// @Param	expr	query 	string	true	"The expression you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 string error
// @router /0 [delete]
func (c *ExpressionController) DeleteExpression0() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	expr := c.GetString("expr")
	if expr == "" {
		c.SendMsg(400, "expr is empty")
		return
	}

	err := op.DeleteExpression0(expr)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}

// @Title DeleteExpression
// @Description delete the expression
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 string error
// @router /:id [delete]
func (c *ExpressionController) DeleteExpression() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteExpression(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
