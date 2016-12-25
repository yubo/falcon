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

// Operations about Tags
type TagController struct {
	BaseController
}

// @Title CreateTag
// @Description create tags
// @Param	body	body 	models.Tag	true	"body for tag content"
// @Success {code:200, data:int} models.Tag.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *TagController) CreateTag() {
	var tag models.Tag
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)
	if id, err := me.AddTag(&tag); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, id)
	}
}

// @Title GetTagsCnt
// @Description get Tags number
// @Success {code:200, data:int} tag number
// @Failure {code:int, msg:string}
// @router /cnt/:query [get]
func (c *TagController) GetTagsCnt() {
	query := strings.TrimSpace(c.GetString(":query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTagsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, cnt)
	}
}

// @Title GetTags
// @Description get all Tags
// @Param   per       query   int  false       "per page number"
// @Param   offset    query   int  false       "offset  number"
// @Success {code:200, data:object} models.Tag
// @Failure {code:int, msg:string}
// @router /search/:query [get]
func (c *TagController) GetTags() {
	query := strings.TrimSpace(c.GetString(":query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	tags, err := me.GetTags(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, tags)
	}
}

// @Title Get
// @Description get tag by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Tag
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *TagController) GetTag() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		tag, err := models.GetTag(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendObj(200, tag)
		}
	}
}

// @Title UpdateTag
// @Description update the tag
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Tag	true		"body for tag content"
// @Success {code:200, data:object} models.Tag
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *TagController) UpdateTag() {
	var tag models.Tag
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	id, err := c.GetInt(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if u, err := me.UpdateTag(id, &tag); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

// @Title DeleteTag
// @Description delete the tag
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TagController) DeleteTag() {
	id, err := c.GetInt(":id")
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

	c.SendObj(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetTag() {
	var tags []*models.Tag

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryTags(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&tags)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Tags"] = tags
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "/tag"}

	c.TplName = "tag/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) EditTag() {
	var tag *models.Tag

	id, err := c.GetInt(":id")
	if err != nil {
		goto out
	}

	tag, err = models.GetTag(id)
	if err != nil {
		goto out
	}

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Tag"] = tag
	c.Data["H1"] = "edit tag"
	c.Data["Method"] = "put"
	c.TplName = "tag/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

func (c *MainController) AddTag() {

	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["Method"] = "post"
	c.Data["H1"] = "add tag"
	c.TplName = "tag/edit.tpl"
}
