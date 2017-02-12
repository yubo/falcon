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
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (u *User) AddToken(o *Token) (id int64, err error) {
	o.Id = 0
	if id, err = orm.NewOrm().Insert(o); err != nil {
		return
	}
	o.Id = id
	cacheModule[CTL_M_TOKEN].set(id, o)
	DbLog(u.Id, CTL_M_TOKEN, id, CTL_A_ADD, jsonStr(o))
	return
}

func (u *User) GetToken(id int64) (o *Token, err error) {
	var ok bool

	if o, ok = cacheModule[CTL_M_TOKEN].get(id).(*Token); ok {
		return
	}
	o = &Token{Id: id}
	err = orm.NewOrm().Read(o, "Id")
	if err == nil {
		cacheModule[CTL_M_TOKEN].set(id, o)
	}
	return
}

func (u *User) GetTokenByName(token string) (o *Token, err error) {
	o = &Token{Name: token}
	err = orm.NewOrm().Read(o, "Name")
	return
}

func (u *User) QueryTokens(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Token))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetTokensCnt(query string) (int64, error) {
	return u.QueryTokens(query).Count()
}

func (u *User) GetTokens(query string, limit, offset int) (tokens []*Token, err error) {
	_, err = u.QueryTokens(query).Limit(limit, offset).All(&tokens)
	return
}

func (u *User) UpdateToken(id int64, _tk *Token) (tk *Token, err error) {
	if tk, err = u.GetToken(id); err != nil {
		return nil, ErrNoToken
	}

	if _tk.Name != "" {
		tk.Name = _tk.Name
	}
	if _tk.Cname != "" {
		tk.Cname = _tk.Cname
	}
	if _tk.Note != "" {
		tk.Note = _tk.Note
	}
	_, err = orm.NewOrm().Update(tk)
	cacheModule[CTL_M_TOKEN].set(id, tk)
	DbLog(u.Id, CTL_M_TOKEN, id, CTL_A_SET, "")
	return tk, err
}

func (u *User) DeleteToken(id int64) error {

	if n, err := orm.NewOrm().Delete(&Token{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_TOKEN].del(id)
	DbLog(u.Id, CTL_M_TOKEN, id, CTL_A_DEL, "")

	return nil
}