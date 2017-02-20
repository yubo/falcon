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

// Operations about Roles
type RoleController struct {
	BaseController
}

// @Title CreateRole
// @Description create roles
// @Param	body	body 	models.Role	true	"body for role content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *RoleController) CreateRole() {
	var role models.Role
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &role)
	role.Id = 0

	if id, err := op.AddRole(&role); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetRolesCnt
// @Description get Roles number
// @Param   query     query   string  false    "role name"
// @Success 200 {total:int} role number
// @Failure 403 string error
// @router /cnt [get]
func (c *RoleController) GetRolesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetRolesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetRoles
// @Description get all Roles
// @Param   query     query   string  false    "role name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} models.Role
// @Failure 403 erorr string
// @router /search [get]
func (c *RoleController) GetRoles() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	roles, err := op.GetRoles(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, roles)
	}
}

// @Title Get
// @Description get role by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object} models.Role
// @Failure 403 error string
// @router /:id [get]
func (c *RoleController) GetRole() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		role, err := op.GetRole(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, role)
		}
	}
}

// @Title UpdateRole
// @Description update the role
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.Role	true	"body for role content"
// @Success 200 {object} models.Role
// @Failure 40x error string
// @router /:id [put]
func (c *RoleController) UpdateRole() {
	var role models.Role

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &role)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if o, err := op.UpdateRole(id, &role); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, o)
	}
}

// @Title DeleteRole
// @Description delete the role
// @Param	id	path 	string	true	"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 error string
// @router /:id [delete]
func (c *RoleController) DeleteRole() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteRole(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
