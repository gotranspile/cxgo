package cxgo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gotranspile/cxgo/types"

	"modernc.org/cc/v3"
)

func (g *translator) convertValue(v cc.Value) Expr {
	switch v := v.(type) {
	case cc.Int64Value:
		return cIntLit(int64(v), 0)
	case cc.Uint64Value:
		return cUintLit(uint64(v), 0)
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
		if val := g.evalMacro(mc.m, ast); val != nil {
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

func (g *translator) evalMacro(m *cc.Macro, ast *cc.AST) Expr {
	toks := m.ReplacementTokens()
	if len(toks) != 1 {
		return evalMacro2(m, ast)
	}

	src := strings.TrimSpace(toks[0].Src.String())
	if len(src) == 0 {
		return nil
	}

	if src[0] == '"' {
		if s, err := strconv.Unquote(src); err == nil {
			if l, err := g.parseCStringLit(s); err == nil {
				return l
			}
		}
	} else {
		if l, err := parseCIntLit(src, g.conf.IntReformat); err == nil {
			return l
		}
		if l, err := parseCFloatLit(src); err == nil {
			return l
		}
	}

	return evalMacro2(m, ast)
}

func evalMacro2(m *cc.Macro, ast *cc.AST) Expr {
	op, err := ast.Eval(m)
	if err != nil {
		return nil
	}

	switch x := op.Value().(type) {
	case cc.Int64Value:
		return cIntLit(int64(x), 0)
	case cc.Uint64Value:
		return cUintLit(uint64(x), 0)
	default:
		return nil
	}
}
