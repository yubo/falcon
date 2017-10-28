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

type Token struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

type TokenCreate struct {
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

type TokenUpdate struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

func (op *Operator) CreateToken(o *TokenCreate) (id int64, err error) {
	id, err = op.SqlInsert("insert token (name, cname, note) values (?, ?, ?)", o.Name, o.Cname, o.Note)
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_TOKEN, id, CTL_A_ADD, jsonStr(o))
	return
}

func (op *Operator) GetToken(id int64) (ret *Token, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_TOKEN].get(id).(*Token); ok {
		return
	}

	ret = &Token{}
	err = op.SqlRow(ret, "select id, name, cname, note, create_time from token where id = ?", id)
	if err == nil {
		moduleCache[CTL_M_TOKEN].set(id, ret, ret.Name)
	}
	return
}

func (op *Operator) GetTokenByName(token string) (ret *Token, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_TOKEN].getByKey(token).(*Token); ok {
		return
	}

	ret = &Token{}
	err = op.SqlRow(ret, "select id, name, cname, note, create_time from token where name = ?", token)
	if err == nil {
		moduleCache[CTL_M_TOKEN].set(ret.Id, ret, ret.Name)
	}
	return
}

func (op *Operator) GetTokensCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlName(query)
	err = op.O.Raw("SELECT count(*) FROM token "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTokens(query string, limit, offset int) (ret []*Token, err error) {
	sql, sql_args := sqlName(query)
	sql = sqlLimit("select id, name, cname, note, create_time from token "+sql+" ORDER BY id", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateToken(t *Token) (ret *Token, err error) {
	_, err = op.SqlExec("update token set name = ?, cname = ?, note = ? where id = ?", t.Name, t.Cname, t.Note, t.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_TOKEN].del(t.Id)
	if ret, err = op.GetToken(t.Id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_TOKEN, t.Id, CTL_A_SET, jsonStr(t))
	return ret, err
}

func (op *Operator) DeleteToken(id int64) error {
	if err := op.RelCheck("SELECT count(*) FROM tpl_rel where sub_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_TOKEN); err != nil {
		return errors.New(err.Error() + "(tag - role - token)")
	}

	if _, err := op.SqlExec("delete from token where id = ?", id); err != nil {
		return err
	}
	moduleCache[CTL_M_TOKEN].del(id)
	DbLog(op.O, op.User.Id, CTL_M_TOKEN, id, CTL_A_DEL, "")

	return nil
}
