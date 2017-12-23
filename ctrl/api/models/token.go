/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Cname string `json:"cname"`
	Note  string `json:"note"`
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
	err = op.SqlRow(ret, "select id, name, cname, note from token where id = ?", id)
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
	err = op.SqlRow(ret, "select id, name, cname, note from token where name = ?", token)
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
	sql = sqlLimit("select id, name, cname, note from token "+sql+" ORDER BY id", limit, offset)
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

/*******************************************************************************
 ************************ tag role token ***************************************
 ******************************************************************************/

type TagRoleTokenApi struct {
	TagId   int64 `json:"tag_id"`
	RoleId  int64 `json:"role_id"`
	TokenId int64 `json:"token_id"`
}

type TagRolesTokensApiAdd struct {
	TagId    int64   `json:"tag_id"`
	RoleIds  []int64 `json:"role_ids"`
	TokenIds []int64 `json:"token_ids"`
}

type TagRolesTokensApiDel struct {
	TagId     int64 `json:"tag_id"`
	RoleToken []struct {
		RoleId  int64 `json:"role_id"`
		TokenId int64 `json:"token_id"`
	} `json:"role_token"`
}

type TagRoleTokenApiGet struct {
	TagName   string `json:"tag_name"`
	RoleName  string `json:"role_name"`
	TokenName string `json:"token_name"`
	TagId     int64  `json:"tag_id"`
	RoleId    int64  `json:"role_id"`
	TokenId   int64  `json:"token_id"`
}

func tagRoleTokenSql(tagId int64, query string, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	sql2 = append(sql2, "a.type_id = ?")
	sql3 = append(sql3, TPL_REL_T_ACL_TOKEN)

	if query != "" {
		sql2 = append(sql2, "tk.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if deep {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tagId))
	} else {
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tagId)
	}

	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTagRoleTokenCnt(tagId int64,
	query string, deep bool) (cnt int64, err error) {
	sql, sql_args := tagRoleTokenSql(tagId, query, deep)
	err = op.O.Raw("SELECT count(*) FROM tpl_rel a JOIN tag t ON t.id = a.tag_id JOIN role r ON r.id = a.tpl_id JOIN token tk ON tk.id = a.sub_id "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTagRoleToken(tagId int64, query string, deep bool,
	limit, offset int) (ret []TagRoleTokenApiGet, err error) {
	sql, sql_args := tagRoleTokenSql(tagId, query, deep)
	sql = "SELECT t.name as tag_name, r.name as role_name, tk.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a JOIN tag t ON t.id = a.tag_id JOIN role r ON r.id = a.tpl_id JOIN token tk ON tk.id = a.sub_id  " + sql + " ORDER BY tk.name, r.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) CreateTagRoleToken(rel *TagRoleTokenApi) (int64, error) {
	return addTplRel(op.O, op.User.Id, rel.TagId, rel.RoleId,
		rel.TokenId, TPL_REL_T_ACL_TOKEN)
}

func (op *Operator) DeleteTagRoleToken(rel *TagRoleTokenApi) (int64, error) {
	return delTplRel(op.O, op.User.Id, rel.TagId, rel.RoleId,
		rel.TokenId, TPL_REL_T_ACL_TOKEN)
}

func (op *Operator) GetTagTags(tagId int64) (nodes []zTreeNode, err error) {
	if tagId == 0 {
		return []zTreeNode{{Id: 1, Name: "/"}}, nil
	}
	_, err = op.O.Raw("SELECT `tag_id` AS `id`, `b`.`name` FROM `tag_rel` `a` LEFT JOIN `tag` `b` ON `a`.`tag_id` = `b`.`id` WHERE `a`.`sup_tag_id` = ? AND `a`.`offset` = 1 and `b`.`type` = 0", tagId).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return
}
