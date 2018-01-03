/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

type ItemInf interface {
	//Abschange(isNum bool, args []float64 ) float64
	//Avg(isNum bool, args []float64 ) float64
	//Band(isNum bool, args []float64 ) float64
	//Change(isNum bool, args []float64 ) float64
	//Count(isNum bool, args []float64 ) float64
	//Delta(isNum bool, args []float64 ) float64
	//Diff(isNum bool, args []float64 ) float64
	//Last(isNum bool, args []float64 ) float64
	Min(isNum bool, args []float64, get GetHandle) float64
	//Max(isNum bool, args []float64 ) float64
	//Nodata(isNum bool, args []float64 ) float64
	//Sum(isNum bool, args []float64 ) float64
	Get(isNum bool, last, shift_time int) []float64
}

type GetHandle func(isNum bool, last, shift_time int) []float64

type ExprItem struct {
}

func (p *ExprItem) Get(isNum bool, last, shift_time int) []float64 {
	return []float64{}
}

//min (sec|#num,<time_shift>)
func (p *ExprItem) Min(isNum bool, args []float64,
	get GetHandle) float64 {
	var (
		ret float64
		i   int
	)

	last := int(args[0])
	shift_time := 0

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, last, shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	for ret, i = vs[0], 1; i < len(vs); i++ {
		if ret > vs[i] {
			ret = vs[i]
		}
	}
	return ret
}

func Exec(p ItemInf, e *Expr) bool {
	switch e.Type {
	case EXPR_TYPE_RAW:
		return e.Objs[0].(bool)
	case EXPR_TYPE_OR:
		return Exec(p, e.Objs[0].(*Expr)) || Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_AND:
		return Exec(p, e.Objs[0].(*Expr)) && Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_OP_GT:
		return e.Objs[0].(*ExprObj).Exec(p) >= e.Objs[1].(*ExprObj).Exec(p)+THRESHOLD
	case EXPR_TYPE_OP_GE:
		return e.Objs[0].(*ExprObj).Exec(p) >= e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_EQ:
		a, b := e.Objs[0].(*ExprObj).Exec(p), e.Objs[1].(*ExprObj).Exec(p)
		return (a > b-THRESHOLD) && (a < b+THRESHOLD)
	case EXPR_TYPE_OP_LE:
		return e.Objs[0].(*ExprObj).Exec(p) <= e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_LT:
		return e.Objs[0].(*ExprObj).Exec(p) <= e.Objs[1].(*ExprObj).Exec(p)-THRESHOLD
	case EXPR_TYPE_OP_NE:
		a, b := e.Objs[0].(*ExprObj).Exec(p), e.Objs[1].(*ExprObj).Exec(p)
		return (a <= b-THRESHOLD) || (a >= b+THRESHOLD)
	default:
		return false
	}
}
