/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"time"
)

type Host struct {
	Id            int64     `json:"id"`
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name"`
	Type          string    `json:"typ"`
	Status        string    `json:"status"`
	Loc           string    `json:"loc"`
	Idc           string    `json:"idc"`
	Pause         int64     `json:"pause"`
	MaintainBegin int64     `json:"maintain_begin"`
	MaintainEnd   int64     `json:"maintain_end"`
	CreateTime    time.Time `json:"ctime"`
}

type HostUpdate struct {
	Id            int64  `json:"id"`
	Uuid          string `json:"uuid"`
	Name          string `json:"name"`
	Type          string `json:"typ"`
	Status        string `json:"status"`
	Loc           string `json:"loc"`
	Idc           string `json:"idc"`
	Pause         int64  `json:"pause"`
	MaintainBegin int64  `json:"maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"`
}

func (op *Operator) AddHost(h *Host) (id int64, err error) {

	id, err = op.SqlInsert("insert host (uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", h.Uuid, h.Name, h.Type, h.Status, h.Loc, h.Idc, h.Pause, h.MaintainBegin, h.MaintainEnd)
	if err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_HOST, id, CTL_A_ADD, jsonStr(h))
	return
}

func (op *Operator) GetHost(id int64) (h *Host, err error) {
	if h, ok := moduleCache[CTL_M_HOST].get(id).(*Host); ok {
		return h, nil
	}

	h = &Host{}
	err = op.SqlRow(h, "select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host where id = ?", id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_HOST].set(id, h, h.Name)
	return
}

func (op *Operator) GetHostsCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlName(query)
	err = op.O.Raw("SELECT count(*) FROM host "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetHosts(query string, limit, offset int) (ret []*Host, err error) {
	sql, sql_args := sqlName(query)
	sql = sqlLimit("select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateHost(h *Host) (ret *Host, err error) {
	_, err = op.SqlExec("update host set uuid = ?, name = ?, type = ?, status = ?, loc = ?, idc = ?, pause = ?, maintain_begin = ?, maintain_end = ? where id = ?", h.Uuid, h.Name, h.Type, h.Status, h.Loc, h.Idc, h.Pause, h.MaintainBegin, h.MaintainEnd, h.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_HOST].del(h.Id, h.Name)
	if ret, err = op.GetHost(h.Id); err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_HOST, h.Id, CTL_A_SET, jsonStr(h))
	return ret, err
}

func (op *Operator) DeleteHost(id int64) error {
	if err := op.RelCheck("SELECT count(*) FROM tag_host where host_id = ?", id); err != nil {
		return errors.New(err.Error() + "(tag - host)")
	}

	if _, err := op.SqlExec("delete from host where id = ?", id); err != nil {
		return err
	}
	if h, ok := moduleCache[CTL_M_HOST].get(id).(*Host); ok {
		moduleCache[CTL_M_HOST].del(id, h.Name)
	}
	DbLog(op.O, op.User.Id, CTL_M_HOST, id, CTL_A_DEL, "")

	return nil
}
