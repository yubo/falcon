/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/golang/glog"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/backend"
	"github.com/yubo/falcon/lb"
	"github.com/yubo/falcon/specs"
)

const (
	MAX_CTX_LEVEL = 16
)

type yyCtx struct {
	text []byte
	pos  int
	lino int
	file string
}

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type yyLex struct {
	ctxData [MAX_CTX_LEVEL]yyCtx
	ctxL    int
	ctx     *yyCtx
	t       []byte
	i       int
}

type Falcon struct {
	PidFile string
	Agent   agent.Agent
	Lb []lb.Lb
	Backend []backend.Backend
	Modules []specs.Module
}

func (p Falcon) String() string {
	ret := fmt.Sprintf("%-17s %s", "pidfile", p.PidFile)
	for _, v := range p.Modules {
		ret += fmt.Sprintf("\n%s (\n%s\n)",
			v.Desc(), specs.IndentLines(1, v.String()))
	}
	return ret
}

var (
	conf               *Falcon
	yy                 *yyLex
	yy_lb         *lb.Lb
	yy_lb_backend *lb.Backend
	yy_backend         *backend.Backend
	yy_mod_name        *string
	yy_mod_disable     *bool
	yy_mod_debug       *int
	yy_mod_http        *bool
	yy_mod_http_addr   *string
	yy_mod_rpc         *bool
	yy_mod_rpc_addr    *string
	yy_ss              map[string]string
	yy_as              []string
)

const (
	eof = 0
)

var (
	f_ip      = regexp.MustCompile(`^[0-9]+\.[0-0]+\.[0-9]+\.[0-9]+`)
	f_num     = regexp.MustCompile(`^0x[0-9a-fA-F]+|^[0-9]+`)
	f_keyword = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]*`)
	f_word    = regexp.MustCompile(`^["]?[^"\n \t;]+["]?`)
	f_env     = regexp.MustCompile(`\$\{[a-zA-Z][0-9a-zA-Z_]+\}`)

	keywords = map[string]int{
		"pid_file":         PID_FILE,
		"agent":            AGENT,
		"debug":            DEBUG,
		"host":             HOST,
		"http":             HTTP,
		"http_addr":        HTTP_ADDR,
		"interval":         INTERVAL,
		"iface_prefix":     IFACE_PREFIX,
		"disabled":         DISABLED,
		"batch":            BATCH,
		"conn_timeout":     CONN_TIMEOUT,
		"call_timeout":     CALL_TIMEOUT,
		"upstreams":        UPSTREAMS,
		"lb":          HANDOFF,
		"rpc":              RPC,
		"rpc_addr":         RPC_ADDR,
		"replicas":         REPLICAS,
		"concurrency":      CONCURRENCY,
		"backends":         BACKENDS,
		"tsdb":             TSDB,
		"sched":            SCHED,
		"consistent":       CONSISTENT,
		"falcon":           FALCON,
		"backend":          BACKEND,
		"rrd_storage":      RRD_STORAGE,
		"dsn":              DSN,
		"db_max_idle":      DB_MAX_IDLE,
		"shm_magic_code":   SHM_MAGIC_CODE,
		"shm_key_start_id": SHM_KEY_START_ID,
		"shm_segment_size": SHM_SEGMENT_SIZE,
		"migrate":          MIGRATE,
		"storage":          STORAGE,
		"rrd":              RRD,
		"hdisks":           HDISKS,
		"on":               ON,
		"yes":              YES,
		"off":              OFF,
		"no":               NO,
		"include":          INCLUDE,
		"root":             ROOT,
	}
)

func init() {
	yy_ss = make(map[string]string)
	yy_as = make([]string, 0)
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

	if es = f_env.FindAllIndex(s, -1); es == nil {
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
		glog.Errorf("%s(curdir:%s)", err.Error(), dir)
		os.Exit(1)
	}
	glog.V(4).Infof("ctx level %d", p.ctxL)
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
			glog.V(4).Infof("ctx level %d", p.ctxL)
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
				glog.V(4).Infof("ctx level %d", p.ctxL)
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

		f = f_ip.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			p.t = f[:]
			return IPA
		}

		f = f_num.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			p.t = f[:]
			i64, _ := strconv.ParseInt(string(f), 0, 0)
			p.i = int(i64)
			glog.V(4).Infof("return NUM %d\n", p.i)
			return NUM
		}

		// find keyword
		f = f_keyword.Find(text)
		if f != nil {
			if val, ok := keywords[string(f)]; ok {
				p.ctx.pos += len(f)
				glog.V(4).Infof("find %s return %d\n", string(f), val)
				return val
			}
		}

		if bytes.IndexByte([]byte(`={}:;,()+*/%<>~\[\]?!\|-`), text[0]) != -1 {
			if !prefix(text, []byte(`//`)) &&
				!prefix(text, []byte(`/*`)) {
				p.ctx.pos++
				glog.V(4).Infof("return '%c'\n", int(text[0]))
				return int(text[0])
			}
		}

		// comm
		if text[0] == '#' || prefix(text, []byte(`//`)) {
			for p.ctx.pos < len(p.ctx.text) {
				//glog.Infof("%c", p.ctx.text[p.ctx.pos])
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
		f = f_word.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			if f[0] == '"' {
				p.t = f[1 : len(f)-1]
			} else {
				p.t = f[:]
			}
			glog.V(4).Infof("return TEXT(%s)", string(p.t))
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

	glog.Errorf("parse file(%s) error: %s\n%s",
		p.ctx.file, s, out)
	os.Exit(1)
}

func Parse(filename string) *Falcon {
	var err error
	conf = &Falcon{}
	yy = &yyLex{}
	yy.ctxL = 0
	yy.ctx = &yy.ctxData[0]
	yy.ctx.file = filename
	yy.ctx.lino = 1
	yy.ctx.pos = 0
	if yy.ctx.text, err = ioutil.ReadFile(filename); err != nil {
		yy.Error(err.Error())
	}
	yyParse(yy)
	return conf
}
