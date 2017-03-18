/*
 * Copyright 2016 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"testing"
)

func TestMdiff(t *testing.T) {
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

func Testintscmp64(t *testing.T) {
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
