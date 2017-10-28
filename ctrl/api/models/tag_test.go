/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"testing"

	"github.com/yubo/falcon"
)

func TestNewTagSchema(t *testing.T) {
	cases := []struct {
		name   string
		schema string
		want   error
	}{
		{name: "1", schema: "cop,owt", want: falcon.ErrParam},
		{name: "2", schema: "", want: nil},
		{name: "3", schema: XIAOMI_SCHEMA, want: nil},
	}
	for _, sc := range cases {
		if _, got := NewTagSchema(sc.schema); got != sc.want {
			t.Errorf("createSchema(%q) = %v; want %v", sc.name, got, sc.want)
		}
	}
}

func TestTagMap(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		want error
	}{
		{name: "empty string", tag: "", want: falcon.ErrParam},
		{name: "incomplete kv", tag: "a=b,", want: falcon.ErrParam},
		{name: "incomplete value1", tag: "a=", want: falcon.ErrParam},
		{name: "incomplete value2", tag: "a=b,c=", want: falcon.ErrParam},
		{name: "incomplete key1", tag: "=b", want: falcon.ErrParam},
		{name: "incomplete key2", tag: "a=b,=d", want: falcon.ErrParam},
		{name: "sample tag1", tag: "a=b", want: nil},
		{name: "sample tag2", tag: "a=b,c=d", want: nil},
		{name: "sample tag2", tag: "a=b=b,c=d", want: nil},
		{name: "sample tag3", tag: " a = b ,c = d ", want: nil},
	}
	for _, tc := range cases {
		if _, got := tagMap(tc.tag); got != tc.want {
			t.Errorf("%q tagMap(%q) = %v; want %v", tc.name, tc.tag, got, tc.want)
		}
	}
}

func TestTagFmtErr(t *testing.T) {

	ts, _ := NewTagSchema("a,b,c;d;e,f,")
	cases := []struct {
		name  string
		tag   string
		force bool
		want  error
	}{
		{name: "1", force: false, tag: "b=1", want: falcon.ErrParam},
		{name: "2", force: false, tag: "a=1,c=1", want: falcon.ErrParam},
		{name: "3", force: false, tag: "g=1", want: falcon.ErrParam},
		{name: "4", force: false, tag: "a=1,g=1", want: falcon.ErrParam},

		{name: "5", force: true, tag: "b=1", want: nil},
		{name: "6", force: true, tag: "a=1,c=1", want: nil},
		{name: "7", force: true, tag: "g=1", want: falcon.ErrParam},
		{name: "8", force: true, tag: "a=1,g=1", want: nil},
	}
	for _, tc := range cases {
		if _, got := ts.Fmt(tc.tag, tc.force); got != tc.want {
			t.Errorf("%q schema.Fmt(%q, %v) = %v; want %v", tc.name, tc.tag, tc.force, got, tc.want)
		}
	}
}

func TestTagFmt(t *testing.T) {

	ts, _ := NewTagSchema("a,b,c;d;e,f,")
	cases := []struct {
		name  string
		tag   string
		force bool
		want  string
	}{
		{name: "1", tag: "a=1", force: false, want: "a=1"},
		{name: "2", tag: "b=2,a=1", force: false, want: "a=1,b=2"},
		{name: "3", tag: "b=2,a=1,e=3", force: false, want: "a=1,b=2,e=3"},
		{name: "4", tag: "d=3,g=4,b=2,a=1", force: true, want: "a=1,b=2,d=3"},
	}
	for _, tc := range cases {
		if got, _ := ts.Fmt(tc.tag, tc.force); got != tc.want {
			t.Errorf("%q schema.Fmt(%q, %v) = %v; want %v", tc.name, tc.tag, tc.force, got, tc.want)
		}
	}
}

func TestTagParent(t *testing.T) {
	cases := []struct {
		tag  string
		want string
	}{
		{tag: "", want: ""},
		{tag: "a=1", want: ""},
		{tag: "a=1,b=1", want: "a=1"},
		{tag: "a=1,b=1,c=1", want: "a=1,b=1"},
	}
	for _, tc := range cases {
		if got, want := TagParent(tc.tag), tc.want; got != want {
			t.Errorf("TagParent(%q) = %v; want %v",
				tc.tag, got, want)
		}
	}

}

func TestTagParents(t *testing.T) {
	cases := []struct {
		tag  string
		want []string
	}{
		{tag: "", want: []string{""}},
		{tag: "a=1", want: []string{""}},
		{tag: "a=1,b=1", want: []string{"", "a=1"}},
		{tag: "a=1,b=1,c=1", want: []string{"", "a=1", "a=1,b=1"}},
	}
	for _, tc := range cases {
		if got, want := TagParents(tc.tag), tc.want; stringscmp(got, want) != 0 {
			t.Errorf("TagParents(%q) = %q; want %q",
				tc.tag, got, want)
		}
	}
}

func TestTagRelation(t *testing.T) {
	cases := []struct {
		tag  string
		want []string
	}{
		{tag: "", want: []string{""}},
		{tag: "a=1", want: []string{"", "a=1"}},
		{tag: "a=1,b=1", want: []string{"", "a=1", "a=1,b=1"}},
		{tag: "a=1,b=1,c=1", want: []string{"", "a=1", "a=1,b=1", "a=1,b=1,c=1"}},
	}
	for _, tc := range cases {
		if got, want := TagRelation(tc.tag), tc.want; stringscmp(got, want) != 0 {
			t.Errorf("TagParents(%q) = %q; want %q",
				tc.tag, got, want)
		}
	}
}
