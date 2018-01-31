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
	entries  []*cacheEntry
	expr     *expr.Expr
}

func (p *EventTrigger) Dispatch(e *cacheEntry) *alarm.Event {
	glog.V(4).Infof("%s dispatch %s expr %s", MODULE_NAME, e.key, p.Expr)

	if !expr.Exec(e, p.expr) {
		return nil
	}

	e.RLock()
	defer e.RUnlock()
	v := e.values[(e.dataId-1)&CACHE_DATA_SIZE_MASK]
	return &alarm.Event{
		TagId:     p.TagId,
		Key:       e.key.Key,
		Expr:      []byte(p.Expr),
		Msg:       []byte(p.Msg),
		Timestamp: v.Timestamp,
		Value:     v.Value,
	}

}

func (p *EventTrigger) exprPrepare() (err error) {
	p.expr, err = expr.Parse(p.Expr)
	return
}
