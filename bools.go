package cxgo

import (
	"go/ast"
	"go/token"

	"github.com/gotranspile/cxgo/types"
)

// BoolExpr is a expression that returns a bool value.
type BoolExpr interface {
	Expr
	// Negate a bool expression. Alternative of !x, but for any expression.
	// It may invert the comparison operator or just return !x.
	Negate() BoolExpr
}

// ToBool converts an expression to bool expression.
func (g *translator) ToBool(x Expr) BoolExpr {
	if x, ok := cUnwrap(x).(BoolExpr); ok {
		return x
	}
	if v, ok := cIsBoolConst(x); ok {
		if v {
			return Bool(true)
		}
		return Bool(false)
	}
	if x, ok := cUnwrap(x).(IntLit); ok {
		if x.IsZero() {
			return Bool(false)
		} else if x.IsOne() {
			return Bool(true)
		}
	}
	if x.CType(nil).Kind().IsBool() {
		if x, ok := x.(Ident); ok {
			return BoolIdent{x.Identifier()}
		}
		return BoolAssert{x}
	}
	if types.IsPtr(x.CType(nil)) {
		return ComparePtrs(
			g.ToPointer(x),
			BinOpNeq,
			g.Nil(),
		)
	}
	return g.Compare(
		x,
		BinOpNeq,
		cIntLit(0),
	)
}

func cIsBoolConst(x Expr) (bool, bool) {
	x = cUnwrap(x)
	switch x := x.(type) {
	case Bool:
		if x {
			return true, true
		}
		return false, true
	case IntLit:
		if x.IsOne() {
			return true, true
		} else if x.IsZero() {
			return false, true
		}
	}
	return false, false
}

var (
	_ BoolExpr = Bool(false)
)

// Bool is a constant bool value.
type Bool bool

func (Bool) Visit(v Visitor) {}

func (e Bool) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e Bool) AsExpr() GoExpr {
	if e {
		return ident("true")
	}
	return ident("false")
}

func (e Bool) IsConst() bool {
	return true
}

func (e Bool) HasSideEffects() bool {
	return false
}

func (e Bool) Negate() BoolExpr {
	return !e
}

func (e Bool) Uses() []types.Usage {
	return nil
}

var (
	_ BoolExpr = BoolIdent{}
	_ Ident    = BoolIdent{}
)

type BoolIdent struct {
	*types.Ident
}

func (BoolIdent) Visit(v Visitor) {}

func (e BoolIdent) Identifier() *types.Ident {
	return e.Ident
}

func (e BoolIdent) IsConst() bool {
	return false
}

func (e BoolIdent) HasSideEffects() bool {
	return false
}

func (e BoolIdent) AsExpr() GoExpr {
	return e.GoIdent()
}

func (e BoolIdent) Negate() BoolExpr {
	return &Not{X: e}
}

func (e BoolIdent) Uses() []types.Usage {
	return []types.Usage{{Ident: e.Ident, Access: types.AccessUnknown}}
}

var _ BoolExpr = (*Not)(nil)

// Not negates a bool expression. It's only useful for identifiers and function calls.
type Not struct {
	X BoolExpr
}

func (e *Not) Visit(v Visitor) {
	v(e.X)
}

func (e *Not) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e *Not) AsExpr() GoExpr {
	return &ast.UnaryExpr{
		Op: token.NOT,
		X:  e.X.AsExpr(),
	}
}

func (e *Not) IsConst() bool {
	return e.X.IsConst()
}

func (e *Not) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *Not) Negate() BoolExpr {
	return e.X
}

func (e *Not) Uses() []types.Usage {
	return e.X.Uses()
}

func (g *translator) cNot(x Expr) BoolExpr {
	if x, ok := x.(BoolExpr); ok {
		return x.Negate()
	}
	return g.ToBool(x).Negate()
}

const (
	BinOpAnd BoolOp = "&&"
	BinOpOr  BoolOp = "||"
)

type BoolOp string

