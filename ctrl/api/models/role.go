/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"database/sql"
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

type Role struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

type RoleCreate struct {
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

type RoleUpdate struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
}

func (op *Operator) CreateRole(r *RoleCreate) (id int64, err error) {
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
	err = op.SqlRow(ret, "select id, name, cname, note from role where id = ?", id)
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
	sql = sqlLimit("select id, name, cname, note from role "+sql+" ORDER BY name", limit, offset)
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

/*******************************************************************************
 ************************ tag role / acl ***************************************
 ******************************************************************************/
/*
- 1 : m
  * tag - tag
  * tag - host
  * tag - role - token
  * tag - role - user
  * tag - template
  * tag - plugindir(Inconsistent)
*/

// db.tpl_rel.type_id
const (
	TPL_REL_T_ACL_USER = iota
	TPL_REL_T_ACL_TOKEN
	TPL_REL_T_RULE_TRIGGER
)

type Tpl_rel struct {
	Id      int64
	TplId   int64
	TagId   int64
	SubId   int64
	Creator int64
	TypeId  int64
}

// Access
// return nil *Tag if tag not exist
func (op *Operator) AccessByStr(tokenId int64, tag string, chkExist bool) (t *Tag, err error) {
	// if not found, err will be return <QuerySeter> no row found
	if t, err = op.GetTagByName(tag); chkExist && err != nil {
		return
	}

	if op.IsAdmin() {
		return
	}

	_, err = access(op.O, op.User.Id, tokenId, t.Id)
	return t, err
}

func (op *Operator) Access(tokenId, tagId int64) (err error) {

	if tokenId != SYS_A_TOKEN && op.IsAdmin() {
		return
	}

	_, err = access(op.O, op.User.Id, tokenId, tagId)
	return err
}

// all tags that the user has token
func userHasTokenTag(o orm.Ormer, userId, tokenId int64) (tag_ids []int64, err error) {
	var n int64

	n, err = o.Raw("SELECT distinct b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id",
		TPL_REL_T_ACL_USER, userId, TPL_REL_T_ACL_TOKEN, tokenId).QueryRows(&tag_ids)
	if err != nil || n == 0 {
		err = falcon.EACCES
	}
	return
}

// all tags that the user has token(include child tag)
func userHasTokenTagExpend(o orm.Ormer, userId, tokenId int64) (tag_ids []int64, err error) {
	var n int64

	n, err = o.Raw("SELECT distinct c1.tag_id FROM tag_rel c1 JOIN ( SELECT distinct b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id) c2 on c1.sup_tag_id = c2.user_tag_id",
		TPL_REL_T_ACL_USER, userId, TPL_REL_T_ACL_TOKEN, tokenId).QueryRows(&tag_ids)
	if err != nil || n == 0 {
		err = falcon.EACCES
	}

	return
}

// if user has token
func userHasToken(o orm.Ormer, userId, tokenId int64) (tagId int64, err error) {
	err = o.Raw("SELECT b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id LIMIT 1",
		TPL_REL_T_ACL_USER, userId, TPL_REL_T_ACL_TOKEN, tokenId).QueryRow(&tagId)
	if err != nil {
		err = falcon.EACCES
	}
	return
}

func access(o orm.Ormer, userId, tokenId, tagId int64) (tid int64, err error) {
	// TODO: test
	err = o.Raw("SELECT c2.token_tag_id FROM tag_rel c1 JOIN ( SELECT MAX(b1.token_tag_id) as token_tag_id, b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id GROUP BY b1.user_tag_id) c2 ON c1.sup_tag_id = c2.user_tag_id WHERE c1.tag_id = ?",
		TPL_REL_T_ACL_USER, userId, TPL_REL_T_ACL_TOKEN, tokenId, tagId).QueryRow(&tid)
	if err != nil {
		err = falcon.EACCES
	}

	return
}

func addTplRel(o orm.Ormer, userId, tagId, tplId, subId, typeId int64) (id int64, err error) {
	var res sql.Result

	res, err = o.Raw("insert tpl_rel (tpl_id, tag_id, sub_id, creator, type_id) values (?, ?, ?, ?, ?)", tplId, tagId, subId, userId, typeId).Exec()
	if err != nil {
		return
	}

	id, err = res.LastInsertId()

	DbLog(o, userId, CTL_M_TPL, id, CTL_A_ADD, jsonStr(Tpl_rel{TplId: tplId, TagId: tagId, SubId: subId, Creator: userId, TypeId: typeId}))
	return
}

func delTplRel(o orm.Ormer, userId, tagId, tplId, subId, typeId int64) (int64, error) {
	t := &Tpl_rel{TplId: tplId, TagId: tagId, SubId: subId,
		TypeId: typeId}

	res, err := o.Raw("DELETE FROM `tpl_rel` "+
		"WHERE tpl_id = ? AND tag_id = ? and sub_id = ? and type_id = ?", tplId, tagId, subId, typeId).Exec()
	if err != nil {
		return 0, err
	}

	DbLog(o, userId, CTL_M_TPL, 0, CTL_A_DEL, jsonStr(t))
	return res.RowsAffected()
}
