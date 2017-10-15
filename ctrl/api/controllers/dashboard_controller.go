/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about dashboard
type DashboardController struct {
	BaseController
}

// @Title CreateTmpGraph
// @Description create tmp graph
// @Param	body	body 	models.APITmpGraph	true	"body for tmpGraph content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /tmpgraph [post]
func (c *DashboardController) CreateTmpGraph() {
	var inputs models.APITmpGraph
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if id, err := op.AddDashboardTmpGraph(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTmpGraph
// @Description get tmp graph
// @Param	id	path 	int	true	"The id of tmp graph"
// @Success 200 {object} models.APITmpGraph tmpgraph info
// @Failure 400 string error
// @router /tmpgraph/:id [get]
func (c *DashboardController) GetTmpGraph() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetDashboardTmpGraph(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title CreateGraph
// @Description create graph, .endpoints and .counters is []string
// @Param	body	body 	models.APIGraph	true	"body for graph content"
// @Success 200 {object} models.Id Id
// @Failure 400 string error
// @router /graph [post]
func (c *DashboardController) CreateGraph() {
	var inputs models.APIGraph
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if id, err := op.AddDashboardGraph(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title update Graphs
// @Description update graph, .endpoints and .counters is []string
// @Param	body	body 	[]models.APIGraph	true	"body for graph content"
// @Success 200 {object}  models.Stats DashboardGraph
// @Failure 400 string error
// @router /graphs [put]
func (c *DashboardController) UpdateGraphs() {
	var (
		gs               []models.APIGraph
		success, failure int64
	)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &gs)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	for _, g := range gs {
		if _, err = op.UpdateDashboardGraph(&g); err != nil {
			failure++
		} else {
			success++
		}
	}

	c.SendMsg(200, statsObj(success, failure))
}

// @Title update Graph
// @Description update graph, .endpoints and .counters is []string
// @Param	body	body 	models.APIGraph	true	"body for graph content"
// @Success 200 {object} models.DashboardGraph DashboardGraph
// @Failure 400 string error
// @router /graph/:id [put]
func (c *DashboardController) UpdateGraph() {
	var (
		inputs models.APIGraph
		err    error
		obj    *models.DashboardGraph
	)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	inputs.GraphId, err = c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if obj, err = op.UpdateDashboardGraph(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title GetGraph
// @Description get graph
// @Param	id	path 	int	true	"The id of graph"
// @Success 200 {object} models.APIGraph graph info
// @Failure 400 string error
// @router /graph/:id [get]
func (c *DashboardController) GetGraph() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetDashboardGraph(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		es := strings.Split(obj.Hosts, models.TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(obj.Counters, models.TMP_GRAPH_FILED_DELIMITER)
		c.SendMsg(200, models.APIGraph{
			Title:      obj.Title,
			Endpoints:  es,
			Counters:   cs,
			ScreenId:   obj.ScreenId,
			GraphType:  obj.GraphType,
			TimeSpan:   obj.Timespan,
			Method:     obj.Method,
			Position:   obj.Position,
			FalconTags: obj.FalconTags,
		})
	}
}

// @Title DeleteGraph
// @Description delete the graph
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /graph/:id [delete]
func (c *DashboardController) DeleteGraph() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteDashboardGraph(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "delete success!")
	}
}

// @Title GetGraph by screen
// @Description get graph by screen
// @Param	screen_id	path 	int	true	"The id of screen"
// @Success 200 {object} []models.APIGraph graph info
// @Failure 400 string error
// @router /graph/screen/:screen_id [get]
func (c *DashboardController) GetGraphByScreen() {
	id, err := c.GetInt64(":screen_id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetDashboardGraphByScreen(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title CreateScreen
// @Description create screen
// @Param	body	body 	models.AddDashboardScreen	true	"body for tmpGraph content"
// @Success 200 {object} models.Id screen id
// @Failure 400 string error
// @router /screen [post]
func (c *DashboardController) CreateScreen() {
	var inputs models.DashboardScreen
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)

	if inputs.Name == "" {
		c.SendMsg(400, "empty name")
		return
	}

	if id, err := op.AddDashboardScreen(&inputs); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title Get screen
// @Description get screen
// @Param	id	path 	int	false	"The id of screen"
// @Param	limit	query   int	false	"limit number"
// @Success 200 {object} models.DashboardScreen screen info
// @Failure 400 string error
// @router /screen/:id [get]
func (c *DashboardController) GetScreen() {
	id, _ := c.GetInt64(":id", 0)
	limit, _ := c.GetInt64("limit", 500)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetDashboardScreen(id, limit)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title Get screen by pid
// @Description get screen by pid
// @Param	pid	path 	int	true	"The pid of screen"
// @Success 200 {object} []models.DashboardScreen screen info
// @Failure 400 string error
// @router /screen/pid/:pid [get]
func (c *DashboardController) GetScreenByPid() {
	pid, _ := c.GetInt64(":pid", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetDashboardScreenByPid(pid)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title DeleteScreen
// @Description delete the screen
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 400 string error
// @router /screen/:id [delete]
func (c *DashboardController) DeleteScreen() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteDashboardScreen(id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, "delete success!")
	}
}

// @Title update Screen
// @Description update Screen
// @Param	id	path	string	true	"The id you want to modify"
// @Param	body	body 	models.DashboardScreen	true	"body for screen content"
// @Success 200 {object} models.DashboardScreen DashboardScreen
// @Failure 400 string error
// @router /screen/:id [put]
func (c *DashboardController) UpdateScreen() {
	var screen models.DashboardScreen
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &screen)

	if obj, err := op.UpdateDashboardScreen(id, &screen); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}
