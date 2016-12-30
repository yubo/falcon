/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"
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

/* dst <- src format tpl */
func kvOverlay(dst, src, tpl map[string]string) map[string]string {
	ret := make(map[string]string)
	for k, _ := range tpl {
		if v, ok := src[k]; ok {
			ret[k] = v
		} else if v, ok := dst[k]; ok {
			ret[k] = v
		}
	}
	return ret
}

func stringscmp(a, b []string) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	for i := 0; i < len(a); i++ {
		if ret = strings.Compare(a[i], b[i]); ret != 0 {
			return
		}

	}
	return
}
