package cxgo

import (
	"go/ast"
	"go/token"

	"github.com/dennwc/cxgo/types"
)

// PtrExpr is an expression that returns a pointer value.
type PtrExpr interface {
	Expr
	// PtrType return an underlying pointer type.
	PtrType(exp types.PtrType) types.PtrType
}

// ToPointer asserts that a value is a pointer. If it's not, it is converted to a void pointer.
func (g *translator) ToPointer(x Expr) PtrExpr {
	if x, ok := x.(PtrExpr); ok {
		return x
	}
	if x, ok := cUnwrap(x).(IntLit); ok && x.IsZero() {
		return g.Nil()
	}
	xt := x.CType(nil)
	xk := xt.Kind()
	if xk.IsFunc() {
		return &FuncToPtr{
			e: g.env.Env, X: g.ToFunc(x, nil),
		}
	}
	if xk.IsPtr() {
		switch x := x.(type) {
		case Ident:
			return PtrIdent{x.Identifier()}
		}
	}
	if xk.Is(types.Array) {
		// &x[0]
		return g.cAddr(g.NewCIndexExpr(x, cIntLit(0), nil))
	}
	if xk.IsInt() {
		return g.cIntToPtr(x)
	}
	if xk.IsBool() {
		return g.cIntToPtr(&BoolToInt{X: g.ToBool(x)})
	}
	if l, ok := x.(StringLit); ok {
		// string -> *wchar_t
		// string -> *byte
		return g.StringToPtr(l)
	}
	if x.CType(nil).Kind().IsPtr() {
		return PtrAssert{X: x}
	}
	return g.cIntToPtr(x)
}

var _ PtrExpr = Nil{}

func NewNil(size int) Nil {
	return Nil{size: size}
}

func unwrapCasts(e Expr) Expr {
	switch e := e.(type) {
	case *CParentExpr:
		return unwrapCasts(e.Expr)
	case *CCastExpr:
		return unwrapCasts(e.Expr)
	}
	return e
}

func IsNil(e Expr) bool {
	e = unwrapCasts(e)
	if e == nil {
		return false
	}
	_, ok := e.(Nil)
	return ok
}

type Nil struct {
	size int
}

func (Nil) Visit(v Visitor) {}

func (e Nil) CType(exp types.Type) types.Type {
	switch t := types.Unwrap(exp).(type) {
	case types.PtrType:
		return e.PtrType(t)
	case *types.FuncType:
		return e.FuncType(t)
	}
	return e.PtrType(nil)
}

func (Nil) AsExpr() GoExpr {
	return ident("nil")
}

func (Nil) IsConst() bool {
	return true
}

func (Nil) HasSideEffects() bool {
	return false
}

func (e Nil) PtrType(exp types.PtrType) types.PtrType {
	if exp != nil {
		return exp
	}
	return types.NilT(e.size)
}

func (e Nil) FuncType(exp *types.FuncType) *types.FuncType {
	return exp
}

func (e Nil) Uses() []types.Usage {
	return nil
}

var (
	_ PtrExpr = PtrIdent{}
	_ Ident   = PtrIdent{}
)

type PtrIdent struct {
	*types.Ident
}

func (PtrIdent) Visit(v Visitor) {}

func (e PtrIdent) Identifier() *types.Ident {
	return e.Ident
}

func (e PtrIdent) IsConst() bool {
	return false
}

func (e PtrIdent) HasSideEffects() bool {
	return false
}

func (e PtrIdent) AsExpr() GoExpr {
	return e.GoIdent()
}

func (e PtrIdent) PtrType(exp types.PtrType) types.PtrType {
	var et types.Type
	if exp != nil {
		et = exp
	}
	return types.Unwrap(e.CType(et)).(types.PtrType)
}

func (e PtrIdent) Uses() []types.Usage {
	return []types.Usage{{Ident: e.Ident, Access: types.AccessUnknown}}
}

var _ PtrExpr = PtrAssert{}

type PtrAssert struct {
	X Expr
}

