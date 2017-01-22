/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Teams
type TeamController struct {
	BaseController
}

// @Title CreateTeam
// @Description create teams
// @Param	body	body 	models.Team	true	"body for team content"
// @Success {code:200, data:int} models.Team.Id
// @Failure {code:int, msg:string}
// @router / [post]
func (c *TeamController) CreateTeamUsers() {
	var tu models.TeamUsers
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &tu)

	id, err := me.AddTeamUsers(&tu)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, id)
	}
}

// @Title GetTeamsCnt
// @Description get Teams number
// @Param   query     query   string  false    "team name"
// @Success {code:200, data:int} team number
// @Failure {code:int, msg:string}
// @router /cnt [get]
func (c *TeamController) GetTeamsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTeamsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title GetTeams
// @Description get all Teams
// @Param   query     query   string  false    "team name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success {code:200, data:object} models.Team
// @Failure {code:int, msg:string}
// @router /search [get]
func (c *TeamController) GetTeams() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	teams, err := me.GetTeams(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, teams)
	}
}

// @Title Get
// @Description get team by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success {code:200, data:object} models.Team
// @Failure {code:int, msg:string}
// @router /:id [get]
func (c *TeamController) GetTeam() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		mt, err := me.GetTeamUsers(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, mt)
		}
	}
}

// @Title UpdateTeam
// @Description update the team
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Team	true		"body for team content"
// @Success {code:200, data:object} models.Team
// @Failure {code:int, msg:string}
// @router /:id [put]
func (c *TeamController) UpdateTeam() {
	var team models.TeamUsers

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &team)

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if u, err := me.UpdateTeamUsers(id, &team); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteTeam
// @Description delete the team
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TeamController) DeleteTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	err = me.DeleteTeam(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}

// #####################################
// #############  render ###############
// #####################################
func (c *MainController) GetTeam() {
	var teams []*models.Team

	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	qs := me.QueryTeams(query)
	total, err := qs.Count()
	if err != nil {
		goto out
	}

	_, err = qs.Limit(per,
		c.SetPaginator(per, total).Offset()).All(&teams)
	if err != nil {
		goto out
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Team")
	c.Data["Teams"] = teams
	c.Data["Query"] = query
	c.Data["Search"] = Search{"query", "team name"}

	c.TplName = "team/list.tpl"
	return

out:
	c.SendMsg(400, err.Error())
}

//	beego.Router("/teamusers/edit/:id([0-9]+)", mc, "get:EditTeamUsers")
func (c *MainController) EditTeamUsers() {
	var (
		ids string
		i   int
		id  int64
		err error
		tu  *models.TeamUsers
		me  *models.User
	)

	id, err = c.GetInt64(":id")
	if err != nil {
		goto out
	}

	me, _ = c.Ctx.Input.GetData("me").(*models.User)
	tu, err = me.GetTeamUsers(id)
	if err != nil {
		goto out
	}

	for i = 0; i < len(tu.Users); i++ {
		ids = fmt.Sprintf("%s%d,", ids, tu.Users[i].Id)
	}

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Team")
	c.Data["TeamUsers"] = tu
	c.Data["User_ids"] = ids[:len(ids)-1]
	c.Data["H1"] = "edit team"
	c.Data["Method"] = "put"
	c.TplName = "team/edit.tpl"
	return
out:
	c.SendMsg(400, err.Error())
}

// beego.Router("/teamusers/add", mc, "get:AddTeamUsers")
func (c *MainController) AddTeamUsers() {

	c.PrepareEnv(headLinks[HEAD_LINK_IDX_META].SubLinks, "Team")
	c.Data["Method"] = "post"
	c.Data["H1"] = "add team"
	c.TplName = "team/edit.tpl"
}
