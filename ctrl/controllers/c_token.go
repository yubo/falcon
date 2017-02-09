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

// Operations about Tokens
type TokenController struct {
	BaseController
}

// @Title CreateToken
// @Description create tokens
// @Param	body	body 	models.Token	true	"body for token content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *TokenController) CreateToken() {
	var token models.Token
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	id, err := me.AddToken(&token)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetTokensCnt
// @Description get Tokens number
// @Param   query     query   string  false       "token name"
// @Success 200  {total:int} token total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TokenController) GetTokensCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTokensCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTokens
// @Description get all Tokens
// @Param   query     query   string  false       "token name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} models.Token
// @Failure 403 error string
// @router /search [get]
func (c *TokenController) GetTokens() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	tokens, err := me.GetTokens(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, tokens)
	}
}

// @Title Get
// @Description get token by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Token
// @Failure 403 error string
// @router /:id [get]
func (c *TokenController) GetToken() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		token, err := me.GetToken(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, token)
		}
	}
}

// @Title UpdateToken
// @Description update the token
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Token	true		"body for token content"
// @Success 200 {object} models.Token
// @Failure 403 error string
// @router /:id [put]
func (c *TokenController) UpdateToken() {
	var token models.Token

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &token)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateToken(id, &token); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteToken
// @Description delete the token
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TokenController) DeleteToken() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteToken(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
