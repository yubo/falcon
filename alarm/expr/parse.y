/*
 * Copyright 2016,2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
%{
package expr

%}

%union {
	num int
	expr  *Expr
	expr_obj *ExprObj
	text string
}

%type <text> text
%type <num> num expr_op expr_logic_op 
%type <expr> expr
%type <expr_obj> expr_obj

%token <num> NUM
%token <text> TEXT 

%token '(' ')' '=' '>' '<' '&' '|' '#' '+' '-' '*' '/' ','
%token INDEX VALUE

%%

trigger:
| expr 		{ yy_trigger = $1 }
;

expr:
   num				{ $$ = &Expr{Type: EXPR_TYPE_RAW, Objs: []interface{}{$1}} }
 | expr_obj expr_op expr_obj	{ $$ = &Expr{Type: uint32($2), Objs: []interface{}{$1, $3}} }
 | expr expr_logic_op expr	{ $$ = &Expr{Type: uint32($2), Objs: []interface{}{$1, $3}} }
 | '(' expr ')'			{ $$ = $2 }
;

text:
  TEXT			{ $$ = string(yy.t) }

num:
   NUM			{ $$ = yy.i }
;

expr_op:
   '>'		{ $$ = EXPR_TYPE_OP_GT }
 | '>' '='	{ $$ = EXPR_TYPE_OP_GE }
 | '='		{ $$ = EXPR_TYPE_OP_EQ }
 | '<' '='	{ $$ = EXPR_TYPE_OP_LE }
 | '<'		{ $$ = EXPR_TYPE_OP_LT }
 | '<' '>'	{ $$ = EXPR_TYPE_OP_NE }
;

expr_logic_op:
   '&' '&' { $$ = EXPR_TYPE_AND}
 | '|' '|' { $$ = EXPR_TYPE_OR }
;

expr_obj: num {
 	$$ = &ExprObj{Type: EXPR_OBJ_TYPE_RAW, I: $1}
}| INDEX '(' text ',' text ')' {
	$$ = &ExprObj{Type: EXPR_OBJ_TYPE_INDEX, S0: $3, S1: $5}
}| VALUE '(' text ')' {
	$$ = &ExprObj{Type: EXPR_OBJ_TYPE_VALUE, S0: $3}
}
;

%%

