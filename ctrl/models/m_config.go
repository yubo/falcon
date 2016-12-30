/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"

	"github.com/astaxie/beego/orm"
)

type Conf struct {
	Key   string
	Value string
}

var (
	config    = map[string]*map[string]string{}
	defconfig = map[string]*map[string]string{
		"global": &map[string]string{
			"host_bind_default_tag": "cop=xiaomi",
			"admin_uid":             "1",
		},
	}
)

func ConfigStart() (err error) {
	var confs []Conf

	_, err = orm.NewOrm().Raw("SELECT `key`, `value` FROM `config`").QueryRows(&confs)
	if err != nil {
		return err
	}
	for _, conf := range confs {
		m := make(map[string]string)
		config[conf.Key] = &m
		err = json.Unmarshal([]byte(conf.Value), config[conf.Key])
		if err != nil {
			return
		}
	}
	for k, _ := range defconfig {
		if _, ok := config[k]; !ok {
			config[k] = defconfig[k]
		} else {
			m := kvOverlay(*defconfig[k], *config[k], *defconfig[k])
			config[k] = &m
		}

	}
	return nil
}

func (u *User) ConfigGet(k string) (map[string]string, error) {
	if v, ok := config[k]; ok {
		return *v, nil
	} else {
		return nil, ErrNoModule
	}
}

func (u *User) ConfigSet(k string, v map[string]string) (err error) {
	var (
		vb []byte
		vs string
	)

	if vb, err = json.Marshal(v); err != nil {
		return
	}
	vs = string(vb)
	config[k] = &v
	_, err = orm.NewOrm().Raw("INSERT INTO `config`(`key`,`value`) VALUES (?,?) ON DUPLICATE KEY UPDATE `value`=?", k, vs, vs).Exec()

	return
}
