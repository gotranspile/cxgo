package cxgo

import (
	"go/ast"
	"go/token"

	"github.com/gotranspile/cxgo/types"
)

type Expr interface {
	Node
	// CType returns a type of this expression, given an expected type from the parent.
	// If any type is acceptable, pass nil.
	CType(exp types.Type) types.Type
	AsExpr() GoExpr
	IsConst() bool
	HasSideEffects() bool
	Uses() []types.Usage
}

func canAssignTo(x Expr) bool {
	x = cUnwrap(x)
	switch x.(type) {
	case Ident, *Deref, *CSelectExpr, *CIndexExpr:
		return true
	}
	return false
}

type Ident interface {
	Expr
	Identifier() *types.Ident
}

func cIsTrue(x Expr) bool {
	v, ok := cIsBoolConst(x)
	if ok {
		return v
	}
	return false
}

type CAsmExpr struct {
	e   *types.Env
	typ types.Type
}

func (e *CAsmExpr) Visit(_ Visitor) {}

func (e *CAsmExpr) CType(types.Type) types.Type {
	if e.typ == nil {
		e.typ = e.e.FuncT(
			nil,
			&types.Field{
				Name: types.NewUnnamed(e.e.C().String()),
			},
		)
	}
	return e.typ
}

func (e *CAsmExpr) IsConst() bool {
	return false
}

func (e *CAsmExpr) HasSideEffects() bool {
	return true
}

func (e *CAsmExpr) AsExpr() GoExpr {
	return ident("asm")
}

func (e *CAsmExpr) Uses() []types.Usage {
	return nil
}

var _ Expr = &CMultiExpr{}

func (g *translator) NewCMultiExpr(exprs ...Expr) Expr {
	if len(exprs) == 1 {
		return exprs[0]
	}
	return &CMultiExpr{
		g:     g,
		Exprs: exprs,
	}
}

// CMultiExpr is a list of C expressions executed one by one, and returning the result of the last one.
type CMultiExpr struct {
	g     *translator
	Exprs []Expr
}

func (e *CMultiExpr) Visit(v Visitor) {
	for _, x := range e.Exprs {
		v(x)
	}
}

func (e *CMultiExpr) CType(exp types.Type) types.Type {
	tp := e.Exprs[len(e.Exprs)-1].CType(exp)
	if k := tp.Kind(); k.IsUntypedInt() {
		return e.g.env.DefIntT()
	}
	return tp
}

func (e *CMultiExpr) IsConst() bool {
	return false
}

func (e *CMultiExpr) HasSideEffects() bool {
	return true
}

func (e *CMultiExpr) AsExpr() GoExpr {
	if len(e.Exprs) == 1 {
		return e.Exprs[0].AsExpr()
	}
	typ := e.CType(nil)
	var stmts []CStmt
	for i, x := range e.Exprs {
		if i == len(e.Exprs)-1 {
			stmts = append(stmts, e.g.NewReturnStmt(x, typ)...)
		} else {
			stmts = append(stmts, NewCExprStmt(x)...)
		}
	}
	return callLambda(typ.GoType(), e.g.NewCBlock(stmts...).AsStmt()...)
}

func (e *CMultiExpr) Uses() []types.Usage {
	var list []types.Usage
	for _, s := range e.Exprs {
		list = append(list, s.Uses()...)
	}
	return list
}

var (
	_ Expr  = IdentExpr{}
	_ Ident = IdentExpr{}
)

type IdentExpr struct {
	*types.Ident
}

func (e IdentExpr) Visit(v Visitor) {}

func (e IdentExpr) Identifier() *types.Ident {
	return e.Ident
}

func (e IdentExpr) IsConst() bool {
	return false // TODO: enums
}

func (e IdentExpr) HasSideEffects() bool {
	return false
}

func (e IdentExpr) AsExpr() GoExpr {
	return e.GoIdent()
}

func (e IdentExpr) Uses() []types.Usage {
	return []types.Usage{{Ident: e.Ident, Access: types.AccessUnknown}}
}

type BinaryOp string

func (op BinaryOp) Precedence() int {
	switch op {
	case BinOpMult, BinOpDiv, BinOpMod:
		return 4
	case BinOpAdd, BinOpSub:
		return 3
	case BinOpLsh, BinOpRsh:
		return 2
	default:
		return 1
	}
}

