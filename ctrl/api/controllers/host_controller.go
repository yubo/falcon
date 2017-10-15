/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
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
// @Success 200 {object} models.Id id
// @Failure 400 string error
// @router / [post]
func (c *HostController) CreateHost() {
	var host models.Host

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)

	if id, err := op.AddHost(&host); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetHostsCnt
// @Description get Hosts number
// @Param   query     query   string  false       "host name"
// @Success 200 {object} models.Total  host total number
// @Failure 400 string error
// @router /cnt [get]
func (c *HostController) GetHostsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetHostsCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetHosts
// @Description get all Hosts
// @Param	query	query   string  false       "host name"
// @Param	limit	query   int     false       "limit page number"
// @Param	offset	query   int     false       "offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 400 string error
// @router /search [get]
func (c *HostController) GetHosts() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	hosts, err := op.GetHosts(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, hosts)
	}
}

// @Title GetHost
// @Description get host by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Host host info
// @Failure 400 string error
// @router /:id [get]
func (c *HostController) GetHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		host, err := op.GetHost(id)
		if err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, host)
		}
	}
}

// @Title UpdateHost
// @Description update the host
// @Param	body		body 	models.Host	true	"body for host content"
// @Success 200 {object} models.Host host info
// @Failure 400 string error
// @router / [put]
func (c *HostController) UpdateHost() {
	host := &models.Host{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, host)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if host, err = op.GetHost(host.Id); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	// overlay entry
	o := *host
	json.Unmarshal(c.Ctx.Input.RequestBody, &o)

	if host, err = op.UpdateHost(&o); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, host)
	}
}

// @Title DeleteHost
// @Description delete the host
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /:id [delete]
func (c *HostController) DeleteHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteHost(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}

// @Title DeleteHosts
// @Description delete the hosts
// @Param	body	[]int	true	"The []id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router / [delete]
func (c *HostController) DeleteHosts() {
	var (
		ids              []int64
		success, failure int64
	)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	for _, id := range ids {
		if _ = op.DeleteHost(id); err != nil {
			failure++
		} else {
			success++
		}
	}
	c.SendMsg(200, statsObj(success, failure))
}
