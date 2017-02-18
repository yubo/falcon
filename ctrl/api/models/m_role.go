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

type Role struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (u *User) AddRole(r *Role) (id int64, err error) {
	r.Id = 0
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	moduleCache[CTL_M_ROLE].set(id, r)
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_ADD, jsonStr(r))
	return
}

func (u *User) GetRole(id int64) (*Role, error) {
	if r, ok := moduleCache[CTL_M_ROLE].get(id).(*Role); ok {
		return r, nil
	}
	r := &Role{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		moduleCache[CTL_M_ROLE].set(id, r)
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

func (u *User) GetRolesCnt(query string) (int64, error) {
	return u.QueryRoles(query).Count()
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
	moduleCache[CTL_M_ROLE].set(id, r)
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_SET, "")
	return r, err
}

func (u *User) DeleteRole(id int64) error {
	if n, err := orm.NewOrm().Delete(&Role{Id: id}); err != nil || n == 0 {
		return err
	}
	moduleCache[CTL_M_ROLE].del(id)
	DbLog(u.Id, CTL_M_ROLE, id, CTL_A_DEL, "")

	return nil
}

func (u *User) BindUserRole(user_id, role_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_role_user` (`tag_id`, `role_id`, `user_id`) VALUES (?, ?, ?)", tag_id, role_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) BindTokenRole(token_id, role_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_role_token` (`tag_id`, `role_id`, `token_id`) VALUES (?, ?, ?)", tag_id, role_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
