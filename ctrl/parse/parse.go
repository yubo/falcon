//line ctrl/parse/parse.y:7
package parse

import __yyfmt__ "fmt"

//line ctrl/parse/parse.y:7
import (
	"fmt"
	"os"

	"github.com/yubo/falcon/ctrl/config"
	"github.com/yubo/falcon/utils"
)

//line ctrl/parse/parse.y:19
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
const METRIC = 57360

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
	"METRIC",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line ctrl/parse/parse.y:104

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 49

var yyAct = [...]int{

	11, 20, 21, 22, 23, 38, 18, 25, 16, 15,
	24, 6, 29, 30, 17, 32, 19, 12, 14, 27,
	31, 9, 8, 10, 13, 7, 16, 15, 28, 36,
	26, 26, 16, 15, 34, 35, 37, 20, 21, 22,
	23, 5, 4, 3, 16, 15, 33, 2, 1,
}
var yyPact = [...]int{

	-1000, 34, -1000, -1000, 3, -1000, 5, -3, -9, 39,
	26, 27, 39, 13, 39, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 21, -1000, 39, -1000, -4, -1000,
}
var yyPgo = [...]int{

	0, 16, 0, 7, 48, 47, 46, 42, 25,
}
var yyR1 = [...]int{

	0, 4, 4, 1, 1, 1, 1, 1, 2, 2,
	3, 6, 6, 6, 5, 5, 7, 7, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 1, 1, 1, 0, 1, 1,
	1, 0, 2, 4, 1, 3, 1, 3, 0, 2,
	2, 1, 2, 2, 2, 2, 2, 4, 2,
}
var yyChk = [...]int{

	-1000, -4, -5, 9, -7, 7, 8, -8, 19, 18,
	20, -2, 14, 21, 15, 6, 5, 9, 9, -1,
	10, 11, 12, 13, -2, -3, 4, -3, -1, -2,
	-2, 7, -2, -6, -2, 14, 8, -2, 9,
}
var yyDef = [...]int{

	1, -2, 2, 14, 18, 16, 0, 0, 7, 0,
	21, 7, 0, 0, 0, 8, 9, 15, 17, 19,
	3, 4, 5, 6, 20, 22, 10, 23, 24, 26,
	25, 11, 28, 0, 12, 0, 27, 0, 13,
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
		//line ctrl/parse/parse.y:43
		{
			yyVAL.b = true
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:44
		{
			yyVAL.b = true
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:45
		{
			yyVAL.b = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:46
		{
			yyVAL.b = false
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ctrl/parse/parse.y:47
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:51
		{
			yyVAL.text = string(yy.t)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:52
		{
			yyVAL.text = exprText(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:56
		{
			yyVAL.num = yy.i
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:60
		{
			yy_as = append(yy_as, yyDollar[2].text)
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/parse/parse.y:61
		{
			yy.include(yyDollar[3].text)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ctrl/parse/parse.y:65
		{
			// end
			conf.Ctrl.Set(utils.APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)

			conf.Name = fmt.Sprintf("ctrl_%s", conf.Name)
			if conf.Host == "" {
				conf.Host, _ = os.Hostname()
			}
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:78
		{
			// begin
			conf = &config.ConfCtrl{Name: "ctrl"}
			conf.Ctrl.Set(utils.APP_CONF_DEFAULT, config.ConfDefault)
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:86
		{
			conf.Disabled = yyDollar[2].b
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:87
		{
			conf.Host = yyDollar[2].text
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ctrl/parse/parse.y:88
		{
			conf.Debug = 1
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:89
		{
			conf.Debug = yyDollar[2].num
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:90
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:91
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:92
		{
			yy.include(yyDollar[2].text)
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:93
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ctrl/parse/parse.y:94
		{
			conf.Metrics = yy_as
			yy_as = make([]string, 0)
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ctrl/parse/parse.y:97
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	}
	goto yystack /* stack new state and value */
}
