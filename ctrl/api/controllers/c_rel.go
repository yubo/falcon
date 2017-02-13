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

var (
	rels = map[string]bool{
		"tag_host":         true,
		"tag_role_user":    true,
		"tag_role_token":   true,
		"tag_rule_trigger": true,
	}
)

// Operations about Relations
type RelController struct {
	BaseController
}

// @Title Get vue tag tree
// @Description get tags for vue tree
// @Param	id		body 	int64	true	"tag id"
// @Success 200 [object] []models.TreeNode
// @Failure 403 string error
// @router /treeNode [get]
func (c *RelController) GetTreeNodes() {
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("id", 0)

	nodes, err := me.GetTreeNodes(tag_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, nodes)
	}
}

// @Title Get vue tag tree
// @Description get tags for vue tree
// @Success 200 [object] []models.TreeNode
// @Failure 403 string error
// @router /tree [get]
func (c *RelController) GetTree() {
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tree, err := me.GetTree()
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, []models.TreeNode{*tree})
	}
}

// @Title GetTagNodes
// @Description get tags for ztree
// @Param	id		body 	int64	true	"tag id"
// @Param	name		body 	string	true	"tag name"
// @Param	lv		body 	int	true	"tag level"
// @Param	otherParam	body 	string	true	"zTreeAsyncTest"
// @Success 200 [object] []models.TagNode
// @Failure 403 error string
// @router /zTreeNodes [post]
func (c *RelController) GetzTreeNodes() {
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("id", 0)
	//name := c.GetString("name")
	//lv, _ := c.GetInt("lv")

	nodes, err := me.GetTagTags(tag_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.Data["json"] = nodes
		c.ServeJSON()
	}
}

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	query	query   string  false	"host name"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {total:int} total number
// @Failure 403 string error
// @router /tag/host/cnt [get]
func (c *RelController) GetTagHostCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	n, err := me.GetTagHostCnt(tag_id, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetHost
// @Description get all Host
// @Param	tag_id		query	int	true	"tag id"
// @Param	query	query	string	false	"host name"
// @Param	per		query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 [object] []models.Host
// @Failure 403 string error
// @router /tag/host/search [get]
func (c *RelController) GetTagHost() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	ret, err := me.GetTagHost(tag_id, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag host relation
// @Description create tag/host relation
// @Param	body	body	models.RelTagHost	true	""
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router /tag/host [post]
func (c *RelController) CreateTagHost() {
	var rel models.RelTagHost

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := me.CreateTagHost(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag host relation
// @Description create tag/hosts relation
// @Param	body	body	models.RelTagHosts	true	""
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router /tag/hosts [post]
func (c *RelController) CreateTagHosts() {
	var rel models.RelTagHosts

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := me.CreateTagHosts(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title delete tag host relation
// @Description delete tag/host relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/host [delete]
func (c *RelController) DelTagHost() {
	var rel models.RelTagHost
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.DeleteTagHost(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag host relation
// @Description delete tag/hosts relation
// @Param	body	body 	models.RelTagHosts	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/hosts [delete]
func (c *RelController) DelTagHosts() {
	var rel models.RelTagHosts
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.DeleteTagHosts(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUserCnt
// @Description get tag role user number
// @Param	query	query   string  false	"user name"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {tatal:int} user total number
// @Failure 403 string error
// @router /tag/role/user/cnt [get]
func (c *RelController) GetTagRoleUserCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	n, err := me.GetTagRoleUserCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUser
// @Description get tag role user
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"user name"
// @Param	per	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 [object] []models.Host
// @Failure 403 string error
// @router /tag/role/user/search [get]
func (c *RelController) GetTagRoleUser() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	ret, err := me.GetTagRoleUser(global, tag_id, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role users relation
// @Description create tag/role/users relation
// @Param	body	body 	models.RelTagRoleUser	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/role/user [post]
func (c *RelController) CreateTagRoleUser() {
	var rel models.RelTagRoleUser
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.CreateTagRoleUser(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag role user relation
// @Description delete tag/role/user relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/role/user [delete]
func (c *RelController) DelTagRoleUser() {
	var rel models.RelTagRoleUser
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.DeleteTagRoleUser(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleTokenCnt
// @Description get tag role token number
// @Param	query	query   string  false	"token name"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {tatal:int} token total number
// @Failure 403 string error
// @router /tag/role/token/cnt [get]
func (c *RelController) GetTagRoleTokenCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	n, err := me.GetTagRoleTokenCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleToken
// @Description get tag role token
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"token name"
// @Param	per	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 [object] []models.Host
// @Failure 403 string error
// @router /tag/role/token/search [get]
func (c *RelController) GetTagRoleToken() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	ret, err := me.GetTagRoleToken(global, tag_id, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role tokens relation
// @Description create tag/role/tokens relation
// @Param	body	body 	models.RelTagRoleToken	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/role/token [post]
func (c *RelController) CreateTagRoleToken() {
	var rel models.RelTagRoleToken
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.CreateTagRoleToken(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag role token relation
// @Description delete tag/role/token relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success 200 {total:int} affected number
// @Failure 403 string error
// @router /tag/role/token [delete]
func (c *RelController) DelTagRoleToken() {
	var rel models.RelTagRoleToken
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.DeleteTagRoleToken(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}
