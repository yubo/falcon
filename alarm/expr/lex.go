/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package expr

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/golang/glog"
)

const (
	eof           = 0
	MAX_CTX_LEVEL = 16
	MODULE_NAME   = "\x1B[32m[EXPR_PARSE]\x1B[0m "
)

var (
	yy         *yyLex
	yy_obj     *ExprObj
	yy_trigger *Expr
	yy_ss      = make(map[string]string)
	yy_ss2     = make(map[string]string)
	yy_as      = make([]string, 0)
	f_num      = regexp.MustCompile(`^(\+|-)?\d+`)
	f_keyword  = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]+`)
	f_word     = regexp.MustCompile(`(^[a-zA-z"][^\),]+)`)
	keywords   = map[string]int{
		// bool
		//"true":  TRUE,
		//"false": FALSE,
		// func
		"index": INDEX,
		"value": VALUE,
	}
)

type yyCtx struct {
	text []byte
	pos  int
	lino int
}

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and falcon.Error(string).
type yyLex struct {
	ctxData [MAX_CTX_LEVEL]yyCtx
	ctxL    int
	ctx     *yyCtx
	t       []byte
	i       int
	err     error
}

func prefix(a, b []byte) bool {
	if len(a) < len(b) {
		return false
	}

	if len(a) == len(b) {
		return bytes.Equal(a, b)
	}

	return bytes.Equal(a[:len(b)], b)
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (p *yyLex) Lex(yylval *yySymType) int {
	var f []byte

begin:
	text := p.ctx.text[p.ctx.pos:]
	for {
		if yy.err != nil {
			return -1
		}
		if p.ctx.pos == len(p.ctx.text) {
			if p.ctxL > 0 {
				p.ctxL--
				p.ctx = &p.ctxData[p.ctxL]
				goto begin
			}
			return eof
		}

		for text[0] == ' ' || text[0] == '\t' || text[0] == '\n' {
			p.ctx.pos += 1
			if p.ctx.pos == len(p.ctx.text) {
				glog.V(5).Infof(MODULE_NAME+"ctx level %d", p.ctxL)
				if p.ctxL > 0 {
					p.ctxL--
					p.ctx = &p.ctxData[p.ctxL]
					goto begin
				}
				return eof
			}
			if text[0] == '\n' {
				p.ctx.lino++
			}
			text = p.ctx.text[p.ctx.pos:]
		}

		f = f_num.Find(text)
		if f != nil {
			p.t = f[:]
			p.ctx.pos += len(f)
			i64, _ := strconv.ParseInt(string(f), 0, 0)
			p.i = int(i64)
			glog.V(5).Infof(MODULE_NAME+"return NUM %d\n", p.i)
			return NUM
		}

		// find keyword
		f = f_keyword.Find(text)
		if f != nil {
			if val, ok := keywords[string(f)]; ok {
				p.ctx.pos += len(f)
				glog.V(5).Infof(MODULE_NAME+"find %s return %d\n", string(f), val)
				return val
			}
		}

		if bytes.IndexByte([]byte(`#={}:;,()+*/%<>~\[\]?!\|-&`), text[0]) != -1 {
			p.ctx.pos++
			glog.V(5).Infof(MODULE_NAME+"return '%c'\n", int(text[0]))
			//fmt.Printf(MODULE_NAME+"return '%c'\n", int(text[0]))
			return int(text[0])
		}

		// find text
		f = f_word.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			if f[0] == '"' {
				p.t = f[1 : len(f)-1]
			} else {
				p.t = f[:]
			}
			glog.V(5).Infof(MODULE_NAME+"return TEXT(%s)", string(p.t))
			//fmt.Printf(MODULE_NAME+"return TEXT(%s)\n", string(p.t))
			return TEXT
		}

		p.Error(fmt.Sprintf("unknown character %c", text[0]))
	}
}

// The parser calls this method on a parse error.
func (p *yyLex) Error(s string) {
	p.ctx.pos--
	err := fmt.Sprintf("*%c*", p.ctx.text[p.ctx.pos])

	for pos := p.ctx.pos - 1; pos >= 0; pos-- {
		err = fmt.Sprintf("%c%s", p.ctx.text[pos], err)
	}

	for pos := p.ctx.pos + 1; pos < len(p.ctx.text); pos++ {
		err = fmt.Sprintf("%s%c", err, p.ctx.text[pos])
	}
	glog.V(3).Infof("expr(%s) failed %s", yy.ctx.text, err)
	p.err = errors.New(err)
}

func Parse(text string) (*Expr, error) {
	yy_trigger = &Expr{}
	yy = &yyLex{ctxL: 0}
	yy.ctx = &yy.ctxData[0]
	yy.ctx.pos = 0
	yy.ctx.text = []byte(text)

	glog.V(5).Infof("trigger parse text %s", string(yy.ctx.text))
	yyParse(yy)
	return yy_trigger, yy.err
}
