/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Tag struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Create_time time.Time `json:"-"`
}

func AddTag(t *Tag) (int, error) {
	id, err := orm.NewOrm().Insert(t)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	t.Id = int(id)
	cacheModule[CTL_M_TAG].set(t.Id, t)
	return t.Id, err
}

func GetTag(id int) (*Tag, error) {
	if t, ok := cacheModule[CTL_M_TAG].get(id).(*Tag); ok {
		return t, nil
	}
	t := &Tag{Id: id}
	err := orm.NewOrm().Read(t, "Id")
	if err == nil {
		cacheModule[CTL_M_TAG].set(id, t)
	}
	return t, err
}

func QueryTags(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Tag))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query))
	}
	return qs
}

func GetTagsCnt(query string) (int, error) {
	cnt, err := QueryTags(query).Count()
	return int(cnt), err
}

func GetTags(query string, limit, offset int) (tags []*Tag, err error) {
	_, err = QueryTags(query).Limit(limit, offset).All(&tags)
	return
}

func UpdateTag(id int, _t *Tag) (t *Tag, err error) {
	if t, err = GetTag(id); err != nil {
		return nil, ErrNoTag
	}

	if _t.Name != "" {
		t.Name = _t.Name
	}
	_, err = orm.NewOrm().Update(t)
	return t, err
}

func DeleteTag(id int) error {

	if n, err := orm.NewOrm().Delete(&Tag{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_TAG].del(id)

	return nil
}
