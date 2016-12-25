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

type Scope struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	System_id   int       `json:"system_id"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func AddScope(s *Scope) (int, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	s.Id = int(id)
	cacheModule[CTL_M_SCOPE].set(s.Id, s)
	return s.Id, err
}

func GetScope(id int) (*Scope, error) {
	if s, ok := cacheModule[CTL_M_SCOPE].get(id).(*Scope); ok {
		return s, nil
	}
	s := &Scope{Id: id}
	err := orm.NewOrm().Read(s, "Id")
	if err == nil {
		cacheModule[CTL_M_SCOPE].set(id, s)
	}
	return s, err
}

func QueryScopes(sysid int, query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Scope)).Filter("System_id", sysid)
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func GetScopesCnt(sysid int, query string) (int, error) {
	cnt, err := QueryScopes(sysid, query).Count()
	return int(cnt), err
}

func GetScopes(sysid int, query string, limit, offset int) (scopes []*Scope, err error) {
	_, err = QueryScopes(sysid, query).Limit(limit, offset).All(&scopes)
	return
}

func UpdateScope(id int, _s *Scope) (s *Scope, err error) {
	if s, err = GetScope(id); err != nil {
		return nil, ErrNoScope
	}

	if _s.Name != "" {
		s.Name = _s.Name
	}
	if _s.Cname != "" {
		s.Cname = _s.Cname
	}
	/* not allowed
	if _s.System_id != 0 {
		s.Developers = _s.Developers
	}
	*/
	if _s.Note != "" {
		s.Note = _s.Note
	}
	_, err = orm.NewOrm().Update(s)
	return s, err
}

func DeleteScope(id int) error {

	if n, err := orm.NewOrm().Delete(&Scope{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_SCOPE].del(id)

	return nil
}
