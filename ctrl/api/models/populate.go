/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"github.com/golang/glog"
)

func (op *Operator) populate() (interface{}, error) {
	var (
		err               error
		items             []string
		id                int64
		tag_idx           = make(map[string]int64)
		user_idx          = make(map[string]int64)
		role_idx          = make(map[string]int64)
		token_idx         = make(map[string]int64)
		host_idx          = make(map[string]int64)
		event_trigger_idx = make(map[string]int64)
		test_user         = "test01"
	)
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
		if id, err = op.CreateUser(&UserApiAdd{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		user_idx[item] = id
		glog.Infof("add user(%s)\n", item)
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
		glog.Infof("add tag(%s)\n", item)
		if tag_idx[item], err = op.CreateTag(&TagCreate{Name: item}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
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
		glog.Infof("add host(%s, %s)\n", item2[1], item2[0])
		if host_idx[item2[1]], err = op.CreateHost(&HostCreate{Name: item2[1]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}

		if _, err = op.CreateTagHost(&TagHostApiAdd{TagId: tag_idx[item2[0]],
			HostId: host_idx[item2[1]]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// trigger
	event_triggers := [][]*EventTrigger{
		[]*EventTrigger{
			&EventTrigger{
				Name:     "cpu",
				Priority: 1,
				Metric:   "cpu.idle",
				Expr:     "all(#3)<5",
				Msg:      "cpu idle below the lower limit(5%)",
			},
			&EventTrigger{
				Priority: 0,
				Metric:   "cpu.idle",
				Expr:     "all(#3)<1",
				Msg:      "cpu idle below the lower limit(1%)",
			},
			&EventTrigger{
				Priority: 3,
				Metric:   "cpu.idle",
				Expr:     "all(#3)<10",
				Msg:      "cpu idle below the lower limit(10%)",
			},
		},
		[]*EventTrigger{
			&EventTrigger{
				Name:     "df",
				Priority: 2,
				Metric:   "df.bytes.free.percent",
				Expr:     "all(#3)<10",
				Msg:      "disk free below the lower limit",
			},
		},
		[]*EventTrigger{
			&EventTrigger{
				Name:     "disk",
				Priority: 1,
				Metric:   "disk.io.util",
				Expr:     "min(#3)>800",
				Msg:      "io over the upper limit",
			},
			&EventTrigger{
				Priority: 1,
				Metric:   "disk.io.util",
				Expr:     "all(#3)>96",
				Msg:      "io busy",
			},
		},
	}

	// step 1 init event_triggers
	for _, ts := range event_triggers {
		parentId := int64(0)
		for _, t := range ts {
			t.ParentId = parentId
			if id, err = op.CreateEventTrigger(t); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
			if parentId == 0 {
				event_trigger_idx[t.Name] = id
				parentId = id
			}
		}
	}

	// clone event triggers to tag
	items2 = [][2]string{
		{"cop=xiaomi", "cpu"},
		{"cop=xiaomi", "df"},
		{"cop=xiaomi,owt=inf", "disk"},
	}
	for _, item2 := range items2 {
		if _, err = op.CloneEventTrigger(event_trigger_idx[item2[1]], tag_idx[item2[0]]); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// step 2 clone triggers tpl to tag
	/*
		items = [][2]string{
			{"cop=xiaomi", "cpu"},
			{"cop=xiaomi", "df"},
			{"cop=xiaomi,owt=inf", "disk"},
		}
		for _, item2 := range items2 {
			glog.Infof("add tag triggers(%s, %s)\n", item2[0], item2[1])
			if _, err = op.CreateTagTpl(&RelTagTpl{TagId: tag_idx[item2[0]],
				TplId: tpl_idx[item2[1]]}); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
		}
	*/

	/*
		// template
		items = []string{
			"tpl1",
			"tpl2",
			"tpl3",
		}
		for _, item := range items {
			glog.Infof("add tag(%s)\n", item)
			if id, err = op.AddAction(&Action{}); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
			if tpl_idx[item], err = op.AddTemplate(&Template{Name: item,
				ActionId: id}); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
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
			glog.Infof("add strategy(%s, %s)\n", item2[0], item2[1])
			if _, err = op.AddStrategy(&Strategy{Metric: item2[1],
				TplId: tpl_idx[item2[0]]}); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
		}

		// clone template
		items = []string{
			"tpl1",
			"tpl2",
			"tpl3",
		}
		for _, item := range items {
			glog.Infof("clone template(%s)\n", item)
			if _, err = op.CloneTemplate(tpl_idx[item]); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
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
			glog.Infof("add tag tpl(%s, %s)\n", item2[0], item2[1])
			if _, err = op.CreateTagTpl(&RelTagTpl{TagId: tag_idx[item2[0]],
				TplId: tpl_idx[item2[1]]}); err != nil {
				glog.Error(err.Error())
				return nil, err
			}
		}
	*/

	// role
	items = []string{
		"adm",
		"sre",
		"dev",
		"usr",
	}
	for _, item := range items {
		glog.Infof("add role(%s)\n", item)
		if role_idx[item], err = op.CreateRole(&RoleCreate{Name: item}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// token
	for i := SYS_R_TOKEN; i < SYS_TOKEN_SIZE; i++ {
		token_idx[tokenName[i]] = int64(i)
	}

	// bind user
	binds := [][3]string{
		{"cop=xiaomi,owt=miliao", test_user, "adm"},
		{"cop=xiaomi,owt=miliao", test_user, "sre"},
		{"cop=xiaomi,owt=miliao", test_user, "dev"},
		{"cop=xiaomi,owt=miliao", test_user, "usr"},
	}
	for _, s := range binds {
		glog.Infof("bind tag(%s) user(%s) role(%s)\n", s[0], s[1], s[2])
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[0]], role_idx[s[2]],
			user_idx[s[1]], TPL_REL_T_ACL_USER); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// bind token
	binds = [][3]string{
		{tokenName[SYS_O_TOKEN], "adm", "/"},
		{tokenName[SYS_R_TOKEN], "adm", "/"},
		{tokenName[SYS_A_TOKEN], "adm", "/"},
		{tokenName[SYS_O_TOKEN], "sre", "/"},
		{tokenName[SYS_R_TOKEN], "sre", "/"},
		{tokenName[SYS_R_TOKEN], "dev", "/"},
		{tokenName[SYS_R_TOKEN], "usr", "/"},
		{tokenName[SYS_O_TOKEN], "adm", "cop=xiaomi,owt=miliao"},
		{tokenName[SYS_O_TOKEN], "dev", "cop=xiaomi,owt=miliao,pdl=op"},
		{tokenName[SYS_O_TOKEN], "usr", "cop=xiaomi"},
		{tokenName[SYS_A_TOKEN], "usr", "cop=xiaomi,owt=miliao"},
	}
	for _, s := range binds {
		glog.Infof("bind tag(%s) token(%s) role(%s)\n", s[2], s[0], s[1])
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[2]], role_idx[s[1]],
			token_idx[s[0]], TPL_REL_T_ACL_TOKEN); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	return "populate db done", nil
}