func (e PtrAssert) Visit(v Visitor) {
	v(e.X)
}

func (e PtrAssert) CType(types.Type) types.Type {
	return e.X.CType(nil)
}

func (e PtrAssert) IsConst() bool {
	return e.X.IsConst()
}

func (e PtrAssert) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e PtrAssert) AsExpr() GoExpr {
	return e.X.AsExpr()
}

func (e PtrAssert) PtrType(types.PtrType) types.PtrType {
	return types.Unwrap(e.CType(nil)).(types.PtrType)
}

func (e PtrAssert) Uses() []types.Usage {
	return e.X.Uses()
}

var _ PtrExpr = (*TakeAddr)(nil)

func (g *translator) cAddr(x Expr) PtrExpr {
	if x, ok := cUnwrap(x).(*Deref); ok {
		return x.X
	}
	xt := x.CType(nil)
	if xt.Kind().IsFunc() {
		return g.ToPointer(x)
	}
	if xt.Kind().Is(types.Array) {
		return g.cAddr(g.NewCIndexExpr(x, cUintLit(0), nil))
	}
	return &TakeAddr{g: g, X: x}
}

var _ PtrExpr = (*TakeAddr)(nil)

type TakeAddr struct {
	g *translator
	X Expr
}

func (e *TakeAddr) Visit(v Visitor) {
	v(e.X)
}

func (e *TakeAddr) CType(exp types.Type) types.Type {
	return e.PtrType(types.ToPtrType(exp))
}

func (e *TakeAddr) PtrType(types.PtrType) types.PtrType {
	return e.g.env.PtrT(e.X.CType(nil))
}

func (e *TakeAddr) AsExpr() GoExpr {
	return addr(e.X.AsExpr())
}

func (e *TakeAddr) IsConst() bool {
	return false
}

func (e *TakeAddr) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *TakeAddr) Uses() []types.Usage {
	return types.UseRead(e.X)
}

func (g *translator) cDeref(x PtrExpr) Expr {
	if x, ok := cUnwrap(x).(*TakeAddr); ok {
		return cParenLazy(x.X)
	}
	switch x := cUnwrap(x).(type) {
	case *FuncToPtr:
		// C has pointers to function, as opposed to Go where function variables are always pointers
		// so we should unwrap the cast and return the function identifier itself
		return x.X
	case *CCastExpr:
		if p, ok := x.Type.(types.PtrType); ok && p.Elem().Kind().IsFunc() {
			panic("TODO")
			//return &TypeAssert{
			//	Type: p.Elem(),
			//	Expr: g.NewCCallExpr(pptrFunc, []Expr{c.Expr}),
			//}
		}
	}
	return &Deref{g: g, X: x}
}

func (g *translator) cDerefT(x Expr, typ types.Type) Expr {
	return g.cDeref(g.ToPointer(x))
}

var _ Expr = (*Deref)(nil)

type Deref struct {
	g *translator
	X PtrExpr
}

func (e *Deref) Visit(v Visitor) {
	v(e.X)
}

func (e *Deref) CType(types.Type) types.Type {
	elem := e.X.PtrType(nil).Elem()
	if elem == nil {
		// cannot deref unsafe pointers directly
		elem = e.g.env.Go().Byte()
	}
	return elem
}

func (e *Deref) AsExpr() GoExpr {
	return deref(e.X.AsExpr())
}

func (e *Deref) IsConst() bool {
	return false
}

func (e *Deref) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *Deref) Uses() []types.Usage {
	return e.X.Uses()
}

var _ PtrExpr = (*PtrToPtr)(nil)

