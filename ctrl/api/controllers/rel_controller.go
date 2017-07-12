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
	"github.com/yubo/falcon/utils"
)

// Operations about Relations
type RelController struct {
	BaseController
}

/*
// @Title Get node's children
// @Description get a node's children
// @Param	id	path 	int64	false	"tag id"
// @Success 200 {object} []models.TreeNode All nodes under the current node, read/operate not include
// @ Failure 400 string error
// @ router /treeNode/:id [get]
func (c *RelController) GetTreeNodes() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := c.GetInt64(":id", 0)

	nodes, err := op.GetTreeNodes(tag_id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, nodes)
	}
}
*/

// @Title Get node
// @Description get node's children
// @Param	id	query 	int64	false	"tag id default root(1)"
// @Param	depth	query   int	false	"depth levels default -1(no limit)"
// @Success 200 {object} []models.TreeNode all nodes of the tree(read)
// @Failure 400 string error
// @router /node [get]
func (c *RelController) GetNode() {
	var ret *models.TreeNode

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	id, _ := c.GetInt64("id", 1)
	depth, _ := c.GetInt("depth", -1)
	direct := false

	if depth < 0 {
		depth = 100
	}

	if err := op.Access(models.SYS_IDX_R_TOKEN, id); err == nil {
		direct = true
	}

	ret = op.GetTree0(id, depth, direct)
	if ret == nil {
		c.SendMsg(400, "tree empty or Permission denied")
	} else {
		c.SendMsg(200, []*models.TreeNode{ret})
	}
}

// @Title Get tag tree
// @Description get whole tree
// @Param	id	query 	int64	false	"tag id default root(1)"
// @Param	real	query   bool	false	"ignore admin token"
// @Param	depth	query   int	false	"depth levels default 0(no limit)"
// @Success 200 {object} []models.TreeNode all nodes of the tree(read)
// @Failure 400 string error
// @router /tree [get]
func (c *RelController) GetTree() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, _ := c.GetInt64("id", 1)
	depth, _ := c.GetInt("depth", 0)
	real, _ := c.GetBool("real", false)

	tree := op.GetTree(id, depth, real)
	if tree == nil {
		c.SendMsg(400, "tree empty")
	} else {
		c.SendMsg(200, []models.TreeNode{*tree})
	}
}

// @Title Get tags(operate)
// @Description get has operate token tags
// @Param	deep	query   bool	false	"include child tag(default:false)"
// @Success 200 {object} []int64 all ids of the node that can be operated
// @Failure 400 string error
// @router /operate/tag [get]
func (c *RelController) GetOpTag() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	deep, _ := c.GetBool("deep", false)
	ret, _ := op.GetOpTag(deep)
	c.SendMsg(200, ret)
}

