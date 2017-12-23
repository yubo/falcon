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
%token SHAREMAP

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
	| ROOT text	{ 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	}
	| SHAREMAP sharemap '}' {
		// check sharemap
		conf.ShareCount = len(conf.ShareMap)	
		for i := 0; i < conf.ShareCount; i++ {
			if _, ok := conf.ShareMap[i]; !ok {
				yy.Error(fmt.Sprintf("miss shareMap[%d]\n", i))
			}
		}
	}
;

sharemap:
	'{' {
		conf.ShareMap = make(map[int]string)
	}| sharemap sharemap_item ';'
;

sharemap_item:
	| num text	{ conf.ShareMap[$1] = $2 }
	| INCLUDE text	{ yy.include($2) }
;


%%

