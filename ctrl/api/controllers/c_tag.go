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
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Tags
type TagController struct {
	BaseController
}

// @Title CreateTag
// @Description create tags
// @Param	body	body 	models.Tag	true	"body for tag content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *TagController) CreateTag() {
	var tag models.Tag
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if id, err := me.AddTag(&tag); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetTagsCnt
// @Description get Tags number
// @Param   query     query   string  false       "tag name"
// @Success 200  {total:int} tag total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TagController) GetTagsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTagsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTags
// @Description get all Tags
// @Param   query     query   string  false       "tag name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 [object] []models.Tag
// @Failure 403 string error
// @router /search [get]
func (c *TagController) GetTags() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	tags, err := me.GetTags(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, tags)
	}
}

// @Title Get
// @Description get tag by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Tag
// @Failure 403 string error
// @router /:id [get]
func (c *TagController) GetTag() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		if tag, err := me.GetTag(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, tag)
		}
	}
}

// @Title UpdateTag
// @Description update the tag
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Tag	true		"body for tag content"
// @Success 200 {object} models.Tag
// @Failure 403 string error
// @router /:id [put]
func (c *TagController) UpdateTag() {
	var tag models.Tag
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if u, err := me.UpdateTag(id, &tag); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteTag
// @Description delete the tag
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 string "delete success!"
// @Failure 403 string error
// @router /:id [delete]
func (c *TagController) DeleteTag() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteTag(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
