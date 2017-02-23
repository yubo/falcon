/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
	_ "github.com/yubo/falcon/ctrl/api/models/session"
)

const (
	DB_PREFIX      = ""
	PAGE_PER       = 10
	SYS_TAG_SCHEMA = "cop,owt,pdl;servicegroup;service,jobgroup;job,sbs;mod;srv;grp;cluster;"
	SYS_R_TOKEN    = "falcon_read"
	SYS_O_TOKEN    = "falcon_operate"
	SYS_A_TOKEN    = "falcon_admin"
)

const (
	SYS_F_R_TOKEN = 1 << iota
	SYS_F_O_TOKEN
	SYS_F_A_TOKEN
)

const (
	_ = iota
	SYS_IDX_R_TOKEN
	SYS_IDX_O_TOKEN
	SYS_IDX_A_TOKEN
)

var (
	dbTables = []string{
		"action",
		"expression",
		"host",
		"kv",
		"log",
		"role",
		"session",
		"strategy",
		"tag",
		"tag_host",
		"tag_rel",
		"tag_tpl",
		"team",
		"team_user",
		"template",
		"token",
		"tpl_rel",
		"trigger",
		"user",
	}
)

// ctl meta name
const (
	CTL_M_HOST = iota
	CTL_M_ROLE
	CTL_M_SYSTEM
	CTL_M_TAG
	CTL_M_USER
	CTL_M_TOKEN
	CTL_M_TPL
	CTL_M_RULE
	CTL_M_TEMPLATE
	CTL_M_TRIGGER
	CTL_M_EXPRESSION
	CTL_M_TEAM
	CTL_M_TAG_HOST
	CTL_M_TAG_TPL
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

type Ids struct {
	Ids []int64 `json:"ids"`
}

type Id struct {
	Id int64 `json:"id"`
}

type Total struct {
	Total int64 `json:"total"`
}

var (
	moduleCache  [CTL_M_SIZE]cache
	sysTagSchema *TagSchema

	moduleName = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "user", "token",
		"template", "rule", "trigger", "expression", "team",
	}

	actionName = [CTL_A_SIZE]string{
		"add", "del", "set", "get",
	}
)

func prestart(conf *falcon.ConfCtrl) (err error) {
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
	return nil
}

func initMetric(c *falcon.ConfCtrl) error {
	for _, m := range c.Metrics {
		metrics = append(metrics, &Metric{Name: m})
	}
	return nil
}

func initAuth(c *falcon.ConfCtrl) error {
	for _, name := range strings.Split(c.Ctrl.Str(falcon.C_AUTH_MODULE), ",") {
		if auth, ok := allAuths[name]; ok {
			if auth.Init(c) == nil {
				Auths[name] = auth
			}
		}
	}
	return nil
}

// called by (p *Ctrl) Init()
func initConfig(conf *falcon.ConfCtrl) error {

	beego.Debug(fmt.Sprintf("%s Init()", conf.Name))
	// config
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = conf.Dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"

	// connect db
	orm.RegisterDataBase("default", "mysql", conf.Dsn, conf.DbMaxIdle, conf.DbMaxConn)
	conf.Agent.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["agent"])
	conf.Lb.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["lb"])
	conf.Backend.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["backend"])
	conf.Ctrl.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["ctrl"])

	// get config from db
	o := orm.NewOrm()
	if c, err := GetDbConfig(o, "ctrl"); err == nil {
		conf.Ctrl.Set(falcon.APP_CONF_DB, c)
	}

	// config -> beego config
	c := &conf.Ctrl
	beego.BConfig.AppName = conf.Name
	beego.BConfig.RunMode = c.Str(falcon.C_RUN_MODE)
	beego.BConfig.Listen.HTTPPort, _ = c.Int(falcon.C_HTTP_PORT)
	beego.BConfig.WebConfig.EnableDocs, _ = c.Bool(falcon.C_ENABLE_DOCS)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime, _ = c.Int64(falcon.C_SEESION_GC_MAX_LIFETIME)
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime, _ = c.Int(falcon.C_SESSION_COOKIE_LIFETIME)

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("orm debug on")
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	}

	return nil
}

func initCache(c *falcon.ConfCtrl) error {
	for _, module := range strings.Split(
		c.Ctrl.Str(falcon.C_CACHE_MODULE), ",") {
		for k, v := range moduleName {
			if v == module {
				moduleCache[k] = cache{
					enable: true,
					data:   make(map[int64]interface{}),
				}
				break
			}
		}
	}
	return nil
}
func init() {
	// tag
	sysTagSchema, _ = NewTagSchema(SYS_TAG_SCHEMA)

	// auth
	allAuths = make(map[string]AuthInterface)
	Auths = make(map[string]AuthInterface)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModelWithPrefix("",
		new(User), new(Host), new(Tag),
		new(Role), new(Token), new(Log),
		new(Tag_rel), new(Tpl_rel), new(Team),
		new(Template), new(Trigger), new(Expression),
		new(Action), new(Strategy))

	ctrl.RegisterInit(prestart)
}
