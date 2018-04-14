/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/alarm"
	"github.com/yubo/falcon/modules/alarm/expr"
	"github.com/yubo/falcon/modules/ctrl/api/models"
)

// Operations about action (trigger/filter)
type ActionController struct {
	BaseController
}

func accessActionTriggerOp(op *models.Operator, tagId int64) error {
	if op.IsAdmin() {
		return nil
	}
	if tagId == 0 && op.IsOperator() {
		return nil
	}
	return op.Access(models.SYS_O_TOKEN, tagId)
}

func actionTriggerAttrToFlag(email, sms, script bool) (flag uint64) {

	if email {
		flag |= alarm.ACTION_F_EMAIL
	}
	if sms {
		flag |= alarm.ACTION_F_SMS
	}
	if script {
		flag |= alarm.ACTION_F_SCRIPT
	}
	return
}

func actionTriggerFlagToAttr(flag uint64) (email, sms, script bool) {
	return flag&alarm.ACTION_F_EMAIL != 0, flag&alarm.ACTION_F_SMS != 0, flag&alarm.ACTION_F_SCRIPT != 0
}

// @Title Create action trigger
// @Description create action trigger
// @Param	body	body 	models.ActionTriggerApiAdd	true	"body for trigger content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /trigger [post]
func (c *ActionController) CreateActionTrigger() {
	var input models.ActionTriggerApiAdd

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	if err := accessActionTriggerOp(op, input.TagId); err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	trigger := models.ActionTrigger{
		TagId:        input.TagId,
		TokenId:      input.TokenId,
		OrderId:      input.OrderId,
		Expr:         input.Expr,
		ActionFlag:   actionTriggerAttrToFlag(input.Email, input.Sms, input.Script),
		ActionScript: input.ActionScript,
	}

	if id, err := op.CreateActionTrigger(&trigger); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title Get action triggers Cnt
// @Description get action triggers number
// @Param	tag_id	query	int	false	"tag id"
// @Param	deep	query   int	false	"0: cur tag, 1: include parent , 2: include child"
// @Param	query	query	string	false	"trigger name"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /trigger/cnt [get]
func (c *ActionController) GetActionTriggersCnt() {
	tagId, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetInt("deep", models.CTL_SEARCH_CUR)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetActionTriggersCnt(tagId, deep)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title Get action triggers
// @Description get all action triggers
// @Param	tag_id	query	int	false	"tag id"
// @Param	deep	query   int	false	"0: cur tag, 1: include parent , 2: include child"
// @Param	query	query	string	false	"trigger name"
// @Param	limit	query	int	false	"limit page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.ActionTriggerApiGet triggers ui info
// @Failure 400 string error
// @router /trigger/search [get]
func (c *ActionController) GetActionTriggers() {
	tagId, _ := c.GetInt64("tag_id", 0)
	deep, _ := c.GetInt("deep", models.CTL_SEARCH_CUR)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetActionTriggers(tagId, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	}

	for k, _ := range ret {
		e := &ret[k]
		e.Email, e.Sms, e.Script = actionTriggerFlagToAttr(e.ActionFlag)
	}

	c.SendMsg(200, ret)
}

// @Title update action trigger
// @Description update the trigger
// @Param	body	body 	models.ActionTriggerApiUpdate	true	"body for action trigger content"
// @Success 200 {object} models.Trigger trigger info
// @Failure 400 string error
// @router /trigger [put]
func (c *ActionController) UpdateActionTrigger() {
	input := models.ActionTriggerApiUpdate{}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	p, err := op.GetActionTrigger(input.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	action := *p
	core.Override(&action, &input)
	action.ActionFlag = actionTriggerAttrToFlag(input.Email, input.Sms, input.Script)

	if _, err := expr.Parse(input.Expr); err != nil {
		c.SendMsg(400, fmt.Sprintf("expr parse failed: %s", err.Error()))
		return
	}

	if ret, err := op.UpdateActionTrigger(&action); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title delete action trigger
// @Description delete the action trigger
// @Param	body	body 	models.ActionTriggerApiDel	true	"body for delete action trigger content"
// @Success 200 {object} models.Stats api call result
// @Failure 400 {object} models.Stats api call result
// @router /trigger [delete]
func (c *ActionController) DeleteActionTrigger() {
	var (
		success int64
		input   models.ActionTriggerApiDel
	)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)

	for _, id := range input.ActionTriggerIds {
		if _, err := op.DeleteActionTrigger(id, input.TagId); err != nil {
			c.SendMsg(400, statsObj(success, err))
			return
		}
		success++
	}
	c.SendMsg(200, statsObj(success, nil))

}
