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

// Operations about porfile/config/info
type SetController struct {
	BaseController
}

type Mss map[string]string

// @Title Get config
// @Description get modules config
// @Param	module	path	string	true	"module name"
// @Success 200 map[string]string {defualt{}, conf{}, configfile{}}
// @Failure 400 string error
// @router /config/:module [get]
func (c *SetController) GetConfig() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigGet(module)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title Get
// @Description get profile
// @Success 200 {object} models.User user info
// @Failure 400 string error
// @router /profile [get]
func (c *SetController) GetUser() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if user, err := op.GetUser(op.User.Id); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, user)
	}
}

// @Title Profile update
// @Description update profile
// @Param	body		body 	models.UserProfileUpdate	true		"body for user content"
// @Success 200 {object} models.User user info
// @Failure 400 string error
// @router /profile [put]
func (c *SetController) UpdateUser() {
	user := &models.UserProfileUpdate{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, user)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	o := *op.User
	o.Cname = user.Cname
	o.Email = user.Email
	o.Phone = user.Phone
	o.Qq = user.Qq
	o.Extra = user.Extra

	if ret, err := op.UpdateUser(&o); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op.User = ret
		c.SendMsg(200, user)
	}
}

// @Title GetLogsCnt
// @Description get logs number
// @Param   begin  query   string  false       "end time(YYYY-MM-DD HH:mm:ss)"
// @Param   end    query   string  false       "begin time(YYYY-MM-DD HH:mm:ss)"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /log/cnt [get]
func (c *SetController) GetLogsCnt() {
	begin := c.GetString("begin")
	end := c.GetString("end")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetLogsCnt(begin, end)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetLogs
// @Description get all Logs
// @Param   begin      query  string  false  "end time(YYYY-MM-DD HH:mm:ss)"
// @Param   end        query  string  false  "begin time(YYYY-MM-DD HH:mm:ss)"
// @Param   limit       query   int     false  "limit page number"
// @Param   offset    query   int     false  "offset  number"
// @Success 200 {object} []models.LogApiGet logs info
// @Failure 400 string error
// @router /log/search [get]
func (c *SetController) GetLogs() {
	begin := c.GetString("begin")
	end := c.GetString("end")
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	logs, err := op.GetLogs(begin, end, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	}

	ret := make([]*models.LogApiGet, len(logs))
	for i := 0; i < len(logs); i++ {
		l := logs[i]
		ret[i] = &models.LogApiGet{
			LogId:  l.LogId,
			Module: models.ModuleName[l.Module],
			Id:     l.Id,
			User:   l.User,
			Action: models.ActionName[l.Action],
			Data:   l.Data,
			Time:   l.Time,
		}
	}

	c.SendMsg(200, ret)
}
