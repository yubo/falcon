/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Member struct {
	Uids []int64 `json:"uids"`
}

type Team struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"ctime"`
	Creator     int64     `json:"-"`
}

func (u *User) AddTeam(t *Team) (id int64, err error) {
	id, err = orm.NewOrm().Insert(t)
	if err != nil {
		beego.Error(err)
		return
	}

	t.Id = id
	cacheModule[CTL_M_TEAM].set(id, t)
	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_ADD, jsonStr(t))
	return
}

func (u *User) GetTeam(id int64) (*Team, error) {
	if t, ok := cacheModule[CTL_M_TEAM].get(id).(*Team); ok {
		return t, nil
	}
	t := &Team{Id: id}
	err := orm.NewOrm().Read(t, "Id")
	if err == nil {
		cacheModule[CTL_M_TEAM].set(id, t)
	}
	return t, err
}

func (u *User) GetMember(id int64) ([]User, error) {
	users := []User{}
	_, err := orm.NewOrm().Raw("SELECT `b`.`id`, `b`.`name` "+
		"FROM `team_user` `a` LEFT JOIN `user` `b` "+
		"ON `a`.`user_id` = `b`.`id` WHERE `a`.`team_id` = ? ",
		id).QueryRows(&users)
	return users, err
}

func (u *User) QueryTeams(query string, own bool) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Team))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	if own {
		qs = qs.Filter("Creator", u.Id)
	}
	return qs
}

func (u *User) GetTeamsCnt(query string, own bool) (int64, error) {
	return u.QueryTeams(query, own).Count()
}

func (u *User) GetTeams(query string, own bool, limit, offset int) (teams []*Team, err error) {
	_, err = u.QueryTeams(query, own).Limit(limit, offset).All(&teams)
	return
}

func (u *User) UpdateTeam(id int64, _t *Team) (t *Team, err error) {
	if t, err = u.GetTeam(id); err != nil {
		return nil, ErrNoTeam
	}

	if _t.Name != "" {
		t.Name = _t.Name
	}
	if _t.Note != "" {
		t.Note = _t.Note
	}
	_, err = orm.NewOrm().Update(t)
	cacheModule[CTL_M_TEAM].set(id, t)
	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_t))
	return t, err
}

func (u *User) UpdateMember(id int64, _m *Member) (m *Member, err error) {
	var users []User

	if users, err = u.GetMember(id); err != nil {
		return nil, ErrNoTeam
	}

	m = &Member{}
	m.Uids = make([]int64, len(users))
	for i, v := range users {
		m.Uids[i] = v.Id
	}

	add, del := MdiffInt(m.Uids, _m.Uids)
	beego.Debug("add", add)
	beego.Debug("del", del)
	if len(add) > 0 {
		vs := make([]string, len(add))
		for i := 0; i < len(vs); i++ {
			vs[i] = fmt.Sprintf("(%d, %d)", id, add[i])
		}
		if _, err = orm.NewOrm().Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
			return
		}
	}
	if len(del) > 0 {
		ids := fmt.Sprintf("%d", del[0])
		for i := 0; i < len(del)-1; i++ {
			ids += fmt.Sprintf("%s, %d", ids, del[i])
		}
		if _, err = orm.NewOrm().Raw("DELETE from `team_user` WHERE team_id = ? and user_id in ("+ids+")", id).Exec(); err != nil {
			return
		}
	}
	m.Uids = _m.Uids
	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_m))
	return
}

func (u *User) DeleteTeam(id int64) error {
	if n, err := orm.NewOrm().Delete(&Team{Id: id}); err != nil || n == 0 {
		return err
	}
	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_DEL, "")

	return nil
}

func (u *User) BindTeamUser(team_id, user_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `team_user` (`team_id`, `user_id`) VALUES (?, ?)", team_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) BindTeamUsers(team_id int64, user_ids []int64) (int64, error) {
	vs := make([]string, len(user_ids))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", team_id, user_ids[i])
	}

	if res, err := orm.NewOrm().Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}
