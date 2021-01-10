package cxgo

import (
	"go/ast"
	"go/token"

	"github.com/dennwc/cxgo/types"
)

const optimizeStatements = false

type CStmt interface {
	Node
	AsStmt() []GoStmt
	Uses() []types.Usage
}

type CStmtFunc func(s CStmt) ([]CStmt, bool)

type CCompStmt interface {
	CStmt
	EachStmt(fnc CStmtFunc) bool
}

type CStmtConv interface {
	ToStmt() []CStmt
}

func cEachBlockStmt(fnc func([]CStmt), stmts []CStmt) {
	fnc(stmts)
	for _, s := range stmts {
		switch s := s.(type) {
		case *BlockStmt:
			cEachBlockStmt(fnc, s.Stmts)
		case *CIfStmt:
			cEachBlockStmt(fnc, s.Then.Stmts)
			if s.Else != nil {
				cEachBlockStmt(fnc, []CStmt{s.Else})
			}
		case *CForStmt:
			cEachBlockStmt(fnc, s.Body.Stmts)
		case *CSwitchStmt:
			for _, c := range s.Cases {
				cEachBlockStmt(fnc, c.Stmts)
			}
		}
	}
}

func cEachStmt(fnc func(CStmt) bool, stmts []CStmt) {
	for _, s := range stmts {
		if !fnc(s) {
			continue
		}
		c, ok := s.(CCompStmt)
		if !ok {
			continue
		}
		c.EachStmt(func(s CStmt) ([]CStmt, bool) {
			cEachStmt(fnc, []CStmt{s})
			return nil, false
		})
	}
}

func cReplaceEachStmt(fnc CStmtFunc, stmts []CStmt) ([]CStmt, bool) {
	b := &BlockStmt{Stmts: stmts}
	var replace CStmtFunc
	replace = func(s CStmt) ([]CStmt, bool) {
		if out, mod := fnc(s); mod {
			return out, mod
		}
		if c, ok := s.(CCompStmt); ok {
			c.EachStmt(replace)
		}
		return nil, false
	}
	if b.EachStmt(replace) {
		return b.Stmts, true
	}
	return stmts, false
}

func cExprStmt(x Expr) []GoStmt {
	if s, ok := x.(CStmt); ok {
		return s.AsStmt()
	}
	switch x := x.(type) {
	case CStmt:
		return x.AsStmt()
	case IntLit:
		return nil
	}
	return []GoStmt{exprStmt(x.AsExpr())}
}

func cOneStmt(x CStmt) GoStmt {
	arr := x.AsStmt()
	if len(arr) != 1 {
		panic("should only one element")
	}
	return arr[0]
}

var (
	_ CStmt = &CExprStmt{}
)

func NewCExprStmt(e Expr) []CStmt {
	e = cUnwrap(e)
	if m, ok := e.(*CMultiExpr); ok {
		out := make([]CStmt, 0, len(m.Exprs))
		for _, e := range m.Exprs {
			out = append(out, NewCExprStmt(e)...)
		}
		return out
	}
	if s, ok := e.(CStmtConv); ok {
		return s.ToStmt()
	}
	return []CStmt{&CExprStmt{Expr: e}}
}

func NewCExprStmt1(e Expr) CStmt {
	if e == nil {
		return nil
	}
	e = cUnwrap(e)
	if s, ok := e.(CStmtConv); ok {
		arr := s.ToStmt()
		if len(arr) == 1 {
			return arr[0]
		}
	}
	return &CExprStmt{Expr: e}
}

type CExprStmt struct {
	Expr Expr
}

func (s *CExprStmt) Visit(v Visitor) {
	v(s.Expr)
}

func (s *CExprStmt) CType() types.Type {
	return s.Expr.CType(nil)
}

func (s *CExprStmt) AsExpr() GoExpr {
	return s.Expr.AsExpr()
}

func (s *CExprStmt) AsStmt() []GoStmt {
	return cExprStmt(s.Expr)
}

func (s *CExprStmt) Uses() []types.Usage {
	return types.UseRead(s.Expr)
}

type CLabelStmt struct {
	Label string
}

func (s *CLabelStmt) Visit(v Visitor) {}

func (s *CLabelStmt) AsStmt() []GoStmt {
	return []GoStmt{&ast.LabeledStmt{
		Label: ident(s.Label),
		Stmt:  &ast.EmptyStmt{Implicit: true},
	}}
}

func (s *CLabelStmt) Uses() []types.Usage {
	return nil
}

func (g *translator) NewCaseStmt(exp Expr, stmts ...CStmt) *CCaseStmt {
	return &CCaseStmt{g: g, Expr: exp, Stmts: stmts}
}

type CCaseStmt struct {
	g     *translator
	Expr  Expr
	Stmts []CStmt
}

func (s *CCaseStmt) Visit(v Visitor) {
	v(s.Expr)
	for _, st := range s.Stmts {
		v(st)
	}
}

func (s *CCaseStmt) GoCaseClause() *ast.CaseClause {
	stmts := s.g.NewCBlock(s.Stmts...).GoBlockStmt()
	if s.Expr == nil {
		return &ast.CaseClause{
			Body: stmts.List,
		}
	}
	return &ast.CaseClause{
		List: []GoExpr{s.Expr.AsExpr()},
		Body: stmts.List,
	}
}

func (s *CCaseStmt) AsStmt() []GoStmt {
	return []GoStmt{s.GoCaseClause()}
}

func (s *CCaseStmt) Uses() []types.Usage {
	var list []types.Usage
	if s.Expr != nil {
		list = append(list, types.UseRead(s.Expr)...)
	}
	for _, c := range s.Stmts {
		list = append(list, c.Uses()...)
	}
	return list
}

