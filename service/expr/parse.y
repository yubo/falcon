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
	text string
}

%type <b> bool
%type <text> text
%type <num> expr_op expr_logic_op time
%type <float> float
%type <expr> expr
%type <expr_obj> expr_obj

%token <num> NUM
%token <text> TEXT 

%token '(' ')' '=' '>' '<' '&' '|' '#' '+' '-' '*' '/' ','
%token TRUE FALSE
%token COUNT

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


text:
  TEXT			{ $$ = string(yy.t) }

float:
   NUM			{ $$ = yy.f }
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
}| COUNT count_args ')' {
	yy.err = yy_obj.reduce("count")
	$$ = yy_obj
}| text obj_args ')' {
	yy.err = yy_obj.reduce($1)
	$$ = yy_obj
}
;

// count(10m,12,gt,1d) → number of values for preceding 10 minutes up to 24 hours ago that were over '12'
count_args: '(' {
	yy_obj = &ExprObj{}
}| count_args arg0 ',' float ',' expr_op ',' time {
	yy_obj.Args = append(yy_obj.Args, $4, float64($6), float64($8))
}| count_args arg0 ',' float ',' expr_op {
	yy_obj.Args = append(yy_obj.Args, $4, float64($6))
}| count_args arg0 ',' float {
	yy_obj.Args = append(yy_obj.Args, $4)
}| count_args arg0
;

obj_args: '(' {
	yy_obj = &ExprObj{}
}| obj_args arg0 func_args
;

arg0:
   float {
   yy_obj.Args = []float64{$1}
}| '#' float {
   yy_obj.Type = 1
   yy_obj.Args = []float64{$2}
}| time {
   yy_obj.Args = []float64{float64($1)}
}
;

time:
   float {
   $$ = int($1)
}| float text {
   switch $2 {
   case "m":
   	$$ = int($1) * 60
   case "h":
   	$$ = int($1) * 60 * 60
   case "d":
   	$$ = int($1) * 60 * 60 * 24
   default:
   	yy.err = EINVAL
   }
}

func_args: 
| func_args ',' float {
	yy_obj.Args = append(yy_obj.Args, $3)
}
;

%%

