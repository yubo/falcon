/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"strings"
	"time"

	"github.com/yubo/falcon/utils"
)

func NewMetricValue(step int, host, metric string,
	val float64, dataType string, tags ...string) *utils.MetaData {

	mv := &utils.MetaData{
		Host:  host,
		Name:  metric,
		Value: val,
		Ts:    time.Now().Unix(),
		Step:  int64(step),
		Type:  dataType,
	}

	if len(tags) > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return mv
}

func GaugeValue(step int, host, metric string, val float64, tags ...string) *utils.MetaData {
	return NewMetricValue(step, host, metric, val, utils.GAUGE, tags...)
}

func CounterValue(step int, host, metric string, val float64, tags ...string) *utils.MetaData {
	return NewMetricValue(step, host, metric, val, utils.COUNTER, tags...)
}
