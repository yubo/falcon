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

func (u *User) AddHost(h *Host, tagstring string) (_id int, err error) {
	var id int64
	var t *Tag

	if tagstring, err = sysTagSchema.Fmt(t.Name,
		false); err != nil {
		return
	}

	if t, err = u.Access(SYS_W_SCOPE,
		TagParent(tagstring)); err != nil {
		return
	}

	id, err = orm.NewOrm().Insert(h)
	if err != nil {
		beego.Error(err)
		return
	}

	// TODO: bind host to t.Name
	_, err = orm.NewOrm().Raw("insert into tag_host(tag_id, host_id) values (?, ?)", id, t.Id).Exec()
	if err != nil {
		return
	}

	h.Id = int(id)
	cacheModule[CTL_M_HOST].set(h.Id, h)
	DbLog(u.Id, CTL_M_HOST, h.Id, CTL_A_ADD, "")

	return h.Id, nil
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

func (u *User) QueryHosts(query string) orm.QuerySeter {
	// TODO: acl filter
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

func (u *User) UpdateHost(id int, _h *Host) (h *Host, err error) {
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
	DbLog(u.Id, CTL_M_HOST, h.Id, CTL_A_SET, "")
	return h, err
}

func (u *User) DeleteHost(host_id int, tag string) (err error) {
	var n, m int64
	var t *Tag

	qs := orm.NewOrm().QueryTable("tag_host").Filter("host_id", host_id)
	n, err = qs.Count()
	if err != nil {
		// not exist
		return
	}

	// ACL
	if t, err = u.Access(SYS_W_SCOPE, tag); err != nil {
		return
	}

	// unbind tag_host
	if m, err = qs.Filter("tag_id", t.Id).Delete(); err != nil {
		return
	}

	// delete host_id if not bound by any tag
	if n-m == 0 {
		if n, err = orm.NewOrm().Delete(&Host{Id: host_id}); err != nil {
			return
		}
		cacheModule[CTL_M_HOST].del(host_id)
		DbLog(u.Id, CTL_M_HOST, host_id, CTL_A_DEL, "")
	}

	return nil
}