func (g *translator) NewCIfStmt(cond BoolExpr, then []CStmt, els IfElseStmt) *CIfStmt {
	then = g.NewCBlock(then...).Stmts
	//if els != nil {
	//	elss := g.NewCBlock(els).Stmts
	//	if len(elss) != 0 && len(then) == 1 {
	//		// TODO: check if both conditions has no side-effects
	//		_, ok := then[0].(*CIfStmt)
	//		_, ok2 := elss[0].(*CIfStmt)
	//		// "then" branch contains a single if statement
	//		// and "else" contains something else
	//		if ok && !ok2 {
	//			// invert if condition and swap branches
	//			return g.NewCIfStmt(g.cNot(cond), elss, then[0])
	//		}
	//	}
	//	if len(elss) == 1 {
	//		if eif, ok := elss[0].(*CIfStmt); ok {
	//			els = eif
	//		}
	//	}
	//}
	return &CIfStmt{
		g:    g,
		Cond: cond,
		Then: g.NewCBlock(then...),
		Else: els,
	}
}

var _ CCompStmt = &CIfStmt{}

type IfElseStmt interface {
	CStmt
	isElseStmt()
}

type CIfStmt struct {
	g    *translator
	Cond BoolExpr
	Then *BlockStmt
	Else IfElseStmt
}

func (s *CIfStmt) isElseStmt() {}

func (s *CIfStmt) Visit(v Visitor) {
	v(s.Cond)
	v(s.Then)
	v(s.Else)
}

func (g *translator) toElseStmt(st ...CStmt) IfElseStmt {
	if len(st) == 0 {
		return nil
	}
	if len(st) == 1 {
		if e, ok := st[0].(IfElseStmt); ok {
			return e
		}
	}
	return g.newBlockStmt(st...)
}

func (s *CIfStmt) EachStmt(fnc CStmtFunc) bool {
	mod := s.Then.EachStmt(fnc)
	if s.Else == nil {
		return mod
	}
	els := s.g.cBlockCast(s.Else)
	mod2 := els.EachStmt(fnc)
	if len(els.Stmts) == 0 {
		s.Else = nil
	} else if len(els.Stmts) == 1 {
		s.Else = s.g.toElseStmt(els.Stmts[0])
	} else {
		s.Else = els
	}
	return mod || mod2
}

func (s *CIfStmt) AsStmt() []GoStmt {
	cond := s.Cond.AsExpr()
	then := s.Then.GoBlockStmt().List
	var els []GoStmt
	if s.Else != nil {
		els = s.Else.AsStmt()
	}
	return []GoStmt{
		ifelse(cond, then, els),
	}
}

func (s *CIfStmt) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(s.Cond)...)
	list = append(list, s.Then.Uses()...)
	if s.Else != nil {
		list = append(list, s.Else.Uses()...)
	}
	return list
}

func (g *translator) NewCSwitchStmt(cond Expr, stmts []CStmt) *CSwitchStmt {
	sw := &CSwitchStmt{g: g, Cond: cond}
	sw.addStmts(stmts)
	// TODO: fix branches (break, fallthrough)
	return sw
}

var _ CCompStmt = &CSwitchStmt{}

type CSwitchStmt struct {
	g     *translator
	Cond  Expr
	Cases []*CCaseStmt
}

func (s *CSwitchStmt) Visit(v Visitor) {
	v(s.Cond)
	for _, c := range s.Cases {
		v(c)
	}
}

func (s *CSwitchStmt) addStmts(stmts []CStmt) {
	for _, st := range s.g.NewCBlock(stmts...).Stmts {
		if c, ok := st.(*CCaseStmt); ok {
			sub := c.Stmts
			c.Stmts = nil
			s.Cases = append(s.Cases, c)
			if len(sub) != 0 {
				s.addStmts(sub)
			}
		} else {
			last := s.Cases[len(s.Cases)-1]
			last.Stmts = append(last.Stmts, st)
		}
	}
}

func (s *CSwitchStmt) EachStmt(fnc CStmtFunc) bool {
	gmod := false
	for _, c := range s.Cases {
		var mod bool
		c.Stmts, mod = eachStmtIn(c.Stmts, fnc)
		gmod = gmod || mod
	}
	return gmod
}

func (s *CSwitchStmt) AsStmt() []GoStmt {
	var stmts []GoStmt
	for i, c := range s.Cases {
		cs := c.GoCaseClause()
		sub := cs.Body
		if len(sub) == 0 {
			if i != len(s.Cases)-1 {
				sub = []GoStmt{&ast.BranchStmt{Tok: token.FALLTHROUGH}}
			}
		} else {
			switch b := sub[len(sub)-1].(type) {
			case *ast.BranchStmt:
				if b.Tok == token.BREAK {
					sub = sub[:len(sub)-1]
				}
			case *ast.ReturnStmt:
			default:
				if i != len(s.Cases)-1 {
					sub = append(sub, &ast.BranchStmt{Tok: token.FALLTHROUGH})
				}
			}
		}
		cs.Body = sub
		stmts = append(stmts, cs)
	}
	return []GoStmt{&ast.SwitchStmt{
		Tag:  s.Cond.AsExpr(),
		Body: block(stmts...),
	}}
}

func (s *CSwitchStmt) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(s.Cond)...)
	for _, c := range s.Cases {
		list = append(list, c.Uses()...)
	}
	return list
}

func (g *translator) NewCForDeclStmt(init CDecl, cond BoolExpr, iter Expr, stmts []CStmt) *CForStmt {
	if cIsTrue(cond) {
		cond = nil
	}
	f := &CForStmt{
		Init: g.NewCDeclStmt1(init),
		Cond: cond,
		Iter: NewCExprStmt1(iter),
	}
	if body := g.NewCBlock(stmts...); body != nil {
		f.Body = *body
	}
	return f
}