func (op BinaryOp) IsCommutative() bool {
	switch op {
	case BinOpAdd, BinOpMult,
		BinOpBitAnd, BinOpBitOr, BinOpBitXor:
		return true
	}
	return false
}
func (op BinaryOp) IsArithm() bool {
	switch op {
	case BinOpAdd, BinOpSub, BinOpMod, BinOpMult, BinOpDiv:
		return true
	}
	return false
}
func (op BinaryOp) GoToken() token.Token {
	var tok token.Token
	switch op {
	case BinOpMult:
		tok = token.MUL
	case BinOpDiv:
		tok = token.QUO
	case BinOpMod:
		tok = token.REM
	case BinOpAdd:
		tok = token.ADD
	case BinOpSub:
		tok = token.SUB
	case BinOpLsh:
		tok = token.SHL
	case BinOpRsh:
		tok = token.SHR
	case BinOpBitAnd:
		tok = token.AND
	case BinOpBitOr:
		tok = token.OR
	case BinOpBitXor:
		tok = token.XOR
	default:
		panic(op)
	}
	return tok
}
func (op BinaryOp) GoAssignToken() token.Token {
	var tok token.Token
	switch op {
	case "":
		tok = token.ASSIGN
	case BinOpMult:
		tok = token.MUL_ASSIGN
	case BinOpDiv:
		tok = token.QUO_ASSIGN
	case BinOpMod:
		tok = token.REM_ASSIGN
	case BinOpAdd:
		tok = token.ADD_ASSIGN
	case BinOpSub:
		tok = token.SUB_ASSIGN
	case BinOpLsh:
		tok = token.SHL_ASSIGN
	case BinOpRsh:
		tok = token.SHR_ASSIGN
	case BinOpBitAnd:
		tok = token.AND_ASSIGN
	case BinOpBitXor:
		tok = token.XOR_ASSIGN
	case BinOpBitOr:
		tok = token.OR_ASSIGN
	default:
		panic(op)
	}
	return tok
}

const (
	BinOpMult BinaryOp = "*"
	BinOpDiv  BinaryOp = "/"
	BinOpMod  BinaryOp = "%"

	BinOpAdd BinaryOp = "+"
	BinOpSub BinaryOp = "-"

	BinOpLsh BinaryOp = "<<"
	BinOpRsh BinaryOp = ">>"

	BinOpBitAnd BinaryOp = "&"
	BinOpBitOr  BinaryOp = "|"
	BinOpBitXor BinaryOp = "^"
)

func (g *translator) NewCBinaryExpr(x Expr, op BinaryOp, y Expr) Expr {
	return g.newCBinaryExpr(nil, x, op, y)
}
func (g *translator) newCBinaryExpr(exp types.Type, x Expr, op BinaryOp, y Expr) Expr {
	if op.IsCommutative() {
		if x.IsConst() && !y.IsConst() {
			// let const to be the second operand
			x, y = y, x
		}
		if x.CType(nil).Kind().IsInt() && y.CType(nil).Kind().IsPtr() {
			x, y = y, x
		}
	}
	if xt := x.CType(nil); xt.Kind().Is(types.Array) {
		at := xt.(types.ArrayType)
		if op == BinOpAdd {
			// index array/slices instead of adding addresses
			return g.cAddr(g.NewCIndexExpr(x, y, nil))
		}
		return g.newCBinaryExpr(exp, g.cCast(g.env.PtrT(at.Elem()), x), op, y)
	}
	if xt := x.CType(nil); xt.Kind().Is(types.Func) && !y.IsConst() {
		return g.newCBinaryExpr(exp, g.cCast(g.env.PtrT(nil), x), op, y)
	}
	if xt, ok := types.Unwrap(x.CType(nil)).(types.PtrType); ok && op.IsArithm() {
		if addr, ok := cUnwrap(x).(*TakeAddr); ok && (op == BinOpAdd || op == BinOpSub) {
			if ind, ok := addr.X.(*CIndexExpr); ok {
				// adding another component to an existing index expression
				return g.cAddr(g.NewCIndexExpr(ind.Expr,
					g.NewCBinaryExpr(ind.Index, op, y),
					nil))
			}
		}
		x := g.ToPointer(x)
		if yk := y.CType(nil).Kind(); yk.IsInt() {
			if op == BinOpAdd {
				return cPtrOffset(x, y)
			} else if op == BinOpSub {
				return cPtrOffset(x, g.NewCUnaryExpr(UnaryMinus, y))
			}
		} else if yk.IsPtr() {
			if op == BinOpSub {
				return cPtrDiff(x, g.ToPointer(y))
			}
		}
		inc := xt.ElemSizeof()
		if inc != 1 {
			y = g.newCBinaryExpr(exp, y, BinOpMult, cIntLit(int64(inc), 10))
		}
		return &CBinaryExpr{
			Left:  x,
			Op:    op,
			Right: y,
		}
	}
	xt := x.CType(exp)
	yt := y.CType(exp)
	typ := g.env.CommonType(xt, yt)
	x = g.cCast(typ, x)
	y = g.cCast(typ, y)
	return g.cCast(typ, &CBinaryExpr{
		Left:  cParenLazyOp(x, op),
		Op:    op,
		Right: cParenLazyOpR(y, op),
	})
}

func (g *translator) NewCBinaryExprT(x Expr, op BinaryOp, y Expr, _ types.Type) Expr {
	return g.NewCBinaryExpr(x, op, y)
}

type CBinaryExpr struct {
	Left  Expr
	Op    BinaryOp
	Right Expr
}

func (e *CBinaryExpr) Visit(v Visitor) {
	v(e.Left)
	v(e.Right)
}

func (e *CBinaryExpr) CType(exp types.Type) types.Type {
	l, r := e.Left.CType(nil), e.Right.CType(nil)
	if l.Kind().IsUntyped() {
		return r
	}
	return l
}

func (e *CBinaryExpr) IsConst() bool {
	return e.Left.IsConst() && e.Right.IsConst()
}

func (e *CBinaryExpr) HasSideEffects() bool {
	return e.Left.HasSideEffects() || e.Right.HasSideEffects()
}

