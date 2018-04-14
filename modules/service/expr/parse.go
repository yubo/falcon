//line parse.y:7
package expr

import __yyfmt__ "fmt"

//line parse.y:7
//line parse.y:11
type yySymType struct {
	yys      int
	num      int
	float    float64
	b        bool
	expr     *Expr
	expr_obj *ExprObj
	text     string
}

const NUM = 57346
const TEXT = 57347
const TRUE = 57348
const FALSE = 57349
const COUNT = 57350

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"TEXT",
	"'('",
	"')'",
	"'='",
	"'>'",
	"'<'",
	"'&'",
	"'|'",
	"'#'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"','",
	"TRUE",
	"FALSE",
	"COUNT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parse.y:135

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 63

var yyAct = [...]int{

	8, 37, 10, 49, 34, 16, 11, 12, 46, 45,
	11, 12, 5, 40, 4, 2, 14, 15, 18, 17,
	19, 20, 35, 9, 35, 6, 7, 9, 39, 25,
	11, 28, 27, 38, 26, 30, 31, 42, 41, 36,
	11, 44, 29, 33, 24, 22, 47, 12, 11, 36,
	51, 50, 48, 32, 41, 43, 23, 14, 15, 21,
	1, 13, 3,
}
var yyPact = [...]int{

	6, -1000, 5, -1000, 10, 6, -1000, -1000, -1000, 39,
	38, -1000, -1000, 6, 23, 20, 2, 34, -1000, 27,
	46, 36, -1000, 26, -1000, 5, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -5, 42, 44, -1000, -1000, -1000,
	44, -1000, -1000, -9, -10, 44, 10, -1000, -15, 44,
	-1000, 42,
}
var yyPgo = [...]int{

	0, 62, 2, 5, 61, 1, 0, 15, 14, 60,
	59, 56, 4, 55,
}
var yyR1 = [...]int{

	0, 9, 9, 7, 7, 7, 7, 1, 1, 2,
	6, 3, 3, 3, 3, 3, 3, 4, 4, 8,
	8, 8, 10, 10, 10, 10, 10, 11, 11, 12,
	12, 12, 5, 5, 13, 13,
}
var yyR2 = [...]int{

	0, 0, 1, 1, 3, 3, 3, 1, 1, 1,
	1, 1, 2, 1, 2, 1, 2, 2, 2, 1,
	3, 3, 1, 8, 6, 4, 2, 1, 3, 1,
	2, 1, 1, 2, 0, 3,
}
var yyChk = [...]int{

	-1000, -9, -7, -1, -8, 6, 19, 20, -6, 21,
	-2, 4, 5, -4, 11, 12, -3, 9, 8, 10,
	-7, -10, 6, -11, 6, -7, 11, 12, -8, 8,
	8, 9, 7, 7, -12, -6, 13, -5, 7, -12,
	18, -2, -6, -13, -6, 18, 18, -6, -3, 18,
	-5, -6,
}
var yyDef = [...]int{

	1, -2, 2, 3, 0, 0, 7, 8, 19, 0,
	0, 10, 9, 0, 0, 0, 0, 11, 13, 15,
	0, 0, 22, 0, 27, 5, 17, 18, 4, 12,
	14, 16, 6, 20, 26, 29, 0, 31, 21, 34,
	0, 33, 30, 28, 25, 0, 0, 35, 24, 0,
	23, 32,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 13, 3, 3, 11, 3,
	6, 7, 16, 14, 18, 15, 3, 17, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	10, 8, 9, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 12,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 19, 20, 21,
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

	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:37
		{
			yy_trigger = yyDollar[1].expr
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:41
		{
			yyVAL.expr = &Expr{Type: EXPR_TYPE_RAW << 1, Objs: []interface{}{yyDollar[1].b}}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:42
		{
			yyVAL.expr = &Expr{Type: uint32(yyDollar[2].num), Objs: []interface{}{yyDollar[1].expr_obj, yyDollar[3].expr_obj}}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:43
		{
			yyVAL.expr = &Expr{Type: uint32(yyDollar[2].num), Objs: []interface{}{yyDollar[1].expr, yyDollar[3].expr}}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:44
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:48
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:49
		{
			yyVAL.b = false
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:54
		{
			yyVAL.text = string(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:57
		{
			yyVAL.float = yy.f
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:61
		{
			yyVAL.num = EXPR_TYPE_OP_GT
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:62
		{
			yyVAL.num = EXPR_TYPE_OP_GE
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:63
		{
			yyVAL.num = EXPR_TYPE_OP_EQ
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:64
		{
			yyVAL.num = EXPR_TYPE_OP_LE
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:65
		{
			yyVAL.num = EXPR_TYPE_OP_LT
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:66
		{
			yyVAL.num = EXPR_TYPE_OP_NE
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:70
		{
			yyVAL.num = EXPR_TYPE_AND
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:71
		{
			yyVAL.num = EXPR_TYPE_OR
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:74
		{
			yyVAL.expr_obj = &ExprObj{Type: EXPR_OBJ_TYPE_RAW, Args: []float64{yyDollar[1].float}}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:76
		{
			yy.err = yy_obj.reduce("count")
			yyVAL.expr_obj = yy_obj
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:79
		{
			yy.err = yy_obj.reduce(yyDollar[1].text)
			yyVAL.expr_obj = yy_obj
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:86
		{
			yy_obj = &ExprObj{}
		}
	case 23:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parse.y:88
		{
			yy_obj.Args = append(yy_obj.Args, yyDollar[4].float, float64(yyDollar[6].num), float64(yyDollar[8].num))
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parse.y:90
		{
			yy_obj.Args = append(yy_obj.Args, yyDollar[4].float, float64(yyDollar[6].num))
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:92
		{
			yy_obj.Args = append(yy_obj.Args, yyDollar[4].float)
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:97
		{
			yy_obj = &ExprObj{}
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:103
		{
			yy_obj.Args = []float64{yyDollar[1].float}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:105
		{
			yy_obj.Type = 1
			yy_obj.Args = []float64{yyDollar[2].float}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:108
		{
			yy_obj.Args = []float64{float64(yyDollar[1].num)}
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:114
		{
			yyVAL.num = int(yyDollar[1].float)
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:116
		{
			switch yyDollar[2].text {
			case "m":
				yyVAL.num = int(yyDollar[1].float) * 60
			case "h":
				yyVAL.num = int(yyDollar[1].float) * 60 * 60
			case "d":
				yyVAL.num = int(yyDollar[1].float) * 60 * 60 * 24
			default:
				yy.err = EINVAL
			}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:130
		{
			yy_obj.Args = append(yy_obj.Args, yyDollar[3].float)
		}
	}
	goto yystack /* stack new state and value */
}
