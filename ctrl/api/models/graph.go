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

	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/transfer"
)

// for api doc

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
	Counter    string    `json:"counter"` // metric/tags/type
	Ts         int       `json:"-"`
	TCreate    time.Time `json:"-"`
	TModify    time.Time `json:"-"`
}

type DataPointApiGet struct {
	Keys      []string `json:"keys"`
	ConsolFun string   `json:"consol_fun"`
	Start     int64    `json:"start"`
	End       int64    `json:"end"`
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
	sql := fmt.Sprintf("select counter from endpoint_counter WHERE %s group by counter order by counter LIMIT %d",
		strings.Join(sql2, " AND "), limit)
	_, err = Db.Idx.Raw(sql, args...).QueryRows(&ret)
	return
}

func GetDataPoints(in *DataPointApiGet) (*transfer.GetResponse, error) {
	req := &transfer.GetRequest{
		Start: in.Start,
		End:   in.End,
	}

	req.Keys = make([][]byte, len(in.Keys))
	for i, key := range in.Keys {
		req.Keys[i] = []byte(key)
	}

	return ctrl.GetDps(req)
}
