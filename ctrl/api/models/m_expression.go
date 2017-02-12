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

type Expression struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (u *User) AddExpression(r *Expression) (id int64, err error) {
	r.Id = 0
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	cacheModule[CTL_M_EXPRESSION].set(id, r)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_ADD, jsonStr(r))
	return
}

func (u *User) GetExpression(id int64) (*Expression, error) {
	if r, ok := cacheModule[CTL_M_EXPRESSION].get(id).(*Expression); ok {
		return r, nil
	}
	r := &Expression{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		cacheModule[CTL_M_EXPRESSION].set(id, r)
	}
	return r, err
}

func (u *User) QueryExpressions(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Expression))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetExpressionsCnt(query string) (int, error) {
	cnt, err := u.QueryExpressions(query).Count()
	return int(cnt), err
}

func (u *User) GetExpressions(query string, limit, offset int) (expressions []*Expression, err error) {
	_, err = u.QueryExpressions(query).Limit(limit, offset).All(&expressions)
	return
}

func (u *User) UpdateExpression(id int64, _r *Expression) (r *Expression, err error) {
	if r, err = u.GetExpression(id); err != nil {
		return nil, ErrNoExpression
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
	cacheModule[CTL_M_EXPRESSION].set(id, r)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_SET, "")
	return r, err
}

func (u *User) DeleteExpression(id int64) error {
	if n, err := orm.NewOrm().Delete(&Expression{Id: id}); err != nil || n == 0 {
		return err
	}
	cacheModule[CTL_M_EXPRESSION].del(id)
	DbLog(u.Id, CTL_M_EXPRESSION, id, CTL_A_DEL, "")

	return nil
}

func (u *User) BindUserExpression(user_id, expression_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_expression_user` (`tag_id`, `expression_id`, `user_id`) VALUES (?, ?, ?)", tag_id, expression_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) BindTokenExpression(token_id, expression_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_expression_token` (`tag_id`, `expression_id`, `token_id`) VALUES (?, ?, ?)", tag_id, expression_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
