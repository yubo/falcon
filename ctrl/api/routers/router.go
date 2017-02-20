// @APIVersion 1.0.0
// @Title falcon ctrl API
// @Description Open-Falcon 是小米运维部开源的一款互联网企业级监控系统解决方案.
// @Contact yubo@xiaomi.com
// @TermsOfServiceUrl http://open-falcon.org/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/ctrl/api/controllers"
	"github.com/yubo/falcon/ctrl/api/models"
)

const (
	ACL = true
)

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, profileFilter)

	if ACL {
		beego.InsertFilter("/v1.0/host/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/role/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/user/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/token/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/team/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/template/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/strategy/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/expression/*", beego.BeforeRouter, authFilter)
		beego.InsertFilter("/v1.0/rel/*", beego.BeforeRouter, authFilter)
	}
	ns := beego.NewNamespace("/v1.0",
		beego.NSNamespace("/auth", beego.NSInclude(&controllers.AuthController{})),
		beego.NSNamespace("/host", beego.NSInclude(&controllers.HostController{})),
		beego.NSNamespace("/role", beego.NSInclude(&controllers.RoleController{})),
		beego.NSNamespace("/tag", beego.NSInclude(&controllers.TagController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/token", beego.NSInclude(&controllers.TokenController{})),
		beego.NSNamespace("/rel", beego.NSInclude(&controllers.RelController{})),
		beego.NSNamespace("/team", beego.NSInclude(&controllers.TeamController{})),
		beego.NSNamespace("/template", beego.NSInclude(&controllers.TemplateController{})),
		beego.NSNamespace("/expression", beego.NSInclude(&controllers.ExpressionController{})),
		beego.NSNamespace("/strategy", beego.NSInclude(&controllers.StrategyController{})),
		beego.NSNamespace("/settings", beego.NSInclude(&controllers.SetController{})),
		beego.NSNamespace("/metric", beego.NSInclude(&controllers.MetricController{})),
	)
	beego.AddNamespace(ns)
}

func authFilter(ctx *context.Context) {
	op, ok := ctx.Input.GetData("op").(*models.Operator)
	if !ok || op.User == nil {
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
	}
}

func adminFiler(ctx *context.Context) {
	if !IsAdmin(ctx) {
		http.Error(ctx.ResponseWriter, "permission denied", 403)
	}
}

func profileFilter(ctx *context.Context) {
	op := &models.Operator{O: orm.NewOrm()}
	ctx.Input.SetData("op", op)
	if id, ok := ctx.Input.Session("uid").(int64); ok {
		u, err := models.GetUser(id, op.O)
		if err != nil {
			beego.Debug("login, but can not found user")
			return
		}

		op.User = u
		op.Token, _ = ctx.Input.Session("token").(int)
	} else {
		beego.Debug("not login 2")
	}
}

func IsAdmin(ctx *context.Context) bool {
	if me, ok := ctx.Input.GetData("me").(*models.User); ok && me.Id == 1 {
		return true
	}
	return false
}