func (op BoolOp) Negate() BoolOp {
	switch op {
	case BinOpAnd:
		return BinOpOr
	case BinOpOr:
		return BinOpAnd
	}
	panic(op)
}

func (op BoolOp) GoToken() token.Token {
	var tok token.Token
	switch op {
	case BinOpAnd:
		tok = token.LAND
	case BinOpOr:
		tok = token.LOR
	default:
		panic(op)
	}
	return tok
}

func And(x, y BoolExpr) BoolExpr {
	return &BinaryBoolExpr{
		X: x, Op: BinOpAnd, Y: y,
	}
}

func Or(x, y BoolExpr) BoolExpr {
	return &BinaryBoolExpr{
		X: x, Op: BinOpOr, Y: y,
	}
}

var _ BoolExpr = (*BinaryBoolExpr)(nil)

type BinaryBoolExpr struct {
	X  BoolExpr
	Op BoolOp
	Y  BoolExpr
}

func (e *BinaryBoolExpr) Visit(v Visitor) {
	v(e.X)
	v(e.Y)
}

func (e *BinaryBoolExpr) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e *BinaryBoolExpr) AsExpr() GoExpr {
	return &ast.BinaryExpr{
		X:  e.X.AsExpr(),
		Op: e.Op.GoToken(),
		Y:  e.Y.AsExpr(),
	}
}

func (e *BinaryBoolExpr) IsConst() bool {
	return e.X.IsConst() && e.Y.IsConst()
}

func (e *BinaryBoolExpr) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Y.HasSideEffects()
}

func (e *BinaryBoolExpr) Negate() BoolExpr {
	return &BinaryBoolExpr{
		X:  e.X.Negate(),
		Op: e.Op.Negate(),
		Y:  e.Y.Negate(),
	}
}

func (e *BinaryBoolExpr) Uses() []types.Usage {
	return types.UseRead(e.X, e.Y)
}

const (
	BinOpEq  ComparisonOp = "=="
	BinOpNeq ComparisonOp = "!="
	BinOpLt  ComparisonOp = "<"
	BinOpGt  ComparisonOp = ">"
	BinOpLte ComparisonOp = "<="
	BinOpGte ComparisonOp = ">="
)

// ComparisonOp is a comparison operator.
type ComparisonOp string

func (op ComparisonOp) IsEquality() bool {
	return op == BinOpEq || op == BinOpNeq
}

func (op ComparisonOp) IsRelational() bool {
	switch op {
	case BinOpLt, BinOpGt, BinOpLte, BinOpGte:
		return true
	}
	return false
}

func (op ComparisonOp) Negate() ComparisonOp {
	switch op {
	case BinOpEq:
		return BinOpNeq
	case BinOpNeq:
		return BinOpEq
	case BinOpLt:
		return BinOpGte
	case BinOpGt:
		return BinOpLte
	case BinOpLte:
		return BinOpGt
	case BinOpGte:
		return BinOpLt
	}
	panic(op)
}

func (op ComparisonOp) GoToken() token.Token {
	var tok token.Token
	switch op {
	case BinOpLt:
		tok = token.LSS
	case BinOpGt:
		tok = token.GTR
	case BinOpLte:
		tok = token.LEQ
	case BinOpGte:
		tok = token.GEQ
	case BinOpEq:
		tok = token.EQL
	case BinOpNeq:
		tok = token.NEQ
	default:
		panic(op)
	}
	return tok
}

