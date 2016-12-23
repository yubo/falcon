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

type Host struct {
	Id          int       `json:"id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Loc         string    `json:"loc"`
	Idc         string    `json:"idc"`
	Create_time time.Time `json:"-"`
}

func AddHost(h *Host) (int, error) {
	id, err := orm.NewOrm().Insert(h)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	h.Id = int(id)
	cacheModule[CTL_M_HOST].set(h.Id, h)
	return h.Id, err
}

func GetHost(id int) (*Host, error) {
	if h, ok := cacheModule[CTL_M_HOST].get(id).(*Host); ok {
		return h, nil
	}
	h := &Host{Id: id}
	err := orm.NewOrm().Read(h, "Id")
	if err == nil {
		cacheModule[CTL_M_HOST].set(id, h)
	}
	return h, err
}

func GetHostByUuid(uuid string) (h *Host, err error) {
	h = &Host{Uuid: uuid}
	err = orm.NewOrm().Read(h, "Uuid")
	return h, err
}

func QueryHosts(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Host))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query))
	}
	return qs
}

func GetHostsCnt(query string) (int, error) {
	cnt, err := QueryHosts(query).Count()
	return int(cnt), err
}

func GetHosts(query string, limit, offset int) (hosts []*Host, err error) {
	_, err = QueryHosts(query).Limit(limit, offset).All(&hosts)
	return
}

func UpdateHost(id int, _h *Host) (h *Host, err error) {
	if h, err = GetHost(id); err != nil {
		return nil, ErrNoHost
	}

	if _h.Uuid != "" {
		h.Uuid = _h.Uuid
	}
	if _h.Name != "" {
		h.Name = _h.Name
	}
	if _h.Type != "" {
		h.Type = _h.Type
	}
	if _h.Type != "" {
		h.Type = _h.Type
	}
	if _h.Status != "" {
		h.Status = _h.Status
	}
	if _h.Loc != "" {
		h.Loc = _h.Loc
	}
	if _h.Idc != "" {
		h.Idc = _h.Idc
	}
	_, err = orm.NewOrm().Update(h)
	return h, err
}

func DeleteHost(id int) error {

	if n, err := orm.NewOrm().Delete(&Host{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	cacheModule[CTL_M_HOST].del(id)

	return nil
}
