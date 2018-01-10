/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

type Event_i interface {
	Index(s, substr string) int
	Value(key string) int
}

type ExprEvent struct {
}

// index(key, cop=xiaomi)
func (p *ExprEvent) Index(s, substr string) int {
	return -1
}

// value(priority)
func (p *ExprEvent) Value(key string) int {
	return -1
}

func Exec(p Event_i, e *Expr) bool {
	switch e.Type {
	case EXPR_TYPE_RAW:
		return e.Objs[0].(bool)
	case EXPR_TYPE_OR:
		return Exec(p, e.Objs[0].(*Expr)) || Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_AND:
		return Exec(p, e.Objs[0].(*Expr)) && Exec(p, e.Objs[1].(*Expr))
	case EXPR_TYPE_OP_GT:
		return e.Objs[0].(*ExprObj).Exec(p) > e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_GE:
		return e.Objs[0].(*ExprObj).Exec(p) >= e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_EQ:
		return e.Objs[0].(*ExprObj).Exec(p) == e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_LE:
		return e.Objs[0].(*ExprObj).Exec(p) <= e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_LT:
		return e.Objs[0].(*ExprObj).Exec(p) < e.Objs[1].(*ExprObj).Exec(p)
	case EXPR_TYPE_OP_NE:
		return e.Objs[0].(*ExprObj).Exec(p) != e.Objs[1].(*ExprObj).Exec(p)
	default:
		return false
	}
}