func (g *translator) cPtrToPtr(typ types.Type, x PtrExpr) PtrExpr {
	ptyp := types.UnwrapPtr(typ)
	if ptyp == nil {
		panic("must not be nil")
	}
	if types.Same(typ, x.CType(nil)) {
		return x
	}
	switch x := x.(type) {
	case Nil:
		return x
	case *IntToPtr:
		return &IntToPtr{
			X:  x.X,
			To: ptyp,
		}
	case *PtrOffset:
		return &PtrOffset{
			X:    x.X,
			Ind:  x.Ind,
			Conv: &ptyp,
		}
	case *PtrVarOffset:
		return &PtrVarOffset{
			X:    x.X,
			Mul:  x.Mul,
			Ind:  x.Ind,
			Conv: &ptyp,
		}
	case *PtrToPtr:
		if typ.Kind().IsUnsafePtr() && x.X.PtrType(nil).Kind().IsUnsafePtr() {
			return x.X
		}
	}
	if ptyp.Elem() != nil {
		if s, ok := types.Unwrap(x.PtrType(nil).Elem()).(*types.StructType); ok {
			fields := s.Fields()
			if len(fields) != 0 && types.Same(ptyp.Elem(), fields[0].Type()) {
				return g.cAddr(NewCSelectExpr(x, fields[0].Name))
			}
		}
	}
	return &PtrToPtr{
		X:  x,
		To: ptyp,
	}
}

type PtrToPtr struct {
	X  PtrExpr
	To types.PtrType
}

func (e *PtrToPtr) Visit(v Visitor) {
	v(e.X)
}

func (e *PtrToPtr) CType(types.Type) types.Type {
	return e.To
}

func (e *PtrToPtr) PtrType(types.PtrType) types.PtrType {
	return e.To
}

func (e *PtrToPtr) AsExpr() GoExpr {
	var tp GoExpr = e.To.GoType()
	if e.To.Elem() != nil {
		switch et := e.X.CType(e.To).(type) {
		case types.IntType:
			panic("IntToPtr must be used instead")
		case types.PtrType:
			if et.Elem() != nil {
				if _, ok := e.To.(types.Named); !ok {
					tp = paren(tp)
				}
				return call(tp, call(unsafePtr(), e.X.AsExpr()))
			}
		}
		tp = paren(tp)
	}
	return call(tp, e.X.AsExpr())
}

func (e *PtrToPtr) IsConst() bool {
	return e.X.IsConst()
}

func (e *PtrToPtr) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *PtrToPtr) Uses() []types.Usage {
	// TODO: use type
	return e.X.Uses()
}

var _ PtrExpr = (*IntToPtr)(nil)

func (g *translator) cIntToPtr(x Expr) PtrExpr {
	switch x.CType(nil).(type) {
	case types.PtrType:
		panic("PtrToPtr must be used")
	}
	if l, ok := cUnwrap(x).(IntLit); ok && !l.IsUint() {
		x = l.OverflowUint(g.env.PtrSize())
	}
	return &IntToPtr{
		X:  x,
		To: g.env.PtrT(nil),
	}
}

type IntToPtr struct {
	X  Expr
	To types.PtrType
}

func (e *IntToPtr) Visit(v Visitor) {
	v(e.X)
}

func (e *IntToPtr) CType(types.Type) types.Type {
	return e.To
}

func (e *IntToPtr) PtrType(types.PtrType) types.PtrType {
	return e.To
}

func (e *IntToPtr) AsExpr() GoExpr {
	tp := e.To.GoType()
	x := e.X.AsExpr()
	x = call(ident("uintptr"), x)
	if e.To.Elem() == nil {
		return call(tp, x)
	}
	switch e.X.CType(nil).(type) {
	case types.PtrType:
		panic("PtrToPtr must be used")
	}
	return call(paren(tp), call(unsafePtr(), x))
}

func (e *IntToPtr) IsConst() bool {
	return e.X.IsConst()
}

func (e *IntToPtr) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *IntToPtr) Uses() []types.Usage {
	// TODO: use type
	return e.X.Uses()
}

var _ Expr = (*PtrToInt)(nil)

func (g *translator) cPtrToInt(typ types.Type, x PtrExpr) Expr {
	if typ == nil {
		typ = g.env.DefUintT()
	}
	return &PtrToInt{
		X:  x,
		To: typ,
	}
}

type PtrToInt struct {
	X  PtrExpr
	To types.Type
}

