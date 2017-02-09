/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"fmt"
	"sort"
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

func DbLog(uid, module, module_id, action int64, data string) {
	log := &Log{
		User_id:   uid,
		Module:    module,
		Module_id: module_id,
		Action:    action,
		Data:      data,
	}
	orm.NewOrm().Insert(log)
}

func array2sql(array []int64) string {
	var ret string
	if len(array) == 0 {
		return "()"
	}

	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("%d,", array[i])
	}
	return fmt.Sprintf("(%s)", ret[:len(ret)-1])
}

func stringscmp(a, b []string) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := 0; i < len(a); i++ {
		if ret = strings.Compare(a[i], b[i]); ret != 0 {
			return
		}
	}
	return
}

func jsonStr(i interface{}) string {
	if ret, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return string(ret)
	}
}

func MdiffStr(src, dst []string) (add, del []string) {
	_src := make(map[string]bool)
	_dst := make(map[string]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}
func MdiffInt(src, dst []int64) (add, del []int64) {
	_src := make(map[int64]bool)
	_dst := make(map[int64]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}
