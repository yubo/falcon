//line parse.y:7
package expr

import __yyfmt__ "fmt"

//line parse.y:7
//line parse.y:11
type yySymType struct {
	yys      int
	num      int
	expr     *Expr
	expr_obj *ExprObj
	text     string
}

const NUM = 57346
const TEXT = 57347
const INDEX = 57348
const VALUE = 57349

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
	"INDEX",
	"VALUE",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parse.y:72

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 3,
	8, 17,
	9, 17,
	10, 17,
	-2, 3,
}

const yyPrivate = 57344

const yyLast = 39

var yyAct = [...]int{

	28, 6, 31, 5, 21, 6, 27, 10, 11, 4,
	10, 11, 3, 14, 13, 15, 7, 8, 20, 30,
	7, 8, 22, 2, 24, 23, 25, 26, 34, 16,
	32, 18, 33, 19, 17, 29, 1, 9, 12,
}
var yyPact = [...]int{

	-3, -1000, -4, -1000, 5, -3, -1000, 28, 25, -3,
	7, -8, 1, 16, -1000, 18, -1, 30, 30, -4,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -16, -1000,
	23, 30, -1000, 21, -1000,
}
var yyPgo = [...]int{

	0, 0, 12, 38, 37, 23, 9, 36,
}
var yyR1 = [...]int{

	0, 7, 7, 5, 5, 5, 5, 1, 2, 3,
	3, 3, 3, 3, 3, 4, 4, 6, 6, 6,
}
var yyR2 = [...]int{

	0, 0, 1, 1, 3, 3, 3, 1, 1, 1,
	2, 1, 2, 1, 2, 2, 2, 1, 6, 4,
}
var yyChk = [...]int{

	-1000, -7, -5, -2, -6, 6, 4, 19, 20, -4,
	11, 12, -3, 9, 8, 10, -5, 6, 6, -5,
	11, 12, -6, -2, 8, 8, 9, 7, -1, 5,
	-1, 18, 7, -1, 7,
}
var yyDef = [...]int{

	1, -2, 2, -2, 0, 0, 8, 0, 0, 0,
	0, 0, 0, 9, 11, 13, 0, 0, 0, 5,
	15, 16, 4, 17, 10, 12, 14, 6, 0, 7,
	0, 0, 19, 0, 18,
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

	2, 3, 4, 5, 19, 20,
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
		//line parse.y:32
		{
			yy_trigger = yyDollar[1].expr
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:36
		{
			yyVAL.expr = &Expr{Type: EXPR_TYPE_RAW, Objs: []interface{}{yyDollar[1].num}}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:37
		{
			yyVAL.expr = &Expr{Type: uint32(yyDollar[2].num), Objs: []interface{}{yyDollar[1].expr_obj, yyDollar[3].expr_obj}}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:38
		{
			yyVAL.expr = &Expr{Type: uint32(yyDollar[2].num), Objs: []interface{}{yyDollar[1].expr, yyDollar[3].expr}}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parse.y:39
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:43
		{
			yyVAL.text = string(yy.t)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:46
		{
			yyVAL.num = yy.i
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:50
		{
			yyVAL.num = EXPR_TYPE_OP_GT
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:51
		{
			yyVAL.num = EXPR_TYPE_OP_GE
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:52
		{
			yyVAL.num = EXPR_TYPE_OP_EQ
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:53
		{
			yyVAL.num = EXPR_TYPE_OP_LE
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:54
		{
			yyVAL.num = EXPR_TYPE_OP_LT
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:55
		{
			yyVAL.num = EXPR_TYPE_OP_NE
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:59
		{
			yyVAL.num = EXPR_TYPE_AND
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parse.y:60
		{
			yyVAL.num = EXPR_TYPE_OR
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parse.y:63
		{
			yyVAL.expr_obj = &ExprObj{Type: EXPR_OBJ_TYPE_RAW, I: yyDollar[1].num}
		}
	case 18:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parse.y:65
		{
			yyVAL.expr_obj = &ExprObj{Type: EXPR_OBJ_TYPE_INDEX, S0: yyDollar[3].text, S1: yyDollar[5].text}
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parse.y:67
		{
			yyVAL.expr_obj = &ExprObj{Type: EXPR_OBJ_TYPE_VALUE, S0: yyDollar[3].text}
		}
	}
	goto yystack /* stack new state and value */
}