func (g *translator) NewCForStmt(init Expr, cond BoolExpr, iter Expr, stmts []CStmt) *CForStmt {
	if cIsTrue(cond) {
		cond = nil
	}
	f := &CForStmt{
		Init: NewCExprStmt1(init),
		Cond: cond,
		Iter: NewCExprStmt1(iter),
	}
	if body := g.NewCBlock(stmts...); body != nil {
		f.Body = *body
	}
	return f
}

var _ CCompStmt = &CForStmt{}

type CForStmt struct {
	Init CStmt
	Cond Expr
	Iter CStmt
	Body BlockStmt
}

func (s *CForStmt) Visit(v Visitor) {
	v(s.Init)
	v(s.Cond)
	v(s.Iter)
	v(&s.Body)
}

func (s *CForStmt) EachStmt(fnc CStmtFunc) bool {
	return s.Body.EachStmt(fnc)
}

func (s *CForStmt) AsStmt() []GoStmt {
	var init GoStmt
	if s.Init != nil {
		init = cOneStmt(s.Init)
		if s, ok := init.(*ast.DeclStmt); ok {
			if g, ok := s.Decl.(*ast.GenDecl); ok && g.Tok == token.VAR && len(g.Specs) == 1 {
				if sp, ok := g.Specs[0].(*ast.ValueSpec); ok && len(sp.Names) == 1 && len(sp.Values) == 1 {
					init = &ast.AssignStmt{
						Lhs: []ast.Expr{sp.Names[0]},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{Fun: sp.Type, Args: []ast.Expr{sp.Values[0]}},
						},
					}
				}
			}
		}
	}
	var iter GoStmt
	if s.Iter != nil {
		iter = cOneStmt(s.Iter)
	}
	var cond GoExpr
	if s.Cond != nil {
		cond = s.Cond.AsExpr()
	}
	return []GoStmt{&ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: iter,
		Body: s.Body.GoBlockStmt(),
	}}
}

func (s *CForStmt) Uses() []types.Usage {
	var list []types.Usage
	if s.Init != nil {
		list = append(list, s.Init.Uses()...)
	}
	if s.Cond != nil {
		list = append(list, types.UseRead(s.Cond)...)
	}
	if s.Iter != nil {
		list = append(list, s.Iter.Uses()...)
	}
	list = append(list, s.Body.Uses()...)
	return list
}

// forCanBreak returns true in case a for statement can be exited at some point to continue execution of the function.
// This means that it will return false for loops that return from the function directly.
func (g *translator) forCanBreak(s *CForStmt) bool {
	if s.Cond != nil {
		v, ok := cIsBoolConst(s.Cond)
		if !ok {
			// not a constant, so should break at some point... most probably
			return true
		}
		if !v {
			// breaks instantly
			return true
		}
	}
	// "infinite", but can still break manually, so find break statements
	return g.hasBreaks(&s.Body)
}

// switchCanBreak returns true in case a switch statement can be exited at some point to continue execution of the function.
// This means that it will return false for switches that return from all the branches.
func (g *translator) switchCanBreak(s *CSwitchStmt) bool {
	if len(s.Cases) == 0 {
		return true
	}
	for i, c := range s.Cases {
		if len(c.Stmts) == 0 {
			if i == len(c.Stmts)-1 {
				// last is empty - will exit
				return true
			}
			// fallthrough - ignore
			continue
		}
		if !g.allBranchesJump(c.Stmts[len(c.Stmts)-1]) {
			return true
		}
	}
	return false
}

func (g *translator) NewCDoWhileStmt(cond Expr, stmts []CStmt) CStmt {
	stmts = append(stmts, g.NewCIfStmt(
		g.cNot(cond), []CStmt{&CBreakStmt{}}, nil,
	))
	return g.NewCForStmt(nil, nil, nil, stmts)
}

type CGotoStmt struct {
	Label string
}

func (s *CGotoStmt) Visit(v Visitor) {}

func (s *CGotoStmt) AsStmt() []GoStmt {
	return []GoStmt{&ast.BranchStmt{Tok: token.GOTO, Label: ident(s.Label)}}
}

func (s *CGotoStmt) Uses() []types.Usage {
	return nil
}

func countContinue(stmts ...CStmt) int {
	var cnt int
	cEachStmt(func(s CStmt) bool {
		switch s.(type) {
		case *CContinueStmt:
			cnt++
		case *CForStmt:
			return false
		}
		return true
	}, stmts)
	return cnt
}

func countBreak(stmts ...CStmt) int {
	var cnt int
	cEachStmt(func(s CStmt) bool {
		switch s.(type) {
		case *CBreakStmt:
			cnt++
		case *CForStmt:
			return false
		}
		return true
	}, stmts)
	return cnt
}

type CContinueStmt struct{}

func (s *CContinueStmt) Visit(v Visitor) {}

func (s *CContinueStmt) AsStmt() []GoStmt {
	return []GoStmt{&ast.BranchStmt{Tok: token.CONTINUE}}
}

func (s *CContinueStmt) Uses() []types.Usage {
	return nil
}

type CBreakStmt struct{}

func (s *CBreakStmt) Visit(v Visitor) {}

func (s *CBreakStmt) AsStmt() []GoStmt {
	return []GoStmt{&ast.BranchStmt{Tok: token.BREAK}}
}

func (s *CBreakStmt) Uses() []types.Usage {
	return nil
}

func (g *translator) NewReturnStmt(x Expr, rtyp types.Type) []CStmt {
	x = cUnwrap(x)
	switch x := x.(type) {
	case *CTernaryExpr:
		// return (v ? a : b) -> if v { return a } return b
		var stmts []CStmt
		stmts = append(stmts, g.NewCIfStmt(x.Cond, g.NewReturnStmt(x.Then, rtyp), nil))
		stmts = append(stmts, g.NewReturnStmt(x.Else, rtyp)...)
		return stmts
	}
	if rtyp != nil {
		x = g.cCast(rtyp, x)
	}
	return []CStmt{&CReturnStmt{
		Expr: x,
	}}
}

