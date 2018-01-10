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

type Event struct {
	ItemId    string
	TimeStamp int
	Tag       string
	EndPoint  string
	Metric    string
	Tags      string
	Msg       string
	Stats     string
}

type EventTrigger struct {
	Id       int64  `json:"id"`
	TagId    int64  `json:"tag_id"`
	ParentId int64  `json:"parent_id"`
	TplId    int64  `json:"tpl_id"`
	Version  int    `json:"version"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Expr     string `json:"expr"`
	Msg      string `json:"msg"`
}

type EventTriggerApiAdd struct {
	TagId    int64  `json:"tag_id"`
	ParentId int64  `json:"parent_id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Expr     string `json:"expr"`
	Msg      string `json:"msg"`
}

type EventTriggerApiClone struct {
	TagId           int64   `json:"tag_id"`
	EventTriggerIds []int64 `json:"event_trigger_ids"`
}

type EventTriggerApiDel struct {
	TagId           int64   `json:"tag_id"`
	EventTriggerIds []int64 `json:"event_trigger_ids"`
}

type EventTriggerApiGet struct {
	Id         int64                `json:"id"`
	ParentId   int64                `json:"parent_id"`
	TplId      int64                `json:"tpl_id"`
	TagId      int64                `json:"tag_id"`
	ParentName string               `json:"parent_name"`
	TplName    string               `json:"tpl_name"`
	TagName    string               `json:"tag_name"`
	Version    int                  `json:"version"`
	Refcnt     int                  `json:"refcnt"`
	Priority   int                  `json:"priority"`
	Name       string               `json:"name"`
	Metric     string               `json:"metric"`
	Tags       string               `json:"tags"`
	Expr       string               `json:"expr"`
	Msg        string               `json:"msg"`
	Children   []EventTriggerApiGet `json:"children"`
}

