/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package core

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestFuncs(t *testing.T) {
	os.Setenv("FOO", "bar")
	cases := map[string]string{
		`{{hello}}`:                                                  "hello!",
		`{{env "FOO"}}`:                                              "bar",
		`{{expandenv "hello ${FOO}"}}`:                               "hello bar",
		`{{expandenv "hello $FOO"}}`:                                 "hello bar",
		`{{base "foo/bar"}}`:                                         "bar",
		`{{dir "foo/bar/baz"}}`:                                      "foo/bar",
		`{{clean "/foo/../foo/../bar"}}`:                             "/bar",
		`{{ext "/foo/bar/baz.txt"}}`:                                 ".txt",
		`{{isAbs "/foo/bar/baz"}}`:                                   "true",
		`{{isAbs "foo/bar/baz"}}`:                                    "false",
		`{{quote "a" "b" "c"}}`:                                      `"a" "b" "c"`,
		`{{quote "\"a\"" "b" "c"}}`:                                  `"\"a\"" "b" "c"`,
		`{{quote 1 2 3 }}`:                                           `"1" "2" "3"`,
		`{{squote "a" "b" "c"}}`:                                     `'a' 'b' 'c'`,
		`{{squote 1 2 3 }}`:                                          `'1' '2' '3'`,
		`{{if contains "cat" "fair catch"}}1{{end}}`:                 "1",
		`{{if hasPrefix "cat" "catch"}}1{{end}}`:                     "1",
		`{{if hasSuffix "cat" "ducat"}}1{{end}}`:                     "1",
		`{{trim "   5.00   "}}`:                                      "5.00",
		`{{trimAll "$" "$5.00$"}}`:                                   "5.00",
		`{{trimPrefix "$" "$5.00"}}`:                                 "5.00",
		`{{trimSuffix "$" "5.00$"}}`:                                 "5.00",
		`{{$v := "foo$bar$baz" | split "$"}}{{$v._0}}`:               "foo",
		`{{toString 1 | kindOf }}`:                                   "string",
		`{{$s := list 1 2 3 | toStrings }}{{ index $s 1 | kindOf }}`: "string",
		`{{list "a" "b" "c" | join "-" }}`:                           "a-b-c",
		`{{list 1 2 3 | join "-" }}`:                                 "1-2-3",
		`{{list "c" "a" "b" | sortAlpha | join "" }}`:                "abc",
		`{{list 2 1 4 3 | sortAlpha | join "" }}`:                    "1234",
		`{{b64enc "hello,world"}}`:                                   "aGVsbG8sd29ybGQ=",
		`{{b64dec "aGVsbG8sd29ybGQ="}}`:                              "hello,world",
		`{{$b := "b"}}{{"c" | cat "a" $b}}`:                          "a b c",
		`{{indent 4 "a\nb\nc"}}`:                                     "    a\n    b\n    c",
		`{{nindent 4 "a\nb\nc"}}`:                                    "\n    a\n    b\n    c",
		`{{"a b c" | replace " " "-"}}`:                              "a-b-c",
		`{{ timeOf "3600" }}`:                                        "3600",
		`{{ timeOf "3600s" }}`:                                       "3600",
		`{{ timeOf "60m" }}`:                                         "3600",
		`{{ timeOf "1h" }}`:                                          "3600",
		`{{ timeOf "1d" }}`:                                          "86400",
	}
	for tpl, want := range cases {
		if err := runt(tpl, want); err != nil {
			t.Error(err)
		}
	}
}

// runt runs a template and checks that the output exactly matches the expected string.
func runt(tpl, expect string) error {
	return runtv(tpl, expect, map[string]string{})
}

// runtv takes a template, and expected return, and values for substitution.
//
// It runs the template and verifies that the output is an exact match.
func runtv(tpl, expect string, vars interface{}) error {
	t := template.Must(template.New("test").Funcs(FuncMap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("%s Expected '%s', got '%s'", tpl, expect, b.String())
	}
	return nil
}
