/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "github.com/astaxie/beego/orm"

type Strategy struct {
	Id        int64  `json:"id"`
	Metric_id int64  `json:"metricId"`
	Tags      string `json:"tags"`
	Max_step  int    `json:"maxStep"`
	Priority  int    `json:"priority"`
	Func      string `json:"fun"`
	Op        string `json:"op"`
	Condition string `json:"condition"`
	Note      string `json:"note"`
	Metric    string `json:"metric"`
	RunBegin  string `json:"runBegin"`
	RunEnd    string `json:"runEnd"`
	TplId     int    `json:"tplId"`
}

func (u *User) AddStrategy(s *Strategy) (id int64, err error) {
	if id, err = orm.NewOrm().Insert(s); err != nil {
		return
	}
	s.Id = id
	return
}

func (u *User) GetStrategy(id int64) (s *Strategy, err error) {
	s = &Strategy{Id: id}
	err = orm.NewOrm().Read(s, "Id")
	return
}

func (u *User) QueryStrategys(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Strategy))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetStrategysCnt(query string) (int64, error) {
	return u.QueryStrategys(query).Count()
}

func (u *User) GetStrategys(query string, limit, offset int) (strategys []*Strategy, err error) {
	_, err = u.QueryStrategys(query).Limit(limit, offset).All(&strategys)
	return
}

func (u *User) UpdateStrategy(id int64, _o *Strategy) (o *Strategy, err error) {
	if o, err = u.GetStrategy(id); err != nil {
		return nil, ErrNoStrategy
	}

	o.Metric_id = _o.Metric_id
	o.Tags = _o.Tags
	o.Max_step = _o.Max_step
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
