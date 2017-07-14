/*
 * Copyright 2016,2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
%{
package parse

import (
	"os"
	"fmt"

	"github.com/yubo/falcon/backend/config"
	fconfig "github.com/yubo/falcon/config"
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
%token MIGRATE UPSTREAM

%%

config: 
| config conf
;

bool:
	ON	{ $$ = true }
	| YES	{ $$ = true }
	| OFF	{ $$ = false }
	| NO	{ $$ = false }
	|	{ $$ = true }
;

text:
	IPA	{ $$ = string(yy.t) }
	| TEXT	{ $$ = exprText(yy.t) }
;

num:
	NUM { $$ = yy.i }
;


conf: ';'
	| backend '}' ';'      {
		// end
	 	conf.Configer.Set(fconfig.APP_CONF_FILE, yy_ss)
		yy_ss = make(map[string]string)
	
		//conf.Name = fmt.Sprintf("backend_%s", conf.Name)
		if conf.Host == "" {
			conf.Host, _ = os.Hostname()
		}
	}
;

backend:
	'{' {
	 	// begin
		conf = &config.ConfBackend{Name: "backend"}
	}| backend backend_item ';'
;

backend_item:
	| DISABLED bool	{ conf.Disabled = $2 }
 	| MIGRATE 	'{' backend_migrate '}'
	| HOST text	{ conf.Host = $2 }
	| DEBUG		{ conf.Debug = 1 }
	| DEBUG num	{ conf.Debug = $2 }
	| text num	{ yy_ss[$1] = fmt.Sprintf("%d", $2) }
	| text bool	{ yy_ss[$1] = fmt.Sprintf("%v", $2) }
	| INCLUDE text	{ yy.include($2) }
	| text text	{ yy_ss[$1] = $2 }
	| ROOT text	{ 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	};

backend_migrate:
	| backend_migrate backend_migrate_item ';'

backend_migrate_item:
	| DISABLED bool { conf.Migrate.Disabled = $2 }
	| UPSTREAM '{' ss2 '}' {
		conf.Migrate.Upstream = yy_ss2
		yy_ss2 = make(map[string]string)
	}
;

ss2:
	| ss2 text text ';'    { yy_ss2[$2] = $3 }
	| ss2 text num ';'     { yy_ss2[$2] = fmt.Sprintf("%d", $3) }
	| ss2 text bool ';'    { yy_ss2[$2] = fmt.Sprintf("%v", $3) }
	| ss2 INCLUDE text ';' { yy.include($3) }
;

%%

