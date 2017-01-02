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

type Token struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddToken(s *Token) (id int64, err error) {
	if id, err = orm.NewOrm().Insert(s); err != nil {
		return
	}
	s.Id = id
	cacheModule[CTL_M_SCOPE].set(id, s)
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_ADD, jsonStr(s))
	return
}

func (u *User) GetToken(id int64) (s *Token, err error) {
	var ok bool

	if s, ok = cacheModule[CTL_M_SCOPE].get(id).(*Token); ok {
		return
	}
	s = &Token{Id: id}
	err = orm.NewOrm().Read(s, "Id")
	if err == nil {
		cacheModule[CTL_M_SCOPE].set(id, s)
	}
	return
}

func (u *User) GetTokenByName(token string) (s *Token, err error) {
	s = &Token{Name: token}
	err = orm.NewOrm().Read(s, "Name")
	return
}

func (u *User) QueryTokens(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Token))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetTokensCnt(query string) (int, error) {
	cnt, err := u.QueryTokens(query).Count()
	return int(cnt), err
}

func (u *User) GetTokens(query string, limit, offset int) (tokens []*Token, err error) {
	_, err = u.QueryTokens(query).Limit(limit, offset).All(&tokens)
	return
}

func (u *User) UpdateToken(id int64, _s *Token) (s *Token, err error) {
	if s, err = u.GetToken(id); err != nil {
		return nil, ErrNoToken
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
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_SET, "")
	return s, err
}

func (u *User) DeleteToken(id int64) error {

	if n, err := orm.NewOrm().Delete(&Token{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_SCOPE].del(id)
	DbLog(u.Id, CTL_M_SCOPE, id, CTL_A_DEL, "")

	return nil
}
