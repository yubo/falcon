/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Scopes
type ScopeController struct {
	BaseController
}

// @Title CreateScope
// @Description create scopes
// @Param	body	body 	models.Scope	true	"body for scope content"
// @Success {code:200, data:int} models.Scope.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *ScopeController) CreateScope() {
	var scope models.Scope
	json.Unmarshal(c.Ctx.Input.RequestBody, &scope)
	beego.Debug(string(c.Ctx.Input.RequestBody))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	id, err := me.AddScope(&scope)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetScopesCnt
// @Description get Scopes number
// @Param   system_id query   int  true        "system id"
// @Success {code:200, data:int} scope number
// @Failure {code:int, msg:string}
// @router /cnt/:query [get]
func (c *ScopeController) GetScopesCnt() {
	query := strings.TrimSpace(c.GetString(":query"))
	sysid, err := c.GetInt64("system_id", 0)

	if err != nil || sysid == 0 {
		c.SendMsg(403, models.ErrParam.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	cnt, err := me.GetScopesCnt(sysid, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetScopes
// @Description get all Scopes
// @Param   per       query   int  false       "per page number"
// @Param   offset    query   int  false       "offset  number"
// @Param   system_id query   int  true        "system id"
// @Success {code:200, data:object} models.Scope
// @Failure {code:int, msg:string}
// @router /search/:query [get]
func (c *ScopeController) GetScopes() {
	query := strings.TrimSpace(c.GetString(":query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)

	sysid, err := c.GetInt64("system_id", 0)
	if err != nil || sysid == 0 {
		c.SendMsg(403, models.ErrParam.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	scopes, err := me.GetScopes(sysid, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, scopes)
	}
}

// @Title Get
// @Description get scope by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Scope
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *ScopeController) GetScope() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		scope, err := me.GetScope(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, scope)
		}
	}
}

// @Title UpdateScope
// @Description update the scope
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Scope	true		"body for scope content"
// @Success {code:200, data:object} models.Scope
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *ScopeController) UpdateScope() {
	var scope models.Scope

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &scope)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateScope(id, &scope); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title DeleteScope
// @Description delete the scope
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *ScopeController) DeleteScope() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteScope(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendObj(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetScope() {
	var scopes []*models.Scope

	sysid, err := c.GetInt64(":sysid")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryScopes(sysid, query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&scopes)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Scopes"] = scopes
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", fmt.Sprintf("/scope/%d", sysid)}

	c.TplName = "scope/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditScope() {
	var scope *models.Scope
	var sys *models.System
	var me *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	scope, err = me.GetScope(id)
	if err != nil {
		goto out
	}
	sys, err = me.GetSystem(scope.System_id)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Scope"] = scope
	c.Data["System"] = sys
	c.Data["H1"] = "edit scope at "
	c.Data["Method"] = "put"
	c.TplName = "scope/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddScope() {
	var sys *models.System
	var me *models.User

	sysid, err := c.GetInt64(":sysid")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	sys, err = me.GetSystem(sysid)
	if err != nil {
		goto out
	}
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Scope"] = &models.Scope{System_id: sysid}
	c.Data["System"] = sys
	c.Data["H1"] = "add scope at "
	c.Data["Method"] = "post"
	c.TplName = "scope/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}
