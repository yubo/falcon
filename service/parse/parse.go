//line parse.y:7
package parse

import __yyfmt__ "fmt"

//line parse.y:7
import (
	"fmt"
	"os"

	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/service/config"
)

//line parse.y:19
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
const MIGRATE = 57360
const UPSTREAM = 57361

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
	"MIGRATE",
	"UPSTREAM",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parse.y:114

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 63

var yyAct = [...]int{

	12, 52, 51, 19, 26, 20, 21, 22, 23, 16,
	15, 25, 6, 30, 31, 32, 29, 28, 13, 14,
	34, 40, 10, 8, 11, 9, 16, 15, 50, 42,
	5, 36, 3, 49, 37, 44, 38, 18, 17, 24,
	39, 27, 43, 41, 45, 48, 35, 47, 46, 27,
	16, 15, 16, 15, 33, 20, 21, 22, 23, 7,
	4, 2, 1,
}
var yyPact = [...]int{

	-1000, 23, -1000, -1000, 4, -1000, 29, 28, -5, 32,
	47, 37, 45, 47, 47, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 12, -1000, 27, -5, 14, -1000, -1000,
	-1000, 21, -1000, 45, 47, 24, 19, -7, -8, -1000,
	-1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 3, 0, 4, 62, 61, 60, 59, 54, 46,
	43,
}
var yyR1 = [...]int{

	0, 4, 4, 1, 1, 1, 1, 1, 2, 2,
	3, 5, 5, 6, 6, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 8, 8, 9, 9,
	9, 10, 10, 10, 10, 10,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 1, 1, 1, 0, 1, 1,
	1, 1, 3, 1, 3, 0, 2, 4, 2, 1,
	2, 2, 2, 2, 2, 2, 0, 3, 0, 2,
	4, 0, 4, 4, 4, 4,
}
var yyChk = [...]int{

	-1000, -4, -5, 9, -6, 7, 8, -7, 19, 21,
	18, 20, -2, 14, 15, 6, 5, 9, 9, -1,
	10, 11, 12, 13, 7, -2, -3, 4, -3, -1,
	-2, -2, -2, -8, 8, -9, 19, 22, 9, -1,
	7, -10, 8, -2, 14, -2, -3, -1, -2, 9,
	9, 9, 9,
}
var yyDef = [...]int{

	1, -2, 2, 11, 15, 13, 0, 0, 7, 0,
	0, 19, 7, 0, 0, 8, 9, 12, 14, 16,
	3, 4, 5, 6, 26, 18, 20, 10, 21, 22,
	24, 23, 25, 28, 17, 0, 7, 0, 27, 29,
	31, 0, 30, 7, 0, 0, 0, 0, 0, 32,
	33, 34, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 9,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 8,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22,
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
		//line parse.y:43
		{
			yyVAL.b = true
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:44
		{
			yyVAL.b = true
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:45
		{
			yyVAL.b = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:46
		{
			yyVAL.b = false
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parse.y:47
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:51
		{
			yyVAL.text = string(yy.t)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:52
		{
			yyVAL.text = exprText(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:56
		{
			yyVAL.num = yy.i
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:61
		{
			// end
			conf.Configer.Set(fconfig.APP_CONF_FILE, yy_ss)
			yy_ss = make(map[string]string)

			//conf.Name = fmt.Sprintf("service_%s", conf.Name)
			if conf.Host == "" {
				conf.Host, _ = os.Hostname()
			}
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:74
		{
			// begin
			conf = &config.Service{Name: "service"}
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:81
		{
			conf.Disabled = yyDollar[2].b
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:83
		{
			conf.Host = yyDollar[2].text
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:84
		{
			conf.Debug = 1
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:85
		{
			conf.Debug = yyDollar[2].num
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:86
		{
			yy_ss[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:87
		{
			yy_ss[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:88
		{
			yy.include(yyDollar[2].text)
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:89
		{
			yy_ss[yyDollar[1].text] = yyDollar[2].text
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:90
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:100
		{
			conf.Migrate.Disabled = yyDollar[2].b
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:101
		{
			conf.Migrate.Upstream = yy_ss2
			yy_ss2 = make(map[string]string)
		}
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:108
		{
			yy_ss2[yyDollar[2].text] = yyDollar[3].text
		}
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:109
		{
			yy_ss2[yyDollar[2].text] = fmt.Sprintf("%d", yyDollar[3].num)
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:110
		{
			yy_ss2[yyDollar[2].text] = fmt.Sprintf("%v", yyDollar[3].b)
		}
	case 35:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:111
		{
			yy.include(yyDollar[3].text)
		}
	}
	goto yystack /* stack new state and value */
}
