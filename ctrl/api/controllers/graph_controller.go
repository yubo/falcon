/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/yubo/falcon"
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
// @Param	body	body 	models.DataPointApiGet	true	"api query graph draw data"
// @Success 200 {object} service.GetResponse "graph query response"
// @Failure 400 string error
// @router /counter_data [post]
func (c *GraphController) GetCounterData() {
	var inputs models.CounterDataApiGet

	//op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	if len(inputs.Counters) == 0 || len(inputs.Endpoints) == 0 {
		c.SendMsg(200, "")
		return
	}

	req := &models.DataPointApiGet{
		ConsolFun: inputs.ConsolFun,
		Start:     inputs.Start,
		End:       inputs.End,
	}

	for _, counter := range inputs.Counters {
		for _, endpoint := range inputs.Endpoints {
			req.Keys = append(req.Keys, endpoint+"/"+counter)
		}
	}

	resp, err := models.GetDataPoints(req)
	if err != nil {
		c.SendMsg(400, err.Error())
		return
	}

	ret := []*models.CounterDataApiGetResponse{}

	for _, data := range resp.Data {
		endpoint, counter, tags, typ, err := falcon.KeyAttr(data.Key)
		if err != nil {
			continue
		}

		if tags != "" {
			counter += "/" + tags
		}

		dp := &models.CounterDataApiGetResponse{
			Endpoint: endpoint,
			Counter:  counter,
			DsType:   typ,
		}
		dp.Values = make([][2]interface{}, len(data.Values))
		for k, v := range data.Values {
			dp.Values[k] = [2]interface{}{v.Timestamp, v.Value}
		}
		ret = append(ret, dp)
	}

	c.SendMsg(200, ret)
}

// @Title Get Counters Data for draw lines
// @Description get tmp graph
// @Param	body	body 	models.DataPointApiGet	true	"api query graph draw data"
// @Success 200 {object} service.GetResponse "graph query response"
// @Failure 400 string error
// @router /datapoint [post]
func (c *GraphController) GetDataPoint() {
	var inputs models.DataPointApiGet

	//op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputs)
	if len(inputs.Keys) == 0 {
		c.SendMsg(200, "")
		return
	}

	obj, err := models.GetDataPoints(&inputs)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}