func (g *translator) NewZeroReturnStmt(rtyp types.Type) []CStmt {
	return g.NewReturnStmt(g.ZeroValue(rtyp), rtyp)
}

type CReturnStmt struct {
	Expr Expr
}

func (s *CReturnStmt) Visit(v Visitor) {
	v(s.Expr)
}

func (s *CReturnStmt) AsStmt() []GoStmt {
	if s.Expr == nil {
		return []GoStmt{returnStmt()}
	}
	return []GoStmt{returnStmt(s.Expr.AsExpr())}
}

func (s *CReturnStmt) Uses() []types.Usage {
	var list []types.Usage
	if s.Expr != nil {
		list = append(list, types.UseRead(s.Expr)...)
	}
	return list
}

func (g *translator) setReturnType(ret types.Type, stmts []CStmt) {
	if ret == nil {
		return
	}
	for _, s := range stmts {
		switch s := s.(type) {
		case *CReturnStmt:
			s.Expr = g.cCast(ret, s.Expr)
		case *BlockStmt:
			g.setReturnType(ret, s.Stmts)
		case *CIfStmt:
			g.setReturnType(ret, s.Then.Stmts)
			if s.Else != nil {
				g.setReturnType(ret, []CStmt{s.Else})
			}
		case *CForStmt:
			g.setReturnType(ret, s.Body.Stmts)
		case *CSwitchStmt:
			for _, c := range s.Cases {
				g.setReturnType(ret, c.Stmts)
			}
		}
	}
}

func (g *translator) cBlockCast(s CStmt) *BlockStmt {
	if b, ok := s.(*BlockStmt); ok {
		return b
	}
	return g.newBlockStmt(s)
}

func flattenBlocks(stmts []CStmt) ([]CStmt, bool) {
	mod := false
	for i := 0; i < len(stmts); i++ {
		s, ok := stmts[i].(*BlockStmt)
		if !ok {
			continue
		}
		arr := append([]CStmt{}, stmts[:i]...)
		arr = append(arr, s.Stmts...)
		arr = append(arr, stmts[i+1:]...)
		stmts = arr
		mod = true
		i += len(s.Stmts) - 1
	}
	return stmts, mod
}

func (g *translator) isReturnOrExit(s CStmt) bool {
	switch s := s.(type) {
	case *CReturnStmt:
		return true
	case *CExprStmt:
		switch e := s.Expr.(type) {
		case *CallExpr:
			if id, ok := e.Fun.(Ident); ok {
				switch id.Identifier() {
				case g.env.Go().OsExitFunc(),
					g.env.Go().PanicFunc():
					return true
				}
			}
		}
	}
	return false
}

func isEmptyReturn(s CStmt) bool {
	r, ok := s.(*CReturnStmt)
	return ok && r.Expr == nil
}

func (g *translator) isHardJump(s CStmt) bool {
	switch s.(type) {
	case *CReturnStmt, *CGotoStmt:
		return true
	}
	return false
}

func (g *translator) isJump(s CStmt) bool {
	if g.isHardJump(s) {
		return true
	}
	switch s.(type) {
	case *CContinueStmt, *CBreakStmt:
		return true
	}
	return false
}

func (g *translator) allBranchesJump(s CStmt) bool {
	if g.isJump(s) {
		return true
	}
	switch s := s.(type) {
	case *BlockStmt:
		if len(s.Stmts) == 0 {
			return false
		}
		return g.allBranchesJump(s.Stmts[len(s.Stmts)-1])
	case *CIfStmt:
		if !g.allBranchesJump(s.Then) {
			return false
		}
		if s.Else == nil {
			return false // yes, in other cases we'll miss false condition jump
		}
		return g.allBranchesJump(s.Else)
	case *CSwitchStmt:
		return !g.switchCanBreak(s)
	case *CForStmt:
		return !g.forCanBreak(s)
	}
	return false
}

func dupStmts(stmts []CStmt) []CStmt {
	out, _ := cReplaceEachStmt(func(s CStmt) ([]CStmt, bool) {
		_, ok := s.(*CLabelStmt)
		return nil, ok
	}, stmts)
	return out
}

// invertLastIf finds an if statements that can be inverted.
// It checks in an statement body is smaller than all the other statements
// that follow this if.
// This function can only work in function bodies, because it relies on the returns.
func (g *translator) invertLastIf(stmts []CStmt) ([]CStmt, bool) {
	n := len(stmts)
	if n < 2 {
		return stmts, false
	}
	var (
		ind   = -1
		iff   *CIfStmt
		isRet = false
	)
	for i := n - 1; i > 0; i-- {
		f, ok := stmts[i].(*CIfStmt)
		if !ok || f.Else != nil || len(f.Then.Stmts) == 0 {
			continue
		}
		fn := len(f.Then.Stmts)
		isRet = g.isHardJump(f.Then.Stmts[fn-1])
		if isRet && fn == 1 {
			continue
		}
		if cost := stmtCost(stmts[i+1:]...); !isRet && cost > 2 {
			continue
		} else if cost >= stmtCost(f.Then.Stmts...) {
			continue
		}
		ind, iff = i, f
		break
	}
	if iff == nil {
		return stmts, false
	}
	ret := append([]CStmt{}, stmts[ind+1:]...)
	if len(ret) == 0 || !g.isHardJump(ret[len(ret)-1]) {
		ret = append(ret, &CReturnStmt{})
	}
	iff.Cond = g.cNot(iff.Cond)
	b := iff.Then
	iff.Then = g.NewCBlock(append([]CStmt{}, stmts[ind+1:]...)...)
	if n := len(iff.Then.Stmts); n == 0 || !g.isHardJump(iff.Then.Stmts[n-1]) {
		if l, ok := ret[0].(*CLabelStmt); ok {
			iff.Then.Stmts = append(iff.Then.Stmts, &CGotoStmt{Label: l.Label})
		} else {
			iff.Then.Stmts = append(iff.Then.Stmts, dupStmts(ret)...)
		}
	}
	stmts = append(stmts[:ind+1], b.Stmts...)
	if !g.isHardJump(stmts[len(stmts)-1]) {
		stmts = append(stmts, ret...)
	}
	if isEmptyReturn(stmts[len(stmts)-1]) {
		stmts = stmts[:len(stmts)-1]
	}
	return stmts, true
}

