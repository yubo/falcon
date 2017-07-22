/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"strings"

	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about porfile/config/info
type PubController struct {
	BaseController
}

// backword  api
type BackwardUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type BackwardUsersWrap struct {
	Msg   string          `json:"msg"`
	Users []*BackwardUser `json:"users"`
}

type BackwardActionWrap struct {
	Msg  string         `json:"msg"`
	Data *models.Action `json:"data"`
}

// @Title Get config
// @Description get ctrl modules config
// @Success 200 {object} [3]map[string]string {defualt{}, conf{}, configfile{}}
// @Failure 400 string error
// @router /config/ctrl [get]
func (c *PubController) GetConfig() {
	var err error
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigerGet("ctrl")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	// config filter
	ret := map[string]interface{}{
		ctrl.C_AUTH_MODULE: conf.Str(ctrl.C_AUTH_MODULE),
		ctrl.C_MASTER_MODE: conf.DefaultBool(ctrl.C_MASTER_MODE, false),
		ctrl.C_DEV_MODE:    conf.DefaultBool(ctrl.C_DEV_MODE, false),
		ctrl.C_MI_MODE:     conf.DefaultBool(ctrl.C_MI_MODE, false),
	}

	c.SendMsg(200, ret)
}

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	tag	query   string	false	"tag string(cop.xiaomi_pdl.inf or cop=xiaomi,pdl=inf)"
// @Param	query	query   string  false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 400 string error
// @router /rel/tag/host/cnt [get]
func (c *PubController) GetTagHostCnt() {
	tag := c.GetString("tag")
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if tag == "/" {
		tag = ""
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
// @router /rel/tag/host/search [get]
func (c *PubController) GetTagHost() {
	tag := c.GetString("tag")
	query := c.GetString("query")
	deep, _ := c.GetBool("deep", true)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if tag == "/" {
		tag = ""
	}
	ret, err := op.GetTagHost(tag, query, deep, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title get action
// @Description get action by action id
// @Param	id	path 	int	true	"the action id for search"
// @Success 200 {object} BackwardActionWrap "action info"
// @Failure 200 {object} BackwardActionWrap error in msg
// @router /api/action/:id [get]
func (c *PubController) GetAction() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(200, &BackwardActionWrap{Msg: err.Error()})
		return
	}

	act, err := models.SysOp.GetAction(id)
	if err != nil {
		c.SendMsg(200, &BackwardActionWrap{Msg: err.Error()})
	} else {
		c.SendMsg(200, &BackwardActionWrap{Data: act})
	}
}

// @Title get team user
// @Description get team user
// @Param	uic	query 	string	true	"the team name for search"
// @Success 200 {object} BackwardUsersWrap "users list"
// @Failure 200 {object} BackwardUsersWrap error in msg
// @router /team/users [get]
func (c *PubController) GetTeamUser() {
	team := c.GetString("uic")
	if mem, err := models.SysOp.GetMember(0, team); err != nil {
		c.SendMsg(200, &BackwardUsersWrap{Msg: err.Error()})
	} else {
		users := make([]*BackwardUser, len(mem.Users))
		for k, v := range mem.Users {
			users[k] = &BackwardUser{
				Name:  v.Name,
				Email: v.Email,
				Phone: v.Phone,
			}
		}
		c.SendMsg(200, &BackwardUsersWrap{Users: users})
	}
}