func (e *PtrToInt) Visit(v Visitor) {
	v(e.X)
}

func (e *PtrToInt) CType(types.Type) types.Type {
	return e.To
}

func (e *PtrToInt) AsExpr() GoExpr {
	x := e.X.AsExpr()
	// TODO: handle functions
	if !e.X.CType(nil).Kind().IsUnsafePtr() {
		x = call(unsafePtr(), x)
	}
	x = call(ident("uintptr"), x)
	return call(e.CType(nil).GoType(), x)
}

func (e *PtrToInt) IsConst() bool {
	return e.X.IsConst()
}

func (e *PtrToInt) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *PtrToInt) Uses() []types.Usage {
	// TODO: use type
	return e.X.Uses()
}

func ComparePtrs(x PtrExpr, op ComparisonOp, y PtrExpr) BoolExpr {
	if op.IsRelational() {
		return &PtrComparison{
			X: x, Op: op, Y: y,
		}
	}
	if !op.IsEquality() {
		panic("must not happen")
	}
	// == and != simplifications
	// always compare with the constant on the right
	if x.IsConst() && !y.IsConst() {
		return ComparePtrs(y, op, x)
	}
	return &PtrComparison{
		X: x, Op: op, Y: y,
	}
}

var _ BoolExpr = (*PtrComparison)(nil)

type PtrComparison struct {
	X  PtrExpr
	Op ComparisonOp
	Y  PtrExpr
}

func (e *PtrComparison) Visit(v Visitor) {
	v(e.X)
	v(e.Y)
}

func (e *PtrComparison) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e *PtrComparison) AsExpr() GoExpr {
	x := e.X.AsExpr()
	y := e.Y.AsExpr()
	if e.Op.IsRelational() {
		// compare via uintptr
		if !e.X.CType(nil).Kind().IsUnsafePtr() {
			x = call(unsafePtr(), x)
		}
		if !e.Y.CType(nil).Kind().IsUnsafePtr() {
			y = call(unsafePtr(), y)
		}
		x = call(ident("uintptr"), x)
		y = call(ident("uintptr"), y)
	} else {
		xt := e.X.CType(nil)
		yt := e.Y.CType(nil)
		if IsNil(e.X) || IsNil(e.Y) {
			// compare directly
			if IsNil(e.Y) {
				// TODO: workaround for slice comparison with nil
				if addr, ok := e.X.(*TakeAddr); ok {
					if ind, ok := addr.X.(*CIndexExpr); ok {
						if ind.IndexZero() {
							x = ind.Expr.AsExpr()
						}
					}
				}
			}
		} else if e.X.IsConst() || e.Y.IsConst() {
			// compare as uintptr in case of consts
			if e.X.IsConst() {
				if xc, ok := e.X.(*IntToPtr); ok {
					x = xc.X.AsExpr()
				}
			} else {
				if !xt.Kind().IsUnsafePtr() {
					x = call(unsafePtr(), x)
				}
			}
			x = call(ident("uintptr"), x)

			if e.Y.IsConst() {
				if yc, ok := e.Y.(*IntToPtr); ok {
					y = yc.X.AsExpr()
				}
			} else {
				if !yt.Kind().IsUnsafePtr() {
					y = call(unsafePtr(), y)
				}
			}
			y = call(ident("uintptr"), y)
		} else {
			if !types.Same(xt, yt) {
				// compare as unsafe pointers
				if !xt.Kind().IsUnsafePtr() {
					x = call(unsafePtr(), x)
				}
				if !yt.Kind().IsUnsafePtr() {
					y = call(unsafePtr(), y)
				}
			}
		}
	}
	return &ast.BinaryExpr{
		X:  x,
		Op: e.Op.GoToken(),
		Y:  y,
	}
}

func (e *PtrComparison) IsConst() bool {
	return e.X.IsConst() && e.Y.IsConst()
}

func (e *PtrComparison) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Y.HasSideEffects()
}

func (e *PtrComparison) Negate() BoolExpr {
	return ComparePtrs(e.X, e.Op.Negate(), e.Y)
}

