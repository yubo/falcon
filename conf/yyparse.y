// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%{

package conf

import (
	"os"
	"fmt"
	"github.com/golang/glog"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/backend"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/lb"
	"github.com/yubo/falcon/specs"
)


// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
//const eof = 0

%}

%union {
	num int
	text string
	b bool
}

%type <b> bool
%type <text> text
%type <num> num

%token <num> NUM
%token <text> TEXT IPA

%token '{' '}' ';'
%token PID_FILE AGENT DEBUG HOST HTTP HTTP_ADDR RPC RPC_ADDR CTRL CTRL_ADDR
%token INTERVAL IFACE_PREFIX DISABLED BATCH CONN_TIMEOUT CALL_TIMEOUT UPSTREAMS LB 
%token CONCURRENCY BACKENDS LBS TSDB FALCON 
%token BACKEND DSN DB_MAX_IDLE SHM_MAGIC_CODE SHM_KEY_START_ID
%token SHM_SEGMENT_SIZE MIGRATE STORAGE RRD HDISKS
%token ON YES OFF NO INCLUDE ROOT

%%



config: 
| config conf
;

bool:
  ON  { $$ = true }
| YES { $$ = true }
| OFF { $$ = false }
| NO  { $$ = false }
|     { $$ = true }
;

text:
  IPA  { $$ = string(yy.t) }
| TEXT { $$ = exprText(yy.t) }
;

num:
NUM { $$ = yy.i }
;

ss:
| ss text text ';' { yy_ss[$2] = $3 }
;

as:
| as text { yy_as = append(yy_as, $2) }
;

conf: ';'
 | PID_FILE text ';' { conf.PidFile = $2 }
 | INCLUDE text ';'  { yy.include($2) }
 | ROOT text ';' { 
 	if err := os.Chdir($2); err != nil {
 		yy.Error(err.Error())
 	}
}| agent_mod '}' {
	if yy_mod_params.Host == ""{
		yy_mod_params.Host, _ = os.Hostname()
	}
	if !yy_mod_params.Disabled || yy.debug {
		conf.Modules = append(conf.Modules, &conf.Agent)
		if conf.Agent.Params.CtrlAddr == ""{
			yy.Error("ctrlAddr empty")
		}
	}
}| ctrl_mod '}' {
	if yy_mod_params.Host == ""{
		yy_mod_params.Host, _ = os.Hostname()
	}
	if !yy_mod_params.Disabled || yy.debug {
		conf.Modules = append([]specs.Module{&conf.Ctrl}, conf.Modules...)
	}
}| lb_mod '}' {
	yy_mod_params.Name = fmt.Sprintf("lb %s", yy_mod_params.Name)
	if yy_mod_params.Host == ""{
		yy_mod_params.Host, _ = os.Hostname()
	}
	if !yy_mod_params.Disabled || yy.debug {
		conf.Modules = append(conf.Modules, yy_lb)
	}
}| backend_mod '}' {
	yy_mod_params.Name = fmt.Sprintf("backend %s", yy_mod_params.Name)
	if yy_mod_params.Host == ""{
		yy_mod_params.Host, _ = os.Hostname()
	}
	if !yy_mod_params.Disabled || yy.debug {
		conf.Modules = append(conf.Modules, yy_backend)
	}
}
;

agent_mod:
  agent_start mod_name '{' 
| agent_mod agent_mod_item ';'
| agent_mod INCLUDE text ';'  { yy.include($3) }
;

agent_start: AGENT {
	conf.Agent       = agent.DefaultAgent
	yy_mod_params    = &conf.Agent.Params
}
;

mod_name:
| text { yy_mod_params.Name = $1 }
;

agent_mod_item:
   mod_item
 | INTERVAL num { conf.Agent.Interval = $2 }
 | BATCH num { conf.Agent.Batch = $2 }
 | IFACE_PREFIX as {
	conf.Agent.IfPre = yy_as
	yy_as = make([]string, 0)
}
;

