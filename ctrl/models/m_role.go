/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func AddRole(r *Role) (int, error) {
	id, err := orm.NewOrm().Insert(r)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	r.Id = int(id)
	cacheModule[CTL_M_ROLE].set(r.Id, r)
	return r.Id, err
}

func GetRole(id int) (*Role, error) {
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

func QueryRoles(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Role))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func GetRolesCnt(query string) (int, error) {
	cnt, err := QueryRoles(query).Count()
	return int(cnt), err
}

func GetRoles(query string, limit, offset int) (roles []*Role, err error) {
	_, err = QueryRoles(query).Limit(limit, offset).All(&roles)
	return
}

func UpdateRole(id int, _r *Role) (r *Role, err error) {
	if r, err = GetRole(id); err != nil {
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
	return r, err
}

func DeleteRole(id int) error {

	if n, err := orm.NewOrm().Delete(&Role{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_ROLE].del(id)

	return nil
}
