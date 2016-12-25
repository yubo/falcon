package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "LoginGet",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "LoginPost",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Callback",
			Router: `/callback/:module`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "CreateHost",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHostsCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHosts",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "UpdateHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "DeleteHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "CreateRole",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRolesCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRoles",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "UpdateRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "DeleteRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "CreateScope",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "GetScopesCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "GetScopes",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "GetScope",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "UpdateScope",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ScopeController"],
		beego.ControllerComments{
			Method: "DeleteScope",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "CreateSystem",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "GetSystemsCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "GetSystems",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "GetSystem",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "UpdateSystem",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:SystemController"],
		beego.ControllerComments{
			Method: "DeleteSystem",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "CreateTag",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTagsCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTags",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "UpdateTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "DeleteTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "CreateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsersCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsers",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "DeleteUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

}
