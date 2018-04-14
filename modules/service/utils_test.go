/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import "testing"

func testTagsMatch(t *testing.T) {
	cases := []struct {
		pattern string
		tags    string
		want    bool
	}{
		{"", "a=1", true},
		{"a=1", "", false},
		{"a=1", "a=1", true},
		{"a=1", "b=1", false},
		{"a=1", "a=1,b=1", true},
		{"a=1,b=1", "a=1", false},
		{"a=1,b=1", "a=1,b=1", true},
		{"a=1,b=1", "a=1,b=1,c=1", true},
	}
	for _, tc := range cases {
		if got := tagsMatch(tc.pattern, tc.tags); got != tc.want {
			t.Errorf("tagsMatch(%s, %s) = %v; want %v", tc.pattern, tc.tags, got, tc.want)
		}
	}
}
