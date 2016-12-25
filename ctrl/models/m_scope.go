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

type Scope struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	System_id   int64     `json:"system_id"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddScope(s *Scope) (id int64, err error) {
	if id, err = orm.NewOrm().Insert(s); err != nil {
		return
	}
	s.Id = id
	cacheModule[CTL_M_SCOPE].set(id, s)
	data, _ := json.Marshal(s)
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_ADD, data)
	return
}

func (u *User) GetScope(id int64) (s *Scope, err error) {
	var ok bool

	if s, ok = cacheModule[CTL_M_SCOPE].get(id).(*Scope); ok {
		return
	}
	s = &Scope{Id: id}
	err = orm.NewOrm().Read(s, "Id")
	if err == nil {
		cacheModule[CTL_M_SCOPE].set(id, s)
	}
	return
}

func (u *User) GetScopeByName(scope string) (s *Scope, err error) {
	s = &Scope{Name: scope}
	err = orm.NewOrm().Read(s, "Name")
	return
}

func (u *User) QueryScopes(sysid int64, query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Scope)).Filter("System_id", sysid)
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetScopesCnt(sysid int64, query string) (int, error) {
	cnt, err := u.QueryScopes(sysid, query).Count()
	return int(cnt), err
}

func (u *User) GetScopes(sysid int64, query string, limit, offset int) (scopes []*Scope, err error) {
	_, err = u.QueryScopes(sysid, query).Limit(limit, offset).All(&scopes)
	return
}

func (u *User) UpdateScope(id int64, _s *Scope) (s *Scope, err error) {
	if s, err = u.GetScope(id); err != nil {
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
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_SET, nil)
	return s, err
}

func (u *User) DeleteScope(id int64) error {

	if n, err := orm.NewOrm().Delete(&Scope{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_SCOPE].del(id)
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_DEL, nil)

	return nil
}
