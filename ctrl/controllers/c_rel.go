/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/models"
)

var (
	rels = map[string]bool{
		"tag_host":         true,
		"tag_role_user":    true,
		"tag_role_token":   true,
		"tag_rule_trigger": true,
	}
)

// Operations about Rel
type RelController struct {
	BaseController
}

// @Title GetTagNodes
// @Description get tags for ztree
// @Param	id		body 	int64	true	"tag id"
// @Param	name		body 	string	true	"tag name"
// @Param	lv		body 	int	true	"tag level"
// @Param	otherParam	body 	string	true	"zTreeAsyncTest"
// @Success [{id:string,name:string,isParent:bool},...] []models.TagNode
// @Failure {code:int, msg:string}
// @router /zTreeNodes [post]
func (c *RelController) GetzTreeNodes() {
	var host models.Host
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("id", 0)
	//name := c.GetString("name")
	//lv, _ := c.GetInt("lv")

	json.Unmarshal(c.Ctx.Input.RequestBody, &host)
	nodes, err := me.GetTagTags(tag_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.Data["json"] = nodes
		c.ServeJSON()
	}
}

// @Title create tag hosts relation
// @Description create tag/hosts relation
// @Param	body		body 	models.RelTagHosts	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/hosts [post]
func (c *RelController) CreateTagHosts() {
	var rel models.RelTagHosts
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.CreateTagHosts(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, n)
	}
}

// @Title create tag role users relation
// @Description create tag/role/users relation
// @Param	body		body 	models.RelTagRoleUsers	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/role/users [post]
func (c *RelController) CreateTagRoleUsers() {
	var rel models.RelTagRoleUsers
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.CreateTagRoleUsers(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, n)
	}
}

// @Title create tag role tokens relation
// @Description create tag/role/tokens relation
// @Param	body		body 	models.RelTagRoleTokens	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/role/tokens [post]
func (c *RelController) CreateTagRoleTokens() {
	var rel models.RelTagRoleTokens
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := me.CreateTagRoleTokens(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, n)
	}
}

// @Title delete tag host relation
// @Description delete tag/host relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/host [delete]
func (c *RelController) DelTagHost() {
	var rel models.RelTagHost
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	err := me.DeleteTagHost(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, nil)
	}
}

// @Title delete tag role user relation
// @Description delete tag/role/user relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/role/user [delete]
func (c *RelController) DelTagRoleUser() {
	var rel models.RelTagRoleUser
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	err := me.DeleteTagRoleUser(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, nil)
	}
}

// @Title delete tag role token relation
// @Description delete tag/role/token relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success {code:200, data:nil}
// @Failure {code:int, msg:string}
// @router /tag/role/token [delete]
func (c *RelController) DelTagRoleToken() {
	var rel models.RelTagRoleToken
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	err := me.DeleteTagRoleToken(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, nil)
	}
}

// /rel/tag/host get:
func (c *MainController) GetTagHost() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("tag_id", 0)

	if tag_id > 0 {
		total, err := me.QueryTagHostCnt(tag_id, query)
		if err != nil || total == 0 {
			goto out
		}
		offset := c.SetPaginator(per, total).Offset()
		c.Data["Hosts"], _ = me.GetTagHost(tag_id, query, per, offset)
		c.Data["TagId"] = tag_id
	}

out:
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Host")
	c.Data["zTree"] = true
	c.Data["Search"] = Search{"query", "host name"}
	c.TplName = "rel/tag_host.tpl"
}

// /rel/tag/role/user get:
func (c *MainController) GetTagRoleUser() {
	var (
		offset int
	)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("tag_id", 0)

	total, err := me.QueryTagRoleUserCnt(tag_id, query)
	if err != nil || total == 0 {
		goto out
	}
	offset = c.SetPaginator(per, total).Offset()
	c.Data["TagRoleUser"], _ = me.GetTagRoleUser(tag_id, query, per, offset)
out:
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Role User")
	c.Data["zTree"] = true
	c.Data["Search"] = Search{"query", "user name"}
	c.TplName = "rel/tag_role_user.tpl"
}

// /rel/tag/role/token get:
func (c *MainController) GetTagRoleToken() {
	var (
		offset int
	)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	tag_id, _ := c.GetInt64("tag_id", 0)

	total, err := me.QueryTagRoleTokenCnt(tag_id, query)
	if err != nil || total == 0 {
		goto out
	}
	offset = c.SetPaginator(per, total).Offset()
	c.Data["TagRoleToken"], _ = me.GetTagRoleToken(tag_id, query, per, offset)
out:

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Role Token")
	c.Data["zTree"] = true
	c.Data["Search"] = Search{"query", "token name"}
	c.TplName = "rel/tag_role_token.tpl"
}

// /rel/tag/rule/trigger get:
func (c *MainController) GetTagRuleTrigger() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Template Trigger")
	c.Data["zTree"] = true
	c.Data["Search"] = Search{"query", "template name"}
	c.TplName = "rel/tag_rule_trigger.tpl"
}
