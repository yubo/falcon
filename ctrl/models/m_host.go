/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
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
	Create_time time.Time `json:"-"`
}

func (u *User) AddHost(h *Host, tag string) (id int64, err error) {
	var t *Tag

	if tag, err = sysTagSchema.Fmt(tag,
		false); err != nil {
		return
	}

	if t, err = u.Access(SYS_W_SCOPE,
		tag, true); err != nil {
		return
	}

	id, err = orm.NewOrm().Insert(h)
	if err != nil {
		beego.Error(err)
		return
	}

	// bind host to tag
	_, err = orm.NewOrm().Raw("insert into tag_host(tag_id, host_id) values (?, ?)", t.Id, id).Exec()
	if err != nil {
		return
	}

	h.Id = id
	cacheModule[CTL_M_HOST].set(id, h)
	data, _ := json.Marshal(h)
	DbLog(u.Id, CTL_M_HOST, id, CTL_A_ADD, data)

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

func (u *User) GetHostsCnt(query string) (int, error) {
	cnt, err := u.QueryHosts(query).Count()
	return int(cnt), err
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
	DbLog(u.Id, CTL_M_HOST, id, CTL_A_SET, nil)
	return h, err
}

func (u *User) DeleteHost(id int64, tag string) (err error) {
	var n, m int64
	var t *Tag

	qs := orm.NewOrm().QueryTable("tag_host").Filter("host_id", id)
	n, err = qs.Count()
	if err != nil {
		// not exist
		return
	}

	// ACL
	if t, err = u.Access(SYS_W_SCOPE, tag, true); err != nil {
		return
	}

	// unbind tag_host
	if m, err = qs.Filter("tag_id", t.Id).Delete(); err != nil {
		return
	}

	// delete host_id if not bound by any tag
	if n-m == 0 {
		if n, err = orm.NewOrm().Delete(&Host{Id: id}); err != nil {
			return
		}
		cacheModule[CTL_M_HOST].del(id)
		DbLog(u.Id, CTL_M_HOST, id, CTL_A_DEL, nil)
	}

	return nil
}
