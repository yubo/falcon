/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Log struct {
	Id        int64
	Module    int64
	Module_id int64
	User_id   int64
	Action    int64
	Data      string
	Time      time.Time
}

type cache struct {
	enable bool
	data   map[int64]interface{}
}

func (c *cache) set(id int64, p interface{}) {
	if c.enable {
		c.data[id] = p
	}
}

func (c *cache) get(id int64) interface{} {
	return c.data[id]
}

func (c *cache) del(id int64) {
	if c.enable {
		delete(c.data, id)
	}
}

func DbLog(uid, module, module_id, action int64, data []byte) {
	log := &Log{
		User_id:   uid,
		Module:    module,
		Module_id: module_id,
		Action:    action,
		Data:      string(data),
	}
	orm.NewOrm().Insert(log)
}
