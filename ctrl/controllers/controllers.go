package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/wetalk/modules/utils"
	"github.com/yubo/falcon/ctrl/models"
)

type Link struct {
	Text     string
	Url      string
	Disabled bool
}

type Links struct {
	Text     string
	Url      string
	Disabled bool
	SubLinks []Link
}

type Search struct {
	Name        string
	Url         string
	Placeholder string
}

type BaseController struct {
	beego.Controller
}

type MainController struct {
	BaseController
}

const (
	HEAD_LINK_IDX_REL = iota
	HEAD_LINK_IDX_META
	HEAD_LINK_IDX_SETTINGS
	HEAD_LINK_IDX_HELP
	HEAD_LINK_IDX_LOGOUT
)

var (
	notLoginheadLinks = []Links{
		{Text: "help", Url: "#",
			SubLinks: []Link{
				{Text: "Doc", Url: "/doc"},
				{},
				{Text: "About Falcon", Url: "/about"}}},
		{Text: "[login]", Url: "/login"},
	}
	headLinks = []Links{
		HEAD_LINK_IDX_REL: {Text: "Relation", Url: "#",
			SubLinks: []Link{
				{Text: "Tag Host", Url: "/rel/tag/host"},
				{Text: "Tag Role User", Url: "/rel/tag/role/user"},
				{Text: "Tag Role Token", Url: "/rel/tag/role/token"},
				{Text: "Tag Rule Trigger", Url: "/rel/tag/rule/trigger"}}},
		HEAD_LINK_IDX_META: {Text: "Meta", Url: "#",
			SubLinks: []Link{
				{Text: "Tag", Url: "/tag"},
				{Text: "Host", Url: "/host"},
				{},
				{Text: "Role", Url: "/role"},
				{Text: "User", Url: "/user"},
				{Text: "Token", Url: "/token"},
				{},
				{Text: "Rule", Url: "/rule", Disabled: true},
				{Text: "Trigger", Url: "/trigger", Disabled: true}}},
		HEAD_LINK_IDX_SETTINGS: {Text: "Settings", Url: "#",
			SubLinks: []Link{
				{Text: "Global", Url: "/settings/config/global"},
				{Text: "Profile", Url: "/settings/profile"},
				{Text: "About Me", Url: "/settings/aboutme"}}},
		HEAD_LINK_IDX_HELP: {Text: "help", Url: "#",
			SubLinks: []Link{
				{Text: "doc", Url: "/doc"},
				{},
				{Text: "About Falcon", Url: "/about"}}},
		HEAD_LINK_IDX_LOGOUT: {Text: "[logout]", Url: "/logout"},
	}
)

func init() {
	// The hookfuncs will run in beego.Run()
	beego.AddAPPStartHook(start)
}

func start() (err error) {
	if beego.BConfig.RunMode == "dev" {
		headLinks[2].SubLinks = append(headLinks[2].SubLinks, Link{Text: "Debug", Url: "/settings/debug"})
	}
	return
}

func (c *BaseController) SendMsg(code int, msg string) {
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	//c.Ctx.ResponseWriter.WriteHeader(code)
	beego.Debug(c.Data["json"])
	c.ServeJSON()
}

func (c *BaseController) SendObj(code int, obj interface{}) {
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"data": obj,
	}
	c.ServeJSON()
}

func (c *BaseController) SetPaginator(per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c.Ctx.Request, per, nums)
	c.Data["paginator"] = p
	return p
}

func (c *BaseController) PrepareEnv(links []Link, curLink string) {
	var ok bool

	if c.Data["Me"], ok = c.Ctx.Input.GetData("me").(*models.User); !ok {
		c.Data["HeadLinks"] = notLoginheadLinks
		return
	}
	c.Data["HeadLinks"] = headLinks
	c.Data["Links"] = links
	c.Data["CurLink"] = curLink

}

func (c *MainController) Get() {
	c.PrepareEnv(nil, "")
	c.TplName = "index.tpl"
}

func (c *MainController) About() {
	c.PrepareEnv(nil, "")
	c.TplName = "about.tpl"
}
