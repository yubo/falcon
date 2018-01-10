/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon/alarm"
	"github.com/yubo/falcon/service/expr"
)

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

func (p *EventTrigger) Exec(item *itemEntry) *alarm.Event {
	glog.V(4).Infof("exec %s expr %s", item.key, p.Expr)

	if !expr.Exec(item, p.expr) {
		return nil
	}

	item.RLock()
	defer item.RUnlock()
	id := (item.dataId - 1) & CACHE_SIZE_MASK
	return &alarm.Event{
		TagId:     p.TagId,
		Key:       []byte(item.key),
		Expr:      []byte(p.Expr),
		Msg:       []byte(p.Msg),
		Timestamp: item.timestamp[id],
		Value:     item.value[id],
	}

}

func (p *EventTrigger) exprPrepare() (err error) {
	p.expr, err = expr.Parse(p.Expr)
	return
}
