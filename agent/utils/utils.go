/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package utils

import (
	"strings"
	"time"

	"github.com/yubo/falcon"
)

func NewMetricValue(step int, endpoint, metric string,
	val float64, dataType falcon.ItemType, tags ...string) *falcon.Item {

	item := &falcon.Item{
		Value:    val,
		Ts:       time.Now().Unix(),
		Type:     dataType,
		Endpoint: []byte(endpoint),
		Metric:   []byte(metric),
	}

	if len(tags) > 0 {
		item.Tags = []byte(strings.Join(tags, ","))
	}

	return item
}

func GaugeValue(step int, host, metric string, val float64, tags ...string) *falcon.Item {
	return NewMetricValue(step, host, metric, val, falcon.ItemType_GAUGE, tags...)
}

func CounterValue(step int, host, metric string, val float64, tags ...string) *falcon.Item {
	return NewMetricValue(step, host, metric, val, falcon.ItemType_COUNTER, tags...)
}
