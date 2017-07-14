/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"

	"github.com/yubo/falcon"
)

type ExpressionActionApiPut struct {
	Expression ExpressionApiPut `json:"expression"`
	Action     Action           `json:"action"`
}

type ExpressionActionApiGet struct {
	Expression Expression `json:"expression"`
	Action     Action     `json:"action"`
}

// for ui
type ExpressionApiGet struct {
	Id              int64  `json:"id"`
	Expression      string `json:"expression"`
	Op              string `json:"op"`
	Condition       string `json:"condition"`
	MaxStep         int    `json:"max_step"`
	Priority        int    `json:"priority"`
	ActionThreshold string `json:"fun"`
	Pause           int    `json:"pause"`
	Uic             string `json:"uic"`
	SendSms         uint   `json:"sendSms"`
	SendMail        uint   `json:"sendMail"`
	Msg             string `json:"msg"`
	Creator         string `json:"creator"`
}

type ExpressionApiPut struct {
	Name            string `json:"name"`
	Expression      string `json:"expression"`
	Op              string `json:"op"`
	Condition       string `json:"condition"`
	MaxStep         int    `json:"maxStep"`
	Priority        int    `json:"priority"`
	Msg             string `json:"msg"`
	ActionThreshold string `json:"fun"`
	Pause           int    `json:"pause"`
	CreateUserId    int64  `json:"createUid"`
	ActionId        int64  `json:"-"`
}

type Expression struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Expression      string `json:"expression"`
	Op              string `json:"op"`
	Condition       string `json:"condition"`
	MaxStep         int    `json:"maxStep"`
	Priority        int    `json:"priority"`
	Msg             string `json:"msg"`
	ActionThreshold string `json:"fun"`
	Pause           int    `json:"pause"`
	CreateUserId    int64  `json:"createUid"`
	ActionId        int64  `json:"-"`
}

func (op *Operator) AddExpression(e *ExpressionApiPut) (id int64, err error) {
	if e.Name != "" {
		id, err = op.SqlInsert("insert expression (name, expression, op, `condition`, max_step, priority, msg, action_threshold, pause, create_user_id, action_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			e.Name, e.Expression, e.Op, e.Condition, e.MaxStep,
			e.Priority, e.Msg, e.ActionThreshold, e.Pause, e.CreateUserId,
			e.ActionId)
	} else {
		id, err = op.SqlInsert("insert expression (expression, op, `condition`, max_step, priority, msg, action_threshold, pause, create_user_id, action_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			e.Expression, e.Op, e.Condition, e.MaxStep, e.Priority,
			e.Msg, e.ActionThreshold, e.Pause, e.CreateUserId, e.ActionId)
	}
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_EXPRESSION, id, CTL_A_ADD, jsonStr(e))
	return
}

func (op *Operator) getExpression(id int64) (*Expression, error) {
	e := &Expression{}
	err := op.SqlRow(e, "select id, name, expression, op, `condition`, max_step, priority, msg, action_threshold, pause, create_user_id, action_id from expression where id = ?", id)
	return e, err
}

func (op *Operator) GetExpressionAction(id int64) (*ExpressionActionApiGet, error) {
	var ret ExpressionActionApiGet

	if e, err := op.getExpression(id); err != nil {
		return nil, err
	} else {
		ret.Expression = *e
	}

	if a, err := op.GetAction(ret.Expression.ActionId); err != nil {
		return nil, err
	} else {
		ret.Action = *a
	}

	return &ret, nil
}

func sqlExpression(query string, user_id int64) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	if query != "" {
		sql2 = append(sql2, "a.expression like ?")
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

func (op *Operator) GetExpressionsCnt(query string, user_id int64) (cnt int64, err error) {
	sql, sql_args := sqlExpression(query, user_id)
	err = op.O.Raw("SELECT count(*) FROM expression a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetExpressions(query string, user_id int64, limit, offset int) (ret []ExpressionApiGet, err error) {
	sql, sql_args := sqlExpression(query, user_id)
	sql = "SELECT a.id as id, a.name as name, a.expression as expression, a.op as op, a.`condition` as `condition`, a.max_step as max_step, a.priority as priority, a.msg as msg, a.action_threshold as action_threshold, a.pause as pause, b.uic as uic, b.send_sms, b.send_mail, c.name as creator from expression a LEFT JOIN action b ON a.action_id = b.id LEFT JOIN user c ON a.create_user_id = c.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateExpressionAction(id int64, ea *ExpressionActionApiPut) (o *Expression, err error) {
	var e *Expression
	e, err = op.UpdateExpression(id, &ea.Expression)
	if err != nil {
		return nil, err
	}

	_, err = op.UpdateAction(e.ActionId, &ea.Action)
	if err != nil {
		return nil, err
	}

	return e, err
}

func (op *Operator) PauseExpression(id int64, pause int) (e *Expression, err error) {
	_, err = op.SqlExec("update expression set pause = ? where id = ?", pause, id)
	if err != nil {
		return
	}

	if e, err = op.getExpression(id); err != nil {
		return nil, falcon.ErrNoExits
	}
	DbLog(op.O, op.User.Id, CTL_M_EXPRESSION, id, CTL_A_SET, jsonStr(pause))
	return e, err
}

func (op *Operator) UpdateExpression(id int64, e *ExpressionApiPut) (o *Expression, err error) {
	_, err = op.SqlExec("update expression set expression = ?, op = ?, `condition` = ?, max_step = ?, priority = ?, msg = ?, action_threshold = ?, pause = ? where id = ?", e.Expression, e.Op, e.Condition, e.MaxStep, e.Priority, e.Msg, e.ActionThreshold, e.Pause, id)

	if o, err = op.getExpression(id); err != nil {
		return nil, falcon.ErrNoExits
	}

	DbLog(op.O, op.User.Id, CTL_M_EXPRESSION, id, CTL_A_SET, jsonStr(e))
	return o, err
}

func (op *Operator) DeleteExpression0(name string) (err error) {
	var e Expression
	err = op.O.Raw("SELECT id, action_id FROM expression where name = ?", name).QueryRow(&e)
	if err != nil {
		return err
	}

	if _, err = op.SqlExec("delete from action where id = ?", e.ActionId); err != nil {
		return err
	}

	if _, err = op.SqlExec("delete from expression where id = ?", e.Id); err != nil {
		return err
	}

	DbLog(op.O, op.User.Id, CTL_M_EXPRESSION, e.Id, CTL_A_DEL, "")
	return nil
}

func (op *Operator) DeleteExpression(id int64) error {
	e, err := op.getExpression(id)
	if err != nil {
		return err
	}

	if _, err = op.SqlExec("delete from action where id = ?", e.ActionId); err != nil {
		return err
	}

	if _, err = op.SqlExec("delete from expression where id = ?", e.Id); err != nil {
		return err
	}

	DbLog(op.O, op.User.Id, CTL_M_EXPRESSION, id, CTL_A_DEL, "")
	return nil
}
