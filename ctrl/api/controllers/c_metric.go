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
// @Success 200 {total:int}  Metric total number
// @Failure 403 string error
// @router /cnt [get]
func (c *MetricController) GetMetricsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	cnt, err := me.GetMetricsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetMetrics
// @Description get all Metrics
// @Param   query     query   string  false       "metric name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 [object] []models.Metric
// @Failure 403 string error
// @router /search [get]
func (c *MetricController) GetMetrics() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	metrics, err := me.GetMetrics(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, metrics)
	}
}
