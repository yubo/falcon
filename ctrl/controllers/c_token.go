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
// @Success {code:200, data:int} models.Token.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *TokenController) CreateToken() {
	var token models.Token
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	id, err := me.AddToken(&token)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, id)
	}
}

// @Title GetTokensCnt
// @Description get Tokens number
// @Param   query     query   string  false       "token name"
// @Success {code:200, data:int} token number
// @Failure {code:int, msg:string}
// @router /cnt [get]
func (c *TokenController) GetTokensCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTokensCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title GetTokens
// @Description get all Tokens
// @Param   query     query   string  false       "token name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success {code:200, data:object} models.Token
// @Failure {code:int, msg:string}
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
// @Success {code:200, data:object} models.Token
// @Failure {code:int, msg:string}
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
// @Success {code:200, data:object} models.Token
// @Failure {code:int, msg:string}
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

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetToken() {
	var tokens []*models.Token

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryTokens(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&tokens)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Token")
	c.Data["Tokens"] = tokens
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "tag name"}

	c.TplName = "token/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditToken() {
	var token *models.Token
	var me *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	token, err = me.GetToken(id)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Token")
	c.Data["Token"] = token
	c.Data["H1"] = "edit token"
	c.Data["Method"] = "put"
	c.TplName = "token/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddToken() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Token")
	c.Data["H1"] = "add token"
	c.Data["Method"] = "post"
	c.TplName = "token/edit.tpl"
	return
}
