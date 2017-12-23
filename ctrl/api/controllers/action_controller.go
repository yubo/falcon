/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

// Operations about action (trigger/filter)
type ActionController struct {
	BaseController
}

// @Title Create action trigger
// @Description create action trigger
// @Param	body	body 	models.ActionTriggerApiAdd	true	"body for trigger content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /trigger [post]
func (c *ActionController) CreateActionTrigger() {
	/*
		var input models.ActionTriggerApiAdd

		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		json.Unmarshal(c.Ctx.Input.RequestBody, &input)

		if err := accessActionTriggerOp(op, input.TagId); err != nil {
			c.SendMsg(403, err.Error())
			return
		}

		if input.ParentId > 0 {
			p, err := op.GetActionTrigger(input.ParentId)
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

		trigger := models.ActionTrigger{
			TagId:    input.TagId,
			ParentId: input.ParentId,
			Priority: input.Priority,
			Name:     input.Name,
			Metric:   input.Metric,
			Tags:     input.Tags,
			Func:     input.Func,
			Op:       input.Op,
			Value:    input.Value,
			Msg:      input.Msg,
		}

		if id, err := op.CreateActionTrigger(&trigger); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, idObj(id))
		}
	*/
}
