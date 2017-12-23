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
// @Param	body	body	models.UserApiAdd	true	"body for user content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *UserController) CreateUser() {
	var user models.UserApiAdd
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
// @Param	body	body	models.UserApiUpdate	true	"body for user content"
// @Success 200 {object} models.User user info
// @Failure 400 string error
// @router / [put]
func (c *UserController) UpdateUser() {
	input := models.UserApiUpdate{}
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

// @Title Delete
// @Description delete the user
// @Param body	body	[]int	true	"The []id you want to delete"
// @Success 200 {object} models.Stats api call result
// @Failure 400 string error
// @router / [delete]
func (c *UserController) DeleteUsers() {
	var (
		ids              []int64
		errs             []string
		success, failure int64
	)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	for _, id := range ids {
		if err := op.DeleteUser(id); err != nil {
			errs = append(errs, err.Error())
			failure++
		} else {
			success++
		}
	}
	c.SendMsg(200, statsObj(success, failure, errs))
}

/*******************************************************************************
 ************************ tag role user ****************************************
 ******************************************************************************/

// @Title GetTagRoleUserCnt
// @Description get tag role user number
// @Param	query	query   string  false	"user name"
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total user total number
// @Failure 400 string error
// @router /tag/role/cnt [get]
func (c *UserController) GetTagRoleUserCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagRoleUserCnt(tagId, query, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUser
// @Description get tag role user
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"user name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagRoleUserApiGet tag role user info
// @Failure 400 string error
// @router /tag/role/search [get]
func (c *UserController) GetTagRoleUser() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagRoleUser(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role users relation
// @Description create tag/roles/users relation
// @Param	body	body 	models.TagRolesUsersApiAdd	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/role [post]
func (c *UserController) CreateTagRolesUsers() {
	var rel models.TagRolesUsersApiAdd
	var cnt int64

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rel)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	for _, role := range rel.RoleIds {
		for _, user := range rel.UserIds {
			if _, err := op.CreateTagRoleUser(
				&models.TagRoleUserApi{
					TagId:  rel.TagId,
					RoleId: role,
					UserId: user,
				}); err == nil {

				cnt++
			}
		}
	}

	c.SendMsg(200, totalObj(cnt))
}

// @Title delete tag role user relation
// @Description delete tag/roles/users relation
// @Param	body		body 	models.TagRolesUsersApiDel	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/role [delete]
func (c *UserController) DelTagRolesUsers() {
	var rel models.TagRolesUsersApiDel
	var cnt int64

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rel)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	for _, ru := range rel.RoleUser {
		if _, err := op.DeleteTagRoleUser(
			&models.TagRoleUserApi{
				TagId:  rel.TagId,
				RoleId: ru.RoleId,
				UserId: ru.UserId,
			}); err == nil {

			cnt++
		}
	}
	c.SendMsg(200, totalObj(cnt))
}
