/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/models"
)

// Operations about Teams
type TeamController struct {
	BaseController
}

// @Title CreateTeam
// @Description create teams
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {id:int} Id
// @Failure 403 string error
// @router / [post]
func (c *TeamController) CreateTeam() {
	var team models.Team
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	json.Unmarshal(c.Ctx.Input.RequestBody, &team)
	team.Creator = me.Id
	id, err := me.AddTeam(&team)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, models.Id{Id: id})
	}
}

// @Title GetTeamsCnt
// @Description get Teams number
// @Param   query     query   string  false    "team name"
// @Success 200  {total:int} team total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TeamController) GetTeamsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	own, _ := c.GetBool("own", false)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetTeamsCnt(query, own)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTeams
// @Description get all Teams
// @Param   query     query   string  false    "team name"
// @Param   own       query   bool    false    "check if the creator is yourself"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} models.Team
// @Failure 403 error string
// @router /search [get]
func (c *TeamController) GetTeams() {
	query := strings.TrimSpace(c.GetString("query"))
	own, _ := c.GetBool("own", false)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	teams, err := me.GetTeams(query, own, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, teams)
	}
}

// @Title Get team by id
// @Description get team by id
// @Param	id	path 	int	true	"team id"
// @Success 200 {object} models.Team
// @Failure 403 string   error
// @router /:id [get]
func (c *TeamController) GetTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		if t, err := me.GetTeam(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, t)
		}
	}
}

// @Title Get team member
// @Description get team by id
// @Param	id	path 	int	true	"team id"
// @Success 200 users models.User
// @Failure 403 string   error
// @router /:id/member [get]
func (c *TeamController) GetMember() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		me, _ := c.Ctx.Input.GetData("me").(*models.User)
		if users, err := me.GetMember(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, map[string]interface{}{"users": users})
		}
	}
}

// @Title UpdateTeam
// @Description update the team
// @Param	id	path 	string		true	"The id you want to update"
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Team
// @Failure 403 string error
// @router /:id [put]
func (c *TeamController) UpdateTeam() {
	var team models.Team

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &team)

	if t, err := me.UpdateTeam(id, &team); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, t)
	}
}

// @Title update Team members
// @Description create teams
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Member
// @Failure 403 string   error
// @router /:id/member [put]
func (c *TeamController) UpdateMember() {
	var member models.Member
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, &member)

	if m, err := me.UpdateMember(id, &member); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, m)
	}
}

// @Title DeleteTeam
// @Description delete the team
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 string "delete success!"
// @Failure 403 string error
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
	c.SendMsg(200, "delete success!")
}
