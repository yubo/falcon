/*
 * Copyright 2016,2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
%{
package config

import (
	"os"
	"fmt"

	"github.com/yubo/falcon"
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
%token <text> TEXT IPA ADDR

%token '{' '}' ';' '*' '(' ')' '+' '-' '<' '>'
%token ON YES OFF NO INCLUDE ROOT PID_FILE LOG HOST DISABLED DEBUG
%token SHARDMAP

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
	| ADDR	{ $$ = string(yy.t) }
	| TEXT	{ $$ = exprText(yy.t) }
;

num:
	NUM			{ $$ = yy.i }
	| '(' num ')'		{ $$ = $2}
	| num '*' num		{ $$ = $1 * $3}
	| num '+' num		{ $$ = $1 + $3}
	| num '<' '<' num	{ $$ = int(uint($1) << uint($4)) }
	| num '>' '>' num	{ $$ = int(uint($1) >> uint($4)) }
;


conf: ';'
	| transfer '}' ';'      {
		// end
	 	conf.Configer.Set(falcon.APP_CONF_FILE, yy_ss)
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
		conf = &Transfer{Name: "transfer"}
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
	| text '=' text {
		if err := os.Setenv($1, $3); err != nil {
	 		yy.Error(err.Error())
		}
	}
	| ROOT text	{ 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	}
	| SHARDMAP shardmap '}' {
		// check shardmap
		for i := 0; i < len(conf.ShardMap); i++ {
			if _, ok := conf.ShardMap[i]; !ok {
				yy.Error(fmt.Sprintf("miss ShardMap[%d]\n", i))
			}
		}
		if len(conf.ShardMap) != falcon.SHARD_NUM {
			yy.Error(fmt.Sprintf("len(ShardMap) %d != SHARD_NUM %d\n", 
				len(conf.ShardMap), falcon.SHARD_NUM))
		}
	}
;

shardmap:
	'{' {
		conf.ShardMap = make(map[int]string)
	}| shardmap shardmap_item ';'
;

shardmap_item:
	| num text	{ conf.ShardMap[$1] = $2 }
	| INCLUDE text	{ yy.include($2) }
;


%%

