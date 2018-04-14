/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"

	"github.com/yubo/falcon/modules/ctrl/api/models"
)

type PluginController struct {
	BaseController
}

// @Title GetTagPluginCnt
// @Description get Tag-plugin number
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"include child tag(default:true)"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/plugindir/cnt [get]
func (c *PluginController) GetTagPluginCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_F_R_TOKEN, tagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	n, err := op.GetPluginDirCnt(tagId, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title Get plugin dir
// @Description get all Template
// @Param	tag_id	query	int	true	"tag id"
// @Param	deep	query   bool	false	"include child tag(default:true)"
// @Param	limit	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.PluginDirGet plugin info
// @Failure 400 string error
// @router /tag/plugindir/search [get]
func (c *PluginController) GetPluginDir() {
	tagId, _ := c.GetInt64("tag_id", 1)
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_F_R_TOKEN, tagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	ret, err := op.GetPluginDir(tagId, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag template relation
// @Description create tag/template relation
// @Param	body	body	models.PluginDirPost	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/plugindir [post]
func (c *PluginController) CreatePluginDir() {
	var input models.PluginDirPost

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := op.Access(models.SYS_F_O_TOKEN, input.TagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if id, err := op.CreatePluginDir(&input); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title Delete tag plugin
// @Description delete the plugin
// @Param	tag_id	query 	int	true	"tag id"
// @Param	id	query 	int	true	"The id you want to delete"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/plugindir [delete]
func (c *PluginController) DeletePluginDir() {
	id, _ := c.GetInt64("id")
	tagId, _ := c.GetInt64("tag_id", 1)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_F_O_TOKEN, tagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	n, err := op.DeletePluginDir(tagId, id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}
