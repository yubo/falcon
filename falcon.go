/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"fmt"

	"github.com/yubo/gotool/flags"
)

const (
	VERSION = "0.0.1"
)

func Version_handle(arg interface{}) {
	fmt.Println(VERSION)
}

func Git_handle(arg interface{}) {
	fmt.Println(COMMIT)
}

func Changelog_handle(arg interface{}) {
	fmt.Println(CHANGELOG)
}

func Help_handle(arg interface{}) {
	flags.Usage()
}
