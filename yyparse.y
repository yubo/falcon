// Copyright 2016 falcon Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
%{

package falcon

import (
	"os"
	"fmt"
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
%token ON YES OFF NO INCLUDE ROOT PID_FILE LOG HOST DISABLED DEBUG
%token CTRL AGENT TRANSFER BACKEND
%token UPSTREAM METRIC MIGRATE

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
| ss text text ';'    { yy_ss[$2] = $3 }
| ss text num ';'     { yy_ss[$2] = fmt.Sprintf("%d", $3) }
| ss text bool ';'    { yy_ss[$2] = fmt.Sprintf("%v", $3) }
| ss INCLUDE text ';' { yy.include($3) }
;

as:
| as text             { yy_as = append(yy_as, $2) }
| as INCLUDE text ';' { yy.include($3) }
;

conf: ';'
 | PID_FILE text ';' { conf.pidFile = $2 }
 | LOG text num ';'  {
 conf.log = $2
 conf.logv = $3
}| INCLUDE text ';'  { yy.include($2) }
 | ROOT text ';'     { 
 	if err := os.Chdir($2); err != nil {
 		yy.Error(err.Error())
 	}
}| ctrl_mod '}'      {
 	yy_ctrl.Ctrl.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_ctrl.Name = fmt.Sprintf("ctrl_%s", yy_ctrl.Name)
	if yy_ctrl.Host == ""{
		yy_ctrl.Host, _ = os.Hostname()
	}

	if !yy_ctrl.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_ctrl)
	}
}| agent_mod '}'      {
 	yy_agent.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_agent.Name = fmt.Sprintf("agent_%s", yy_agent.Name)
	if yy_agent.Host == ""{
		yy_agent.Host, _ = os.Hostname()
	}

	if !yy_agent.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_agent)
	}
}| transfer_mod '}'   {
 	yy_transfer.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_transfer.Name = fmt.Sprintf("transfer_%s", yy_transfer.Name)
	if yy_transfer.Host == ""{
		yy_transfer.Host, _ = os.Hostname()
	}
	if !yy_transfer.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_transfer)
	}
}| backend_mod '}'      {
 	yy_backend.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_backend.Name = fmt.Sprintf("backend_%s", yy_backend.Name)
	if yy_backend.Host == ""{
		yy_backend.Host, _ = os.Hostname()
	}
	if !yy_backend.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_backend)
	}
}
;

////////////////////// ctrl /////////////////////////
ctrl_mod: CTRL text '{' {
	yy_ctrl      = &ConfCtrl{}
	yy_ctrl.Ctrl.Set(APP_CONF_DEFAULT, ConfDefault["ctrl"])
	yy_ctrl.Name = $2
}| ctrl_mod ctrl_mod_item ';'
;

ctrl_mod_item:
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DISABLED bool   { yy_ctrl.Disabled = $2 }
 | DEBUG           { yy_ctrl.Debug = 1 }
 | DEBUG num       { yy_ctrl.Debug = $2 }
 | HOST text       { yy_ctrl.Host = $2 }
 | METRIC '{' as '}' {
 	yy_ctrl.Metrics = yy_as
	yy_as = make([]string, 0)
}| text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }
;

////////////////////// agent /////////////////////////
agent_mod: AGENT text '{' {
	yy_agent      = &ConfAgent{}
	yy_agent.Configer.Set(APP_CONF_DEFAULT, ConfDefault["agent"])
	yy_agent.Name = $2
}| agent_mod agent_mod_item ';'
;

agent_mod_item:
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DISABLED bool   { yy_agent.Disabled = $2 }
 | DEBUG           { yy_agent.Debug = 1 }
 | DEBUG num       { yy_agent.Debug = $2 }
 | HOST text       { yy_agent.Host = $2 }
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | UPSTREAM text   { yy_ss2["upstream"] = $2 }
 | INCLUDE text    { yy.include($2) }
;


////////////////////// transfer  /////////////////////////
transfer_mod: TRANSFER text '{' {
	yy_transfer      = &ConfTransfer{}
	yy_transfer.Configer.Set(APP_CONF_DEFAULT, ConfDefault["transfer"])
	yy_transfer.Name = $2
}| transfer_mod transfer_mod_item ';'
;

transfer_mod_item:
 | DISABLED bool   { yy_transfer.Disabled = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG           { yy_transfer.Debug = 1 }
 | DEBUG num       { yy_transfer.Debug = $2 }
 | HOST text       { yy_transfer.Host = $2 }
 | BACKEND '{' transfer_backend '}'
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }

;

transfer_backend:
| transfer_backend transfer_backend_item ';'
;

transfer_backend_item:
| text text '{' transfer_backend_obj '}' { 
	yy_transfer_backend.Type = $1
	yy_transfer_backend.Name = $2
	if !yy_transfer_backend.Disabled || yy.debug {
		yy_transfer.Backend = append(yy_transfer.Backend, *yy_transfer_backend)
	}
	yy_transfer_backend = &TransferBackend{}
}
;

transfer_backend_obj:
| transfer_backend_obj transfer_backend_obj_item ';'
;

transfer_backend_obj_item:
| DISABLED bool { yy_transfer_backend.Disabled = $2 }
| UPSTREAM '{' ss '}' { 
	yy_transfer_backend.Upstream = yy_ss
	yy_ss = make(map[string]string)
}
;

////////////////////// backend  /////////////////////////
backend_mod: BACKEND text '{'   {
	yy_backend      = &ConfBackend{}
	yy_backend.Configer.Set(APP_CONF_DEFAULT, ConfDefault["backend"])
	yy_backend.Name = $2
}| backend_mod backend_mod_item ';'
;


backend_mod_item:
 | DISABLED bool   { yy_backend.Disabled = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG           { yy_backend.Debug = 1 }
 | DEBUG num       { yy_backend.Debug = $2 }
 | HOST text       { yy_backend.Host = $2 }
 | MIGRATE '{' backend_migrate '}'
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }
;

backend_migrate:
| backend_migrate backend_migrate_item ';'

backend_migrate_item:
| DISABLED bool { yy_backend.Migrate.Disabled = $2 }
| UPSTREAM '{' ss '}' {
	yy_backend.Migrate.Upstream = yy_ss
	yy_ss = make(map[string]string)
}
;

%%

