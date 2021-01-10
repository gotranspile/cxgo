package cxgo

import (
	"fmt"
	"go/ast"

	"github.com/gotranspile/cxgo/types"
)

func ToFuncExpr(exp types.Type) *types.FuncType {
	if t, ok := types.Unwrap(exp).(*types.FuncType); ok {
		return t
	}
	return nil
}

var _ FuncExpr = Nil{}

// FuncExpr is an expression that returns a function value.
type FuncExpr interface {
	Expr
	// FuncType return an underlying function type.
	FuncType(exp *types.FuncType) *types.FuncType
}

var (
	_ FuncExpr = FuncIdent{}
	_ Ident    = FuncIdent{}
)

type FuncIdent struct {
	*types.Ident
}

func (FuncIdent) Visit(v Visitor) {}

func (e FuncIdent) Identifier() *types.Ident {
	return e.Ident
}

func (e FuncIdent) IsConst() bool {
	return false
}

func (e FuncIdent) HasSideEffects() bool {
	return false
}

func (e FuncIdent) AsExpr() GoExpr {
	return e.GoIdent()
}

func (e FuncIdent) FuncType(exp *types.FuncType) *types.FuncType {
	var et types.Type
	if exp != nil {
		et = exp
	}
	return types.Unwrap(e.CType(et)).(*types.FuncType)
}

func (e FuncIdent) Uses() []types.Usage {
	return []types.Usage{{Ident: e.Ident, Access: types.AccessUnknown}}
}

var _ FuncExpr = FuncAssert{}

type FuncAssert struct {
	X Expr
}

func (e FuncAssert) Visit(v Visitor) {
	v(e.X)
}

func (e FuncAssert) CType(types.Type) types.Type {
	return e.X.CType(nil)
}

func (e FuncAssert) IsConst() bool {
	return e.X.IsConst()
}

func (e FuncAssert) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e FuncAssert) AsExpr() GoExpr {
	return e.X.AsExpr()
}

func (e FuncAssert) FuncType(*types.FuncType) *types.FuncType {
	return types.Unwrap(e.CType(nil)).(*types.FuncType)
}

func (e FuncAssert) Uses() []types.Usage {
	return e.X.Uses()
}

func asSizeofT(e Expr) types.Type {
	e = unwrapCasts(e)
	sz, ok := e.(*CSizeofExpr)
	if !ok {
		return nil
	}
	return sz.Type
}

func (g *translator) NewCCallExpr(fnc FuncExpr, args []Expr) Expr {
	if id, ok := cUnwrap(fnc).(IdentExpr); ok {
		// TODO: another way to hook into it?
		switch id.Ident {
		case g.env.C().MallocFunc():
			// malloc(sizeof(T)) -> new(T)
			if len(args) == 1 {
				if tp := asSizeofT(args[0]); tp != nil {
					return &NewExpr{
						e:    g.env.Env,
						Elem: tp,
					}
				}
			}
		case g.env.C().CallocFunc():
			// calloc(n, sizeof(T)) -> make([]T, n)
			if len(args) == 2 {
				if tp := asSizeofT(args[1]); tp != nil {
					return &MakeExpr{
						e:    g.env.Env,
						Elem: tp,
						Size: g.cCast(g.env.Go().Int(), args[0]),
					}
				}
			}
		case g.env.C().MemsetFunc():
			// memset(p, 0, sizeof(T)) -> *p = T{}
			if len(args) == 3 {
				if lit, ok := unwrapCasts(args[1]).(IntLit); ok && lit.IsZero() {
					if tp := asSizeofT(args[2]); tp != nil {
						p := g.cDeref(g.ToPointer(g.cCast(g.env.PtrT(tp), args[0])))
						return g.NewCAssignExpr(p, "", g.ZeroValue(tp))
					}
				}
			}
		}
		//	kf, ok := knownCFuncs[id.Name]
		//	if !ok {
		//		kf, ok = known[id.Ident]
		//	}
		//	if ok && (id.Ident != kf.Name || (kf.This != nil && len(args) > len(kf.This.Args))) {
		//		if kf.This == nil {
		//			return NewCCallExpr(IdentExpr{kf.Name}, args)
		//		}
		//		fnc = &CSelectExpr{
		//			ctype: kf.This,
		//			Expr: cCast(kf.Type.Args[0].Type, args[0]),
		//			Sel: kf.Name,
		//		}
		//		return NewCCallExpr(fnc, args[1:])
		//	}
	}
	t := fnc.CType(nil)
	if p, ok := t.(types.PtrType); ok && p.Elem() != nil {
		t = p.Elem()
	}
	ft := types.Unwrap(t).(*types.FuncType)
	ftargs := ft.Args()
	for i, a := range args {
		var atyp types.Type
		if i < len(ftargs) {
			atyp = ftargs[i].Type()
		} else if ft.Variadic() {
			atyp = types.UnkT(1)
		} else {
			break
		}
		args[i] = g.cCast(atyp, a)
	}
	return &CallExpr{
		Fun:  fnc,
		Args: args,
	}
}