func (e *PtrComparison) Uses() []types.Usage {
	return types.UseRead(e.X, e.Y)
}

func cPtrOffset(x PtrExpr, ind Expr) PtrExpr {
	elem := x.PtrType(nil).ElemSizeof()
	if elem != 1 {
		return &PtrElemOffset{
			X: x, Ind: ind,
		}
	}
	y := ind
	sub := false
	if u, ok := cUnwrap(y).(*CUnaryExpr); ok && u.Op == UnaryMinus {
		y = u.Expr
		sub = true
	}
	mul := elem
	if sub {
		mul = -mul
	}
	if y, ok := cUnwrap(ind).(IntLit); ok {
		y = y.MulLit(int64(mul))
		return &PtrOffset{
			X: x, Ind: y.Int(),
		}
	}
	// TODO: check the ind for a const multiplier
	return &PtrVarOffset{
		X: x, Mul: mul, Ind: y,
	}
}

var _ PtrExpr = (*PtrOffset)(nil)

// PtrOffset acts like a pointer arithmetic with a constant value.
// The index is NOT multiplied by the pointer element size.
// It optionally converts the pointer to Conv type.
type PtrOffset struct {
	X    PtrExpr
	Ind  int64
	Conv *types.PtrType
}

func (e *PtrOffset) Visit(v Visitor) {
	v(e.X)
}

func (e *PtrOffset) CType(exp types.Type) types.Type {
	return e.PtrType(types.ToPtrType(exp))
}

func (e *PtrOffset) toType() types.PtrType {
	if e.Conv == nil {
		return e.X.PtrType(nil)
	}
	return *e.Conv
}

func (e *PtrOffset) parts() (GoExpr, BinaryOp, GoExpr) {
	x := e.X.AsExpr()
	ind := e.Ind
	op := BinOpAdd
	if ind < 0 {
		ind = -ind
		op = BinOpSub
	}
	y := intLit64(ind)
	if e.X.PtrType(nil).Elem() != nil {
		x = call(unsafePtr(), x)
	}
	x = call(ident("uintptr"), x)
	return x, op, y
}

func (e *PtrOffset) AsExpr() GoExpr {
	x, tok, y := e.parts()
	//if !VirtualPtrs {
	x = &ast.BinaryExpr{
		X: x, Op: tok.GoToken(), Y: y,
	}
	to := e.toType()
	if to.Elem() != nil {
		x = call(unsafePtr(), x)
	}
	return call(to.GoType(), x)
	//}
	//panic("implement me")
}

func (e *PtrOffset) IsConst() bool {
	return e.X.IsConst()
}

func (e *PtrOffset) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *PtrOffset) PtrType(exp types.PtrType) types.PtrType {
	if e.Conv != nil {
		return *e.Conv
	}
	return e.X.PtrType(exp)
}

func (e *PtrOffset) Uses() []types.Usage {
	// TODO: use the type
	return e.X.Uses()
}

var _ PtrExpr = (*PtrElemOffset)(nil)

// PtrElemOffset acts like a pointer arithmetic with a variable value.
// The index is always multiplied only by the pointer elements size, as opposed to PtrVarOffset.
// This operation is preferable over PtrVarOffset because the size of C and Go structs may not match.
type PtrElemOffset struct {
	X    PtrExpr
	Ind  Expr
	Conv *types.PtrType
}

func (e *PtrElemOffset) Visit(v Visitor) {
	v(e.X)
	v(e.Ind)
}

func (e *PtrElemOffset) CType(exp types.Type) types.Type {
	return e.PtrType(types.ToPtrType(exp))
}

