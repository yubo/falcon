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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "CreateExpression",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpressionsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpressions",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpression",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "UpdateExpression",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "DeleteExpression",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
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
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHosts",
			Router: `/search`,
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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetzTreeNodes",
			Router: `/zTreeNodes`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagHosts",
			Router: `/tag/hosts`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagRoleUsers",
			Router: `/tag/role/users`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagRoleTokens",
			Router: `/tag/role/tokens`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagHost",
			Router: `/tag/host`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagRoleUser",
			Router: `/tag/role/user`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagRoleToken",
			Router: `/tag/role/token`,
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
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRoles",
			Router: `/search`,
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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "CreateRule",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "GetRulesCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "GetRules",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "GetRule",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "UpdateRule",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:RuleController"],
		beego.ControllerComments{
			Method: "DeleteRule",
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
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTags",
			Router: `/search`,
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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "CreateTeamUsers",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeamsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeams",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeam",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "UpdateTeam",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TeamController"],
		beego.ControllerComments{
			Method: "DeleteTeam",
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
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokens",
			Router: `/search`,
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

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "CreateTrigger",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "GetTriggersCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "GetTriggers",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "GetTrigger",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "UpdateTrigger",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:TriggerController"],
		beego.ControllerComments{
			Method: "DeleteTrigger",
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
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsers",
			Router: `/search`,
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
