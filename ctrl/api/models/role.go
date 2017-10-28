/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"time"
)

type Role struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

type RoleUpdate struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

func (op *Operator) AddRole(r *Role) (id int64, err error) {
	id, err = op.SqlInsert("insert role (name, cname, note) values (?, ?, ?)", r.Name, r.Cname, r.Note)
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_ROLE, id, CTL_A_ADD, jsonStr(r))
	return
}

func (op *Operator) GetRole(id int64) (ret *Role, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_ROLE].get(id).(*Role); ok {
		return
	}

	ret = &Role{}
	err = op.SqlRow(ret, "select id, name, cname, note, create_time from role where id = ?", id)
	if err == nil {
		moduleCache[CTL_M_ROLE].set(id, ret)
	}
	return
}

func (op *Operator) GetRolesCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlName(query)
	err = op.O.Raw("SELECT count(*) FROM role "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetRoles(query string, limit, offset int) (ret []*Role, err error) {
	sql, sql_args := sqlName(query)
	sql = sqlLimit("select id, name, cname, note, create_time from role "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateRole(role *Role) (ret *Role, err error) {
	_, err = op.SqlExec("update role set name = ?, cname = ?, note = ? where id = ?", role.Name, role.Cname, role.Note, role.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_ROLE].del(role.Id)
	if ret, err = op.GetRole(role.Id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_ROLE, role.Id, CTL_A_SET, "")
	return ret, err
}

func (op *Operator) DeleteRole(id int64) error {

	if err := op.RelCheck("SELECT count(*) FROM tpl_rel where tpl_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_USER); err != nil {
		return errors.New(err.Error() + "(tag - role - user)")
	}

	if err := op.RelCheck("SELECT count(*) FROM tpl_rel where tpl_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_TOKEN); err != nil {
		return errors.New(err.Error() + "(tag - role - token)")
	}

	if _, err := op.SqlExec("delete from role where id = ?", id); err != nil {
		return err
	}
	moduleCache[CTL_M_ROLE].del(id)
	DbLog(op.O, op.User.Id, CTL_M_ROLE, id, CTL_A_DEL, "")

	return nil
}
