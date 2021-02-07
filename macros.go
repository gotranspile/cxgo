package cxgo

import (
	"fmt"
	"github.com/gotranspile/cxgo/types"
	"sort"

	"modernc.org/cc/v3"
)

func (g *translator) convertValue(v cc.Value) Expr {
	switch v := v.(type) {
	case cc.Int64Value:
		return cIntLit(int64(v))
	case cc.Uint64Value:
		return cUintLit(uint64(v))
	case cc.Float32Value:
		return FloatLit{val: float64(v)}
	case cc.Float64Value:
		return FloatLit{val: float64(v)}
	case cc.StringValue:
		e, err := g.parseCStringLit(string(v))
		if err != nil {
			panic(err)
		}
		return e
	default:
		panic(fmt.Errorf("unsupported value type: %T", v))
	}
}

func (g *translator) convertMacros(ast *cc.AST) []CDecl {
	type macro struct {
		name string
		m    *cc.Macro
	}
	var arr []macro
	for name, mc := range ast.Macros {
		if !g.inCurFile(mc) {
			continue
		}
		if mc.IsFnLike() {
			continue // we don't support function macros yet
		}
		if len(mc.ReplacementTokens()) == 0 {
			continue // no value
		}
		arr = append(arr, macro{name.String(), mc})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].m.Position().Offset < arr[j].m.Position().Offset
	})
	var decls []CDecl
	for _, mc := range arr {
		if op, err := ast.Eval(mc.m); err == nil {
			val := g.convertValue(op.Value())
			typ := val.CType(nil)
			id := types.NewIdent(mc.name, typ)
			decls = append(decls, &CVarDecl{Const: true, CVarSpec: CVarSpec{
				g: g, Type: typ,
				Names: []*types.Ident{id},
				Inits: []Expr{val},
			}})
		}
	}
	return decls
}
