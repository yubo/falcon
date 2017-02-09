/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package specs

import (
	"flag"
	"fmt"
	"strings"

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
	Debug       int    `json:"debug"`
	ConnTimeout int    `json:"connTimeout"`
	CallTimeout int    `json:"callTimeout"`
	Concurrency int    `json:"concurrency"`
	Disabled    bool   `json:"disabled"`
	Http        bool   `json:"http"`
	Rpc         bool   `json:"rpc"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	HttpAddr    string `json:"httpAddr"`
	RpcAddr     string `json:"rpcAddr"`
	CtrlAddr    string `json:"ctrlAddr"`
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

/*******************************
	CONFIG
********************************/
type ConfLb struct {
	Params      ModuleParams `json:"params"`
	PayloadSize int          `json:"payloadSize"`
	Backends    []Backend    `json:"backends"`
}

type ConfAgent struct {
	Params      ModuleParams `json:"params"`
	Interval    int          `json:"interval"`
	PayloadSize int          `json:"payloadSize"`
	IfPre       []string     `json:"ifPre"`
	Upstreams   []string     `json:"upstreams"`
}

func (c ConfAgent) String() string {
	return fmt.Sprintf("%s (\n%s\n)\n"+
		"%-17s %s\n"+
		"%-17s %d\n"+
		"%-17s %d\n"+
		"%-17s %s",
		"params", IndentLines(1, c.Params.String()),
		"ifprefix", strings.Join(c.IfPre, ", "),
		"interval", c.Interval,
		"payloadSize", c.PayloadSize,
		"upstreams", strings.Join(c.Upstreams, ", "))
}

type Storage struct {
	Type   string   `json:"type"`
	Hdisks []string `json:"hdisks"`
}

func (o Storage) String() string {
	return fmt.Sprintf("%-12s %s\n%-12s %s",
		"type", o.Type,
		"hdisks", strings.Join(o.Hdisks, ", "))
}

type ConfBackend struct {
	Params          ModuleParams `json:"params"`
	Migrate         Migrate      `json:"migrate"`
	Idx             bool         `json:"idx"`
	IdxInterval     int          `json:"idxInterval"`
	IdxFullInterval int          `json:"idxFullInterval"`
	Dsn             string       `json:"dsn"`
	DbMaxIdle       int          `json:"dbMaxIdle"`
	ShmMagic        uint32       `json:"shmMagic"`
	ShmKey          int          `json:"shmKey"`
	ShmSize         int          `json:"shmSize"`
	Storage         Storage      `json:"storage"`
}

func (c ConfBackend) String() string {
	http := c.Params.HttpAddr
	rpc := c.Params.RpcAddr

	if !c.Params.Http {
		http += "(disabled)"
	}
	if !c.Params.Rpc {
		rpc += "(disabled)"
	}
	return fmt.Sprintf("%s (\n%s\n)\n"+
		"%-17s %v\n"+
		"%-17s %d\n"+
		"%-17s %d\n"+
		"%-17s %s\n"+
		"%-17s %d\n"+
		"%-17s 0x%x\n"+
		"%-17s 0x%x\n"+
		"%-17s %d\n"+
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)",
		"params", IndentLines(1, c.Params.String()),
		"idx", c.Idx,
		"idx_interval", c.IdxInterval,
		"idx_full_interval", c.IdxFullInterval,
		"dsn", c.Dsn,
		"dbmaxidle", c.DbMaxIdle,
		"magic_code", c.ShmMagic,
		"key_start_id", c.ShmKey,
		"segment_size", c.ShmSize,
		"migrate", IndentLines(1, c.Migrate.String()),
		"storage", IndentLines(1, c.Storage.String()))
}

type ConfCtrl struct {
	Params ModuleParams `json:"params"`
}

var (
	ConfAgentDef = ConfAgent{
		Params: ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Name:        "Agent Module",
			Disabled:    false,
			Http:        true,
			Rpc:         true,
			HttpAddr:    "127.0.0.1:1988",
			RpcAddr:     "127.0.0.1:1989",
			CtrlAddr:    "127.0.0.1:8001",
		},
		Interval:    60,
		PayloadSize: 16,
		IfPre:       []string{"eth", "em"},
		Upstreams:   []string{},
	}

	ConfBackendDef = ConfBackend{
		Params: ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Disabled:    false,
			Http:        true,
			Rpc:         true,
			HttpAddr:    "0.0.0.0:7021",
			RpcAddr:     "0.0.0.0:7020",
			CtrlAddr:    "",
		},
		Migrate: Migrate{
			Disabled:  false,
			Upstreams: map[string]string{},
		},
		Idx:             true,
		IdxInterval:     30,
		IdxFullInterval: 86400,
		Dsn:             "falcon:1234@tcp(127.0.0.1:3306)/falcon?loc=Local&parseTime=true",
		DbMaxIdle:       4,
		ShmMagic:        0x80386,
		ShmKey:          0x7020,
		ShmSize:         256 * (1 << 20), // 256m
		Storage: Storage{
			Type: "rrd",
		},
	}

	ConfLbDef = ConfLb{
		Params: ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Disabled:    false,
			Http:        true,
			HttpAddr:    "0.0.0.0:6060",
			Rpc:         true,
			RpcAddr:     "0.0.0.0:8433",
			CtrlAddr:    "",
		},
		PayloadSize: 16,
		Backends:    make([]Backend, 0),
	}

	ConfCtrlDef = ConfCtrl{
		Params: ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Name:        "Control Module",
			Disabled:    false,
			Http:        true,
			HttpAddr:    "0.0.0.0:6060",
			Rpc:         true,
			RpcAddr:     "0.0.0.0:8433",
			CtrlAddr:    "N/A",
		},
	}
)
