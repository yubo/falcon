/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"fmt"
	"strings"
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

const (
	C_RUN_MODE                = "runmode"
	C_HTTP_PORT               = "httpport"
	C_ENABLE_DOCS             = "enabledocs"
	C_SEESION_GC_MAX_LIFETIME = "sessiongcmaxlifetime"
	C_SESSION_COOKIE_LIFETIME = "sessioncookielifetime"
	C_AUTH_MODULE             = "authmodule"
	C_CACHE_MODULE            = "cachemodule"
	C_LDAP_ADDR               = "ldapaddr"
	C_LDAP_BASE_DN            = "ldapbasedn"
	C_LDAP_BIND_DN            = "ldapbinddn"
	C_LDAP_BIND_PWD           = "ldapbindpwd"
	C_LDAP_FILTER             = "ldapfilter"
	C_MISSO_REDIRECT_URL      = "missoredirecturl"
	C_GITHUB_CLIENT_ID        = "githubclientid"
	C_GITHUB_CLIENT_SECRET    = "githubclientsecret"
	C_GITHUB_REDIRECT_URL     = "githubredirecturl"
	C_GOOGLE_CLIENT_ID        = "googleclientid"
	C_GOOGLE_CLIENT_SECRET    = "googleclientsecret"
	C_GOOGLE_REDIRECT_URL     = "googleredirecturl"
)

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
		},
		PayloadSize: 16,
		Backends:    make([]Backend, 0),
	}

	ConfCtrlDef = ConfCtrl{Name: "ctrl"}

	ConfDefault = map[string]map[string]string{
		"agent":   map[string]string{},
		"lb":      map[string]string{},
		"backend": map[string]string{},
		"ctrl": map[string]string{
			C_RUN_MODE:                "dev",
			C_HTTP_PORT:               "8001",
			C_ENABLE_DOCS:             "true",
			C_SEESION_GC_MAX_LIFETIME: "86400",
			C_SESSION_COOKIE_LIFETIME: "86400",
			C_AUTH_MODULE:             "ldap",
			C_CACHE_MODULE:            "host,role,system,tag,user",
		},
	}
)

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
	// only in falcon.conf
	Debug     int
	Disabled  bool
	Name      string
	Host      string
	Dsn       string
	DbMaxIdle int
	DbMaxConn int
	Metrics   []string
	Ctrl      Configer
	Agent     Configer
	Lb        Configer
	Backend   Configer
	// 1: default, 2: db, 3: ConfCtrl.Container
	// height will cover low
}

func (c ConfCtrl) String() string {
	var s string
	for k, v := range c.Metrics {
		s += fmt.Sprintf("%s ", v)
		if k%5 == 0 {
			s += "\n"
		}
	}
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"Metrics", IndentLines(1, s),
		"Container", IndentLines(1, c.Ctrl.String()),
	)
}
