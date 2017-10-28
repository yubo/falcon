/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.UserCreate	true		"body for user content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *UserController) CreateUser() {
	var user models.UserCreate
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)

	if id, err := op.CreateUser(&user); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetUsersCnt
// @Description get Users number
// @Param   query     query   string  false       "user name/email"
// @Success 200 {object} models.Total user total number
// @Failure 400 string error
// @router /cnt [get]
func (c *UserController) GetUsersCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if cnt, err := op.GetUsersCnt(query); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetUsers
// @Description get all Users
// @Param   query	query   string  false       "user name/email"
// @Param   limit	query   int     false       "limit page number"
// @Param   offset	query   int     false       "offset  number"
// @Success 200 {object} []models.User users info
// @Failure 400 string error
// @router /search [get]
func (c *UserController) GetUsers() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if users, err := op.GetUsers(query, limit, offset); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, users)
	}
}

// @Title unbind User
// @Description unbind user
// @Param id	path	int	true	"user id"
// @Success 200 string success
// @Failure 400 string error
// @router /unbind/:id [get]
func (c *UserController) UnBindUser() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if err := op.UnBindUser(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, "success")
		}
	}
}

// @Title GetBindedUsers
// @Description get all Users
// @Param   id     path   int  true       "user id"
// @Success 200 {object} []models.User users info
// @Failure 400 string error
// @router /binded/:id [get]
func (c *UserController) GetBindedUsers() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if users, err := op.GetBindedUsers(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, users)
		}
	}
}

// @Title Get
// @Description get user by id
// @Param	id	path	int	true	"user id"
// @Success 200 {object} models.User user info
// @Failure 400 string error
// @router /:id [get]
func (c *UserController) GetUser() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if user, err := op.GetUser(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, user)
		}
	}
}

// @Title Update
// @Description update user information
// @Param	body	body	models.UserUpdate	true	"body for user content"
// @Success 200 {object} models.User user info
// @Failure 400 string error
// @router / [put]
func (c *UserController) UpdateUser() {
	input := models.UserUpdate{}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	p, err := op.GetUser(input.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	user := *p
	falcon.Override(&user, &input)

	if ret, err := op.UpdateUser(&user); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Delete
// @Description delete the user
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:string} delete success!
// @Failure 400 string error
// @router /:id [delete]
func (c *UserController) DeleteUser() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if err = op.DeleteUser(id); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	c.SendMsg(200, "delete success!")
}
