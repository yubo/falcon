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

// Operations about idx
type GraphController struct {
	BaseController
}

// @Title GetEndpoint
// @Description get endpoint
// @Param   query	query   string  false    "endpoint name"
// @Param   tag		query   string  false    "tag name(a=b[,c=d])"
// @Param   limit	query   int     false    "limit number"
// @Success 200 {object} []models.Endpoint endpoints info
// @Failure 400 string error
// @router /endpoint [get]
func (c *GraphController) GetEndpoint() {

	query := strings.TrimSpace(c.GetString("query"))
	tag := strings.TrimSpace(c.GetString("tag"))
	if query == "" && tag == "" {
		c.SendMsg(400, "query and tag is empty")
		return
	}

	limit, _ := c.GetInt("limit", 500)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetEndpoint(strings.Split(query, " "),
		strings.Split(tag, ","), limit)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title GetEndpointCounter
// @Description get endpoint's counter
// @Param   ids		query   string  false    "endpoint id ... eg(1,2,4,5)"
// @Param   query	query   string  false    "counter name"
// @Param   limit	query   int     false    "limit number"
// @Success 200 {object} []models.EndpointCounter endpointsCounter info
// @Failure 400 string error
// @router /endpoint_counter [get]
func (c *GraphController) GetEndpointCounter() {

	ids := strings.TrimSpace(c.GetString("ids"))
	if ids == "" {
		c.SendMsg(400, "ids is empty")
		return
	}

	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", 500)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	obj, err := op.GetEndpointCounter(query, strings.Split(ids, ","), limit)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title Get Counters Data for draw lines
// @Description get tmp graph
// @Param	body	body 	models.APIQueryGraphDrawData	true	"api query graph draw data"
// @Success 200 {object} []models.GraphQueryResponse "graph query response"
// @Failure 400 string error
// @router /counter_data [post]
func (c *GraphController) GetCounterData() {
	var inputs models.APIQueryGraphDrawData

	//op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	/*
		obj, err := models.GetCounterData(&inputs)
		if err != nil {
			c.SendMsg(400, err.Error())
		} else {
			c.SendMsg(200, obj)
		}
	*/
	c.SendMsg(200, "")
}
