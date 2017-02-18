/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"flag"
	"fmt"

	"github.com/yubo/gotool/flags"
)

const (
	IndentSize   = 4
	DEFAULT_STEP = 60 //s
	MIN_STEP     = 30 //s
	VERSION      = "0.0.2"
	REPLICAS     = 500
)

const (
	_ = iota
	ROUTINE_EVENT_M_EXIT
	ROUTINE_EVENT_M_RELOAD
)

const (
	APP_STATUS_INIT = iota
	APP_STATUS_PENDING
	APP_STATUS_RUNING
	APP_STATUS_EXIT
	APP_STATUS_RELOAD
)

const (
	GAUGE       = "GAUGE"
	DERIVE      = "DERIVE"
	COUNTER     = "COUNTER"
	MODULE_NAME = "\x1B[32m[SPECS]\x1B[0m "
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func init() {
	flags.NewCommand("version", "show falcon version information",
		Version, flag.ExitOnError)

	flags.NewCommand("git", "show falcon git version information",
		Git, flag.ExitOnError)

	flags.NewCommand("changelog", "show falcon changelog information",
		Changelog, flag.ExitOnError)
}

func Version(arg interface{}) {
	fmt.Println(VERSION)
}

func Git(arg interface{}) {
	fmt.Println(COMMIT)
}

func Changelog(arg interface{}) {
	fmt.Println(CHANGELOG)
}
