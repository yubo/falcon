/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Teams
type TeamController struct {
	BaseController
}

// @Title CreateTeam
// @Description create teams
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router / [post]
func (c *TeamController) CreateTeam() {
	var team models.Team
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &team)
	team.Creator = op.User.Id
	id, err := op.AddTeam(&team)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTeamsCnt
// @Description get Teams number
// @Param   query     query   string  false    "team name"
// @Param   mine      query   bool    false    "only show mine team, default false"
// @Success 200 {object} models.Total team total number
// @Failure 400 string error
// @router /cnt [get]
func (c *TeamController) GetTeamsCnt() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", false)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	cnt, err := op.GetTeamsCnt(query, user_id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTeams
// @Description get all Teams
// @Param   query     query   string  false    "team name"
// @Param   mine      query   bool    false    "only show mine team, default false"
// @Param   limit       query   int     false    "limit page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.TeamUi teams info
// @Failure 400 string error
// @router /search [get]
func (c *TeamController) GetTeams() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", false)
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	teams, err := op.GetTeams(query, user_id, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, teams)
	}
}

// @Title Get team by id
// @Description get team by id
// @Param	id	path 	int	true	"team id"
// @Success 200 {object} models.Team team info
// @Failure 400 string error
// @router /:id [get]
func (c *TeamController) GetTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if t, err := op.GetTeam(id); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, t)
		}
	}
}

// @Title Get team member by teamname or teamid
// @Description get team by id
// @Param	id	query 	int	false	"team id"
// @Param	name	query 	string	false	"team name"
// @Success 200 {object} models.TeamMembers user info
// @Failure 400 string error
// @router /member [get]
func (c *TeamController) GetMember() {
	id, _ := c.GetInt64("id", 0)
	name := c.GetString("name")

	if id == 0 && name == "" {
		c.SendMsg(400, "id and name are empty")
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if obj, err := op.GetMember(id, name); err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, obj)
		}
	}
}

// @Title UpdateTeam
// @Description update the team
// @Param	id	path 	string		true	"The id you want to update"
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Team team info
// @Failure 400 string error
// @router /:id [put]
func (c *TeamController) UpdateTeam() {
	var team models.Team

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &team)

	if t, err := op.UpdateTeam(id, &team); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, t)
	}
}

// @Title update Team op.bers
// @Description create teams
// @Param	body	body 	models.TeamMemberIds	true	"body for team content"
// @Success 200 {object} models.TeamMemberIds member info
// @Failure 400 string error
// @router /:id/member [put]
func (c *TeamController) UpdateMember() {
	var member models.TeamMemberIds
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &member)

	if m, err := op.UpdateMember(id, &member); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, m)
	}
}

// @Title DeleteTeam
// @Description delete the team
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /:id [delete]
func (c *TeamController) DeleteTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteTeam(id)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}
