package cxgo

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type GoExpr = ast.Expr
type GoStmt = ast.Stmt
type GoDecl = ast.Decl
type GoType = ast.Expr

type GoField = ast.Field

func ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func boolLit(v bool) GoExpr {
	if v {
		return ident("true")
	}
	return ident("false")
}

func intLit(v int) GoExpr {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(v),
	}
}

func formatInt(v int64, base int) string {
	neg := v < 0
	if neg {
		base = 10
	}
	if neg {
		v = -v
	}
	s := formatUint(uint64(v), base)
	if neg {
		s = "-" + s
	}
	return s
}

func formatUint(v uint64, base int) string {
	s := strconv.FormatUint(v, base)
	switch base {
	case 2:
		s = "0b" + s
	case 8:
		s = "0o" + s
	case 16:
		s = "0x" + strings.ToUpper(s)
	}
	return s
}

func intLit64(v int64, base int) GoExpr {
	if base <= 0 {
		base = 10
	}
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: formatInt(v, base),
	}
}

func uintLit64(v uint64, base int) GoExpr {
	if base > 0 {
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: formatUint(v, base),
		}
	}
	s10 := strconv.FormatUint(v, 10)
	s := s10
	if len(s) > 4 {
		s16 := strconv.FormatUint(v, 16)
		if len(strings.Trim(s10, "0")) > len(strings.Trim(s16, "0")) {
			s = "0x" + strings.ToUpper(s16)
		}
	}
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: s,
	}
}

func noFields() *ast.FieldList {
	return &ast.FieldList{}
}

func fields(arr []*GoField) *ast.FieldList {
	return &ast.FieldList{List: arr}
}

func fieldTypes(arr ...GoType) *ast.FieldList {
	var out []*GoField
	for _, t := range arr {
		out = append(out, &ast.Field{Type: t})
	}
	return fields(out)
}

func funcTypeRet(ret GoType) *ast.FuncType {
	var fret *ast.FieldList
	if ret != nil {
		fret = fieldTypes(ret)
	}
	return &ast.FuncType{
		Params:  noFields(),
		Results: fret,
	}
}

func fixLabels(stmts []GoStmt) []GoStmt {
	for i := 0; i < len(stmts); i++ {
		l, ok := stmts[i].(*ast.LabeledStmt)
		if !ok {
			continue
		}
		switch l.Stmt.(type) {
		case nil, *ast.EmptyStmt:
			if i == len(stmts)-1 {
				l.Stmt = &ast.EmptyStmt{Implicit: true}
			} else {
				l.Stmt = stmts[i+1]
				switch l.Stmt.(type) {
				case *ast.CaseClause:
					l.Stmt = &ast.EmptyStmt{Implicit: true}
					continue
				}
				stmts = append(stmts[:i+1], stmts[i+2:]...)
			}
		}
	}
	return stmts
}

func mergeDecls(stmts []GoStmt) []GoStmt {
	if len(stmts) < 2 {
		return stmts
	}
	d1, ok := stmts[0].(*ast.DeclStmt)
	if !ok {
		return stmts
	}
	g1, ok := d1.Decl.(*ast.GenDecl)
	if !ok {
		return stmts
	}
	j := 0
	for i := 0; i < len(stmts); i++ {
		if i == 0 {
			continue
		}
		d2, ok := stmts[i].(*ast.DeclStmt)
		if !ok {
			break
		}
		g2, ok := d2.Decl.(*ast.GenDecl)
		if !ok || g2.Tok != g1.Tok {
			break
		}
		g1.Specs = append(g1.Specs, g2.Specs...)
		g2.Specs = nil
		j = i
	}
	if j == 0 {
		return stmts
	}
	return append([]GoStmt{d1}, stmts[j+1:]...)
}

func block(stmts ...GoStmt) *ast.BlockStmt {
	if len(stmts) == 1 {
		if b, ok := stmts[0].(*ast.BlockStmt); ok {
			return b
		}
	}
	stmts = fixLabels(stmts)
	stmts = mergeDecls(stmts)
	return &ast.BlockStmt{List: stmts}
}

func ifelse(cond GoExpr, then, els []GoStmt) *ast.IfStmt {
	s := &ast.IfStmt{
		Cond: cond,
		Body: block(then...),
	}
	if len(els) != 0 {
		if len(els) == 1 {
			if s2, ok := els[0].(*ast.IfStmt); ok {
				s.Else = s2
			} else {
				s.Else = block(els...)
			}
		} else {
			s.Else = block(els...)
		}
	}
	return s
}

func returnStmt(expr ...GoExpr) *ast.ReturnStmt {
	return &ast.ReturnStmt{Results: expr}
}

func call(fnc GoExpr, args ...GoExpr) *ast.CallExpr {
	for i, a := range args {
		if p, ok := a.(*ast.ParenExpr); ok {
			args[i] = p.X
		}
	}
	return &ast.CallExpr{
		Fun:  fnc,
		Args: args,
	}
}

func callVari(fnc GoExpr, args ...GoExpr) *ast.CallExpr {
	e := call(fnc, args...)
	e.Ellipsis = 1
	return e
}

func callLambda(ret GoType, stmts ...GoStmt) *ast.CallExpr {
	return call(&ast.FuncLit{
		Type: funcTypeRet(ret),
		Body: block(stmts...),
	})
}

func assignTok(x GoExpr, tok token.Token, y GoExpr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: tok,
		Lhs: []GoExpr{x},
		Rhs: []GoExpr{y},
	}
}

func assign(x, y GoExpr) *ast.AssignStmt {
	return assignTok(x, token.ASSIGN, y)
}

func define(x, y GoExpr) *ast.AssignStmt {
	return assignTok(x, token.DEFINE, y)
}

func addr(x GoExpr) GoExpr {
	if e1, ok := x.(*ast.StarExpr); ok {
		return e1.X
	}
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  x,
	}
}

func deref(x GoExpr) *ast.StarExpr {
	return &ast.StarExpr{
		X: x,
	}
}

func index(x, ind GoExpr) *ast.IndexExpr {
	return &ast.IndexExpr{
		X: x, Index: ind,
	}
}

func exprStmt(x GoExpr) GoStmt {
	return &ast.ExprStmt{
		X: x,
	}
}

func paren(x GoExpr) *ast.ParenExpr {
	if p, ok := x.(*ast.ParenExpr); ok {
		return p
	}
	return &ast.ParenExpr{
		X: x,
	}
}

func unsafePtr() GoType {
	return ident("unsafe.Pointer")
}

func typAssert(x GoExpr, t GoType) GoExpr {
	return &ast.TypeAssertExpr{
		X:    x,
		Type: t,
	}
}
