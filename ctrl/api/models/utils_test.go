/*
 * Copyright 2016 yubo. All rights reserved.
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
