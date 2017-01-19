/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Tpl_rel struct {
	Id      int64
	Tpl_id  int64
	Tag_id  int64
	Sub_id  int64
	Creator int64
	Type_id int64
}

type zTreeNode struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type RelTagHost struct {
	Tag_id  int64 `json:"tag_id"`
	Host_id int64 `json:"host_id"`
}

type RelTagRoleUser struct {
	Tag_id  int64 `json:"tag_id"`
	Role_id int64 `json:"role_id"`
	User_id int64 `json:"user_id"`
}

type RelTagRoleToken struct {
	Tag_id   int64 `json:"tag_id"`
	Role_id  int64 `json:"role_id"`
	Token_id int64 `json:"token_id"`
}

type RelTagHosts struct {
	Tag_id   int64   `json:"tag_id"`
	Host_ids []int64 `json:"host_ids"`
}

type RelTagRoleUsers struct {
	Tag_id   int64   `json:"tag_id"`
	Role_id  int64   `json:"role_id"`
	User_ids []int64 `json:"user_ids"`
}

type RelTagRoleTokens struct {
	Tag_id    int64   `json:"tag_id"`
	Role_id   int64   `json:"role_id"`
	Token_ids []int64 `json:"token_ids"`
}

type TagRoleUser struct {
	TagName  string
	RoleName string
	UserName string
	TagId    int64
	RoleId   int64
	UserId   int64
}

type TagRoleToken struct {
	TagName   string
	RoleName  string
	TokenName string
	TagId     int64
	RoleId    int64
	TokenId   int64
}

func (u *User) CreateTagHosts(rel RelTagHosts) (int64, error) {
	vs := make([]string, len(rel.Host_ids))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", rel.Tag_id, rel.Host_ids[i])
	}

	if res, err := orm.NewOrm().Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}

}

func (u *User) CreateTagRoleUsers(rel RelTagRoleUsers) (int64, error) {
	vs := make([]string, len(rel.User_ids))
	for i := 0; i < len(rel.User_ids); i++ {
		vs[i] = fmt.Sprintf("(%d, %d, %d, %d)", rel.Tag_id,
			rel.Role_id, rel.User_ids[i], TPL_REL_T_ACL_USER)
	}
	beego.Debug(vs)

	res, err := orm.NewOrm().Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, `sub_id`, `type_id`)" +
		"VALUES " + strings.Join(vs, ", ")).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (u *User) CreateTagRoleTokens(rel RelTagRoleTokens) (int64, error) {
	vs := make([]string, len(rel.Token_ids))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d, %d, %d)", rel.Tag_id,
			rel.Role_id, rel.Token_ids[i], TPL_REL_T_ACL_TOKEN)
	}

	res, err := orm.NewOrm().Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, `sub_id`, `type_id`)" +
		"VALUES " + strings.Join(vs, ", ")).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (u *User) CreateTagHost(rel RelTagHost) error {
	_, err := orm.NewOrm().Raw("INSERT `tag_host` (`tag_id`, `host_id`)"+
		"VALUES (?, ?)", rel.Tag_id, rel.Host_id).Exec()
	return err

}

func (u *User) DeleteTagHost(rel RelTagHost) error {
	_, err := orm.NewOrm().Raw("DELETE FROM `tag_host` "+
		"WHERE tag_id = ? and host_id = ?", rel.Tag_id,
		rel.Host_id).Exec()
	return err
}

func (u *User) DeleteTagRoleUser(rel RelTagRoleUser) error {
	_, err := orm.NewOrm().Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ?", rel.Tag_id,
		rel.Role_id, rel.User_id).Exec()
	return err
}

func (u *User) DeleteTagRoleToken(rel RelTagRoleToken) error {
	_, err := orm.NewOrm().Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ?", rel.Tag_id,
		rel.Role_id, rel.Token_id).Exec()
	return err
}

func (u *User) GetTagTags(tag_id int64) (nodes []zTreeNode, err error) {
	if tag_id == 0 {
		return []zTreeNode{{Id: 1, Name: "/"}}, nil
	}
	_, err = orm.NewOrm().Raw("SELECT `tag_id` AS `id`, `b`.`name` "+
		"FROM `tag_rel` `a` LEFT JOIN `tag` `b` "+
		"ON `a`.`tag_id` = `b`.`id` WHERE `a`.`sup_tag_id` = ? "+
		"AND `a`.`offset` = 1", tag_id).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return
}

/**************************
 *       tag host
 *************************/

