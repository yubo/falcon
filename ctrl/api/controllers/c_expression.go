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
// @Param	body	body 	models.ExpressionAction	true	"body for expression content"
// @Success 200 {id:int} models.Expression.Id
// @Failure 403 string error
// @router / [post]
func (c *ExpressionController) CreateExpression() {
	var ea models.ExpressionAction
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)
	beego.Debug("pause", ea.Expression.Pause)

	id, err := me.AddAction(&ea.Action)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	ea.Expression.ActionId = id
	ea.Expression.CreateUserId = me.Id
	id, err = me.AddExpression(&ea.Expression)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetExpressionsCnt
// @Description get Expressions number
// @Param   query	query   string	false    "expression name"
// @Param   mine	query   bool	false    "only show mine expressions, default true"
// @Success 200 {int} expression number
// @Failure 403 string error
// @router /cnt [get]
func (c *ExpressionController) GetExpressionsCnt() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	if mine {
		user_id = me.Id
	}
	cnt, err := me.GetExpressionsCnt(query, user_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetExpressions
// @Description get all Expressions
// @Param   query	query   string  false    "expression name"
// @Param   mine	query   bool	false    "only show mine expressions, default true"
// @Param   per		query   int     false    "per page number"
// @Param   offset	query   int     false    "offset  number"
// @Success 200 [object] []models.ExpressionUi
// @Failure 403 error string
// @router /search [get]
func (c *ExpressionController) GetExpressions() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	if mine {
		user_id = me.Id
	}
	ret, err := me.GetExpressions(query, user_id, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Get
// @Description get expression by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object models.ExpressionAction
// @Failure 403 error string
// @router /:id [get]
func (c *ExpressionController) GetExpressionAction() {

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if ret, err := me.GetExpressionAction(id); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title UpdateExpressionAction
// @Description update the expression
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.Expression	true	"body for expression content"
// @Success 200 {object} models.Expression
// @Failure 403 error string
// @router /:id [put]
func (c *ExpressionController) UpdateExpressionAction() {
	var ea models.ExpressionAction

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateExpressionAction(id, &ea); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title UpdateExpression
// @Description update the expression
// @Param	id	query 	string	true	"The id you want to update"
// @Param	pause	query 	int	true	"1: pause, 0: not pause"
// @Success 200 null success
// @Failure 403 error string
// @router /pause [put]
func (c *ExpressionController) PauseExpression() {
	var pause int

	id, err := c.GetInt64("id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	pause, err = c.GetInt("pause")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if _, err := me.PauseExpression(id, pause); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, nil)
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