type EventTriggerApiUpdate struct {
	Id       int64  `json:"id"`
	TplId    int64  `json:"tpl_id"`
	TagId    int64  `json:"tag_id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Expr     string `json:"expr"`
	Msg      string `json:"msg"`
}

func (op *Operator) CreateEventTrigger(e *EventTrigger) (id int64, err error) {

	id, err = op.SqlInsert("insert event_trigger (parent_id, tpl_id, tag_id, priority, name, metric, tags, expr, msg) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		e.ParentId, e.TplId, e.TagId, e.Priority, iif(e.Name == "", nil, e.Name),
		e.Metric, e.Tags, e.Expr, e.Msg)
	if err != nil {
		return
	}

	if e.ParentId > 0 {
		op.syncEventTriggerRefcnt(e.ParentId)
	}
	DbLog(op.O, op.User.Id, CTL_M_EVENT_TRIGGER, id, CTL_A_ADD, jsonStr(e))
	return
}

func (op *Operator) GetEventTrigger(id int64) (*EventTrigger, error) {
	e := &EventTrigger{}
	err := op.SqlRow(e, "select id, parent_id, tpl_id, tag_id, version, priority, name, metric, tags, expr, msg from event_trigger where id = ?", id)
	return e, err
}

func (op *Operator) syncEventTriggerRefcnt(id int64) (int64, error) {
	return op.SqlExec("UPDATE event_trigger a LEFT JOIN (select a2.id, count(b2.id) as refcnt FROM event_trigger a2  LEFT JOIN event_trigger b2 ON a2.id = b2.parent_id WHERE a2.id = ? GROUP BY a2.id) b ON a.id = b.id SET a.refcnt = b.refcnt WHERE a.id = ?", id, id)
}

func (op *Operator) syncEventTriggerRefcntAll() (int64, error) {
	return op.SqlExec("update event_trigger a join (select a2.id, count(b2.id) as refcnt from event_trigger a2  left join event_trigger b2 on a2.id = b2.parent_id and a2.parent_id = 0 group by a2.id) b on a.id = b.id set a.refcnt = b.refcnt where a.parent_id = 0")
}

func sqlTrigger(tagId int64, query string, deep int) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	sql2 = append(sql2, "a.parent_id = 0")

	if query != "" {
		sql2 = append(sql2, "a.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if deep == CTL_SEARCH_CHILD {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tagId))
	} else if deep == CTL_SEARCH_PARENT {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select sup_tag_id from tag_rel where tag_id = %d)", tagId))
	} else { // CTL_SEARCH_CUR
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tagId)
	}

	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetEventTriggersCnt(tagId int64, query string, deep int) (cnt int64, err error) {
	sql, sql_args := sqlTrigger(tagId, query, deep)
	err = op.O.Raw("SELECT count(*) FROM event_trigger a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetEventTriggers(tagId int64, query string, deep, limit, offset int) (ret []EventTriggerApiGet, err error) {
	sql, sql_args := sqlTrigger(tagId, query, deep)
	sql = "SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.refcnt as refcnt, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.expr as expr, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from event_trigger a LEFT JOIN event_trigger b ON a.parent_id = b.id LEFT JOIN event_trigger c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) GetEventTriggerChilds(id int64) (ret []EventTriggerApiGet, err error) {
	_, err = op.O.Raw("SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.refcnt as refcnt, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.expr as expr, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from event_trigger a LEFT JOIN event_trigger b ON a.parent_id = b.id LEFT JOIN event_trigger c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id WHERE a.parent_id = ?  ORDER BY a.name", id).QueryRows(&ret)
	return
}

func (op *Operator) GetEventTriggersTpl(id int64) (ret []EventTriggerApiGet, err error) {
	_, err = op.O.Raw("SELECT a.id as id, a.parent_id as parent_id, a.tpl_id as tpl_id, a.tag_id as tag_id, a.version as version, a.refcnt as refcnt, a.priority as priority, a.name as name, a.metric as metric, a.tags as tags, a.expr as expr, a.msg as msg, b.name as grp_name, c.name as tpl_name, d.name as tag_name from event_trigger a LEFT JOIN event_trigger b ON a.parent_id = b.id LEFT JOIN event_trigger c ON a.tpl_id = c.id LEFT JOIN tag d ON a.tag_id = d.id WHERE a.id = ? or a.parent_id = ?  ORDER BY a.name", id, id).QueryRows(&ret)
	return
}

func (op *Operator) UpdateEventTrigger(e *EventTrigger) (ret *EventTrigger, err error) {
	_, err = op.SqlExec("update event_trigger set parent_id = ?, tpl_id = ?, tag_id = ?, version = ?, priority = ?, name = ?, metric = ?, tags = ?, expr = ?, msg = ? where id = ?", e.ParentId, e.TplId, e.TagId, e.Version, e.Priority, iif(e.Name == "", nil, e.Name), e.Metric, e.Tags, e.Expr, e.Msg, e.Id)
	if err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_EVENT_TRIGGER, e.Id, CTL_A_SET, jsonStr(e))
	return ret, err

}

func (op *Operator) DeleteEventTrigger(id int64, tagId int64) (n int64, err error) {
	t, err := op.GetEventTrigger(id)
	if err != nil {
		return 0, err
	}
	if n, err = op.SqlExec("delete from event_trigger where tag_id = ? and (id = ? or parent_id = ?)",
		tagId, id, id); err != nil {
		return
	}
	if t.ParentId > 0 {
		op.syncEventTriggerRefcnt(t.ParentId)
	}

	DbLog(op.O, op.User.Id, CTL_M_EVENT_TRIGGER, id, CTL_A_DEL, fmt.Sprintf("RowsAffected %d", n))
	return
}

func (op *Operator) CloneEventTrigger(srcEventTriggerId int64, dstTagId int64) (id int64, err error) {
	var e *EventTrigger
	var childs []EventTrigger

	e, err = op.GetEventTrigger(srcEventTriggerId)
	if err != nil {
		return
	}

	e.TagId = dstTagId
	if e.TplId == 0 {
		e.TplId = srcEventTriggerId
	}
	e.ParentId = 0
	e.Id = 0
	e.Id, err = op.CreateEventTrigger(e)
	if err != nil {
		return
	}

	id = e.Id

	_, err = op.O.Raw("SELECT id, parent_id, tpl_id, tag_id, version, priority, name, metric, tags, expr, msg from event_trigger WHERE parent_id = ?", srcEventTriggerId).QueryRows(&childs)
	if err != nil || len(childs) == 0 {
		return
	}

	for _, child := range childs {
		child.TagId = dstTagId
		child.TplId = e.TplId
		child.ParentId = e.Id
		_, err = op.CreateEventTrigger(&child)
		if err != nil {
			// impossible
			return
		}
	}
	DbLog(op.O, op.User.Id, CTL_M_EVENT_TRIGGER, id, CTL_A_ADD, fmt.Sprintf("clone from %d", srcEventTriggerId))
	return
}