mod_item:
 | DISABLED bool  { yy_mod_params.Disabled = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG            { yy_mod_params.Debug = 1 }
 | DEBUG num        { yy_mod_params.Debug = $2 }
 | HOST text        { yy_mod_params.Host = $2 }
 | HTTP bool        { yy_mod_params.Http = $2 }
 | RPC bool         { yy_mod_params.Rpc = $2 }
 | HTTP_ADDR text   { yy_mod_params.HttpAddr = $2 }
 | RPC_ADDR text    { yy_mod_params.RpcAddr = $2 }
 | CTRL_ADDR text   { yy_mod_params.CtrlAddr = $2 }
 | CONN_TIMEOUT num { yy_mod_params.ConnTimeout = $2 }
 | CALL_TIMEOUT num { yy_mod_params.CallTimeout = $2 }
 | CONCURRENCY num  { yy_mod_params.Concurrency = $2 }
;

ctrl_mod:
  ctrl_start mod_name '{' 
| ctrl_mod ctrl_mod_item ';'
| ctrl_mod INCLUDE text ';'  { yy.include($3) }
;

ctrl_start: CTRL {
	conf.Ctrl        = ctrl.DefaultCtrl
	yy_mod_params    = &conf.Ctrl.Params
	yy_specs_backend = &specs.Backend{}
}
;

ctrl_mod_item:
   mod_item
 | BACKENDS '{' ctrl_backends '}'
 | MIGRATE '{' ctrl_migrate '}'
 | LBS as {
	conf.Ctrl.Lbs = yy_as
	yy_as = make([]string, 0)
}
;

ctrl_backends:
|ctrl_backends ctrl_backends_item ';'
;

ctrl_backends_item:
| TSDB backend_name '{' ctrl_backend '}' { 
	yy_specs_backend.Type = "tsdb"
	if !yy_specs_backend.Disabled || yy.debug {
		conf.Ctrl.Backends = append(conf.Ctrl.Backends, *yy_specs_backend)
	}
	yy_specs_backend = &specs.Backend{}
}
| FALCON backend_name '{' ctrl_backend '}'{
	yy_specs_backend.Type = "falcon"
	if !yy_specs_backend.Disabled || yy.debug {
		conf.Ctrl.Backends = append(conf.Ctrl.Backends, *yy_specs_backend)
	}
	yy_specs_backend = &specs.Backend{}
}
;

backend_name:
| text {
	yy_specs_backend.Name = $1
}
;

ctrl_backend:
| ctrl_backend ctrl_backend_item ';'
;

ctrl_backend_item:
| DISABLED bool { yy_specs_backend.Disabled = $2 }
| UPSTREAMS '{' ss '}' { 
	yy_specs_backend.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", yy_specs_backend.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}
;

ctrl_migrate:
| ctrl_migrate ctrl_migrate_item ';'

ctrl_migrate_item:
| DISABLED bool { conf.Ctrl.Migrate.Disabled = $2 }
| UPSTREAMS '{' ss '}' {
	conf.Ctrl.Migrate.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", conf.Ctrl.Migrate.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}
;

lb_mod: lb_start mod_name '{'
| lb_mod lb_mod_item ';'
| lb_mod INCLUDE text ';'  { yy.include($3) }
;

lb_start: LB {
	conf.Lb       = append(conf.Lb, lb.DefaultLb)
	yy_lb         = &conf.Lb[len(conf.Lb)-1]
	yy_mod_params = &yy_lb.Params
}
;

lb_mod_item:
   mod_item
 | BATCH num { yy_lb.Batch = $2 }
;


backend_mod: backend_start mod_name '{'
| backend_mod backend_mod_item ';'
| backend_mod INCLUDE text ';'  { yy.include($3) }
;

backend_start: BACKEND {
	conf.Backend     = append(conf.Backend, backend.DefaultBackend)
	yy_backend       = &conf.Backend[len(conf.Backend)-1]
	yy_mod_params    = &yy_backend.Params
}
;

backend_mod_item:
  mod_item
| DSN text { yy_backend.Dsn = $2 }
| DB_MAX_IDLE num { yy_backend.DbMaxIdle = $2 }
| SHM_MAGIC_CODE num { yy_backend.ShmMagic = uint32($2) }
| SHM_KEY_START_ID num { yy_backend.ShmKey = $2 }
| SHM_SEGMENT_SIZE num { yy_backend.ShmSize = $2 }
| STORAGE RRD '{' backend_storage '}' { yy_backend.Storage.Type = "rrd" }
;



backend_storage:
| HDISKS as ';' {
	yy_backend.Storage.Hdisks = yy_as
	yy_as = make([]string, 0)
}

%%

