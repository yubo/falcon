/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/yubo/falcon/utils"
)

type TemplateAction struct {
	Template Template `json:"template"`
	Action   Action   `json:"action"`
	Pname    string   `json:"pname"`
}

// for ui
type TemplateUi struct {
	Id      int64  `json:"id"`
	Pid     int64  `json:"pid"`
	Name    string `json:"name"`
	Pname   string `json:"pname"`
	Creator string `json:"creator"`
}

// for db
type Template struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	ParentId     int64     `json:"pid"`
	ActionId     int64     `json:"-"`
	CreateUserId int64     `json:"-"`
	CreateTime   time.Time `json:"ctime"`
}

// for db
type Action struct {
	Id                 int64  `json:"id"`
	Uic                string `json:"uic"`
	Url                string `json:"url"`
	SendSms            uint   `json:"send_sms"`
	SendMail           uint   `json:"send_mail"`
	Callback           uint   `json:"callback"`
	BeforeCallbackSms  uint   `json:"before_callback_sms"`
	BeforeCallbackMail uint   `json:"before_callback_mail"`
	AfterCallbackSms   uint   `json:"after_callback_sms"`
	AfterCallbackMail  uint   `json:"after_callback_mail"`
}

func (op *Operator) AddTemplate(o *Template) (id int64, err error) {
	id, err = op.SqlInsert("insert template (name, parent_id, action_id, create_user_id) values (?, ?, ?, ?)",
		o.Name, o.ParentId, o.ActionId, op.User.Id)
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_TEMPLATE, id, CTL_A_ADD, jsonStr(o))
	return
}

func (op *Operator) AddAction(o *Action) (id int64, err error) {
	id, err = op.SqlInsert("insert action (uic, url, send_sms, send_mail, callback, before_callback_sms, before_callback_mail, after_callback_sms, after_callback_mail) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", o.Uic, o.Url, o.SendSms, o.SendMail, o.Callback, o.BeforeCallbackSms, o.BeforeCallbackMail, o.AfterCallbackSms, o.AfterCallbackMail)
	return
}

func (op *Operator) CloneTemplate(tpl_id int64) (*TemplateAction, error) {
	var (
		ret     TemplateAction
		src_tpl *Template
		src_act *Action
		err     error
		objs    []*Strategy
		id      int64
	)

	if src_tpl, err = op.getTemplate(tpl_id); err != nil {
		return nil, err
	}
	if src_act, err = op.GetAction(src_tpl.ActionId); err != nil {
		return nil, err
	}

	ret.Template = *src_tpl
	ret.Action = *src_act
	ret.Template.Name += "_copy" + fmt.Sprintf("%08X", rand.New(rand.NewSource(time.Now().Unix())).Uint32())

	ret.Action.Id, err = op.AddAction(&ret.Action)
	if err != nil {
		return nil, err
	}

	ret.Template.ActionId = ret.Action.Id
	ret.Template.Id, err = op.AddTemplate(&ret.Template)
	if err != nil {
		return nil, err
	}

	objs, _ = op.GetStrategys(src_tpl.Id, "", 0, 0)
	for _, obj := range objs {
		obj.TplId = id
		op.AddStrategy(obj)
	}

	if t, err := op.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (op *Operator) GetTemplate(id int64) (*TemplateAction, error) {
	var ret TemplateAction

	if t, err := op.getTemplate(id); err != nil {
		return nil, err
	} else {
		ret.Template = *t
	}

	if a, err := op.GetAction(ret.Template.ActionId); err != nil {
		return nil, err
	} else {
		ret.Action = *a
	}

	if t, err := op.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (op *Operator) getTemplate(id int64) (t *Template, err error) {
	t = &Template{}
	err = op.SqlRow(t, "select id, name, parent_id, action_id, create_user_id, create_time from template where id = ?", id)
	return
}

func (op *Operator) getTemplateIdByName(name string) (id int64, err error) {
	err = op.SqlRow(&id, "select id from template where name = ?", name)
	return
}

func (op *Operator) GetAction(id int64) (ret *Action, err error) {
	ret = &Action{}
	err = op.SqlRow(ret, "select id, uic, url, send_sms, send_mail, callback, before_callback_sms, before_callback_mail, after_callback_sms, after_callback_mail from action where id = ?", id)
	return
}

func sqlTemplate(query string, user_id int64) (where string, args []interface{}) {
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

func (op *Operator) GetTemplatesCnt(query string, user_id int64) (cnt int64, err error) {
	sql, sql_args := sqlTemplate(query, user_id)
	err = op.O.Raw("SELECT count(*) FROM template a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTemplates(query string, user_id int64, limit, offset int) (ret []TemplateUi, err error) {
	sql, sql_args := sqlTemplate(query, user_id)
	sql = "SELECT a.id as id, b.id as pid, a.name as name, b.name as pname, c.name as creator FROM template a LEFT JOIN template b ON a.parent_id = b.id LEFT JOIN user c ON a.create_user_id = c.id " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)

	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateTemplate(id int64, _o *TemplateAction) (o *Template, err error) {
	var t *Template

	t, err = op.getTemplate(id)
	if err != nil {
		return nil, err
	}

	if !op.IsOperator() && op.User.Id != t.CreateUserId {
		return nil, utils.EACCES
	}

	t, err = op.updateTemplate(id, &_o.Template)
	if err != nil {
		return nil, err
	}
	_, err = op.UpdateAction(t.ActionId, &_o.Action)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (op *Operator) updateTemplate(id int64, _o *Template) (o *Template, err error) {
	_, err = op.SqlExec("update template set name = ?, parent_id = ? where id = ?", _o.Name, _o.ParentId, id)

	if o, err = op.getTemplate(id); err != nil {
		return nil, utils.ErrNoExits
	}
	return o, err
}

func (op *Operator) UpdateAction(id int64, act *Action) (o *Action, err error) {
	_, err = op.SqlExec("update action set uic = ?, url = ?, send_sms = ?, send_mail = ?, callback = ?, before_callback_sms = ?, before_callback_mail = ?, after_callback_sms = ?, after_callback_mail = ? where id = ?", act.Uic, act.Url, act.SendSms, act.SendMail, act.Callback, act.BeforeCallbackSms, act.BeforeCallbackMail, act.AfterCallbackSms, act.AfterCallbackMail, id)

	if err != nil {
		return
	}

	return op.GetAction(id)
}

func (op *Operator) DeleteTemplate(id int64) error {
	t, err := op.getTemplate(id)
	if err != nil {
		return err
	}

	if !op.IsAdmin() && op.User.Id != t.CreateUserId {
		return utils.EACCES
	}

	if err = op.RelCheck("SELECT count(*) FROM tag_tpl where tpl_id = ?", id); err != nil {
		return errors.New(err.Error() + "(tag - template)")
	}

	if _, err = op.SqlExec("delete from action where id = ?", t.ActionId); err != nil {
		return err
	}

	if _, err = op.SqlExec("delete from template where id = ?", id); err != nil {
		return err
	}

	DbLog(op.O, op.User.Id, CTL_M_TEMPLATE, id, CTL_A_DEL, "")

	return nil
}
