/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserProfileUpdate struct {
	Cname string `json:"cname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Im    string `json:"im"`
	Qq    string `json:"qq"`
	Extra string `json:"extra"`
}

type User struct {
	Id         int64     `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Im         string    `json:"im"`
	Qq         string    `json:"qq"`
	Extra      string    `json:"extra"`
	Disabled   int       `json:"disabled"`
	CreateTime time.Time `json:"ctime"`
}

type Operator struct {
	User  *User
	Token int
	O     orm.Ormer
}

type OperatorInfo struct {
	User     *User `json:"user"`
	Reader   bool  `json:"reader"`
	Operator bool  `json:"operator"`
	Admin    bool  `json:"admin"`
}

func (op *Operator) Info() *OperatorInfo {
	if op.User == nil {
		return &OperatorInfo{}
	}
	return &OperatorInfo{
		User:     op.User,
		Reader:   op.IsReader(),
		Operator: op.IsOperator(),
		Admin:    op.IsAdmin(),
	}
}

func (op *Operator) IsAdmin() bool {
	return (op.Token & SYS_F_A_TOKEN) != 0
}

func (op *Operator) IsOperator() bool {
	return (op.Token & SYS_F_O_TOKEN) != 0
}

func (op *Operator) IsReader() bool {
	return (op.Token & SYS_F_R_TOKEN) != 0
}

func (op *Operator) UserTokens() (token int) {
	var (
		tids []int64
	)

	if op.User == nil {
		return 0
	}

	if admin[op.User.Name] {
		return SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_R_TOKEN
	}

	_, err := op.O.Raw("SELECT b1.token_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id in (?, ?, ?) AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id GROUP BY b1.token_id",
		TPL_REL_T_ACL_USER, op.User.Id, TPL_REL_T_ACL_TOKEN,
		SYS_IDX_R_TOKEN, SYS_IDX_O_TOKEN,
		SYS_IDX_A_TOKEN).QueryRows(&tids)
	if err != nil {
		return 0
	}
	for _, tid := range tids {
		switch tid {
		case SYS_IDX_R_TOKEN:
			token |= SYS_F_R_TOKEN
		case SYS_IDX_O_TOKEN:
			token |= SYS_F_O_TOKEN
		case SYS_IDX_A_TOKEN:
			token |= SYS_F_A_TOKEN
		}
	}

	if op.User.Id < 3 {
		token |= SYS_F_A_TOKEN
	}

	if token&SYS_F_A_TOKEN != 0 {
		token |= SYS_F_O_TOKEN
	}
	if token&SYS_F_O_TOKEN != 0 {
		token |= SYS_F_R_TOKEN
	}

	return token
}

func (op *Operator) AddUser(user *User) (*User, error) {
	user.Id = 0
	user.CreateTime = time.Now()
	id, err := op.SqlInsert("insert user (uuid, name, cname, email, phone, im, qq, disabled) values (?, ?, ?, ?, ?, ?, ?, ?)", user.Uuid, user.Name, user.Cname, user.Email, user.Phone, user.Im, user.Qq, user.Disabled)
	if err != nil {
		return nil, err
	}
	user.Id = id
	moduleCache[CTL_M_USER].set(id, user)

	DbLog(op.O, op.User.Id, CTL_M_USER, id, CTL_A_ADD, jsonStr(user))
	return user, nil
}

// just called from profileFilter()
func GetUser(id int64, o orm.Ormer) (ret *User, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_USER].get(id).(*User); ok {
		return ret, nil
	}

	ret = &User{}
	err = o.Raw("select id, uuid, name, cname, email, phone, im, qq, disabled, extra, create_time from user where id = ?", id).QueryRow(ret)
	if err == nil {
		moduleCache[CTL_M_USER].set(id, ret)
	}
	return
}

func (op *Operator) GetUser(id int64) (*User, error) {
	return GetUser(id, op.O)
}

func (op *Operator) GetUserByUuid(uuid string) (ret *User, err error) {
	ret = &User{}
	err = op.SqlRow(ret, "select id, uuid, name, cname, email, phone, im, qq, disabled, create_time from user where uuid = ?", uuid)
	return ret, err
}

func sqlUser(query string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "name like ? or email like ?")
		sql3 = append(sql3, "%"+query+"%", "%"+query+"%")
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetUsersCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlUser(query)
	err = op.O.Raw("select count(*) from user "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetUsers(query string, limit, offset int) (ret []*User, err error) {
	sql, sql_args := sqlUser(query)
	sql = sqlLimit("select id, uuid, name, cname, email, phone, im, qq, disabled, create_time, extra from user "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)

	return
}

func (op *Operator) UpdateUser(user *User) (ret *User, err error) {
	_, err = op.SqlExec("update user set name = ?, cname = ?, email = ?, phone = ?, im = ?, qq = ?, disabled = ?, extra = ? where id = ?", user.Name, user.Cname, user.Email, user.Phone, user.Im, user.Qq, user.Disabled, user.Extra, user.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_USER].del(user.Id)
	if ret, err = op.GetUser(user.Id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_USER, user.Id, CTL_A_SET, "")
	return
}

func (op *Operator) DeleteUser(id int64) error {
	if err := op.RelCheck("SELECT count(*) FROM tpl_rel where sub_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_USER); err != nil {
		return errors.New(err.Error() + "(tag - role - user)")
	}

	if _, err := op.SqlExec("delete from user where id = ?", id); err != nil {
		return err
	}
	moduleCache[CTL_M_USER].del(id)
	DbLog(op.O, op.User.Id, CTL_M_USER, id, CTL_A_DEL, "")

	return nil
}
