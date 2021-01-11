package cxgo

import (
	"fmt"

	"github.com/gotranspile/cxgo/types"
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
	var decls []CDecl
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
		sname := name.String()
		if op, err := ast.Eval(mc); err == nil {
			val := g.convertValue(op.Value())
			typ := val.CType(nil)
			id := types.NewIdent(sname, typ)
			decls = append(decls, &CVarDecl{Const: true, CVarSpec: CVarSpec{
				g: g, Type: typ,
				Names: []*types.Ident{id},
				Inits: []Expr{val},
			}})
		}
	}
	return decls
}
