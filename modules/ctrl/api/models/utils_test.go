/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"testing"
)

func testMdiff(t *testing.T) {
	cases := []struct {
		src   []string
		dst   []string
		wanta []string
		wantd []string
	}{
		{src: []string{}, dst: []string{}, wanta: []string{}, wantd: []string{}},
		{src: []string{"1", "2"}, dst: []string{"3", "4"}, wanta: []string{"3", "4"}, wantd: []string{"1", "2"}},
		{src: []string{"1", "2"}, dst: []string{"2", "3"}, wanta: []string{"3"}, wantd: []string{"1"}},
	}
	for _, c := range cases {
		if gota, gotd := MdiffStr(c.src, c.dst); stringscmp(gota,
			c.wanta) != 0 || stringscmp(gotd, c.wantd) != 0 {
			t.Errorf("Mdiff(%v,%v) = %v, %v; want %v %v",
				c.src, c.dst, gota, gotd, c.wanta, c.wantd)
		}
	}
}

func testintscmp64(t *testing.T) {
	cases := []struct {
		a    []int64
		b    []int64
		want bool
	}{
		{a: []int64{2, 3, 4}, b: []int64{3, 4, 2}, want: true},
		{a: []int64{2, 3, 4, 4}, b: []int64{3, 4, 4, 2}, want: true},
		{a: []int64{2, 3, 4, 4, 5}, b: []int64{3, 4, 4, 2, 6}, want: false},
		{a: []int64{2, 3, 4, 6}, b: []int64{3, 4, 2}, want: false},
	}
	for _, c := range cases {
		if got := intscmp64(c.a, c.b); (got == 0) != c.want {
			t.Errorf("intscmp64(%v,%v) = %v; want %v",
				c.a, c.b, got, c.want)
		}
	}
}

func test_t(t *testing.T) {
	cases := []struct {
		a string
		b string
	}{
		{a: "a=1,b=2", b: "a.1_b.2"},
		{a: "a=1=2,b=2", b: "a.1=2_b.2"},
	}
	for _, c := range cases {
		if got := TagToOld(c.a); got != c.b {
			t.Errorf("_t(%s) = %s; want %v", c.a, got, c.b)
		}
		if got := TagToNew(c.b); got != c.a {
			t.Errorf("_T(%s) = %s; want %v", c.b, got, c.a)
		}
	}
}

func TestUtils(t *testing.T) {
	testMdiff(t)
	testintscmp64(t)
	test_t(t)
}
