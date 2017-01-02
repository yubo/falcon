package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "GetLogin",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:AuthController"],
		beego.ControllerComments{
			Method: "PostLogin",
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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "CreateToken",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokensCnt",
			Router: `/cnt/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokens",
			Router: `/search/:query`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "UpdateToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "DeleteToken",
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
