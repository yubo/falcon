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

type System struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Developers  string    `json:"developers"`
	Email       string    `json:"email"`
	Create_time time.Time `json:"-"`
}

func AddSystem(s *System) (int, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	s.Id = int(id)
	cacheModule[CTL_M_SYSTEM].set(s.Id, s)
	return s.Id, err
}

func GetSystem(id int) (*System, error) {
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

func QuerySystems(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(System))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query))
	}
	return qs
}

func GetSystemsCnt(query string) (int, error) {
	cnt, err := QuerySystems(query).Count()
	return int(cnt), err
}

func GetSystems(query string, limit, offset int) (systems []*System, err error) {
	_, err = QuerySystems(query).Limit(limit, offset).All(&systems)
	return
}

func UpdateSystem(id int, _s *System) (s *System, err error) {
	if s, err = GetSystem(id); err != nil {
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
	return s, err
}

func DeleteSystem(id int) error {

	if n, err := orm.NewOrm().Delete(&System{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_SYSTEM].del(id)

	return nil
}
