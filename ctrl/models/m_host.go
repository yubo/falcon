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
	Id          int64     `json:"id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Loc         string    `json:"loc"`
	Idc         string    `json:"idc"`
	Create_time time.Time `json:"ctime"`
}

func (u *User) AddHost(h *Host) (id int64, err error) {
	id, err = orm.NewOrm().Insert(h)
	if err != nil {
		beego.Error(err)
		return
	}

	h.Id = id
	cacheModule[CTL_M_HOST].set(id, h)
	DbLog(u.Id, CTL_M_HOST, id, CTL_A_ADD, jsonStr(h))
	return
}

func (u *User) GetHost(id int64) (*Host, error) {
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

func (u *User) GetHostByUuid(uuid string) (h *Host, err error) {
	h = &Host{Uuid: uuid}
	err = orm.NewOrm().Read(h, "Uuid")
	return h, err
}

func (u *User) QueryHosts(query string) orm.QuerySeter {
	// TODO: acl filter
	// just for admin?
	qs := orm.NewOrm().QueryTable(new(Host))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetHostsCnt(query string) (int64, error) {
	return u.QueryHosts(query).Count()
}

func (u *User) GetHosts(query string, limit, offset int) (hosts []*Host, err error) {
	_, err = u.QueryHosts(query).Limit(limit, offset).All(&hosts)
	return
}

func (u *User) UpdateHost(id int64, _h *Host) (h *Host, err error) {
	if h, err = u.GetHost(id); err != nil {
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
	cacheModule[CTL_M_HOST].set(id, h)
	DbLog(u.Id, CTL_M_HOST, id, CTL_A_SET, "")
	return h, err
}

func (u *User) DeleteHost(id int64) error {
	if n, err := orm.NewOrm().Delete(&Host{Id: id}); err != nil || n == 0 {
		return err
	}
	cacheModule[CTL_M_HOST].del(id)
	DbLog(u.Id, CTL_M_HOST, id, CTL_A_DEL, "")

	return nil
}
