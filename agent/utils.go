/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"sort"
	"strings"
	"time"
)

func NewMetricValue(metric string,
	val float64, typ string, tags ...string) *Item {

	item := &Item{
		Value:     val,
		Timestamp: time.Now().Unix(),
		Type:      []byte(typ),
		Metric:    []byte(metric),
	}

	if len(tags) > 0 {
		sort.Strings(tags)
		item.Tags = []byte(strings.Join(tags, ","))
	}

	return item
}

func GaugeValue(metric string, val float64, tags ...string) *Item {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val float64, tags ...string) *Item {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}
