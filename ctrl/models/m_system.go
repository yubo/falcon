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

type System struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Developers  string    `json:"developers"`
	Email       string    `json:"email"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddSystem(s *System) (id int64, err error) {

	if u.IsAdmin() {
		return 0, EACCES
	}

	if id, err = orm.NewOrm().Insert(s); err != nil {
		return
	}
	s.Id = id
	cacheModule[CTL_M_SYSTEM].set(id, s)
	data, _ := json.Marshal(s)
	DbLog(u.Id, CTL_M_SYSTEM, id, CTL_A_ADD, data)
	return
}

func (u *User) GetSystem(id int64) (*System, error) {
	if s, ok := cacheModule[CTL_M_SYSTEM].get(id).(*System); ok {
		return s, nil
	}
	s := &System{Id: id}
	err := orm.NewOrm().Read(s, "Id")
	if err == nil {
		cacheModule[CTL_M_SYSTEM].set(id, s)
	}
	return s, err
}

func (u *User) QuerySystems(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(System))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetSystemsCnt(query string) (int, error) {
	cnt, err := u.QuerySystems(query).Count()
	return int(cnt), err
}

func (u *User) GetSystems(query string, limit, offset int) (systems []*System, err error) {
	_, err = u.QuerySystems(query).Limit(limit, offset).All(&systems)
	return
}

func (u *User) UpdateSystem(id int64, _s *System) (s *System, err error) {
	if s, err = u.GetSystem(id); err != nil {
		return nil, ErrNoSystem
	}

	if _s.Name != "" {
		s.Name = _s.Name
	}
	if _s.Cname != "" {
		s.Cname = _s.Cname
	}
	if _s.Developers != "" {
		s.Developers = _s.Developers
	}
	if _s.Email != "" {
		s.Email = _s.Email
	}
	_, err = orm.NewOrm().Update(s)
	DbLog(u.Id, CTL_M_SYSTEM, id, CTL_A_SET, nil)
	return s, err
}

func (u *User) DeleteSystem(id int64) error {

	if n, err := orm.NewOrm().Delete(&System{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_SYSTEM].del(id)
	DbLog(u.Id, CTL_M_SYSTEM, id, CTL_A_DEL, nil)

	return nil
}
