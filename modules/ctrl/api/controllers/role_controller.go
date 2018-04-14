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
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/ctrl/api/models"
)

// Operations about Roles
type RoleController struct {
	BaseController
}

// @Title CreateRole
// @Description create roles
// @Param	body	body 	models.RoleCreate	true	"body for role content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *RoleController) CreateRole() {
	var role models.RoleApiAdd
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &role)

	if id, err := op.CreateRole(&role); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetRolesCnt
// @Description get Roles number
// @Param   query     query   string  false    "role name"
// @Success 200 {object} models.Total role number
// @Failure 400 string error
// @router /cnt [get]
func (c *RoleController) GetRolesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetRolesCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetRoles
// @Description get all Roles
// @Param   query     query   string  false    "role name"
// @Param   limit       query   int     false    "limit page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Role roles info
// @Failure 400 string error
// @router /search [get]
func (c *RoleController) GetRoles() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	roles, err := op.GetRoles(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, roles)
	}
}

// @Title Get
// @Description get role by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object} models.Role role info
// @Failure 400 string error
// @router /:id [get]
func (c *RoleController) GetRole() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		role, err := op.GetRole(id)
		if err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, role)
		}
	}
}

// @Title UpdateRole
// @Description update the role
// @Param	body	body 	models.RoleUpdate	true	"body for role content"
// @Success 200 {object} models.Role role info
// @Failure 400 string error
// @router / [put]
func (c *RoleController) UpdateRole() {
	input := models.RoleApiUpdate{}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	p, err := op.GetRole(input.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	role := *p
	core.Override(&role, &input)

	if ret, err := op.UpdateRole(&role); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title DeleteRole
// @Description delete the role
// @Param	id	path 	string	true	"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 string error
// @router /:id [delete]
func (c *RoleController) DeleteRole() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteRole(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
