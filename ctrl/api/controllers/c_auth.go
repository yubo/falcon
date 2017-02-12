/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Auth
type AuthController struct {
	BaseController
}

// @Title AuthLogin
// @Description auth login, such as ldap auth
// @Param	username	query	string	true	"username for login"
// @Param	password	query	string	true	"passworld for login"
// @Param	method		query	string	true	"login method"
// @Success 200 {object} models.User
// @Failure 406 error
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
			goto out
			/*
				err = fmt.Errorf("%s(%s/%s)",
					models.ErrLogged.Error(), user.Name, user.Uuid)
				goto out_err
			*/
		}
	}

	if method = c.GetString("method"); method == "" {
		err = models.ErrParam
		goto out_err
	}

	if auth, ok = models.AuthMap[method]; !ok {
		err = models.ErrNoExits
		goto out_err
	}

	if ok, uuid, err = auth.Verify(c); !ok {
		err = models.ErrAuthFailed
		goto out_err
	}

	user, _ = c.Access(uuid)
out:
	c.SendMsg(200, user)
	return

out_err:
	beego.Debug(err)
	c.SendMsg(405, err.Error())
}

// @Title Auth module callback handle
// @Description Auth module callback handle
// @Param	module	path	string	true	"the module you want to use"
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
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	user, err = me.GetUserByUuid(uuid)
	if err != nil {
		sys, _ := models.GetUser(1)
		user, err = sys.AddUser(&models.User{Uuid: uuid, Name: uuid})
		if err != nil {
			return
		}
	}
	c.SetSession("uid", user.Id)
	return
}