func (u *User) QueryTagHostCnt(tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if query == "" {
		err = orm.NewOrm().Raw("SELECT count(*) as cnt "+
			"FROM `tag_host` `a` LEFT JOIN `host` `b` "+
			"ON `a`.`host_id` = `b`.`id` WHERE `a`.`tag_id` = ?",
			tag_id).QueryRow(&cnt)
	} else {
		err = orm.NewOrm().Raw("SELECT count(*) as cnt "+
			"FROM `tag_host` `a` LEFT JOIN `host` `b` "+
			"ON `a`.`host_id` = `b`.`id` "+
			"WHERE `a`.`tag_id` = ? and `b`.`name` like ? ",
			tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagHost(tag_id int64, query string,
	limit, offset int) (hosts []Host, err error) {
	if query == "" {
		_, err = orm.NewOrm().Raw("SELECT `b`.* "+
			"FROM `tag_host` `a` LEFT JOIN `host` `b` "+
			"ON `a`.`host_id` = `b`.`id` WHERE `a`.`tag_id` = ? "+
			"LIMIT ? OFFSET ?", tag_id, limit,
			offset).QueryRows(&hosts)
	} else {
		_, err = orm.NewOrm().Raw("SELECT `b`.* "+
			"FROM `tag_host` `a` LEFT JOIN `host` `b` "+
			"ON `a`.`host_id` = `b`.`id` WHERE `a`.`tag_id` = ? "+
			"AND `b`.`name` LIKE ? LIMIT ? OFFSET ?",
			tag_id, "%"+query+"%", limit, offset).QueryRows(&hosts)
	}
	return
}

/**************************
 *     tag role user
 *************************/
func (u *User) QueryTagRoleUserCnt(tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if query == "" && tag_id == 0 {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_USER).QueryRow(&cnt)
	} else if query != "" { // show global query(user name)
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN user b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_USER, "%"+query+"%").QueryRow(&cnt)
	} else { // show tag's user
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_USER, tag_id).QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagRoleUser(tag_id int64, query string,
	limit, offset int) (ret []TagRoleUser, err error) {
	if query == "" && tag_id == 0 { // show all
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, limit, offset).QueryRows(&ret)
	} else if query != "" { // show global query(user name)
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and d.name, c.name like ? ORDER BY d.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else { // show tag's user
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, tag_id, limit, offset).QueryRows(&ret)
	}
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	return
}

/**************************
 *     tag role token
 *************************/
func (u *User) QueryTagRoleTokenCnt(tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if query == "" && tag_id == 0 {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_TOKEN).QueryRow(&cnt)
	} else if query != "" { // show global query(token name)
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN token b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%").QueryRow(&cnt)
	} else { // show tag's token
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_TOKEN, tag_id).QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagRoleToken(tag_id int64, query string,
	limit, offset int) (ret []TagRoleToken, err error) {
	if query == "" && tag_id == 0 { // show all
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, limit, offset).QueryRows(&ret)
	} else if query != "" { // show global query(token name)
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else { // show tag's token
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, tag_id, limit, offset).QueryRows(&ret)
	}
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	return
}

/**************************
 *     tag rule trigger
 *************************/

func (u *User) BindAclUser(tag_id, role_id, user_id int64) error {
	return addTplRel(u.Id, tag_id, role_id,
		user_id, TPL_REL_T_ACL_USER)
}

func (u *User) BindAclToken(tag_id, role_id, token_id int64) (err error) {
	return addTplRel(u.Id, tag_id, role_id,
		token_id, TPL_REL_T_ACL_TOKEN)
}

/*
func (u *User) BindRuleHost(tag_id, rule_id, host_id int64) error {
	return addTplRel(u.Id, tag_id, rule_id, host_id, TPL_REL_T_RULE_HOST)
}
*/

func (u *User) BindRuleTrigger(tag_id, rule_id, trigger_id int64) error {
	return addTplRel(u.Id, tag_id, rule_id,
		trigger_id, TPL_REL_T_RULE_TRIGGER)
}

// Access
// return nil *Tag if tag not exist
func (u *User) Access(token, tag string, chkExist bool) (t *Tag, err error) {
	var tk *Token

	// if not found, err will be return <QuerySeter> no row found
	if t, err = u.GetTagByName(tag); chkExist && err != nil {
		return
	}

	if u.IsAdmin() {
		return
	}

	if tk, err = u.GetTokenByName(token); err != nil {
		return
	}

	_, err = access(u.Id, tk.Id, t.Id)
	return t, err
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

func addTplRel(user_id, tag_id, tpl_id, sub_id, type_id int64) error {
	t := &Tpl_rel{Tpl_id: tpl_id, Tag_id: tag_id, Sub_id: sub_id,
		Creator: user_id, Type_id: type_id}

	id, err := orm.NewOrm().Insert(t)
	if err != nil {
		return err
	}
	DbLog(user_id, CTL_M_TPL, id, CTL_A_ADD, jsonStr(t))
	return nil
}
