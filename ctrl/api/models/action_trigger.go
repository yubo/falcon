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

/*
 * 以第一个命中的规则为准
 * action trigger 优先级
 *  - 路径越大，优先级越高
 *  - 相同路径长度下，OrderId 越小，优先级越高, Id 越小，优先级越高
 * 匹配规则
 *  - metric 留空，匹配所有情况， 不为空时，字符串比较
 *  - tags 留空，匹配所有情况， 不为空时，比较tag k v
 */
type ActionTrigger struct {
	Id           int64
	TagId        int64
	TokenId      int64
	OrderId      int
	Expr         string
	ActionFlag   uint64
	ActionScript string
}

type ActionTriggerApiAdd struct {
	TagId        int64  `json:"tag_id"`
	TokenId      int64  `json:"token_id"`
	OrderId      int    `json:"order_id"`
	Expr         string `json:"expr"`
	Email        bool   `json:"email"`
	Sms          bool   `json:"sms"`
	Script       bool   `json:"script"`
	ActionScript string `json:"action_script"`
}

type ActionTriggerApiDel struct {
	TagId            int64   `json:"tag_id"`
	ActionTriggerIds []int64 `json:"action_trigger_ids"`
}

type ActionTriggerApiGet struct {
	Id           int64  `json:"id"`
	TagId        int64  `json:"tag_id"`
	TagName      string `json:"tag_name"`
	TokenId      int64  `json:"token_id"`
	TokenName    string `json:"token_name"`
	OrderId      int    `json:"order_id"`
	Expr         string `json:"expr"`
	Email        bool   `json:"email"`
	Sms          bool   `json:"sms"`
	Script       bool   `json:"script"`
	ActionScript string `json:"action_script"`
	ActionFlag   uint64 `json:"-"`
}

type ActionTriggerApiUpdate struct {
	Id           int64  `json:"id"`
	TagId        int64  `json:"tag_id"`
	TokenId      int64  `json:"token_id"`
	OrderId      int    `json:"order_id"`
	Expr         string `json:"expr"`
	Email        bool   `json:"email"`
	Sms          bool   `json:"sms"`
	Script       bool   `json:"script"`
	ActionScript string `json:"action_script"`
}

func (op *Operator) CreateActionTrigger(e *ActionTrigger) (id int64, err error) {

	id, err = op.SqlInsert("insert action_trigger (tag_id, token_id, order_id, expr, action_flag, action_script) values (?, ?, ?, ?, ?, ?)",
		e.TagId, e.TokenId, e.OrderId, e.Expr, e.ActionFlag, e.ActionScript)
	if err != nil {
		return
	}

	op.log(CTL_M_ACTION_TRIGGER, id, CTL_A_ADD, jsonStr(e))
	return
}

func (op *Operator) GetActionTrigger(id int64) (*ActionTrigger, error) {
	e := &ActionTrigger{}
	err := op.SqlRow(e, "select id, tag_id, token_id, order_id, expr, action_flag, action_script from action_trigger where id = ?", id)
	return e, err
}

func sqlActionTrigger(tagId int64, deep int) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

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

func (op *Operator) GetActionTriggersCnt(tagId int64, deep int) (cnt int64, err error) {
	sql, sql_args := sqlActionTrigger(tagId, deep)
	err = op.O.Raw("SELECT count(*) FROM action_trigger a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetActionTriggers(tagId int64, deep, limit, offset int) (ret []ActionTriggerApiGet, err error) {
	sql, sql_args := sqlActionTrigger(tagId, deep)
	sql = "SELECT a.id, a.tag_id, a.token_id, a.order_id, a.expr, a.action_flag, a.action_script, b.name as tag_name, c.name as token_name from action_trigger a LEFT JOIN tag b ON a.tag_id = b.id LEFT JOIN token c ON a.token_id = c.id " + sql + " ORDER BY a.order_id, a.id LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateActionTrigger(e *ActionTrigger) (ret *ActionTrigger, err error) {
	_, err = op.SqlExec("update action_trigger set  tag_id = ?, token_id = ?, order_id = ?, expr = ?, action_flag = ?, action_script = ? where id = ?", e.TagId, e.TokenId, e.OrderId, e.Expr, e.ActionFlag, e.ActionScript, e.Id)
	if err != nil {
		return
	}

	op.log(CTL_M_ACTION_TRIGGER, e.Id, CTL_A_SET, jsonStr(e))
	return e, err

}

func (op *Operator) DeleteActionTrigger(id int64, tagId int64) (n int64, err error) {
	if n, err = op.SqlExec("delete from action_trigger where tag_id = ? and id = ?",
		tagId, id); err != nil {
		return
	}

	op.log(CTL_M_ACTION_TRIGGER, id, CTL_A_DEL, fmt.Sprintf("RowsAffected %d", n))
	return
}
