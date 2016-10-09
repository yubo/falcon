/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package specs

import (
	"flag"
	"fmt"

	"github.com/yubo/gotool/flags"
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
	IndentSize   = 4
	DEFAULT_STEP = 60 //s
	MIN_STEP     = 30 //s
	GAUGE        = "GAUGE"
	DERIVE       = "DERIVE"
	COUNTER      = "COUNTER"
	VERSION      = "0.0.2"
	REPLICAS     = 500
	MODULE_NAME  = "\x1B[32m[SPECS]\x1B[0m "
)

type CmdOpts struct {
	ConfigFile string
	Args       []string
}

type ModuleParams struct {
	Debug       int
	ConnTimeout int
	CallTimeout int
	Concurrency int
	Disabled    bool
	Http        bool
	Rpc         bool
	Name        string
	Host        string
	HttpAddr    string
	RpcAddr     string
	CtrlAddr    string
}

func (p ModuleParams) String() string {
	http := p.HttpAddr
	rpc := p.RpcAddr

	if !p.Http {
		http += "(disabled)"
	}
	if !p.Rpc {
		rpc += "(disabled)"
	}

	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%-17s %d\n"+
		"%-17s %d\n"+
		"%-17s %d",
		"debug", p.Debug,
		"disabled", p.Disabled,
		"Name", p.Name,
		"Host", p.Host,
		"http", http,
		"rpc", rpc,
		"ctrl", p.CtrlAddr,
		"concurrency", p.Concurrency,
		"conntimeout", p.ConnTimeout,
		"callTimeout", p.CallTimeout)
}

type Backend struct {
	Disabled  bool
	Name      string
	Type      string
	Upstreams map[string]string
}

func (p Backend) String() string {
	var s1, s2 string

	s1 = fmt.Sprintf("%s %s", p.Type, p.Name)
	if p.Disabled {
		s1 += "(Disable)"
	}

	for k, v := range p.Upstreams {
		s2 += fmt.Sprintf("%-17s %s\n", k, v)
	}
	return fmt.Sprintf("%s upstreams (\n%s\n)", s1, IndentLines(1, s2))
}

type Migrate struct {
	Disabled  bool
	Upstreams map[string]string
}

func (p Migrate) String() string {
	var s string

	for k, v := range p.Upstreams {
		s += fmt.Sprintf("%-17s %s\n", k, v)
	}
	if s != "" {
		s = fmt.Sprintf("\n%s\n", IndentLines(1, s))
	}

	return fmt.Sprintf("%-17s %v\n"+
		"%s (%s)",
		"disable", p.Disabled,
		"upstreams", s)
}

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
