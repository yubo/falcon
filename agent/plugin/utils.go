/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package plugin

import (
	"strings"

	"github.com/yubo/falcon/specs"
)

func NewMetricValue(step int, host, metric string,
	val float64, dataType string, tags ...string) *specs.MetaData {

	mv := &specs.MetaData{
		Host: host,
		K:    metric,
		V:    val,
		Step: int64(step),
		Type: dataType,
	}

	if len(tags) > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return mv
}

func GaugeValue(step int, host, metric string, val float64, tags ...string) *specs.MetaData {
	return NewMetricValue(step, host, metric, val, specs.GAUGE, tags...)
}

func CounterValue(step int, host, metric string, val float64, tags ...string) *specs.MetaData {
	return NewMetricValue(step, host, metric, val, specs.COUNTER, tags...)
}
