/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package utils

import (
	"sort"
	"strings"
	"time"

	"github.com/yubo/falcon"
)

func NewMetricValue(metric string,
	val float64, dataType falcon.ItemType, tags ...string) *falcon.Item {

	item := &falcon.Item{
		Value:     val,
		Timestamp: time.Now().Unix(),
		Type:      dataType,
		Metric:    []byte(metric),
	}

	if len(tags) > 0 {
		sort.Strings(tags)
		item.Tags = []byte(strings.Join(tags, ","))
	}

	return item
}

func GaugeValue(metric string, val float64, tags ...string) *falcon.Item {
	return NewMetricValue(metric, val, falcon.ItemType_GAUGE, tags...)
}

func CounterValue(metric string, val float64, tags ...string) *falcon.Item {
	return NewMetricValue(metric, val, falcon.ItemType_COUNTER, tags...)
}
