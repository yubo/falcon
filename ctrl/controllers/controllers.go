package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/wetalk/modules/utils"
	"github.com/yubo/falcon/ctrl/models"
)

type Link struct {
	Text string
	Url  string
}

type Search struct {
	Name string
	Url  string
}

type BaseController struct {
	beego.Controller
}

type MainController struct {
	BaseController
}

var (
	settingsLinks = []Link{
		{"Profile", "/settings/profile"},
		{"About", "/settings/about"},
	}
)

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

func (c *BaseController) PrepareEnv() {
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
}

func (c *MainController) Get() {
	c.PrepareEnv()
	c.TplName = "index.tpl"
}
