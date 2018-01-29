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
%token METRIC AGENT TRANSFER BACKEND

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

as:
	| as text             { yy_as = append(yy_as, $2) }
	| as INCLUDE text ';' { yy.include($3) }
;

conf: ';'
	| ctrl '}' ';'      {
		// end
	 	conf.Ctrl.Set(falcon.APP_CONF_FILE, yy_ss)
		yy_ss = make(map[string]string)
	
		//conf.Name = fmt.Sprintf("ctrl_%s", conf.Name)
		if conf.Host == "" {
			conf.Host, _ = os.Hostname()
		}
	}
;

ctrl:
	'{' {
	 	// begin
		conf = &Ctrl{Name: "ctrl"}
	}| ctrl ctrl_item ';'
;

ctrl_item:
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
	| METRIC '{' as '}' {
	 	conf.Metrics = yy_as
		yy_as = make([]string, 0)
	}| ROOT text { 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	}| agent '}' {
	 	conf.Agent.Set(falcon.APP_CONF_FILE, yy_ss2)
		yy_ss2 = make(map[string]string)
	}| transfer '}' {
	 	conf.Transfer.Set(falcon.APP_CONF_FILE, yy_ss2)
		yy_ss2 = make(map[string]string)
	}| backend '}' {
	 	conf.Backend.Set(falcon.APP_CONF_FILE, yy_ss2)
		yy_ss2 = make(map[string]string)
	}
;

agent: AGENT '{'
	| agent kv_item ';'

transfer: TRANSFER '{'
	| transfer kv_item ';'

backend: BACKEND '{'
       | backend kv_item ';'


kv_item:
	| ROOT text { 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	}
	| text num	{ yy_ss2[$1] = fmt.Sprintf("%d", $2) }
	| text bool	{ yy_ss2[$1] = fmt.Sprintf("%v", $2) }
	| text text	{ yy_ss2[$1] = $2 }
	| INCLUDE text	{ yy.include($2) }


%%

