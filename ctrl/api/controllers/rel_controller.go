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

// Operations about Relations
type RelController struct {
	BaseController
}

// @Title Get tree's node
// @Description get node and it's children
// @Param	id	query 	int64	false	"tag id default root(1)"
// @Param	depth	query   int	false	"depth levels default -1(no limit)"
// @Success 200 {object} models.TreeNode all nodes of the tree(read)
// @Failure 400 string error
// @router /node [get]
func (c *RelController) GetTreeNode() {
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
// @router /operate/tag [get]
func (c *RelController) GetOpTag() {
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
// @router /read/tag [get]
func (c *RelController) GetReadTag() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	expand, _ := c.GetBool("expand", false)
	ret, _ := op.GetReadTag(expand)
	c.SendMsg(200, ret)
}

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	tag_id	query   int	false	"tag id"
// @Param	query	query   string  false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/host/cnt [get]
func (c *RelController) GetTagHostCnt() {
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
// @Success 200 {object} []models.RelTagHost tag host info
// @Failure 400 string error
// @router /tag/host/search [get]
func (c *RelController) GetTagHost() {
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
// @Param	body	body	models.RelTagHostApiAdd	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/host [post]
func (c *RelController) CreateTagHost() {
	var rel models.RelTagHostApiAdd

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
// @Param	body	body	models.RelTagHostsAdd	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/hosts [post]
func (c *RelController) CreateTagHosts() {
	var rel models.RelTagHostsApiAdd

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
// @Param	body	body	models.RelTagHostApiDel	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/host [delete]
func (c *RelController) DelTagHost() {
	var rel models.RelTagHostApiDel

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
// @Param	body	body	models.RelTagHostsApiDel	true	"unbind tag_id host_id relation"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/hosts [delete]
func (c *RelController) DelTagHosts() {
	var rel models.RelTagHostsApiDel

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

// @Title GetTagTriggerCnt by tag string
// @Description get Tag-Template number
// @Param	tag_id	query   int	false	"tag id"
// @Param	query	query   string  false	"trigger name"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /tag/trigger/cnt [get]
func (c *RelController) GetTagTriggerCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if tagId > 0 {
		if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}
	}

	n, err := op.GetTagTriggersCnt(tagId, query, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTriggers by tag id
// @Description get all Template by tag string
// @Param	tag_id	query   int	false	"tag id"
// @Param	query	query	string	false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TriggerApiGet templates info
// @Failure 400 string error
// @router /tag/template/search/0 [get]
func (c *RelController) GetTagTpl0() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if tagId > 0 {
		if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}
	}

	ret, err := op.GetTagTriggers(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag trigger relation
// @Description create tag/template relation
// @Param	body	body	models.RelTagTpl	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/template [post]
func (c *RelController) CreateTagTrigger() {
}

// @Title create tag template relation by tag string & template name
// @Description create tag/templates relation by tag string & template name
// @Param	body	body	models.RelTagTpl0	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/template/0 [post]
func (c *RelController) CreateTagTpl0() {

}

// @Title create tag template relation
// @Description create tag/templates relation
// @Param	body	body	models.RelTagTpls	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/templates [post]
func (c *RelController) CreateTagTpls() {
}

// @Title delete tag template relation by tag string & template name
// @Description delete tag/template relation by tag string & template name
// @Param	body		body 	models.RelTagTpl0	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/template/0 [delete]
func (c *RelController) DelTagTpl0() {
}

// @Title delete tag template relation
// @Description delete tag/template relation
// @Param	body		body 	models.RelTagTpl	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/template [delete]
func (c *RelController) DelTagTpl() {
}

// @Title delete tag template relation
// @Description delete tag/templates relation
// @Param	body	body 	models.RelTagTpls	true	""
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/templates [delete]
func (c *RelController) DelTagTpls() {
}

// @Title GetTagRoleUserCnt
// @Description get tag role user number
// @Param	query	query   string  false	"user name"
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total user total number
// @Failure 400 string error
// @router /tag/role/user/cnt [get]
func (c *RelController) GetTagRoleUserCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagRoleUserCnt(tagId, query, deep)
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
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.TagRoleUserApiGet tag role user info
// @Failure 400 string error
// @router /tag/role/user/search [get]
func (c *RelController) GetTagRoleUser() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagRoleUser(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role users relation
// @Description create tag/role/users relation
// @Param	body	body 	models.TagRoleUserApi	true	""
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tag/role/user [post]
func (c *RelController) CreateTagRoleUser() {
	var rel models.TagRoleUserApi
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if rel.RoleId == 0 || rel.TagId == 0 || rel.UserId == 0 {
		c.SendMsg(400, falcon.ErrParam.Error())
		return
	}

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
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
// @Param	body		body 	models.TagRoleUserApi	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/user [delete]
func (c *RelController) DelTagRoleUser() {
	var rel models.TagRoleUserApi
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_A_TOKEN, rel.TagId); err != nil {
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
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total affected number
// @Failure 400 string error
// @router /tag/role/token/cnt [get]
func (c *RelController) GetTagRoleTokenCnt() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	n, err := op.GetTagRoleTokenCnt(tagId, query, deep)
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
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 400 string error
// @router /tag/role/token/search [get]
func (c *RelController) GetTagRoleToken() {
	tagId, _ := c.GetInt64("tag_id", 1)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetTagRoleToken(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role tokens relation
// @Description create tag/role/tokens relation
// @Param	body	body 	models.TagRoleTokenApi	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/token [post]
func (c *RelController) CreateTagRoleToken() {
	var rel models.TagRoleTokenApi
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if rel.RoleId == 0 || rel.TagId == 0 || rel.TokenId == 0 {
		c.SendMsg(400, falcon.ErrParam.Error())
		return
	}

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
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
// @Param	body		body 	models.TagRoleTokenApi	true	""
// @Success 200 {object} models.Id affected id
// @Failure 400 string error
// @router /tag/role/token [delete]
func (c *RelController) DelTagRoleToken() {
	var rel models.TagRoleTokenApi
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if err := op.Access(models.SYS_O_TOKEN, rel.TagId); err != nil {
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
	tagId, _ := c.GetInt64("tag_id", 1)
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
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
func (c *RelController) GetPluginDir() {
	tagId, _ := c.GetInt64("tag_id", 1)
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_R_TOKEN, tagId); err != nil {
		c.SendMsg(403, err.Error())
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
func (c *RelController) CreatePluginDir() {
	var input models.PluginDirPost

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := op.Access(models.SYS_O_TOKEN, input.TagId); err != nil {
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
	tagId, _ := c.GetInt64("tag_id", 1)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_O_TOKEN, tagId); err != nil {
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
