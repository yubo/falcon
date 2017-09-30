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
	b bool
}

%type <b> bool
%type <text> text
%type <num> num

%token <num> NUM
%token <text> TEXT IPA

%token '{' '}' ';'
%token ON YES OFF NO INCLUDE ROOT PID_FILE LOG HOST DISABLED DEBUG
%token MODULE

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
	NUM	{ $$ = yy.i }
;

conf: ';'
	| PID_FILE text ';' { conf.PidFile = $2 }
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
	| module '}' '}' {
		p1, _ := falcon.PreByte(yy.ctx.text, yy.ctx.pos)
		yy.ctx.text[p1] = ';'

		conf.Conf  = append(conf.Conf, yy_module_parse(
		yy.ctx.text[yy_module.pos : yy.ctx.pos],
		yy_module.file, yy_module.lino, yy_module.debug))
	}
;

module: text '{' {
	 	yy_module = &yyModule {
			level: 1,
			debug: yy.debug,
			file: yy.ctx.file,
			lino: yy.ctx.lino,
			pos: yy.ctx.pos - 1,
		}
		if m, ok := falcon.Modules[$1]; ok {
			yy_module_parse = m.Parse
		} else {
			yy.Error(fmt.Sprintf("module [%s] not exists", $1))
		}

	}
	| module ';' {
		if (yy_module.level == 0) {
			p1, c1 := falcon.PreByte(yy.ctx.text, yy.ctx.pos)
			p2, c2 := falcon.PreByte(yy.ctx.text, p1 - 1)
			if c1 == ';' && c2 == '}' {
				yy.ctx.text[p1] = '}'
			}
			yy.ctx.pos = p2 - 1
		}
	}
	| module '{' { yy_module.level++ }
	| module '}' { yy_module.level-- }
	| module text
	| module num 
	| module bool
	| module ROOT
	| module LOG
	| module HOST
	| module DEBUG
	| module INCLUDE
	| module PID_FILE
	| module DISABLED
%%

