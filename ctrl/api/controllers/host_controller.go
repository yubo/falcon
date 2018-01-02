/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Hosts
type HostController struct {
	BaseController
}

// @Title CreateHost
// @Description create hosts
// @Param	body	body 	models.HostCreate	true	"body for host content"
// @Success 200 {object} models.Id id
// @Failure 400 string error
// @router / [post]
func (c *HostController) CreateHost() {
	var host models.HostCreate

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)

	if id, err := op.CreateHost(&host); err != nil {
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
// @Param	body	body 	models.HostUpdate	true	"body for host content"
// @Success 200 {object} models.Host host info
// @Failure 400 string error
// @router / [put]
func (c *HostController) UpdateHost() {
	input := models.HostUpdate{}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	p, err := op.GetHost(input.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	host := *p
	falcon.Override(&host, &input)

	if ret, err := op.UpdateHost(&host); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
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
// @Param body	body	[]int	true	"The []id you want to delete"
// @Success 200 {string} "delete success!"
// @Success 200 {object} models.Stats api call result
// @Failure 400 string error
// @router / [delete]
func (c *HostController) DeleteHosts() {
	var (
		ids              []int64
		errs             []string
		success, failure int64
	)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	for _, id := range ids {
		if err := op.DeleteHost(id); err != nil {
			errs = append(errs, err.Error())
			failure++
		} else {
			success++
		}
	}
	c.SendMsg(200, statsObj(success, failure, errs))
}

/*******************************************************************************
 ************************ tag host *********************************************
 ******************************************************************************/

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	tag_id	query   int	false	"tag id"
// @Param	query	query   string  false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/cnt [get]
func (c *HostController) GetTagHostCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagHostCnt(tagId, query, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetHost
// @Description get all Host
// @Param	tag_id	query   int	false	"tag id"
// @Param	query	query	string	false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagHostApiGet tag host info
// @Failure 400 string error
// @router /tag/search [get]
func (c *HostController) GetTagHost() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := c.GetString("query")
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagHost(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag host relation
// @Description create tag/host relation
// @Param	body	body	models.TagHostApiAdd	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag [post]
func (c *HostController) CreateTagHost() {
	var rel models.TagHostApiAdd

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if !op.IsAdmin() {
		if err := op.Access(models.SYS_O_TOKEN,
			rel.TagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}

		if err := op.Access(models.SYS_O_TOKEN,
			rel.SrcTagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}

		if cnt, err := op.ChkTagHostCnt(rel.SrcTagId,
			[]int64{rel.HostId}); err != nil || cnt != 1 {
			c.SendMsg(403, falcon.EPERM)
			return
		}
	}

	if _, err := op.CreateTagHost(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(1))
	}
}

// @Title create tag host relation
// @Description create tag/hosts relation
// @Param	body	body	models.TagHostsApiAdd	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag [post]
func (c *HostController) CreateTagHosts() {
	var rel models.TagHostsApiAdd

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rel)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if !op.IsAdmin() {
		if err := op.Access(models.SYS_O_TOKEN,
			rel.TagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}

		if err := op.Access(models.SYS_O_TOKEN,
			rel.SrcTagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}

		if cnt, err := op.ChkTagHostCnt(rel.SrcTagId,
			rel.HostIds); err != nil || cnt != int64(len(rel.HostIds)) {
			c.SendMsg(403, falcon.EPERM)
			return
		}

	}

	if n, err := op.CreateTagHosts(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag host relation
// @Description delete tag/host relation
// @Param	body	body	models.TagHostApiDel	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag [delete]
func (c *HostController) DelTagHost() {
	var rel models.TagHostApiDel

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagHost(&rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag host relation
// @Description delete tag/hosts relation
// @Param	body	body	models.TagHostsApiDel	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag [delete]
func (c *HostController) DelTagHosts() {
	var rel models.TagHostsApiDel

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagHosts(&rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}
