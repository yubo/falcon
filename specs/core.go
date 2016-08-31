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
	ErrEmpty       = errors.New("empty items")
	EPERM          = errors.New("Operation not permitted")
	ENOENT         = errors.New("No such file or directory")
	ESRCH          = errors.New("No such process")
	EINTR          = errors.New("Interrupted system call")
	EIO            = errors.New("I/O error")
	ENXIO          = errors.New("No such device or address")
	E2BIG          = errors.New("Argument list too long")
	ENOEXEC        = errors.New("Exec format error")
	EBADF          = errors.New("Bad file number")
	ECHILD         = errors.New("No child processes")
	EAGAIN         = errors.New("Try again")
	ENOMEM         = errors.New("Out of memory")
	EACCES         = errors.New("Permission denied")
	EFAULT         = errors.New("Bad address")
	ENOTBLK        = errors.New("Block device required")
	EBUSY          = errors.New("Device or resource busy")
	EEXIST         = errors.New("File exists")
	EXDEV          = errors.New("Cross-device link")
	ENODEV         = errors.New("No such device")
	ENOTDIR        = errors.New("Not a directory")
	EISDIR         = errors.New("Is a directory")
	EINVAL         = errors.New("Invalid argument")
	ENFILE         = errors.New("File table overflow")
	EMFILE         = errors.New("Too many open files")
	ENOTTY         = errors.New("Not a typewriter")
	ETXTBSY        = errors.New("Text file busy")
	EFBIG          = errors.New("File too large")
	ENOSPC         = errors.New("No space left on device")
	ESPIPE         = errors.New("Illegal seek")
	EROFS          = errors.New("Read-only file system")
	EMLINK         = errors.New("Too many links")
	EPIPE          = errors.New("Broken pipe")
	EDOM           = errors.New("Math argument out of domain of func")
	ERANGE         = errors.New("Math result not representable")
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
