/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/specs"
)

type ConfigEntry struct {
	Key   string
	Note  string
	Value interface{}
}

type _ConfigEntry struct {
	Key   string
	Note  string
	Value string
}

type __ConfigEntry struct {
	Key   string
	Note  string
	Value []ConfigEntry
}

func ConfigGet(k string, def interface{}) interface{} {
	var row _ConfigEntry
	err := orm.NewOrm().Raw("SELECT `key`, `note`, `value` FROM `kv` where "+
		"`key` = ? and `type_id` = ?", k, KV_T_CONFIG).QueryRow(&row)
	if err != nil {
		return def
	}

	ptr := reflect.New(reflect.ValueOf(def).Elem().Type()).
		Elem().Addr().Interface()

	err = json.Unmarshal([]byte(row.Value), ptr)
	if err != nil {
		return def
	}
	return ptr
}

func (u *User) ConfigGet(k string) (interface{}, error) {
	switch k {
	case "ctrl":
		return ConfigGet(k, &specs.ConfCtrlDef), nil
	case "agent":
		return ConfigGet(k, &specs.ConfAgentDef), nil
	case "lb":
		return ConfigGet(k, &specs.ConfLbDef), nil
	case "backend":
		return ConfigGet(k, &specs.ConfBackendDef), nil
	default:
		return nil, ErrNoModule
	}
}

func (u *User) ConfigSet(k string, v []byte) (err error) {
	var conf interface{}
	switch k {
	case "ctrl":
		conf = &specs.ConfCtrl{}
		err = json.Unmarshal(v, conf)
	case "agent":
		conf = &specs.ConfAgent{}
		err = json.Unmarshal(v, conf)
	case "lb":
		conf = &specs.ConfLb{}
		err = json.Unmarshal(v, conf)
	case "backend":
		conf = &specs.ConfBackend{}
		err = json.Unmarshal(v, conf)
	default:
		return ErrNoModule
	}

	v, err = json.Marshal(conf)
	if err != nil {
		return err
	}
	s := string(v)
	_, err = orm.NewOrm().Raw("INSERT INTO `kv`(`key`, `value`, `type_id`)"+
		" VALUES (?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
		k, s, KV_T_CONFIG, s).Exec()
	return err
}

// ugly hack
// should called by main package
func ConfInit(conf map[string]string) {
	var (
		dsn     string
		maxIdle int
		maxConn int
	)

	if _, err := os.Stat("./conf/app.conf"); os.IsNotExist(err) {
		dsn = conf["mysqldsn"]
		maxIdle, _ = strconv.Atoi(conf["mysqlmaxidle"])
		maxConn, _ = strconv.Atoi(conf["mysqlmaxconn"])

		cfg := ConfigGet("ctrl", &specs.ConfCtrlDef).(*specs.ConfCtrl)
		beego.BConfig.AppName = cfg.AppName
		beego.BConfig.RunMode = cfg.RunMode
		beego.BConfig.Listen.HTTPPort = cfg.HttpPort
		beego.BConfig.WebConfig.EnableDocs = cfg.EnableDocs
		beego.BConfig.WebConfig.Session.SessionName = cfg.SessionName
		beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = cfg.SessionGCMaxLifetime
		beego.BConfig.WebConfig.Session.SessionCookieLifeTime = cfg.SessionCookieLifeTime
	} else {
		dsn = beego.AppConfig.String("mysqldsn")
		maxIdle, _ = beego.AppConfig.Int("mysqlmaxidle")
		maxConn, _ = beego.AppConfig.Int("mysqlmaxconn")
	}

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxConn)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	}
}
