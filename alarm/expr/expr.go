/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ARG_TYPE_DATA = iota
	ARG_TYPE_TIME
)

var (
	EINVAL = errors.New("Invalid argument")
)

const (
	EXPR_TYPE_RAW   = iota
	EXPR_TYPE_OR    // || Logical OR
	EXPR_TYPE_AND   // && Logical AND
	EXPR_TYPE_OP_GT // >  More than. The operator is defined as: A>B ⇔ (A>=B+0.000001)
	EXPR_TYPE_OP_GE // >= greater than or equal
	EXPR_TYPE_OP_EQ // =  Is equal. The operator is defined as: A=B ⇔ (A>B-0.000001) & (A<B+0.000001)
	EXPR_TYPE_OP_LE // <= less than or equal
	EXPR_TYPE_OP_LT // <  Less than. The operator is defined as: A<B ⇔ (A<=B-0.000001)
	EXPR_TYPE_OP_NE // <> Not equal. The operator is defined as: A#B ⇔ (A<=B-0.000001) | (A>=B+0.000001)

)

// eg. min(#3) > 3
type Expr struct {
	Type uint32
	Objs []interface{}
}

func (p *Expr) String() string {
	if len(p.Objs) == 0 {
		return ""
	}
	switch p.Type {
	case EXPR_TYPE_RAW:
		return fmt.Sprintf("%s", p.Objs[0].(*int))
	case EXPR_TYPE_OR:
		return fmt.Sprintf("%s || %s",
			p.Objs[0].(*Expr).String(),
			p.Objs[1].(*Expr).String())
	case EXPR_TYPE_AND:
		return fmt.Sprintf("%s && %s",
			p.Objs[0].(*Expr).String(),
			p.Objs[1].(*Expr).String())
	case EXPR_TYPE_OP_GT:
		return fmt.Sprintf("%s > %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	case EXPR_TYPE_OP_GE:
		return fmt.Sprintf("%s >= %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	case EXPR_TYPE_OP_EQ:
		return fmt.Sprintf("%s = %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	case EXPR_TYPE_OP_LE:
		return fmt.Sprintf("%s <= %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	case EXPR_TYPE_OP_LT:
		return fmt.Sprintf("%s < %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	case EXPR_TYPE_OP_NE:
		return fmt.Sprintf("%s <> %s",
			p.Objs[0].(*ExprObj).String(),
			p.Objs[1].(*ExprObj).String())
	}
	return ""
}

func (p *Expr) Json() string {
	s, _ := json.Marshal(p)
	return string(s)
}

const (
	EXPR_OBJ_TYPE_RAW = iota
	EXPR_OBJ_TYPE_INDEX
	EXPR_OBJ_TYPE_VALUE
	EXPR_OBJ_TYPE_SIZE
)

var (
	obj_type_name = [EXPR_OBJ_TYPE_SIZE]string{
		"raw",
		"index",
		"value",
	}

	obj_type_map = map[string]int{
		"raw":   EXPR_OBJ_TYPE_RAW,
		"index": EXPR_OBJ_TYPE_INDEX,
		"value": EXPR_OBJ_TYPE_VALUE,
	}

	obj_argc_thresholds = [EXPR_OBJ_TYPE_SIZE][2]int{
		{0, 0}, // "raw"
		{2, 2}, // "index",
		{1, 1}, // "value",
	}
)

// eg. min(#3)
type ExprObj struct {
	Type uint32
	I    int
	S0   string
	S1   string
}

func (p *ExprObj) Invoke(e Event_i) int {
	switch p.Type {
	case EXPR_OBJ_TYPE_RAW:
		return p.I
	case EXPR_OBJ_TYPE_INDEX:
		return e.Index(p.S0, p.S1)
	case EXPR_OBJ_TYPE_VALUE:
		return e.Value(p.S0)
	default:
		return 0
	}
}

func (p *ExprObj) String() string {
	switch p.Type {
	case EXPR_OBJ_TYPE_RAW:
		return fmt.Sprintf("%d", p.I)
	case EXPR_OBJ_TYPE_INDEX:
		return fmt.Sprintf("index(%s,%s)", p.S0, p.S1)
	case EXPR_OBJ_TYPE_VALUE:
		return fmt.Sprintf("value(%s)", p.S0)
	default:
		return ""
	}
}
