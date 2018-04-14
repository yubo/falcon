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
	if len(p.Objs) == 0 {
		return ""
	}
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
	EXPR_OBJ_TYPE_RAW = iota
	EXPR_OBJ_TYPE_ABSCHANGE
	EXPR_OBJ_TYPE_AVG
	EXPR_OBJ_TYPE_BAND
	EXPR_OBJ_TYPE_CHANGE
	EXPR_OBJ_TYPE_COUNT
	EXPR_OBJ_TYPE_DELTA
	EXPR_OBJ_TYPE_DIFF
	EXPR_OBJ_TYPE_LAST
	EXPR_OBJ_TYPE_MIN
	EXPR_OBJ_TYPE_MAX
	EXPR_OBJ_TYPE_NODATA
	EXPR_OBJ_TYPE_SUM
	EXPR_OBJ_TYPE_SIZE
)

var (
	obj_type_name = [EXPR_OBJ_TYPE_SIZE]string{
		"raw",
		"abschange",
		"avg",
		"band",
		"change",
		"count",
		"delta",
		"diff",
		"last",
		"min",
		"max",
		"nodata",
		"sum",
	}

	obj_type_map = map[string]int{
		"raw":       EXPR_OBJ_TYPE_RAW,
		"abschange": EXPR_OBJ_TYPE_ABSCHANGE,
		"avg":       EXPR_OBJ_TYPE_AVG,
		"band":      EXPR_OBJ_TYPE_BAND,
		"change":    EXPR_OBJ_TYPE_CHANGE,
		"count":     EXPR_OBJ_TYPE_COUNT,
		"delta":     EXPR_OBJ_TYPE_DELTA,
		"diff":      EXPR_OBJ_TYPE_DIFF,
		"last":      EXPR_OBJ_TYPE_LAST,
		"min":       EXPR_OBJ_TYPE_MIN,
		"max":       EXPR_OBJ_TYPE_MAX,
		"nodata":    EXPR_OBJ_TYPE_NODATA,
		"sum":       EXPR_OBJ_TYPE_SUM,
	}

	obj_argc_thresholds = [EXPR_OBJ_TYPE_SIZE][2]int{
		{0, 0}, // "raw",
		{0, 0}, // "abschange",
		{1, 2}, // "avg",
		{1, 2}, // "band",
		{0, 0}, // "change",
		{1, 4}, // "count",
		{1, 2}, // "delta",
		{0, 0}, // "diff",
		{1, 2}, // "last",
		{1, 2}, // "min",
		{1, 2}, // "max",
		{1, 1}, // "nodata",
		{1, 2}, // "sum",
	}
)

// eg. min(#3)
type ExprObj struct {
	Type uint32
	Args []float64
}

func (p *ExprObj) reduce(name string) error {

	typ, ok := obj_type_map[name]
	if !ok {
		return EINVAL
	}

	if len(p.Args) < obj_argc_thresholds[typ][0] ||
		len(p.Args) > obj_argc_thresholds[typ][1] {
		return EINVAL
	}
	p.Type |= uint32(typ) << 1
	return nil
}

func (p *ExprObj) Invoke(item Item) float64 {
	var isNum bool
	typ := p.Type >> 1

	if p.Type&0x01 == 0x01 {
		isNum = true
	}

	switch typ {
	case EXPR_OBJ_TYPE_RAW:
		return p.Args[0]
	case EXPR_OBJ_TYPE_ABSCHANGE:
		return item.Abschange(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_AVG:
		return item.Avg(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_BAND:
		return item.Band(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_CHANGE:
		return item.Change(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_COUNT:
		return item.Count(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_DELTA:
		return item.Delta(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_DIFF:
		return item.Diff(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_LAST:
		return item.Last(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_MIN:
		return item.Min(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_MAX:
		return item.Max(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_NODATA:
		return item.Nodata(isNum, p.Args, item.Get)
	case EXPR_OBJ_TYPE_SUM:
		return item.Sum(isNum, p.Args, item.Get)
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
