/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/yubo/falcon/ctrl/api/models"
	"time"
)

// Operations about Matters
type MatterController struct {
	BaseController
}

// @Title GetMatters
// @Description get Matters
// @Param   status    query   int     true    "matter status"
// @Param   limit       query   int     false    "limit page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Matter matters
// @Failure 400 string error
// @router /search [get]
func (c *MatterController) GetMatters() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	//status := alarmModels.STATUS_PENDING
	status, err := c.GetInt("status")
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	matters, err := op.QueryMatters(status, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, matters)
	}
}

// @Title GetMattersCnt
// @Description get Matters number
// @Param   status     query   int  false    "matter status"
// @Success 200 {object} int matter total number
// @Failure 400 string error
// @router /cnt [get]
func (c *MatterController) GetMattersCnt() {
	status, _ := c.GetInt("status")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetMatterCnt(status)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title update Matter
// @Description update matter
// @Param	id	path 	int          true	"The id you want to update"
// @Param	body	body 	models.Matter       true	"body for matter content"
// @Success 200 {object} models.Matter matter  info
// @Failure 400 string error
// @router /:id [put]
func (c *MatterController) UpdateMatter() {
	var matter models.Matter
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &matter)
	fmt.Println(matter)

	if err := op.UpdateMatter(id, matter); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "")
	}
}

// @Title GetEvents
// @Description get Events
// @Param   matter	query   int     true    "matter id"
// @Param   limit	query   int     false    "limit page number"
// @Param   offset	query   int     false    "offset  number"
// @Success 200 {object} []models.Events matters
// @Failure 400 string error
// @router /event/search [get]
func (c *MatterController) GetEvents() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	//status := alarmModels.STATUS_PENDING
	matter, err := c.GetInt64("matter")
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	matters := op.QueryEventsByMatter(matter, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, matters)
	}
}

// @Title GetEventCnt
// @Description get Event number
// @Param   matter   query   int  false    "matter id"
// @Success 200 {object} int event total number
// @Failure 400 string error
// @router /event/cnt [get]
func (c *MatterController) GetEventCnt() {
	matter, _ := c.GetInt64("matter")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.QueryEventsCntByMatter(matter)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title Claim matter
// @Description claim matter
// @Param   body    body    models.Claim    true    "body for team content"
// @Success 200 string ok
// @Failure 400 string error
// @router /claim [post]
func (c *MatterController) CreateClaim() {
	var claim models.Claim
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &claim)
	claim.Timestamp = time.Now().Unix()
	claim.User = op.User.Name
	err := op.AddClaim(claim)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "ok")
	}
}
