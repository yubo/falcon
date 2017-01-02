/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"fmt"

	"github.com/yubo/falcon/ctrl/models"
)

type AuthController struct {
	BaseController
}

// @Title Auth Login page
// @Description auth login page render, for test
// @router /login [get]
func (c *AuthController) GetLogin() {
	modules := make([]*models.AuthModule, 0)

	for _, auth := range models.Auths {
		modules = append(modules,
			&models.AuthModule{Name: auth.GetName(),
				Html: auth.LoginHtml(c),
			})
	}

	c.PrepareEnv(nil, "")
	c.Data["Modules"] = modules
	c.TplName = "login.tpl"
}

// @Title AuthLogin
// @Description auth login, such as ldap auth
// @Param	username	query	string	true	"username for login"
// @Param	password	query	string	true	"passworld for login"
// @Param	method		query	string	true	"login method"
// @Success 200 {uid:string, uuid:string} models.User.Id, models.User.Uuid
// @Failure 406 {string} error
// @router /login [post]
func (c *AuthController) PostLogin() {
	var (
		user         *models.User
		err          error
		ok           bool
		uuid, method string
		auth         models.AuthInterface
	)
	if id, ok := c.GetSession("uid").(int64); ok {
		if user, err = models.GetUser(id); err == nil {
			c.Data["json"] = fmt.Errorf("%s(%s/%s)",
				models.ErrLogged.Error(), user.Name, user.Uuid)
			goto out_err
		}
	}

	if method = c.GetString("method"); method == "" {
		c.Data["json"] = models.ErrNoExits.Error()
		goto out_err
	}

	if auth, ok = models.AuthMap[method]; !ok {
		c.Data["json"] = models.ErrNoExits.Error()
		goto out_err
	}

	if ok, uuid, err = auth.Verify(c); !ok {
		c.Data["json"] = models.ErrAuthFailed.Error()
		goto out_err
	}

	c.Access(uuid)
	c.Ctx.Redirect(302, "/")
	return

out_err:
	c.Ctx.ResponseWriter.WriteHeader(406)
	c.ServeJSON()
}

// @Title Auth module callback handle
// @Description Auth module callback handle
// @router /callback/:module [get]
func (c *AuthController) Callback() {
	module := c.Ctx.Input.Param(":module")

	if auth, ok := models.AuthMap[module]; !ok {
		c.Ctx.Redirect(302, "/")
	} else {
		auth.CallBack(c)
	}
}

// @Title Auth Logout
// @Description user logout, reset cookie
// @Success {code:200, data:string} logout success!
// @Failure {code:405, msg:string} Method Not Allowed
// @router /logout [get]
func (c *AuthController) Logout() {
	if uid := c.GetSession("uid"); uid != nil {
		c.DelSession("uid")
		c.SendObj(200, "logout success!")
	} else {
		c.SendMsg(405, models.ErrNoLogged.Error())
	}
}

func (c *AuthController) Access(uuid string) (id int64) {
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	user, err := me.GetUserByUuid(uuid)
	if err != nil {
		sys, _ := models.GetUser(1)
		id, err = sys.AddUser(&models.User{Uuid: uuid})
		if err != nil {
			return -1
		}
		user, err = me.GetUser(id)
	}
	c.SetSession("uid", user.Id)
	return user.Id
}
