package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "UpdateConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetEtcdMap",
			Router: `/config/list/etcd`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetModuleMap",
			Router: `/config/list/module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetDebugAction",
			Router: `/debug/:action`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetExpansion",
			Router: `/expansion/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "SetExpansion",
			Router: `/expansion/:module`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetOnline",
			Router: `/online/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "CreateAggreator",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "CreateAggreator0",
			Router: `/0`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "DeleteAggreator0",
			Router: `/0`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "DeleteAggreator",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "GetAggreator",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "UpdateAggreator",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "GetAggreatorsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "GetAggreatorsCnt0",
			Router: `/cnt/0`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "GetAggreators",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AggreatorController"],
		beego.ControllerComments{
			Method: "GetAggreators0",
			Router: `/search/0`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Callback",
			Router: `/callback/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Info",
			Router: `/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "PostLogin",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Authorize",
			Router: `/login/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Modules",
			Router: `/modules`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateGraph",
			Router: `/graph`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "DeleteGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetGraphByScreen",
			Router: `/graph/screen/:screen_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateGraphs",
			Router: `/graphs`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateScreen",
			Router: `/screen`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "DeleteScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetScreenByPid",
			Router: `/screen/pid/:pid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateTmpGraph",
			Router: `/tmpgraph`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetTmpGraph",
			Router: `/tmpgraph/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "CreateExpression",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "DeleteExpression0",
			Router: `/0`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "DeleteExpression",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpressionAction",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "UpdateExpressionAction",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpressionsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "PauseExpression",
			Router: `/pause`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:ExpressionController"],
		beego.ControllerComments{
			Method: "GetExpressions",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetCounterData",
			Router: `/counter_data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetEndpoint",
			Router: `/endpoint`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetEndpointCounter",
			Router: `/endpoint_counter`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "CreateHost",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "UpdateHost",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DeleteHosts",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DeleteHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHostsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHosts",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "UpdateMatter",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "CreateClaim",
			Router: `/claim`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "GetMattersCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "GetEventCnt",
			Router: `/event/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "GetEvents",
			Router: `/event/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MatterController"],
		beego.ControllerComments{
			Method: "GetMatters",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MetricController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MetricController"],
		beego.ControllerComments{
			Method: "GetMetricsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MetricController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MetricController"],
		beego.ControllerComments{
			Method: "GetMetrics",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "CreateMockcfg",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "GetMockcfg",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "UpdateMockcfg",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "DeleteMockcfg",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "GetMockcfgsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:MockcfgController"],
		beego.ControllerComments{
			Method: "GetMockcfgs",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/ctrl`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"],
		beego.ControllerComments{
			Method: "GetTagHostCnt",
			Router: `/rel/tag/host/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:PubController"],
		beego.ControllerComments{
			Method: "GetTagHost",
			Router: `/rel/tag/host/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetNode",
			Router: `/node`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetOpTag",
			Router: `/operate/tag`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetReadTag",
			Router: `/read/tag`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagHost",
			Router: `/tag/host`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagHost",
			Router: `/tag/host`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagHostCnt",
			Router: `/tag/host/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagHost",
			Router: `/tag/host/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagHosts",
			Router: `/tag/hosts`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagHosts",
			Router: `/tag/hosts`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreatePluginDir",
			Router: `/tag/plugindir`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DeletePluginDir",
			Router: `/tag/plugindir`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagPluginCnt",
			Router: `/tag/plugindir/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetPluginDir",
			Router: `/tag/plugindir/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagRoleToken",
			Router: `/tag/role/token`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagRoleToken",
			Router: `/tag/role/token`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagRoleTokenCnt",
			Router: `/tag/role/token/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagRoleToken",
			Router: `/tag/role/token/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagRoleUser",
			Router: `/tag/role/user`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagRoleUser",
			Router: `/tag/role/user`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagRoleUserCnt",
			Router: `/tag/role/user/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagRoleUser",
			Router: `/tag/role/user/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagTpl",
			Router: `/tag/template`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagTpl",
			Router: `/tag/template`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagTpl0",
			Router: `/tag/template/0`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagTpl0",
			Router: `/tag/template/0`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagTplCnt",
			Router: `/tag/template/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagTplCnt0",
			Router: `/tag/template/cnt/0`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagTpl",
			Router: `/tag/template/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTagTpl0",
			Router: `/tag/template/search/0`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "DelTagTpls",
			Router: `/tag/templates`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "CreateTagTpls",
			Router: `/tag/templates`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RelController"],
		beego.ControllerComments{
			Method: "GetTree",
			Router: `/tree`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "CreateRole",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "UpdateRole",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "DeleteRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRolesCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRoles",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetLogsCnt",
			Router: `/log/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetLogs",
			Router: `/log/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/profile`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/profile`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "CreateStrategy",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "GetStrategy",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "UpdateStrategy",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "DeleteStrategy",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "GetStrategysCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:StrategyController"],
		beego.ControllerComments{
			Method: "GetStrategys",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "CreateTag",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "DeleteTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTagsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTags",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "CreateTeam",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "DeleteTeam",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeam",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "UpdateTeam",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "UpdateMember",
			Router: `/:id/member`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeamsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetMember",
			Router: `/member`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TeamController"],
		beego.ControllerComments{
			Method: "GetTeams",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "CreateTemplate",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "GetTemplate",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "UpdateTemplate",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "DeleteTemplate",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "GetTemplatesCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TemplateController"],
		beego.ControllerComments{
			Method: "GetTemplates",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "CreateToken",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "UpdateToken",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "DeleteToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokensCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokens",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "CreateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "DeleteUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetBindedUsers",
			Router: `/binded/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsersCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsers",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "UnBindUser",
			Router: `/unbind/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Authqr",
			Router: `/authqr`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Bindqr",
			Router: `/bindqr`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Openid",
			Router: `/openid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Task",
			Router: `/task`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Taskack",
			Router: `/taskack`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "TemplateMessage",
			Router: `/templateMessage`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "TestRequest",
			Router: `/testRequest`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Tunnel",
			Router: `/tunnel`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Upload",
			Router: `/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
