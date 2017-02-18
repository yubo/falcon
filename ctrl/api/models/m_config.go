/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

type Kv struct {
	Key     string
	Section string
	Value   string
}

func GetDbConfig(module string) (ret map[string]string, err error) {
	var row Kv

	err = orm.NewOrm().Raw("SELECT `section`, `key`, `value` FROM `kv` where "+
		"`section` = ? and `key` = 'config'", module).QueryRow(&row)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(row.Value), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetDbConfig(module string, conf map[string]string) error {
	kv := make(map[string]string)
	for k, v := range conf {
		if v != "" {
			kv[k] = v
		}
	}
	v, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	s := string(v)
	_, err = orm.NewOrm().Raw("INSERT INTO `kv`(`section`, `key`, `value`)"+
		" VALUES (?,'config',?) ON DUPLICATE KEY UPDATE `value`=?",
		module, s, s).Exec()

	return err
}

func (u *User) ConfigGet(module string) (interface{}, error) {
	var c *falcon.Configer

	switch module {
	case "ctrl":
		c = &config.Ctrl
	case "agent":
		c = &config.Agent
	case "lb":
		c = &config.Lb
	case "backend":
		c = &config.Backend
	default:
		return nil, ErrNoModule
	}

	conf, err := GetDbConfig(module)
	if err == nil {
		c.Set(falcon.APP_CONF_DB, conf)
	}
	return c.Get(), nil
}

func (u *User) ConfigSet(module string, conf map[string]string) error {
	switch module {
	case "ctrl", "agent", "lb", "backend":
		return SetDbConfig(module, conf)
	default:
		return ErrNoModule
	}
}
