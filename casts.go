package cxgo

import (
	"fmt"
	"math"

	"github.com/gotranspile/cxgo/types"
)

func (g *translator) cCast(toType types.Type, x Expr) Expr {
	if toType == nil {
		panic("no type")
	}
	toKind := toType.Kind()
	xType := x.CType(toType)
	xKind := xType.Kind()
	if xType == g.env.Go().Any() {
		return &CCastExpr{Assert: true, Type: toType, Expr: x}
	}
	if toType == g.env.Go().Any() {
		return x
	}
	if at, ok := toType.(types.ArrayType); ok && at.IsSlice() {
		switch x := x.(type) {
		case Nil:
			return g.Nil()
		case IntLit:
			if x.IsZero() {
				return g.Nil()
			}
		case *TakeAddr:
			if ind, ok := x.X.(*CIndexExpr); ok && types.Same(ind.Expr.CType(nil), toType) {
				if ind.IndexZero() {
					// special case: unwrap unnecessary cast to slice
					return ind.Expr
				}
			}
		}
		if fc, ok := x.(*CallExpr); ok && len(fc.Args) >= 1 {
			if f, ok := fc.Fun.(Ident); ok {
				gg := g.env.Go()
				switch f.Identifier() {
				case gg.SliceFunc(),
					gg.AppendFunc():
					if types.Same(toType, fc.Args[0].CType(nil)) {
						return x
					}
				}
			}
		}
	}
	if xKind.Is(types.Array) && !toKind.Is(types.Array) {
		x = g.cAddr(x)
		return g.cCast(toType, x)
	}
	// equal or same type: no conversion
	if types.Same(toType, xType) {
		return x
	}
	// unknown types: bypass
	if toKind.Is(types.Unknown) {
		// special cases for well-known types
		switch toType {
		case g.env.Go().String():
			var fnc *types.Ident
			if types.Same(xType, g.env.C().String()) {
				fnc = g.env.StringC2Go()
			} else if types.Same(xType, g.env.C().WString()) {
				fnc = g.env.WStringC2Go()
			} else {
				return g.cCast(g.env.C().String(), x)
			}
			return g.NewCCallExpr(FuncIdent{fnc}, []Expr{x})
		}
		return x
	}
	if c1, ok := cUnwrap(x).(*CCastExpr); ok {
		// casts A(A(x)) -> A(x)
		if types.Same(c1.Type, toType) {
			return c1
		}
	}
	// conversions to bool - we have a specialized function for that
	if toKind.IsBool() {
		return g.ToBool(x)
	}
	// nil should be first, because it's an "untyped ptr"
	if xKind.Is(types.Nil) {
		if toKind.IsPtr() || toKind.IsFunc() {
			return cUnwrap(x)
		}
	}
	// strings are immutable, so call a specialized function for conversion
	if types.Same(xType, g.env.Go().String()) {
		// string -> []byte
		if at, ok := types.Unwrap(toType).(types.ArrayType); ok && at.IsSlice() && at.Elem() == g.env.Go().Byte() {
			return &CCastExpr{Type: at, Expr: x}
		}
		// [N]byte = "xyz"
		if at, ok := types.Unwrap(toType).(types.ArrayType); ok && (types.Same(at.Elem(), g.env.Go().Byte()) || xKind == types.Unknown) {
			if !at.IsSlice() {
				tmp := types.NewIdent("t", at)
				copyF := FuncIdent{g.env.Go().CopyFunc()}
				// declare a function literal returning an array of the same size
				var body []CStmt
				// declare temp variable with an array type
				body = append(body, g.NewCDeclStmt(&CVarDecl{CVarSpec: CVarSpec{
					g:     g,
					Type:  at,
					Names: []*types.Ident{tmp},
				}})...)
				body = append(body,
					// copy string into it
					&CExprStmt{Expr: &CallExpr{
						Fun: copyF, Args: []Expr{
							&SliceExpr{Expr: IdentExpr{tmp}},                   // tmp[:]
							&CCastExpr{Type: types.SliceT(at.Elem()), Expr: x}, // ([]TYPE)("xyz")
						},
					}},
				)
				// return temp variable
				body = append(body, g.NewReturnStmt(IdentExpr{tmp}, at)...)
				lit := g.NewFuncLit(g.env.FuncTT(at), body...)
				return g.NewCCallExpr(lit, nil)
			}
		}
		if types.Same(toType, g.env.C().WString()) {
			return g.cCast(toType, g.NewCCallExpr(FuncIdent{g.env.WStringGo2C()}, []Expr{x}))
		}
		return g.cCast(toType, g.NewCCallExpr(FuncIdent{g.env.StringGo2C()}, []Expr{x}))
	}
	if xType == g.env.Go().String() {
		var conv *types.Ident
		if types.Same(toType, g.env.C().String()) {
			conv = g.env.StringGo2C()
		} else if types.Same(toType, g.env.C().WString()) {
			conv = g.env.WStringGo2C()
		}
		if conv != nil {
			return g.cCast(toType, g.NewCCallExpr(FuncIdent{conv}, []Expr{x}))
		}
	}
	// any casts from array to other types should go through pointer to an array
	if xKind.Is(types.Unknown) {
		return &CCastExpr{
			Type: toType,
			Expr: x,
		}
	}
	switch {
	case toKind.IsPtr():
		return g.cPtrToPtr(toType, g.ToPointer(x))
	case toKind.IsInt():
		if l, ok := cUnwrap(x).(IntLit); ok {
			ti, ok := types.Unwrap(toType).(types.IntType)
			if l.IsUint() && ok && ti.Signed() && !litCanStore(ti, l) {
				// try overflowing it
				return l.OverflowInt(ti.Sizeof())
			}
			if l.IsUint() || (ok && ti.Signed()) {
				return &CCastExpr{
					Type: toType,
					Expr: x,
				}
			}
			sz := toType.Sizeof()
			var uv uint64
			switch sz {
			case 1:
				uv = math.MaxUint8
			case 2:
				uv = math.MaxUint16
			case 4:
				uv = math.MaxUint32
			case 8:
				uv = math.MaxUint64
			default:
				return &CCastExpr{
					Type: toType,
					Expr: x,
				}
			}
			uv -= uint64(-l.Int()) - 1
			return cUintLit(uv, l.base)
		}
		if xKind.IsFunc() {
			// func() -> int
			return &FuncToInt{
				X:  g.ToFunc(x, nil),
				To: types.Unwrap(toType).(types.IntType),
			}
		}
		if xKind.IsPtr() {
			// *some -> int
			return g.cPtrToInt(toType, g.ToPointer(x))
		}
		if x.IsConst() && xKind.IsUntypedInt() {
			return x
		}
		if xKind.IsBool() {
			return g.cCast(toType, &BoolToInt{X: g.ToBool(x)})
		}
		xi, ok1 := types.Unwrap(xType).(types.IntType)
		ti, ok2 := types.Unwrap(toType).(types.IntType)
		if ok1 && ok2 && xi.Signed() != ti.Signed() {
			if ti.Sizeof() > xi.Sizeof() {
				return &CCastExpr{
					Type: toType,
					Expr: x,
				}
			} else if ti.Sizeof() < xi.Sizeof() {
				var t2 types.Type
				if !ti.Signed() {
					t2 = types.IntT(ti.Sizeof())
				} else {
					t2 = types.UintT(ti.Sizeof())
				}
				return &CCastExpr{
					Type: toType,
					Expr: g.cCast(t2, x),
				}
			}
		}
		return &CCastExpr{
			Type: toType,
			Expr: x,
		}
	case toKind.IsFunc():
		switch x := cUnwrap(x).(type) {
		case Nil:
			return x
		case IntLit:
			if x.IsZero() {
				return g.Nil()
			}
		}
		if !xKind.IsFunc() {
			x = g.ToFunc(x, types.Unwrap(toType).(*types.FuncType))
			return g.cCast(toType, x)
		}
		ft, fx := types.Unwrap(toType).(*types.FuncType), types.Unwrap(xType).(*types.FuncType)
		if (ft.Variadic() == fx.Variadic() || !ft.Variadic()) && ft.ArgN() >= fx.ArgN() && ((ft.Return() != nil) == (fx.Return() != nil) || (ft.Return() == nil && fx.Return() != nil)) {
			// cannot cast directly, but can return lambda instead
			callArgs := make([]Expr, 0, ft.ArgN())
			funcArgs := make([]*types.Field, 0, ft.ArgN())
			for i, a := range ft.Args() {
				at := a.Type()
				name := types.NewIdent(fmt.Sprintf("arg%d", i+1), at)
				if a.Name != nil && !a.Name.IsUnnamed() {
					name = types.NewIdent(a.Name.Name, at)
				}
				funcArgs = append(funcArgs, &types.Field{
					Name: name,
				})
				callArgs = append(callArgs, IdentExpr{name})
			}
			argn := fx.ArgN()
			if ft.Variadic() {
				callArgs = append(callArgs, &ExpandExpr{X: IdentExpr{types.NewIdent("_rest", types.UnkT(1))}})
				argn++
			}
			e := g.NewCCallExpr(g.ToFunc(x, nil), callArgs[:argn])
			var stmts []CStmt
			if ft.Return() != nil {
				stmts = g.NewReturnStmt(e, ft.Return())
			} else {
				stmts = NewCExprStmt(e)
			}
			var litT *types.FuncType
			if ft.Variadic() {
				litT = g.env.VarFuncT(
					ft.Return(),
					funcArgs...,
				)
			} else {
				litT = g.env.FuncT(
					ft.Return(),
					funcArgs...,
				)
			}
			return g.NewFuncLit(litT, stmts...)
		}
		// incompatible function types - force error
		return x
	case toKind.IsFloat():
		if xKind.IsUntypedFloat() {
			return x
		}
		if xKind.IsUntypedInt() {
			if !toKind.IsUntypedFloat() {
				toType = types.AsUntypedFloatT(types.Unwrap(toType).(types.FloatType))
				toKind = toType.Kind()
			}
		}
		if xKind.IsBool() {
			return g.cCast(toType, &BoolToInt{X: g.ToBool(x)})
		}
	case toKind.Is(types.Array):
		ta := types.Unwrap(toType).(types.ArrayType)
		if xa, ok := types.Unwrap(xType).(types.ArrayType); ok && ta.Len() == 0 && xa.Len() != 0 {
			return &SliceExpr{Expr: x}
		}
	case toKind.Is(types.Struct):
		isZero := false
		switch x := cUnwrap(x).(type) {
		case Nil:
			isZero = true
		case IntLit:
			if x.IsZero() {
				isZero = true
			}
		}
		if isZero {
			return g.NewCCompLitExpr(toType, nil)
		}
	}
	return &CCastExpr{
		Type: toType,
		Expr: x,
	}
}
