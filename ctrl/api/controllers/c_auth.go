/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Auth
type AuthController struct {
	BaseController
}

// @Title get support auth modules
// @Description get support auth modules
// @Success 200 [string] []strins{module, ...}
// @Failure 405 error
// @router /modules [get]
func (c *AuthController) Modules() {
	m := []string{}
	for k, _ := range models.Auths {
		m = append(m, k)
	}

	c.SendMsg(200, m)
}

// @Title OAuth Login
// @Description auth login
// @Param	module	path	string	true	"the module you want to use(github/google)"
// @Success 302 redirect
// @Failure 405 error
// @router /login/:module [get]
func (c *AuthController) Authorize() {
	module := c.GetString(":module")

	auth, ok := models.Auths[module]
	if !ok {
		c.SendMsg(405, models.ErrNoModule.Error())
		return
	}

	URL := auth.AuthorizeUrl(c.Ctx)
	if URL == "" {
		c.SendMsg(405, nil)
		return
	}

	c.Ctx.Redirect(302, URL)
}

// @Title OAuth module callback handle
// @Description Auth module callback handle
// @Param	module	path	string	true	"the module you want to use"
// @Success 302 redirect to RedirectUrl(default "/")
// @Failure 406 not acceptable
// @router /callback/:module [get]
func (c *AuthController) Callback() {
	auth, ok := models.Auths[c.GetString(":module")]
	if !ok {
		c.SendMsg(406, models.ErrNoModule.Error())
		return
	}
	cb := c.GetString("cb")

	uuid, err := auth.CallBack(c.Ctx)
	if err != nil {
		c.SendMsg(406, err.Error())
		return
	}

	if _, err = c.Access(uuid); err != nil {
		c.SendMsg(406, err.Error())
		return
	}

	c.Ctx.Redirect(302, "/#"+cb)
}

// @Title AuthLogin
// @Description auth login, such as ldap auth
// @Param	username	query	string	false	"username for login"
// @Param	password	query	string	false	"passworld for login"
// @Param	method		query	string	false	"login method"
// @Success 200 {object} models.User
// @Failure 406 not acceptable
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

		if user, err = models.GetUser(id, orm.NewOrm()); err == nil {
			goto out
		}
	}

	if method = c.GetString("method"); method == "" {
		err = models.ErrParam
		goto out_err
	}

	if auth, ok = models.Auths[method]; !ok {
		err = models.ErrNoExits
		goto out_err
	}

	if ok, uuid, err = auth.Verify(c); !ok {
		err = models.ErrLogin
		goto out_err
	}

	user, _ = c.Access(uuid)
out:
	c.SendMsg(200, user)
	return

out_err:
	c.SendMsg(406, err.Error())
}

// @Title Auth Logout
// @Description user logout, reset cookie
// @Success 200 {string} logout success!
// @Failure 405 {string} Method Not Allowed
// @router /logout [get]
func (c *AuthController) Logout() {
	if uid := c.GetSession("uid"); uid != nil {
		c.DelSession("uid")
		c.SendMsg(200, "logout success!")
	} else {
		c.SendMsg(405, models.ErrNoLogged.Error())
	}
}

func (c *AuthController) Access(uuid string) (user *models.User, err error) {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	user, err = op.GetUserByUuid(uuid)
	if err != nil {
		beego.Debug("can't get user by uuid ", uuid)
		sys, _ := models.GetUser(1, op.O)
		sysOp := &models.Operator{
			User:  sys,
			O:     op.O,
			Token: models.SYS_F_A_TOKEN,
		}
		user, err = sysOp.AddUser(&models.User{Uuid: uuid, Name: uuid})
		if err != nil {
			beego.Debug("add user failed ", err.Error())
			return
		}
	}
	beego.Debug("get login user ", user)
	c.SetSession("uid", user.Id)
	c.SetSession("token", models.Tokens(user.Id, op.O))
	return
}
