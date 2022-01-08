package types

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type GoType = ast.Expr

func ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func goName(name string) string {
	switch name {
	case "len", "cap",
		"var", "const", "func",
		"type", "chan", "range",
		"copy", "close", "make",
		"string", "byte", "rune",
		"interface", "map":
		name += "_"
	}
	name = strings.ReplaceAll(name, "$", "_")
	if name == "_" {
		name += "1"
	}
	return name
}

func (e *Ident) convertName() {
	if e.GoName != "" {
		return
	}
	e.GoName = goName(e.Name)
}

func (e *Ident) GoIdent() *ast.Ident {
	e.convertName()
	return ident(e.GoName)
}

func (t *unkType) GoType() GoType {
	if t.isStruct {
		return ident("struct{}") // TODO
	}
	return ident("interface{}") // TODO
}

func (t IntType) GoType() GoType {
	u := "u"
	if t.signed {
		u = ""
	}
	return ident(fmt.Sprintf("%sint%d", u, t.Sizeof()*8))
}

func (t FloatType) GoType() GoType {
	return ident(fmt.Sprintf("float%d", t.size*8))
}

func (t BoolType) GoType() GoType {
	return ident("bool")
}

func (t *ptrType) GoType() GoType {
	elem := t.Elem()
	if elem == nil {
		return ident("unsafe.Pointer")
	}
	return &ast.StarExpr{
		X: elem.GoType(),
	}
}

func (t namedPtr) GoType() GoType {
	return t.name.GoIdent()
}

func (t ArrayType) GoType() GoType {
	var sz ast.Expr
	if !t.slice {
		sz = &ast.BasicLit{
			Kind:  token.INT,
			Value: strconv.Itoa(t.size),
		}
	}
	return &ast.ArrayType{
		Len: sz,
		Elt: t.elem.GoType(),
	}
}

func (t *namedType) GoType() GoType {
	return t.name.GoIdent()
}

func (f *Field) GoField() *ast.Field {
	var names []*ast.Ident
	if f.Name != nil && !f.Name.IsUnnamed() {
		names = append(names, f.Name.GoIdent())
	}
	return &ast.Field{Names: names, Type: f.Type().GoType()}
}

func (t *StructType) GoType() GoType {
	fields := &ast.FieldList{}
	if t.union {
		fields.List = append(fields.List, &ast.Field{Type: ident("// union")})
	}
	for _, f := range t.fields {
		fields.List = append(fields.List, f.GoField())
	}
	return &ast.StructType{Fields: fields}
}

func (t *FuncType) GoType() GoType {
	return t.GoFuncType()
}

func (t *FuncType) GoFuncType() *ast.FuncType {
	var ret ast.Expr
	if t.Return() != nil {
		ret = t.Return().GoType()
	}
	p := &ast.FieldList{}
	f := &ast.FuncType{
		Params: p,
	}
	if ret != nil {
		f.Results = &ast.FieldList{List: []*ast.Field{{Type: ret}}}
	}
	hasNames := t.ArgN() == 0 // default to names, if only has variadic
	for _, a := range t.Args() {
		if a.Name.Name != "" {
			hasNames = true
		}
		p.List = append(p.List, a.GoField())
	}
	if t.Variadic() {
		var names []*ast.Ident
		if hasNames {
			names = []*ast.Ident{ident("_rest")}
		}
		p.List = append(p.List, &ast.Field{
			Names: names,
			Type:  &ast.Ellipsis{Elt: ident("interface{}")},
		})
	}
	return f
}
