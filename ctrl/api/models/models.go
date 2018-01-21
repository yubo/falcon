/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
)

const (
	DB_PREFIX   = ""
	PAGE_LIMIT  = 10
	MODULE_NAME = "\x1B[33m[CTRL_MODELS]\x1B[0m "
)

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

var (
	tokenName = [SYS_TOKEN_SIZE]string{
		"",
		"falcon_read",
		"falcon_operate",
		"falcon_admin",
	}
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

var (
	ModuleName = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "tpl",
		"user", "token", "rule", "event_trigger", "action_trigger",
		"tag_host", "dashboard_graph", "dashboard_screen", "tmp_graph",
	}
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

var (
	Db           *ctrl.OrmModule
	sysTagSchema *TagSchema
	admin        map[string]bool
	SysOp        *Operator

	// weixin app
	wxappid     string
	wxappsecret string
)

func (op *Operator) log(module, module_id, action int64, data_ interface{}) {
	data, ok := data_.(string)
	if !ok {
		b, _ := json.Marshal(data_)
		data = string(b)
	}
	op.O.Raw("insert log (user_id, module, module_id, action, data) values (?, ?, ?, ?, ?)", op.User.Id, module, module_id, action, data).Exec()
}

func initAuth(conf *falcon.Configer) error {
	Auths = make(map[string]AuthInterface)
	for _, name := range strings.Split(conf.Str(ctrl.C_AUTH_MODULE), ",") {
		if auth, ok := allAuths[name]; ok {
			if auth.Init(conf) == nil {
				Auths[name] = auth
			}
		}
	}
	return nil
}

func initConfigAdmin(conf *falcon.Configer) map[string]bool {
	ret := make(map[string]bool)
	for _, u := range strings.Split(conf.Str(ctrl.C_ADMIN), ",") {
		ret[u] = true
	}
	return ret
}

// called by (p *Ctrl) Init()
// already load file config and def config
// will load db config
func initConfig(conf *falcon.Configer) error {
	var err error

	// set default
	//conf.Agent.Set(falcon.APP_CONF_DEFAULT, config.ConfDefault["agent"])
	//conf.Loadbalance.Set(fconPfig.APP_CONF_DEFAULT, config.ConfDefault["loadbalance"])
	//conf.Backend.Set(falcon.APP_CONF_DEFAULT, config.ConfDefault["backend"])
	sysTagSchema, err = NewTagSchema(conf.Str(ctrl.C_TAG_SCHEMA))

	// admin
	admin = make(map[string]bool)
	for _, u := range strings.Split(conf.Str(ctrl.C_ADMIN), ",") {
		admin[u] = true
	}

	wxappid = conf.Str(ctrl.C_WEIXIN_APP_ID)
	wxappsecret = conf.Str(ctrl.C_WEIXIN_APP_SECRET)

	return err
}

/* called by initModels() */
func initCache(conf *falcon.Configer) error {
	for _, module := range strings.Split(
		conf.Str(ctrl.C_CACHE_MODULE), ",") {
		for k, v := range ModuleName {
			if v == module {
				moduleCache[k] = cache{
					enable: true,
					id:     make(map[int64]interface{}),
					key:    make(map[string]interface{}),
				}
				break
			}
		}
	}

	// build host cache
	if moduleCache[CTL_M_HOST].enable {
		var items []*Host
		for i := 0; ; i++ {
			n, err := Db.Ctrl.Raw("select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host limit ? offset ?", 100, 100*i).QueryRows(&items)
			if err != nil || n == 0 {
				break
			}

			for _, h := range items {
				moduleCache[CTL_M_HOST].set(h.Id, h, h.Name)
			}
		}
	}
	// build tag cache
	if moduleCache[CTL_M_TAG].enable {
		var items []*Tag
		for i := 0; ; i++ {
			n, err := Db.Ctrl.Raw("select id, name, type, create_time from tag LIMIT ? OFFSET ?", 100, 100*i).QueryRows(&items)
			if err != nil || n == 0 {
				break
			}

			for _, t := range items {
				moduleCache[CTL_M_TAG].set(t.Id, t, t.Name)
			}
		}
	}

	go func() {
		cacheTree = &cacheTreeT{
			m: make(map[int64]*TreeNode),
			c: make(chan struct{}, 1),
		}
		cacheTree._build()
		ticker := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ticker.C:
				cacheTree._build()
			case <-cacheTree.c:
				cacheTree._build()
			}
		}
	}()

	return nil
}

func putEtcdConfig() error {
	put := make(map[string]string)
	for module, kv := range EtcdMap {
		ks := make(map[string]bool)
		for _, k := range kv {
			ks[k] = false
		}

		prefix := fmt.Sprintf("/open-falcon/%s/config/", module)
		resp, err := ctrl.EtcdGetPrefix(prefix)
		if err != nil {
			return err
		}
		for _, v := range resp.Kvs {
			if _, ok := ks[string(v.Key)]; ok {
				ks[string(v.Key)] = true
			}
		}
		for k, exist := range ks {
			if !exist {
				put[k] = ""
			}
		}
	}
	return ctrl.EtcdPuts(put)
}
