/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"

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

type TreeNode struct {
	Id    int64       `json:"id"`
	Name  string      `json:"name"`
	Label string      `json:"label"`
	Child []*TreeNode `json:"child"`
}

type tagNode struct {
	Tag_id     int64
	Sup_tag_id int64
	Name       string
	Child      []*tagNode
}

type RelTagHost struct {
	Tag_id  int64 `json:"tag_id"`
	Host_id int64 `json:"host_id"`
}

type RelTagHosts struct {
	Tag_id   int64   `json:"tag_id"`
	Host_ids []int64 `json:"host_ids"`
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

type TagRoleUser struct {
	Tag_name  string `json:"tag_name"`
	Role_name string `json:"role_name"`
	User_name string `json:"user_name"`
	Tag_id    int64  `json:"tag_id"`
	Role_id   int64  `json:"role_id"`
	User_id   int64  `json:"user_id"`
}

type TagRoleToken struct {
	Tag_name   string `json:"tag_name"`
	Role_name  string `json:"role_name"`
	Token_name string `json:"token_name"`
	Tag_id     int64  `json:"tag_id"`
	Role_id    int64  `json:"role_id"`
	Token_id   int64  `json:"token_id"`
}

/*******************************************************************************
 ************************ tag host *********************************************
 ******************************************************************************/

const (
	tagHostCntSql   = "SELECT count(*) as cnt FROM tag_host WHERE tag_id = ?"
	tagHostSql      = " FROM tag_host a LEFT JOIN host b ON a.host_id = b.id WHERE a.tag_id = ?"
	tagHostQuerySql = tagHostSql + " AND b.name LIKE ?"
)

func (u *User) GetTagHostCnt(tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if query == "" {
		err = orm.NewOrm().Raw(tagHostCntSql, tag_id).QueryRow(&cnt)
	} else {
		err = orm.NewOrm().Raw("SELECT count(*) as cnt"+
			tagHostQuerySql, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagHost(tag_id int64, query string,
	limit, offset int) (hosts []Host, err error) {
	if query == "" {
		_, err = orm.NewOrm().Raw("SELECT b.*"+
			tagHostSql+" LIMIT ? OFFSET ?",
			tag_id, limit, offset).QueryRows(&hosts)
	} else {
		_, err = orm.NewOrm().Raw("SELECT b.*"+
			tagHostQuerySql+" LIMIT ? OFFSET ?",
			tag_id, "%"+query+"%", limit, offset).QueryRows(&hosts)
	}
	return
}

func (u *User) CreateTagHost(rel RelTagHost) (int64, error) {
	val := fmt.Sprintf("(%d, %d)", rel.Tag_id, rel.Host_id)

	if res, err := orm.NewOrm().Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + val).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
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

func (u *User) DeleteTagHost(rel RelTagHost) (int64, error) {
	res, err := orm.NewOrm().Raw("DELETE FROM `tag_host` "+
		"WHERE tag_id = ? and host_id = ? ",
		rel.Tag_id, rel.Host_id).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (u *User) DeleteTagHosts(rel RelTagHosts) (int64, error) {
	res, err := orm.NewOrm().Raw("DELETE FROM `tag_host` "+
		"WHERE tag_id = ? and host_id IN "+array2sql(rel.Host_ids),
		rel.Tag_id).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

/*******************************************************************************
 ************************ tag role user ****************************************
 ******************************************************************************/
func (u *User) GetTagRoleUserCnt(global bool, tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if global && query == "" {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_USER).QueryRow(&cnt)
	} else if global && query != "" {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN user b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_USER, "%"+query+"%").QueryRow(&cnt)
	} else if !global && query == "" {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_USER, tag_id).QueryRow(&cnt)
	} else {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN user b ON b.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and b.name like ?", TPL_REL_T_ACL_USER, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagRoleUser(global bool, tag_id int64, query string,
	limit, offset int) (ret []TagRoleUser, err error) {
	if global && query == "" { // show all
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, limit, offset).QueryRows(&ret)
	} else if global && query != "" { // show global query(user name)
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and d.name like ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else if !global && query == "" { // show tag's user
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, tag_id, limit, offset).QueryRows(&ret)
	} else {
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and  a.tag_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, tag_id, "%"+query+"%", limit, offset).QueryRows(&ret)

	}
	if err != nil {
		return nil, err
	}
	return
}

func (u *User) CreateTagRoleUser(rel RelTagRoleUser) (int64, error) {
	var val string
	val = fmt.Sprintf("(%d, %d, %d, %d)", rel.Tag_id,
		rel.Role_id, rel.User_id, TPL_REL_T_ACL_USER)

	res, err := orm.NewOrm().Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, " +
		"`sub_id`, `type_id`) VALUES " + val).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (u *User) DeleteTagRoleUser(rel RelTagRoleUser) (int64, error) {
	res, err := orm.NewOrm().Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ? and "+
		"type_id = ?", rel.Tag_id, rel.Role_id, rel.User_id,
		TPL_REL_T_ACL_USER).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

/*******************************************************************************
 ************************ tag role token ***************************************
 ******************************************************************************/
func (u *User) GetTagRoleTokenCnt(global bool, tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if global && query == "" {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_TOKEN).QueryRow(&cnt)
	} else if global && query != "" { // show global query(token name)
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN token b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%").QueryRow(&cnt)
	} else if !global && query == "" { // show tag's token
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_TOKEN, tag_id).QueryRow(&cnt)
	} else {
		err = orm.NewOrm().Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN token b ON b.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and b.name like ?", TPL_REL_T_ACL_TOKEN, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (u *User) GetTagRoleToken(global bool, tag_id int64, query string,
	limit, offset int) (ret []TagRoleToken, err error) {
	if global && query == "" { // show all
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, limit, offset).QueryRows(&ret)
	} else if global && query != "" { // show global query(token name)
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else if !global && query == "" { // show tag's token
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, tag_id, limit, offset).QueryRows(&ret)
	} else {
		_, err = orm.NewOrm().Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, tag_id, "%"+query+"%", limit, offset).QueryRows(&ret)
	}
	if err != nil {
		return nil, err
	}
	return
}

