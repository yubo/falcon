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
	"math"
)

const (
	ARG_TYPE_DATA = iota
	ARG_TYPE_TIME
)

var (
	EINVAL    = errors.New("Invalid argument")
	THRESHOLD = 0.000001
	UNKNOWN   = math.NaN()
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
	switch p.Type {
	case EXPR_TYPE_RAW:
		return fmt.Sprintf("%s", p.Objs[0].(*bool))
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
	EXPR_OBJ_TYPE_RAW    = iota
	EXPR_OBJ_TYPE_ALL    // all (sec|#num,<time_shift>)
	EXPR_OBJ_TYPE_AVG    // avg (sec|#num,<time_shift>)
	EXPR_OBJ_TYPE_MAX    // max (sec|#num,<time_shift>)
	EXPR_OBJ_TYPE_MIN    // min (sec|#num,<time_shift>)
	EXPR_OBJ_TYPE_SUM    // sum (sec|#num,<time_shift>)
	EXPR_OBJ_TYPE_NODATA // nodata (sec)
	EXPR_OBJ_TYPE_DIFF
	EXPR_OBJ_TYPE_PDIFF
	EXPR_OBJ_TYPE_SIZE
)

var (
	obj_type_name = [EXPR_OBJ_TYPE_SIZE]string{
		"raw",
		"all",
		"avg",
		"max",
		"min",
		"sum",
		"nodata",
		"diff",
		"pdiff",
	}
)

// eg. min(#3)
type ExprObj struct {
	Type uint32
	Args []float64
}

func (p *ExprObj) reduce(typ uint32, argc0, argc1 int) error {
	if len(p.Args) < argc0 || len(p.Args) > argc1 {
		fmt.Printf("%s [%d, %d]\n", p, argc0, argc1)
		return EINVAL
	}
	p.Type |= typ << 1
	return nil
}

func (p *ExprObj) Exec(item ItemInf) float64 {
	var isNum bool
	typ := p.Type >> 1

	if p.Type&0x01 == 0x01 {
		isNum = true
	}

	switch typ {
	case EXPR_OBJ_TYPE_RAW:
		return p.Args[0]
	case EXPR_OBJ_TYPE_MIN:
		return item.Min(isNum, p.Args, item.Get)
	default:
		return 0
	}
}

func floatArrayToString(in []float64) (out string) {
	if len(in) == 0 {
		return
	}

	for _, v := range in {
		out += fmt.Sprintf(", %.0f", v)
	}
	return out[2:]
}

func (p *ExprObj) String() string {
	var isNum, args string
	typ := p.Type >> 1

	if typ == EXPR_OBJ_TYPE_RAW {
		return fmt.Sprintf("%.0f", p.Args[0])
	}

	if p.Type&0x01 == 0x01 {
		isNum = "#"
	}
	args = floatArrayToString(p.Args)

	return fmt.Sprintf("%s(%s%s)", obj_type_name[typ], isNum, args)
}
