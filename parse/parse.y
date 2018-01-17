/*
 * Copyright 2016,2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
%{

package parse

import (
	"fmt"
	"os"

	"github.com/yubo/falcon"
)


%}

%union {
	num int
	text string
}

%type <num> num
%type <text> text module_text

%token <num> NUM
%token <text> TEXT IPA MODULE_TEXT ADDR

%token '{' '}' ';' '*' '+' '>' '<' '+' '(' ')'
%token INCLUDE ROOT PID_FILE LOG HOST DISABLED DEBUG

%%

config: 
	| config conf_item
;

num:
	NUM			{ $$ = yy.i }
	| '(' num ')'		{ $$ = $2}
	| num '*' num		{ $$ = $1 * $3}
	| num '+' num		{ $$ = $1 + $3}
	| num '<' '<' num	{ $$ = int(uint($1) << uint($4)) }
	| num '>' '>' num	{ $$ = int(uint($1) >> uint($4)) }
;

text:
	IPA	{ $$ = string(yy.t) }
	| ADDR	{ $$ = string(yy.t) }
	| TEXT	{ $$ = exprText(yy.t) }
;

module_text:
  MODULE_TEXT { $$ = string(yy.t) }	
;

conf_item: ';'
	| PID_FILE text ';' { conf.PidFile = $2 }
	| text '=' text ';' {
		if err := os.Setenv($1, $3); err != nil {
	 		yy.Error(err.Error())
		}
	}
	| LOG text num ';'  {
		conf.Log = $2
		conf.Logv = $3
	}
	| INCLUDE text ';'  { yy.include($2) }
	| ROOT text ';'     { 
	 	if err := os.Chdir($2); err != nil {
	 		yy.Error(err.Error())
	 	}
	}
	| module ';'
;

module: text '{' {
	 	yy_module = &yyModule {
			file: yy.ctx.file,
			lino: yy.ctx.lino,
		}
		if m, ok := falcon.Modules[$1]; ok {
			yy_module_parse = m.Parse
		} else {
			yy.Error(fmt.Sprintf("module [%s] not exists", $1))
		}

	}
	| module module_text '}' {
		conf.Conf  = append(conf.Conf, yy_module_parse([]byte(fmt.Sprintf("{ %s };", $2)),
			yy_module.file, yy_module.lino))
		yy_module = nil
	}
;

%%

