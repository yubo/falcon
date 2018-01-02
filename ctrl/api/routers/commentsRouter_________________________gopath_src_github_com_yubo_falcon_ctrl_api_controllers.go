package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:ActionController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:ActionController"],
		beego.ControllerComments{
			Method: "CreateActionTrigger",
			Router: `/trigger`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "UpdateConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetEtcdMap",
			Router: `/config/list/etcd`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetModuleMap",
			Router: `/config/list/module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetDebugAction",
			Router: `/debug/:action`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetExpansion",
			Router: `/expansion/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "SetExpansion",
			Router: `/expansion/:module`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AdminController"],
		beego.ControllerComments{
			Method: "GetOnline",
			Router: `/online/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Callback",
			Router: `/callback/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Info",
			Router: `/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "PostLogin",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Authorize",
			Router: `/login/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Modules",
			Router: `/modules`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateGraph",
			Router: `/graph`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "DeleteGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetGraph",
			Router: `/graph/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetGraphByScreen",
			Router: `/graph/screen/:screen_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateGraphs",
			Router: `/graphs`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateScreen",
			Router: `/screen`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "UpdateScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "DeleteScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetScreen",
			Router: `/screen/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetScreenByPid",
			Router: `/screen/pid/:pid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "CreateTmpGraph",
			Router: `/tmpgraph`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:DashboardController"],
		beego.ControllerComments{
			Method: "GetTmpGraph",
			Router: `/tmpgraph/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "CreateEventTrigger",
			Router: `/trigger`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "UpdateEventTrigger",
			Router: `/trigger`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "DeleteEventTrigger",
			Router: `/trigger`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "CloneEventTrigger",
			Router: `/trigger/clone`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "GetEventTriggersCnt",
			Router: `/trigger/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:EventController"],
		beego.ControllerComments{
			Method: "GetEventTriggers",
			Router: `/trigger/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetCounterData",
			Router: `/counter_data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetEndpoint",
			Router: `/endpoint`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:GraphController"],
		beego.ControllerComments{
			Method: "GetEndpointCounter",
			Router: `/endpoint_counter`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "UpdateHost",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DeleteHosts",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "CreateHost",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DeleteHost",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHostsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetHosts",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "CreateTagHost",
			Router: `/tag`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "CreateTagHosts",
			Router: `/tag`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DelTagHost",
			Router: `/tag`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "DelTagHosts",
			Router: `/tag`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetTagHostCnt",
			Router: `/tag/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:HostController"],
		beego.ControllerComments{
			Method: "GetTagHost",
			Router: `/tag/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:MetricController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:MetricController"],
		beego.ControllerComments{
			Method: "GetMetricsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:MetricController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:MetricController"],
		beego.ControllerComments{
			Method: "GetMetrics",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"],
		beego.ControllerComments{
			Method: "CreatePluginDir",
			Router: `/tag/plugindir`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"],
		beego.ControllerComments{
			Method: "DeletePluginDir",
			Router: `/tag/plugindir`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"],
		beego.ControllerComments{
			Method: "GetTagPluginCnt",
			Router: `/tag/plugindir/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PluginController"],
		beego.ControllerComments{
			Method: "GetPluginDir",
			Router: `/tag/plugindir/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PubController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:PubController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/ctrl`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "CreateRole",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "UpdateRole",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "DeleteRole",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRolesCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:RoleController"],
		beego.ControllerComments{
			Method: "GetRoles",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/config/:module`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetLogsCnt",
			Router: `/log/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetLogs",
			Router: `/log/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/profile`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:SetController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/profile`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "CreateTag",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "DeleteTag",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTagsCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTreeNode",
			Router: `/node`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetOpTag",
			Router: `/operate`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetReadTag",
			Router: `/read`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TagController"],
		beego.ControllerComments{
			Method: "GetTags",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "CreateToken",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "UpdateToken",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "DeleteToken",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokensCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "CreateTagRolesTokens",
			Router: `/m/tag/roles`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "DelTagRoleToken",
			Router: `/m/tag/roles`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTokens",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTagRoleTokenCnt",
			Router: `/tag/role/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:TokenController"],
		beego.ControllerComments{
			Method: "GetTagRoleToken",
			Router: `/tag/role/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "DeleteUsers",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "CreateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "DeleteUser",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetBindedUsers",
			Router: `/binded/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsersCnt",
			Router: `/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUsers",
			Router: `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "CreateTagRolesUsers",
			Router: `/tag/role`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "DelTagRolesUsers",
			Router: `/tag/role`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetTagRoleUserCnt",
			Router: `/tag/role/cnt`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetTagRoleUser",
			Router: `/tag/role/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:UserController"],
		beego.ControllerComments{
			Method: "UnBindUser",
			Router: `/unbind/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Authqr",
			Router: `/authqr`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Bindqr",
			Router: `/bindqr`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Openid",
			Router: `/openid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Task",
			Router: `/task`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Taskack",
			Router: `/taskack`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "TemplateMessage",
			Router: `/templateMessage`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "TestRequest",
			Router: `/testRequest`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Tunnel",
			Router: `/tunnel`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"] = append(beego.GlobalControllerRouter["github.com/yubo/falcon/cmd/vendor/github.com/yubo/falcon/ctrl/api/controllers:WeappController"],
		beego.ControllerComments{
			Method: "Upload",
			Router: `/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
