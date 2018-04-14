/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"strings"
)

type PluginDirPost struct {
	Id           int64  `json:"-"`
	TagId        int64  `json:"tag_id"`
	Dir          string `json:"dir"`
	CreateUserId int64  `json:"-"`
}

type PluginDirGet struct {
	Id      int64  `json:"id"`
	TagId   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
	Dir     string `json:"dir"`
	Creator string `json:"creator"`
}

type PluginDir struct {
	Id           int64  `json:"id"`
	TagId        int64  `json:"tag_id"`
	Dir          string `json:"dir"`
	CreateUserId int64  `json:"create_user_id"`
}

func sqlTagPlugin(tag_id int64, deep bool) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if deep {
		sql2 = append(sql2, fmt.Sprintf("a.tag_id in (select tag_id from tag_rel where sup_tag_id = %d)", tag_id))
	} else {
		sql2 = append(sql2, "a.tag_id = ?")
		sql3 = append(sql3, tag_id)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetPluginDirCnt(tag_id int64, deep bool) (cnt int64, err error) {
	sql, sql_args := sqlTagPlugin(tag_id, deep)
	err = op.O.Raw("SELECT count(*) FROM plugin_dir a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetPluginDir(tag_id int64, deep bool, limit, offset int) (ret []*PluginDirGet, err error) {
	sql, sql_args := sqlTagPlugin(tag_id, deep)

	sql = "SELECT a.id, a.tag_id, a.dir, u.name as creator, t.name as tag_name FROM plugin_dir a left join user u on u.id = a.create_user_id left join tag t on t.id = a.tag_id " + sql + " LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) CreatePluginDir(input *PluginDirPost) (id int64, err error) {
	return op.SqlInsert("insert plugin_dir (dir, tag_id, create_user_id) values (?, ?, ?)", input.Dir, input.TagId, op.User.Id)
}

func (op *Operator) DeletePluginDir(tag_id, id int64) (int64, error) {
	return op.SqlExec("delete from plugin_dir where tag_id = ? and id = ?", tag_id, id)
}
