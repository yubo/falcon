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

type User struct {
	Id          int       `json:"id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	IM          string    `json:"im" orm:"column(im)"`
	QQ          string    `json:"qq" orm:"column(qq)"`
	Create_time time.Time `json:"-"`
}

func (u *User) AccessTid(scope string, tag int) error {
	return nil
}

func (u *User) Access(scope, tag string) (*Tag, error) {
	return nil, nil
}

func AddUser(u *User) (int, error) {
	id, err := orm.NewOrm().Insert(u)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	u.Id = int(id)
	cacheModule[CTL_M_USER].set(u.Id, u)
	return u.Id, err
}

func GetUser(id int) (*User, error) {
	if u, ok := cacheModule[CTL_M_USER].get(id).(*User); ok {
		return u, nil
	}
	u := &User{Id: id}
	err := orm.NewOrm().Read(u, "Id")
	if err == nil {
		cacheModule[CTL_M_USER].set(id, u)
	}
	return u, err
}

func GetUserByUuid(uuid string) (u *User, err error) {
	u = &User{Uuid: uuid}
	err = orm.NewOrm().Read(u, "Uuid")
	return u, err
}

func QueryUsers(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(User))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query).Or("Email__icontains", query))
	}
	return qs
}

func GetUsersCnt(query string) (int, error) {
	cnt, err := QueryUsers(query).Count()
	return int(cnt), err
}

func GetUsers(query string, limit, offset int) (users []*User, err error) {
	_, err = QueryUsers(query).Limit(limit, offset).All(&users)
	return
}

func UpdateUser(id int, _u *User) (u *User, err error) {
	if u, err = GetUser(id); err != nil {
		return nil, ErrNoUsr
	}

	if _u.Name != "" {
		u.Name = _u.Name
	}
	if _u.Cname != "" {
		u.Cname = _u.Cname
	}
	if _u.Email != "" {
		u.Email = _u.Email
	}
	if _u.Phone != "" {
		u.Phone = _u.Phone
	}
	if _u.IM != "" {
		u.IM = _u.IM
	}
	if _u.QQ != "" {
		u.QQ = _u.QQ
	}
	_, err = orm.NewOrm().Update(u)
	return u, err
}

func DeleteUser(id int) error {
	if n, err := orm.NewOrm().Delete(&User{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_USER].del(id)

	return nil
}
