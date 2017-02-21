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
	TagId    int64
	SupTagId int64
	Name     string
	Child    []*tagNode
}

type RelTagHost struct {
	Id       int64  `json:"id"`
	TagId    int64  `json:"tag_id"`
	HostId   int64  `json:"host_id"`
	TagName  string `json:"tag_name"`
	HostName string `json:"host_name"`
}

type RelTagHosts struct {
	TagId   int64   `json:"tag_id"`
	HostIds []int64 `json:"host_ids"`
}

type RelTagTpl struct {
	TagId int64 `json:"tag_id"`
	TplId int64 `json:"tpl_id"`
}

type RelTagTpls struct {
	TagId  int64   `json:"tag_id"`
	TplIds []int64 `json:"tpl_ids"`
}

type RelTagRoleUser struct {
	TagId  int64 `json:"tag_id"`
	RoleId int64 `json:"role_id"`
	UserId int64 `json:"user_id"`
}

type RelTagRoleToken struct {
	TagId   int64 `json:"tag_id"`
	RoleId  int64 `json:"role_id"`
	TokenId int64 `json:"token_id"`
}

type TagRoleUser struct {
	TagName  string `json:"tag_name"`
	RoleName string `json:"role_name"`
	UserName string `json:"user_name"`
	TagId    int64  `json:"tag_id"`
	RoleId   int64  `json:"role_id"`
	UserId   int64  `json:"user_id"`
}

type TagRoleToken struct {
	TagName   string `json:"tag_name"`
	RoleName  string `json:"role_name"`
	TokenName string `json:"token_name"`
	TagId     int64  `json:"tag_id"`
	RoleId    int64  `json:"role_id"`
	TokenId   int64  `json:"token_id"`
}

/*******************************************************************************
 ************************ tag host *********************************************
 ******************************************************************************/

