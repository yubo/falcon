/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon/service/expr"
)

type event struct {
	trigger   *EventTrigger
	timestamp int64
}

type EventTrigger struct {
	Id       int64
	ParentId int64
	TagId    int64
	Priority int
	Name     string
	Metric   string
	Tags     string
	Expr     string
	Msg      string
	Child    []*EventTrigger
	items    []*itemEntry
	expr     *expr.Expr
}

func (p *EventTrigger) Exec(item *itemEntry) *event {
	glog.V(4).Infof("exec endpoint %s metric %s expr %s",
		string(item.endpoint), string(item.metric), p.Expr)
	if expr.Exec(item, p.expr) {
		return &event{trigger: p, timestamp: timer.now()}
	}

	return nil
}

func (p *EventTrigger) exprPrepare() (err error) {
	p.expr, err = expr.Parse(p.Expr)
	return
}
