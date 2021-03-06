/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon/modules/alarm/expr"
)

type ActionTrigger struct {
	Id           int64
	TagId        int64
	TokenId      int64
	OrderId      int64
	Expr         string
	ActionFlag   uint64
	ActionScript string
	events       map[string]*Event
	expr         *expr.Expr
}

func (p *ActionTrigger) Dispatch(item *eventEntry) bool {
	glog.V(5).Infof("dispatch %s expr %s", item.key, p.Expr)
	return expr.Exec(item, p.expr)
}

func (p *ActionTrigger) exprPrepare() (err error) {
	p.expr, err = expr.Parse(p.Expr)
	return
}