// moveLastReturnToLastIf finds an if statement with else branch, that is before the return and
// moves the return to all the if branches.
// It helps apply other simplifications later.
func (g *translator) moveLastReturnToLastIf(stmts []CStmt) ([]CStmt, bool) {
	n := len(stmts)
	if n < 2 {
		return stmts, false
	}
	ret, ok := stmts[n-1].(*CReturnStmt)
	if !ok {
		return stmts, false
	}
	iff, ok := stmts[n-2].(*CIfStmt)
	if !ok || iff.Else == nil {
		return stmts, false
	}
	stmts = stmts[:n-1]
	n--
	allReturns := true
	mod := false
	for iff != nil {
		if !g.isJump(iff.Then.Stmts[len(iff.Then.Stmts)-1]) {
			iff.Then.Stmts = append(iff.Then.Stmts, ret)
			mod = true
		}
		if iff.Else == nil {
			allReturns = false
			break
		}
		els := iff.Else
		switch els := els.(type) {
		case *BlockStmt:
			if !g.isJump(els.Stmts[len(els.Stmts)-1]) {
				els.Stmts = append(els.Stmts, ret)
				mod = true
			}
			iff = nil
		case *CIfStmt:
			iff = els
		default:
			if !g.isJump(els) {
				iff.Else = g.NewCBlock(els, ret)
				mod = true
			}
			iff = nil
		}
	}
	if !allReturns {
		stmts = append(stmts, ret)
	}
	return stmts, mod
}

// inlineElseInLastReturn finds a return in the end and inlines its else branch
// if all the branches contain jumps at the end.
func (g *translator) inlineElseInLastReturn(stmts []CStmt) ([]CStmt, bool) {
	n := len(stmts)
	if n < 1 {
		return stmts, false
	}
	iff, ok := stmts[n-1].(*CIfStmt)
	if !ok || iff.Else == nil {
		return stmts, false
	}
	if !g.allBranchesJump(iff.Then) || !g.allBranchesJump(iff.Else) {
		return stmts, false
	}
	els := iff.Else
	iff.Else = nil
	switch els := els.(type) {
	case *BlockStmt:
		stmts = append(stmts, els.Stmts...)
	default:
		stmts = append(stmts, els)
	}
	return stmts, true
}

func (g *translator) inlineSmallGotos(stmts []CStmt) ([]CStmt, bool) {
	const maxCost = 5
	// first, collect all paths from the label
	// we only care about direct paths that lead to a "hard jump" (return or another goto)
	labels := make(map[string][]CStmt)
	seen := make(map[*CLabelStmt]struct{})
	cEachBlockStmt(func(stmts []CStmt) {
		for i, s := range stmts {
			l, ok := s.(*CLabelStmt)
			if !ok {
				continue
			} else if _, ok := seen[l]; ok {
				continue
			}
			seen[l] = struct{}{}
			cnt := 0
			for j, s := range stmts[i:] {
				if g.isHardJump(s) {
					cnt = j + 1
					break
				}
			}
			if cnt == 0 {
				continue
			}
			body := append([]CStmt{}, stmts[i+1:i+cnt]...)
			if stmtCost(body...) > maxCost {
				continue
			}
			labels[l.Label] = body
		}
	}, stmts)
	seen = nil
	// rescan label bodies and replace gotos in them as well
	for name, body := range labels {
		mod, del := false, false
		body, _ = cReplaceEachStmt(func(s CStmt) ([]CStmt, bool) {
			if del {
				return nil, false
			}
			g, ok := s.(*CGotoStmt)
			if !ok {
				return nil, false
			}
			if name == g.Label {
				// loop - delete the label
				delete(labels, name)
				del = true
				return nil, false
			}
			if body := labels[g.Label]; len(body) != 0 {
				mod = true
				return body, true
			}
			return nil, false
		}, body)
		if mod && !del {
			labels[name] = body
		}
	}
	if len(labels) == 0 {
		return stmts, false
	}
	// replace gotos with those paths and remove labels
	stmts, _ = cReplaceEachStmt(func(s CStmt) ([]CStmt, bool) {
		if l, ok := s.(*CLabelStmt); ok {
			if len(labels[l.Label]) != 0 {
				return nil, true
			}
		}
		g, ok := s.(*CGotoStmt)
		if !ok {
			return nil, false
		}
		if body := labels[g.Label]; len(body) != 0 {
			return body, true
		}
		return nil, false
	}, stmts)
	return stmts, true
}

