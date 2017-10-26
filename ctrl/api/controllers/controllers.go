package controllers

import (
	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/api/models"
)

type Search struct {
	Name        string
	Placeholder string
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) SendMsg(code int, msg interface{}) {
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.Data["json"] = models.PruneNilMsg(msg)
	c.ServeJSON()
}

func (c *BaseController) SendText(code int, msg []byte) {
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.Ctx.Output.Header("Content-Type", "text/plain")
	c.Ctx.Output.Body(msg)
}

func totalObj(n int64) models.Total {
	return models.Total{Total: n}
}

func idObj(n int64) models.Id {
	return models.Id{Id: n}
}

func statsObj(success, failure int64) models.Stats {
	return models.Stats{Success: success, Failure: failure}
}
