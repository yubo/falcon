/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"strings"

	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about Metrics
type MetricController struct {
	BaseController
}

// @Title GetMetricsCnt
// @Description get Metrics number
// @Param   query     query   string  false       "metric name"
// @Success 200 {object} models.Total  Metric total number
// @Failure 400 string error
// @router /cnt [get]
func (c *MetricController) GetMetricsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetMetricsCnt(query)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetMetrics
// @Description get all Metrics
// @Param   query     query   string  false       "metric name"
// @Param   limit       query   int     false       "limit page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Metric metrics info
// @Failure 400 string error
// @router /search [get]
func (c *MetricController) GetMetrics() {
	query := strings.TrimSpace(c.GetString("query"))
	limit, _ := c.GetInt("limit", models.PAGE_LIMIT)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	metrics, err := op.GetMetrics(query, limit, offset)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, metrics)
	}
}
