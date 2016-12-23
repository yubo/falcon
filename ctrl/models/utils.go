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
	Id        int
	Module    int
	Module_id int
	User_id   int
	Action    int
	Data      string
	Time      time.Time
}

type cache struct {
	enable bool
	data   map[int]interface{}
}

func (c *cache) set(id int, p interface{}) {
	if c.enable {
		c.data[id] = p
	}
}

func (c *cache) get(id int) interface{} {
	return c.data[id]
}

func (c *cache) del(id int) {
	if c.enable {
		delete(c.data, id)
	}
}

func DbLog(uid, module, module_id, action int, data string) {
	log := &Log{
		User_id:   uid,
		Module:    module,
		Module_id: module_id,
		Action:    action,
		Data:      data,
	}
	orm.NewOrm().Insert(log)
}
