/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/service/expr"
)

// Operations about event
type EventController struct {
	BaseController
}

func accessEventTriggerOp(op *models.Operator, tagId int64) error {
	if op.IsAdmin() {
		return nil
	}
	if tagId == 0 && op.IsOperator() {
		return nil
	}
	return op.Access(models.SYS_O_TOKEN, tagId)
}

// @Title Create event trigger
// @Description create event trigger
// @Param	body	body 	models.EventTriggerApiAdd	true	"body for trigger content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /trigger [post]
func (c *EventController) CreateEventTrigger() {
	var input models.EventTriggerApiAdd

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := accessEventTriggerOp(op, input.TagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	if input.ParentId > 0 {
		p, err := op.GetEventTrigger(input.ParentId)
		if err != nil {
			c.SendMsg(400, fmt.Sprintf("parent %d not exist",
				p.Id))
			return
		} else if p.ParentId > 0 {
			c.SendMsg(400, fmt.Sprintf("trigger %d can not "+
				"be a parent node", input.ParentId))
			return
		}
	}

	trigger := models.EventTrigger{
		TagId:    input.TagId,
		ParentId: input.ParentId,
		Priority: input.Priority,
		Name:     input.Name,
		Metric:   input.Metric,
		Tags:     input.Tags,
		Expr:     input.Expr,
		Msg:      input.Msg,
	}

	if _, err := expr.Parse(input.Expr); err != nil {
		c.SendMsg(400, fmt.Sprintf("expr parse failed: %s", err.Error()))
		return
	}

	if id, err := op.CreateEventTrigger(&trigger); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title clone event triggers to tag node
// @Description clone event triggers to tag node
// @Param	body	body 	models.EventTriggerApiClone	true	"body for clone event trigger content"
// @Success 200 {object} models.Stats api call result
// @Success 400 {object} models.Stats api call result
// @router /trigger/clone [post]
func (c *EventController) CloneEventTrigger() {
	var (
		success int64
		input   models.EventTriggerApiClone
	)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := accessEventTriggerOp(op, input.TagId); err != nil {
		c.SendMsg(400, statsObj(success, err))
		return
	}

	for _, id := range input.EventTriggerIds {
		if _, err := op.CloneEventTrigger(id, input.TagId); err != nil {
			c.SendMsg(400, statsObj(success, err))
			return
		}
		success++
	}
	c.SendMsg(200, statsObj(success, nil))
}

// @Title Get event triggers Cnt
// @Description get event triggers number
// @Param	tag_id	query	int	false	"tag id"
// @Param	deep	query   int	false	"0: cur tag, 1: include parent , 2: include child"
// @Param	query	query	string	false	"trigger name"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /trigger/cnt [get]
func (c *EventController) GetEventTriggersCnt() {
	tagId, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetInt("deep", models.CTL_SEARCH_CUR)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetEventTriggersCnt(tagId, query, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title Get event triggers
// @Description get all event triggers
// @Param	tag_id	query	int	false	"tag id"
// @Param	deep	query   int	false	"0: cur tag, 1: include parent , 2: include child"
// @Param	query	query	string	false	"trigger name"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.EventTriggerApiGet triggers ui info
// @Failure 400 string error
// @router /trigger/search [get]
func (c *EventController) GetEventTriggers() {
	tagId, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetInt("deep", models.CTL_SEARCH_CUR)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetEventTriggers(tagId, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	}

	for i, _ := range ret {
		if ret[i].Refcnt > 0 {
			ret[i].Children, _ = op.GetEventTriggerChilds(ret[i].Id)
		}
	}
	c.SendMsg(200, ret)
}

// @Title update event trigger
// @Description update the trigger
// @Param	body	body 	models.EventTriggerApiUpdate	true	"body for event trigger content"
// @Success 200 {object} models.Trigger trigger info
// @Failure 400 string error
// @router /trigger [put]
func (c *EventController) UpdateEventTrigger() {
	input := models.EventTriggerApiUpdate{}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	p, err := op.GetEventTrigger(input.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	event := *p
	falcon.Override(&event, &input)

	if _, err := expr.Parse(input.Expr); err != nil {
		c.SendMsg(400, fmt.Sprintf("expr parse failed: %s", err.Error()))
		return
	}

	if ret, err := op.UpdateEventTrigger(&event); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title delete event trigger
// @Description delete the event trigger
// @Param	body	body 	models.EventTriggerApiDel	true	"body for delete event trigger content"
// @Success 200 {object} models.Stats api call result
// @Failure {code:400, msg:string}
// @router /trigger [delete]
func (c *EventController) DeleteEventTrigger() {
	var (
		success int64
		input   models.EventTriggerApiDel
	)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	for _, id := range input.EventTriggerIds {
		if _, err := op.DeleteEventTrigger(id, input.TagId); err != nil {
			c.SendMsg(400, statsObj(success, err))
			return
		}
		success++
	}

	c.SendMsg(200, statsObj(success, nil))
}
