/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yubo/falcon/lib/core"
)

type Host struct {
	Id     int64  `json:"id"`
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Type   string `json:"typ"`
	Status string `json:"status"`
	Loc    string `json:"loc"`
	Idc    string `json:"idc"`
	Pause  int64  `json:"pause"`
}

type HostApiAdd struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Type   string `json:"typ"`
	Status string `json:"status"`
	Loc    string `json:"loc"`
	Idc    string `json:"idc"`
	Pause  int64  `json:"pause"`
}

type HostApiUpdate struct {
	Id     int64  `json:"id"`
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Type   string `json:"typ"`
	Status string `json:"status"`
	Loc    string `json:"loc"`
	Idc    string `json:"idc"`
	Pause  int64  `json:"pause"`
}

func (op *Operator) CreateHost(h *HostApiAdd) (id int64, err error) {
	id, err = op.SqlInsert("insert host (uuid, name, type, status, loc, idc, pause) values (?, ?, ?, ?, ?, ?, ?)",
		h.Uuid, h.Name, h.Type, h.Status, h.Loc, h.Idc, h.Pause)
	if err != nil {
		return
	}

	op.log(CTL_M_HOST, id, CTL_A_ADD, jsonStr(h))
	return
}

func (op *Operator) GetHost(id int64) (h *Host, err error) {
	if h, ok := moduleCache[CTL_M_HOST].get(id).(*Host); ok {
		return h, nil
	}

	h = &Host{}
	err = op.SqlRow(h, "select id, uuid, name, type, status, loc, idc, pause from host where id = ?", id)
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
	sql = sqlLimit("select id, uuid, name, type, status, loc, idc, pause from host "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) UpdateHost(h *Host) (ret *Host, err error) {
	_, err = op.SqlExec("update host set uuid = ?, name = ?, type = ?, status = ?, loc = ?, idc = ?, pause = ? where id = ?", h.Uuid, h.Name, h.Type, h.Status, h.Loc, h.Idc, h.Pause, h.Id)
	if err != nil {
		return
	}

	moduleCache[CTL_M_HOST].del(h.Id, h.Name)
	if ret, err = op.GetHost(h.Id); err != nil {
		return
	}

	op.log(CTL_M_HOST, h.Id, CTL_A_SET, jsonStr(h))
	return ret, err
}

func (op *Operator) DeleteHost(id int64) error {
	if err := op.RelCheck("SELECT count(*) FROM tag_host where host_id = ?", id); err != nil {
		return errors.New(err.Error() + "(tag - host)")
	}

	if _, err := op.SqlExec("DELETE FROM host WHERE id = ?", id); err != nil {
		return err
	}
	if h, ok := moduleCache[CTL_M_HOST].get(id).(*Host); ok {
		moduleCache[CTL_M_HOST].del(id, h.Name)
	}
	op.log(CTL_M_HOST, id, CTL_A_DEL, "")

	return nil
}

/*******************************************************************************
 ************************ tag host *********************************************
 ******************************************************************************/

type TagHostApiGet struct {
	Id            int64  `json:"id"`
	TagId         int64  `json:"tag_id"`
	HostId        int64  `json:"host_id"`
	TagName       string `json:"tag_name"`
	HostName      string `json:"host_name"`
	Pause         int64  `json:"pause"`
	MaintainBegin int64  `json:"maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"`
}

type TagHostApiAdd struct {
	SrcTagId int64 `json:"src_tag_id"`
	TagId    int64 `json:"tag_id"`
	HostId   int64 `json:"host_id"`
}

type TagHostsApiAdd struct {
	SrcTagId int64   `json:"src_tag_id"`
	TagId    int64   `json:"tag_id"`
	HostIds  []int64 `json:"host_ids"`
}

type TagHostApiDel struct {
	TagId  int64 `json:"tag_id"`
	HostId int64 `json:"host_id"`
}

type TagHostsApiDel struct {
	TagId   int64   `json:"tag_id"`
	HostIds []int64 `json:"host_ids"`
}

func (op *Operator) ChkTagHostCnt(tagId int64, host_ids []int64) (cnt int64, err error) {
	err = op.O.Raw("SELECT count(*) FROM tag_host where tag_id = ? and host_id IN "+
		array2sql(host_ids), tagId).QueryRow(&cnt)
	return
}

func tagHostSql(tagId int64, query string, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "h.name like ?")
		sql3 = append(sql3, "%"+query+"%")
	}
	if deep {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tagId))
	} else {
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tagId)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetTagHostCnt(tagId int64, query string, deep bool) (cnt int64, err error) {
	//TODO: acl filter just for admin?
	sql, sql_args := tagHostSql(tagId, query, deep)
	err = op.O.Raw("SELECT count(*) FROM tag_host a left join host h on a.host_id = h.id left join tag t on a.tag_id = t.id "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTagHost(tagId int64, query string, deep bool,
	limit, offset int) (ret []*TagHostApiGet, err error) {

	sql, sql_args := tagHostSql(tagId, query, deep)
	sql = "select a.id, a.tag_id as tag_id, a.host_id as host_id, h.name as host_name, h.pause as pause, t.name as tag_name from tag_host a left join host h on a.host_id = h.id left join tag t on a.tag_id = t.id " + sql + " ORDER BY h.name"
	if limit > 0 {
		sql += " LIMIT ? OFFSET ?"
		sql_args = append(sql_args, limit, offset)
	}
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) CreateTagHost(rel *TagHostApiAdd) (id int64, err error) {
	id, err = op.SqlInsert("insert tag_host (tag_id, host_id) values (?, ?)",
		rel.TagId, rel.HostId)
	if err == nil {
		op.log(CTL_M_TAG_HOST, id, CTL_A_ADD, jsonStr(rel))
	}
	return
}

func (op *Operator) CreateTagHosts(rel *TagHostsApiAdd) (int64, error) {
	vs := make([]string, len(rel.HostIds))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", rel.TagId, rel.HostIds[i])
	}

	res, err := op.O.Raw("INSERT `tag_host` (`tag_id`, " +
		"`host_id`) VALUES " + strings.Join(vs, ", ")).Exec()
	if err != nil {
		return 0, err
	}
	op.log(CTL_M_TAG_HOST, 0, CTL_A_ADD, strings.Join(vs, ", "))
	return res.RowsAffected()
}

func (op *Operator) DeleteTagHost(rel *TagHostApiDel) (int64, error) {
	res, err := op.O.Raw("DELETE FROM `tag_host` WHERE tag_id = ? and host_id = ?",
		rel.TagId, rel.HostId).Exec()
	if err != nil {
		return 0, err
	}
	op.log(CTL_M_TAG_HOST, rel.TagId, CTL_A_DEL, jsonStr(rel))
	return res.RowsAffected()
}

func (op *Operator) DeleteTagHosts(rel *TagHostsApiDel) (int64, error) {
	if len(rel.HostIds) == 0 {
		return 0, core.ErrEmpty
	}
	ids := array2sql(rel.HostIds)
	res, err := op.O.Raw("DELETE FROM `tag_host` "+
		"WHERE tag_id = ? and host_id  IN "+ids, rel.TagId).Exec()
	if err != nil {
		return 0, err
	}
	op.log(CTL_M_TAG_HOST, 0, CTL_A_DEL, ids)
	return res.RowsAffected()
}
