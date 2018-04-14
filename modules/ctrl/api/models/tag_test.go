/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"testing"

	"github.com/yubo/falcon/lib/core"
)

func testNewTagSchema(t *testing.T) {
	cases := []struct {
		name   string
		schema string
		want   error
	}{
		{"1", "cop,owt", core.ErrParam},
		{"2", "", nil},
		{"3", XIAOMI_SCHEMA, nil},
	}
	for _, sc := range cases {
		if _, got := NewTagSchema(sc.schema); got != sc.want {
			t.Errorf("createSchema(%q) = %v; want %v", sc.name, got, sc.want)
		}
	}
}

func testTagMap(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		want error
	}{
		{"empty string", "", core.ErrParam},
		{"incomplete kv", "a=b,", core.ErrParam},
		{"incomplete value1", "a=", core.ErrParam},
		{"incomplete value2", "a=b,c=", core.ErrParam},
		{"incomplete key1", "=b", core.ErrParam},
		{"incomplete key2", "a=b,=d", core.ErrParam},
		{"sample tag1", "a=b", nil},
		{"sample tag2", "a=b,c=d", nil},
		{"sample tag2", "a=b=b,c=d", nil},
		{"sample tag3", " a = b ,c = d ", nil},
	}
	for _, tc := range cases {
		if _, got := tagMap(tc.tag); got != tc.want {
			t.Errorf("%q tagMap(%q) = %v; want %v", tc.name, tc.tag, got, tc.want)
		}
	}
}

func testTagFmtErr(t *testing.T) {

	ts, _ := NewTagSchema("a,b,c;d;e,f,")
	cases := []struct {
		name  string
		force bool
		tag   string
		want  error
	}{
		{"1", false, "b=1", core.ErrParam},
		{"2", false, "a=1,c=1", core.ErrParam},
		{"3", false, "g=1", core.ErrParam},
		{"4", false, "a=1,g=1", core.ErrParam},

		{"5", true, "b=1", nil},
		{"6", true, "a=1,c=1", nil},
		{"7", true, "g=1", core.ErrParam},
		{"8", true, "a=1,g=1", nil},
	}
	for _, tc := range cases {
		if _, got := ts.Fmt(tc.tag, tc.force); got != tc.want {
			t.Errorf("%q schema.Fmt(%q, %v) = %v; want %v", tc.name, tc.tag, tc.force, got, tc.want)
		}
	}
}

func testTagFmt(t *testing.T) {

	ts, _ := NewTagSchema("a,b,c;d;e,f,")
	cases := []struct {
		name  string
		tag   string
		force bool
		want  string
	}{
		{"1", "a=1", false, "a=1"},
		{"2", "b=2,a=1", false, "a=1,b=2"},
		{"3", "b=2,a=1,e=3", false, "a=1,b=2,e=3"},
		{"4", "d=3,g=4,b=2,a=1", true, "a=1,b=2,d=3"},
	}
	for _, tc := range cases {
		if got, _ := ts.Fmt(tc.tag, tc.force); got != tc.want {
			t.Errorf("%q schema.Fmt(%q, %v) = %v; want %v", tc.name, tc.tag, tc.force, got, tc.want)
		}
	}
}

func testTagParent(t *testing.T) {
	cases := []struct {
		tag  string
		want string
	}{
		{"", ""},
		{"a=1", ""},
		{"a=1,b=1", "a=1"},
		{"a=1,b=1,c=1", "a=1,b=1"},
	}
	for _, tc := range cases {
		if got, want := TagParent(tc.tag), tc.want; got != want {
			t.Errorf("TagParent(%q) = %v; want %v",
				tc.tag, got, want)
		}
	}

}

func testTagParents(t *testing.T) {
	cases := []struct {
		tag  string
		want []string
	}{
		{"", []string{""}},
		{"a=1", []string{""}},
		{"a=1,b=1", []string{"", "a=1"}},
		{"a=1,b=1,c=1", []string{"", "a=1", "a=1,b=1"}},
	}
	for _, tc := range cases {
		if got, want := TagParents(tc.tag), tc.want; stringscmp(got, want) != 0 {
			t.Errorf("TagParents(%q) = %q; want %q",
				tc.tag, got, want)
		}
	}
}

func testTagRelation(t *testing.T) {
	cases := []struct {
		tag  string
		want []string
	}{
		{"", []string{""}},
		{"a=1", []string{"", "a=1"}},
		{"a=1,b=1", []string{"", "a=1", "a=1,b=1"}},
		{"a=1,b=1,c=1", []string{"", "a=1", "a=1,b=1", "a=1,b=1,c=1"}},
	}
	for _, tc := range cases {
		if got, want := TagRelation(tc.tag), tc.want; stringscmp(got, want) != 0 {
			t.Errorf("TagParents(%q) = %q; want %q",
				tc.tag, got, want)
		}
	}
}

func TestTag(t *testing.T) {
	testNewTagSchema(t)
	testTagMap(t)
	testTagFmt(t)
	testTagFmt(t)
	testTagParent(t)
	testTagParents(t)
	testTagRelation(t)
}
