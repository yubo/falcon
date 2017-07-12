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

// Operations about Tokens
type TokenController struct {
	BaseController
}

// @Title CreateToken
// @Description create tokens
// @Param	body	body 	models.Token	true	"body for token content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *TokenController) CreateToken() {
	var token models.Token
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	token.Id = 0

	id, err := op.AddToken(&token)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTokensCnt
// @Description get Tokens number
// @Param   query     query   string  false       "token name"
// @Success 200 {object} models.Total token total number
// @Failure 400 string error
// @router /cnt [get]
func (c *TokenController) GetTokensCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTokensCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTokens
// @Description get all Tokens
// @Param   query     query   string  false       "token name"
// @Param   limit       query   int     false       "limit page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Token tokens info
// @Failure 400 string error
// @router /search [get]
func (c *TokenController) GetTokens() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	tokens, err := op.GetTokens(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, tokens)
	}
}

// @Title Get
// @Description get token by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Token token info
// @Failure 400 string error
// @router /:id [get]
func (c *TokenController) GetToken() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		token, err := op.GetToken(id)
		if err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, token)
		}
	}
}

// @Title UpdateToken
// @Description update the token
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Token	true		"body for token content"
// @Success 200 {object} models.Token token info
// @Failure 400 string error
// @router /:id [put]
func (c *TokenController) UpdateToken() {
	var token models.Token

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &token)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateToken(id, &token); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteToken
// @Description delete the token
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:400, msg:string}
// @router /:id [delete]
func (c *TokenController) DeleteToken() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteToken(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
