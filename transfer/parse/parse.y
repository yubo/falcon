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

	"github.com/yubo/falcon/transfer/config"
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
%token BACKEND UPSTREAM

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
	| transfer '}' ';'      {
		// end
	 	conf.Configer.Set(fconfig.APP_CONF_FILE, yy_ss)
		yy_ss = make(map[string]string)
	
		//conf.Name = fmt.Sprintf("transfer_%s", conf.Name)
		if conf.Host == "" {
			conf.Host, _ = os.Hostname()
		}
	}
;

transfer:
	'{' {
	 	// begin
		conf = &config.ConfTransfer{Name: "transfer"}
	}| transfer transfer_item ';'
;

transfer_item:
	| DISABLED bool	{ conf.Disabled = $2 }
	| HOST text	{ conf.Host = $2 }
	| DEBUG		{ conf.Debug = 1 }
	| DEBUG num	{ conf.Debug = $2 }
	| text num	{ yy_ss[$1] = fmt.Sprintf("%d", $2) }
	| text bool	{ yy_ss[$1] = fmt.Sprintf("%v", $2) }
	| INCLUDE text	{ yy.include($2) }
	| text text	{ yy_ss[$1] = $2 }
	| ROOT text { 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	}
	| BACKEND '{' transfer_backend '}'
;

transfer_backend:
	| transfer_backend transfer_backend_item ';'
;

transfer_backend_item:
	| text text '{' transfer_backend_obj '}' { 
		yy_transfer_backend.Type = $1
		yy_transfer_backend.Name = $2
		if !yy_transfer_backend.Disabled || yy.debug {
			conf.Backend = append(conf.Backend, *yy_transfer_backend)
		}
		yy_transfer_backend = &config.TransferBackend{}
	}
;

transfer_backend_obj:
	| transfer_backend_obj transfer_backend_obj_item ';'
;

transfer_backend_obj_item:
	| DISABLED bool { yy_transfer_backend.Disabled = $2 }
	| UPSTREAM '{' ss2 '}' { 
		yy_transfer_backend.Upstream = yy_ss2
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

