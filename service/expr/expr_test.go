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
		} else {
			t.Logf("Parse(%s) = %s %v", tc.text, expr, err)
		}
	}
}

/*
## resource
https://www.zabbix.com/documentation/4.0/manual/appendix/triggers/functions
https://www.zabbix.com/documentation/4.0/manual/config/triggers/expression
*/
