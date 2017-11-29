/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
)

type TriggerApiAdd struct {
	ParentId int64  `json:"parent_id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type TriggerApiGet struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"`
	TplId      int64  `json:"tpl_id"`
	TagId      int64  `json:"tag_id"`
	ParentName string `json:"parent_name"`
	TplName    string `json:"tpl_name"`
	TagName    string `json:"tag_name"`
	Version    int    `json:"version"`
	Refcnt     int    `json:"refcnt"`
	Priority   int    `json:"priority"`
	Name       string `json:"name"`
	Metric     string `json:"metric"`
	Tags       string `json:"tags"`
	Func       string `json:"func"`
	Op         string `json:"op"`
	Value      string `json:"value"`
	Msg        string `json:"msg"`
}

type TagTriggerApiAdd struct {
	ParentId int64  `json:"parent_id"`
	TplId    int64  `json:"tpl_id"`
	TagId    int64  `json:"tag_id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type TagTriggerApiGet struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parent_id"`
	TplId    int64  `json:"tpl_id"`
	TagId    int64  `json:"tag_id"`
	GrpName  string `json:"grp_name"`
	TplName  string `json:"tpl_name"`
	TagName  string `json:"tag_name"`
	Version  int    `json:"version"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type TriggerApiEdit struct {
	Id       int64  `json:"id"`
	TplId    int64  `json:"tpl_id"`
	TagId    int64  `json:"tag_id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type Trigger struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parent_id"`
	TplId    int64  `json:"tpl_id"`
	TagId    int64  `json:"tag_id"`
	Version  int    `json:"version"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

func (op *Operator) CreateTrigger(e *Trigger) (id int64, err error) {
	if e.Name == "" {
		id, err = op.SqlInsert("insert triggers (parent_id, tpl_id, tag_id, priority, metric, tags, func, op, value, msg) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			e.ParentId, e.TplId, e.TagId, e.Priority,
			e.Metric, e.Tags, e.Func, e.Op, e.Value, e.Msg)
	} else {
		id, err = op.SqlInsert("insert triggers (parent_id, tpl_id, tag_id, priority, name, metric, tags, func, op, value, msg) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			e.ParentId, e.TplId, e.TagId, e.Priority, e.Name,
			e.Metric, e.Tags, e.Func, e.Op, e.Value, e.Msg)
	}
	if err != nil {
		return
	}

	if e.ParentId > 0 {
		op.syncTriggerRefcnt(e.ParentId)
	}
	DbLog(op.O, op.User.Id, CTL_M_TRIGGER, id, CTL_A_ADD, jsonStr(e))
	return
}

func (op *Operator) GetTrigger(id int64) (*Trigger, error) {
	e := &Trigger{}
	err := op.SqlRow(e, "select id, parent_id, tpl_id, tag_id, version, priority, name, metric, tags, func, op, value, msg from triggers where id = ?", id)
	return e, err
}

func sqlTrigger(tagId int64, query string, deep, child bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	if query != "" {
		sql2 = append(sql2, "a.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if !child {
		sql2 = append(sql2, "a.parent_id = 0")
	}

	if deep && tagId > 0 {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tagId))
	} else {
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tagId)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTriggersCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlTrigger(0, query, false, false)
	err = op.O.Raw("SELECT count(*) FROM triggers a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) syncTriggerRefcnt(id int64) (int64, error) {
	return op.SqlExec("UPDATE triggers a LEFT JOIN (select a2.id, count(b2.id) as refcnt FROM triggers a2  LEFT JOIN triggers b2 ON a2.id = b2.parent_id WHERE a2.id = ? GROUP BY a2.id) b ON a.id = b.id SET a.refcnt = b.refcnt WHERE a.id = ?", id, id)
}

func (op *Operator) syncTriggerRefcntAll() (int64, error) {
	return op.SqlExec("update triggers a join (select a2.id, count(b2.id) as refcnt from triggers a2  left join triggers b2 on a2.id = b2.parent_id and a2.parent_id = 0 group by a2.id) b on a.id = b.id set a.refcnt = b.refcnt where a.parent_id = 0")
}

func (op *Operator) GetTriggers(query string, limit, offset int) (ret []TriggerApiGet, err error) {
	sql, sql_args := sqlTrigger(0, query, false, false)
	sql = "SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.refcnt as refcnt, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.func as func, a.op as op, a.value as value, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from triggers a LEFT JOIN triggers b ON a.parent_id = b.id LEFT JOIN triggers c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) GetTriggersTpl(id int64) (ret []TriggerApiGet, err error) {
	_, err = op.O.Raw("SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.refcnt as refcnt, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.func as func, a.op as op, a.value as value, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from triggers a LEFT JOIN triggers b ON a.parent_id = b.id LEFT JOIN triggers c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id WHERE a.id = ? or a.parent_id = ?  ORDER BY a.name", id, id).QueryRows(&ret)
	return
}

func (op *Operator) DeleteTrigger(id int64) (n int64, err error) {
	t, err := op.GetTrigger(id)
	if err != nil {
		return 0, err
	}
	if n, err = op.SqlExec("delete from triggers where id = ? or parent_id = ?", id); err != nil {
		return
	}
	if t.ParentId > 0 {
		op.syncTriggerRefcnt(t.ParentId)
	}

	DbLog(op.O, op.User.Id, CTL_M_TRIGGER, id, CTL_A_DEL, fmt.Sprintf("RowsAffected %d", n))
	return
}

// ############ tag trigger ###########

func (op *Operator) GetTagTriggersCnt(tag_id int64, query string, deep bool) (cnt int64, err error) {
	sql, sql_args := sqlTrigger(tag_id, query, deep, false)
	err = op.O.Raw("SELECT count(*) FROM trigger a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTagTriggers(tag_id int64, query string, deep bool, limit, offset int) (ret []TriggerApiGet, err error) {
	sql, sql_args := sqlTrigger(tag_id, query, deep, false)
	sql = "SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.func as func, a.op as op, a.value as value, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from triggers a LEFT JOIN triggers b ON a.parent_id = b.id LEFT JOIN triggers c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}
