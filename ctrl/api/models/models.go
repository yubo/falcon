/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/config"
	"github.com/yubo/gotool/ratelimits"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yubo/falcon/ctrl/api/models/session"
)

const (
	DB_PREFIX   = ""
	PAGE_LIMIT  = 10
	MODULE_NAME = "\x1B[32m[CTRL_MODELS]\x1B[0m "
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

var (
	SysOp    *Operator
	dbTables = []string{
		"agents_info",
		"plugin_dir",
		"dashboard_graph",
		"dashboard_screen",
		"tmp_graph",
		"kv",
		"host",
		"token",
		"role",
		"session",
		"tag",
		"tag_rel",
		"tag_host",
		"user",
		"log",
		"tpl_rel",
		"triggers",
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
	CTL_M_TRIGGER
	CTL_M_EXPRESSION

	CTL_M_TAG_HOST
	CTL_M_DASHBOARD_GRAPH
	CTL_M_DASHBOARD_SCREEN
	CTL_M_TMP_GRAPH
	CTL_M_SIZE
)

var (
	ModuleName = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "tpl",
		"user", "token", "rule", "trigger", "expression",
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

// ctl runmode name
const (
	CTL_RUNMODE_MASTER = 1 << iota
	CTL_RUNMODE_DEV
	CTL_RUNMODE_MI
)

type Obj struct {
}

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
	Success int64 `json:"success"`
	Failure int64 `json:"failure"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Falcon_db struct {
	Ctrl  orm.Ormer
	Idx   orm.Ormer
	Alarm orm.Ormer
}

var (
	//RunMode      string
	Db           Falcon_db
	sysTagSchema *TagSchema
	transferUrl  string
	RunMode      uint32
	ApiRl        *ratelimits.RateLimits
	admin        map[string]bool

	// weixin app
	wxappid     string
	wxappsecret string
)

func PreStart(conf *config.ConfCtrl) (err error) {
	if err = initConfig(conf); err != nil {
		panic(err)
	}
	if err = initAuth(conf); err != nil {
		panic(err)
	}
	if err = initCache(conf); err != nil {
		panic(err)
	}
	if err = initMetric(conf); err != nil {
		panic(err)
	}
	prepareEtcdConfig()

	SysOp = &Operator{
		O:     orm.NewOrm(),
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_A_TOKEN,
	}
	SysOp.User, _ = GetUser(1, SysOp.O)

	return nil
}

func initMetric(c *config.ConfCtrl) error {
	for _, m := range c.Metrics {
		metrics = append(metrics, &Metric{Name: m})
	}
	return nil
}

func initAuth(c *config.ConfCtrl) error {
	Auths = make(map[string]AuthInterface)
	for _, name := range strings.Split(c.Ctrl.Str(ctrl.C_AUTH_MODULE), ",") {
		if auth, ok := allAuths[name]; ok {
			if auth.Init(c) == nil {
				Auths[name] = auth
			}
		}
	}
	return nil
}

func initConfigRl(c *fconfig.Configer) *ratelimits.RateLimits {
	limit, _ := c.Int(ctrl.C_RL_LIMIT)
	accuracy, _ := c.Int(ctrl.C_RL_ACCURACY)
	if limit <= 0 {
		return nil

	}
	rl, err := ratelimits.New(uint32(limit), uint32(accuracy))
	if err != nil {
		return nil
	}
	timeout, _ := c.Int(ctrl.C_RL_GC_TIMEOUT)
	interval, _ := c.Int(ctrl.C_RL_GC_INTERVAL)
	err = rl.GcStart(time.Duration(timeout)*time.Millisecond, time.Duration(interval)*time.Millisecond)
	if err != nil {
		return nil
	}
	return rl
}

func initConfigAdmin(c *fconfig.Configer) map[string]bool {
	ret := make(map[string]bool)
	for _, u := range strings.Split(c.Str(ctrl.C_ADMIN), ",") {
		ret[u] = true
	}
	return ret
}

// called by (p *Ctrl) Init()
// already load file config and def config
// will load db config
func initConfig(conf *config.ConfCtrl) error {
	var err error

	glog.V(4).Infof(MODULE_NAME+"%s Init()", conf.Name)

	orm.RegisterDriver("mysql", orm.DRMySQL)

	// set default
	//conf.Agent.Set(fconfig.APP_CONF_DEFAULT, config.ConfDefault["agent"])
	//conf.Loadbalance.Set(fconPfig.APP_CONF_DEFAULT, config.ConfDefault["loadbalance"])
	//conf.Backend.Set(fconfig.APP_CONF_DEFAULT, config.ConfDefault["backend"])

	// ctrl config
	conf.Ctrl.Set(fconfig.APP_CONF_DEFAULT, ctrl.ConfDefault)
	cf := &conf.Ctrl
	dsn := cf.Str(ctrl.C_DSN)
	dbMaxConn, _ := cf.Int(ctrl.C_DB_MAX_CONN)
	dbMaxIdle, _ := cf.Int(ctrl.C_DB_MAX_IDLE)
	// connect db, can not register db twice  :(
	orm.RegisterDataBase("default", "mysql", dsn, dbMaxIdle, dbMaxConn)
	// get ctrl config from db
	Db.Ctrl = orm.NewOrm()
	if c, err := GetDbConfig(Db.Ctrl, "ctrl"); err == nil {
		conf.Ctrl.Set(fconfig.APP_CONF_DB, c)
	}
	glog.V(4).Infof(MODULE_NAME+"initConfig get config %s", cf.String())

	// ctrl config
	if err = orm.RegisterDataBase("idx", "mysql", cf.Str(ctrl.C_IDX_DSN), dbMaxIdle, dbMaxConn); err != nil {
		return err
	}
	Db.Idx = orm.NewOrm()
	Db.Idx.Using("idx")

	if err = orm.RegisterDataBase("alarm", "mysql", cf.Str(ctrl.C_ALARM_DSN), dbMaxIdle, dbMaxConn); err != nil {
		return err
	}
	Db.Alarm = orm.NewOrm()
	Db.Alarm.Using("alarm")

	sysTagSchema, err = NewTagSchema(cf.Str(ctrl.C_TAG_SCHEMA))
	transferUrl = cf.Str(ctrl.C_TRANSFER_URL)

	if cf.DefaultBool(ctrl.C_MASTER_MODE, false) {
		RunMode |= CTL_RUNMODE_MASTER
	}
	if cf.DefaultBool(ctrl.C_MI_MODE, false) {
		RunMode |= CTL_RUNMODE_MI
	}
	if cf.DefaultBool(ctrl.C_DEV_MODE, false) {
		RunMode |= CTL_RUNMODE_DEV
	}

	if RunMode&CTL_RUNMODE_MI != 0 {
		url := cf.Str(ctrl.C_MI_NORNS_URL)
		interval, _ := cf.Int(ctrl.C_MI_NORNS_INTERVAL)
		miStart(url, interval)
	}

	// ratelimits
	ApiRl = initConfigRl(cf)

	// admin
	admin = initConfigAdmin(cf)

	// ctrl beggo config
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"
	beego.BConfig.AppName = conf.Name
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime, _ = cf.Int64(ctrl.C_SESSION_GC_MAX_LIFETIME)
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime, _ = cf.Int(ctrl.C_SESSION_COOKIE_LIFETIME)
	if RunMode&CTL_RUNMODE_DEV != 0 {
		orm.Debug = true
		beego.BConfig.RunMode = "dev"
		beego.BConfig.WebConfig.EnableDocs = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	} else {
		beego.BConfig.RunMode = "prod"
	}
	if addr := strings.Split(cf.Str(ctrl.C_HTTP_ADDR), ":"); len(addr) == 2 {
		beego.BConfig.Listen.HTTPAddr = addr[0]
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[1])
	} else if len(addr) == 1 {
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[0])
	}

	wxappid = cf.Str(ctrl.C_WEIXIN_APP_ID)
	wxappsecret = cf.Str(ctrl.C_WEIXIN_APP_SECRET)

	return err
}
