/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Hosts
type HostController struct {
	BaseController
}

// @Title CreateHost
// @Description create hosts
// @Param	body	body 	models.Host	true	"body for host content"
// @Success 200 {int} models.Host.Id
// @Failure 403 string error
// @router / [post]
func (c *HostController) CreateHost() {
	var host models.Host

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)
	host.Id = 0

	if id, err := me.AddHost(&host); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetHostsCnt
// @Description get Hosts number
// @Param   query     query   string  false       "host name"
// @Success 200 {total:int}  host total number
// @Failure 403 string error
// @router /cnt [get]
func (c *HostController) GetHostsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetHostsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetHosts
// @Description get all Hosts
// @Param   query     query   string  false       "host name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 [object] []models.Host
// @Failure 403 string error
// @router /search [get]
func (c *HostController) GetHosts() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	hosts, err := me.GetHosts(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, hosts)
	}
}

// @Title GetHost
// @Description get host by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Host
// @Failure 403 string error
// @router /:id [get]
func (c *HostController) GetHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		host, err := me.GetHost(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, host)
		}
	}
}

// @Title UpdateHost
// @Description update the host
// @Param	id		path 	string	true	"The id you want to update"
// @Param	body		body 	models.Host	true	"body for host content"
// @Success 200 {object} models.Host
// @Failure 403 string error
// @router /:id [put]
func (c *HostController) UpdateHost() {
	var host models.Host
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &host)

	if u, err := me.UpdateHost(id, &host); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteHost
// @Description delete the host
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 string "delete success!"
// @Failure 403 string error
// @router /:id [delete]
func (c *HostController) DeleteHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteHost(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}