func (e *CBinaryExpr) AsExpr() GoExpr {
	return &ast.BinaryExpr{
		X:  e.Left.AsExpr(),
		Op: e.Op.GoToken(),
		Y:  e.Right.AsExpr(),
	}
}

func (e *CBinaryExpr) Uses() []types.Usage {
	return types.UseRead(e.Left, e.Right)
}

var (
	_ Expr      = &CTernaryExpr{}
	_ CStmtConv = &CTernaryExpr{}
)

func (g *translator) NewCTernaryExpr(cond BoolExpr, then, els Expr) Expr {
	return &CTernaryExpr{
		g:    g,
		Cond: cond,
		Then: then,
		Else: els,
	}
}

type CTernaryExpr struct {
	g    *translator
	Cond BoolExpr
	Then Expr
	Else Expr
}

func (e *CTernaryExpr) Visit(v Visitor) {
	v(e.Cond)
	v(e.Then)
	v(e.Else)
}

func (e *CTernaryExpr) CType(types.Type) types.Type {
	tt := e.Then.CType(nil)
	et := e.Else.CType(nil)
	tk := tt.Kind()
	ek := et.Kind()
	if tk.IsUntypedInt() && ek.IsUntypedInt() {
		return e.g.env.CommonType(tt, et)
	}
	if tk.IsUntyped() {
		return et
	}
	if tk.Is(types.Array) && et.Kind().IsPtr() {
		return et
	}
	if _, ok := types.Unwrap(tt).(types.IntType); ok {
		if _, ok := types.Unwrap(et).(types.IntType); ok {
			return e.g.env.CommonType(tt, et)
		}
	}
	return tt
}

func (e *CTernaryExpr) HasSideEffects() bool {
	return e.Cond.HasSideEffects() ||
		e.Then.HasSideEffects() ||
		e.Else.HasSideEffects()
}

func (e *CTernaryExpr) IsConst() bool {
	return false
}

func (e *CTernaryExpr) AsExpr() GoExpr {
	ret := e.CType(nil)
	var stmts []GoStmt
	stmts = append(stmts, ifelse(
		e.Cond.AsExpr(),
		asStmts(e.g.NewReturnStmt(e.Then, ret)),
		nil,
	))
	stmts = append(stmts, asStmts(e.g.NewReturnStmt(e.Else, ret))...)
	return callLambda(
		ret.GoType(),
		stmts...,
	)
}

func (e *CTernaryExpr) ToStmt() []CStmt {
	return []CStmt{e.g.NewCIfStmt(
		e.Cond,
		NewCExprStmt(e.Then),
		e.g.NewCBlock(NewCExprStmt(e.Else)...),
	)}
}

func (e *CTernaryExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(e.Cond)...)
	list = append(list, types.UseWrite(e.Then, e.Else)...)
	return list
}

func cUnwrap(e Expr) Expr {
	switch e := e.(type) {
	case *CParentExpr:
		return cUnwrap(e.Expr)
	case PtrAssert:
		return cUnwrap(e.X)
	case Ident:
		return IdentExpr{e.Identifier()}
	}
	return e
}

func cUnwrapConst(e Expr) Expr {
	switch e := e.(type) {
	case *CParentExpr:
		return e.Expr
	case *CCastExpr:
		return e.Expr
	}
	return e
}

func cParenLazyOp(x Expr, pop BinaryOp) Expr {
	if x, ok := x.(*CBinaryExpr); ok {
		if pop.Precedence() <= x.Op.Precedence() {
			return x
		} else if pop.Precedence() > x.Op.Precedence() {
			return cParen(x)
		}
	}
	return cParenLazy(x)
}

func cParenLazyOpR(x Expr, pop BinaryOp) Expr {
	if x, ok := cUnwrap(x).(*CBinaryExpr); ok {
		if pop.Precedence() < x.Op.Precedence() {
			return x
		} else if pop.Precedence() > x.Op.Precedence() {
			return cParen(x)
		}
	}
	return cParenLazyR(x)
}

func cParenLazy(x Expr) Expr {
	switch x := x.(type) {
	case *CBinaryExpr, *CAssignExpr, *CTernaryExpr:
		return cParen(x)
	case *CCastExpr:
		switch x.CType(nil).(type) {
		case types.Named, types.IntType, types.FloatType:
			// nop
		default:
			return cParen(x)
		}
	}
	return x
}

func cParenLazyR(x Expr) Expr {
	switch x := x.(type) {
	case *CUnaryExpr:
		if x.Op == UnaryMinus {
			return cParen(x)
		}
	case IntLit:
		if x.IsNeg() {
			return cParen(x)
		}
	}
	return cParenLazy(x)
}

func cParen(x Expr) Expr {
	switch x := x.(type) {
	case IdentExpr:
		return x
	case PtrIdent:
		return x
	case FuncIdent:
		return x
	case IntLit:
		if !x.IsNeg() {
			return x
		}
	case FloatLit:
		return x
	case *CLiteral:
		return x
	case *CParentExpr:
		return x
	case *CCompLitExpr:
		return x
	case *CallExpr:
		return x
	case *CSelectExpr:
		return x
	case *CUnaryExpr:
		if p, ok := x.Expr.(*CParentExpr); ok {
			x.Expr = p.Expr
		}
	}
	return &CParentExpr{Expr: x}
}

