package cxgo

import (
	"go/ast"
	"go/format"
	token2 "go/token"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/gotranspile/cxgo/libs"
)

func PrintGo(w io.Writer, pkg string, decls []GoDecl) error {
	return format.Node(w, token2.NewFileSet(), &ast.File{Decls: decls, Name: ident(pkg)})
}

type usageVisitor struct {
	used map[string]struct{}
}

func (v *usageVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Ident:
		sub := strings.SplitN(n.Name, ".", 2)
		if len(sub) == 2 {
			v.used[sub[0]] = struct{}{}
		}
	}
	return v
}

func goUsedImports(used map[string]struct{}, decls []GoDecl) {
	for _, d := range decls {
		ast.Walk(&usageVisitor{used: used}, d)
	}
}

// ImportsFor generates import specs for well-known imports required for given declarations.
func ImportsFor(e *libs.Env, decls []GoDecl) []GoDecl {
	used := make(map[string]struct{})
	goUsedImports(used, decls)
	var list []string
	for k := range used {
		list = append(list, k)
	}
	sort.Strings(list)
	var specs []ast.Spec
	for _, name := range list {
		path := e.ResolveImport(name)
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{
			Kind:  token2.STRING,
			Value: strconv.Quote(path),
		}})
	}
	if len(specs) == 0 {
		return nil
	}
	return []GoDecl{
		&ast.GenDecl{
			Tok:   token2.IMPORT,
			Specs: specs,
		},
	}
}
