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
	"strconv"
	"strings"

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
		beego.InsertFilter("/v1.0/*", beego.BeforeRouter, accessFilter)
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
		beego.NSNamespace("/admin", beego.NSInclude(&controllers.AdminController{})),
	)
	beego.AddNamespace(ns)
}

func accessFilter(ctx *context.Context) {
	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/auth") {
		return
	}

	op, ok := ctx.Input.GetData("op").(*models.Operator)
	if !ok || op.User == nil {
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/admin") {
		if !op.IsAdmin() {
			http.Error(ctx.ResponseWriter, "permission denied", 403)
		}
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/user") {
		if id, err := strconv.ParseInt(
			ctx.Request.RequestURI[len("/v1.0/user/"):],
			10, 64); err != nil {

			if id == op.User.Id {
				// if target is self, is ok
				return
			}
		}
	}

	switch ctx.Request.Method {
	case "GET":
		if !op.IsReader() {
			http.Error(ctx.ResponseWriter, "permission denied", 403)
		}
	case "POST", "PUT", "DELETE":
		if !op.IsOperator() {
			http.Error(ctx.ResponseWriter, "permission denied", 403)
		}
	default:
		http.Error(ctx.ResponseWriter, "Method Not Allowed", 405)
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
