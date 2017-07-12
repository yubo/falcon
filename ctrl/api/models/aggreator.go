/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
)

type GetAggreator struct {
	Id          int64  `json:"id"`
	TagId       int64  `json:"tag_id"`
	TagName     string `json:"tag_name"`
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Tags        string `json:"tags"`
	DsType      string `json:"ds_type"`
	Step        int    `json:"step"`
	Creator     string `json:"creator"`
}

type Aggreator struct {
	Id          int64  `json:"id"`
	TagId       int64  `json:"tag_id"`
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Tags        string `json:"tags"`
	DsType      string `json:"ds_type"`
	Step        int    `json:"step"`
	Creator     string `json:"creator"`
}

type DeleteAggreator0Entry struct {
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
}

type DeleteAggreator0 struct {
	TagString string                  `json:"tag_string"`
	Entries   []DeleteAggreator0Entry `json:"entries"`
}

type APICreateAggregatorInput0 struct {
	TagString   string `json:"tag_string"`
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Tags        string `json:"tags"`
	Step        int    `json:"step"`
	// DsType      string `json:"ds_type" binding:"exists"`
}

type APICreateAggregatorInput struct {
	TagId       int64  `json:"tag_id"`
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Tags        string `json:"tags"`
	Step        int    `json:"step"`
	// DsType      string `json:"ds_type" binding:"exists"`
}

type APIUpdateAggregatorInput struct {
	ID          int64  `json:"id"`
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Tags        string `json:"tags"`
	Step        int    `json:"step"`
	// DsType      string `json:"ds_type" binding:"exists"`
}

func (op *Operator) AddAggreator(in *APICreateAggregatorInput) (id int64, err error) {
	id, err = op.SqlInsert("insert aggreator (tag_id, numerator, denominator, endpoint, metric, tags, ds_type, step, creator) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", in.TagId, in.Numerator, in.Denominator, in.Endpoint, in.Metric, in.Tags, "GAUGE", in.Step, op.User.Name)
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_AGGREATOR, id, CTL_A_ADD, jsonStr(in))
	return
}

func (op *Operator) GetAggreator(id int64) (ret *Aggreator, err error) {
	ret = &Aggreator{}
	err = op.SqlRow(ret, "select id, tag_id, numerator, denominator, endpoint, metric, tags, ds_type, step, creator from aggreator where id = ?", id)
	return
}

func aggreatorSql(tag_id int64, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if deep {
		sql2 = append(sql2, fmt.Sprintf("tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tag_id))
	} else {
		sql2 = append(sql2, "tag_id = ?")
		sql3 = append(sql3, tag_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetAggreatorsCnt(tag_id int64, deep bool) (cnt int64, err error) {
	sql, sql_args := aggreatorSql(tag_id, deep)
	err = op.O.Raw("SELECT count(*) FROM aggreator "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetAggreators(tag_id int64, deep bool, limit, offset int) (ret []*GetAggreator, err error) {
	sql, sql_args := aggreatorSql(tag_id, deep)
	sql = "SELECT a.*, b.name as tag_name from aggreator a left join tag b on a.tag_id = b.id " + sql + " ORDER BY id"
	if limit > 0 {
		sql += " LIMIT ? OFFSET ?"
		sql_args = append(sql_args, limit, offset)
	}
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateAggreator(id int64, a *APIUpdateAggregatorInput) (ret *Aggreator, err error) {

	_, err = op.SqlExec("update aggreator set  numerator = ?, denominator = ?, endpoint = ?, metric = ?, tags = ?, step = ? where id = ?", a.Numerator, a.Denominator, a.Endpoint, a.Metric, a.Tags, a.Step, id)
	if err != nil {
		return
	}

	if ret, err = op.GetAggreator(id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_AGGREATOR, id, CTL_A_SET, jsonStr(a))
	return
}

func (op *Operator) DeleteAggreator0(tag_id int64, inputs *DeleteAggreator0) (total int64, err error) {

	for _, entry := range inputs.Entries {

		res, err := op.O.Raw("DELETE FROM `aggreator` "+
			"WHERE endpoint = ? and metric = ? and tags = ? and tag_id = ?",
			entry.Endpoint, entry.Metric, entry.Tags, tag_id).Exec()
		if err == nil {
			n, _ := res.RowsAffected()
			total += n
		}

	}
	DbLog(op.O, op.User.Id, CTL_M_AGGREATOR, 0, CTL_A_DEL, jsonStr(inputs))

	return
}

func (op *Operator) DeleteAggreator(id int64) (err error) {
	if _, err = op.SqlExec("delete from aggreator where id = ?", id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_AGGREATOR, id, CTL_A_DEL, "")
	return
}
