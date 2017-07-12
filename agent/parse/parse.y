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

	"github.com/yubo/falcon/agent/config"
	"github.com/yubo/falcon/utils"
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
%token METRIC

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

as:
	| as text             { yy_as = append(yy_as, $2) }
	| as INCLUDE text ';' { yy.include($3) }
;

conf: ';'
	| agent '}' ';'      {
		// end
	 	conf.Configer.Set(utils.APP_CONF_FILE, yy_ss2)
		yy_ss2 = make(map[string]string)
	
		conf.Name = fmt.Sprintf("ctrl_%s", conf.Name)
		if conf.Host == "" {
			conf.Host, _ = os.Hostname()
		}
	}
;

agent:
	'{' {
	 	// begin
		conf = &config.ConfAgent{Name: "agent"}
		conf.Configer.Set(utils.APP_CONF_DEFAULT, config.ConfDefault)
	}| agent agent_item ';'
;

agent_item:
	| DISABLED bool	{ conf.Disabled = $2 }
	| HOST text	{ conf.Host = $2 }
	| DEBUG		{ conf.Debug = 1 }
	| DEBUG num	{ conf.Debug = $2 }
	| text num	{ yy_ss2[$1] = fmt.Sprintf("%d", $2) }
	| text bool	{ yy_ss2[$1] = fmt.Sprintf("%v", $2) }
	| INCLUDE text	{ yy.include($2) }
	| text text	{ yy_ss2[$1] = $2 }
	| ROOT text { 
		if err := os.Chdir($2); err != nil {
			yy.Error(err.Error())
		}
	};

%%

