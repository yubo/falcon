/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
	"time"
	//cmodel "github.com/yubo/falcon/common/model"
	//fhttp "github.com/yubo/falcon/transfer/http"
)

// for api doc
type JsonFloat float64

type RRDData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type GraphQueryResponse struct {
	Endpoint string     `json:"endpoint"`
	Counter  string     `json:"counter"`
	DsType   string     `json:"dstype"`
	Step     int        `json:"step"`
	Values   []*RRDData `json:"Values"` //大写为了兼容已经再用这个api的用户
}

type Endpoint struct {
	Id       int64     `json:"id"`
	Endpoint string    `json:"endpoint"`
	Ts       int       `json:"-"`
	TCreate  time.Time `json:"-"`
	TModify  time.Time `json:"-"`
}

type EndpointCounter struct {
	Id         uint      `json:"-"`
	EndpointId int64     `json:"-"`
	Counter    string    `json:"counter"`
	Step       int       `json:"step"`
	Type       string    `json:"type"`
	Ts         int       `json:"-"`
	TCreate    time.Time `json:"-"`
	TModify    time.Time `json:"-"`
}

type APIQueryGraphDrawData struct {
	HostNames []string `json:"hostnames"`
	Counters  []string `json:"counters"`
	ConsolFun string   `json:"consol_fun"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	Step      int      `json:"step"`
}

func (op *Operator) GetEndpoint(qs []string, tags []string,
	limit int) (ret []Endpoint, err error) {
	sql2 := []string{}
	args := []interface{}{}

	if qs[0] != "" {
		for _, q := range qs {
			sql2 = append(sql2, "endpoint regexp ?")
			args = append(args, strings.TrimSpace(q))
		}
	}

	if tags[0] != "" {
		sql2 = append(sql2,
			fmt.Sprintf("id in (select distinct endpoint_id from tag_endpoint where tag in ('%s'))", strings.Join(tags, "','")))
	}

	sql := fmt.Sprintf("select endpoint, id from endpoint WHERE %s LIMIT %d", strings.Join(sql2, " AND "), limit)
	_, err = Db.Idx.Raw(sql, args...).QueryRows(&ret)
	return
}

func (op *Operator) GetEndpointCounter(query string, ids []string, limit int) (ret []EndpointCounter, err error) {
	sql2 := []string{}
	args := []interface{}{}

	sql2 = append(sql2, fmt.Sprintf("endpoint_id IN (%s)", strings.Join(ids, ", ")))

	if query != "" {
		for _, q := range strings.Split(query, ",") {
			sql2 = append(sql2, "counter regexp ?")
			args = append(args, strings.TrimSpace(q))
		}
	}
	sql := fmt.Sprintf("select counter, step, type from endpoint_counter WHERE %s group by counter order by counter LIMIT %d",
		strings.Join(sql2, " AND "), limit)
	_, err = Db.Idx.Raw(sql, args...).QueryRows(&ret)
	return
}

/*
func GetCounterData(inputs *APIQueryGraphDrawData) (ret []*cmodel.GraphQueryResponse, err error) {
	items := make([]cmodel.GraphInfoParam, 0)

	for _, host := range inputs.HostNames {
		for _, counter := range inputs.Counters {
			items = append(items, cmodel.GraphInfoParam{Endpoint: host, Counter: counter})
		}
	}
	query := fhttp.GraphHistoryParam{
		Start:            int(inputs.StartTime),
		End:              int(inputs.EndTime),
		CF:               inputs.ConsolFun,
		EndpointCounters: items,
	}
	err = postJson(transferUrl+"/graph/history", query, &ret)
	return
}
*/
