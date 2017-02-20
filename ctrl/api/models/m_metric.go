/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "strings"

type Metric struct {
	Name string `json:"name"`
}

var (
	metrics []*Metric
)

func queryMetrics(query string, limit, offset int) (ret []*Metric) {
	for k, v := range metrics {
		if strings.Contains(v.Name, query) {
			if offset == 0 {
				ret = append(ret, metrics[k])
			} else {
				offset--
			}
			if limit == 0 {
				return
			} else {
				limit--
			}
		}
	}
	return
}

func (u *User) GetMetricsCnt(query string) (int64, error) {
	return int64(len(queryMetrics(query, -1, 0))), nil
}

func (u *User) GetMetrics(query string, limit, offset int) (metrics []*Metric, err error) {
	return queryMetrics(query, limit, offset), nil
}
