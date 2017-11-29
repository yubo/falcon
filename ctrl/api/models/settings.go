/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"
	"time"
)

type LogApi struct {
	LogId  int64
	Module int64
	Id     int64
	User   string
	Action int64
	Data   string
	Time   time.Time
}

type LogApiGet struct {
	LogId  int64     `json:"id"`
	Module string    `json:"module"`
	Id     int64     `json:"tid"`
	User   string    `json:"user"`
	Action string    `json:"action"`
	Data   string    `json:"data"`
	Time   time.Time `json:"time"`
}

func logSql(begin, end string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if begin != "" {
		sql2 = append(sql2, "a.time >= ?")
		sql3 = append(sql3, begin)
	}
	if end != "" {
		sql2 = append(sql2, "a.time <= ?")
		sql3 = append(sql3, end)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetLogsCnt(begin, end string) (cnt int64, err error) {
	sql, sql_args := logSql(begin, end)
	err = op.O.Raw("SELECT count(*) FROM log a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetLogs(begin, end string, limit, offset int) (ret []*LogApi, err error) {
	sql, sql_args := logSql(begin, end)
	sql = "select a.id as log_id, a.module, a.module_id as id, b.name as user, a.action, a.data, a.time from log a left join user b on a.user_id = b.id " + sql + " ORDER BY a.id DESC LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}
