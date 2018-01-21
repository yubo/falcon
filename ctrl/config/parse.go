//line ctrl/config/parse.y:7
package config

import __yyfmt__ "fmt"

//line ctrl/config/parse.y:7
import (
	"fmt"
	"os"

	"github.com/yubo/falcon"
)

//line ctrl/config/parse.y:18
type yySymType struct {
	yys  int
	num  int
	text string
	b    bool
}

const NUM = 57346
const TEXT = 57347
const IPA = 57348
const ADDR = 57349
const ON = 57350
const YES = 57351
const OFF = 57352
const NO = 57353
const INCLUDE = 57354
const ROOT = 57355
const PID_FILE = 57356
const LOG = 57357
const HOST = 57358
const DISABLED = 57359
const DEBUG = 57360
const METRIC = 57361
const AGENT = 57362
const TRANSFER = 57363
const BACKEND = 57364

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"TEXT",
	"IPA",
	"ADDR",
	"'{'",
	"'}'",
	"';'",
	"'*'",
	"'('",
	"')'",
	"'+'",
	"'-'",
	"'<'",
	"'>'",
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
	"METRIC",
	"AGENT",
	"TRANSFER",
	"BACKEND",
	"'='",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line ctrl/config/parse.y:144

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 124

var yyAct = [...]int{

	45, 72, 32, 71, 26, 11, 80, 20, 18, 19,
	31, 6, 37, 39, 35, 41, 36, 68, 20, 18,
	19, 67, 49, 53, 12, 14, 61, 25, 9, 8,
	10, 13, 21, 22, 23, 46, 44, 58, 43, 59,
	27, 28, 29, 30, 24, 62, 65, 66, 63, 52,
	64, 33, 20, 18, 19, 48, 50, 69, 70, 34,
	5, 74, 3, 51, 40, 27, 28, 29, 30, 33,
	20, 18, 19, 33, 77, 78, 79, 34, 17, 16,
	38, 34, 15, 27, 28, 29, 30, 20, 18, 19,
	7, 47, 20, 18, 19, 4, 42, 20, 18, 19,
	60, 76, 2, 1, 46, 44, 20, 18, 19, 46,
	44, 0, 0, 54, 75, 73, 55, 54, 56, 57,
	55, 0, 56, 57,
}
var yyPact = [...]int{

	-1000, 52, -1000, -1000, 2, -1000, 34, 17, 22, 101,
	69, 47, 101, 56, 101, 87, 82, 13, -1000, -1000,
	-1000, 55, 41, 15, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 106, -1000, 69, 106, -1000, -1000, 101, -1000,
	-1000, -1000, -1000, 16, 101, 65, 101, -1000, 11, -1000,
	7, -1000, -1000, -1000, 69, 69, -13, -16, 102, -1000,
	92, -1000, -1000, 106, -1000, -1000, -1000, -1000, -1000, 106,
	106, 69, 69, -1000, -1000, 101, -1000, 106, 106, -4,
	-1000,
}
var yyPgo = [...]int{

	0, 4, 0, 2, 103, 102, 100, 95, 90, 82,
	79, 78, 38,
}
var yyR1 = [...]int{

	0, 4, 4, 1, 1, 1, 1, 1, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 6, 6, 6,
	5, 5, 7, 7, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 9,
	9, 10, 10, 11, 11, 12, 12, 12, 12, 12,
	12,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 1, 1, 1, 0, 1, 1,
	1, 1, 3, 3, 3, 4, 4, 0, 2, 4,
	1, 3, 1, 3, 0, 2, 2, 1, 2, 2,
	2, 2, 2, 3, 4, 2, 2, 2, 2, 2,
	3, 2, 3, 2, 3, 0, 2, 2, 2, 2,
	2,
}
var yyChk = [...]int{

	-1000, -4, -5, 10, -7, 8, 9, -8, 27, 26,
	28, -2, 22, 29, 23, -9, -10, -11, 6, 7,
	5, 30, 31, 32, 10, 10, -1, 18, 19, 20,
	21, -2, -3, 4, 12, -3, -1, -2, 33, -2,
	8, -2, 9, -12, 23, -2, 22, 9, -12, 9,
	-12, 8, 8, 8, 11, 14, 16, 17, -3, -2,
	-6, 10, -2, -3, -1, -2, -2, 10, 10, -3,
	-3, 16, 17, 13, -2, 22, 9, -3, -3, -2,
	10,
}
var yyDef = [...]int{

	1, -2, 2, 20, 24, 22, 0, 0, 7, 0,
	27, 7, 0, 0, 0, 45, 45, 45, 8, 9,
	10, 0, 0, 0, 21, 23, 25, 3, 4, 5,
	6, 26, 28, 11, 0, 29, 30, 32, 0, 31,
	17, 35, 36, 0, 0, 7, 0, 37, 0, 38,
	0, 39, 41, 43, 0, 0, 0, 0, 0, 33,
	0, 40, 46, 47, 48, 49, 50, 42, 44, 13,
	14, 0, 0, 12, 18, 0, 34, 15, 16, 0,
	19,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	12, 13, 11, 14, 3, 15, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 10,
	16, 33, 17, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 8, 3, 9,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32,
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
		//line ctrl/config/parse.y:42
		{
			yyVAL.b = true
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:43
		{
			yyVAL.b = true
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:44
		{
			yyVAL.b = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:45
		{
			yyVAL.b = false
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ctrl/config/parse.y:46
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:50
		{
			yyVAL.text = string(yy.t)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:51
		{
			yyVAL.text = string(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:52
		{
			yyVAL.text = exprText(yy.t)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:56
		{
			yyVAL.num = yy.i
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/config/parse.y:57
		{
			yyVAL.num = yyDollar[2].num
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/config/parse.y:58
		{
			yyVAL.num = yyDollar[1].num * yyDollar[3].num
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/config/parse.y:59
		{
			yyVAL.num = yyDollar[1].num + yyDollar[3].num
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/config/parse.y:60
		{
			yyVAL.num = int(uint(yyDollar[1].num) << uint(yyDollar[4].num))
		}
	case 16:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/config/parse.y:61
		{
			yyVAL.num = int(uint(yyDollar[1].num) >> uint(yyDollar[4].num))
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:65
		{
			yy_as = append(yy_as, yyDollar[2].text)
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/config/parse.y:66
		{
			yy.include(yyDollar[3].text)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/config/parse.y:70
		{
			// end
			conf.Ctrl.Set(falcon.APP_CONF_FILE, yy_ss)
			yy_ss = make(map[string]string)

			//conf.Name = fmt.Sprintf("ctrl_%s", conf.Name)
			if conf.Host == "" {
				conf.Host, _ = os.Hostname()
			}
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:83
		{
			// begin
			conf = &Ctrl{Name: "ctrl"}
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:90
		{
			conf.Disabled = yyDollar[2].b
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:91
		{
			conf.Host = yyDollar[2].text
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/config/parse.y:92
		{
			conf.Debug = 1
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:93
		{
			conf.Debug = yyDollar[2].num
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:94
		{
			yy_ss[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:95
		{
			yy_ss[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:96
		{
			yy.include(yyDollar[2].text)
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:97
		{
			yy_ss[yyDollar[1].text] = yyDollar[2].text
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/config/parse.y:98
		{
			if err := os.Setenv(yyDollar[1].text, yyDollar[3].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/config/parse.y:103
		{
			conf.Metrics = yy_as
			yy_as = make([]string, 0)
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:106
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:110
		{
			conf.Agent.Set(falcon.APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:113
		{
			conf.Transfer.Set(falcon.APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:116
		{
			conf.Backend.Set(falcon.APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:133
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:138
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:139
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:140
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/config/parse.y:141
		{
			yy.include(yyDollar[2].text)
		}
	}
	goto yystack /* stack new state and value */
}
