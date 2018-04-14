/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

import "math"

type GetHandle func(isNum bool, num, shift_time int) []float64

type Item interface {
	Get(isNum bool, num, shift_time int) []float64
	Abschange(isNum bool, args []float64, get GetHandle) float64
	Avg(isNum bool, args []float64, get GetHandle) float64
	Band(isNum bool, args []float64, get GetHandle) float64
	Change(isNum bool, args []float64, get GetHandle) float64
	Count(isNum bool, args []float64, get GetHandle) float64
	Delta(isNum bool, args []float64, get GetHandle) float64
	Diff(isNum bool, args []float64, get GetHandle) float64
	Last(isNum bool, args []float64, get GetHandle) float64
	Min(isNum bool, args []float64, get GetHandle) float64
	Max(isNum bool, args []float64, get GetHandle) float64
	Nodata(isNum bool, args []float64, get GetHandle) float64
	Sum(isNum bool, args []float64, get GetHandle) float64
}

type ExprItem struct {
}

func (p *ExprItem) Get(isNum bool, num, shift_time int) []float64 {
	return []float64{}
}

// TODO
// (previous value;last value=abschange)
func (p *ExprItem) Abschange(isNum bool, args []float64, get GetHandle) float64 {
	vs := get(true, 2, 0)
	if len(vs) == 2 {
		return math.Abs(vs[1] - vs[0])
	}
	return UNKNOWN
}

// (sec|#num,<time_shift>)
func (p *ExprItem) Avg(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0
	sum := float64(0)

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	for _, v := range vs {
		sum += v
	}
	return sum / float64(len(vs))
}

// (sec|#num,mask,<time_shift>)
// float64 <-> int64 not safe
func (p *ExprItem) Band(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	return float64(uint64(vs[0]) & uint64(args[0]))
}

// (previous value;last value=change)
func (p *ExprItem) Change(isNum bool, args []float64, get GetHandle) float64 {
	vs := get(true, 2, 0)
	if len(vs) == 2 {
		return vs[1] - vs[0]
	}
	return UNKNOWN
}

// (sec|#num,<pattern>,<operator>,<time_shift>)
// ⇒ count(10m) → number of values for last 10 minutes
// ⇒ count(10m,"error",eq) → number of values for last 10 minutes that equal 'error'
// ⇒ count(10m,12) → number of values for last 10 minutes that equal '12'
// ⇒ count(10m,12,gt) → number of values for last 10 minutes that are over '12'
// ⇒ count(#10,12,gt) → number of values within last 10 values that are over '12'
// ⇒ count(10m,12,gt,1d) → number of values for preceding 10 minutes up to 24 hours ago that were over '12'
// ⇒ count(10m,6/7,band) → number of values for last 10 minutes having '110' (in binary) in the 3 least significant bits.
// ⇒ count(10m,,,1d) → number of values for preceding 10 minutes up to 24 hours ago
func (p *ExprItem) Count(isNum bool, args []float64, get GetHandle) float64 {
	var (
		ret, pattern        float64
		num, op, shift_time int
	)

	switch len(args) {
	case 1:
		num = int(args[0])
	case 2:
		num = int(args[0])
		pattern = args[1]
		op = EXPR_TYPE_OP_EQ
	case 3:
		num = int(args[0])
		pattern = args[1]
		op = int(args[2])
	case 4:
		num = int(args[0])
		pattern = args[1]
		op = int(args[2])
		shift_time = int(args[3])
	default:
		return UNKNOWN
	}

	vs := get(isNum, num, shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	switch op {
	case 0:
		return float64(len(vs))
	case EXPR_TYPE_OP_GT:
		pattern += THRESHOLD
		for _, v := range vs {
			if v >= pattern {
				ret++
			}
		}
	case EXPR_TYPE_OP_GE:
		for _, v := range vs {
			if v >= pattern {
				ret++
			}
		}
	case EXPR_TYPE_OP_EQ:
		a := pattern - THRESHOLD
		b := pattern + THRESHOLD
		for _, v := range vs {
			if v > a && v < b {
				ret++
			}
		}
	case EXPR_TYPE_OP_LE:
		for _, v := range vs {
			if v <= pattern {
				ret++
			}
		}
	case EXPR_TYPE_OP_LT:
		pattern -= THRESHOLD
		for _, v := range vs {
			if v <= pattern {
				ret++
			}
		}
	case EXPR_TYPE_OP_NE:
		a := pattern - THRESHOLD
		b := pattern + THRESHOLD
		for _, v := range vs {
			if v <= a && v >= b {
				ret++
			}
		}
	default:
		return UNKNOWN
	}
	return ret
}

// (sec|#num,<time_shift>)
// Difference between the maximum and minimum values within the defined evaluation period ('max()' minus 'min()').
func (p *ExprItem) Delta(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0
	min := math.MaxFloat64
	max := float64(0)

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	for _, v := range vs {
		if min > v {
			min = v
		}
		if max < v {
			max = v
		}
	}
	return max - min
}

// Checking if last and previous values differ.
// Returns:
// 1 - last and previous values differ
// 0 - otherwise
func (p *ExprItem) Diff(isNum bool, args []float64, get GetHandle) float64 {
	vs := get(true, 2, 0)
	if len(vs) == 2 && vs[0] >= vs[1]-THRESHOLD && vs[0] <= vs[1]+THRESHOLD {
		return 0
	}
	return 1
}

// (sec|#num,<time_shift>)
func (p *ExprItem) Last(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	return vs[0]
}

//min (sec|#num,<time_shift>)
func (p *ExprItem) Min(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0
	ret := math.MaxFloat64

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)
	if len(vs) == 0 {
		return UNKNOWN
	}

	for _, v := range vs {
		if ret > v {
			ret = v
		}
	}
	return ret
}

