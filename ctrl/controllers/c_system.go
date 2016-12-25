/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Systems
type SystemController struct {
	BaseController
}

// @Title CreateSystem
// @Description create systems
// @Param	body	body 	models.System	true	"body for system content"
// @Success {code:200, data:int} models.System.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *SystemController) CreateSystem() {
	var system models.System
	json.Unmarshal(c.Ctx.Input.RequestBody, &system)
	id, err := models.AddSystem(&system)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetSystemsCnt
// @Description get Systems number
// @Success {code:200, data:int} system number
// @Failure {code:int, msg:string}
// @router /cnt/:query [get]
func (c *SystemController) GetSystemsCnt() {
	query := strings.TrimSpace(c.GetString(":query"))

	cnt, err := models.GetSystemsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetSystems
// @Description get all Systems
// @Param   per       query   int  false       "per page number"
// @Param   offset    query   int  false       "offset  number"
// @Success {code:200, data:object} models.System
// @Failure {code:int, msg:string}
// @router /search/:query [get]
func (c *SystemController) GetSystems() {
	query := strings.TrimSpace(c.GetString(":query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)

	systems, err := models.GetSystems(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, systems)
	}
}

// @Title Get
// @Description get system by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.System
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *SystemController) GetSystem() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		system, err := models.GetSystem(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, system)
		}
	}
}

// @Title UpdateSystem
// @Description update the system
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.System	true		"body for system content"
// @Success {code:200, data:object} models.System
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *SystemController) UpdateSystem() {
	var system models.System

	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &system)

	if u, err := models.UpdateSystem(id, &system); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title DeleteSystem
// @Description delete the system
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *SystemController) DeleteSystem() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	err = models.DeleteSystem(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendObj(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetSystem() {
	var systems []*models.System

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)

	qs := models.QuerySystems(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&systems)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Systems"] = systems
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "/system"}

	c.TplName = "system/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditSystem() {
	var system *models.System

	id, err := c.GetInt(":id")
	if err != nil {
		goto out
	}

	system, err = models.GetSystem(id)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["System"] = system
	c.Data["H1"] = "edit system"
	c.Data["Method"] = "put"
	c.TplName = "system/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddSystem() {

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Method"] = "post"
	c.Data["H1"] = "add system"
	c.TplName = "system/edit.tpl"
}
