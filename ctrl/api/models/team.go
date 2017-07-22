/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/yubo/falcon"
)

type TeamUi struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Note       string    `json:"note"`
	Creator    string    `json:"creator"`
	CreateTime time.Time `json:"ctime"`
}

type Team struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Note       string    `json:"note"`
	Creator    int64     `json:"-"`
	CreateTime time.Time `json:"ctime"`
}

type TeamMemberIds struct {
	Uids []int64 `json:"uids"`
}

type TeamMembers struct {
	Users []User `json:"users"`
}

func (op *Operator) AddTeam(t *Team) (id int64, err error) {

	id, err = op.SqlInsert("insert team (name, note, creator) values (?, ?, ?)", t.Name, t.Note, t.Creator)
	if err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_ADD, jsonStr(t))
	return
}

func (op *Operator) GetTeam(id int64) (ret *Team, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_TEAM].get(id).(*Team); ok {
		return
	}

	ret = &Team{}
	err = op.SqlRow(ret, "select id, name, note, creator, create_time from team where id = ?", id)
	if err == nil {
		moduleCache[CTL_M_TEAM].set(id, ret)
	}
	return
}

func (op *Operator) GetMember(id int64, name string) (*TeamMembers, error) {
	var m TeamMembers
	var err error

	if id != 0 {
		_, err = op.O.Raw("SELECT u.id, u.name, u.email, u.phone "+
			"FROM team_user a JOIN user u "+
			"ON a.user_id = u.id WHERE a.team_id = ? ",
			id).QueryRows(&m.Users)
	} else {
		_, err = op.O.Raw("SELECT u.id, u.name, u.email, u.phone "+
			"FROM team_user a JOIN user u "+
			"ON a.user_id = u.id JOIN team t ON t.id = a.team_id WHERE t.name = ? ",
			name).QueryRows(&m.Users)
	}
	return &m, err
}

func teamSql(query string, user_id int64) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}

	if query != "" {
		sql2 = append(sql2, "a.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}

	if user_id != 0 {
		sql2 = append(sql2, "a.creator = ?")
		sql3 = append(sql3, user_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTeamsCnt(query string, user_id int64) (cnt int64, err error) {
	sql, sql_args := teamSql(query, user_id)
	err = op.O.Raw("SELECT count(*) FROM team a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTeams(query string, user_id int64, limit, offset int) (ret []*TeamUi, err error) {
	sql, sql_args := teamSql(query, user_id)
	sql = "SELECT a.id as id, a.name as name, a.note as note, a.create_time as create_time, u.name as creator from team a LEFT JOIN user u ON u.id = a.creator " + sql + " ORDER BY a.name LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateTeam(id int64, team *Team) (ret *Team, err error) {

	_, err = op.SqlExec("update team set name = ?, note = ? where id = ?", team.Name, team.Note, id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_TEAM].del(id)
	if ret, err = op.GetTeam(id); err != nil {
		return nil, falcon.ErrNoExits
	}

	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(team))
	return ret, err
}

func (op *Operator) UpdateMember(id int64, _m *TeamMemberIds) (m *TeamMemberIds, err error) {
	var tm *TeamMembers

	if tm, err = op.GetMember(id, ""); err != nil {
		return nil, falcon.ErrNoExits
	}

	m = &TeamMemberIds{}
	m.Uids = make([]int64, len(tm.Users))
	for i, v := range tm.Users {
		m.Uids[i] = v.Id
	}

	add, del := MdiffInt(m.Uids, _m.Uids)
	if len(add) > 0 {
		vs := make([]string, len(add))
		for i := 0; i < len(vs); i++ {
			vs[i] = fmt.Sprintf("(%d, %d)", id, add[i])
		}
		if _, err = op.O.Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
			return
		}
	}
	if len(del) > 0 {
		ids := fmt.Sprintf("%d", del[0])
		for i := 0; i < len(del)-1; i++ {
			ids += fmt.Sprintf("%s, %d", ids, del[i])
		}
		if _, err = op.O.Raw("DELETE from `team_user` WHERE team_id = ? and user_id in ("+ids+")", id).Exec(); err != nil {
			return
		}
	}
	m.Uids = _m.Uids
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_m))
	return
}

func (op *Operator) DeleteTeam(id int64) error {
	if _, err := op.SqlExec("delete from team where id = ?", id); err != nil {
		return err
	}
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_DEL, "")

	return nil
}

func (op *Operator) BindTeamUser(team_id, user_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `team_user` (`team_id`, `user_id`) VALUES (?, ?)", team_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (op *Operator) BindTeamUsers(team_id int64, user_ids []int64) (int64, error) {
	vs := make([]string, len(user_ids))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", team_id, user_ids[i])
	}

	if res, err := op.O.Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}
