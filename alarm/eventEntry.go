/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/yubo/falcon/alarm/expr"
	"github.com/yubo/gotool/list"
)

type eventEntry struct {
	sync.RWMutex
	expr.ExprEvent
	list      list.ListHead // point to newQueue or lruQueue
	lastTs    int64
	tagId     int64
	key       string
	expr      string
	msg       string
	timestamp int64
	value     float64
	priority  int
}

func (p *eventEntry) Index(s, substr string) int {
	switch s {
	case "key":
		s = p.key
	}
	glog.V(3).Infof("%s index(%s, %s) = %d", MODULE_NAME, s, substr, strings.Index(s, substr))
	return strings.Index(s, substr)
}

func (p *eventEntry) Value(key string) int {
	switch key {
	case "priority":
		return p.priority
	}
	return -1
}