func tagHostSql(tag_id int64, query string, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "b.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}
	if deep {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tag_id))
	} else {
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tag_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTagHostCnt(tag_id int64,
	query string, deep bool) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	sql, sql_args := tagHostSql(tag_id, query, deep)
	err = op.O.Raw("SELECT count(*) FROM tag_host a left join host b on a.host_id = b.id left join tag c on a.tag_id = c.id "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTagHost(tag_id int64, query string, deep bool,
	limit, offset int) (ret []RelTagHost, err error) {
	sql, sql_args := tagHostSql(tag_id, query, deep)
	sql = "select a.tag_id as tag_id, a.host_id as host_id, b.name as host_name, c.name as tag_name from tag_host a left join host b on a.host_id = b.id left join tag c on a.tag_id = c.id " + sql + " ORDER BY b.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) CreateTagHost(rel RelTagHost) (int64, error) {
	val := fmt.Sprintf("(%d, %d)", rel.TagId, rel.HostId)

	if res, err := op.O.Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + val).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

func (op *Operator) CreateTagHosts(rel RelTagHosts) (int64, error) {
	vs := make([]string, len(rel.HostIds))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", rel.TagId, rel.HostIds[i])
	}

	if res, err := op.O.Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

func (op *Operator) DeleteTagHost(rel Id) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_host` "+
		"WHERE id = ?", rel.Id).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (op *Operator) DeleteTagHosts(rel Ids) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_host` " +
		"WHERE id IN " + array2sql(rel.Ids)).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

/*******************************************************************************
 ************************ tag template *********************************************
 ******************************************************************************/

const (
	tagTplCntSql   = "SELECT count(*) as cnt FROM tag_tpl WHERE tag_id = ?"
	tagTplSql      = " FROM tag_tpl a LEFT JOIN template b ON a.tpl_id = b.id WHERE a.tag_id = ?"
	tagTplQuerySql = tagTplSql + " AND b.name LIKE ?"
)

func tagTplSql1(query string, deep bool, user_id int64) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "a.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}
	if deep {
		sql2 = append(sql2, "a.create_user_id = ?")
		sql3 = append(sql3, user_id)
	}
	if user_id != 0 {
		sql2 = append(sql2, "a.create_user_id = ?")
		sql3 = append(sql3, user_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTagTplCnt(tag_id int64,
	query string, deep, mine bool) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?

	if query == "" {
		err = op.O.Raw(tagTplCntSql, tag_id).QueryRow(&cnt)
	} else {
		err = op.O.Raw("SELECT count(*) as cnt"+
			tagTplQuerySql, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (op *Operator) GetTagTpl(tag_id int64, query string, deep, mine bool,
	limit, offset int) (tpls []Template, err error) {
	if query == "" {
		_, err = op.O.Raw("SELECT b.*"+
			tagTplSql+" LIMIT ? OFFSET ?",
			tag_id, limit, offset).QueryRows(&tpls)
	} else {
		_, err = op.O.Raw("SELECT b.*"+
			tagTplQuerySql+" LIMIT ? OFFSET ?",
			tag_id, "%"+query+"%", limit, offset).QueryRows(&tpls)
	}
	return
}

func (op *Operator) CreateTagTpl(rel RelTagTpl) (int64, error) {
	val := fmt.Sprintf("(%d, %d, %d)", rel.TagId, rel.TplId, op.User.Id)

	if res, err := op.O.Raw("INSERT `tag_tpl` (`tag_id`, " +
		"`tpl_id`, `creator`) VALUES " + val).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

func (op *Operator) CreateTagTpls(rel RelTagTpls) (int64, error) {
	vs := make([]string, len(rel.TplIds))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", rel.TagId, rel.TplIds[i])
	}

	if res, err := op.O.Raw("INSERT `tag_tpl` (`tag_id`, " +
		"`tpl_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

func (op *Operator) DeleteTagTpl(rel RelTagTpl) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_tpl` "+
		"WHERE tag_id = ? and tpl_id = ? ",
		rel.TagId, rel.TplId).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (op *Operator) DeleteTagTpls(rel RelTagTpls) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_tpl` "+
		"WHERE tag_id = ? and tpl_id IN "+array2sql(rel.TplIds),
		rel.TagId).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

/*******************************************************************************
 ************************ tag role user ****************************************
 ******************************************************************************/
func (op *Operator) GetTagRoleUserCnt(global bool, tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if global && query == "" {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_USER).QueryRow(&cnt)
	} else if global && query != "" {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN user b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_USER, "%"+query+"%").QueryRow(&cnt)
	} else if !global && query == "" {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_USER, tag_id).QueryRow(&cnt)
	} else {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN user b ON b.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and b.name like ?", TPL_REL_T_ACL_USER, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (op *Operator) GetTagRoleUser(global bool, tag_id int64, query string,
	limit, offset int) (ret []TagRoleUser, err error) {
	if global && query == "" { // show all
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, limit, offset).QueryRows(&ret)
	} else if global && query != "" { // show global query(user name)
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and d.name like ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else if !global && query == "" { // show tag's user
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, tag_id, limit, offset).QueryRows(&ret)
	} else {
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN user d ON d.id = a.sub_id WHERE a.type_id = ? and  a.tag_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_USER, tag_id, "%"+query+"%", limit, offset).QueryRows(&ret)

	}
	if err != nil {
		return nil, err
	}
	return
}

func (op *Operator) CreateTagRoleUser(rel RelTagRoleUser) (int64, error) {
	var val string
	val = fmt.Sprintf("(%d, %d, %d, %d)", rel.TagId,
		rel.RoleId, rel.UserId, TPL_REL_T_ACL_USER)

	res, err := op.O.Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, " +
		"`sub_id`, `type_id`) VALUES " + val).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (op *Operator) DeleteTagRoleUser(rel RelTagRoleUser) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ? and "+
		"type_id = ?", rel.TagId, rel.RoleId, rel.UserId,
		TPL_REL_T_ACL_USER).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

/*******************************************************************************
 ************************ tag role token ***************************************
 ******************************************************************************/
func (op *Operator) GetTagRoleTokenCnt(global bool, tag_id int64,
	query string) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	if global && query == "" {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ?", TPL_REL_T_ACL_TOKEN).QueryRow(&cnt)
	} else if global && query != "" { // show global query(token name)
		err = op.O.Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN token b ON b.id = a.sub_id WHERE a.type_id = ? and b.name like ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%").QueryRow(&cnt)
	} else if !global && query == "" { // show tag's token
		err = op.O.Raw("SELECT count(*) FROM tpl_rel WHERE type_id = ? and tag_id = ?", TPL_REL_T_ACL_TOKEN, tag_id).QueryRow(&cnt)
	} else {
		err = op.O.Raw("SELECT count(*) FROM tpl_rel a LEFT JOIN token b ON b.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and b.name like ?", TPL_REL_T_ACL_TOKEN, tag_id, "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (op *Operator) GetTagRoleToken(global bool, tag_id int64, query string,
	limit, offset int) (ret []TagRoleToken, err error) {
	if global && query == "" { // show all
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? ORDER BY d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, limit, offset).QueryRows(&ret)
	} else if global && query != "" { // show global query(token name)
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, "%"+query+"%", limit, offset).QueryRows(&ret)
	} else if !global && query == "" { // show tag's token
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and a.tag_id = ? ORDER by d.name, c.name LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, tag_id, limit, offset).QueryRows(&ret)
	} else {
		_, err = op.O.Raw("SELECT b.name as tag_name, c.name as role_name, d.name as token_name, a.tag_id, a.tpl_id as role_id, a.sub_id as token_id FROM tpl_rel a LEFT JOIN tag b ON b.id = a.tag_id LEFT JOIN role c ON c.id = a.tpl_id LEFT JOIN token d ON d.id = a.sub_id WHERE a.type_id = ? and tag_id = ? and d.name like ? ORDER BY d.name, c.name  LIMIT ? OFFSET ?", TPL_REL_T_ACL_TOKEN, tag_id, "%"+query+"%", limit, offset).QueryRows(&ret)
	}
	if err != nil {
		return nil, err
	}
	return
}

