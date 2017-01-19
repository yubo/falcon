/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

var (
	config = map[string]*ConfigEntry{}

	defconfig = map[string]*ConfigEntry{
		"ctrl": &ConfigEntry{"ctrl", "管理端配置", []ConfigEntry{
			{"host_bind_default_tag", "机器默认绑定节点名", "cop=xiaomi"},
			{"admin_uid", "管理员id号", 2},
		}},
		"agent": &ConfigEntry{"agent", "agent配置", []ConfigEntry{
			{"string", "配置string", "cop=xiaomi"},
			{"int", "配置int", 2},
			{"float64", "配置float64", float64(2.3)},
			{"slices", "配置slices", []string{"s1", "s2", "s3"}},
		}},
		"graph": &ConfigEntry{"graph", "graph配置", []ConfigEntry{
			{"sleep", "sleep second", 20},
			{"sleep1", "sleep second", "asdf"},
		}},
	}
)

func ConfigStart() (err error) {
	var rows []_ConfigEntry
	var conf __ConfigEntry

	_, err = orm.NewOrm().Raw("SELECT `key`, `note`, `value` FROM `kv` where "+
		"`type_id` = ?", KV_T_CONFIG).QueryRows(&rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		beego.Debug(row.Value)
		err = json.Unmarshal([]byte(row.Value), &conf)
		if err != nil {
			return
		}
		config[row.Key] = &ConfigEntry{Key: row.Key, Note: row.Note, Value: conf.Value}
	}
	for k, _ := range defconfig {
		if _, ok := config[k]; !ok {
			config[k] = defconfig[k]
		} else {
			beego.Debug(k)
			m := kvOverlay(defconfig[k], config[k], defconfig[k])
			config[k] = &m
		}

	}
	return nil
}

func (u *User) ConfigGet(k string) (*ConfigEntry, error) {
	if entry, ok := config[k]; ok {
		return entry, nil
	} else {
		return nil, ErrNoModule
	}
}

func (u *User) ConfigSet(name string, value map[string]string) (err error) {
	var (
		vb  []byte
		vc  string
		def *ConfigEntry
		ok  bool
	)

	// prepare
	if def, ok = defconfig[name]; !ok {
		return ErrNoModule
	}
	defValue := def.Value.([]ConfigEntry)
	cValue := make([]ConfigEntry, len(defValue))
	for idx, v := range defValue {
		if _, ok := value[v.Key]; !ok {
			cValue[idx] = v
			continue
		}
		cValue[idx].Key = v.Key
		cValue[idx].Note = v.Note

		switch v.Value.(type) {
		case int:
			cValue[idx].Value, err = strconv.Atoi(value[v.Key])
		case string:
			cValue[idx].Value = value[v.Key]
		case float32:
			cValue[idx].Value, err = strconv.ParseFloat(value[v.Key], 32)
		case float64:
			cValue[idx].Value, err = strconv.ParseFloat(value[v.Key], 64)
		case []string:
			cValue[idx].Value = strings.Split(value[v.Key], ",")
		default:
			beego.Debug("don't know the type")
		}
		if err != nil {
			return err
		}
	}
	c := ConfigEntry{
		Key:   def.Key,
		Note:  def.Note,
		Value: cValue,
	}

	if vb, err = json.Marshal(&c); err != nil {
		return err
	}
	config[c.Key] = &c
	vc = string(vb)
	_, err = orm.NewOrm().Raw("INSERT INTO `kv`(`key`, `note`, `value`, `type_id`) "+
		"VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
		c.Key, c.Note, vc, KV_T_CONFIG, vc).Exec()

	return
}