func (g *translator) NewCBlock(stmts ...CStmt) *BlockStmt {
	if len(stmts) == 1 {
		if b, ok := stmts[0].(*BlockStmt); ok {
			return b
		}
	}
	// TODO:
	//if len(stmts) >= 2 {
	//	l := len(stmts)
	//	ret, ok1 := stmts[l-1].(*ast.ReturnStmt)
	//	ifs, ok2 := stmts[l-2].(*ast.IfStmt)
	//	if ok1 && ok2 && ifs.Else == nil && len(ifs.Body.List) > 1 {
	//		ifs.Cond = not(ifs.Cond)
	//		body := ifs.Body.List
	//		ifs.Body.List = []GoStmt{ret}
	//		if _, ok := body[len(body)-1].(*ast.ReturnStmt); !ok {
	//			body = append(body, ret)
	//		}
	//		stmts = append(stmts[:l-1], body...)
	//		return block(stmts...)
	//	}
	//}
	for {
		var mod bool
		stmts, mod = flattenBlocks(stmts)
		if mod {
			continue
		}
		if optimizeStatements {
			stmts, mod = replaceBreaks(stmts)
			if mod {
				continue
			}
			stmts, mod = rebuildFors(stmts)
			if mod {
				continue
			}
		}
		break
	}
	return g.newBlockStmt(stmts...)
}

func isIfReturn(s CStmt) (Expr, *CReturnStmt, bool) {
	iff, ok := s.(*CIfStmt)
	if !ok || iff.Else != nil || len(iff.Then.Stmts) != 1 {
		return nil, nil, false
	}
	ret, ok := iff.Then.Stmts[0].(*CReturnStmt)
	if !ok {
		return nil, nil, false
	}
	return iff.Cond, ret, true
}

//func tryMergeForPreCond(cond Expr, ret *CReturnStmt, st CStmt) bool {
//	forr, ok := st.(*CForStmt)
//	if !ok || forr.Cond != nil || forr.Iter != nil || forr.Body == nil {
//		return false
//	}
//	n := len(forr.Body.Stmts)
//	if n == 0 {
//		return false
//	}
//	cond2, ret2, ok := isIfReturn(forr.Body.Stmts[n-1])
//	if !ok {
//		return false
//	}
//	if countContinue(forr.Body.Stmts...) != 0 {
//		return false
//	}
//	if
//}

func revFindInfiniteFor(stmts []CStmt) int {
	for i := len(stmts) - 1; i >= 0; i-- {
		if f, ok := stmts[i].(*CForStmt); ok && f.Cond == nil {
			return i
		}
	}
	return -1
}

func replaceBreaks(stmts []CStmt) ([]CStmt, bool) {
	n := len(stmts)
	if n < 2 {
		return stmts, false
	}
	_, ok := stmts[n-1].(*CReturnStmt)
	if !ok {
		return stmts, false
	}
	for i := n - 1; i >= 0 && n-i < 5; i = revFindInfiniteFor(stmts[:i]) {
		forr, ok := stmts[i].(*CForStmt)
		if !ok {
			continue
		}
		part := stmts[i+1:]
		b := countBreak(&forr.Body)
		if b != 1 {
			continue
		}
		r := 0
		var replace CStmtFunc
		replace = func(s CStmt) ([]CStmt, bool) {
			switch s := s.(type) {
			case *CBreakStmt:
				r++
				return part, true
			case *CForStmt:
				return nil, false
			case CCompStmt:
				s.EachStmt(replace)
			}
			return nil, false
		}
		forr.EachStmt(replace)
		if r != b {
			panic(r)
		}
		return stmts[:i+1], true
	}
	return stmts, false
}

func rebuildFors(stmts []CStmt) ([]CStmt, bool) {
	n := len(stmts)
	if n < 2 {
		return stmts, false
	}
	//for i, s := range stmts[:n-1] {
	//	j := i+1
	//	s2 := stmts[j]
	//	if cond, ret, ok := isIfReturn(s); ok && tryMergeForPreCond(cond, ret, s2) {
	//		stmts = append(stmts[:i], stmts[i+1:]...)
	//		continue rules
	//	}
	//}
	return stmts, false
	// TODO:
	//for i, s := range stmts[:len(stmts)-1] {
	//	as, ok := s.(*ast.AssignStmt)
	//	if !ok || len(as.Lhs) != 1 || as.Tok != token.ASSIGN {
	//		continue
	//	}
	//	x, ok := as.Lhs[0].(*ast.Ident)
	//	if !ok {
	//		continue
	//	}
	//	if lit, ok := as.Rhs[0].(*ast.BasicLit); !ok || lit.Kind != token.INT {
	//		continue
	//	}
	//	loop, ok := stmts[i+1].(*ast.ForStmt)
	//	if !ok || loop.Init != nil || loop.Cond != nil || loop.Post != nil || loop.Body == nil || len(loop.Body.List) < 2 {
	//		continue
	//	}
	//	body := loop.Body.List
	//	l := len(body)
	//	inc, ok := body[l-2].(*ast.IncDecStmt)
	//	if !ok {
	//		continue
	//	} else if x2, ok := inc.X.(*ast.Ident); !ok || x2.Name != x.Name {
	//		continue
	//	}
	//	ifc, ok := body[l-1].(*ast.IfStmt)
	//	if !ok || ifc.Init != nil || ifc.Else != nil || len(ifc.Body.List) != 1 {
	//		continue
	//	} else if b, ok := ifc.Body.List[0].(*ast.BranchStmt); !ok || b.Tok != token.BREAK || b.Label != nil {
	//		continue
	//	}
	//	loop.Init = as
	//	loop.Cond = not(ifc.Cond)
	//	loop.Post = inc
	//	loop.Body.List = body[:l-2]
	//	return append(stmts[:i], rebuildFors(stmts[i+1:])...)
	//}
}

var _ CCompStmt = &BlockStmt{}

func (g *translator) newBlockStmt(stmts ...CStmt) *BlockStmt {
	return &BlockStmt{
		g:     g,
		Stmts: stmts,
	}
}

type BlockStmt struct {
	g     *translator
	Stmts []CStmt
}

func (s *BlockStmt) isElseStmt() {}

func (s *BlockStmt) Visit(v Visitor) {
	for _, st := range s.Stmts {
		v(st)
	}
}