type CParentExpr struct {
	Expr Expr
}

func (e *CParentExpr) Visit(v Visitor) {
	v(e.Expr)
}

func (e *CParentExpr) CType(types.Type) types.Type {
	return e.Expr.CType(nil)
}

func (e *CParentExpr) IsConst() bool {
	return e.Expr.IsConst()
}

func (e *CParentExpr) HasSideEffects() bool {
	return e.Expr.HasSideEffects()
}

func (e *CParentExpr) AsExpr() GoExpr {
	return paren(e.Expr.AsExpr())
}

func (e *CParentExpr) Uses() []types.Usage {
	return e.Expr.Uses()
}

func (g *translator) NewCIndexExpr(x, ind Expr, typ types.Type) Expr {
	ind = cUnwrap(ind)
	if types.Same(x.CType(nil), g.env.Go().String()) {
		return &CIndexExpr{
			ctype: g.env.Go().Byte(),
			Expr:  x,
			Index: ind,
		}
	}
	if addr, ok := cUnwrap(x).(*TakeAddr); ok {
		if ind2, ok := addr.X.(*CIndexExpr); ok {
			// add another component to the index expr
			return g.NewCIndexExpr(ind2.Expr,
				g.NewCBinaryExpr(ind2.Index, BinOpAdd, ind),
				typ)
		}
	}
	if types.IsPtr(x.CType(nil)) {
		x := g.ToPointer(x)
		return g.cDeref(cPtrOffset(x, ind))
	}
	if arr, ok := types.Unwrap(x.CType(nil)).(types.ArrayType); ok {
		typ = arr.Elem()
	} else if typ == nil {
		panic("no type for an element")
	}
	return &CIndexExpr{
		ctype: typ,
		Expr:  x,
		Index: ind,
	}
}

type CIndexExpr struct {
	ctype types.Type
	Expr  Expr
	Index Expr
}

func (e *CIndexExpr) IndexZero() bool {
	l, ok := unwrapCasts(e.Index).(IntLit)
	if !ok {
		return false
	}
	return l.IsZero()
}

func (e *CIndexExpr) Visit(v Visitor) {
	v(e.Expr)
	v(e.Index)
}

func (e *CIndexExpr) CType(types.Type) types.Type {
	return e.ctype
}

func (e *CIndexExpr) IsConst() bool {
	return false
}

func (e *CIndexExpr) HasSideEffects() bool {
	return e.Expr.HasSideEffects() || e.Index.HasSideEffects()
}

func (e *CIndexExpr) AsExpr() GoExpr {
	return &ast.IndexExpr{
		X:     e.Expr.AsExpr(),
		Index: e.Index.AsExpr(),
	}
}

func (e *CIndexExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, e.Expr.Uses()...)
	list = append(list, types.UseRead(e.Index)...)
	return list
}

func NewCSelectExpr(x Expr, f *types.Ident) Expr {
	x2 := cUnwrap(x)
	switch x2 := x2.(type) {
	case Ident:
		x = x2
	case *TakeAddr:
		x = x2.X
	}
	return &CSelectExpr{
		Expr: x, Sel: f,
	}
}

type CSelectExpr struct {
	Expr Expr
	Sel  *types.Ident
}

func (e *CSelectExpr) Visit(v Visitor) {
	v(e.Expr)
	v(IdentExpr{e.Sel})
}

func (e *CSelectExpr) CType(exp types.Type) types.Type {
	return e.Sel.CType(exp)
}

func (e *CSelectExpr) IsConst() bool {
	return false
}

func (e *CSelectExpr) HasSideEffects() bool {
	return e.Expr.HasSideEffects()
}

func (e *CSelectExpr) AsExpr() GoExpr {
	return &ast.SelectorExpr{
		X:   e.Expr.AsExpr(),
		Sel: e.Sel.GoIdent(),
	}
}

func (e *CSelectExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.Usage{Ident: e.Sel, Access: types.AccessUnknown})
	list = append(list, types.UseRead(e.Expr)...)
	return list
}

var (
	_ CStmtConv = &CIncrExpr{}
)

func (g *translator) NewCPrefixExpr(x Expr, decr bool) *CIncrExpr {
	return &CIncrExpr{
		g:    g,
		Expr: x, Prefix: true,
		Decr: decr,
	}
}

func (g *translator) NewCPostfixExpr(x Expr, decr bool) *CIncrExpr {
	return &CIncrExpr{
		g:    g,
		Expr: x, Prefix: false,
		Decr: decr,
	}
}

type CIncrExpr struct {
	g      *translator
	Expr   Expr
	Prefix bool
	Decr   bool
}

func (e *CIncrExpr) Visit(v Visitor) {
	v(e.Expr)
}

func (e *CIncrExpr) CType(types.Type) types.Type {
	return e.Expr.CType(nil)
}

func (e *CIncrExpr) IsConst() bool {
	return false
}

func (e *CIncrExpr) HasSideEffects() bool {
	return true
}

