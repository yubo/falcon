/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Strategys
type StrategyController struct {
	BaseController
}

// @Title CreateStrategy
// @Description create strategys
// @Param	body	body 	models.Strategy	true	"body for strategy content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *StrategyController) CreateStrategy() {
	var strategy models.Strategy
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &strategy)
	id, err := me.AddStrategy(&strategy)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetStrategysCnt
// @Description get Strategys number
// @Param   query     query   string  false       "strategy name"
// @Success 200  {total:int} strategy total number
// @Failure 403 string error
// @router /cnt [get]
func (c *StrategyController) GetStrategysCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetStrategysCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetStrategys
// @Description get all Strategys
// @Param   query     query   string  false       "strategy name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} models.Strategy
// @Failure 403 error string
// @router /search [get]
func (c *StrategyController) GetStrategys() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	strategys, err := me.GetStrategys(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, strategys)
	}
}

// @Title Get
// @Description get strategy by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Strategy
// @Failure 403 error string
// @router /:id [get]
func (c *StrategyController) GetStrategy() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		strategy, err := me.GetStrategy(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, strategy)
		}
	}
}

// @Title UpdateStrategy
// @Description update the strategy
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Strategy	true		"body for strategy content"
// @Success 200 {object} models.Strategy
// @Failure 403 error string
// @router /:id [put]
func (c *StrategyController) UpdateStrategy() {
	var strategy models.Strategy

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &strategy)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateStrategy(id, &strategy); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteStrategy
// @Description delete the strategy
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *StrategyController) DeleteStrategy() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteStrategy(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
