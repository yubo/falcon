/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"sort"
	"strings"
	"time"

	"github.com/yubo/falcon"
)

const TMP_GRAPH_FILED_DELIMITER = "|"

////////////////// tmp graph

type APITmpGraph struct {
	Endpoints []string `json:"endpoints"`
	Counters  []string `json:"counters"`
}

type TmpGraph struct {
	Id         int64     `json:"id"`
	Endpoints  string    `json:"endpoints"`
	Counters   string    `json:"counters"`
	Ck         string    `json:"ck"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) AddDashboardTmpGraph(inputs *APITmpGraph) (id int64, err error) {
	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)

	es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
	ck := falcon.Md5sum(es_string + ":" + cs_string)

	res, err := op.O.Raw("insert ignore into `tmp_graph` (endpoints, counters, ck) values(?, ?, ?) on duplicate key update time_=?", es_string, cs_string, ck, time.Now()).Exec()
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()

	DbLog(op.O, op.User.Id, CTL_M_TMP_GRAPH, id, CTL_A_ADD, jsonStr(inputs))
	return id, nil
}

func (op *Operator) GetDashboardTmpGraph(id int64) (*APITmpGraph, error) {
	o := &TmpGraph{}

	err := op.SqlRow(o, "select id, endpoints, counters, ck from tmp_graph where id = ?", id)
	if err != nil {
		return nil, err
	}

	return &APITmpGraph{
		Endpoints: strings.Split(o.Endpoints,
			TMP_GRAPH_FILED_DELIMITER),
		Counters: strings.Split(o.Counters,
			TMP_GRAPH_FILED_DELIMITER),
	}, nil
}

////////////////// graph

type APIGraph struct {
	ScreenId   int64    `json:"screen_id"`
	GraphId    int64    `json:"graph_id"`
	Title      string   `json:"title"`
	Endpoints  []string `json:"endpoints"`
	Counters   []string `json:"counters"`
	TimeSpan   int      `json:"timespan"`
	GraphType  string   `json:"graph_type"`
	Method     string   `json:"method"`
	Position   int      `json:"position"`
	FalconTags string   `json:"falcon_tags"`
}

type DashboardGraph struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Hosts      string `json:"hosts"`
	Counters   string `json:"counters"`
	ScreenId   int64  `json:"screen_id"`
	Timespan   int    `json:"timespan"`
	GraphType  string `json:"graph_type"`
	Method     string `json:"method"`
	Position   int    `json:"position"`
	FalconTags string `json:"falcon_tags"`
}

func (op *Operator) AddDashboardGraph(in *APIGraph) (id int64, err error) {

	es := in.Endpoints
	cs := in.Counters
	sort.Strings(es)
	sort.Strings(cs)
	es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)

	if in.TimeSpan == 0 {
		in.TimeSpan = 3600
	}
	if in.GraphType == "" {
		in.GraphType = "h"
	}

	id, err = op.SqlInsert("insert dashboard_graph (title, hosts, counters, screen_id, timespan, graph_type, method, position) values (?, ?, ?, ?, ?, ?, ?, ?)", in.Title, es_string, cs_string, int64(in.ScreenId), in.TimeSpan, in.GraphType, in.Method, in.Position)
	if err != nil {
		return
	}

	DbLog(op.O, op.User.Id, CTL_M_DASHBOARD_GRAPH, id, CTL_A_ADD, jsonStr(in))
	return
}

func (op *Operator) UpdateDashboardGraph(inputs *APIGraph) (ret *DashboardGraph, err error) {

	ret, err = op.GetDashboardGraph(inputs.GraphId)
	if err != nil {
		return
	}

	if inputs.Title != "" {
		ret.Title = inputs.Title
	}
	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		ret.Hosts = es_string
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		ret.Counters = cs_string
	}
	if inputs.ScreenId != 0 {
		ret.ScreenId = int64(inputs.ScreenId)
	}
	if inputs.TimeSpan != 0 {
		ret.Timespan = inputs.TimeSpan
	}
	if inputs.GraphType != "" {
		ret.GraphType = inputs.GraphType
	}
	if inputs.Method != "" {
		ret.Method = inputs.Method
	}
	if inputs.Position != 0 {
		ret.Position = inputs.Position
	}
	if inputs.FalconTags != "" {
		ret.FalconTags = inputs.FalconTags
	}

	_, err = op.SqlExec("update dashboard_graph set title = ?, hosts = ?, counters = ?, screen_id = ?, timespan = ?, graph_type = ?, method = ?, position = ?, falcon_tags = ? where id = ?", ret.Title, ret.Hosts, ret.Counters, ret.ScreenId, ret.Timespan, ret.GraphType, ret.Method, ret.Position, ret.FalconTags, ret.Id)

	DbLog(op.O, op.User.Id, CTL_M_DASHBOARD_GRAPH, ret.Id, CTL_A_SET, jsonStr(ret))
	return
}

func (op *Operator) GetDashboardGraph(id int64) (ret *DashboardGraph, err error) {
	ret = &DashboardGraph{}
	err = op.SqlRow(ret, "select id, title, hosts, counters, screen_id, timespan, graph_type, method, position, falcon_tags from dashboard_graph where id = ?", id)
	return
}

func (op *Operator) DeleteDashboardGraph(id int64) (err error) {
	if _, err = op.SqlExec("delete from dashboard_graph where id = ?", id); err != nil {
		return err
	}
	DbLog(op.O, op.User.Id, CTL_M_DASHBOARD_GRAPH, id, CTL_A_DEL, "")
	return nil
}

func (op *Operator) GetDashboardGraphByScreen(id int64) (ret []*APIGraph, err error) {
	var graphs []*DashboardGraph

	_, err = op.O.Raw("select id, title, hosts, counters, screen_id, timespan, graph_type, method, position, falcon_tags from dashboard_graph where screen_id = ? LIMIT ? OFFSET ?", id, 500, 0).QueryRows(&graphs)
	if err != nil {
		return
	}

	ret = make([]*APIGraph, len(graphs))
	for i, graph := range graphs {
		ret[i] = &APIGraph{
			Title:      graph.Title,
			GraphId:    graph.Id,
			Endpoints:  strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER),
			Counters:   strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER),
			ScreenId:   graph.ScreenId,
			GraphType:  graph.GraphType,
			TimeSpan:   graph.Timespan,
			Method:     graph.Method,
			Position:   graph.Position,
			FalconTags: graph.FalconTags,
		}
	}
	return
}

/////////////////////// screen

type AddDashboardScreen struct {
	Id   int64  `json:"-"`
	Pid  int64  `json:"pid"`
	Name string `json:"name"`
}

type DashboardScreen struct {
	Id   int64  `json:"id"`
	Pid  int64  `json:"pid"`
	Name string `json:"name"`
}

func (op *Operator) AddDashboardScreen(in *DashboardScreen) (id int64, err error) {

	id, err = op.SqlInsert("insert dashboard_screen (pid, name ) values (?, ?)", in.Pid, in.Name)
	if err != nil {
		return
	}
	DbLog(op.O, op.User.Id, CTL_M_DASHBOARD_SCREEN, id, CTL_A_ADD, jsonStr(in))
	return
}

func (op *Operator) GetDashboardScreen(id, limit int64) (ret []*DashboardScreen, err error) {
	if id != 0 {
		_, err = op.O.Raw("select id, pid, name from dashboard_screen where id = ? ", id).QueryRows(&ret)
	} else {
		_, err = op.O.Raw("select id, pid, name from dashboard_screen limit ? order by id", limit).QueryRows(&ret)
	}
	return

}

func (op *Operator) GetDashboardScreenByPid(pid int64) (ret []*DashboardScreen, err error) {
	_, err = op.O.Raw("select id, pid, name from dashboard_screen where pid = ? ", pid).QueryRows(&ret)
	return
}

func (op *Operator) DeleteDashboardScreen(id int64) (err error) {

	if _, err = op.SqlExec("delete from dashboard_screen where id = ?", id); err != nil {
		return err
	}
	DbLog(op.O, op.User.Id, CTL_M_DASHBOARD_SCREEN, id, CTL_A_DEL, "")
	return
}

func (op *Operator) UpdateDashboardScreen(id int64, inputs *DashboardScreen) (d *DashboardScreen, err error) {
	d = &DashboardScreen{}

	err = op.O.Raw("select id, pid, name from dashboard_screen where id = ? ", id).QueryRow(d)
	if err != nil {
		return
	}

	if inputs.Name != "" {
		d.Name = inputs.Name
	}

	if inputs.Pid != 0 {
		d.Pid = inputs.Pid
	}

	_, err = op.SqlExec("update dashboard_screen set name = ?, pid = ? where id = ?", d.Name, d.Pid, id)
	return
}
