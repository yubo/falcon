// @APIVersion 1.0.0
// @Title falcon ctrl API
// @Description Open-Falcon 是小米运维部开源的一款互联网企业级监控系统解决方案.
// @Contact yubo@xiaomi.com
// @TermsOfServiceUrl http://open-falcon.org/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/api/controllers"
	"github.com/yubo/falcon/ctrl/api/models"
)

var (
	miMode           bool
	RateLimitsAccess func(*context.Context) bool
)

func Init(miMode_ bool) {

	miMode = miMode_

	beego.InsertFilter("/v1.0/*", beego.BeforeRouter, profileFilter)
	beego.InsertFilter("/v1.0/*", beego.BeforeRouter, accessFilter)

	ns := beego.NewNamespace("/v1.0",
		beego.NSNamespace("/auth", beego.NSInclude(&controllers.AuthController{})),
		beego.NSNamespace("/host", beego.NSInclude(&controllers.HostController{})),
		beego.NSNamespace("/role", beego.NSInclude(&controllers.RoleController{})),
		beego.NSNamespace("/tag", beego.NSInclude(&controllers.TagController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/token", beego.NSInclude(&controllers.TokenController{})),
		beego.NSNamespace("/settings", beego.NSInclude(&controllers.SetController{})),
		beego.NSNamespace("/metric", beego.NSInclude(&controllers.MetricController{})),
		beego.NSNamespace("/admin", beego.NSInclude(&controllers.AdminController{})),
		beego.NSNamespace("/dashboard", beego.NSInclude(&controllers.DashboardController{})),
		beego.NSNamespace("/graph", beego.NSInclude(&controllers.GraphController{})),
		beego.NSNamespace("/pub", beego.NSInclude(&controllers.PubController{})),
		beego.NSNamespace("/weapp", beego.NSInclude(&controllers.WeappController{})),
		beego.NSNamespace("/plugin", beego.NSInclude(&controllers.PluginController{})),
		beego.NSNamespace("/event", beego.NSInclude(&controllers.EventController{})),
		beego.NSNamespace("/action", beego.NSInclude(&controllers.ActionController{})),
	)
	beego.AddNamespace(ns)

}

func accessFilter(ctx *context.Context) {

	//time.Sleep(time.Millisecond * 100)

	if RateLimitsAccess != nil {
		if !RateLimitsAccess(ctx) {
			return
		}
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/auth") ||
		strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/pub") ||
		strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/weapp") {
		return
	}

	op, ok := ctx.Input.GetData("op").(*models.Operator)
	if !ok || op.User == nil {
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/settings") {
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/admin") {
		if !op.IsAdmin() {
			http.Error(ctx.ResponseWriter, "permission denied, admin only", 403)
		}
		return
	}

	switch ctx.Request.Method {
	case "GET":
		if !op.IsReader() {
			http.Error(ctx.ResponseWriter, "permission denied, read", 403)
		}
	case "POST", "PUT", "DELETE":
		if !op.IsOperator() {
			beego.Debug("TOKEN ", op.Token)
			http.Error(ctx.ResponseWriter, "permission denied, operate", 403)
		}
		if miMode {
			if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/host") &&
				ctx.Request.Method == "PUT" {
				return
			}
			if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/host") ||
				strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/role") ||
				strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/tag") ||
				strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/token") ||
				//strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/token/tag") ||
				strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/host/tag") ||
				strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/user/tag") {
				http.Error(ctx.ResponseWriter, "permission denied, xiaomi runmode", 403)
			}
		}
	default:
		http.Error(ctx.ResponseWriter, "Method Not Allowed", 405)
		return
	}
}

func profileFilter(ctx *context.Context) {
	//需要 切换数据库 和 事务处理 的话，不要使用全局保存的 Ormer 对象。
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
	} else if skey := ctx.Input.Header(models.WX_HEADER_SKEY); skey != "" {
		// if set skey, chk session available
		sess, err := models.WeappGetSession(skey)
		glog.V(4).Infof("get skey %s sess %#v", skey, sess)
		if err != nil {
			http.Error(ctx.ResponseWriter,
				fmt.Sprintf(`{"%s":true}`, models.WX_SESSION_MAGIC_ID), 401)
			return
		}
		// for weapp session
		op.User = sess.User
		op.Token = sess.Token
	}
}