func eachStmtIn(stmts []CStmt, fnc CStmtFunc) ([]CStmt, bool) {
	gmod := false
	for i := 0; i < len(stmts); i++ {
		arr, mod := fnc(stmts[i])
		if !mod {
			continue
		}
		if !gmod {
			stmts = append([]CStmt{}, stmts...)
		}
		gmod = true
		if d := len(arr); d == 0 {
			stmts = append(stmts[:i], stmts[i+1:]...)
			i--
		} else if d == 1 {
			stmts[i] = arr[0]
		} else if n := len(stmts); i == n-1 {
			stmts = append(stmts[:i], arr...)
			break
		} else {
			tmp := append([]CStmt{}, stmts[:i]...)
			tmp = append(tmp, arr...)
			tmp = append(tmp, stmts[i+1:]...)
			stmts = tmp
			i += len(arr) - 1
		}
	}
	return stmts, gmod
}

func (s *BlockStmt) EachStmt(fnc CStmtFunc) bool {
	var mod bool
	s.Stmts, mod = eachStmtIn(s.Stmts, fnc)
	return mod
}

func (s *BlockStmt) In(ft *types.FuncType) *BlockStmt {
	if ft.Return() != nil {
		s.g.setReturnType(ft.Return(), s.Stmts)
	}
	stmts := s.Stmts
	for optimizeStatements {
		var mod bool
		stmts, mod = s.g.moveLastReturnToLastIf(stmts)
		if mod {
			continue
		}
		stmts, mod = s.g.inlineElseInLastReturn(stmts)
		if mod {
			continue
		}
		stmts, mod = s.g.invertLastIf(stmts)
		if mod {
			continue
		}
		stmts, mod = s.g.inlineSmallGotos(stmts)
		if mod {
			continue
		}
		break
	}
	s.Stmts = stmts
	return s
}

func asStmts(arr []CStmt) []GoStmt {
	var out []GoStmt
	for _, s := range arr {
		out = append(out, s.AsStmt()...)
	}
	return out
}

func (s *BlockStmt) GoBlockStmt() *ast.BlockStmt {
	if s == nil {
		return nil
	}
	return block(asStmts(s.Stmts)...)
}

func (s *BlockStmt) AsStmt() []GoStmt {
	return []GoStmt{s.GoBlockStmt()}
}

func (s *BlockStmt) Uses() []types.Usage {
	var list []types.Usage
	for _, s := range s.Stmts {
		list = append(list, s.Uses()...)
	}
	return list
}

func (g *translator) NewCDeclStmt1(decl CDecl) CStmt {
	return &CDeclStmt{Decl: decl}
}

func (g *translator) NewCDeclStmt(decl CDecl) []CStmt {
	switch d := decl.(type) {
	case *CVarDecl:
		for i, v := range d.Inits {
			switch v.(type) {
			case *CTernaryExpr:
				var stmts []CStmt
				if pre := d.Inits[:i:i]; len(pre) > 0 {
					dc := *d
					dc.Inits = pre
					dc.Names = dc.Names[:i:i]
					stmts = append(stmts, g.NewCDeclStmt(&dc)...)
				}
				dt := *d
				dt.Inits = nil
				dt.Names = []*types.Ident{d.Names[i]}
				stmts = append(stmts, g.NewCDeclStmt(&dt)...)
				stmts = append(stmts, g.NewCAssignStmt(IdentExpr{d.Names[i]}, "", v)...)
				if post := d.Inits[i+1:]; len(post) > 0 {
					dc := *d
					dc.Inits = post
					dc.Names = dc.Names[i+1:]
					stmts = append(stmts, g.NewCDeclStmt(&dc)...)
				}
				return stmts
			}
		}
	}
	return []CStmt{&CDeclStmt{Decl: decl}}
}

type CDeclStmt struct {
	Decl CDecl
}

func (s *CDeclStmt) Visit(v Visitor) {
	v(s.Decl)
}

func (s *CDeclStmt) AsStmt() []GoStmt {
	var out []GoStmt
	for _, d := range s.Decl.AsDecl() {
		out = append(out, &ast.DeclStmt{Decl: d})
	}
	return out
}

func (s *CDeclStmt) Uses() []types.Usage {
	return s.Decl.Uses()
}

func (g *translator) NewCIncStmt(x Expr, decr bool) *CIncrStmt {
	return &CIncrStmt{
		g:    g,
		Expr: x,
		Decr: decr,
	}
}

type CIncrStmt struct {
	g    *translator
	Expr Expr
	Decr bool
}

func (s *CIncrStmt) Visit(v Visitor) {
	v(s.Expr)
}

func (s *CIncrStmt) AsStmt() []GoStmt {
	if s.Expr.CType(nil).Kind().IsPtr() {
		var arg Expr
		if s.Decr {
			arg = cIntLit(-1)
		} else {
			arg = cIntLit(+1)
		}
		x := cPtrOffset(s.g.ToPointer(s.Expr), arg)
		return asStmts(s.g.NewCAssignStmt(s.Expr, "", x))
	}
	var tok token.Token
	if s.Decr {
		tok = token.DEC
	} else {
		tok = token.INC
	}
	return []GoStmt{&ast.IncDecStmt{
		X:   s.Expr.AsExpr(),
		Tok: tok,
	}}
}

func (s *CIncrStmt) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseRead(s.Expr)...)
	list = append(list, types.UseWrite(s.Expr)...)
	return list
}

func (g *translator) NewCAssignStmtP(x Expr, op BinaryOp, y Expr) *CAssignStmt {
	x = cUnwrap(x)
	y = cUnwrap(y)
	r := g.cCast(x.CType(nil), y)
	return &CAssignStmt{
		g:     g,
		Left:  x,
		Op:    op,
		Right: r,
	}
}

