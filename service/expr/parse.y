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
	float float64
	b bool
	expr  *Expr
	expr_obj *ExprObj
}

%type <b> bool
%type <num> expr_op expr_logic_op
%type <float> float
%type <expr> expr
%type <expr_obj> expr_obj

%token <num> NUM

%token '(' ')' '=' '>' '<' '&' '|' '#' '+' '-' '*' '/' ','
%token TRUE FALSE
%token ALL MAX MIN AVG SUM DIFF PDIFF

%%

trigger:
| expr 		{ yy_trigger = $1 }
;

expr:
   bool				{ $$ = &Expr{Type: EXPR_TYPE_RAW << 1, Objs: []interface{}{$1}} }
 | expr_obj expr_op expr_obj	{ $$ = &Expr{Type: uint32($2), Objs: []interface{}{$1, $3}} }
 | expr expr_logic_op expr	{ $$ = &Expr{Type: uint32($2), Objs: []interface{}{$1, $3}} }
 | '(' expr ')'			{ $$ = $2 }
;

bool:
   TRUE			{ $$ = true }
 | FALSE		{ $$ = false }
;

float:
   NUM { $$ = yy.f }
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

expr_obj: float {
 	$$ = &ExprObj{Type: EXPR_OBJ_TYPE_RAW, Args: []float64{$1}}
}| MIN obj_args ')' {
	yy.err = yy_obj.reduce(EXPR_OBJ_TYPE_MIN, 2, 3)
 	$$ = yy_obj
}
;

obj_args: '(' {
	yy_obj = &ExprObj{}
}| obj_args arg0_type func_args
;


arg0_type:
 | '#' { yy_obj.Type = 1 }
;

func_args: 
   float {
	yy_obj.Args = append(yy_obj.Args, $1)
}| func_args ',' float {
	yy_obj.Args = append(yy_obj.Args, $3)
}
;

%%

