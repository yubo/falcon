/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"reflect"

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

func configGet(k string, def interface{}) interface{} {
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
		return configGet(k, &specs.ConfCtrlDef), nil
	case "agent":
		return configGet(k, &specs.ConfAgentDef), nil
	case "lb":
		return configGet(k, &specs.ConfLbDef), nil
	case "backend":
		return configGet(k, &specs.ConfBackendDef), nil
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
