/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"github.com/astaxie/beego/orm"
)

const (
	PAGE_LIMIT = 10
)

type CtrlDb struct {
	Ctrl  orm.Ormer
	Idx   orm.Ormer
	Alarm orm.Ormer
}

const (
	CTL_SEARCH_CUR = iota
	CTL_SEARCH_PARENT
	CTL_SEARCH_CHILD
)

const (
	SYS_F_R_TOKEN = 1 << iota
	SYS_F_O_TOKEN
	SYS_F_A_TOKEN
)

const (
	_ = iota
	SYS_R_TOKEN
	SYS_O_TOKEN
	SYS_A_TOKEN
	SYS_TOKEN_SIZE
)

// ctl meta name
const (
	CTL_M_HOST = iota
	CTL_M_ROLE
	CTL_M_SYSTEM
	CTL_M_TAG
	CTL_M_TPL

	CTL_M_USER
	CTL_M_TOKEN
	CTL_M_RULE
	CTL_M_EVENT_TRIGGER
	CTL_M_ACTION_TRIGGER

	CTL_M_TAG_HOST
	CTL_M_DASHBOARD_GRAPH
	CTL_M_DASHBOARD_SCREEN
	CTL_M_TMP_GRAPH
	CTL_M_SIZE
)

// ctl method name
const (
	CTL_A_ADD = iota
	CTL_A_DEL
	CTL_A_SET
	CTL_A_GET
	CTL_A_SIZE
)

var (
	ActionName = [CTL_A_SIZE]string{
		"add", "del", "set", "get",
	}
)

type Ids struct {
	Ids []int64 `json:"ids"`
}

type Id struct {
	Id int64 `json:"id"`
}

type Total struct {
	Total int64 `json:"total"`
}

type Stats struct {
	Success int64  `json:"success"`
	Err     string `json:"err"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
