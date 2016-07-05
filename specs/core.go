/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package specs

import "errors"

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
	IndentSize   = 4
	DEFAULT_STEP = 60 //s
	MIN_STEP     = 30 //s
	GAUGE        = "GAUGE"
	DERIVE       = "DERIVE"
	COUNTER      = "COUNTER"
)

type CmdOpts struct {
	ConfigFile string
	Args       []string
}

var (
	ErrUnsupported = errors.New("unsupported")
	ErrExist       = errors.New("entry exists")
	ErrNoent       = errors.New("entry not exists")
	ErrParam       = errors.New("param error")
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