func (op *Operator) CreateTagRoleToken(rel RelTagRoleToken) (int64, error) {
	var val string
	val = fmt.Sprintf("(%d, %d, %d, %d)", rel.TagId,
		rel.RoleId, rel.TokenId, TPL_REL_T_ACL_TOKEN)

	res, err := op.O.Raw("INSERT `tpl_rel` (`tag_id`, `tpl_id`, " +
		"`sub_id`, `type_id`) VALUES " + val).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (op *Operator) DeleteTagRoleToken(rel RelTagRoleToken) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tpl_rel` "+
		"WHERE tag_id = ? and tpl_id = ? and sub_id = ?"+
		" and type_id = ?", rel.TagId, rel.RoleId,
		rel.TokenId, TPL_REL_T_ACL_TOKEN).Exec()
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (op *Operator) GetTagTags(tag_id int64) (nodes []zTreeNode, err error) {
	if tag_id == 0 {
		return []zTreeNode{{Id: 1, Name: "/"}}, nil
	}
	_, err = op.O.Raw("SELECT `tag_id` AS `id`, `b`.`name` "+
		"FROM `tag_rel` `a` LEFT JOIN `tag` `b` "+
		"ON `a`.`tag_id` = `b`.`id` WHERE `a`.`sup_tag_id` = ? "+
		"AND `a`.`offset` = 1", tag_id).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return
}

func (op *Operator) GetTreeNodes(id int64) (nodes []TreeNode, err error) {
	if id == 0 {
		return []TreeNode{{Id: 1, Name: "/"}}, nil
	}
	_, err = op.O.Raw("SELECT `tag_id` AS `id`, `b`.`name` "+
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
		Id:    n.TagId,
		Name:  n.Name,
		Label: n.Name[strings.LastIndexAny(n.Name, ",")+1:],
	}
	for _, v := range n.Child {
		tree.Child = append(tree.Child, pruneTagTree(nodes, v.TagId))
	}
	return tree
}

func (op *Operator) GetTree() (tree *TreeNode, err error) {
	var ns []tagNode
	nodes := make(map[int64]*tagNode)
	nodes[1] = &tagNode{
		Name:  "/",
		TagId: 1,
	}

	_, err = op.O.Raw("SELECT `a`.`tag_id`, `a`.`sup_tag_id`, `b`.`name` " +
		"FROM `tag_rel` `a` LEFT JOIN `tag` `b` " +
		"ON `a`.`tag_id` = `b`.`id` WHERE " +
		"`a`.`offset` = 1 ORDER BY `tag_id`").QueryRows(&ns)
	if err != nil {
		return
	}
	for idx, _ := range ns {
		n := &ns[idx]
		nodes[n.TagId] = n
		nodes[n.SupTagId].Child = append(nodes[n.SupTagId].Child, n)
	}
	if t := pruneTagTree(nodes, 1); t != nil {
		return t, nil
	}
	return
}

/**************************
 *     tag rule trigger
 *************************/

func (op *Operator) BindAclUser(tag_id, role_id, user_id int64) error {
	return addTplRel(op.O, op.User.Id, tag_id, role_id,
		user_id, TPL_REL_T_ACL_USER)
}

func (op *Operator) BindAclToken(tag_id, role_id, token_id int64) (err error) {
	return addTplRel(op.O, op.User.Id, tag_id, role_id,
		token_id, TPL_REL_T_ACL_TOKEN)
}

/*
func (op *Operator) BindRuleHost(tag_id, rule_id, host_id int64) error {
	return addTplRel(op.O, op.User.Id, tag_id, rule_id, host_id, TPL_REL_T_RULE_HOST)
}
*/

func (op *Operator) BindRuleTrigger(tag_id, rule_id, trigger_id int64) error {
	return addTplRel(op.O, op.User.Id, tag_id, rule_id,
		trigger_id, TPL_REL_T_RULE_TRIGGER)
}

// Access
// return nil *Tag if tag not exist
func (op *Operator) Access(token, tag string, chkExist bool) (t *Tag, err error) {
	var tk *Token

	// if not found, err will be return <QuerySeter> no row found
	if t, err = op.GetTagByName(tag); chkExist && err != nil {
		return
	}

	if op.IsAdmin() {
		return
	}

	if tk, err = op.GetTokenByName(token); err != nil {
		return
	}

	_, err = access(op.O, op.User.Id, tk.Id, t.Id)
	return t, err
}

func access(o orm.Ormer, user_id, token_id, tag_id int64) (tid int64, err error) {
	// TODO: test
	if tag_id == 0 {
		err = o.Raw(`
    SELECT MIN(b1.token_tag_id) as token_tag_id
    FROM (SELECT a1.tag_id AS user_tag_id,
                a2.tag_id AS token_tag_id,
                a1.tpl_id AS role_id,
                a1.sub_id AS user_id,
                a2.sub_id AS token_id
          FROM tpl_rel a1
          JOIN tpl_rel a2 
          ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?
	      AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1
    JOIN tag_rel b2
    ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id
    GROUP BY b1.user_tag_id
`,
			TPL_REL_T_ACL_USER, user_id, TPL_REL_T_ACL_TOKEN, token_id).QueryRow(&tid)
	} else {
		err = o.Raw(`
SELECT c2.token_tag_id
FROM tag_rel c1
JOIN (
    SELECT MAX(b1.token_tag_id) as token_tag_id, b1.user_tag_id
    FROM (SELECT a1.tag_id AS user_tag_id,
                a2.tag_id AS token_tag_id,
                a1.tpl_id AS role_id,
                a1.sub_id AS user_id,
                a2.sub_id AS token_id
          FROM tpl_rel a1
          JOIN tpl_rel a2 
          ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?
	      AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1
    JOIN tag_rel b2
    ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id
    GROUP BY b1.user_tag_id) c2
ON  c1.sup_tag_id = c2.user_tag_id 
WHERE c1.tag_id = ?
`,
			TPL_REL_T_ACL_USER, user_id, TPL_REL_T_ACL_TOKEN, token_id, tag_id).QueryRow(&tid)
	}
	if err != nil {
		err = EACCES
	}

	return
}

func addTplRel(o orm.Ormer, user_id, tag_id, tpl_id, sub_id, type_id int64) error {
	t := &Tpl_rel{TplId: tpl_id, TagId: tag_id, SubId: sub_id,
		Creator: user_id, TypeId: type_id}

	id, err := o.Insert(t)
	if err != nil {
		return err
	}
	DbLog(o, user_id, CTL_M_TPL, id, CTL_A_ADD, jsonStr(t))
	return nil
}
