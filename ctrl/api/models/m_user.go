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

type User struct {
	Id         int64     `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Im         string    `json:"im"`
	Qq         string    `json:"qq"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) IsAdmin() bool {
	// 1: system
	// 2: admin(first user)
	return (op.Token & SYS_F_A_TOKEN) != 0
}

func (op *Operator) AddUser(user *User) (*User, error) {
	user.Id = 0
	id, err := op.O.Insert(user)
	if err != nil {
		return nil, err
	}
	user.Id = id
	moduleCache[CTL_M_USER].set(id, user)

	DbLog(op.User.Id, CTL_M_USER, id, CTL_A_ADD, jsonStr(user))
	return user, nil
}

// just called from profileFilter()
func GetUser(id int64, o orm.Ormer) (*User, error) {
	if user, ok := moduleCache[CTL_M_USER].get(id).(*User); ok {
		return user, nil
	}
	user := &User{Id: id}
	err := o.Read(user, "Id")
	if err == nil {
		moduleCache[CTL_M_USER].set(id, user)
	}
	return user, err
}

func (op *Operator) GetUser(id int64) (*User, error) {
	return GetUser(id, op.O)
}

func (op *Operator) GetUserByUuid(uuid string) (user *User, err error) {
	user = &User{Uuid: uuid}
	err = op.O.Read(user, "Uuid")
	return user, err
}

func (op *Operator) QueryUsers(query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(User))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query).Or("Email__icontains", query))
	}
	return qs
}

func (op *Operator) GetUsersCnt(query string) (int64, error) {
	return op.QueryUsers(query).Count()
}

func (op *Operator) GetUsers(query string, limit, offset int) (users []*User, err error) {
	_, err = op.QueryUsers(query).Limit(limit, offset).All(&users)
	return
}

func (op *Operator) UpdateUser(id int64, _u *User) (user *User, err error) {
	if user, err = op.GetUser(id); err != nil {
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
	if _u.Im != "" {
		user.Im = _u.Im
	}
	if _u.Qq != "" {
		user.Qq = _u.Qq
	}
	_, err = op.O.Update(user)
	moduleCache[CTL_M_USER].set(id, user)
	DbLog(op.User.Id, CTL_M_USER, id, CTL_A_SET, "")
	return user, err
}

func (op *Operator) DeleteUser(id int64) error {
	if n, err := op.O.Delete(&User{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	moduleCache[CTL_M_USER].del(id)
	DbLog(op.User.Id, CTL_M_USER, id, CTL_A_DEL, "")

	return nil
}