func (u *User) CreateTagRoleToken(rel RelTagRoleToken) (int64, error) {
	var val string
	val = fmt.Sprintf("(%d, %d, %d, %d)", rel.Tag_id,
		rel.Role_id, rel.Token_id, TPL_REL_T_ACL_TOKEN)

	res, err := orm.NewOrm().Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, " +
		"`sub_id`, `type_id`) VALUES " + val).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (u *User) DeleteTagRoleToken(rel RelTagRoleToken) (int64, error) {
	res, err := orm.NewOrm().Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ?"+
		" and type_id = ?", rel.Tag_id, rel.Role_id,
		rel.Token_id, TPL_REL_T_ACL_TOKEN).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
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

func (u *User) GetTreeNodes(id int64) (nodes []TreeNode, err error) {
	if id == 0 {
		return []TreeNode{{Id: 1, Name: "/"}}, nil
	}
	_, err = orm.NewOrm().Raw("SELECT `tag_id` AS `id`, `b`.`name` "+
		"FROM `tag_rel` `a` LEFT JOIN `tag` `b` "+
		"ON `a`.`tag_id` = `b`.`id` WHERE `a`.`sup_tag_id` = ? "+
		"AND `a`.`offset` = 1", id).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return
}

func pruneTagTree(nodes map[int64]*tagNode, idx int64) (tree *TreeNode) {
	n, ok := nodes[idx]
	if !ok {
		return nil
	}

	tree = &TreeNode{
		Id:    n.Tag_id,
		Name:  n.Name,
		Label: n.Name[strings.LastIndexAny(n.Name, ",")+1:],
	}
	for _, v := range n.Child {
		tree.Child = append(tree.Child, pruneTagTree(nodes, v.Tag_id))
	}
	return tree
}

func (u *User) GetTree() (tree *TreeNode, err error) {
	var ns []tagNode
	nodes := make(map[int64]*tagNode)
	nodes[1] = &tagNode{
		Name:   "/",
		Tag_id: 1,
	}

	_, err = orm.NewOrm().Raw("SELECT `a`.`tag_id`, `a`.`sup_tag_id`, `b`.`name` " +
		"FROM `tag_rel` `a` LEFT JOIN `tag` `b` " +
		"ON `a`.`tag_id` = `b`.`id` WHERE " +
		"`a`.`offset` = 1 ORDER BY `tag_id`").QueryRows(&ns)
	if err != nil {
		return
	}
	for idx, _ := range ns {
		n := &ns[idx]
		nodes[n.Tag_id] = n
		nodes[n.Sup_tag_id].Child = append(nodes[n.Sup_tag_id].Child, n)
	}
	if t := pruneTagTree(nodes, 1); t != nil {
		return t, nil
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
