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
	"github.com/yubo/falcon/lb"
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
%token PID_FILE AGENT DEBUG HOST HTTP HTTP_ADDR INTERVAL IFACE_PREFIX HANDOFF
%token DISABLED BATCH CONN_TIMEOUT CALL_TIMEOUT UPSTREAMS HANDOFF RPC RPC_ADDR
%token REPLICAS CONCURRENCY BACKENDS TSDB SCHED CONSISTENT FALCON
%token BACKEND RRD_STORAGE DSN DB_MAX_IDLE SHM_MAGIC_CODE SHM_KEY_START_ID
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

conf: ';'
 | PID_FILE text ';' { conf.PidFile = $2 }
 | INCLUDE text ';'  { yy.include($2) }
 | ROOT text ';' { 
 	if err := os.Chdir($2); err != nil {
 		yy.Error(err.Error())
 	}
}| agent_mod '}' {
	if !*yy_mod_disable {
		conf.Modules = append(conf.Modules, &conf.Agent)
	}
}| lb_mod '}' {
	*yy_mod_name = fmt.Sprintf("lb %s", *yy_mod_name)
	if !*yy_mod_disable {
		conf.Modules = append(conf.Modules, yy_lb)
	}
}| backend_mod '}' {
	*yy_mod_name = fmt.Sprintf("backend %s", *yy_mod_name)
	if !*yy_mod_disable {
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
	yy_mod_name      = &conf.Agent.Name
	yy_mod_disable   = &conf.Agent.Disabled
	yy_mod_debug     = &conf.Agent.Debug
	yy_mod_http      = &conf.Agent.Http
	yy_mod_http_addr = &conf.Agent.HttpAddr
	yy_mod_rpc       = &conf.Agent.Rpc
	yy_mod_rpc_addr  = &conf.Agent.RpcAddr
}
;

mod_name:
| text { *yy_mod_name = $1 }
;

agent_mod_item:
   mod_item
 | HOST text { conf.Agent.Host = $2 }
 | IFACE_PREFIX as {
	conf.Agent.IfPre = yy_as
	yy_as = make([]string, 0)
}| INTERVAL num { conf.Agent.Interval = $2 }
 | HANDOFF '{' agent_lb '}'{
 	conf.Agent.Lb.Upstreams = yy_as
	yy_as = make([]string, 0)
 }
 ;

agent_lb:
| agent_lb agent_lb_item ';'

agent_lb_item:
| BATCH num { conf.Agent.Lb.Batch = $2 }
| CONN_TIMEOUT num { conf.Agent.Lb.ConnTimeout = $2 }
| CALL_TIMEOUT num { conf.Agent.Lb.CallTimeout = $2 }
| UPSTREAMS as 
;

ss:
| ss text text ';' { yy_ss[$2] = $3 }
;

as:
| as text { yy_as = append(yy_as, $2) }
;

mod_item:
 | DISABLED bool  { *yy_mod_disable = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG          { *yy_mod_debug = 1 }
 | DEBUG num      { *yy_mod_debug = $2 }
 | HTTP bool      { *yy_mod_http = $2 }
 | HTTP_ADDR text { *yy_mod_http_addr = $2 }
 | RPC bool       { *yy_mod_rpc = $2 }
 | RPC_ADDR text  { *yy_mod_rpc_addr = $2 }
;

lb_mod: lb_start mod_name '{'
| lb_mod lb_mod_item ';'
| lb_mod INCLUDE text ';'  { yy.include($3) }
;

lb_start: HANDOFF {
	conf.Lb     = append(conf.Lb, lb.DefaultLb)
	yy_lb       = &conf.Lb[len(conf.Lb)-1]
	yy_lb_backend = &lb.Backend{}
	yy_mod_name      = &yy_lb.Name
	yy_mod_disable   = &yy_lb.Disabled
	yy_mod_debug     = &yy_lb.Debug
	yy_mod_http      = &yy_lb.Http
	yy_mod_http_addr = &yy_lb.HttpAddr
	yy_mod_rpc       = &yy_lb.Rpc
	yy_mod_rpc_addr  = &yy_lb.RpcAddr
}
;

lb_mod_item:
   mod_item
 | REPLICAS num { yy_lb.Replicas = $2 }
 | CONCURRENCY num { yy_lb.Concurrency = $2 }
 | BACKENDS '{' lb_backends '}'
 ;


lb_backends:
|lb_backends lb_backends_item ';'
;

lb_backends_item:
| TSDB backend_name '{' lb_backend '}' { 
	yy_lb_backend.Type = "tsdb"
	yy_lb.Backends = append(yy_lb.Backends, *yy_lb_backend)
	yy_lb_backend = &lb.Backend{}
}
| FALCON backend_name '{' lb_backend '}'{
	yy_lb_backend.Type = "falcon"
	yy_lb.Backends = append(yy_lb.Backends, *yy_lb_backend)
	yy_lb_backend = &lb.Backend{}
}
;

backend_name:
| text {
	yy_lb_backend.Name = $1
}
;

lb_backend:
| lb_backend lb_backend_item ';'
;

lb_backend_item:
| DISABLED bool { yy_lb_backend.Disabled = $2 }
| SCHED CONSISTENT { yy_lb_backend.Sched = "consistent" }
| BATCH num { yy_lb_backend.Batch = $2 }
| CONN_TIMEOUT num { yy_lb_backend.ConnTimeout = $2 }
| CALL_TIMEOUT num { yy_lb_backend.CallTimeout = $2 }
| UPSTREAMS '{' ss '}' { 
	yy_lb_backend.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", yy_lb_backend.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}
;

backend_mod: backend_start mod_name '{'
| backend_mod backend_mod_item ';'
| backend_mod INCLUDE text ';'  { yy.include($3) }
;

backend_start: BACKEND {
	conf.Backend     = append(conf.Backend, backend.DefaultBackend)
	yy_backend       = &conf.Backend[len(conf.Backend)-1]
	yy_mod_name      = &yy_backend.Name
	yy_mod_disable   = &yy_backend.Disabled
	yy_mod_debug     = &yy_backend.Debug
	yy_mod_http      = &yy_backend.Http
	yy_mod_http_addr = &yy_backend.HttpAddr
	yy_mod_rpc       = &yy_backend.Rpc
	yy_mod_rpc_addr  = &yy_backend.RpcAddr
}
;

backend_mod_item:
  mod_item
| DSN text { yy_backend.Dsn = $2 }
| DB_MAX_IDLE num { yy_backend.DbMaxIdle = $2 }
| SHM_MAGIC_CODE num { yy_backend.Shm.Magic = uint32($2) }
| SHM_KEY_START_ID num { yy_backend.Shm.Key = $2 }
| SHM_SEGMENT_SIZE num { yy_backend.Shm.Size = $2 }
| MIGRATE '{' backend_migrate '}'
| STORAGE RRD '{' backend_storage '}' { yy_backend.Storage.Type = "rrd" }
;

backend_migrate:
| backend_migrate backend_migrate_item ';'

backend_migrate_item:
| DISABLED bool { yy_backend.Migrate.Disabled = $2 }
| CONCURRENCY num { yy_backend.Migrate.Concurrency = $2 }
| REPLICAS num { yy_backend.Migrate.Replicas = $2 }
| CONN_TIMEOUT num { yy_backend.Migrate.ConnTimeout = $2 }
| CALL_TIMEOUT num { yy_backend.Migrate.CallTimeout = $2 }
| UPSTREAMS '{' ss '}' {
	yy_backend.Migrate.Upstreams = yy_ss
	glog.V(4).Infof("upstreams %v yy_ss %v", yy_backend.Migrate.Upstreams, yy_ss)
	yy_ss = make(map[string]string)
}

backend_storage:
| HDISKS as ';' {
	yy_backend.Storage.Hdisks = yy_as
	yy_as = make([]string, 0)
}

%%

