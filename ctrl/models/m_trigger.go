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

type Trigger struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddTrigger(r *Trigger) (id int64, err error) {
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	cacheModule[CTL_M_TRIGGER].set(id, r)
	DbLog(u.Id, CTL_M_TRIGGER, id, CTL_A_ADD, jsonStr(r))
	return
}

func (u *User) GetTrigger(id int64) (*Trigger, error) {
	if r, ok := cacheModule[CTL_M_TRIGGER].get(id).(*Trigger); ok {
		return r, nil
	}
	r := &Trigger{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		cacheModule[CTL_M_TRIGGER].set(id, r)
	}
	return r, err
}

func (u *User) QueryTriggers(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Trigger))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetTriggersCnt(query string) (int, error) {
	cnt, err := u.QueryTriggers(query).Count()
	return int(cnt), err
}

func (u *User) GetTriggers(query string, limit, offset int) (triggers []*Trigger, err error) {
	_, err = u.QueryTriggers(query).Limit(limit, offset).All(&triggers)
	return
}

func (u *User) UpdateTrigger(id int64, _r *Trigger) (r *Trigger, err error) {
	if r, err = u.GetTrigger(id); err != nil {
		return nil, ErrNoTrigger
	}

	if _r.Name != "" {
		r.Name = _r.Name
	}
	if _r.Cname != "" {
		r.Cname = _r.Cname
	}
	if _r.Note != "" {
		r.Note = _r.Note
	}
	_, err = orm.NewOrm().Update(r)
	cacheModule[CTL_M_TRIGGER].set(id, r)
	DbLog(u.Id, CTL_M_TRIGGER, id, CTL_A_SET, "")
	return r, err
}

func (u *User) DeleteTrigger(id int64) error {
	if n, err := orm.NewOrm().Delete(&Trigger{Id: id}); err != nil || n == 0 {
		return err
	}
	cacheModule[CTL_M_TRIGGER].del(id)
	DbLog(u.Id, CTL_M_TRIGGER, id, CTL_A_DEL, "")

	return nil
}

func (u *User) BindUserTrigger(user_id, trigger_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_trigger_user` (`tag_id`, `trigger_id`, `user_id`) VALUES (?, ?, ?)", tag_id, trigger_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) BindTokenTrigger(token_id, trigger_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_trigger_token` (`tag_id`, `trigger_id`, `token_id`) VALUES (?, ?, ?)", tag_id, trigger_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
