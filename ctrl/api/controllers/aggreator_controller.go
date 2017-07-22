/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Aggreators
type AggreatorController struct {
	BaseController
}

// @Title CreateAggreator
// @Description create aggreators by tag_string
// @Param	body	body 	models.APICreateAggregatorInput0	true	"body for aggreator content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /0 [post]
func (c *AggreatorController) CreateAggreator0() {
	var inputs0 models.APICreateAggregatorInput0
	var inputs models.APICreateAggregatorInput

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs0)

	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	inputs.TagId, _ = op.GetTagIdByName(inputs0.TagString)

	if err := op.Access(models.SYS_IDX_O_TOKEN, inputs.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if id, err := op.AddAggreator(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title CreateAggreator
// @Description create aggreators
// @Param	body	body 	models.APICreateAggregatorInput	true	"body for aggreator content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *AggreatorController) CreateAggreator() {
	var inputs models.APICreateAggregatorInput
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if err := op.Access(models.SYS_IDX_O_TOKEN, inputs.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	if id, err := op.AddAggreator(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GeaggreatorgreatorsCnt
// @Description get Aggreators number
// @Param	tag_string	query	string	true	"tag string"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total aggreator total number
// @Failure 400 string error
// @router /cnt/0 [get]
func (c *AggreatorController) GetAggreatorsCnt0() {
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := op.GetTagIdByName(c.GetString("tag_string"))

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	cnt, err := op.GetAggreatorsCnt(tag_id, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title Geaggreatorgreators by tag string
// @Description get all Aggreators
// @Param	tag_string	query	string	true	"tag string"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.GetAggreator  aggreators info
// @Failure 400 string error
// @router /search/0 [get]
func (c *AggreatorController) GetAggreators0() {
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := op.GetTagIdByName(c.GetString("tag_string"))

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	ret, err := op.GetAggreators(tag_id, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title GeaggreatorgreatorsCnt
// @Description get Aggreators number
// @Param	tag_id	query	int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total aggreator total number
// @Failure 400 string error
// @router /cnt [get]
func (c *AggreatorController) GetAggreatorsCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	cnt, err := op.GetAggreatorsCnt(tag_id, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title Geaggreatorgreators
// @Description get all Aggreators
// @Param	tag_id	query	int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.Aggreator  aggreators info
// @Failure 400 string error
// @router /search [get]
func (c *AggreatorController) GetAggreators() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if err := op.Access(models.SYS_IDX_R_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	aggreators, err := op.GetAggreators(tag_id, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, aggreators)
	}
}

// @Title Get
// @Description get aggreator by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Aggreator aggreator info
// @Failure 400 string error
// @router /:id [get]
func (c *AggreatorController) GetAggreator() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	ret, err := op.GetAggreator(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if err := op.Access(models.SYS_IDX_R_TOKEN, ret.TagId); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title UpdateAggreator
// @Description update the aggreator
// @Param	id	path 	string	true		"The id you want to update"
// @Param	body	body 	models.APIUpdateAggregatorInput	true		"body for aggreator content"
// @Success 200 {object} models.Aggreator aggreator info
// @Failure 400 string error
// @router /:id [put]
func (c *AggreatorController) UpdateAggreator() {
	var inputs models.APIUpdateAggregatorInput

	id, _ := c.GetInt64(":id")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetAggreator(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if err := op.Access(models.SYS_IDX_O_TOKEN, ret.TagId); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	if u, err := op.UpdateAggreator(id, &inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteAggreators by tag string
// @Description delete the aggreator by tag string
// @Param	rule	body	models.DeleteAggreator0	true	"The aggreators you want to delete"
// @Success 200 {object} models.Total delete aggreator total number
// @Failure 400 string error
// @router /0 [delete]
func (c *AggreatorController) DeleteAggreator0() {
	var inputs models.DeleteAggreator0

	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := op.GetTagIdByName(inputs.TagString)

	if err := op.Access(models.SYS_IDX_O_TOKEN, tag_id); err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	cnt, err := op.DeleteAggreator0(tag_id, &inputs)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}

}

// @Title DeleteAggreator
// @Description delete the aggreator
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /:id [delete]
func (c *AggreatorController) DeleteAggreator() {
	id, _ := c.GetInt64(":id")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetAggreator(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	err = op.Access(models.SYS_IDX_O_TOKEN, ret.TagId)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	err = op.DeleteAggreator(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "delete success!")
	}

}
