/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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
	SendSms            uint   `json:"sendSms"`
	SendMail           uint   `json:"sendMail"`
	Callback           uint   `json:"callback"`
	BeforeCallbackSms  uint   `json:"beforeCallbackSms"`
	BeforeCallbackMail uint   `json:"beforeCallbackMail"`
	AfterCallbackSms   uint   `json:"afterCallbackSms"`
	AfterCallbackMail  uint   `json:"afterCallbackMail"`
}

func (u *User) AddTemplate(o *Template) (id int64, err error) {
	o.CreateUserId = u.Id
	o.Id = 0
	id, err = orm.NewOrm().Insert(o)
	if err != nil {
		return
	}
	o.Id = id
	cacheModule[CTL_M_TEMPLATE].set(id, o)
	DbLog(u.Id, CTL_M_TEMPLATE, id, CTL_A_ADD, jsonStr(o))
	return
}

func (u *User) AddAction(o *Action) (id int64, err error) {
	o.Id = 0
	id, err = orm.NewOrm().Insert(o)
	if err != nil {
		return
	}
	o.Id = id
	return
}

func (u *User) CloneTemplate(id int64) (*TemplateAction, error) {
	var (
		ret     TemplateAction
		src_tpl *Template
		src_act *Action
		err     error
		objs    []*Strategy
		tid     int64
	)

	if src_tpl, err = u.getTemplate(id); err != nil {
		return nil, err
	}
	if src_act, err = u.GetAction(src_tpl.ActionId); err != nil {
		return nil, err
	}

	ret.Template = *src_tpl
	ret.Action = *src_act
	ret.Template.Name += "_copy"

	_, err = u.AddAction(&ret.Action)
	if err != nil {
		return nil, err
	}

	ret.Template.ActionId = ret.Action.Id
	tid, err = u.AddTemplate(&ret.Template)
	if err != nil {
		return nil, err
	}

	objs, _ = u.GetStrategys(src_tpl.Id, "", 0, 0)
	for _, obj := range objs {
		obj.TplId = tid
		u.AddStrategy(obj)
	}

	if t, err := u.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (u *User) GetTemplate(id int64) (*TemplateAction, error) {
	var ret TemplateAction

	if t, err := u.getTemplate(id); err != nil {
		return nil, err
	} else {
		ret.Template = *t
	}

	if a, err := u.GetAction(ret.Template.ActionId); err != nil {
		return nil, err
	} else {
		ret.Action = *a
	}

	if t, err := u.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (u *User) getTemplate(id int64) (*Template, error) {
	if r, ok := cacheModule[CTL_M_TEMPLATE].get(id).(*Template); ok {
		return r, nil
	}
	r := &Template{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		cacheModule[CTL_M_TEMPLATE].set(id, r)
	}
	return r, err
}

func (u *User) GetAction(id int64) (*Action, error) {
	a := &Action{Id: id}
	err := orm.NewOrm().Read(a, "Id")
	return a, err
}

func (u *User) GetTemplatesCnt(query string) (cnt int64, err error) {
	if query == "" {
		err = orm.NewOrm().Raw("SELECT count(*) FROM template").QueryRow(&cnt)
	} else {
		err = orm.NewOrm().Raw("SELECT count(*) FROM template WHERE name like ?", "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (u *User) GetTemplates(query string, limit, offset int) (ret []TemplateUi, err error) {
	if query == "" {
		_, err = orm.NewOrm().Raw("SELECT a.id, b.id as pid, a.name, b.name as pname, c.name as creator  FROM template a LEFT JOIN template b ON a.parent_id = b.id LEFT JOIN user c ON a.create_user_id = c.id ORDER BY a.name LIMIT ? OFFSET ?", limit, offset).QueryRows(&ret)
	} else {
		_, err = orm.NewOrm().Raw("SELECT a.id as id, b.id as pid, a.name as name, b.name as pname, c.name as creator  FROM template a LEFT JOIN template b ON a.parent_id = b.id LEFT JOIN user c ON a.create_user_id = c.id WHERE a.name like ? ORDER BY a.name LIMIT ? OFFSET ?", "%"+query+"%", limit, offset).QueryRows(&ret)
	}
	return
}

func (u *User) UpdateTemplate(id int64, _o *TemplateAction) (o *Template, err error) {
	var t *Template
	t, err = u.updateTemplate(id, &_o.Template)
	if err != nil {
		return nil, err
	}
	_, err = u.UpdateAction(t.ActionId, &_o.Action)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (u *User) updateTemplate(id int64, _o *Template) (o *Template, err error) {
	if o, err = u.getTemplate(id); err != nil {
		return nil, ErrNoTemplate
	}

	if _o.Name != "" {
		o.Name = _o.Name
	}
	if _o.ParentId != 0 {
		o.ParentId = _o.ParentId
	}
	_, err = orm.NewOrm().Update(o)
	return o, err
}

func (u *User) UpdateAction(id int64, _o *Action) (o *Action, err error) {
	if o, err = u.GetAction(id); err != nil {
		return nil, ErrNoTemplate
	}
	if _o.Uic != "" {
		o.Uic = _o.Uic
	}
	if _o.Url != "" {
		o.Url = _o.Url
	}
	o.SendSms = _o.SendSms
	o.SendMail = _o.SendMail
	o.Callback = _o.Callback
	o.BeforeCallbackSms = _o.BeforeCallbackSms
	o.BeforeCallbackMail = _o.BeforeCallbackMail
	o.AfterCallbackSms = _o.AfterCallbackSms
	o.AfterCallbackMail = _o.AfterCallbackMail
	_, err = orm.NewOrm().Update(o)
	return o, err
}

func (u *User) DeleteTemplate(id int64) error {
	template, err := u.getTemplate(id)
	if err != nil {
		return err
	}

	if _, err = orm.NewOrm().Delete(&Action{Id: template.ActionId}); err != nil {
		return err
	}

	if _, err = orm.NewOrm().Delete(&Template{Id: id}); err != nil {
		return err
	}

	cacheModule[CTL_M_TEMPLATE].del(id)
	DbLog(u.Id, CTL_M_TEMPLATE, id, CTL_A_DEL, "")

	return nil
}