func (e *CIncrExpr) ToStmt() []CStmt {
	return []CStmt{e.g.NewCIncStmt(e.Expr, e.Decr)}
}

func (e *CIncrExpr) AsExpr() GoExpr {
	pi := types.NewIdent("p_", e.g.env.PtrT(e.Expr.CType(nil)))
	p := pi.GoIdent()
	y := e.g.cAddr(e.Expr).AsExpr()
	stmts := []GoStmt{
		define(p, y),
	}
	inc := (&CIncrStmt{
		g:    e.g,
		Expr: e.g.cDeref(PtrIdent{pi}),
		Decr: e.Decr,
	}).AsStmt()
	if e.Prefix {
		stmts = append(stmts, inc...)
		stmts = append(stmts,
			returnStmt(deref(p)),
		)
	} else {
		x := ident("x")
		stmts = append(stmts,
			define(x, deref(p)),
		)
		stmts = append(stmts, inc...)
		stmts = append(stmts,
			returnStmt(x),
		)
	}
	return callLambda(
		e.CType(nil).GoType(),
		stmts...,
	)
}

func (e *CIncrExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(e.Expr)...)
	list = append(list, types.UseWrite(e.Expr)...)
	return list
}

type CCastExpr struct {
	Assert bool
	Expr   Expr
	Type   types.Type
}

func (e *CCastExpr) Visit(v Visitor) {
	v(e.Expr)
}

func (e *CCastExpr) CType(_ types.Type) types.Type {
	return e.Type
}

func (e *CCastExpr) IsConst() bool {
	return e.Expr.IsConst()
}

func (e *CCastExpr) HasSideEffects() bool {
	return e.Expr.HasSideEffects()
}

func isParenSafe(v ast.Node) bool {
	switch v := v.(type) {
	case *ast.Ident:
		return true
	case *ast.ArrayType:
		return isParenSafe(v.Elt)
	default:
		return false
	}
}

func (e *CCastExpr) AsExpr() GoExpr {
	tp := e.Type.GoType()
	if e.Assert {
		return typAssert(e.Expr.AsExpr(), tp)
	}
	switch to := e.Type.(type) {
	case types.IntType:
		if from := e.Expr.CType(to); from.Kind().Is(types.UntypedFloat) {
			return call(tp, call(ident("math.Floor"), e.Expr.AsExpr()))
		}
	case types.FloatType:
		if from := e.Expr.CType(to); from.Kind().Is(types.UntypedInt) {
			return e.Expr.AsExpr()
		}
	case types.PtrType:
		if to.Elem() != nil {
			switch et := e.Expr.CType(nil).(type) {
			case types.IntType:
				return call(paren(tp), call(unsafePtr(), call(ident("uintptr"), e.Expr.AsExpr())))
			case types.PtrType:
				if et.Elem() != nil {
					return call(paren(tp), call(unsafePtr(), e.Expr.AsExpr()))
				}
			}
		}
	}
	if !isParenSafe(tp) {
		tp = paren(tp)
	}
	return call(tp, e.Expr.AsExpr())
}

func (e *CCastExpr) Uses() []types.Usage {
	return types.UseRead(e.Expr)
}

type UnaryOp string

const (
	UnarySizeof UnaryOp = "sizeof"

	UnaryPlus  UnaryOp = "+"
	UnaryMinus UnaryOp = "-"
	UnaryXor   UnaryOp = "^"
)

func (g *translator) cSizeofE(x Expr) Expr {
	switch t := types.Unwrap(x.CType(nil)).(type) {
	case types.ArrayType:
		if t.Len() <= 0 {
			return g.newUnaryExpr(UnarySizeof, x)
		}
	}
	switch x := cUnwrap(x).(type) {
	case Bool, BoolExpr:
		// workaround for C bools (they should be reported as int in some cases)
		return g.SizeofT(g.env.DefIntT(), nil)
	case StringLit:
		return &CBinaryExpr{
			Left: &CallExpr{
				Fun:  FuncIdent{g.env.Go().LenFunc()},
				Args: []Expr{x},
			},
			Op:    BinOpAdd,
			Right: cIntLit(1, 10),
		}
	}
	return g.SizeofT(x.CType(nil), nil)
}

func (g *translator) NewCUnaryExpr(op UnaryOp, x Expr) Expr {
	switch op {
	case UnarySizeof:
		return g.cSizeofE(x)
	case UnaryXor, UnaryMinus, UnaryPlus:
		if xk := x.CType(nil).Kind(); xk.IsBool() {
			x = g.cCast(g.env.DefIntT(), x)
		} else if v, ok := x.(IntLit); ok && op == UnaryMinus {
			return v.Negate()
		}
	}
	return g.newUnaryExpr(op, x)
}

func (g *translator) NewCUnaryExprT(op UnaryOp, x Expr, typ types.Type) Expr {
	if _, ok := x.(IntLit); ok && op == UnaryXor {
		x = &CCastExpr{Expr: x, Type: typ}
	}
	x = g.NewCUnaryExpr(op, x)
	if typ.Kind().IsBool() {
		return x
	}
	if op == UnarySizeof {
		return x
	}
	return g.cCast(typ, x)
}

func (g *translator) newUnaryExpr(op UnaryOp, x Expr) Expr {
	return &CUnaryExpr{
		g: g, Op: op, Expr: x,
	}
}