// @Title Get tags(read)
// @Description get has read token tags
// @Param	deep	query   bool	false	"include child tag(default:false)"
// @Success 200 {object} []int64 all ids of the node that can be read
// @Failure 400 string error
// @router /read/tag [get]
func (c *RelController) GetReadTag() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	deep, _ := c.GetBool("deep", false)
	ret, _ := op.GetReadTag(deep)
	c.SendMsg(200, ret)
}

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	tag	query   string	false	"tag string(cop.xiaomi_pdl.inf or cop=xiaomi,pdl=inf)"
// @Param	query	query   string  false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/host/cnt [get]
func (c *RelController) GetTagHostCnt() {
	tag := models.TagToNew(c.GetString("tag"))
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if _, err := op.AccessByStr(models.SYS_IDX_R_TOKEN, tag, false); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagHostCnt(tag, query, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetHost
// @Description get all Host
// @Param	tag	query   string	false	"tag string(cop.xiaomi_pdl.inf or cop=xiaomi,pdl=inf)"
// @Param	query	query	string	false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.RelTagHost tag host info
// @Failure 400 string error
// @router /tag/host/search [get]
func (c *RelController) GetTagHost() {
	tag := models.TagToNew(c.GetString("tag"))
	query := c.GetString("query")
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if _, err := op.AccessByStr(models.SYS_IDX_R_TOKEN, tag, false); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagHost(tag, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag host relation
// @Description create tag/host relation
// @Param	body	body	models.RelTagHost	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/host [post]
func (c *RelController) CreateTagHost() {
	var rel models.RelTagHost

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if id, err := op.CreateTagHost(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag host relation
// @Description create tag/hosts relation
// @Param	body	body	models.RelTagHosts	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/hosts [post]
func (c *RelController) CreateTagHosts() {
	var rel models.RelTagHosts

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	if n, err := op.CreateTagHosts(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag host relation
// @Description delete tag/host relation
// @Param	body	body	models.RelTagHostApiDel	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/host [delete]
func (c *RelController) DelTagHost() {
	var rel models.RelTagHostApiDel

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
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
// @Param	body	body	models.RelTagHosts	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/hosts [delete]
func (c *RelController) DelTagHosts() {
	var rel models.RelTagHosts

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
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

// @Title GetTagTemplateCnt by tag string
// @Description get Tag-Template number
// @Param	query	query   string  false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	tag_string	query   string	true	"tag string"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/template/cnt/0 [get]
func (c *RelController) GetTagTplCnt0() {
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := op.GetTagIdByName(c.GetString("tag_string"))

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagTplCnt(tag_id, query, deep, mine)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTemplate by tag string
// @Description get all Template by tag string
// @Param	tag_string	query   string	true	"tag string"
// @Param	query	query	string	false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	limit	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagTplGet templates info
// @Failure 400 string error
// @router /tag/template/search/0 [get]
func (c *RelController) GetTagTpl0() {
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := op.GetTagIdByName(c.GetString("tag_string"))

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagTpl(tag_id, query, deep, mine, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title GetTagTemplateCnt
// @Description get Tag-Template number
// @Param	query	query   string  false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/template/cnt [get]
func (c *RelController) GetTagTplCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagTplCnt(tag_id, query, deep, mine)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTemplate
// @Description get all Template
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	limit	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagTplGet templates info
// @Failure 400 string error
// @router /tag/template/search [get]
func (c *RelController) GetTagTpl() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagTpl(tag_id, query, deep, mine, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag template relation
// @Description create tag/template relation
// @Param	body	body	models.RelTagTpl	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/template [post]
func (c *RelController) CreateTagTpl() {
	var rel models.RelTagTpl

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if id, err := op.CreateTagTpl(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag template relation by tag string & template name
// @Description create tag/templates relation by tag string & template name
// @Param	body	body	models.RelTagTpl0	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/template/0 [post]
func (c *RelController) CreateTagTpl0() {
	var rel models.RelTagTpl0

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if _, err := op.AccessByStr(models.SYS_IDX_O_TOKEN,
		rel.TagString, false); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if !op.IsAdmin() {
		c.SendMsg(400, utils.EACCES.Error())
		return
	}

	if id, err := op.CreateTagTpl0(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag template relation
// @Description create tag/templates relation
// @Param	body	body	models.RelTagTpls	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/templates [post]
func (c *RelController) CreateTagTpls() {
	var rel models.RelTagTpls

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if id, err := op.CreateTagTpls(&rel); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title delete tag template relation by tag string & template name
// @Description delete tag/template relation by tag string & template name
// @Param	body		body 	models.RelTagTpl0	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/template/0 [delete]
func (c *RelController) DelTagTpl0() {
	var rel models.RelTagTpl0
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if _, err := op.AccessByStr(models.SYS_IDX_O_TOKEN, rel.TagString,
		true); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagTpl0(&rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag template relation
// @Description delete tag/template relation
// @Param	body		body 	models.RelTagTpl	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/template [delete]
func (c *RelController) DelTagTpl() {
	var rel models.RelTagTpl
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagTpl(&rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag template relation
// @Description delete tag/templates relation
// @Param	body	body 	models.RelTagTpls	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/templates [delete]
func (c *RelController) DelTagTpls() {
	var rel models.RelTagTpls
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagTpls(&rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUserCnt
// @Description get tag role user number
// @Param	query	query   string  false	"user name"
// @Param	global	query   bool	false	"ignore tag id(default false)"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Total user total number
// @Failure 400 string error
// @router /tag/role/user/cnt [get]
func (c *RelController) GetTagRoleUserCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagRoleUserCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUser
// @Description get tag role user
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"user name"
// @Param	global	query   bool	false	"ignore tag id(default false)"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagRoleUser tag role user info
// @Failure 400 string error
// @router /tag/role/user/search [get]
func (c *RelController) GetTagRoleUser() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagRoleUser(global, tag_id, query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role users relation
// @Description create tag/role/users relation
// @Param	body	body 	models.RelTagRoleUser	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/role/user [post]
func (c *RelController) CreateTagRoleUser() {
	var rel models.RelTagRoleUser
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if rel.RoleId == 0 || rel.TagId == 0 || rel.UserId == 0 {
		c.SendMsg(400, utils.ErrParam.Error())
		return
	}

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.CreateTagRoleUser(rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag role user relation
// @Description delete tag/role/user relation
// @Param	body		body 	models.RelTagRoleUser	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/user [delete]
func (c *RelController) DelTagRoleUser() {
	var rel models.RelTagRoleUser
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_A_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagRoleUser(rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}

// @Title GetTagRoleTokenCnt
// @Description get tag role token number
// @Param	query	query   string  false	"token name"
// @Param	global	query   bool	false	"ignore tag id(default false)"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/role/token/cnt [get]
func (c *RelController) GetTagRoleTokenCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagRoleTokenCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleToken
// @Description get tag role token
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"token name"
// @Param	global	query   bool	false	"ignore tag id(default false)"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 400 string error
// @router /tag/role/token/search [get]
func (c *RelController) GetTagRoleToken() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagRoleToken(global, tag_id, query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role tokens relation
// @Description create tag/role/tokens relation
// @Param	body	body 	models.RelTagRoleToken	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/token [post]
func (c *RelController) CreateTagRoleToken() {
	var rel models.RelTagRoleToken
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if rel.RoleId == 0 || rel.TagId == 0 || rel.TokenId == 0 {
		c.SendMsg(400, utils.ErrParam.Error())
		return
	}

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.CreateTagRoleToken(rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}

// @Title delete tag role token relation
// @Description delete tag/role/token relation
// @Param	body		body 	models.RelTagRoleToken	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/token [delete]
func (c *RelController) DelTagRoleToken() {
	var rel models.RelTagRoleToken
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_IDX_O_TOKEN, rel.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeleteTagRoleToken(rel)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}

// @Title GetTagPluginCnt
// @Description get Tag-plugin number
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"include child tag(default:true)"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/plugindir/cnt [get]
func (c *RelController) GetTagPluginCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetPluginDirCnt(tag_id, deep)
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
func (c *RelController) GetPluginDir() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetPluginDir(tag_id, deep, limit, offset)
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
func (c *RelController) CreatePluginDir() {
	var input models.PluginDirPost

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := op.Access(models.SYS_IDX_O_TOKEN, input.TagId); err != nil {
		c.SendMsg(403, err.Error())
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
func (c *RelController) DeletePluginDir() {
	id, _ := c.GetInt64("id")
	tagId, _ := c.GetInt64("tag_id")

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_O_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.DeletePluginDir(tagId, id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}
