/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

type ExpressionAction struct {
	Expression Expression `json:"expression"`
	Action     Action     `json:"action"`
}

// for ui
type ExpressionUi struct {
	Id              int64  `json:"id"`
	Expression      string `json:"expression"`
	Op              string `json:"op"`
	Condition       string `json:"condition"`
	MaxStep         int    `json:"max_step"`
	Priority        int    `json:"priority"`
	ActionThreshold string `json:"fun"`
	Pause           int    `json:"pause"`
	Uic             string `json:"uic"`
	Creator         string `json:"creator"`
}

type Expression struct {
	Id              int64  `json:"id"`
	Expression      string `json:"expression"`
	Op              string `json:"op"`
	Condition       string `json:"condition"`
	MaxStep         int    `json:"maxStep"`
	Priority        int    `json:"priority"`
	Msg             string `json:"msg"`
	ActionThreshold string `json:"fun"`
	Pause           int    `json:"pause"`
	ActionId        int64  `json:"-"`
	CreateUserId    int64  `json:"-"`
}

func (u *User) AddExpression(r *Expression) (id int64, err error) {
	r.Id = 0
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	moduleCache[CTL_M_EXPRESSION].set(id, r)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_ADD, jsonStr(r))
	return
}

func (u *User) getExpression(id int64) (*Expression, error) {
	if r, ok := moduleCache[CTL_M_EXPRESSION].get(id).(*Expression); ok {
		return r, nil
	}
	r := &Expression{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		moduleCache[CTL_M_EXPRESSION].set(id, r)
	}
	return r, err
}

func (u *User) GetExpressionAction(id int64) (*ExpressionAction, error) {
	var ret ExpressionAction

	if e, err := u.getExpression(id); err != nil {
		return nil, err
	} else {
		ret.Expression = *e
	}

	if a, err := u.GetAction(ret.Expression.ActionId); err != nil {
		return nil, err
	} else {
		ret.Action = *a
	}

	return &ret, nil
}

func queryExpressions(query string, user_id int64) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	if query != "" {
		sql2 = append(sql2, "a.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if user_id != 0 {
		sql2 = append(sql2, "a.create_user_id = ?")
		sql3 = append(sql3, user_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (u *User) GetExpressionsCnt(query string, user_id int64) (cnt int64, err error) {
	sql, sql_args := queryExpressions(query, user_id)
	err = orm.NewOrm().Raw("SELECT count(*) FROM expression a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (u *User) GetExpressions(query string, user_id int64, limit, offset int) (ret []ExpressionUi, err error) {
	sql, sql_args := queryExpressions(query, user_id)
	sql = "SELECT a.id as id, a.name as name, a.expression as expression, a.op as op, a.`condition` as `condition`, a.max_step as max_step, a.priority as priority, a.msg as msg, a.action_threshold as action_threshold, a.pause as pause, b.uic as uic, c.name as creator from expression a LEFT JOIN action b ON a.action_id = b.id LEFT JOIN user c ON a.create_user_id = c.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = orm.NewOrm().Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (u *User) UpdateExpressionAction(id int64, _o *ExpressionAction) (o *Expression, err error) {
	var e *Expression
	e, err = u.UpdateExpression(id, &_o.Expression)
	if err != nil {
		return nil, err
	}

	_, err = u.UpdateAction(e.ActionId, &_o.Action)
	if err != nil {
		return nil, err
	}

	return e, err
}

func (u *User) PauseExpression(id int64, pause int) (o *Expression, err error) {
	if o, err = u.getExpression(id); err != nil {
		return nil, ErrNoExpression
	}
	o.Pause = pause
	_, err = orm.NewOrm().Update(o)
	moduleCache[CTL_M_EXPRESSION].set(id, o)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_SET, "")
	return o, err
}

func (u *User) UpdateExpression(id int64, _o *Expression) (o *Expression, err error) {
	if o, err = u.getExpression(id); err != nil {
		return nil, ErrNoExpression
	}

	o.Expression = _o.Expression
	o.Op = _o.Op
	o.Condition = _o.Condition
	o.MaxStep = _o.MaxStep
	o.Priority = _o.Priority
	o.Msg = _o.Msg
	o.ActionThreshold = _o.ActionThreshold
	o.Pause = _o.Pause

	_, err = orm.NewOrm().Update(o)
	moduleCache[CTL_M_EXPRESSION].set(id, o)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_SET, "")
	return o, err
}

func (u *User) DeleteExpression(id int64) error {
	expression, err := u.getExpression(id)
	if err != nil {
		return err
	}

	if _, err = orm.NewOrm().Delete(&Action{Id: expression.ActionId}); err != nil {
		return err
	}

	if _, err := orm.NewOrm().Delete(&Expression{Id: id}); err != nil {
		return err
	}

	moduleCache[CTL_M_EXPRESSION].del(id)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_DEL, "")

	return nil
}
