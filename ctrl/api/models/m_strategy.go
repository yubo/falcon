/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "github.com/astaxie/beego/orm"

type Strategy struct {
	Id        int64  `json:"id"`
	MetricId  int64  `json:"metricId"`
	Tags      string `json:"tags"`
	MaxStep   int    `json:"maxStep"`
	Priority  int    `json:"priority"`
	Func      string `json:"fun"`
	Op        string `json:"op"`
	Condition string `json:"condition"`
	Note      string `json:"note"`
	Metric    string `json:"metric"`
	RunBegin  string `json:"runBegin"`
	RunEnd    string `json:"runEnd"`
	TplId     int64  `json:"tplId"`
}

func (u *User) AddStrategy(o *Strategy) (id int64, err error) {
	o.Id = 0
	if id, err = orm.NewOrm().Insert(o); err != nil {
		return
	}
	o.Id = id
	return
}

func (u *User) GetStrategy(id int64) (s *Strategy, err error) {
	s = &Strategy{Id: id}
	err = orm.NewOrm().Read(s, "Id")
	return
}

func (u *User) QueryStrategys(tid int64, query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Strategy))
	if tid != 0 {
		qs = qs.Filter("TplId", tid)
	}
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetStrategysCnt(tid int64, query string) (int64, error) {
	return u.QueryStrategys(tid, query).Count()
}

func (u *User) GetStrategys(tid int64, query string, limit, offset int) (strategys []*Strategy, err error) {
	_, err = u.QueryStrategys(tid, query).Limit(limit, offset).All(&strategys)
	return
}

func (u *User) UpdateStrategy(id int64, _o *Strategy) (o *Strategy, err error) {
	if o, err = u.GetStrategy(id); err != nil {
		return nil, ErrNoStrategy
	}

	o.MetricId = _o.MetricId
	o.Tags = _o.Tags
	o.MaxStep = _o.MaxStep
	o.Priority = _o.Priority
	o.Func = _o.Func
	o.Op = _o.Op
	o.Condition = _o.Condition
	o.Note = _o.Note
	o.Metric = _o.Metric
	o.RunBegin = _o.RunBegin
	o.RunEnd = _o.RunEnd
	o.TplId = _o.TplId

	_, err = orm.NewOrm().Update(o)
	return o, err
}

func (u *User) DeleteStrategy(id int64) error {

	if n, err := orm.NewOrm().Delete(&Strategy{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}

	return nil
}
