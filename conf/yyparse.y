// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%{

package conf

import (
	"os"
	"fmt"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend"
	"github.com/yubo/falcon/lb"
)

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
%token INTERVAL IFACE_PREFIX DISABLED PAYLOAD_SIZE CONN_TIMEOUT CALL_TIMEOUT UPSTREAMS LB 
%token CONCURRENCY BACKENDS LBS TSDB FALCON 
%token BACKEND DSN DB_MAX_IDLE DB_MAX_CONN CONTAINER SHM_MAGIC_CODE SHM_KEY_START_ID
%token SHM_SEGMENT_SIZE MIGRATE STORAGE RRD HDISKS
%token ON YES OFF NO INCLUDE ROOT METRICS

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
| ss text num ';'  { yy_ss[$2] = fmt.Sprintf("%d", $3) }
| ss text bool ';' { yy_ss[$2] = fmt.Sprintf("%s", $3) }
| ss INCLUDE text ';' { yy.include($3) }
;

as:
| as text          { yy_as = append(yy_as, $2) }
| as INCLUDE text ';' { yy.include($3) }
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
		if conf.Agent.Conf.Params.CtrlAddr == ""{
			yy.Error("ctrlAddr empty")
		}
	}
}| ctrl_mod '}' {
	if conf.Ctrl.Conf.Host == ""{
		conf.Ctrl.Conf.Host, _ = os.Hostname()
	}
	if !conf.Ctrl.Conf.Disabled || yy.debug {
		conf.Modules = append([]falcon.Module{&conf.Ctrl}, conf.Modules...)
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
	conf.Agent.Conf  = falcon.ConfAgentDef
	yy_mod_params    = &conf.Agent.Conf.Params
}
;

mod_name:
| text { yy_mod_params.Name = $1 }
;

agent_mod_item:
   mod_item
 | INTERVAL num { conf.Agent.Conf.Interval = $2 }
 | PAYLOAD_SIZE num { conf.Agent.Conf.PayloadSize = $2 }
 | IFACE_PREFIX as {
	conf.Agent.Conf.IfPre = yy_as
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
  ctrl_start text '{'        { conf.Ctrl.Conf.Name = $2 }
| ctrl_mod ctrl_mod_item ';'
| ctrl_mod INCLUDE text ';'  { yy.include($3) }
;

ctrl_start: CTRL {
	conf.Ctrl.Conf   = falcon.ConfCtrlDef
}
;

ctrl_mod_item:
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DISABLED bool   { conf.Ctrl.Conf.Disabled = $2 }
 | DEBUG           { conf.Ctrl.Conf.Debug = 1 }
 | DEBUG num       { conf.Ctrl.Conf.Debug = $2 }
 | HOST text       { conf.Ctrl.Conf.Host = $2 }
 | DSN text        { conf.Ctrl.Conf.Dsn = $2 }
 | DB_MAX_IDLE num { conf.Ctrl.Conf.DbMaxIdle = $2 }
 | DB_MAX_CONN num { conf.Ctrl.Conf.DbMaxConn = $2 }
 | CONTAINER '{' ss '}' { 
 	conf.Ctrl.Conf.Ctrl.Set(falcon.APP_CONF_FILE, yy_ss)
	yy_ss = make(map[string]string)
}| METRICS '{' as '}' {
 	conf.Ctrl.Conf.Metrics = yy_as
	yy_as = make([]string, 0)
}
;

lb_mod: lb_start mod_name '{'
| lb_mod lb_mod_item ';'
| lb_mod INCLUDE text ';'  { yy.include($3) }
;

lb_start: LB {
	conf.Lb          = append(conf.Lb, lb.Lb{})
	yy_lb            = &conf.Lb[len(conf.Lb)-1]
	yy_specs_backend = &falcon.Backend{}
	yy_lb.Conf       = falcon.ConfLbDef
	yy_mod_params    = &yy_lb.Conf.Params
}
;

lb_mod_item:
   mod_item
 | PAYLOAD_SIZE num { yy_lb.Conf.PayloadSize = $2 }
 | BACKENDS '{' lb_backends '}'
;

lb_backends:
| lb_backends lb_backends_item ';'
;

lb_backends_item:
| TSDB backend_name '{' lb_backend '}' { 
	yy_specs_backend.Type = "tsdb"
	if !yy_specs_backend.Disabled || yy.debug {
		yy_lb.Conf.Backends = append(yy_lb.Conf.Backends, *yy_specs_backend)
	}
	yy_specs_backend = &falcon.Backend{}
}
| FALCON backend_name '{' lb_backend '}'{
	yy_specs_backend.Type = "falcon"
	if !yy_specs_backend.Disabled || yy.debug {
		yy_lb.Conf.Backends = append(yy_lb.Conf.Backends, *yy_specs_backend)
	}
	yy_specs_backend = &falcon.Backend{}
}
;

backend_name:
| text {
	yy_specs_backend.Name = $1
}
;

lb_backend:
| lb_backend lb_backend_item ';'
;

lb_backend_item:
| DISABLED bool { yy_specs_backend.Disabled = $2 }
| UPSTREAMS '{' ss '}' { 
	yy_specs_backend.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", yy_specs_backend.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}
;

backend_mod: backend_start mod_name '{'
| backend_mod backend_mod_item ';'
| backend_mod INCLUDE text ';'  { yy.include($3) }
;

backend_start: BACKEND {
	conf.Backend    = append(conf.Backend, backend.Backend{})
	yy_backend      = &conf.Backend[len(conf.Backend)-1]
	yy_backend.Conf = falcon.ConfBackendDef
	yy_mod_params   = &yy_backend.Conf.Params
}
;

backend_mod_item:
  mod_item
| DSN text { yy_backend.Conf.Dsn = $2 }
| DB_MAX_IDLE num { yy_backend.Conf.DbMaxIdle = $2 }
| SHM_MAGIC_CODE num { yy_backend.Conf.ShmMagic = uint32($2) }
| SHM_KEY_START_ID num { yy_backend.Conf.ShmKey = $2 }
| SHM_SEGMENT_SIZE num { yy_backend.Conf.ShmSize = $2 }
| MIGRATE '{' backend_migrate '}'
| STORAGE RRD '{' backend_storage '}' { yy_backend.Conf.Storage.Type = "rrd" }
;

backend_migrate:
| backend_migrate backend_migrate_item ';'

backend_migrate_item:
| DISABLED bool { yy_backend.Conf.Migrate.Disabled = $2 }
| UPSTREAMS '{' ss '}' {
	yy_backend.Conf.Migrate.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", yy_backend.Conf.Migrate.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}
;

backend_storage:
| HDISKS as ';' {
	yy_backend.Conf.Storage.Hdisks = yy_as
	yy_as = make([]string, 0)
}

%%

