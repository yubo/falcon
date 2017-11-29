/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

// relation
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

type zTreeNode struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type TreeNode struct {
	TagId    int64  `json:"id"`
	SupTagId int64  `json:"-"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	Read     bool   `json:"read"`
	//Operate  bool        `json:"operate"`
	Child []*TreeNode `json:"children"`
}

type tagNode struct {
	TagId    int64
	SupTagId int64
	Name     string
	Child    []*tagNode
}

type RelTagHostApiAdd struct {
	SrcTagId int64 `json:"src_tag_id"`
	TagId    int64 `json:"tag_id"`
	HostId   int64 `json:"host_id"`
}

type RelTagHostApiDel struct {
	TagId  int64 `json:"tag_id"`
	HostId int64 `json:"host_id"`
}

type RelTagHostsApiDel struct {
	TagId   int64   `json:"tag_id"`
	HostIds []int64 `json:"host_ids"`
}

type RelTagHost struct {
	Id            int64  `json:"id"`
	TagId         int64  `json:"tag_id"`
	HostId        int64  `json:"host_id"`
	TagName       string `json:"tag_name"`
	HostName      string `json:"host_name"`
	Pause         int64  `json:"pause"`
	MaintainBegin int64  `json:"maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"`
}

type RelTagHostsApiAdd struct {
	SrcTagId int64   `json:"src_tag_id"`
	TagId    int64   `json:"tag_id"`
	HostIds  []int64 `json:"host_ids"`
}

type RelTagTpl0 struct {
	TagString string `json:"tag_string"`
	TplName   string `json:"tpl_name"`
}

type RelTagTpl struct {
	TagId int64 `json:"tag_id"`
	TplId int64 `json:"tpl_id"`
}

type RelTagTpls struct {
	TagId  int64   `json:"tag_id"`
	TplIds []int64 `json:"tpl_ids"`
}

type TagTplGet struct {
	Id       int64  `json:"id"`
	TagId    int64  `json:"tag_id"`
	TagName  string `json:"tag_name"`
	TplId    int64  `json:"tpl_id"`
	TplName  string `json:"tpl_name"`
	TplPid   int64  `json:"tpl_pid"`
	TplPname string `json:"tpl_pname"`
	Creator  string `json:"creator"`
}

type TagRoleUserApi struct {
	TagId  int64 `json:"tag_id"`
	RoleId int64 `json:"role_id"`
	UserId int64 `json:"user_id"`
}

type TagRoleUserApiGet struct {
	TagName  string `json:"tag_name"`
	RoleName string `json:"role_name"`
	UserName string `json:"user_name"`
	TagId    int64  `json:"tag_id"`
	RoleId   int64  `json:"role_id"`
	UserId   int64  `json:"user_id"`
}

type TagRoleTokenApi struct {
	TagId   int64 `json:"tag_id"`
	RoleId  int64 `json:"role_id"`
	TokenId int64 `json:"token_id"`
}

type TagRoleTokenApiGet struct {
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

func (op *Operator) ChkTagHostCnt(tagId int64, host_ids []int64) (cnt int64, err error) {
	err = op.O.Raw("SELECT count(*) FROM tag_host where tag_id = ? and host_id IN "+
		array2sql(host_ids), tagId).QueryRow(&cnt)
	return
}

func tagHostSql(tagId int64, query string, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "h.name like ?")
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

func (op *Operator) GetTagHostCnt(tagId int64, query string, deep bool) (cnt int64, err error) {
	var tag *Tag

	//TODO: acl filter just for admin?
	if RunMode&CTL_RUNMODE_MI != 0 {
		tag, err = op.GetTag(tagId)
		if err != nil {
			return
		}
		return miGetTagHostCnt(tag.Name, query, deep)
	} else {
		sql, sql_args := tagHostSql(tagId, query, deep)
		err = op.O.Raw("SELECT count(*) FROM tag_host a left join host h on a.host_id = h.id left join tag t on a.tag_id = t.id "+sql, sql_args...).QueryRow(&cnt)
		return
	}
}

func (op *Operator) GetTagHost(tagId int64, query string, deep bool,
	limit, offset int) (ret []*RelTagHost, err error) {
	var tag *Tag

	if RunMode&CTL_RUNMODE_MI != 0 {
		tag, err = op.GetTag(tagId)
		if err != nil {
			return
		}
		return miGetTagHost(tag.Name, query, deep, limit, offset)
	} else {
		sql, sql_args := tagHostSql(tagId, query, deep)
		sql = "select a.id, a.tag_id as tag_id, a.host_id as host_id, h.name as host_name, h.pause as pause, h.maintain_begin as maintain_begin, h.maintain_end as maintain_end, t.name as tag_name from tag_host a left join host h on a.host_id = h.id left join tag t on a.tag_id = t.id " + sql + " ORDER BY h.name"
		if limit > 0 {
			sql += " LIMIT ? OFFSET ?"
			sql_args = append(sql_args, limit, offset)
		}
		_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
		return
	}
}

func (op *Operator) CreateTagHost(rel *RelTagHostApiAdd) (id int64, err error) {
	id, err = op.SqlInsert("insert tag_host (tag_id, host_id) values (?, ?)",
		rel.TagId, rel.HostId)
	if err == nil {
		DbLog(op.O, op.User.Id, CTL_M_TAG_HOST, id, CTL_A_ADD, jsonStr(rel))
	}
	return
}

func (op *Operator) CreateTagHosts(rel *RelTagHostsApiAdd) (int64, error) {
	vs := make([]string, len(rel.HostIds))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", rel.TagId, rel.HostIds[i])
	}

	res, err := op.O.Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + strings.Join(vs, ", ")).Exec()
	if err != nil {
		return 0, err
	}
	DbLog(op.O, op.User.Id, CTL_M_TAG_HOST, 0, CTL_A_ADD, strings.Join(vs, ", "))
	return res.RowsAffected()
}

func (op *Operator) DeleteTagHost(rel *RelTagHostApiDel) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_host` WHERE tag_id = ? and host_id = ?",
		rel.TagId, rel.HostId).Exec()
	if err != nil {
		return 0, err
	}
	DbLog(op.O, op.User.Id, CTL_M_TAG_HOST, rel.TagId, CTL_A_DEL, jsonStr(rel))
	return res.RowsAffected()
}

