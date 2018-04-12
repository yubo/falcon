/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package core

import (
	"io/ioutil"
	"testing"

	"github.com/golang/glog"
)

// get config  ParseConfigFile(values, config.yaml)
// merge baseconf, config

var (
	testBaseYaml string
)

func init() {
	b, err := ioutil.ReadFile("./test/base.yaml")
	if err != nil {
		panic(err)
	}
	testBaseYaml = string(b)
}

func TestConfig(t *testing.T) {
	config, err := NewConfiger("./test/conf.yaml", testBaseYaml, []string{"./test/values.yaml"})
	if err != nil {
		t.Error(t)
	}

	err = config.Parse()
	if err != nil {
		t.Error(t)
	}

	glog.V(3).Infof("%s", config)
}

func TestRaw(t *testing.T) {
	config, _ := NewConfiger("./test/conf.yaml", testBaseYaml,
		[]string{"./test/values.yaml"})
	config.Parse()

	var cases = []struct {
		path string
		want interface{}
	}{
		{"foo1", "b_bar1"},
		{"foo2", "v_bar2"},
		{"foo3", "b_bar3"},
		{"fooo.foo", "bar"},
		{"na", nil},
		{"na.na", nil},
	}

	for _, c := range cases {
		if got := config.GetRaw(c.path); got != c.want {
			t.Errorf("config.GetRaw(%s) expected %#v got %#v", c.path, c.want, got)
		}
	}
}

func TestRead(t *testing.T) {
	config, _ := NewConfiger("./test/conf.yaml", testBaseYaml,
		[]string{"./test/values.yaml"})
	config.Parse()

	var (
		got  []string
		path = "fooo.foos"
	)

	if err := config.Read(path, &got); err != nil {
		t.Error(err)
	} else {
		t.Logf("config.Read(%s) got %#v", path, got)
	}
}
