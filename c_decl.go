package cxgo

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gotranspile/cxgo/types"
)

type Range struct {
	Start     int
	StartLine int
	End       int
}

type CDecl interface {
	Node
	AsDecl() []GoDecl
	Uses() []types.Usage
}

type CVarSpec struct {
	g     *translator
	Type  types.Type
	Names []*types.Ident
	Inits []Expr
}

func (d *CVarSpec) Visit(v Visitor) {
	for _, name := range d.Names {
		v(IdentExpr{name})
	}
	for _, x := range d.Inits {
		v(x)
	}
}

func (d *CVarSpec) GoSpecs(isConst bool) *ast.ValueSpec {
	var (
		names []*ast.Ident
		init  []GoExpr
	)
	var typ GoType
	if d.Type != nil {
		typ = d.Type.GoType()
	}
	dropType := 0
	for i, name := range d.Names {
		if name.Name == "__func__" {
			continue
		}
		names = append(names, name.GoIdent())
		if len(d.Inits) != 0 {
			v := d.Inits[i]
			if d.Type != nil && v != nil {
				v = d.g.cCast(d.Type, v)
			}
			var e GoExpr
			if v != nil {
				e = v.AsExpr()
			}
			if ci, ok := e.(*ast.CompositeLit); ok && ci.Type == nil {
				ci.Type = typ
				dropType++
			} else if _, ok = e.(*ast.BasicLit); ok && isConst {
				if _, ok := d.Type.(types.Named); !ok {
					dropType++
				}
			}
			init = append(init, e)
		}
	}
	if dropType == len(names) {
		typ = nil
	}
	if len(init) != 0 && len(init) < len(names) {
		panic("partial init")
	}
	return &ast.ValueSpec{
		Type:   typ,
		Names:  names,
		Values: init,
	}
}

func (d *CVarSpec) Uses() []types.Usage {
	var list []types.Usage
	for i, name := range d.Names {
		acc := types.AccessDefine
		if i < len(d.Inits) {
			acc = types.AccessWrite
		}
		list = append(list, types.Usage{Ident: name, Access: acc})
	}
	for _, e := range d.Inits {
		list = append(list, types.UseRead(e)...)
	}
	return list
}

type CVarDecl struct {
	Const  bool
	Single bool
	CVarSpec
}

func (d *CVarDecl) Visit(v Visitor) {
	v(&d.CVarSpec)
}

func (d *CVarDecl) GoField() *GoField {
	if len(d.Names) > 1 {
		panic("too large")
	}
	if len(d.Inits) != 0 {
		panic("FIXME")
	}
	var names []*ast.Ident
	if len(d.Names) != 0 {
		names = append(names, d.Names[0].GoIdent())
	}
	return &GoField{
		Names: names,
		Type:  d.Type.GoType(),
		// FIXME: init
	}
}

func (d *CVarDecl) AsDecl() []GoDecl {
	sp := d.GoSpecs(d.Const)
	if sp == nil || len(sp.Names) == 0 {
		return nil
	}
	single := d.Single
	tok := token.VAR
	if d.Const {
		tok = token.CONST
		// no complex consts in Go
		if d.Type != nil && (d.Type.Kind().Is(types.Array) || d.Type.Kind().Is(types.Struct)) {
			tok = token.VAR
		} else {
			for _, v := range d.Inits {
				if v == nil {
					single = false
					continue
				}
				if vk := v.CType(nil).Kind(); vk.Is(types.Array) || vk.Is(types.Struct) {
					tok = token.VAR
					break
				}
			}
		}
	}
	var specs []ast.Spec
	if single {
		specs = []ast.Spec{sp}
	} else {
		for i, name := range sp.Names {
			var vals []ast.Expr
			if len(sp.Values) != 0 && sp.Values[i] != nil {
				vals = []ast.Expr{sp.Values[i]}
			}
			specs = append(specs, &ast.ValueSpec{
				Names:  []*ast.Ident{name},
				Type:   sp.Type,
				Values: vals,
			})
		}
	}
	return []GoDecl{&ast.GenDecl{
		Tok:   tok,
		Specs: specs,
	}}
}

type CTypeDef struct {
	types.Named
}

func (d *CTypeDef) Visit(v Visitor) {
	v(IdentExpr{d.Name()})
}

func (d *CTypeDef) AsDecl() []GoDecl {
	var decls []GoDecl
	if false && d.Kind().Is(types.Struct) {
		decls = append(decls, &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{ident("_")},
					Values: []GoExpr{ident(fmt.Sprintf("([1]struct{}{})[%d-unsafe.Sizeof(%s{})]", d.Sizeof(), d.Name().Name))},
				},
			},
		})
	}
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: d.Name().GoIdent(),
				Type: d.Underlying().GoType(),
			},
		},
	})
	return decls
}

func (d *CTypeDef) Uses() []types.Usage {
	// TODO: use type
	return nil
}

type CFuncDecl struct {
	Name  *types.Ident
	Type  *types.FuncType
	Body  *BlockStmt
	Range *Range
}

func (d *CFuncDecl) Visit(v Visitor) {
	v(IdentExpr{d.Name})
	v(d.Body)
}

func (d *CFuncDecl) AsDecl() []GoDecl {
	return []GoDecl{
		&ast.FuncDecl{
			Name: d.Name.GoIdent(),
			Type: d.Type.GoFuncType(),
			Body: d.Body.GoBlockStmt(),
		},
	}
}

func (d *CFuncDecl) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.Usage{Ident: d.Name, Access: types.AccessDefine})
	if d.Body != nil {
		list = append(list, d.Body.Uses()...)
	}
	// TODO: use type
	return nil
}