func (op *Operator) DeleteTagHosts(rel *RelTagHostsApiDel) (int64, error) {
	if len(rel.HostIds) == 0 {
		return 0, falcon.ErrEmpty
	}
	ids := array2sql(rel.HostIds)
	res, err := op.O.Raw("DELETE FROM `tag_host` "+
		"WHERE tag_id = ? and host_id  IN "+ids, rel.TagId).Exec()
	if err != nil {
		return 0, err
	}
	DbLog(op.O, op.User.Id, CTL_M_TAG_HOST, 0, CTL_A_DEL, ids)
	return res.RowsAffected()
}

/*******************************************************************************
 ************************ tag role user ****************************************
 ******************************************************************************/

func tagRoleUserSql(tagId int64, query string, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	sql2 = append(sql2, "a.type_id = ?")
	sql3 = append(sql3, TPL_REL_T_ACL_USER)

	if query != "" {
		sql2 = append(sql2, "u.name like ?")
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

func (op *Operator) GetTagRoleUserCnt(tagId int64,
	query string, deep bool) (cnt int64, err error) {
	// TODO: acl filter
	// just for admin?
	sql, sql_args := tagRoleUserSql(tagId, query, deep)
	err = op.O.Raw("SELECT count(*) FROM tpl_rel a JOIN tag t ON t.id = a.tag_id JOIN role r ON r.id = a.tpl_id JOIN user u ON u.id = a.sub_id "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTagRoleUser(tagId int64, query string, deep bool,
	limit, offset int) (ret []TagRoleUserApiGet, err error) {
	sql, sql_args := tagRoleUserSql(tagId, query, deep)
	sql = "SELECT t.name as tag_name, r.name as role_name, u.name as user_name, a.tag_id, a.tpl_id as role_id, a.sub_id as user_id FROM tpl_rel a JOIN tag t ON t.id = a.tag_id JOIN role r ON r.id = a.tpl_id JOIN user u ON u.id = a.sub_id " + sql + " ORDER BY u.name, r.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) CreateTagRoleUser(rel TagRoleUserApi) (int64, error) {
	return addTplRel(op.O, op.User.Id, rel.TagId, rel.RoleId,
		rel.UserId, TPL_REL_T_ACL_USER)
}

func (op *Operator) DeleteTagRoleUser(rel TagRoleUserApi) (int64, error) {
	return delTplRel(op.O, op.User.Id, rel.TagId, rel.RoleId,
		rel.UserId, TPL_REL_T_ACL_USER)
}

/*******************************************************************************
 ************************ tag role token ***************************************
 ******************************************************************************/

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

func (op *Operator) CreateTagRoleToken(rel TagRoleTokenApi) (int64, error) {
	return addTplRel(op.O, op.User.Id, rel.TagId, rel.RoleId,
		rel.TokenId, TPL_REL_T_ACL_TOKEN)
}

func (op *Operator) DeleteTagRoleToken(rel TagRoleTokenApi) (int64, error) {
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

func cloneTreeNode(t *TreeNode, depth int) (out *TreeNode) {

	if t == nil {
		return
	}

	out = &TreeNode{
		TagId:    t.TagId,
		SupTagId: t.SupTagId,
		Name:     t.Name,
		Label:    t.Label,
		Read:     t.Read,
	}

	if depth == 0 {
		return
	}

	out.Child = make([]*TreeNode, len(t.Child))

	for i, _ := range t.Child {
		out.Child[i] = cloneTreeNode(t.Child[i], depth-1)
	}

	return
}

func pruneTagTree(nodes map[int64]*TreeNode, idx int64) (tree *TreeNode) {
	n, ok := nodes[idx]
	if !ok {
		return nil
	}

	n.Label = n.Name[strings.LastIndexAny(n.Name, ",")+1:]

	return tree
}

func (op *Operator) GetOpTag(expand bool) ([]int64, error) {
	if expand {
		return userHasTokenTagExpend(op.O, op.User.Id, SYS_O_TOKEN)
	} else {
		return userHasTokenTag(op.O, op.User.Id, SYS_O_TOKEN)
	}
}

func (op *Operator) GetReadTag(expand bool) ([]int64, error) {
	if expand {
		return userHasTokenTagExpend(op.O, op.User.Id, SYS_R_TOKEN)
	} else {
		return userHasTokenTag(op.O, op.User.Id, SYS_R_TOKEN)
	}
}

func (op *Operator) GetTreeNode(tagId int64, depth int, direct bool) (tree *TreeNode) {
	var (
		ids, names []string
		isChild    []bool
	)

	if direct {
		return cloneTreeNode(cacheTree.get(tagId), depth)
	}

	_, err := op.O.Raw("SELECT group_concat(d1.sup_tag_id order by d1.`offset` desc) as ids, group_concat(d3.name order by d1.`offset` desc SEPARATOR ',,') as tags from tag_rel d1 join (SELECT c1.tag_id, c1.sup_tag_id, c1.offset from tag_rel c1 join (SELECT distinct b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id ) c2 on c1.tag_id = c2.user_tag_id WHERE c1.sup_tag_id = ?) d2 left join tag d3 on d1.sup_tag_id = d3.id WHERE d1.tag_id = d2.tag_id AND d1.offset <= d2.offset GROUP BY d1.tag_id ", TPL_REL_T_ACL_USER, op.User.Id, TPL_REL_T_ACL_TOKEN, SYS_R_TOKEN, tagId).QueryRows(&ids, &names)
	if err != nil {
		return nil
	}

	// set child falg to remove from tree
	isChild = make([]bool, len(names))
	for i := 0; i < len(names); i++ {
		for j := 0; j < i; j++ {
			if strings.HasPrefix(names[i], names[j]+",,") {
				isChild[i] = true
			}
		}
	}

	nmap := make(map[int64]*TreeNode)

	for i := 0; i < len(ids); i++ {
		if isChild[i] {
			continue
		}
		name := strings.Split(names[i], ",,")
		_id := strings.Split(ids[i], ",")
		id := make([]int64, len(_id))

		for j := 0; j < len(id); j++ {
			id[j], _ = strconv.ParseInt(_id[j], 10, 0)
		}

		if !(len(id) > 0 && id[0] > 0) {
			return nil
		}

		for j := 0; j < len(id); j++ {
			if _, ok := nmap[id[j]]; ok {
				continue
			}

			n := &TreeNode{
				TagId: id[j],
				Name:  name[j],
				Label: name[j][strings.LastIndexAny(name[j], ",")+1:],
			}

			if j == len(id)-1 {
				if m := depth - len(id); m > 0 {
					n = cloneTreeNode(cacheTree.get(id[len(id)-1]), m)
				} else {
					n.Read = true
				}
			}
			nmap[id[j]] = n

			if j-1 >= 0 {
				nmap[id[j-1]].Child = append(nmap[id[j-1]].Child, n)
			}
		}

	}
	return nmap[tagId]
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
