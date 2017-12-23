//line parse/parse.y:7
package parse

import __yyfmt__ "fmt"

//line parse/parse.y:8
import (
	"fmt"
	"os"

	"github.com/yubo/falcon"
)

//line parse/parse.y:20
type yySymType struct {
	yys  int
	num  int
	text string
	b    bool
}

const NUM = 57346
const TEXT = 57347
const IPA = 57348
const ON = 57349
const YES = 57350
const OFF = 57351
const NO = 57352
const INCLUDE = 57353
const ROOT = 57354
const PID_FILE = 57355
const LOG = 57356
const HOST = 57357
const DISABLED = 57358
const DEBUG = 57359
const MODULE = 57360

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"TEXT",
	"IPA",
	"'{'",
	"'}'",
	"';'",
	"ON",
	"YES",
	"OFF",
	"NO",
	"INCLUDE",
	"ROOT",
	"PID_FILE",
	"LOG",
	"HOST",
	"DISABLED",
	"DEBUG",
	"MODULE",
	"'='",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parse/parse.y:124

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 48

var yyAct = [...]int{

	31, 11, 10, 20, 18, 19, 32, 33, 34, 35,
	28, 24, 29, 25, 26, 30, 27, 14, 11, 10,
	43, 22, 3, 42, 40, 39, 5, 7, 8, 4,
	6, 12, 13, 15, 16, 17, 21, 38, 36, 41,
	37, 11, 10, 31, 9, 2, 1, 23,
}
var yyPact = [...]int{

	-1000, 13, -1000, -1000, 36, 10, 36, 36, 36, -4,
	-1000, -1000, 29, 36, -1000, 39, 16, 15, 31, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 14, 11, -1000,
	-1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 47, 26, 21, 46, 45, 44,
}
var yyR1 = [...]int{

	0, 4, 4, 1, 1, 1, 1, 1, 2, 2,
	3, 5, 5, 5, 5, 5, 5, 5, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 1, 1, 1, 0, 1, 1,
	1, 1, 3, 4, 4, 3, 3, 3, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2,
}
var yyChk = [...]int{

	-1000, -4, -5, 9, 16, -2, 17, 14, 15, -6,
	6, 5, -2, 22, 7, -2, -2, -2, 8, 9,
	7, -2, -3, -1, 15, 17, 18, 20, 14, 16,
	19, 4, 10, 11, 12, 13, 9, -2, -3, 9,
	9, 8, 9, 9,
}
var yyDef = [...]int{

	1, -2, 2, 11, 0, 0, 0, 0, 0, 0,
	8, 9, 0, 0, 18, 0, 0, 0, 21, 19,
	20, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	31, 10, 3, 4, 5, 6, 12, 0, 0, 15,
	16, 17, 13, 14,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 9,
	3, 22, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 8,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:44
		{
			yyVAL.b = true
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:45
		{
			yyVAL.b = true
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:46
		{
			yyVAL.b = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:47
		{
			yyVAL.b = false
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parse/parse.y:48
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:52
		{
			yyVAL.text = string(yy.t)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:53
		{
			yyVAL.text = exprText(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse/parse.y:57
		{
			yyVAL.num = yy.i
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse/parse.y:61
		{
			conf.PidFile = yyDollar[2].text
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse/parse.y:62
		{
			if err := os.Setenv(yyDollar[1].text, yyDollar[3].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 14:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse/parse.y:67
		{
			conf.Log = yyDollar[2].text
			conf.Logv = yyDollar[3].num
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse/parse.y:71
		{
			yy.include(yyDollar[2].text)
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse/parse.y:72
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse/parse.y:77
		{
			p1, _ := falcon.PreByte(yy.ctx.text, yy.ctx.pos)
			yy.ctx.text[p1] = ';'

			conf.Conf = append(conf.Conf, yy_module_parse(
				yy.ctx.text[yy_module.pos:yy.ctx.pos],
				yy_module.file, yy_module.lino, yy_module.debug))
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse/parse.y:87
		{
			yy_module = &yyModule{
				level: 1,
				debug: yy.debug,
				file:  yy.ctx.file,
				lino:  yy.ctx.lino,
				pos:   yy.ctx.pos - 1,
			}
			if m, ok := falcon.Modules[yyDollar[1].text]; ok {
				yy_module_parse = m.Parse
			} else {
				yy.Error(fmt.Sprintf("module [%s] not exists", yyDollar[1].text))
			}

		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse/parse.y:102
		{
			if yy_module.level == 0 {
				p1, c1 := falcon.PreByte(yy.ctx.text, yy.ctx.pos)
				p2, c2 := falcon.PreByte(yy.ctx.text, p1-1)
				if c1 == ';' && c2 == '}' {
					yy.ctx.text[p1] = '}'
				}
				yy.ctx.pos = p2 - 1
			}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse/parse.y:112
		{
			yy_module.level++
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse/parse.y:113
		{
			yy_module.level--
		}
	}
	goto yystack /* stack new state and value */
}