func (e *PtrElemOffset) AsExpr() GoExpr {
	var to types.PtrType
	if e.Conv == nil {
		to = e.X.PtrType(nil)
	} else {
		to = *e.Conv
	}
	x := e.X.AsExpr()
	ind := e.Ind.AsExpr()
	//if !VirtualPtrs {
	op := BinOpAdd
	if u, ok := cUnwrap(e.Ind).(*CUnaryExpr); ok && u.Op == UnaryMinus {
		op = BinOpSub
		ind = u.Expr.AsExpr()
	} else if l, ok := cUnwrap(e.Ind).(IntLit); ok && !l.IsUint() {
		op = BinOpSub
		ind = l.NegateLit().AsExpr()
	}
	y := ind
	if e.X.PtrType(nil).Elem() != nil {
		x = call(unsafePtr(), x)
	}
	indTyp := e.Ind.CType(nil)
	szof := sizeOf(e.X.PtrType(nil).Elem())
	if e.Ind.IsConst() && indTyp.Kind().IsInt() {
		y = &ast.BinaryExpr{
			X:  szof,
			Op: token.MUL,
			Y:  y,
		}
	} else {
		y = &ast.BinaryExpr{
			X:  szof,
			Op: token.MUL,
			Y:  call(ident("uintptr"), y),
		}
	}
	x = &ast.BinaryExpr{
		X:  call(ident("uintptr"), x),
		Op: op.GoToken(),
		Y:  y,
	}
	if to.Elem() != nil {
		x = call(unsafePtr(), x)
	}
	return call(to.GoType(), x)
	//}
	//panic("implement me")
}

func (e *PtrElemOffset) IsConst() bool {
	return e.X.IsConst() && e.Ind.IsConst()
}

func (e *PtrElemOffset) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Ind.HasSideEffects()
}

func (e *PtrElemOffset) PtrType(exp types.PtrType) types.PtrType {
	if e.Conv != nil {
		return *e.Conv
	}
	return e.X.PtrType(exp)
}

func (e *PtrElemOffset) Uses() []types.Usage {
	// TODO: use the type
	return types.UseRead(e.X, e.Ind)
}

var _ PtrExpr = (*PtrOffset)(nil)

// PtrVarOffset acts like a pointer arithmetic with a variable value.
// The index is multiplied only by a Mul, but not the pointer element size.
// It optionally converts the pointer to Conv type.
type PtrVarOffset struct {
	X    PtrExpr
	Mul  int
	Ind  Expr
	Conv *types.PtrType
}

func (e *PtrVarOffset) Visit(v Visitor) {
	v(e.X)
	v(e.Ind)
}

func (e *PtrVarOffset) CType(exp types.Type) types.Type {
	return e.PtrType(types.ToPtrType(exp))
}

func (e *PtrVarOffset) AsExpr() GoExpr {
	var to types.PtrType
	if e.Conv == nil {
		to = e.X.PtrType(nil)
	} else {
		to = *e.Conv
	}
	x := e.X.AsExpr()
	ind := e.Ind.AsExpr()
	//if !VirtualPtrs {
	mul := e.Mul
	op := BinOpAdd
	if mul < 0 {
		mul = -mul
		op = BinOpSub
	}
	y := ind
	if e.X.PtrType(nil).Elem() != nil {
		x = call(unsafePtr(), x)
	}
	if mul != 1 {
		y = &ast.BinaryExpr{
			X:  intLit(mul),
			Op: token.MUL,
			Y:  call(ident("uintptr"), y),
		}
	} else {
		y = call(ident("uintptr"), y)
	}
	x = &ast.BinaryExpr{
		X:  call(ident("uintptr"), x),
		Op: op.GoToken(),
		Y:  y,
	}
	if to.Elem() != nil {
		x = call(unsafePtr(), x)
	}
	return call(to.GoType(), x)
	//}
	//panic("implement me")
}

func (e *PtrVarOffset) IsConst() bool {
	return e.X.IsConst() && e.Ind.IsConst()
}

func (e *PtrVarOffset) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Ind.HasSideEffects()
}

func (e *PtrVarOffset) PtrType(exp types.PtrType) types.PtrType {
	if e.Conv != nil {
		return *e.Conv
	}
	return e.X.PtrType(exp)
}

func (e *PtrVarOffset) Uses() []types.Usage {
	// TODO: use the type
	return types.UseRead(e.X, e.Ind)
}

