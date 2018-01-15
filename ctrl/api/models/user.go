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

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id        int64  `json:"id"`
	Muid      int64  `json:"muid"` /* master uid */
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Cname     string `json:"cname"`
	Mname     string `json:"mname"` /* master user name */
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Qq        string `json:"qq"`
	Extra     string `json:"extra"`
	Avatarurl string `json:"avatarurl"`
	Disabled  int    `json:"disabled"`
}

type UserApiAdd struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Cname     string `json:"cname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Qq        string `json:"qq"`
	Extra     string `json:"extra"`
	Avatarurl string `json:"avatarurl"`
}

type UserApiUpdate struct {
	Id        int64  `json:"id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Cname     string `json:"cname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Qq        string `json:"qq"`
	Extra     string `json:"extra"`
	Avatarurl string `json:"avatarurl"`
	Disabled  int    `json:"disabled"`
}

type UserProfileUpdate struct {
	Cname string `json:"cname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Qq    string `json:"qq"`
	Extra string `json:"extra"`
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

func UserTokens(uid int64, username string, o orm.Ormer) (token int) {
	var (
		tids []int64
	)

	if admin[username] {
		return SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_R_TOKEN
	}

	_, err := o.Raw("SELECT b1.token_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id in (?, ?, ?) AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id GROUP BY b1.token_id",
		TPL_REL_T_ACL_USER, uid, TPL_REL_T_ACL_TOKEN,
		SYS_R_TOKEN, SYS_O_TOKEN,
		SYS_A_TOKEN).QueryRows(&tids)
	if err != nil {
		return 0
	}
	for _, tid := range tids {
		switch tid {
		case SYS_R_TOKEN:
			token |= SYS_F_R_TOKEN
		case SYS_O_TOKEN:
			token |= SYS_F_O_TOKEN
		case SYS_A_TOKEN:
			token |= SYS_F_A_TOKEN
		}
	}

	if uid < 3 {
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

func (op *Operator) UserTokens() (token int) {
	if op.User == nil {
		return 0
	}

	return UserTokens(op.User.Id, op.User.Name, op.O)
}

func (op *Operator) CreateUser(user *UserApiAdd) (id int64, err error) {
	id, err = op.SqlInsert("insert user (uuid, name, cname, email, phone, qq, extra, avatarurl) values (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Uuid, user.Name, user.Cname, user.Email,
		user.Phone, user.Qq, user.Extra, user.Avatarurl)
	if err != nil {
		return
	}

	op.log(CTL_M_USER, id, CTL_A_ADD, jsonStr(user))
	return
}

// just called from profileFilter()
func GetUser(id int64, o orm.Ormer) (ret *User, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_USER].get(id).(*User); ok {
		return ret, nil
	}

	ret = &User{}
	err = o.Raw("select a.id, a.muid, a.uuid, a.name, a.cname, a.email, a.phone, a.qq, a.disabled, a.extra, a.avatarurl, b.name as mname from user a left join user b on a.muid = b.id where a.id = ?", id).QueryRow(ret)
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
	err = op.SqlRow(ret, "select a.id, a.muid, a.uuid, a.name, a.cname, a.email, a.phone, a.qq, a.disabled, a.extra, a.avatarurl, b.name as mname from user a left join user b on a.muid = b.id where a.uuid = ?", uuid)
	return ret, err
}

func sqlUser(query string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "user.name like ? or user.email like ?")
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
	sql = sqlLimit("select user.id, user.muid, user.uuid, user.name, user.cname, user.email, user.phone, user.qq, user.disabled, user.extra, user.avatarurl, user.extra, b.name as mname from user left join user b on user.muid = b.id "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)

	return
}

func (op *Operator) GetBindedUsers(id int64) (ret []*User, err error) {
	_, err = op.O.Raw("select id, muid, uuid, name, cname, email, phone, qq, disabled, extra, avatarurl, extra from user where muid = ?", id).QueryRows(&ret)
	return
}

func (op *Operator) UpdateUser(user *User) (ret *User, err error) {
	_, err = op.SqlExec("update user set name = ?, cname = ?, email = ?, phone = ?, qq = ?, disabled = ?, extra = ?, avatarurl = ? where id = ?",
		user.Name, user.Cname, user.Email, user.Phone, user.Qq,
		user.Disabled, user.Extra, user.Avatarurl, user.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_USER].del(user.Id)
	if ret, err = op.GetUser(user.Id); err != nil {
		return
	}

	op.log(CTL_M_USER, user.Id, CTL_A_SET, "")
	return
}

func (op *Operator) UnBindUser(id int64) (err error) {
	_, err = op.SqlExec("update user set muid = 0 where muid = ?", id)
	if err != nil {
		return
	}
	op.log(CTL_M_USER, id, CTL_A_SET, fmt.Sprintf("unbind %d ", id))
	return
}

func (op *Operator) BindUser(src, dst int64) (err error) {
	_, err = op.SqlExec("update user set muid = ? where id = ?", dst, src)
	if err != nil {
		return
	}
	moduleCache[CTL_M_USER].del(src)
	op.log(CTL_M_USER, src, CTL_A_SET, fmt.Sprintf("bind %d to %d", src, dst))
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
	op.log(CTL_M_USER, id, CTL_A_DEL, "")

	return nil
}

/*******************************************************************************
 ************************ tag role user ****************************************
 ******************************************************************************/

type TagRoleUserApi struct {
	TagId  int64 `json:"tag_id"`
	RoleId int64 `json:"role_id"`
	UserId int64 `json:"user_id"`
}

type TagRolesUsersApiAdd struct {
	TagId   int64   `json:"tag_id"`
	RoleIds []int64 `json:"role_ids"`
	UserIds []int64 `json:"user_ids"`
}

type TagRolesUsersApiDel struct {
	TagId    int64 `json:"tag_id"`
	RoleUser []struct {
		RoleId int64 `json:"role_id"`
		UserId int64 `json:"user_id"`
	} `json:"role_user"`
}

type TagRoleUserApiGet struct {
	TagName  string `json:"tag_name"`
	RoleName string `json:"role_name"`
	UserName string `json:"user_name"`
	TagId    int64  `json:"tag_id"`
	RoleId   int64  `json:"role_id"`
	UserId   int64  `json:"user_id"`
}

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

func (op *Operator) CreateTagRoleUser(rel *TagRoleUserApi) (int64, error) {
	return op.addTplRel(rel.TagId, rel.RoleId, rel.UserId, TPL_REL_T_ACL_USER)
}

func (op *Operator) DeleteTagRoleUser(rel *TagRoleUserApi) (int64, error) {
	return op.delTplRel(rel.TagId, rel.RoleId, rel.UserId, TPL_REL_T_ACL_USER)
}
