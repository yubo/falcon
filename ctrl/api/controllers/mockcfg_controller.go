/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Mockcfgs
type MockcfgController struct {
	BaseController
}

// @Title CreateMockcfg
// @Description create mockcfgs
// @Param	body	body 	models.NoDataApiPut	true	"body for mockcfg content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *MockcfgController) CreateMockcfg() {
	var inputs models.NoDataApiPut
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if id, err := op.AddMockcfg(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetMockcfgsCnt
// @Description get Mockcfgs number
// @Param   query	query   string  false	"mockcfg name"
// @Param   mine	query   bool	false	"only show mine expressions, default true"
// @Success 200 {object} models.Total mockcfg total number
// @Failure 400 string error
// @router /cnt [get]
func (c *MockcfgController) GetMockcfgsCnt() {
	var user_name string
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_name = op.User.Name
	}
	cnt, err := op.GetMockcfgsCnt(query, user_name)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetMockcfgs
// @Description get all Mockcfgs
// @Param   query     query   string  false       "mockcfg name"
// @Param   mine	query   bool	false	"only show mine expressions, default true"
// @Param   limit       query   int     false       "limit page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Mockcfg  mockcfgs info
// @Failure 400 string error
// @router /search [get]
func (c *MockcfgController) GetMockcfgs() {
	var user_name string
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_name = op.User.Name
	}
	mockcfgs, err := op.GetMockcfgs(query, user_name, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, mockcfgs)
	}
}

// @Title Get
// @Description get mockcfg by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Mockcfg mockcfg info
// @Failure 400 string error
// @router /:id [get]
func (c *MockcfgController) GetMockcfg() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if mockcfg, err := op.GetMockcfg(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, mockcfg)
		}
	}
}

// @Title UpdateMockcfg
// @Description update the mockcfg
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.NoDataApiPut	true	"body for mockcfg content"
// @Success 200 {object} models.NoDataApiPut mockcfg info
// @Failure 400 string error
// @router /:id [put]
func (c *MockcfgController) UpdateMockcfg() {
	var mockcfg models.NoDataApiPut
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &mockcfg)

	if u, err := op.UpdateMockcfg(id, &mockcfg); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteMockcfg
// @Description delete the mockcfg
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /:id [delete]
func (c *MockcfgController) DeleteMockcfg() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteMockcfg(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}
