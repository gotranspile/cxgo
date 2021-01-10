package cxgo

import "github.com/gotranspile/cxgo/types"

func (g *translator) adaptMain(decl []CDecl) []CDecl {
	for _, d := range decl {
		f, ok := d.(*CFuncDecl)
		if !ok || f.Body == nil {
			continue
		}
		switch f.Name.Name {
		case "main":
			g.translateMain(f)
		}
	}
	return decl
}

func (g *translator) flatten(decl []CDecl) {
	for _, d := range decl {
		f, ok := d.(*CFuncDecl)
		if !ok || f.Body == nil {
			continue
		}
		flatten := g.conf.FlattenAll
		if c, ok := g.idents[f.Name.Name]; ok && c.Flatten != nil {
			flatten = *c.Flatten
		}
		if !flatten {
			continue
		}
		cf := g.NewControlFlow(f.Body.Stmts)
		f.Body.Stmts = cf.Flatten()
	}
}

func (g *translator) fixUnusedVars(decl []CDecl) {
	for _, d := range decl {
		switch d := d.(type) {
		case *CFuncDecl:
			g.fixUnusedVarsBlock(d.Body)
		}
	}
}

func (g *translator) fixUnusedVarsBlock(b *BlockStmt) {
	if b != nil {
		b.Stmts = g.fixUnusedVarsStmts(b.Stmts)
	}
}

func (g *translator) fixUnusedVarsStmts(stmts []CStmt) []CStmt {
	for i := 0; i < len(stmts); i++ {
		st := stmts[i]
		var next []CStmt
		if i != len(stmts)-1 {
			next = stmts[i+1:]
		}
		add := g.unusedVarsStmt(st, next)
		if len(add) == 0 {
			continue
		}
		arr := make([]CStmt, 0, len(stmts)+len(add))
		arr = append(arr, stmts[:i+1]...)
		arr = append(arr, add...)
		arr = append(arr, next...)
		stmts = arr
		i += len(add)
	}
	return stmts
}

func (g *translator) unusedVarsStmt(st CStmt, next []CStmt) []CStmt {
	switch st := st.(type) {
	case *BlockStmt:
		g.fixUnusedVarsBlock(st)
		return nil
	case *CForStmt:
		g.fixUnusedVarsBlock(&st.Body)
		return nil
	case *CIfStmt:
		g.fixUnusedVarsBlock(st.Then)
		g.unusedVarsStmt(st.Else, nil)
		return nil
	case *CDeclStmt:
		vd, ok := st.Decl.(*CVarDecl)
		if !ok || vd.Const {
			return nil
		}
		var add []CStmt
		if len(next) == 0 {
			// everything in the last declaration is always unused
			for _, name := range vd.Names {
				if name.Name == "__func__" {
					continue
				}
				add = append(add, &UnusedVar{Name: name})
			}
			return add
		}
		for i, name := range vd.Names {
			if name.Name == "__func__" {
				continue
			}
			used := false
			if i+1 < len(vd.Inits) {
				// we should scan inits of following variables in this block
				for _, e := range vd.Inits[i+1:] {
					for _, u := range types.UseRead(e) {
						if name == u.Ident && u.Access == types.AccessRead {
							used = true
							break
						}
					}
				}
			}
			if !used {
				for _, st := range next {
					for _, u := range st.Uses() {
						if name == u.Ident && u.Access == types.AccessRead {
							used = true
							break
						}
					}
				}
			}
			if !used {
				add = append(add, &UnusedVar{Name: name})
			}
		}
		return add
	default:
		return nil
	}
}

func (g *translator) fixImplicitReturns(decl []CDecl) {
	for _, d := range decl {
		switch d := d.(type) {
		case *CFuncDecl:
			ret := d.Type.Return()
			if ret == nil || d.Body == nil {
				continue
			}
			d.Body.Stmts = g.fixImplicitReturnStmts(ret, d.Body.Stmts)
		}
	}
}

func (g *translator) fixImplicitReturnStmts(ret types.Type, stmts []CStmt) []CStmt {
	if len(stmts) == 0 {
		return g.NewZeroReturnStmt(ret)
	}
	last := stmts[len(stmts)-1]
	if g.isHardJump(last) {
		return stmts
	}
	if !g.fixImplicitReturnStmt(ret, last) {
		stmts = append(stmts, g.NewZeroReturnStmt(ret)...)
	}
	return stmts
}

func (g *translator) fixImplicitReturnStmt(ret types.Type, st CStmt) bool {
	if g.isReturnOrExit(st) {
		return true
	}
	switch st := st.(type) {
	case *CGotoStmt:
		return true
	case *BlockStmt:
		st.Stmts = g.fixImplicitReturnStmts(ret, st.Stmts)
		return true
	case *CIfStmt:
		if st.Else == nil {
			return false
		}
		if !g.fixImplicitReturnStmt(ret, st.Else) {
			return false
		}
		st.Then.Stmts = g.fixImplicitReturnStmts(ret, st.Then.Stmts)
		return true
	case *CForStmt:
		return !g.forCanBreak(st)
	case *CSwitchStmt:
		return !g.switchCanBreak(st)
	default:
		return false
	}
}
