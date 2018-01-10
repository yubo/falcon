/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
https://www.zabbix.com/documentation/4.0/manual/appendix/triggers/functions
https://www.zabbix.com/documentation/4.0/manual/config/triggers/expression
*/
package expr

import (
	"strings"
	"testing"

	"github.com/golang/glog"
)

func TestParse(t *testing.T) {
	cases := []struct {
		text string
		want error
	}{
		{"index(key, cop=xiaomi) >= 0", nil},
		{"value(priority) >= 1", nil},
		{"test", EINVAL},
	}
	for _, tc := range cases {
		expr, err := Parse(tc.text)
		if tc.want == nil && err != nil {
			t.Errorf("Parse(%s) = %s %v; want %v", tc.text, expr, err, tc.want)
		} else if tc.want != nil && err == nil {
			t.Errorf("Parse(%s) = %s %v; want %v", tc.text, expr, err, tc.want)
		}
	}
}

type testEvent struct {
	ExprEvent
	key      string
	priority int
}

func (p *testEvent) Index(s, substr string) int {
	switch s {
	case "key":
		s = p.key
	}
	glog.V(3).Infof("index(%s, %s) = %d", s, substr, strings.Index(s, substr))
	return strings.Index(s, substr)
}

func (p *testEvent) Value(key string) int {
	switch key {
	case "priority":
		return p.priority
	}
	return -1
}

func TestExec(t *testing.T) {
	cases := []struct {
		key      string
		priority int
		expr     string
		want     bool
	}{
		{"cop=xiaomi,owt=miliao,pdl=op", 1, "index(key, owt=miliao) >= 0", true},
		{"cop=xiaomi,owt=miliao,pdl=op", 1, "index(key, owt=inf) >= 0", false},
		{"cop=xiaomi,owt=miliao,pdl=op", 1, "value(priority) = 1", true},
		{"cop=xiaomi,owt=miliao,pdl=op", 1, "value(priority) > 1", false},
		{"cop=xiaomi,owt=miliao,pdl=op", 1, "value(priority) = 1 && index(key, pdl=op) >= 0", true},
	}
	for _, tc := range cases {
		expr, err := Parse(tc.expr)
		if err != nil {
			t.Errorf("Parse(%s) err %v", tc.expr, err)
		}
		event := &testEvent{key: tc.key, priority: tc.priority}
		if got := Exec(event, expr); got != tc.want {
			t.Errorf("%s = %v; want %v", tc.expr, got, tc.want)
		}
	}
}