type CUnaryExpr struct {
	g    *translator
	Op   UnaryOp
	Expr Expr
}

func (e *CUnaryExpr) Visit(v Visitor) {
	v(e.Expr)
}

func (e *CUnaryExpr) CType(exp types.Type) types.Type {
	switch e.Op {
	case UnarySizeof:
		return e.g.env.UintPtrT()
	}
	return e.Expr.CType(exp)
}

func (e *CUnaryExpr) IsConst() bool {
	switch e.Op {
	case UnaryPlus, UnaryMinus, UnaryXor:
		return e.Expr.IsConst()
	}
	return false
}

func (e *CUnaryExpr) HasSideEffects() bool {
	return e.Expr.HasSideEffects()
}

func (e *CUnaryExpr) AsExpr() GoExpr {
	var tok token.Token
	switch e.Op {
	case UnarySizeof:
		return call(e.CType(nil).GoType(), call(ident("unsafe.Sizeof"), e.Expr.AsExpr()))
	case UnaryPlus:
		tok = token.ADD
	case UnaryMinus:
		tok = token.SUB
	case UnaryXor:
		tok = token.XOR
	default:
		panic(e.Op)
	}
	return &ast.UnaryExpr{
		Op: tok,
		X:  e.Expr.AsExpr(),
	}
}

func (e *CUnaryExpr) Uses() []types.Usage {
	return e.Expr.Uses()
}

func (g *translator) SizeofT(t types.Type, typ types.Type) Expr {
	e := &CSizeofExpr{
		rtyp: g.env.Go().Uintptr(),
		std:  true,
		Type: t,
	}
	if typ != nil && typ != e.rtyp {
		if !typ.Kind().IsInt() {
			return g.cCast(typ, e)
		}
		e.rtyp = typ
		e.std = false
	}
	return e
}

func (g *translator) AlignofT(t types.Type, typ types.Type) Expr {
	e := &CAlignofExpr{
		rtyp: g.env.Go().Uintptr(),
		Type: t,
	}
	if typ != nil && typ != e.rtyp {
		if !typ.Kind().IsInt() {
			return g.cCast(typ, e)
		}
		e.rtyp = typ
		e.std = false
	}
	return e
}

type CSizeofExpr struct {
	rtyp types.Type
	std  bool
	Type types.Type
}

func (e *CSizeofExpr) Visit(v Visitor) {}

func (e *CSizeofExpr) CType(types.Type) types.Type {
	return e.rtyp
}

func (e *CSizeofExpr) IsConst() bool {
	return true
}

func (e *CSizeofExpr) HasSideEffects() bool {
	return false
}

func (e *CSizeofExpr) Uses() []types.Usage {
	// TODO: use type
	return nil
}

func sizeOf(typ types.Type) GoExpr {
	switch typ := types.Unwrap(typ).(type) {
	case types.ArrayType:
		if !typ.IsSlice() && types.Unwrap(typ.Elem()) == types.UintT(1) {
			return intLit(typ.Len())
		}
	}
	x := typ.GoType()
	switch types.Unwrap(typ).(type) {
	case *types.StructType, types.ArrayType:
		x = &ast.CompositeLit{Type: x}
	case types.PtrType:
		x = call(x, Nil{}.AsExpr())
	case types.BoolType:
		x = call(x, boolLit(false))
	case *types.FuncType:
		x = call(ident("uintptr"), intLit(0))
	default:
		x = call(x, intLit(0))
	}
	return call(ident("unsafe.Sizeof"), x)
}

func (e *CSizeofExpr) AsExpr() GoExpr {
	x := sizeOf(e.Type)
	if !e.std {
		x = call(e.CType(nil).GoType(), x)
	}
	return x
}

type CAlignofExpr struct {
	rtyp types.Type
	std  bool
	Type types.Type
}

func (e *CAlignofExpr) Visit(v Visitor) {}

func (e *CAlignofExpr) CType(types.Type) types.Type {
	return e.rtyp
}

func (e *CAlignofExpr) IsConst() bool {
	return true
}

func (e *CAlignofExpr) HasSideEffects() bool {
	return false
}

func (e *CAlignofExpr) Uses() []types.Usage {
	// TODO: use type
	return nil
}

func zeroOf(typ types.Type) GoExpr {
	x := typ.GoType()
	switch types.Unwrap(typ).(type) {
	case *types.StructType, types.ArrayType:
		x = &ast.CompositeLit{Type: x}
	case types.PtrType:
		x = call(x, Nil{}.AsExpr())
	case types.BoolType:
		x = call(x, boolLit(false))
	default:
		x = call(x, intLit(0))
	}
	return x
}

func (e *CAlignofExpr) AsExpr() GoExpr {
	x := call(ident("unsafe.Alignof"), zeroOf(e.Type))
	if !e.std {
		x = call(e.CType(nil).GoType(), x)
	}
	return x
}

var _ CStmtConv = &CAssignExpr{}

func (g *translator) NewCAssignExpr(x Expr, op BinaryOp, y Expr) Expr {
	return &CAssignExpr{
		Stmt: g.NewCAssignStmtP(x, op, y),
	}
}

