/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"
)

type Strategy struct {
	Id        int64  `json:"id"`
	MetricId  int64  `json:"metricId"`
	Tags      string `json:"tags"`
	MaxStep   int    `json:"maxStep"`
	Priority  int    `json:"priority"`
	Func      string `json:"fun"`
	Op        string `json:"op"`
	Condition string `json:"condition"`
	Note      string `json:"note"`
	Metric    string `json:"metric"`
	RunBegin  string `json:"runBegin"`
	RunEnd    string `json:"runEnd"`
	TplId     int64  `json:"tplId"`
}

func (op *Operator) AddStrategy(o *Strategy) (id int64, err error) {
	id, err = op.SqlInsert("insert strategy (metric_id, tags, max_step, priority, func, op, `condition`, note, metric, run_begin, run_end, tpl_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", o.MetricId, o.Tags, o.MaxStep, o.Priority, o.Func, o.Op, o.Condition, o.Note, o.Metric, o.RunBegin, o.RunEnd, o.TplId)
	if err != nil {
		return
	}
	return
}

func (op *Operator) GetStrategy(id int64) (ret *Strategy, err error) {
	ret = &Strategy{}
	err = op.SqlRow(ret, "select id, metric_id, tags, max_step, priority, func, op, `condition`, note, metric, run_begin, run_end, tpl_id from strategy where id = ?", id)
	return
}

func sqlStrategy(tid int64, query string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}
	if tid != 0 {
		sql2 = append(sql2, "tpl_id = ?")
		sql3 = append(sql3, tid)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetStrategysCnt(tid int64, query string) (cnt int64, err error) {
	sql, sql_args := sqlStrategy(tid, query)
	err = op.O.Raw("SELECT count(*) FROM strategy "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetStrategys(tid int64, query string, limit, offset int) (ret []*Strategy, err error) {
	sql, sql_args := sqlStrategy(tid, query)
	sql = sqlLimit("select id, metric_id, tags, max_step, priority, func, op, `condition`, note, metric, run_begin, run_end, tpl_id from strategy "+sql+" ORDER BY id", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateStrategy(id int64, s *Strategy) (ret *Strategy, err error) {

	_, err = op.SqlExec("update strategy set metric_id = ?, tags = ?, max_step = ?, priority = ?, func = ?, op = ?, `condition` = ?, note = ?, metric = ?, run_begin = ?, run_end = ?, tpl_id = ? where id = ?", s.MetricId, s.Tags, s.MaxStep, s.Priority, s.Func, s.Op, s.Condition, s.Note, s.Metric, s.RunBegin, s.RunEnd, s.TplId, id)
	if err != nil {
		return
	}

	ret, err = op.GetStrategy(id)
	return ret, err
}

func (op *Operator) DeleteStrategy(id int64) (err error) {

	_, err = op.SqlExec("delete from strategy where id = ?", id)
	return
}