var _ Expr = (*CallExpr)(nil)

type CallExpr struct {
	Fun  FuncExpr
	Args []Expr
}

func (e *CallExpr) Visit(v Visitor) {
	v(e.Fun)
	for _, a := range e.Args {
		v(a)
	}
}

func (e *CallExpr) CType(types.Type) types.Type {
	typ := e.Fun.FuncType(nil)
	if typ.Return() == nil {
		panic("function doesn't return")
	}
	return typ.Return()
}

func (e *CallExpr) IsConst() bool {
	return false
}

func (e *CallExpr) HasSideEffects() bool {
	return true
}

func (e *CallExpr) AsExpr() GoExpr {
	var args []GoExpr
	vari := false
	for _, a := range e.Args {
		if _, ok := a.(*ExpandExpr); ok {
			vari = true
			args = append(args, ident("_rest"))
			continue
		}
		args = append(args, a.AsExpr())
	}
	if vari {
		return callVari(e.Fun.AsExpr(), args...)
	}
	return call(e.Fun.AsExpr(), args...)
}

func (e *CallExpr) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(e.Fun)...)
	for _, a := range e.Args {
		list = append(list, types.UseRead(a)...)
	}
	return list
}

var _ FuncExpr = (*FuncLit)(nil)

func (g *translator) NewFuncLit(typ *types.FuncType, body ...CStmt) *FuncLit {
	if typ.Return() != nil && typ.Return().Kind().IsUntyped() {
		panic("untyped")
	}
	return &FuncLit{
		Type: typ,
		Body: g.NewCBlock(body...),
	}
}

type FuncLit struct {
	Type *types.FuncType
	Body *BlockStmt
}

func (e *FuncLit) Visit(v Visitor) {
	v(e.Body)
}

func (e *FuncLit) IsConst() bool {
	return false
}

func (e *FuncLit) HasSideEffects() bool {
	return true
}

func (e *FuncLit) CType(types.Type) types.Type {
	return e.Type
}

func (e *FuncLit) FuncType(*types.FuncType) *types.FuncType {
	return e.Type
}

func (e *FuncLit) AsExpr() GoExpr {
	return &ast.FuncLit{
		Type: e.Type.GoFuncType(),
		Body: e.Body.GoBlockStmt(),
	}
}

func (e *FuncLit) Uses() []types.Usage {
	var list []types.Usage
	// TODO: use the type
	if e.Body != nil {
		list = append(list, e.Body.Uses()...)
	}
	return list
}

var _ FuncExpr = (*PtrToFunc)(nil)

type PtrToFunc struct {
	X  PtrExpr
	To *types.FuncType
}

func (e *PtrToFunc) Visit(v Visitor) {
	v(e.X)
}

func (e *PtrToFunc) CType(types.Type) types.Type {
	return e.To
}

func (e *PtrToFunc) AsExpr() GoExpr {
	x := e.X.AsExpr()
	ft := e.To.GoType()
	asFunc := call(
		ident("libc.AsFunc"), x,
		call(paren(deref(ft)), ident("nil")),
	)
	return typAssert(asFunc, ft)
}

func (e *PtrToFunc) IsConst() bool {
	return false
}

func (e *PtrToFunc) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *PtrToFunc) Uses() []types.Usage {
	return e.X.Uses()
}

func (e *PtrToFunc) FuncType(*types.FuncType) *types.FuncType {
	return e.To
}

var _ PtrExpr = (*FuncToPtr)(nil)

type FuncToPtr struct {
	e *types.Env
	X FuncExpr
}

func (e *FuncToPtr) Visit(v Visitor) {
	v(e.X)
}

func (e *FuncToPtr) CType(types.Type) types.Type {
	return e.PtrType(nil)
}

func (e *FuncToPtr) AsExpr() GoExpr {
	x := e.X.AsExpr()
	return call(unsafePtr(), call(ident("libc.FuncAddr"), x))
}

func (e *FuncToPtr) IsConst() bool {
	return false
}

func (e *FuncToPtr) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *FuncToPtr) Uses() []types.Usage {
	return e.X.Uses()
}

func (e *FuncToPtr) PtrType(types.PtrType) types.PtrType {
	return e.e.PtrT(nil)
}

var _ FuncExpr = (*IntToFunc)(nil)

type IntToFunc struct {
	X  Expr
	To *types.FuncType
}

func (e *IntToFunc) Visit(v Visitor) {
	v(e.X)
}

func (e *IntToFunc) CType(types.Type) types.Type {
	return e.To
}

