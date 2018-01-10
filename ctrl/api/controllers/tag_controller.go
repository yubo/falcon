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

// Operations about Tags
type TagController struct {
	BaseController
}

// @Title CreateTag
// @Description create tags
// @Param	body	body 	models.TagCreate	true	"body for tag content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *TagController) CreateTag() {
	var tag models.TagCreate
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)
	tag.Type = 0

	// TODO: check parent exist/acl
	if _, err := op.AccessByStr(models.SYS_O_TOKEN, models.TagParent(tag.Name),
		false); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if id, err := op.CreateTag(&tag); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTagsCnt
// @Description get Tags number
// @Param	query	query	string	false	"tag name"
// @Success 200	{object}	models.Total tag total number
// @Failure 400	string error
// @router /cnt [get]
func (c *TagController) GetTagsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTagsCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTags
// @Description get all Tags
// @Param	query	query	string	false	"tag name"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200	{object}	[]models.Tag	tags info
// @Failure 400	string	error
// @router /search [get]
func (c *TagController) GetTags() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	tags, err := op.GetTags(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, tags)
	}
}

// @Title Get
// @Description get tag by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Tag tag info
// @Failure 400 string error
// @router /:id [get]
func (c *TagController) GetTag() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if tag, err := op.GetTag(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, tag)
		}
	}
}

/*
Title UpdateTag
Description update the tag
Param	id		path 	string	true		"The id you want to update"
Param	body		body 	models.Tag	true		"body for tag content"
Success 200 {object} models.Tag tag info
Failure 400 string error
*/
/*
func (c *TagController) UpdateTag() {
	var tag models.Tag
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if _, err = op.AccessByStr(models.SYS_O_TOKEN, models.TagParent(tag.Name),
		true); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if u, err := op.UpdateTag(id, &tag); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}
*/

// @Title DeleteTag
// @Description delete the tag
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /:id [delete]
func (c *TagController) DeleteTag() {
	var (
		tag *models.Tag
		id  int64
		err error
		op  *models.Operator
	)

	id, err = c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ = c.Ctx.Input.GetData("op").(*models.Operator)

	if tag, err = op.GetTag(id); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if _, err = op.AccessByStr(models.SYS_O_TOKEN,
		models.TagParent(tag.Name), false); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	err = op.DeleteTag(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "delete success!")
	}
}

/*******************************************************************************
 ************************ relation *********************************************
 ******************************************************************************/

// @Title Get tree's node
// @Description get node and it's children
// @Param	id	query 	int64	false	"tag id default root(1)"
// @Param	depth	query   int	false	"depth levels default -1(no limit)"
// @Success 200 {object} models.TreeNode all nodes of the tree(read)
// @Failure 400 string error
// @router /node [get]
func (c *TagController) GetTreeNode() {
	var ret *models.TreeNode

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	id, _ := c.GetInt64("id", 1)
	depth, _ := c.GetInt("depth", -1)
	direct := false

	if depth < 0 {
		depth = 100
	}

	if err := op.Access(models.SYS_R_TOKEN, id); err == nil {
		direct = true
	}

	ret = op.GetTreeNode(id, depth, direct)
	if ret == nil {
		c.SendMsg(400, "tree empty or Permission denied")
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Get tags(operate)
// @Description get has operate token tags
// @Param	expand	query   bool	false	"include child tag(default:false)"
// @Success 200 {object} []int64 all ids of the node that can be operated
// @Failure 400 string error
// @router /operate [get]
func (c *TagController) GetOpTag() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	expand, _ := c.GetBool("expand", false)
	ret, _ := op.GetOpTag(expand)
	c.SendMsg(200, ret)
}

// @Title Get tags(read)
// @Description get has read token tags
// @Param	expand	query   bool	false	"include child tag(default:false)"
// @Success 200 {object} []int64 all ids of the node that can be read
// @Failure 400 string error
// @router /read [get]
func (c *TagController) GetReadTag() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	expand, _ := c.GetBool("expand", false)
	ret, _ := op.GetReadTag(expand)
	c.SendMsg(200, ret)
}