func cPtrDiff(x, y PtrExpr) Expr {
	return &PtrDiff{X: x, Y: y}
}

var _ Expr = (*PtrDiff)(nil)

type PtrDiff struct {
	X PtrExpr
	Y PtrExpr
}

func (e *PtrDiff) Visit(v Visitor) {
	v(e.X)
	v(e.Y)
}

func (e *PtrDiff) CType(exp types.Type) types.Type {
	return types.IntT(e.X.PtrType(nil).Sizeof())
}

func (e *PtrDiff) AsExpr() GoExpr {
	x, y := e.X.AsExpr(), e.Y.AsExpr()
	if e.X.PtrType(nil).Elem() != nil {
		x = call(unsafePtr(), x)
	}
	if e.Y.PtrType(nil).Elem() != nil {
		y = call(unsafePtr(), y)
	}
	return call(e.CType(nil).GoType(), &ast.BinaryExpr{
		X:  call(ident("uintptr"), x),
		Op: token.SUB,
		Y:  call(ident("uintptr"), y),
	})
}

func (e *PtrDiff) IsConst() bool {
	return false
}

func (e *PtrDiff) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Y.HasSideEffects()
}

func (e *PtrDiff) Uses() []types.Usage {
	return types.UseRead(e.X, e.Y)
}

var _ PtrExpr = (*StringToPtr)(nil)

func (g *translator) StringToPtr(x StringLit) PtrExpr {
	return &StringToPtr{e: g.env.Env, X: x}
}

type StringToPtr struct {
	e *types.Env
	X StringLit
}

func (e *StringToPtr) Visit(v Visitor) {
	v(e.X)
}

func (e *StringToPtr) CType(exp types.Type) types.Type {
	return e.PtrType(types.ToPtrType(exp))
}

func (e *StringToPtr) AsExpr() GoExpr {
	if e.X.IsWide() {
		return call(ident("cstd.CWString"), e.X.AsExpr())
	}
	return call(ident("cstd.CString"), e.X.AsExpr())
}

func (e *StringToPtr) IsConst() bool {
	return false
}

func (e *StringToPtr) HasSideEffects() bool {
	return false
}

func (e *StringToPtr) PtrType(types.PtrType) types.PtrType {
	if e.X.IsWide() {
		return e.e.C().WString()
	}
	return e.e.C().String()
}

func (e *StringToPtr) Uses() []types.Usage {
	return e.X.Uses()
}

var _ PtrExpr = (*NewExpr)(nil)

type NewExpr struct {
	e    *types.Env
	Elem types.Type
}

func (e *NewExpr) Visit(_ Visitor) {
}

func (e *NewExpr) CType(_ types.Type) types.Type {
	return e.e.PtrT(e.Elem)
}

func (e *NewExpr) AsExpr() GoExpr {
	t := e.Elem.GoType()
	return call(ident("new"), t)
}

func (e *NewExpr) IsConst() bool {
	return false
}

func (e *NewExpr) HasSideEffects() bool {
	return true
}

func (e *NewExpr) Uses() []types.Usage {
	// TODO: use type
	return nil
}

func (e *NewExpr) PtrType(_ types.PtrType) types.PtrType {
	return e.e.PtrT(e.Elem)
}

var _ Expr = (*MakeExpr)(nil)

type MakeExpr struct {
	e    *types.Env
	Elem types.Type
	Size Expr
}

func (e *MakeExpr) Visit(v Visitor) {
	v(e.Size)
}

func (e *MakeExpr) CType(_ types.Type) types.Type {
	return types.SliceT(e.Elem)
}

func (e *MakeExpr) AsExpr() GoExpr {
	tp := e.CType(nil).GoType()
	return call(ident("make"), tp, e.Size.AsExpr())
}

func (e *MakeExpr) IsConst() bool {
	return false
}

func (e *MakeExpr) HasSideEffects() bool {
	return true
}

func (e *MakeExpr) Uses() []types.Usage {
	// TODO: use type
	return types.UseRead(e.Size)
}
