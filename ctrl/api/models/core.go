/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/transfer"
)

const (
	//MODULE_NAME = "\x1B[33m[CTRL_MODELS]\x1B[0m "
	MODULE_NAME = "models"
)

var (
	tokenName = [SYS_TOKEN_SIZE]string{
		"",
		"falcon_read",
		"falcon_operate",
		"falcon_admin",
	}

	ModuleName = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "tpl",
		"user", "token", "rule", "event_trigger", "action_trigger",
		"tag_host", "dashboard_graph", "dashboard_screen", "tmp_graph",
	}
)

type ModelsConfig struct {
	Configer    *core.Configer `json:"-"`
	CacheModule []string       `json:"cache_module"`
	Admin       []string       `json:"admin"`
	TagSchema   string         `json:"tag_schema"`
	WxAppId     string         `json:"wx_app_id"`
	WxAppSecret string         `json:"wx_app_secret"`
	AuthModule  []string       `json:"auth_module"`
	DbSchema    string         `json:"db_schema"`
	MasterMode  bool           `json:"master_mode"`
	DevMode     bool           `json:"dev_mode"`
	MiMode      bool           `json:"mi_mode"`
	CallTimeout int            `json:"call_timeout"`
}

type Models struct {
	conf         *ModelsConfig
	db           *CtrlDb
	sysTagSchema *TagSchema
	admin        map[string]bool
	sysOp        *Operator
	etcdCli      *core.EtcdCli
	transferCli  transfer.TransferClient
}

var (
	_models *Models
	//Db           *CtrlDb
	//sysTagSchema *TagSchema
	//admin        map[string]bool
	//SysOp        *Operator
)

func Init(configer *core.Configer, db *CtrlDb, cli *core.EtcdCli,
	tCli transfer.TransferClient) error {
	conf := &ModelsConfig{}
	if err := configer.Unmarshal(conf); err != nil {
		return err
	}
	conf.Configer = configer

	if err := InitConfig(conf); err != nil {
		return err
	}
	if err := InitAuth(conf); err != nil {
		return err
	}
	if err := InitCache(conf); err != nil {
		return err
	}
	_models = &Models{
		db:          db,
		conf:        conf,
		etcdCli:     cli,
		transferCli: tCli,
	}

	op := &Operator{
		O:     db.Ctrl,
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_A_TOKEN,
	}
	op.User, _ = GetUser(1, op.O)
	_models.sysOp = op

	return PutEtcdConfig()
}

func (op *Operator) log(module, module_id, action int64, data_ interface{}) {
	data, ok := data_.(string)
	if !ok {
		b, _ := json.Marshal(data_)
		data = string(b)
	}
	op.O.Raw("insert log (user_id, module, module_id, action, data) values (?, ?, ?, ?, ?)", op.User.Id, module, module_id, action, data).Exec()
}

func InitAuth(conf *ModelsConfig) error {
	Auths = make(map[string]AuthInterface)
	for _, name := range conf.AuthModule {
		if auth, ok := allAuths[name]; ok {
			if auth.Init(conf.Configer.GetConfiger("auth."+name)) == nil {
				Auths[name] = auth
			}
		}
	}
	return nil
}

// called by (p *Ctrl) Init()
// already load file config and def config
// will load db config
func InitConfig(conf *ModelsConfig) error {
	var err error

	// set default
	//conf.Agent.Set(falcon.APP_CONF_DEFAULT, config.ConfDefault["agent"])
	//conf.Loadbalance.Set(fconPfig.APP_CONF_DEFAULT, config.ConfDefault["loadbalance"])
	//conf.Backend.Set(falcon.APP_CONF_DEFAULT, config.ConfDefault["backend"])
	_models.sysTagSchema, err = NewTagSchema(conf.TagSchema)

	// admin
	_models.admin = make(map[string]bool)
	for _, u := range conf.Admin {
		_models.admin[u] = true
	}

	return err
}

/* called by initModels() */
func InitCache(conf *ModelsConfig) error {
	for _, module := range conf.CacheModule {
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
			n, err := _models.db.Ctrl.Raw("select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host limit ? offset ?", 100, 100*i).QueryRows(&items)
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
			n, err := _models.db.Ctrl.Raw("select id, name, type, create_time from tag LIMIT ? OFFSET ?", 100, 100*i).QueryRows(&items)
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

func PutEtcdConfig() error {
	put := make(map[string]string)
	for module, kv := range EtcdMap {
		ks := make(map[string]bool)
		for _, k := range kv {
			ks[k] = false
		}

		prefix := fmt.Sprintf("/open-falcon/%s/config/", module)
		resp, err := _models.etcdCli.GetPrefix(prefix)
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
	return _models.etcdCli.Puts(put)
}

// call by api /pub/config/ctrl
func GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"auth_module": _models.conf.AuthModule,
		"master_mode": _models.conf.MasterMode,
		"dev_mode":    _models.conf.DevMode,
		"mi_mode":     _models.conf.MiMode,
		"tag_schema":  _models.conf.TagSchema,
	}
}
