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

// Operations about Hosts
type HostController struct {
	BaseController
}

// @Title CreateHost
// @Description create hosts
// @Param	body	body 	models.Host	true	"body for host content"
// @Success {code:200, data:int} models.Host.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *HostController) CreatHost() {
	var host models.Host
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)
	id, err := models.AddHost(&host)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetHostsCnt
// @Description get Hosts number
// @Success {code:200, data:int} host number
// @Failure {code:int, msg:string}
// @router /cnt/:query [get]
func (c *HostController) GetHostsCnt() {
	query := strings.TrimSpace(c.GetString(":query"))

	cnt, err := models.GetHostsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetHosts
// @Description get all Hosts
// @Param   per       query   int  false       "per page number"
// @Param   offset    query   int  false       "offset  number"
// @Success {code:200, data:object} models.Host
// @Failure {code:int, msg:string}
// @router /search/:query [get]
func (c *HostController) GetHosts() {
	query := strings.TrimSpace(c.GetString(":query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)

	hosts, err := models.GetHosts(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, hosts)
	}
}

// @Title GetHost
// @Description get host by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Host
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *HostController) GetHost() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		host, err := models.GetHost(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, host)
		}
	}
}

// @Title UpdateHost
// @Description update the host
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Host	true		"body for host content"
// @Success {code:200, data:object} models.Host
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *HostController) UpdateHost() {
	var host models.Host

	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &host)

	if u, err := models.UpdateHost(id, &host); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title DeleteHost
// @Description delete the host
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *HostController) DeleteHost() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	err = models.DeleteHost(id)
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
func (c *MainController) GetHost() {
	var hosts []*models.Host

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)

	qs := models.QueryHosts(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&hosts)
	if err != nil {
		goto out
	}

	c.PrepareEnv()
	c.Data["Hosts"] = hosts
	c.Data["Query"] = query
	c.Data["Search"] = hostSearch

	c.TplName = "host/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditHost() {
	var host *models.Host

	id, err := c.GetInt(":id")
	if err != nil {
		goto out
	}

	host, err = models.GetHost(id)
	if err != nil {
		goto out
	}

	c.PrepareEnv()
	c.Data["Host"] = host
	c.Data["H1"] = "edit host"
	c.Data["Method"] = "put"
	c.TplName = "host/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddHost() {

	c.PrepareEnv()
	c.Data["Method"] = "post"
	c.Data["H1"] = "add host"
	c.TplName = "host/edit.tpl"
}
