// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/yubo/falcon/ctrl/controllers"
	"github.com/yubo/falcon/ctrl/models"
)

const (
	ACL = false
)

func init() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"ApiToken"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

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
	beego.Router("/login", ac, "get:GetLogin;post:PostLogin")
	beego.Router("/logout", ac, "get:Logout")

	beego.Router("/settings/profile", mc, "get:GetProfile")
	beego.Router("/settings/aboutme", mc, "get:GetAboutMe")
	beego.Router("/settings/config/:module", mc, "get:GetConfig;post:PostConfig")
	if beego.BConfig.RunMode == "dev" {
		beego.Router("/settings/debug", mc, "get:GetDebug")
		beego.Router("/settings/debug/:action", mc, "get:GetDebugAction")
	}

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

	beego.Router("/token", mc, "get:GetToken")
	beego.Router("/token/edit/:id([0-9]+)", mc, "get:EditToken")
	beego.Router("/token/add", mc, "get:AddToken")

	beego.Router("/rel/tag/host", mc, "get:GetTagHost")
	beego.Router("/rel/tag/role/user", mc, "get:GetTagRoleUser")
	beego.Router("/rel/tag/role/token", mc, "get:GetTagRoleToken")
	beego.Router("/rel/tag/rule/trigger", mc, "get:GetTagRuleTrigger")

	beego.Router("/team", mc, "get:GetTeam")
	beego.Router("/teamusers/edit/:id([0-9]+)", mc, "get:EditTeamUsers")
	beego.Router("/teamusers/add", mc, "get:AddTeamUsers")

	beego.Router("/rule", mc, "get:GetRule")
	beego.Router("/rule/edit/:id([0-9]+)", mc, "get:EditRule")
	beego.Router("/rule/add", mc, "get:AddRule")

	beego.Router("/about", mc, "get:About")
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
	if id, ok := ctx.Input.Session("uid").(int64); ok {
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
