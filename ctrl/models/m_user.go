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

type Operator struct {
	id int64
	o  orm.Ormer
	u  *User
}

type User struct {
	Id          int64     `json:"id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	IM          string    `json:"im" orm:"column(im)"`
	QQ          string    `json:"qq" orm:"column(qq)"`
	Create_time time.Time `json:"ctime"`
}

func (u *User) IsAdmin() bool {
	// 1: system
	// 2: admin(first user)
	return u.Id < 1000
}

func (u *User) AddUser(user *User) (*User, error) {
	id, err := orm.NewOrm().Insert(user)
	if err != nil {
		return nil, err
	}
	user.Id = id
	cacheModule[CTL_M_USER].set(id, user)
	beego.Debug(u)

	DbLog(u.Id, CTL_M_USER, id, CTL_A_ADD, jsonStr(user))
	return user, nil
}

// just called from profileFilter()
func GetUser(id int64) (*User, error) {
	if user, ok := cacheModule[CTL_M_USER].get(id).(*User); ok {
		return user, nil
	}
	user := &User{Id: id}
	err := orm.NewOrm().Read(user, "Id")
	if err == nil {
		cacheModule[CTL_M_USER].set(id, user)
	}
	return user, err
}

func (u *User) GetUser(id int64) (*User, error) {
	return GetUser(id)
}

func (u *User) GetUserByUuid(uuid string) (user *User, err error) {
	user = &User{Uuid: uuid}
	err = orm.NewOrm().Read(user, "Uuid")
	return user, err
}

func (u *User) QueryUsers(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(User))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query).Or("Email__icontains", query))
	}
	return qs
}

func (u *User) GetUsersCnt(query string) (int64, error) {
	return u.QueryUsers(query).Count()
}

func (u *User) GetUsers(query string, limit, offset int) (users []*User, err error) {
	_, err = u.QueryUsers(query).Limit(limit, offset).All(&users)
	return
}

func (u *User) UpdateUser(id int64, _u *User) (user *User, err error) {
	if user, err = u.GetUser(id); err != nil {
		return nil, ErrNoUsr
	}

	if _u.Name != "" && user.Name == "" {
		user.Name = _u.Name
	}
	if _u.Cname != "" {
		user.Cname = _u.Cname
	}
	if _u.Email != "" {
		user.Email = _u.Email
	}
	if _u.Phone != "" {
		user.Phone = _u.Phone
	}
	if _u.IM != "" {
		user.IM = _u.IM
	}
	if _u.QQ != "" {
		user.QQ = _u.QQ
	}
	_, err = orm.NewOrm().Update(user)
	cacheModule[CTL_M_USER].set(id, user)
	DbLog(u.Id, CTL_M_USER, id, CTL_A_SET, "")
	return user, err
}

func (u *User) DeleteUser(id int64) error {
	if n, err := orm.NewOrm().Delete(&User{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_USER].del(id)
	DbLog(u.Id, CTL_M_USER, id, CTL_A_DEL, "")

	return nil
}
