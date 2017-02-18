/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

func Populate() (interface{}, error) {
	var (
		ret       string
		err       error
		items     []string
		user      *User
		id        int64
		tag_idx   = make(map[string]int64)
		user_idx  = make(map[string]int64)
		team_idx  = make(map[string]int64)
		role_idx  = make(map[string]int64)
		token_idx = make(map[string]int64)
		host_idx  = make(map[string]int64)
		tpl_idx   = make(map[string]int64)
		test_user = "test01"
	)
	admin, _ := GetUser(1)
	tag_idx["/"] = 1

	// user
	items = []string{
		test_user,
		"user0",
		"user1",
		"user2",
		"user3",
		"user4",
		"user5",
		"user6",
	}
	for _, item := range items {
		if user, err = admin.AddUser(&User{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		user_idx[item] = user.Id
		ret = fmt.Sprintf("%sadd user(%s)\n", ret, item)
	}

	// team
	items = []string{
		"team1",
		"team2",
		"team3",
		"team4",
	}
	for _, item := range items {
		if id, err = admin.AddTeam(&Team{Name: item, Creator: admin.Id}); err != nil {
			fmt.Printf("add team(%s)\n", item)
			return nil, err
		}
		team_idx[item] = id
		ret = fmt.Sprintf("%sadd team(%s)\n", ret, item)
	}
	teamMembers := []struct {
		team  string
		users []string
	}{
		{"team1", []string{"user0", "user1"}},
		{"team2", []string{"user2", "user3"}},
		{"team3", []string{"user4", "user5"}},
		{"team4", []string{"user0", "user1", "user2", "user3", "user4", "user5"}},
	}
	for _, item := range teamMembers {
		uids := make([]int64, len(item.users))
		for i := 0; i < len(uids); i++ {
			uids[i] = user_idx[item.users[i]]
		}
		if _, err = admin.UpdateMember(team_idx[item.team],
			&Member{Uids: uids}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd teamMembers(%v)\n", ret, item)
	}

	// tag
	items = []string{
		"cop=xiaomi",
		"cop=xiaomi,owt=inf",
		"cop=xiaomi,owt=miliao",
		"cop=xiaomi,owt=miliao,pdl=op",
		"cop=xiaomi,owt=miliao,pdl=micloud",
	}
	for _, item := range items {
		if tag_idx[item], err = admin.AddTag(&Tag{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag(%s)\n", ret, item)
	}

	// tag host
	items2 := [][2]string{
		{"cop=xiaomi", "mi1.bj"},
		{"cop=xiaomi", "mi2.bj"},
		{"cop=xiaomi", "mi3.bj"},
		{"cop=xiaomi,owt=inf", "inf1.bj"},
		{"cop=xiaomi,owt=inf", "inf2.bj"},
		{"cop=xiaomi,owt=inf", "inf3.bj"},
		{"cop=xiaomi,owt=miliao", "miliao1.bj"},
		{"cop=xiaomi,owt=miliao", "miliao2.bj"},
		{"cop=xiaomi,owt=miliao", "miliao3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud3.bj"},
	}
	for _, item2 := range items2 {
		if host_idx[item2[1]], err = admin.AddHost(&Host{Name: item2[1]}); err != nil {
			return nil, err
		}

		if _, err = admin.CreateTagHost(RelTagHost{TagId: tag_idx[item2[0]],
			HostId: host_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd host(%s, %s)\n", ret, item2[1], item2[0])
	}

	// template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		if id, err = admin.AddAction(&Action{}); err != nil {
			return nil, err
		}
		if tpl_idx[item], err = admin.AddTemplate(&Template{Name: item,
			ActionId: id}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag(%s)\n", ret, item)
	}
	// template strategy
	items2 = [][2]string{
		{"tpl1", "cpu.busy"},
		{"tpl1", "cpu.cnt"},
		{"tpl1", "cpu.idle"},
		{"tpl2", "cpu.busy"},
		{"tpl2", "cpu.cnt"},
		{"tpl2", "cpu.idle"},
		{"tpl3", "cpu.busy"},
		{"tpl3", "cpu.cnt"},
		{"tpl3", "cpu.idle"},
	}
	for _, item2 := range items2 {
		if _, err = admin.AddStrategy(&Strategy{Metric: item2[1],
			TplId: tpl_idx[item2[0]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd strategy(%s, %s)\n",
			ret, item2[0], item2[1])
	}

	// clone template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		if _, err = admin.CloneTemplate(tpl_idx[item]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%s clone template(%s)\n", ret, item)
	}

	// bind tag template
	items2 = [][2]string{
		{"cop=xiaomi", "tpl1"},
		{"cop=xiaomi", "tpl2"},
		{"cop=xiaomi", "tpl3"},
		{"cop=xiaomi,owt=inf", "tpl1"},
		{"cop=xiaomi,owt=inf", "tpl2"},
		{"cop=xiaomi,owt=inf", "tpl3"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl1"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl2"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl3"},
	}
	for _, item2 := range items2 {
		if _, err = admin.CreateTagTpl(RelTagTpl{TagId: tag_idx[item2[0]],
			TplId: tpl_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag tpl(%s, %s)\n", ret, item2[0], item2[1])
	}

	// role
	items = []string{
		"adm",
		"sre",
		"dev",
		"usr",
	}
	for _, item := range items {
		if role_idx[item], err = admin.AddRole(&Role{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd role(%s)\n", ret, item)
	}

	// token
	items = []string{
		SYS_W_SCOPE,
		SYS_R_SCOPE,
		SYS_B_SCOPE,
		SYS_O_SCOPE,
		SYS_A_SCOPE,
	}
	for _, item := range items {
		if token_idx[item], err = admin.AddToken(&Token{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd token(%s)\n", ret, item)
	}

	// bind user
	binds := [][3]string{
		{"cop=xiaomi,owt=miliao", test_user, "adm"},
		{"cop=xiaomi,owt=miliao", test_user, "sre"},
		{"cop=xiaomi,owt=miliao", test_user, "dev"},
		{"cop=xiaomi,owt=miliao", test_user, "usr"},
	}
	for _, s := range binds {
		if err := admin.BindAclUser(tag_idx[s[0]], role_idx[s[2]],
			user_idx[s[1]]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) user(%s) role(%s)\n",
			ret, s[0], s[1], s[2])
	}

	// bind token
	binds = [][3]string{
		{SYS_W_SCOPE, "adm", "/"},
		{SYS_R_SCOPE, "adm", "/"},
		{SYS_W_SCOPE, "sre", "/"},
		{SYS_R_SCOPE, "sre", "/"},
		{SYS_R_SCOPE, "dev", "/"},
		{SYS_R_SCOPE, "usr", "/"},
		{SYS_W_SCOPE, "adm", "cop=xiaomi,owt=miliao"},
		{SYS_R_SCOPE, "sre", "cop=xiaomi"},
		{SYS_B_SCOPE, "dev", "cop=xiaomi,owt=miliao,pdl=op"},
		{SYS_O_SCOPE, "usr", "cop=xiaomi"},
		{SYS_A_SCOPE, "usr", "cop=xiaomi,owt=miliao"},
	}
	for _, s := range binds {
		if err := admin.BindAclToken(tag_idx[s[2]], role_idx[s[1]],
			token_idx[s[0]]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) token(%s) role(%s)\n",
			ret, s[1], s[2], s[0])
	}

	return ret, nil
}

func ResetDb() (interface{}, error) {
	var err error
	o := orm.NewOrm()

	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range dbTables {
		if _, err = o.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return nil, err
		}
	}
	o.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	o.Insert(&User{Name: "system"})

	// init root tree tag
	o.Insert(&Tag{Name: ""})

	// reset cache
	cacheInit(config.Ctrl.Str(falcon.C_CACHE_MODULE))

	return "reset db done", nil
}