type CAssignExpr struct {
	Stmt *CAssignStmt
}

func (e *CAssignExpr) Visit(v Visitor) {
	v(e.Stmt)
}

func (e *CAssignExpr) CType(types.Type) types.Type {
	return e.Stmt.Left.CType(nil)
}

func (e *CAssignExpr) IsConst() bool {
	return false
}

func (e *CAssignExpr) HasSideEffects() bool {
	return true
}

func (e *CAssignExpr) ToStmt() []CStmt {
	return e.Stmt.g.NewCAssignStmt(e.Stmt.Left, e.Stmt.Op, e.Stmt.Right)
}

func (e *CAssignExpr) AsExpr() GoExpr {
	ret := e.CType(nil)
	var stmts []GoStmt
	if id, ok := e.Stmt.Left.(IdentExpr); ok {
		stmts = append(stmts,
			e.Stmt.g.NewCAssignStmtP(id, e.Stmt.Op, e.Stmt.Right).AsStmt()...,
		)
		stmts = append(stmts,
			returnStmt(id.GoIdent()),
		)
		return callLambda(
			ret.GoType(),
			stmts...,
		)
	}
	x, y := &TakeAddr{g: e.Stmt.g, X: e.Stmt.Left}, e.Stmt.Right
	p := types.NewIdent("p_", x.CType(nil))
	pg := p.GoIdent()
	stmts = append(stmts,
		define(pg, x.AsExpr()),
	)
	stmts = append(stmts,
		e.Stmt.g.NewCAssignStmtP(e.Stmt.g.cDeref(PtrIdent{p}), e.Stmt.Op, y).AsStmt()...,
	)
	stmts = append(stmts,
		returnStmt(deref(pg)),
	)
	return callLambda(
		ret.GoType(),
		stmts...,
	)
}

func (e *CAssignExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(e.Stmt.Left)...)
	list = append(list, types.UseWrite(e.Stmt.Left)...)
	list = append(list, types.UseRead(e.Stmt.Right)...)
	return list
}

func mergeStructInitializers(items []*CompLitField) []*CompLitField {
	out := make([]*CompLitField, 0, len(items))
	byField := make(map[string]*CompLitField)
	for _, it := range items {
		if it.Field == nil {
			out = append(out, it)
			continue
		}
		c1, ok := cUnwrap(it.Value).(*CCompLitExpr)
		if !ok {
			out = append(out, it)
			continue
		}
		it2 := byField[it.Field.Name]
		if it2 == nil {
			byField[it.Field.Name] = it
			out = append(out, it)
			continue
		}
		c2, ok := cUnwrap(it2.Value).(*CCompLitExpr)
		if !ok {
			out = append(out, it)
			continue
		}
		if !types.Same(c1.Type, c2.Type) {
			out = append(out, it)
			continue
		}
		sub := make([]*CompLitField, 0, len(c1.Fields)+len(c2.Fields))
		sub = append(sub, c2.Fields...)
		sub = append(sub, c1.Fields...)
		c2.Fields = sub
	}
	return out
}

func (g *translator) NewCCompLitExpr(typ types.Type, items []*CompLitField) Expr {
	if k := typ.Kind(); k.Is(types.Struct) {
		st := types.Unwrap(typ).(*types.StructType)
		next := 0
		for _, it := range items {
			var ft types.Type
			if it.Field != nil {
				ft = it.Field.CType(nil)
			} else if it.Index == nil {
				f := st.Fields()[next]
				next++
				it.Field = f.Name
				ft = f.Type()
			} else if l, ok := cUnwrap(it.Index).(IntLit); ok {
				next = int(l.Uint())
				f := st.Fields()[next]
				next++
				it.Field = f.Name
				ft = f.Type()
			}
			if ft != nil {
				it.Value = g.cCast(ft, it.Value)
			}
		}
		items = mergeStructInitializers(items)
	} else if k.Is(types.Array) {
		arr := types.Unwrap(typ).(types.ArrayType)
		tp := arr.Elem()
		for _, it := range items {
			it.Value = g.cCast(tp, it.Value)
		}
	}
	return &CCompLitExpr{Type: typ, Fields: items}
}

type CompLitField struct {
	Index Expr
	Field *types.Ident
	Value Expr
}

func (e *CompLitField) Visit(v Visitor) {
	v(e.Index)
	if e.Field != nil {
		v(IdentExpr{e.Field})
	}
	v(e.Value)
}

type CCompLitExpr struct {
	Type   types.Type
	Fields []*CompLitField
}

func (e *CCompLitExpr) Visit(v Visitor) {
	for _, f := range e.Fields {
		v(f)
	}
}

func (e *CCompLitExpr) CType(types.Type) types.Type {
	return e.Type
}

func (e *CCompLitExpr) IsConst() bool {
	return false
}

func (e *CCompLitExpr) HasSideEffects() bool {
	has := false
	for _, f := range e.Fields {
		has = has || f.Value.HasSideEffects()
	}
	return has
}

