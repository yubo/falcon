/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package parse

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "github.com/yubo/falcon/agent/modules"
	_ "github.com/yubo/falcon/ctrl/modules"
	_ "github.com/yubo/falcon/service/modules"
	_ "github.com/yubo/falcon/transfer/modules"
)

type Test struct {
	in  string
	out string
}

func init() {
	flag.Lookup("logtostderr").Value.Set("true")
	//flag.Lookup("v").Value.Set("5")
}

func TestExprText(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	goroot := os.Getenv("GOROOT")

	for i, test := range []Test{
		{"hello,{world}", "hello,{world}"},
		{"${GOPATH}", gopath},
		{"gopath:${GOPATH}", fmt.Sprintf("gopath:%s", gopath)},
		{"gopath:${GOPATH};goroot:${GOROOT}",
			fmt.Sprintf("gopath:%s;goroot:%s", gopath, goroot)},
	} {
		out := exprText([]byte(test.in))
		if out != test.out {
			t.Errorf(`#%d: exprText("%s")="%s"; want "%s"`,
				i, test.in, out, test.out)
		}
	}
}

func TestParse(t *testing.T) {
	conf := Parse("../docs/etc/falcon.example.conf")
	fmt.Printf("conf:\n%s\n", conf)
}