// Compare two expression values.
func (g *translator) Compare(x Expr, op ComparisonOp, y Expr) BoolExpr {
	// compare pointers and functions separately
	if xt := x.CType(nil); xt.Kind().IsFunc() {
		fx := g.ToFunc(x, nil)
		return CompareFuncs(fx, op, g.ToFunc(y, fx.FuncType(nil)))
	}
	if yt := y.CType(nil); yt.Kind().IsFunc() {
		fy := g.ToFunc(y, nil)
		return CompareFuncs(g.ToFunc(x, fy.FuncType(nil)), op, fy)
	}
	if x.CType(nil).Kind().IsPtr() || y.CType(nil).Kind().IsPtr() {
		return ComparePtrs(g.ToPointer(x), op, g.ToPointer(y))
	}
	if x.CType(nil).Kind().Is(types.Array) || y.CType(nil).Kind().Is(types.Array) {
		return ComparePtrs(g.ToPointer(x), op, g.ToPointer(y))
	}
	if op.IsRelational() {
		typ := g.env.CommonType(x.CType(nil), y.CType(nil))
		x = g.cCast(typ, x)
		y = g.cCast(typ, y)
		return &Comparison{
			g: g,
			X: x, Op: op, Y: y,
		}
	}
	if !op.IsEquality() {
		panic("must not happen")
	}
	// always check equality with the constant on the right
	if x.IsConst() && !y.IsConst() {
		return g.Compare(y, op, x)
	}
	// optimizations for bool equality
	if v, ok := cIsBoolConst(y); ok {
		x := cUnwrap(x)
		if x.CType(nil).Kind().IsBool() {
			x = g.ToBool(x)
		}
		if x, ok := x.(BoolExpr); ok {
			if v == (op == BinOpEq) {
				// (bool(x) == true) -> (x)
				// (bool(x) != false) -> (x)
				return x
			}
			// (bool(x) != true) -> (!x)
			// (bool(x) == false) -> (!x)
			return x.Negate()
		}
	}
	typ := g.env.CommonType(x.CType(nil), y.CType(nil))
	x = g.cCast(typ, x)
	y = g.cCast(typ, y)
	return &Comparison{
		g: g,
		X: x, Op: op, Y: y,
	}
}

var _ BoolExpr = (*Comparison)(nil)

type Comparison struct {
	g  *translator
	X  Expr
	Op ComparisonOp
	Y  Expr
}

func (e *Comparison) Visit(v Visitor) {
	v(e.X)
	v(e.Y)
}

func (e *Comparison) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e *Comparison) AsExpr() GoExpr {
	x := e.X.AsExpr()
	if _, ok := x.(*ast.CompositeLit); ok {
		x = paren(x)
	}
	y := e.Y.AsExpr()
	if _, ok := y.(*ast.CompositeLit); ok {
		y = paren(y)
	}
	return &ast.BinaryExpr{
		X:  x,
		Op: e.Op.GoToken(),
		Y:  y,
	}
}

func (e *Comparison) IsConst() bool {
	return e.X.IsConst() && e.Y.IsConst()
}

func (e *Comparison) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Y.HasSideEffects()
}

func (e *Comparison) Negate() BoolExpr {
	return e.g.Compare(e.X, e.Op.Negate(), e.Y)
}

func (e *Comparison) Uses() []types.Usage {
	return types.UseRead(e.X, e.Y)
}

var _ BoolExpr = BoolAssert{}

type BoolAssert struct {
	X Expr
}

func (e BoolAssert) Visit(v Visitor) {
	v(e.X)
}

func (e BoolAssert) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e BoolAssert) AsExpr() GoExpr {
	return e.X.AsExpr()
}

func (e BoolAssert) IsConst() bool {
	return false
}

func (e BoolAssert) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e BoolAssert) Negate() BoolExpr {
	return &Not{X: e}
}

func (e BoolAssert) Uses() []types.Usage {
	return e.X.Uses()
}

var _ Expr = (*BoolToInt)(nil)

type BoolToInt struct {
	X BoolExpr
}

func (e *BoolToInt) Visit(v Visitor) {
	v(e.X)
}

func (e *BoolToInt) CType(types.Type) types.Type {
	return types.IntT(4)
}

func (e *BoolToInt) AsExpr() GoExpr {
	return call(ident("libc.BoolToInt"), e.X.AsExpr())
}

func (e *BoolToInt) IsConst() bool {
	return false
}

func (e *BoolToInt) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *BoolToInt) Uses() []types.Usage {
	return e.X.Uses()
}