func (e *CCompLitExpr) isZero() bool {
	for _, f := range e.Fields {
		switch v := cUnwrapConst(f.Value).(type) {
		case IntLit:
			if !v.IsZero() {
				return false
			}
		case Nil:
			// nop
		case *CCompLitExpr:
			if !v.isZero() {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (e *CCompLitExpr) AsExpr() GoExpr {
	if e.isZero() {
		// special case: MyStruct{0} usually means the same in C as MyStruct{} in Go
		return &ast.CompositeLit{
			Type: e.Type.GoType(),
		}
	}
	kind := e.CType(nil).Kind()
	if len(e.Fields) == 1 && e.Fields[0].Value != nil && (kind.IsInt() || kind.IsFloat() || kind.IsBool()) {
		return tmpVar(e.Type.GoType(), e.Fields[0].Value.AsExpr(), false)
	}
	var items []GoExpr
	isArr := kind.Is(types.Array)
	ordered := false
	if isArr {
		// check if array elements are ordered so we can skip indexes in the init
		at, ok := types.Unwrap(e.CType(nil)).(types.ArrayType)
		if ok && at.Len() == len(e.Fields) {
			ordered = true
			for i, f := range e.Fields {
				if f.Index == nil {
					continue
				}
				v, ok := cUnwrap(f.Index).(IntLit)
				if !ok {
					ordered = false
					break
				}
				if !v.IsUint() || i != int(v.Uint()) {
					ordered = false
					break
				}
			}
		}
	}
	for _, f := range e.Fields {
		v := f.Value.AsExpr()
		if v, ok := v.(*ast.CompositeLit); ok && isArr {
			v.Type = nil
		}
		if f.Field != nil {
			items = append(items, &ast.KeyValueExpr{
				Key:   f.Field.GoIdent(),
				Value: v,
			})
		} else if f.Index != nil && !ordered {
			items = append(items, &ast.KeyValueExpr{
				Key:   f.Index.AsExpr(),
				Value: v,
			})
		} else {
			items = append(items, v)
		}
	}
	return &ast.CompositeLit{
		Type: e.Type.GoType(),
		Elts: items,
	}
}

func (e *CCompLitExpr) Uses() []types.Usage {
	var list []types.Usage
	// TODO: use the type
	for _, e := range e.Fields {
		if e.Field != nil {
			list = append(list, types.Usage{Ident: e.Field, Access: types.AccessWrite})
		}
		if e.Value != nil {
			list = append(list, types.UseRead(e.Value)...)
		}
		if e.Index != nil {
			list = append(list, types.UseRead(e.Index)...)
		}
	}
	return list
}

type TypeAssert struct {
	Type types.Type
	Expr Expr
}

func (e *TypeAssert) IsConst() bool {
	return false
}

func (e *TypeAssert) HasSideEffects() bool {
	return true
}

func (e *TypeAssert) CType(types.Type) types.Type {
	return e.Type
}

func (e *TypeAssert) AsExpr() GoExpr {
	return &ast.TypeAssertExpr{
		X:    e.Expr.AsExpr(),
		Type: e.Type.GoType(),
	}
}

type SliceExpr struct {
	Expr   Expr
	Low    Expr
	High   Expr
	Max    Expr
	Slice3 bool
}

func (e *SliceExpr) Visit(v Visitor) {
	v(e.Expr)
}

func (e *SliceExpr) IsConst() bool {
	return false
}

func (e *SliceExpr) HasSideEffects() bool {
	return false
}

func (e *SliceExpr) CType(exp types.Type) types.Type {
	return e.Expr.CType(exp)
}

func (e *SliceExpr) AsExpr() GoExpr {
	var low, high, max GoExpr
	if e.Low != nil {
		low = e.Low.AsExpr()
	}
	if e.High != nil {
		high = e.High.AsExpr()
	}
	if e.Max != nil {
		max = e.Max.AsExpr()
	}
	return &ast.SliceExpr{
		X:   e.Expr.AsExpr(),
		Low: low, High: high, Max: max,
		Slice3: e.Slice3,
	}
}

func (e *SliceExpr) Uses() []types.Usage {
	uses := e.Expr.Uses()
	uses = append(uses, types.UseRead(e.Low, e.High, e.Max)...)
	return uses
}

func exprCost(e Expr) int {
	if e == nil {
		return 0
	}
	switch e := e.(type) {
	case *CMultiExpr:
		c := 0
		for _, s := range e.Exprs {
			c += exprCost(s)
		}
		return c
	case *CallExpr:
		c := 1
		for _, s := range e.Args {
			c += exprCost(s)
		}
		return c
	case *CTernaryExpr:
		return 1 + exprCost(e.Cond) + exprCost(e.Then) + exprCost(e.Else)
	case *CAssignExpr:
		return 1 + stmtCost(e.Stmt)
	case *CBinaryExpr:
		return 1 + exprCost(e.Left) + exprCost(e.Right)
	case *CIndexExpr:
		return 1 + exprCost(e.Expr) + exprCost(e.Index)
	case *CSelectExpr:
		return 1 + exprCost(e.Expr)
	case *CCastExpr:
		return 1 + exprCost(e.Expr)
	case *CIncrExpr:
		return 1 + exprCost(e.Expr)
	case *CUnaryExpr:
		return 1 + exprCost(e.Expr)
	case *CParentExpr:
		return exprCost(e.Expr)
	default:
		return 1
	}
}
