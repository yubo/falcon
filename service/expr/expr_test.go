/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

import (
	"regexp"
	"testing"
)

/*
min(#3) < 4
min(360) < 4
*/

func TestRegexp(t *testing.T) {
	cases := []struct {
		name    string
		pattern *regexp.Regexp
		text    string
		want    string
	}{
		{"f_num", f_num, "12.34.3 1234", "12.34"},
		{"f_num", f_num, "12.3.3 1234", "12.3"},
		{"f_num", f_num, "12 123 4 2", "12"},
	}
	for _, tc := range cases {
		if f := tc.pattern.Find([]byte(tc.text)); f == nil || string(f[:]) != tc.want {
			t.Errorf("%s.find(%s) = '%v'; want '%v'", tc.name, tc.pattern, string(f), tc.want)
		}
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		text string
		want error
	}{
		{"min(#1, 3) < 3", nil},
		{"(min(#1, 3) < 3) && (min(3600, 4) < 4)", nil},
	}
	for _, tc := range cases {
		if expr, err := Parse([]byte(tc.text), 0); err != tc.want {
			t.Errorf("Parse(%s) = %s %v; want %v", tc.text, expr, err, tc.want)
		} /* else {
			t.Logf("Parse(%s) = %s %v", tc.text, expr, err)
			t.Logf("json %s", expr.Json())
		}*/
	}
}

type testItem struct {
	ExprItem
	dps []dataPoint
}

type dataPoint struct {
	ts    int64
	value float64
}

func dpsToArr(in []dataPoint) (ret []float64) {
	for _, v := range in {
		ret = append(ret, v.value)
	}
	return
}

func (p *testItem) Get(isNum bool, last, shift_time int) []float64 {
	if isNum {
		if shift_time > 0 {
			//skip
		}
		size := len(p.dps)
		if size > last {
			return dpsToArr(p.dps[size-last:])
		}
		return dpsToArr(p.dps)
	} else {
		// TODO
		return []float64{}
	}
}

func TestExec(t *testing.T) {
	dps := []dataPoint{
		{1, 1.0}, {2, 2.0}, {3, 3.0},
		{4, 4.0}, {5, 5.0}, {6, 6.0},
		{7, 7.0}, {8, 8.0}, {9, 9.0},
		{10, 10.0},
	}
	cases := []struct {
		dps  []dataPoint
		expr string
		want bool
	}{
		{dps, "min(#5) < 3", false},
		{dps, "min(#8) < 3", false},
		{dps, "min(#9) < 3", true},
		{dps, "min(#8) < 3 || min(#9) < 3", true},
		{dps, "min(#8) < 3 && min(#9) < 3", false},
	}
	for _, tc := range cases {
		expr, err := Parse([]byte(tc.expr), 0)
		if err != nil {
			t.Errorf("Parse(%s) err %v", tc.expr, err)
		}
		item := &testItem{dps: dps}
		if got := Exec(item, expr); got != tc.want {
			t.Errorf("%s = %v; want %v", tc.expr, got, tc.want)
		}
	}
}

/*
## resource
https://www.zabbix.com/documentation/4.0/manual/appendix/triggers/functions
https://www.zabbix.com/documentation/4.0/manual/config/triggers/expression
*/
