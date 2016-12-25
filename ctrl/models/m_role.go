/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddRole(r *Role) (id int64, err error) {
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	cacheModule[CTL_M_ROLE].set(id, r)
	data, _ := json.Marshal(r)
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_ADD, data)
	return
}

func (u *User) GetRole(id int64) (*Role, error) {
	if r, ok := cacheModule[CTL_M_ROLE].get(id).(*Role); ok {
		return r, nil
	}
	r := &Role{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		cacheModule[CTL_M_ROLE].set(id, r)
	}
	return r, err
}

func (u *User) QueryRoles(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Role))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetRolesCnt(query string) (int, error) {
	cnt, err := u.QueryRoles(query).Count()
	return int(cnt), err
}

func (u *User) GetRoles(query string, limit, offset int) (roles []*Role, err error) {
	_, err = u.QueryRoles(query).Limit(limit, offset).All(&roles)
	return
}

func (u *User) UpdateRole(id int64, _r *Role) (r *Role, err error) {
	if r, err = u.GetRole(id); err != nil {
		return nil, ErrNoRole
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
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_SET, nil)
	return r, err
}

func (u *User) DeleteRole(id int64) error {

	if n, err := orm.NewOrm().Delete(&Role{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_ROLE].del(id)
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_DEL, nil)

	return nil
}
