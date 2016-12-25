/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success {code:200, data:int} models.User.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *UserController) CreateUser() {
	var user models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if id, err := me.AddUser(&user); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetUsersCnt
// @Description get Users number
// @Success {code:200, data:int} user number
// @Failure {code:int, msg:string}
// @router /cnt/:query [get]
func (c *UserController) GetUsersCnt() {
	query := strings.TrimSpace(c.GetString(":query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	if cnt, err := me.GetUsersCnt(query); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetUsers
// @Description get all Users
// @Param   per       query   int  false       "per page number"
// @Param   offset    query   int  false       "offset  number"
// @Success {code:200, data:object} models.User
// @Failure {code:int, msg:string}
// @router /search/:query [get]
func (c *UserController) GetUsers() {
	query := strings.TrimSpace(c.GetString(":query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	if users, err := me.GetUsers(query, per, offset); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, users)
	}
}

// @Title Get
// @Description get user by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.User
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *UserController) GetUser() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		if user, err := me.GetUser(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, user)
		}
	}
}

// @Title Update
// @Description update the user
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success {code:200, data:object} models.User
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *UserController) UpdateUser() {
	var user models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &user)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateUser(id, &user); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title Delete
// @Description delete the user
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:string} delete success!
// @Failure {code:int, msg:string}
// @router /:id [delete]
func (c *UserController) DeleteUser() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if err = me.DeleteUser(id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	c.SendObj(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetUser() {
	var users []*models.User

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryUsers(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&users)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Users"] = users
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "/user"}

	c.TplName = "user/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditUser() {
	var me, user *models.User

	id, err := c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	if user, err = me.GetUser(id); err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["User"] = user
	c.Data["H1"] = "edit user"
	c.Data["Method"] = "put"
	c.TplName = "user/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddUser() {

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Method"] = "post"
	c.Data["H1"] = "add user"
	c.TplName = "user/edit.tpl"
}
