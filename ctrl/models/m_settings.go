/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

func Populate() (interface{}, error) {
	var (
		ret       string
		err       error
		items     []string
		user      *User
		tag_idx   = make(map[string]int64)
		user_idx  = make(map[string]int64)
		role_idx  = make(map[string]int64)
		token_idx = make(map[string]int64)
		host_idx  = make(map[string]int64)
		test_user = "test01"
	)
	admin, _ := GetUser(1)
	tag_idx["/"] = 1

	// user
	items = []string{
		test_user,
	}
	for _, item := range items {
		if user, err = admin.AddUser(&User{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		user_idx[item] = user.Id
		ret = fmt.Sprintf("%sadd user(%s)\n", ret, item)
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

	// host
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

		if _, err = admin.CreateTagHost(RelTagHost{Tag_id: tag_idx[item2[0]],
			Host_id: host_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd host(%s, %s)\n", ret, item2[1], item2[0])
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
	tables := []string{
		"host",
		"kv",
		"log",
		"role",
		"tag",
		"tag_host",
		"tag_rel",
		"token",
		"tpl_rel",
		"user",
		"team",
		"team_user",
		"rule",
		"action",
		"trigger",
		"expression",
	}
	o := orm.NewOrm()

	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range tables {
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
	CacheInit()

	return "reset db done", nil
}
