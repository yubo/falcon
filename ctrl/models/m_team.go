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

type TeamUsers struct {
	Team     Team    `json:"team"`
	User_ids []int64 `json:"user_ids"`
	Users    []User  `json:"-"`
}

type Team struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Note        string    `json:"note"`
	Create_time time.Time `json:"-"`
}

func (u *User) AddTeamUsers(tu *TeamUsers) (id int64, err error) {
	id, err = orm.NewOrm().Insert(&tu.Team)
	if err != nil {
		return
	}

	_, err = u.BindTeamUsers(id, tu.User_ids)
	if err != nil {
		orm.NewOrm().Delete(&Team{Id: id})
		return
	}
	tu.Team.Id = id
	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_ADD, jsonStr(tu))
	return
}

func (u *User) GetTeamUsers(id int64) (tu *TeamUsers, err error) {
	t := &Team{Id: id}
	err = orm.NewOrm().Read(t, "Id")
	if err != nil {
		return
	}

	tu = &TeamUsers{Team: *t}
	_, err = orm.NewOrm().Raw("SELECT `b`.`id`, `b`.`uuid`, `b`.`name`, `b`.`cname`, `b`.`email`, `b`.`phone`, `b`.`im`, `b`.`qq`, `b`.`create_time`  "+
		"FROM `team_user` `a` LEFT JOIN `user` `b` "+
		"ON `a`.`user_id` = `b`.`id` WHERE `a`.`team_id` = ? ",
		id).QueryRows(&tu.Users)
	if err == nil {
		for i := 0; i < len(tu.Users); i++ {
			tu.User_ids = append(tu.User_ids, tu.Users[i].Id)
		}
	}
	return
}

func (u *User) QueryTeams(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Team))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetTeamsCnt(query string) (int, error) {
	cnt, err := u.QueryTeams(query).Count()
	return int(cnt), err
}

func (u *User) GetTeams(query string, limit, offset int) (teams []*Team, err error) {
	_, err = u.QueryTeams(query).Limit(limit, offset).All(&teams)
	return
}

func (u *User) UpdateTeamUsers(id int64, _tu *TeamUsers) (tu *TeamUsers, err error) {
	if tu, err = u.GetTeamUsers(id); err != nil {
		return nil, ErrNoTeam
	}

	if _tu.Team.Name != "" {
		tu.Team.Name = _tu.Team.Name
	}
	if _tu.Team.Note != "" {
		tu.Team.Note = _tu.Team.Note
	}
	if _, err = orm.NewOrm().Update(&tu.Team); err != nil {
		return
	}
	add, del := MdiffInt(tu.User_ids, _tu.User_ids)
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
	tu.User_ids = _tu.User_ids

	DbLog(u.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_tu))

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
