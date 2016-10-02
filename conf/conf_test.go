/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package conf

import (
	"fmt"
	"os"
	"testing"
)

type Test struct {
	in  string
	out string
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
