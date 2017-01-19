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

// Operations about Roles
type RoleController struct {
	BaseController
}

// @Title CreateRole
// @Description create roles
// @Param	body	body 	models.Role	true	"body for role content"
// @Success {code:200, data:int} models.Role.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *RoleController) CreateRole() {
	var role models.Role
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &role)
	id, err := me.AddRole(&role)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetRolesCnt
// @Description get Roles number
// @Param   query     query   string  false    "role name"
// @Success {code:200, data:int} role number
// @Failure {code:int, msg:string}
// @router /cnt [get]
func (c *RoleController) GetRolesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetRolesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetRoles
// @Description get all Roles
// @Param   query     query   string  false    "role name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success {code:200, data:object} models.Role
// @Failure {code:int, msg:string}
// @router /search [get]
func (c *RoleController) GetRoles() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	roles, err := me.GetRoles(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, roles)
	}
}

// @Title Get
// @Description get role by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Role
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *RoleController) GetRole() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		role, err := me.GetRole(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, role)
		}
	}
}

// @Title UpdateRole
// @Description update the role
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Role	true		"body for role content"
// @Success {code:200, data:object} models.Role
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *RoleController) UpdateRole() {
	var role models.Role

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &role)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateRole(id, &role); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title DeleteRole
// @Description delete the role
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *RoleController) DeleteRole() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteRole(id)
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
func (c *MainController) GetRole() {
	var roles []*models.Role

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryRoles(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&roles)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Role")
	c.Data["Roles"] = roles
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "role name"}

	c.TplName = "role/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditRole() {
	var role *models.Role
	var me *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	role, err = me.GetRole(id)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Role")
	c.Data["Role"] = role
	c.Data["H1"] = "edit role"
	c.Data["Method"] = "put"
	c.TplName = "role/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddRole() {

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Role")
	c.Data["Method"] = "post"
	c.Data["H1"] = "add role"
	c.TplName = "role/edit.tpl"
}