func (g *translator) NewCAssignStmt(x Expr, op BinaryOp, y Expr) []CStmt {
	x = cUnwrap(x)
	y = cUnwrap(y)
	if !x.HasSideEffects() {
		switch y := unwrapCasts(y).(type) {
		case *CUnaryExpr:
			switch z := unwrapCasts(y.Expr).(type) {
			case *CTernaryExpr:
				// v = -(x ? a : b) -> if x { v = -a } else { v = -b }
				return []CStmt{
					g.NewCIfStmt(z.Cond,
						g.NewCAssignStmt(x, op, g.NewCUnaryExpr(y.Op, z.Then)),
						g.toElseStmt(g.NewCAssignStmt(x, op, g.NewCUnaryExpr(y.Op, z.Else))...),
					),
				}
			}
		case *CTernaryExpr:
			// v = (x ? a : b) -> if x { v = a } else { v = b }
			return []CStmt{
				g.NewCIfStmt(y.Cond,
					g.NewCAssignStmt(x, op, y.Then),
					g.toElseStmt(g.NewCAssignStmt(x, op, y.Else)...),
				),
			}
		}
	}
	if x.CType(nil).Kind().IsPtr() {
		switch op {
		case BinOpAdd:
			if e, ok := y.(*IntToPtr); ok {
				y = e.X
			}
			op = ""
			y = cPtrOffset(g.ToPointer(x), y)
		case BinOpSub:
			if e, ok := y.(*IntToPtr); ok {
				y = e.X
			}
			op = ""
			y = g.NewCUnaryExpr(UnaryMinus, y)
			y = cPtrOffset(g.ToPointer(x), y)
		}
	}
	return []CStmt{&CAssignStmt{
		g:     g,
		Left:  x,
		Op:    op,
		Right: g.cCast(x.CType(nil), y),
	}}
}

type CAssignStmt struct {
	g     *translator
	Left  Expr
	Op    BinaryOp
	Right Expr
}

func (s *CAssignStmt) Visit(v Visitor) {
	v(s.Left)
	v(s.Right)
}

func (s *CAssignStmt) AsStmt() []GoStmt {
	x := s.Left.AsExpr()
	y := s.Right.AsExpr()
	return []GoStmt{
		assignTok(x, s.Op.GoAssignToken(), y),
	}
}

func (s *CAssignStmt) Uses() []types.Usage {
	var list []types.Usage
	list = append(list, types.UseWrite(s.Left)...)
	list = append(list, types.UseRead(s.Right)...)
	return list
}

type UnusedVar struct {
	Name *types.Ident
}

func (s *UnusedVar) Visit(v Visitor) {
}

func (s *UnusedVar) AsStmt() []GoStmt {
	return []GoStmt{
		assign(ident("_"), s.Name.GoIdent()),
	}
}

func (s *UnusedVar) Uses() []types.Usage {
	return []types.Usage{{Ident: s.Name, Access: types.AccessRead}}
}

func stmtCost(stmts ...CStmt) int {
	c := 0
	for _, s := range stmts {
		if s == nil {
			continue
		}
		switch s := s.(type) {
		case *CForStmt:
			c += 3 + stmtCost(&s.Body)
		case *CSwitchStmt:
			c++
			for _, s2 := range s.Cases {
				c += stmtCost(s2.Stmts...)
			}
		case *CIfStmt:
			c += 1 + stmtCost(s.Then) + stmtCost(s.Else)
		case *CCaseStmt:
			c += 1 + stmtCost(s.Stmts...)
		case *BlockStmt:
			c += stmtCost(s.Stmts...)
		case *CDeclStmt:
			c += 2
		default:
			c++
		}
	}
	return c
}

func (g *translator) hasReturns(stmt ...CStmt) bool {
	for _, st := range stmt {
		if g.isReturnOrExit(st) {
			return true
		}
		switch st := st.(type) {
		case *BlockStmt:
			if g.hasReturns(st.Stmts...) {
				return true
			}
		case *CIfStmt:
			if g.hasReturns(st.Then) {
				return true
			}
			if st.Else != nil {
				if g.hasReturns(st.Else) {
					return true
				}
			}
		case *CForStmt:
			if g.hasReturns(&st.Body) {
				return true
			}
		case *CSwitchStmt:
			for _, c := range st.Cases {
				if g.hasReturns(c.Stmts...) {
					return true
				}
			}
		}
	}
	return false
}

func (g *translator) hasBreaks(stmt ...CStmt) bool {
	for _, st := range stmt {
		switch st := st.(type) {
		case *CBreakStmt:
			return true
		case *CGotoStmt:
			// TODO: check if it causes exit from the loop
		case *BlockStmt:
			if g.hasBreaks(st.Stmts...) {
				return true
			}
		case *CIfStmt:
			if g.hasBreaks(st.Then) {
				return true
			}
			if st.Else != nil {
				if g.hasBreaks(st.Else) {
					return true
				}
			}
		}
	}
	return false
}

func (g *translator) hasBreaksOrReturns(stmt ...CStmt) bool {
	for _, st := range stmt {
		if g.hasReturns(st) {
			return true
		}
		switch st := st.(type) {
		case *CBreakStmt:
			return true
		case *CGotoStmt:
			// TODO: check if it causes exit from the loop
		case *BlockStmt:
			if g.hasBreaksOrReturns(st.Stmts...) {
				return true
			}
		case *CIfStmt:
			if g.hasBreaksOrReturns(st.Then) {
				return true
			}
			if st.Else != nil {
				if g.hasBreaksOrReturns(st.Else) {
					return true
				}
			}
		case *CForStmt:
			// don't consider breaks since it's a new loop
			if g.hasReturns(&st.Body) {
				return true
			}
		case *CSwitchStmt:
			// don't consider breaks since they affect the switch
			for _, c := range st.Cases {
				if g.hasReturns(c.Stmts...) {
					return true
				}
			}
		}
	}
	return false
}