// (sec|#num,<time_shift>)
func (p *ExprItem) Max(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0
	ret := float64(0)

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)

	for _, v := range vs {
		if ret < v {
			ret = v
		}
	}
	return ret
}

// (sec)
func (p *ExprItem) Nodata(isNum bool, args []float64, get GetHandle) float64 {
	vs := get(false, int(args[0]), 0)
	if len(vs) > 0 {
		return 1
	}
	return 0
}

// (sec|#num,<time_shift>)
func (p *ExprItem) Sum(isNum bool, args []float64, get GetHandle) float64 {
	shift_time := 0
	ret := float64(0)

	if len(args) == 2 {
		shift_time = int(args[1])
	}

	vs := get(isNum, int(args[0]), shift_time)

	for _, v := range vs {
		ret += v
	}
	return ret
}

func Exec(p Item, e *Expr) bool {
	switch e.Type {
	case EXPR_TYPE_RAW:
		return e.Objs[0].(bool)
	case EXPR_TYPE_OR:
		return Exec(p, e.Objs[0].(*Expr)) || Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_AND:
		return Exec(p, e.Objs[0].(*Expr)) && Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_OP_GT:
		return e.Objs[0].(*ExprObj).Invoke(p) >= e.Objs[1].(*ExprObj).Invoke(p)+THRESHOLD
	case EXPR_TYPE_OP_GE:
		return e.Objs[0].(*ExprObj).Invoke(p) >= e.Objs[1].(*ExprObj).Invoke(p)
	case EXPR_TYPE_OP_EQ:
		a, b := e.Objs[0].(*ExprObj).Invoke(p), e.Objs[1].(*ExprObj).Invoke(p)
		return (a > b-THRESHOLD) && (a < b+THRESHOLD)
	case EXPR_TYPE_OP_LE:
		return e.Objs[0].(*ExprObj).Invoke(p) <= e.Objs[1].(*ExprObj).Invoke(p)
	case EXPR_TYPE_OP_LT:
		return e.Objs[0].(*ExprObj).Invoke(p) <= e.Objs[1].(*ExprObj).Invoke(p)-THRESHOLD
	case EXPR_TYPE_OP_NE:
		a, b := e.Objs[0].(*ExprObj).Invoke(p), e.Objs[1].(*ExprObj).Invoke(p)
		return (a <= b-THRESHOLD) || (a >= b+THRESHOLD)
	default:
		return false
	}
}
