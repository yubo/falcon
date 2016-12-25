// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/ctrl/controllers"
	"github.com/yubo/falcon/ctrl/models"
)

const (
	ACL = false
)

func init() {

	beego.InsertFilter("/*", beego.BeforeRouter, profileFilter)

	if ACL {
		// web
		beego.InsertFilter("/settings/*", beego.BeforeRouter, authFilter)
		// api
		beego.InsertFilter("/v1.0/((user)|(ugroup))/*", beego.BeforeRouter, adminFiler)
	}

	ac := &controllers.AuthController{}
	mc := &controllers.MainController{}

	beego.Router("/", mc, "get:Get")
	beego.Router("/login", ac, "get:LoginGet;post:LoginPost")
	beego.Router("/logout", ac, "get:Logout")

	beego.Router("/settings/profile", mc, "get:GetProfile")
	beego.Router("/settings/about", mc, "get:GetAbout")

	beego.Router("/user", mc, "get:GetUser")
	beego.Router("/user/edit/:id([0-9]+)", mc, "get:EditUser")
	beego.Router("/user/add", mc, "get:AddUser")

	beego.Router("/host", mc, "get:GetHost")
	beego.Router("/host/edit/:id([0-9]+)", mc, "get:EditHost")
	beego.Router("/host/add", mc, "get:AddHost")

	beego.Router("/tag", mc, "get:GetTag")
	beego.Router("/tag/edit/:id([0-9]+)", mc, "get:EditTag")
	beego.Router("/tag/add", mc, "get:AddTag")

	beego.Router("/role", mc, "get:GetRole")
	beego.Router("/role/edit/:id([0-9]+)", mc, "get:EditRole")
	beego.Router("/role/add", mc, "get:AddRole")

	beego.Router("/system", mc, "get:GetSystem")
	beego.Router("/system/edit/:id([0-9]+)", mc, "get:EditSystem")
	beego.Router("/system/add", mc, "get:AddSystem")

	beego.Router("/scope/:sysid([0-9]+)", mc, "get:GetScope")
	beego.Router("/scope/edit/:id([0-9]+)", mc, "get:EditScope")
	beego.Router("/scope/add/:sysid([0-9]+)", mc, "get:AddScope")
}

/*
 * filter
 */
func authFilter(ctx *context.Context) {
	if ctx.Input.GetData("me") == nil {
		beego.Debug("not login")
		ctx.Redirect(302, "/login")
	}
}

func adminFiler(ctx *context.Context) {
	beego.Debug("checkAdmin")
	if !IsAdmin(ctx) {
		ctx.Redirect(302, "/")
	}
}

func profileFilter(ctx *context.Context) {
	if id, ok := ctx.Input.Session("uid").(int); ok {
		me, err := models.GetUser(id)
		if err != nil {
			return
		}
		ctx.Input.SetData("me", me)
		/*
			if me.Name == "" &&
				!strings.HasPrefix(ctx.Request.URL.String(),
					"/settings") {
				beego.Debug("Redirect /settings/profile")
				ctx.Redirect(302, "/settings/profile")
				return
			}
		*/
	}
}

func IsAdmin(ctx *context.Context) bool {
	if me, ok := ctx.Input.GetData("me").(*models.User); ok && me.Id == 1 {
		return true
	}
	return false
}
