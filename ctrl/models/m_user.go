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

type User struct {
	Id          int64     `json:"id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Cname       string    `json:"cname"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	IM          string    `json:"im" orm:"column(im)"`
	QQ          string    `json:"qq" orm:"column(qq)"`
	Create_time time.Time `json:"-"`
}

// return nil *Tag if tag not exist
func (u *User) Access(scope, tag string, chkExist bool) (t *Tag, err error) {
	var s *Scope
	var tag_id int64

	// if not found, err will be return <QuerySeter> no row found
	if t, err = u.GetTagByName(tag); chkExist && err != nil {
		return
	}

	if u.IsAdmin() {
		return
	}

	if s, err = u.GetScopeByName(scope); err != nil {
		return
	}

	// TODO: test
	err = orm.NewOrm().Raw(`
SELECT sup_tag_id
FROM tag_rel
WHERE sup_tag_id IN (
    SELECT user_scope.user_tag_id
    FROM (SELECT a.tag_id AS user_tag_id,
                b.tag_id AS scope_tag_id,
                a.role_id AS role_id,
                a.user_id AS user_id,
                b.scope_id AS scope_id
          FROM tag_role_user a
          JOIN tag_role_scope b
          ON a.user_id = ? AND b.scope_id = ? AND a.role_id = b.role_id) user_scope
   JOIN tag_rel 
   ON user_scope.user_tag_id = tag_rel.tag_id AND user_scope.scope_tag_id = tag_rel.sup_tag_id
   GROUP BY user_scope.user_tag_id)
AND tag_id = ?
`,
		u.Id, s.Id, t.Id).QueryRow(&tag_id)

	return
}

func access(user_id, scope_id, tag_id int64) (tid int64, err error) {
	// TODO: test
	err = orm.NewOrm().Raw(`
SELECT b2.scope_tag_id
FROM tag_rel a2
JOIN (
    SELECT a1.scope_tag_id, a1.user_tag_id
    FROM (SELECT a0.tag_id AS user_tag_id,
                b0.tag_id AS scope_tag_id,
                a0.role_id AS role_id,
                a0.user_id AS user_id,
                b0.scope_id AS scope_id
          FROM tag_role_user a0
          JOIN tag_role_scope b0
          ON a0.user_id = ? AND b0.scope_id = ? AND a0.role_id = b0.role_id) a1
    JOIN tag_rel b1
    ON a1.user_tag_id = b1.tag_id AND a1.scope_tag_id = b1.sup_tag_id
    GROUP BY a1.user_tag_id 
    HAVING a1.scope_tag_id = MAX(a1.scope_tag_id)) b2
ON  a2.sup_tag_id = b2.user_tag_id 
WHERE a2.tag_id = ?
`,
		user_id, scope_id, tag_id).QueryRow(&tid)
	if err != nil {
		err = EACCES
	}

	return
}

func (u *User) IsAdmin() bool {
	return u.Id == 1
}

func (u *User) AddUser(user *User) (id int64, err error) {
	if id, err = orm.NewOrm().Insert(user); err != nil {
		return
	}
	user.Id = id
	cacheModule[CTL_M_USER].set(id, user)

	data, _ := json.Marshal(user)
	DbLog(u.Id, CTL_M_USER, id, CTL_A_ADD, data)
	return
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

func (u *User) GetUsersCnt(query string) (int, error) {
	cnt, err := u.QueryUsers(query).Count()
	return int(cnt), err
}

func (u *User) GetUsers(query string, limit, offset int) (users []*User, err error) {
	_, err = u.QueryUsers(query).Limit(limit, offset).All(&users)
	return
}

func (u *User) UpdateUser(id int64, _u *User) (user *User, err error) {
	if user, err = u.GetUser(id); err != nil {
		return nil, ErrNoUsr
	}

	if _u.Name != "" {
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
	DbLog(u.Id, CTL_M_USER, id, CTL_A_SET, nil)
	return user, err
}

func (u *User) DeleteUser(id int64) error {
	if n, err := orm.NewOrm().Delete(&User{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_USER].del(id)
	DbLog(u.Id, CTL_M_USER, id, CTL_A_DEL, nil)

	return nil
}