func (e *IntToFunc) AsExpr() GoExpr {
	x := e.X.AsExpr()
	ft := e.To.GoType()
	asFunc := call(
		ident("libc.AsFunc"), x,
		call(paren(deref(ft)), ident("nil")),
	)
	return typAssert(asFunc, ft)
}

func (e *IntToFunc) IsConst() bool {
	return false
}

func (e *IntToFunc) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *IntToFunc) Uses() []types.Usage {
	return e.X.Uses()
}

func (e *IntToFunc) FuncType(*types.FuncType) *types.FuncType {
	return e.To
}

var _ Expr = (*FuncToInt)(nil)

type FuncToInt struct {
	X  FuncExpr
	To types.IntType
}

func (e *FuncToInt) Visit(v Visitor) {
	v(e.X)
}

func (e *FuncToInt) CType(types.Type) types.Type {
	return e.To
}

func (e *FuncToInt) AsExpr() GoExpr {
	x := e.X.AsExpr()
	x = call(ident("libc.FuncAddr"), x)
	x = call(e.To.GoType(), x)
	return x
}

func (e *FuncToInt) IsConst() bool {
	return false
}

func (e *FuncToInt) HasSideEffects() bool {
	return e.X.HasSideEffects()
}

func (e *FuncToInt) Uses() []types.Usage {
	return e.X.Uses()
}

func (g *translator) ToFunc(x Expr, exp *types.FuncType) FuncExpr {
	if x, ok := x.(FuncExpr); ok {
		return x
	}
	if x, ok := cUnwrap(x).(IntLit); ok && x.IsZero() {
		return g.Nil()
	}
	xt := x.CType(nil)
	xk := xt.Kind()
	if xk.IsFunc() {
		switch x := x.(type) {
		case Ident:
			return FuncIdent{x.Identifier()}
		}
		return FuncAssert{x}
	}
	if xk.IsPtr() {
		if exp == nil {
			panic("expected type must be set")
		}
		xp := g.ToPointer(x)
		switch xp := xp.(type) {
		case *FuncToPtr:
			return xp.X
		}
		return &PtrToFunc{
			X:  xp,
			To: exp,
		}
	}
	if xk.IsInt() {
		if exp == nil {
			panic("expected type must be set")
		}
		return &IntToFunc{
			X:  x,
			To: exp,
		}
	}
	panic(fmt.Errorf("not a function: %T, %T", xt, x))
}

func CompareFuncs(x FuncExpr, op ComparisonOp, y FuncExpr) BoolExpr {
	if op.IsRelational() {
		return &FuncComparison{
			X: x, Op: op, Y: y,
		}
	}
	if !op.IsEquality() {
		panic("must not happen")
	}
	// == and != simplifications
	// always compare with the constant on the right
	if x.IsConst() && !y.IsConst() {
		return CompareFuncs(y, op, x)
	}
	return &FuncComparison{
		X: x, Op: op, Y: y,
	}
}

var _ BoolExpr = (*FuncComparison)(nil)

type FuncComparison struct {
	X  FuncExpr
	Op ComparisonOp
	Y  FuncExpr
}

func (e *FuncComparison) Visit(v Visitor) {
	v(e.X)
	v(e.Y)
}

func (e *FuncComparison) CType(types.Type) types.Type {
	return types.BoolT()
}

func (e *FuncComparison) AsExpr() GoExpr {
	x := e.X.AsExpr()
	y := e.Y.AsExpr()
	if _, ok := e.Y.(Nil); !ok {
		// compare via uintptr
		x = call(ident("libc.FuncAddr"), x)
		y = call(ident("libc.FuncAddr"), y)
	}
	return &ast.BinaryExpr{
		X:  x,
		Op: e.Op.GoToken(),
		Y:  y,
	}
}

func (e *FuncComparison) IsConst() bool {
	return e.X.IsConst() && e.Y.IsConst()
}

func (e *FuncComparison) HasSideEffects() bool {
	return e.X.HasSideEffects() || e.Y.HasSideEffects()
}

func (e *FuncComparison) Negate() BoolExpr {
	return CompareFuncs(e.X, e.Op.Negate(), e.Y)
}

func (e *FuncComparison) Uses() []types.Usage {
	return types.UseRead(e.X, e.Y)
}

type ExpandExpr struct {
	// TODO: support properly
}

func (e *ExpandExpr) Visit(v Visitor) {
}

func (e *ExpandExpr) CType(_ types.Type) types.Type {
	return types.UnkT(1)
}

func (e *ExpandExpr) AsExpr() GoExpr {
	panic("should be handled by CallExpr")
}

func (e *ExpandExpr) IsConst() bool {
	return false
}

func (e *ExpandExpr) HasSideEffects() bool {
	return false
}

func (e *ExpandExpr) Uses() []types.Usage {
	return nil
}
