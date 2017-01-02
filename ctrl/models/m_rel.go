/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "github.com/astaxie/beego/orm"

type Tpl_rel struct {
	Id      int64
	Tpl_id  int64
	Tag_id  int64
	Sub_id  int64
	Creator int64
	Type_id int64
}

func addTplRel(user_id, tpl_id, tag_id, sub_id, type_id int64) error {
	t := &Tpl_rel{Tpl_id: tpl_id, Tag_id: tag_id, Sub_id: sub_id,
		Creator: user_id, Type_id: type_id}

	id, err := orm.NewOrm().Insert(t)
	if err != nil {
		return err
	}

	DbLog(user_id, CTL_M_TPL, id, CTL_A_ADD, jsonStr(t))

	return nil
}

func (u *User) BindAclUser(role_id, tag_id, user_id int64) error {
	return addTplRel(u.Id, role_id, tag_id, user_id, TPL_REL_T_ACL_USER)
}

func (u *User) BindAclToken(role_id, tag_id, token_id int64) (err error) {
	return addTplRel(u.Id, role_id, tag_id, token_id, TPL_REL_T_ACL_TOKEN)
}

func (u *User) BindRuleHost(rule_id, tag_id, host_id int64) error {
	return addTplRel(u.Id, rule_id, tag_id, host_id, TPL_REL_T_RULE_HOST)
}

func (u *User) BindRuleTrigger(rule_id, tag_id, trigger_id int64) (err error) {
	return addTplRel(u.Id, rule_id, tag_id, trigger_id, TPL_REL_T_RULE_TRIGGER)
}

// Access
// return nil *Tag if tag not exist
func (u *User) Access(token, tag string, chkExist bool) (t *Tag, err error) {
	var s *Token
	var tag_id int64

	// if not found, err will be return <QuerySeter> no row found
	if t, err = u.GetTagByName(tag); chkExist && err != nil {
		return
	}

	if u.IsAdmin() {
		return
	}

	if s, err = u.GetTokenByName(token); err != nil {
		return
	}

	// TODO: test
	err = orm.NewOrm().Raw(`
SELECT sup_tag_id
FROM tag_rel
WHERE sup_tag_id IN (
    SELECT user_token.user_tag_id
    FROM (SELECT a.tag_id AS user_tag_id,
                b.tag_id AS token_tag_id,
                a.role_id AS role_id,
                a.user_id AS user_id,
                b.token_id AS token_id
          FROM tag_role_user a
          JOIN tag_role_token b
          ON a.user_id = ? AND b.token_id = ? AND a.role_id = b.role_id) user_token
   JOIN tag_rel 
   ON user_token.user_tag_id = tag_rel.tag_id AND user_token.token_tag_id = tag_rel.sup_tag_id
   GROUP BY user_token.user_tag_id)
AND tag_id = ?
`,
		u.Id, s.Id, t.Id).QueryRow(&tag_id)

	return
}

func access(user_id, token_id, tag_id int64) (tid int64, err error) {
	// TODO: test
	err = orm.NewOrm().Raw(`
SELECT b2.token_tag_id
FROM tag_rel a2
JOIN (
    SELECT a1.token_tag_id, a1.user_tag_id
    FROM (SELECT a0.tag_id AS user_tag_id,
                b0.tag_id AS token_tag_id,
                a0.tpl_id AS role_id,
                a0.sub_id AS user_id,
                b0.sub_id AS token_id
          FROM tpl_rel a0
          JOIN tpl_rel b0
          ON a0.type_id = ? AND a0.sub_id = ? AND b0.type_id = ?
	      AND b0.sub_id = ? AND a0.tpl_id = b0.tpl_id) a1
    JOIN tag_rel b1
    ON a1.user_tag_id = b1.tag_id AND a1.token_tag_id = b1.sup_tag_id
    GROUP BY a1.user_tag_id 
    HAVING a1.token_tag_id = MAX(a1.token_tag_id)) b2
ON  a2.sup_tag_id = b2.user_tag_id 
WHERE a2.tag_id = ?
`,
		TPL_REL_T_ACL_USER, user_id, TPL_REL_T_ACL_TOKEN, token_id, tag_id).QueryRow(&tid)
	if err != nil {
		err = EACCES
	}

	return
}
