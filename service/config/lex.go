/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

const (
	eof           = 0
	MAX_CTX_LEVEL = 16
	MODULE_NAME   = "\x1B[37m[SERVICE_PARSE]\x1B[0m "
)

var (
	yy    *yyLex
	conf  *Service
	yy_ss = make(map[string]string)
	yy_as = make([]string, 0)

	keyChars = []byte(`={}:;,()+*/%<>~\[\]?!\|-`)
	keyWords = map[string]int{
		//general
		"on":       ON,
		"yes":      YES,
		"off":      OFF,
		"no":       NO,
		"include":  INCLUDE,
		"root":     ROOT,
		"pidFile":  PID_FILE,
		"log":      LOG,
		"host":     HOST,
		"disabled": DISABLED,
		"debug":    DEBUG,
	}
)

type yyCtx struct {
	text []byte
	pos  int
	lino int
	file string
}

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and falcon.Error(string).
type yyLex struct {
	ctxData [MAX_CTX_LEVEL]yyCtx
	ctxL    int
	ctx     *yyCtx
	t       []byte
	i       int
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

func exprText(s []byte) (ret string) {
	var i int
	var es [][]int

	if es = falcon.F_env.FindAllIndex(s, -1); es == nil {
		return string(s)
	}

	for j := 0; j < len(es); j++ {
		if i < es[j][0] {
			ret += string(s[i:es[j][0]])
		}
		ret += os.Getenv(string(s[es[j][0]+2 : es[j][1]-1]))
		i = es[j][1]
	}
	return ret + string(s[i:])
}

func (p *yyLex) include(filename string) (err error) {
	p.ctxL++
	p.ctx = &p.ctxData[p.ctxL]
	p.ctx.lino = 1
	p.ctx.pos = 0
	p.ctx.file = filename
	if p.ctx.text, err = ioutil.ReadFile(filename); err != nil {
		dir, _ := os.Getwd()
		glog.Errorf("%s %s(curdir:%s)", MODULE_NAME, err.Error(), dir)
		os.Exit(1)
	}
	glog.V(6).Infof("%s ctx level %d", MODULE_NAME, p.ctxL)
	return nil
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (p *yyLex) Lex(yylval *yySymType) int {
	var f []byte
	var b bool

begin:
	text := p.ctx.text[p.ctx.pos:]
	for {
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
				glog.V(6).Infof("%s ctx level %d", MODULE_NAME, p.ctxL)
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

		b = prefix(text, []byte("include"))
		if b {
			p.ctx.pos += len("include")
			return INCLUDE
		}

		f = falcon.F_addr.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			p.t = f
			glog.V(6).Infof("%s return ADDR %s\n", MODULE_NAME, p.t)
			return ADDR
		}

		f = falcon.F_ip.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			p.t = f
			glog.V(6).Infof("%s return IPA %s\n", MODULE_NAME, p.t)
			return IPA
		}

		f = falcon.F_num.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			p.t = f
			i64, _ := strconv.ParseInt(string(f), 0, 0)
			p.i = int(i64)
			glog.V(6).Infof("%s return NUM %d\n", MODULE_NAME, p.i)
			return NUM
		}

		// find keyword
		f = falcon.F_keyword.Find(text)
		if f != nil {
			if val, ok := keyWords[string(f)]; ok {
				p.ctx.pos += len(f)
				glog.V(6).Infof("%s find %s return %d\n", MODULE_NAME, string(f), val)
				return val
			}
		}

		if bytes.IndexByte(keyChars, text[0]) != -1 {
			if !prefix(text, []byte(`//`)) &&
				!prefix(text, []byte(`/*`)) {
				p.ctx.pos++
				glog.V(6).Infof("%s return '%c'\n", MODULE_NAME, int(text[0]))
				return int(text[0])
			}
		}

		// comm
		if text[0] == '#' || prefix(text, []byte(`//`)) {
			for p.ctx.pos < len(p.ctx.text) {
				//glog.Infof("%s %c", MODULE_NAME, p.ctx.text[p.ctx.pos])
				if p.ctx.text[p.ctx.pos] == '\n' {
					p.ctx.pos++
					p.ctx.lino++
					goto begin
				}
				p.ctx.pos++
			}
			return eof
		}

		// ccomm
		if prefix(text, []byte(`/*`)) {
			p.ctx.pos += 2
			for p.ctx.pos < len(p.ctx.text) {
				if p.ctx.text[p.ctx.pos] == '\n' {
					p.ctx.lino++
				}
				if p.ctx.text[p.ctx.pos] == '*' {
					if p.ctx.text[p.ctx.pos-1] == '/' {
						p.Error("Comment nesting not supported")
					}
					if p.ctx.text[p.ctx.pos+1] == '/' {
						p.ctx.pos += 2
						goto begin
					}
				}
				p.ctx.pos++
			}
		}

		// find text
		f = falcon.F_text.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			if f[0] == '"' {
				p.t = f[1 : len(f)-1]
			} else {
				p.t = f[:]
			}
			glog.V(6).Infof("%s return TEXT(%s)", MODULE_NAME, string(p.t))
			return TEXT
		}
		p.Error(fmt.Sprintf("unknown character %c", text[0]))
	}
}

// The parser calls this method on a parse error.
func (p *yyLex) Error(s string) {
	bline := 3
	aline := 3
	p.ctx.pos--
	out := fmt.Sprintf("\x1B[31m%c\x1B[0m", p.ctx.text[p.ctx.pos])

	lino := p.ctx.lino
	for pos := p.ctx.pos - 1; pos > 0; pos-- {
		if p.ctx.text[pos] == '\n' {
			if p.ctx.lino-lino < bline {
				out = fmt.Sprintf("%3d%s", lino, out)
				lino--
			} else {
				out = fmt.Sprintf("%3d%s", lino, out)
				break
			}
		}
		out = fmt.Sprintf("%c%s", p.ctx.text[pos], out)
	}

	lino = p.ctx.lino
	for pos := p.ctx.pos + 1; pos < len(p.ctx.text); pos++ {
		out = fmt.Sprintf("%s%c", out, p.ctx.text[pos])
		if p.ctx.text[pos] == '\n' {
			lino++
			if lino-p.ctx.lino < aline {
				out = fmt.Sprintf("%s%3d", out, lino)
			} else {
				break
			}
		}
	}

	glog.Errorf("%s parse file(%s) error: %s\n%s", MODULE_NAME,
		p.ctx.file, s, out)
	os.Exit(1)
}

func Parse(text []byte, filename string, lino int) falcon.ModuleConf {
	yy = &yyLex{}
	yy.ctx = &yy.ctxData[0]
	yy.ctx.file = filename
	yy.ctx.lino = lino
	yy.ctx.pos = 0
	yy.ctx.text = text

	glog.V(6).Infof("service parse text %s", string(yy.ctx.text))
	yyParse(yy)
	return conf
}
