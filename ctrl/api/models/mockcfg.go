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

type Mockcfg struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Obj     string    `json:"obj"`
	ObjType string    `json:"obj_type"` //group, host, other
	Metric  string    `json:"metric"`
	Tags    string    `json:"tags"`
	Dstype  string    `json:"dstype"`
	Step    int       `json:"step"`
	Mock    float64   `json:"mock"`
	Creator string    `json:"creator"`
	TCreate time.Time `json:"ctime"`
}

type APICreateNoDataInputs struct {
	Name    string  `json:"name"`
	Obj     string  `json:"obj"`
	ObjType string  `json:"obj_type"` //group, host, other
	Metric  string  `json:"metric"`
	Tags    string  `json:"tags"`
	DsType  string  `json:"dstype"`
	Step    int     `json:"step"`
	Mock    float64 `json:"mock"`
}

type NoDataApiPut struct {
	Name    string  `json:"name"`
	Obj     string  `json:"obj"`
	ObjType string  `json:"obj_type"` //group, host, other
	Metric  string  `json:"metric"`
	Tags    string  `json:"tags"`
	DsType  string  `json:"dstype"`
	Step    int     `json:"step"`
	Mock    float64 `json:"mock"`
}

func (op *Operator) AddMockcfg(m *NoDataApiPut) (id int64, err error) {
	id, err = op.SqlInsert("insert mockcfg (name, obj, obj_type, metric, tags, dstype, step, mock, creator) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.Name, m.Obj, m.ObjType, m.Metric, m.Tags, m.DsType, m.Step, m.Mock, op.User.Name)
	DbLog(op.O, op.User.Id, CTL_M_MOCKCFG, id, CTL_A_ADD, jsonStr(m))
	return
}

func (op *Operator) GetMockcfg(id int64) (t *Mockcfg, err error) {
	t = &Mockcfg{}
	err = op.SqlRow(t, "select id, name, obj, obj_type, metric, tags, dstype, step, mock, creator, t_create from mockcfg where id = ?", id)
	return
}

func sqlMockcfg(query, user_name string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	if query != "" {
		sql2 = append(sql2, "name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if user_name != "" {
		sql2 = append(sql2, "creator = ?")
		sql3 = append(sql3, user_name)
	}

	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetMockcfgsCnt(query, user_name string) (cnt int64, err error) {
	sql, sql_args := sqlMockcfg(query, user_name)
	err = op.O.Raw("SELECT count(*) FROM mockcfg "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetMockcfgs(query, user_name string, limit, offset int) (ret []Mockcfg, err error) {
	sql, sql_args := sqlMockcfg(query, user_name)
	sql = "SELECT id, name, obj, obj_type, metric, tags, dstype, step, mock, creator, t_create FROM mockcfg  " + sql + " ORDER BY name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateMockcfg(id int64, m *NoDataApiPut) (o *Mockcfg, err error) {

	_, err = op.SqlExec("update mockcfg set name = ?, obj = ?, obj_type = ?, metric = ?, tags = ?, dstype = ?, step = ?, mock = ? where id = ?", m.Name, m.Obj, m.ObjType, m.Metric, m.Tags, m.DsType, m.Step, m.Mock, id)

	if err != nil {
		return
	}

	o, err = op.GetMockcfg(id)

	DbLog(op.O, op.User.Id, CTL_M_MOCKCFG, id, CTL_A_SET, jsonStr(o))
	return
}

func (op *Operator) DeleteMockcfg(id int64) (err error) {

	if _, err = op.SqlExec("delete from mockcfg where id = ?", id); err != nil {
		return err
	}
	DbLog(op.O, op.User.Id, CTL_M_MOCKCFG, id, CTL_A_DEL, "")

	return nil
}
