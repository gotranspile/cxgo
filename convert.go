package cxgo

import (
	"fmt"
	"strings"

	"modernc.org/cc/v3"
	"modernc.org/token"

	"github.com/gotranspile/cxgo/types"
)

func (g *translator) newIdent(name string, t types.Type) *types.Ident {
	if id, ok := g.env.IdentByName(name); ok {
		if it := id.CType(nil); it.Kind().Major() == t.Kind().Major() || (it.Kind().IsBool() && t.Kind().IsInt()) {
			return id // FIXME: this is invalid, we should consult scopes instead
		}
	}
	return types.NewIdent(name, t)
}

func (g *translator) convMacro(name string, fnc func() Expr) Expr {
	if g.env.ForceMacro(name) {
		return fnc()
	}
	id, ok := g.macros[name]
	if ok {
		return IdentExpr{id}
	}
	x := fnc()
	typ := x.CType(nil)
	id = g.newIdent(name, typ)
	g.macros[name] = id
	return IdentExpr{id}
}

func (g *translator) convertIdent(scope cc.Scope, tok cc.Token, t types.Type) IdentExpr {
	var decl []cc.Node
	for len(scope) != 0 {
		if nodes, ok := scope[tok.Value]; ok {
			decl = nodes
			break
		}
		scope = scope.Parent()
	}
	if len(decl) == 0 {
		panic(fmt.Errorf("unresolved identifier: %s (%s)", tok, tok.Position()))
	}
	return g.convertIdentWith(tok.String(), t, decl...)
}

func (g *translator) convertIdentWith(name string, t types.Type, decls ...cc.Node) IdentExpr {
	for _, d := range decls {
		if id, ok := g.decls[d]; ok && id.Name == name {
			return IdentExpr{id}
		}
	}
	id := g.newIdent(name, t)
	if to, ok := g.idents[name]; ok && to.Rename != "" {
		id.GoName = to.Rename
	}
	for _, d := range decls {
		g.decls[d] = id
	}
	return IdentExpr{id}
}

func (g *translator) replaceIdentWith(id *types.Ident, decls ...cc.Node) {
	for _, d := range decls {
		g.decls[d] = id
	}
}

func (g *translator) tryConvertIdentOn(t types.Type, tok cc.Token) (*types.Ident, bool) {
loop:
	for {
		switch s := t.(type) {
		case types.PtrType:
			t = s.Elem()
		case types.ArrayType:
			t = s.Elem()
		case types.Named:
			t = s.Underlying()
		default:
			break loop
		}
	}
	switch t := t.(type) {
	case *types.StructType:
		name := tok.Value.String()
		for _, f := range t.Fields() {
			if name == f.Name.Name {
				return f.Name, true
			}
			if f.Name.Name == "" {
				if id, ok := g.tryConvertIdentOn(f.Type(), tok); ok {
					return id, true
				}
			}
		}
	}
	return nil, false
}

func (g *translator) convertIdentOn(t types.Type, tok cc.Token) *types.Ident {
	id, ok := g.tryConvertIdentOn(t, tok)
	if ok {
		return id
	}
	panic(fmt.Errorf("%#v.%q (%s)", t, tok.Value.String(), tok.Position()))
}

func (g *translator) convertFuncDef(d *cc.FunctionDefinition) []CDecl {
	decl := d.Declarator
	switch dd := decl.DirectDeclarator; dd.Case {
	case cc.DirectDeclaratorFuncParam, cc.DirectDeclaratorFuncIdent:
		sname := decl.Name().String()
		conf := g.idents[sname]
		ft := g.convertFuncType(conf, decl, decl.Type(), decl.Position())
		if !g.inCurFile(d) {
			return nil
		}
		name := g.convertIdentWith(sname, ft, decl)
		return []CDecl{
			&CFuncDecl{
				Name: name.Ident,
				Type: ft,
				Body: g.convertCompBlockStmt(d.CompoundStatement).In(ft),
				Range: &Range{
					Start:     d.Position().Offset,
					StartLine: d.Position().Line,
				},
			},
		}
	default:
		panic(dd.Case.String() + " " + dd.Position().String())
	}
}

type positioner interface {
	Position() token.Position
}

func (g *translator) inCurFile(p positioner) bool {
	name := strings.TrimLeft(p.Position().Filename, "./")
	if g.cur == name {
		return true
	} else if !strings.HasSuffix(g.cur, ".c") {
		return false
	}
	return g.cur[:len(g.cur)-2]+".h" == name
}

func (g *translator) convertInitList(typ types.Type, list *cc.InitializerList) Expr {
	var items []*CompLitField
	var (
		prev int64 = -1 // previous array init index
		pi   int64 = 0  // relative index added to the last seen item; see below
	)
	for it := list; it != nil; it = it.InitializerList {
		val := g.convertInitExpr(it.Initializer)
		var f *CompLitField
		if it.Designation == nil {
			// no index in the initializer - assign automatically
			pi++
			f = &CompLitField{Index: cIntLit(prev+pi, 10), Value: val}
			items = append(items, f)
			continue
		}
		f = g.convertOneDesignator(typ, it.Designation.DesignatorList, val)
		if lit, ok := f.Index.(IntLit); ok {
			if prev == -1 {
				// first item - note that we started initializing indexes
				prev = 0
			} else if prev == lit.Int() {
				// this was an old bug in CC where it returned stale indexes
				// for items without any index designators
				// it looks like it is fixed now, but we keep the workaround just in case
				pi++
				f.Index = cIntLit(prev+pi, 10)
			} else {
				// valid index - set previous and reset relative index
				prev = lit.Int()
				pi = 0
			}
		}
		items = append(items, f)
	}
	return g.NewCCompLitExpr(
		typ,
		items,
	)
}

func (g *translator) convertInitExpr(d *cc.Initializer) Expr {
	switch d.Case {
	case cc.InitializerExpr:
		return g.convertAssignExpr(d.AssignmentExpression)
	case cc.InitializerInitList:
		return g.convertInitList(
			g.convertTypeRoot(IdentConfig{}, d.Type(), d.Position()),
			d.InitializerList,
		)
	default:
		panic(d.Case.String() + " " + d.Position().String())
	}
}

func (g *translator) convertEnum(b *cc.Declaration, typ types.Type, d *cc.EnumSpecifier) []CDecl {
	if d.EnumeratorList == nil {
		return nil
	}
	if typ == nil {
		typ = types.UntypedIntT(g.env.IntSize())
	}
	vd := &CVarDecl{
		Const:    true,
		Single:   false,
		CVarSpec: CVarSpec{g: g, Type: typ},
	}
	var (
		autos  = 0 // number of implicit inits
		values = 0 // number of explicit inits
	)
	for it := d.EnumeratorList; it != nil; it = it.EnumeratorList {
		e := it.Enumerator
		if e.Case == cc.EnumeratorExpr {
			init := g.convertConstExpr(e.ConstantExpression)
			vd.Inits = append(vd.Inits, init)
			values++
			continue
		}
		vd.Inits = append(vd.Inits, nil)
		autos++
	}
	if autos == 1 && vd.Inits[0] == nil {
		autos--
		values++
		vd.Inits[0] = cIntLit(0, 10)
	}

	// use iota if there is only one explicit init (the first one), or no explicit values are set
	isIota := (vd.Inits[0] == nil && autos == 1) || autos == len(vd.Inits)
	if len(vd.Inits) > 1 && vd.Inits[0] != nil && values == 1 {
		if _, ok := cUnwrap(vd.Inits[0]).(IntLit); ok {
			isIota = true
			values--
			autos++
		}
	}

	if len(vd.Inits) != 0 && isIota && values != 0 {
		panic("TODO: mixed enums")
	}
	var next int64
	for it, i := d.EnumeratorList, 0; it != nil; it, i = it.EnumeratorList, i+1 {
		e := it.Enumerator
		if isIota {
			if i == 0 {
				iot := g.Iota()
				if val := vd.Inits[0]; val != nil {
					if l, ok := cUnwrap(val).(IntLit); !ok || !l.IsZero() {
						iot = &CBinaryExpr{Left: iot, Op: BinOpAdd, Right: val}
					}
				}
				if !typ.Kind().IsUntypedInt() {
					iot = &CCastExpr{Type: typ, Expr: iot}
				}
				vd.Inits[0] = iot
			}
		} else {
			if vd.Inits[i] == nil {
				vd.Inits[i] = cIntLit(next, 10)
				next++
			} else if l, ok := cUnwrap(vd.Inits[i]).(IntLit); ok {
				next = l.Int() + 1
			}
		}
		vd.Names = append(vd.Names, g.convertIdentWith(e.Token.Value.String(), typ, e).Ident)
	}
	if len(vd.Names) == 0 {
		return nil
	}
	if isIota {
		vd.Type = nil
	}
	return []CDecl{vd}
}

func (g *translator) convertTypedefName(d *cc.Declaration) (cc.Token, *cc.Declarator) {
	if d.InitDeclaratorList == nil || d.InitDeclaratorList.InitDeclarator == nil {
		panic("no decl")
	}
	if d.InitDeclaratorList.InitDeclaratorList != nil {
		panic("should have one decl")
	}
	id := d.InitDeclaratorList.InitDeclarator
	if id.Case != cc.InitDeclaratorDecl {
		panic(id.Case.String())
	}
	dd := id.Declarator
	if dd.DirectDeclarator.Case != cc.DirectDeclaratorIdent {
		panic(dd.DirectDeclarator.Case)
	}
	return dd.DirectDeclarator.Token, dd
}

func (g *translator) convertDecl(d *cc.Declaration) []CDecl {
	inCur := g.inCurFile(d)
	var (
		isConst    bool
		isVolatile bool
		isTypedef  bool
		isStatic   bool
		isExtern   bool
		isForward  bool
		isFunc     bool
		isPrim     bool
		isAuto     bool
		typeSpec   types.Type
		enumSpec   *cc.EnumSpecifier
		names      []string // used only for the hooks
	)
	for il := d.InitDeclaratorList; il != nil; il = il.InitDeclaratorList {
		id := il.InitDeclarator
		switch id.Case {
		case cc.InitDeclaratorDecl, cc.InitDeclaratorInit:
			dd := id.Declarator
			if name := dd.Name().String(); name != "" {
				names = append(names, name)
			}
		}
	}
	spec := d.DeclarationSpecifiers
	if spec != nil && spec.Case == cc.DeclarationSpecifiersStorage &&
		spec.StorageClassSpecifier.Case == cc.StorageClassSpecifierTypedef {
		isTypedef = true
		spec = spec.DeclarationSpecifiers
	}
	for sp := spec; sp != nil; sp = sp.DeclarationSpecifiers {
		switch sp.Case {
		case cc.DeclarationSpecifiersTypeQual:
			ds := sp.TypeQualifier
			switch ds.Case {
			case cc.TypeQualifierConst:
				isConst = true
			case cc.TypeQualifierVolatile:
				isVolatile = true
			default:
				panic(ds.Case.String())
			}
		case cc.DeclarationSpecifiersStorage:
			ds := sp.StorageClassSpecifier
			switch ds.Case {
			case cc.StorageClassSpecifierStatic:
				isStatic = true
			case cc.StorageClassSpecifierExtern:
				isExtern = true
			case cc.StorageClassSpecifierAuto:
				isAuto = true
			case cc.StorageClassSpecifierRegister:
				// ignore
			default:
				panic(ds.Case.String())
			}
			if isTypedef {
				panic("wrong type")
			}
		case cc.DeclarationSpecifiersTypeSpec:
			ds := sp.TypeSpecifier
			switch ds.Case {
			case cc.TypeSpecifierStructOrUnion:
				su := ds.StructOrUnionSpecifier
				var conf IdentConfig
				for _, name := range names {
					if c, ok := g.idents[name]; ok {
						conf = c
						break
					}
				}
				switch su.Case {
				case cc.StructOrUnionSpecifierTag:
					// struct/union forward declaration
					isForward = true
					typeSpec = g.convertType(conf, su.Type(), d.Position()).(types.Named)
				case cc.StructOrUnionSpecifierDef:
					// struct/union declaration
					if isForward {
						panic("already marked as a forward decl")
					}
					typeSpec = g.convertType(conf, su.Type(), d.Position())
				default:
					panic(su.Case.String())
				}
			case cc.TypeSpecifierEnum:
				enumSpec = ds.EnumSpecifier
			case cc.TypeSpecifierVoid:
				isFunc = true
			default:
				isPrim = true
			}
		case cc.DeclarationSpecifiersFunc:
			isFunc = true
			// TODO: use specifiers
		default:
			panic(sp.Case.String() + " " + sp.Position().String())
		}
	}
	_ = isStatic // FIXME: static
	_ = isVolatile
	_ = isAuto // FIXME: auto
	var decls []CDecl
	if enumSpec != nil {
		if isForward {
			panic("TODO")
		}
		if isPrim || isFunc {
			panic("wrong type")
		}
		if typeSpec != nil {
			panic("should have no type")
		}
		if !inCur {
			return nil
		}
		var (
			typ           types.Type
			hasOtherDecls = false
		)
		if isTypedef {
			name, dd := g.convertTypedefName(d)
			und := g.convertType(IdentConfig{}, dd.Type(), name.Position())
			nt := g.newOrFindNamedType(name.Value.String(), func() types.Type {
				return und
			})
			typ = nt
			decls = append(decls, &CTypeDef{nt})
		} else if d.InitDeclaratorList != nil {
			hasOtherDecls = true
		} else if name := enumSpec.Token2; name.Value != 0 {
			nt := g.newOrFindNamedType(name.Value.String(), func() types.Type {
				return g.env.DefIntT()
			})
			typ = nt
			decls = append(decls, &CTypeDef{nt})
		}
		if !hasOtherDecls {
			decls = append(decls, g.convertEnum(d, typ, enumSpec)...)
		}
	}
	if d.InitDeclaratorList == nil || d.InitDeclaratorList.InitDeclarator == nil {
		if typeSpec == nil && enumSpec != nil {
			return decls
		}
		if isTypedef && isForward {
			panic("wrong type")
		}
		if isPrim || isFunc {
			panic("wrong type")
		}
		if isForward {
			if typeSpec == nil {
				panic("no type for forward decl")
			}
			if !inCur || !g.conf.ForwardDecl {
				return nil
			}
		} else {
			if !inCur {
				return nil
			}
			if isTypedef {
				panic("TODO")
			}
		}
		nt, ok := typeSpec.(types.Named)
		if !ok {
			if isForward {
				panic("forward declaration of unnamed type")
			} else {
				panic(fmt.Errorf("declaration of unnamed type: %T", typeSpec))
			}
		}
		decls = append(decls, &CTypeDef{nt})
		return decls
	}
	var (
		added   = 0
		skipped = 0
	)
	for il := d.InitDeclaratorList; il != nil; il = il.InitDeclaratorList {
		id := il.InitDeclarator
		switch id.Case {
		case cc.InitDeclaratorDecl, cc.InitDeclaratorInit:
			dd := id.Declarator
			dname := dd.Name().String()
			conf := g.idents[dname]
			vt := g.convertTypeRootOpt(conf, dd.Type(), id.Position())
			if isTypedef && vt == nil {
				vt = types.StructT(nil)
			}
			var init Expr
			if id.Initializer != nil && inCur {
				if isTypedef {
					panic("init in typedef: " + id.Position().String())
				}
				init = g.convertInitExpr(id.Initializer)
			}
			if isConst && propagateConst(vt) {
				isConst = false
			}
			if isTypedef {
				if enumSpec != nil {
					continue
				}
				nt, ok := vt.(types.Named)
				// TODO: this case is "too smart", we already handle those kind of double typedefs on a lower level
				if !ok || nt.Name().Name != dd.Name().String() {
					// we don't call a *From version of the method here because dd.Type() is an underlying type,
					// not a typedef type
					if ok && !strings.HasPrefix(nt.Name().Name, "_cxgo_") {
						decls = append(decls, &CTypeDef{nt})
					}
					if vt == nil {
						panic("TODO: typedef of void? " + id.Position().String())
					}
					nt = g.newOrFindNamedTypedef(dd.Name().String(), func() types.Type {
						return vt
					})
					if nt == nil {
						// typedef suppressed
						skipped++
						continue
					}
				}
				decls = append(decls, &CTypeDef{nt})
				continue
			}
			name := g.convertIdentWith(dd.NameTok().String(), vt, dd)
			isDecl := false
			for di := dd.DirectDeclarator; di != nil; di = di.DirectDeclarator {
				if di.Case == cc.DirectDeclaratorDecl {
					isDecl = true
					break
				}
			}
			if !isDecl && !isForward {
				if nt, ok := typeSpec.(types.Named); ok {
					decls = append(decls, &CTypeDef{nt})
				}
			}
			if ft, ok := vt.(*types.FuncType); ok && !isDecl {
				// forward declaration
				if l, id, ok := g.tenv.LibIdentByName(name.Name); ok && id.CType(nil).Kind().IsFunc() {
					// forward declaration of stdlib function
					// we must first load the corresponding library to the real env
					l, ok = g.env.GetLibrary(l.Name)
					if !ok {
						panic("cannot load stdlib")
					}
					id, ok = l.Idents[name.Name]
					if !ok {
						panic("cannot find stdlib ident")
					}
					g.replaceIdentWith(id, dd)
					skipped++
				} else if g.conf.ForwardDecl {
					decls = append(decls, &CFuncDecl{
						Name: name.Ident,
						Type: ft,
						Body: nil,
					})
				} else {
					skipped++
				}
			} else {
				decls = decls[:len(decls)-added]
				if !isExtern {
					var inits []Expr
					if init != nil {
						inits = []Expr{init}
					}
					decls = append(decls, &CVarDecl{
						// There is no real const in C
						Const: false, // Const: isConst,
						CVarSpec: CVarSpec{
							g:     g,
							Type:  vt,
							Names: []*types.Ident{name.Ident},
							Inits: inits,
						},
					})
				} else {
					skipped++
				}
			}
		default:
			panic(id.Case.String())
		}
	}
	if !inCur {
		return nil
	}
	if len(decls) == 0 && skipped == 0 {
		panic("no declarations converted: " + d.Position().String())
	}
	return decls
}

func (g *translator) convertCompStmt(d *cc.CompoundStatement) []CStmt {
	var stmts []CStmt
	for it := d.BlockItemList; it != nil; it = it.BlockItemList {
		st := it.BlockItem
		switch st.Case {
		case cc.BlockItemDecl:
			for _, dec := range g.convertDecl(st.Declaration) {
				stmts = append(stmts, g.NewCDeclStmt(dec)...)
			}
		case cc.BlockItemStmt:
			stmts = append(stmts, g.convertStmt(st.Statement)...)
		default:
			panic(st.Case.String())
		}
	}
	// TODO: shouldn't it return statements without a block? or call an optimizing version of block constructor?
	return []CStmt{g.newBlockStmt(stmts...)}
}

func (g *translator) convertCompBlockStmt(d *cc.CompoundStatement) *BlockStmt {
	stmts := g.convertCompStmt(d)
	if len(stmts) == 1 {
		if b, ok := stmts[0].(*BlockStmt); ok {
			return b
		}
	}
	// TODO: shouldn't it call a version that applies optimizations?
	return g.newBlockStmt(stmts...)
}

func (g *translator) convertExpr(d *cc.Expression) Expr {
	if d.Expression == nil {
		return g.convertAssignExpr(d.AssignmentExpression)
	}
	var exprs []*cc.AssignmentExpression
	for ; d != nil; d = d.Expression {
		exprs = append(exprs, d.AssignmentExpression)
	}
	var m []Expr
	for i := len(exprs) - 1; i >= 0; i-- {
		m = append(m, g.convertAssignExpr(exprs[i]))
	}
	return g.NewCMultiExpr(m...)
}

func (g *translator) convertExprOpt(d *cc.Expression) Expr {
	if d == nil {
		return nil
	}
	return g.convertExpr(d)
}

func (g *translator) convertMulExpr(d *cc.MultiplicativeExpression) Expr {
	switch d.Case {
	case cc.MultiplicativeExpressionCast:
		return g.convertCastExpr(d.CastExpression)
	}
	x := g.convertMulExpr(d.MultiplicativeExpression)
	y := g.convertCastExpr(d.CastExpression)
	var op BinaryOp
	switch d.Case {
	case cc.MultiplicativeExpressionMul:
		op = BinOpMult
	case cc.MultiplicativeExpressionDiv:
		op = BinOpDiv
	case cc.MultiplicativeExpressionMod:
		op = BinOpMod
	default:
		panic(d.Case.String())
	}
	return g.NewCBinaryExprT(
		x, op, y,
		g.convertTypeOper(d.Operand, d.Position()),
	)
}

func (g *translator) convertAddExpr(d *cc.AdditiveExpression) Expr {
	switch d.Case {
	case cc.AdditiveExpressionMul:
		return g.convertMulExpr(d.MultiplicativeExpression)
	}
	x := g.convertAddExpr(d.AdditiveExpression)
	y := g.convertMulExpr(d.MultiplicativeExpression)
	var op BinaryOp
	switch d.Case {
	case cc.AdditiveExpressionAdd:
		op = BinOpAdd
	case cc.AdditiveExpressionSub:
		op = BinOpSub
	default:
		panic(d.Case.String())
	}
	return g.NewCBinaryExprT(
		x, op, y,
		g.convertTypeOper(d.Operand, d.Position()),
	)
}

func (g *translator) convertShiftExpr(d *cc.ShiftExpression) Expr {
	switch d.Case {
	case cc.ShiftExpressionAdd:
		return g.convertAddExpr(d.AdditiveExpression)
	}
	x := g.convertShiftExpr(d.ShiftExpression)
	y := g.convertAddExpr(d.AdditiveExpression)
	var op BinaryOp
	switch d.Case {
	case cc.ShiftExpressionLsh:
		op = BinOpLsh
	case cc.ShiftExpressionRsh:
		op = BinOpRsh
	default:
		panic(d.Case.String())
	}
	return g.NewCBinaryExprT(
		x, op, y,
		g.convertTypeOper(d.Operand, d.Position()),
	)
}

func (g *translator) convertRelExpr(d *cc.RelationalExpression) Expr {
	switch d.Case {
	case cc.RelationalExpressionShift:
		return g.convertShiftExpr(d.ShiftExpression)
	}
	x := g.convertRelExpr(d.RelationalExpression)
	y := g.convertShiftExpr(d.ShiftExpression)
	var op ComparisonOp
	switch d.Case {
	case cc.RelationalExpressionLt:
		op = BinOpLt
	case cc.RelationalExpressionGt:
		op = BinOpGt
	case cc.RelationalExpressionLeq:
		op = BinOpLte
	case cc.RelationalExpressionGeq:
		op = BinOpGte
	default:
		panic(d.Case.String())
	}
	return g.Compare(x, op, y)
}

func (g *translator) convertEqExpr(d *cc.EqualityExpression) Expr {
	switch d.Case {
	case cc.EqualityExpressionRel:
		return g.convertRelExpr(d.RelationalExpression)
	}
	x := g.convertEqExpr(d.EqualityExpression)
	y := g.convertRelExpr(d.RelationalExpression)
	var op ComparisonOp
	switch d.Case {
	case cc.EqualityExpressionEq:
		op = BinOpEq
	case cc.EqualityExpressionNeq:
		op = BinOpNeq
	default:
		panic(d.Case.String())
	}
	return g.Compare(x, op, y)
}

func (g *translator) convertAndExpr(d *cc.AndExpression) Expr {
	switch d.Case {
	case cc.AndExpressionEq:
		return g.convertEqExpr(d.EqualityExpression)
	case cc.AndExpressionAnd:
		x := g.convertAndExpr(d.AndExpression)
		y := g.convertEqExpr(d.EqualityExpression)
		return g.NewCBinaryExprT(
			x, BinOpBitAnd, y,
			g.convertTypeOper(d.Operand, d.Position()),
		)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertLOrExcExpr(d *cc.ExclusiveOrExpression) Expr {
	switch d.Case {
	case cc.ExclusiveOrExpressionAnd:
		return g.convertAndExpr(d.AndExpression)
	case cc.ExclusiveOrExpressionXor:
		x := g.convertLOrExcExpr(d.ExclusiveOrExpression)
		y := g.convertAndExpr(d.AndExpression)
		return g.NewCBinaryExprT(
			x, BinOpBitXor, y,
			g.convertTypeOper(d.Operand, d.Position()),
		)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertLOrIncExpr(d *cc.InclusiveOrExpression) Expr {
	switch d.Case {
	case cc.InclusiveOrExpressionXor:
		return g.convertLOrExcExpr(d.ExclusiveOrExpression)
	case cc.InclusiveOrExpressionOr:
		x := g.convertLOrIncExpr(d.InclusiveOrExpression)
		y := g.convertLOrExcExpr(d.ExclusiveOrExpression)
		return g.NewCBinaryExprT(
			x, BinOpBitOr, y,
			g.convertTypeOper(d.Operand, d.Position()),
		)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertLAndExpr(d *cc.LogicalAndExpression) Expr {
	switch d.Case {
	case cc.LogicalAndExpressionOr:
		return g.convertLOrIncExpr(d.InclusiveOrExpression)
	case cc.LogicalAndExpressionLAnd:
		x := g.convertLAndExpr(d.LogicalAndExpression)
		y := g.convertLOrIncExpr(d.InclusiveOrExpression)
		return And(g.ToBool(x), g.ToBool(y))
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertLOrExpr(d *cc.LogicalOrExpression) Expr {
	switch d.Case {
	case cc.LogicalOrExpressionLAnd:
		return g.convertLAndExpr(d.LogicalAndExpression)
	case cc.LogicalOrExpressionLOr:
		x := g.convertLOrExpr(d.LogicalOrExpression)
		y := g.convertLAndExpr(d.LogicalAndExpression)
		return Or(g.ToBool(x), g.ToBool(y))
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertCondExpr(d *cc.ConditionalExpression) Expr {
	switch d.Case {
	case cc.ConditionalExpressionLOr:
		return g.convertLOrExpr(d.LogicalOrExpression)
	case cc.ConditionalExpressionCond:
		cond := g.convertLOrExpr(d.LogicalOrExpression)
		return g.NewCTernaryExpr(
			g.ToBool(cond),
			g.convertExpr(d.Expression),
			g.convertCondExpr(d.ConditionalExpression),
		)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertPriExpr(d *cc.PrimaryExpression) Expr {
	switch d.Case {
	case cc.PrimaryExpressionIdent: // x
		if d.Token.String() == "asm" {
			return &CAsmExpr{e: g.env.Env}
		}
		if d.Operand == nil {
			panic(ErrorfWithPos(d.Position(), "empty operand for %q", d.Token.String()))
		}
		return g.convertIdent(d.ResolvedIn(), d.Token, g.convertTypeOper(d.Operand, d.Position()))
	case cc.PrimaryExpressionEnum: // X
		return g.convertIdent(d.ResolvedIn(), d.Token, g.convertTypeOper(d.Operand, d.Position()))
	case cc.PrimaryExpressionInt: // 1
		fnc := func() Expr {
			v, err := parseCIntLit(d.Token.String(), g.conf.IntReformat)
			if err != nil {
				panic(err)
			}
			return v
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionFloat: // 0.0
		fnc := func() Expr {
			v, err := parseCFloatLit(d.Token.String())
			if err != nil {
				panic(err)
			}
			if d.Operand == nil {
				return v
			}
			return g.cCast(
				g.convertTypeOper(d.Operand, d.Position()),
				v,
			)
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionChar: // 'x'
		fnc := func() Expr {
			return cLitT(
				d.Token.String(), CLitChar,
				g.convertTypeOper(d.Operand, d.Position()),
			)
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionLChar: // 'x'
		fnc := func() Expr {
			return cLitT(
				d.Token.String(), CLitWChar,
				g.convertTypeOper(d.Operand, d.Position()),
			)
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionString: // "x"
		fnc := func() Expr {
			v, err := g.parseCStringLit(d.Token.String())
			if err != nil {
				panic(err)
			}
			return v
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionLString: // L"x"
		fnc := func() Expr {
			v, err := g.parseCWStringLit(d.Token.String())
			if err != nil {
				panic(err)
			}
			return v
		}
		if m := d.Token.Macro(); m != 0 {
			return g.convMacro(m.String(), fnc)
		}
		return fnc()
	case cc.PrimaryExpressionExpr: // "(x)"
		e := g.convertExpr(d.Expression)
		return cParen(e)
	case cc.PrimaryExpressionStmt: // "({...; x})"
		stmt := g.convertCompStmt(d.CompoundStatement)
		if len(stmt) != 1 {
			panic("TODO")
		}
		stmt = stmt[0].(*BlockStmt).Stmts
		last, ok := stmt[len(stmt)-1].(*CExprStmt)
		if !ok {
			// let it cause a compilation error in Go
			return &CallExpr{
				Fun: g.NewFuncLit(g.env.FuncTT(g.env.DefIntT()), stmt...),
			}
		}
		typ := last.Expr.CType(nil)
		stmt = append(stmt[:len(stmt)-1], g.NewReturnStmt(last.Expr, typ)...)
		return &CallExpr{
			Fun: g.NewFuncLit(g.env.FuncTT(typ), stmt...),
		}
	default:
		panic(fmt.Errorf("%v (%v)", d.Case, d.Position()))
	}
}

func (g *translator) convertOneDesignator(typ types.Type, list *cc.DesignatorList, val Expr) *CompLitField {
	d := list.Designator
	var (
		f   *CompLitField
		sub types.Type
	)
	switch d.Case {
	case cc.DesignatorIndex:
		f = &CompLitField{Index: g.convertConstExpr(d.ConstantExpression)}
		sub = typ.(types.ArrayType).Elem()
	case cc.DesignatorField:
		f = &CompLitField{Field: g.convertIdentOn(typ, d.Token2)}
		sub = f.Field.CType(nil)
	case cc.DesignatorField2:
		f = &CompLitField{Field: g.convertIdentOn(typ, d.Token)}
		sub = f.Field.CType(nil)
	default:
		panic(d.Case.String() + " " + d.Position().String())
	}
	if list.DesignatorList == nil {
		f.Value = val
		return f
	}
	f2 := g.convertOneDesignator(sub, list.DesignatorList, val)
	f.Value = g.NewCCompLitExpr(sub, []*CompLitField{f2})
	return f
}

func (g *translator) convertPostfixExpr(d *cc.PostfixExpression) Expr {
	switch d.Case {
	case cc.PostfixExpressionPrimary:
		return g.convertPriExpr(d.PrimaryExpression)
	case cc.PostfixExpressionIndex: // "x[y]"
		return g.NewCIndexExpr(
			g.convertPostfixExpr(d.PostfixExpression),
			g.convertExpr(d.Expression),
			g.convertTypeOper(d.Operand, d.Position()),
		)
	case cc.PostfixExpressionCall: // x([args])
		fnc := g.convertPostfixExpr(d.PostfixExpression)
		var args []Expr
		for it := d.ArgumentExpressionList; it != nil; it = it.ArgumentExpressionList {
			args = append(args, g.convertAssignExpr(it.AssignmentExpression))
		}
		return g.NewCCallExpr(g.ToFunc(fnc, ToFuncExpr(fnc.CType(nil))), args)
	case cc.PostfixExpressionPSelect: // x->y
		exp := g.convertPostfixExpr(d.PostfixExpression)
		if _, ok := exp.CType(nil).(types.ArrayType); ok { // pointer accesses might be an array
			return NewCSelectExpr(
				g.NewCIndexExpr(
					exp,
					cUintLit(0, 10), // index the first element
					g.convertTypeOper(d.Operand, d.Position()),
				), g.convertIdentOn(exp.CType(nil), d.Token2),
			)
		}
		return NewCSelectExpr(
			exp, g.convertIdentOn(exp.CType(nil), d.Token2),
		)
	case cc.PostfixExpressionSelect: // x.y
		exp := g.convertPostfixExpr(d.PostfixExpression)
		return NewCSelectExpr(
			exp, g.convertIdentOn(exp.CType(nil), d.Token2),
		)
	case cc.PostfixExpressionInc: // x++
		x := g.convertPostfixExpr(d.PostfixExpression)
		return g.NewCPostfixExpr(x, false)
	case cc.PostfixExpressionDec: // x--
		x := g.convertPostfixExpr(d.PostfixExpression)
		return g.NewCPostfixExpr(x, true)
	case cc.PostfixExpressionComplit:
		return g.convertInitList(
			g.convertType(IdentConfig{}, d.TypeName.Type(), d.Position()),
			d.InitializerList,
		)
	default:
		panic(d.Case.String() + " " + d.Position().String())
	}
}

func (g *translator) convertCastExpr(d *cc.CastExpression) Expr {
	switch d.Case {
	case cc.CastExpressionUnary:
		return g.convertUnaryExpr(d.UnaryExpression)
	case cc.CastExpressionCast:
		x := g.convertCastExpr(d.CastExpression)
		if k := d.Operand.Type().Kind(); k == cc.Invalid || k == cc.Void {
			return x
		}
		return g.cCast(
			g.convertTypeOper(d.Operand, d.Position()),
			x,
		)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertUnaryExpr(d *cc.UnaryExpression) Expr {
	switch d.Case {
	case cc.UnaryExpressionPostfix:
		return g.convertPostfixExpr(d.PostfixExpression)
	case cc.UnaryExpressionInc: // ++x
		x := g.convertUnaryExpr(d.UnaryExpression)
		return g.NewCPrefixExpr(x, false)
	case cc.UnaryExpressionDec: // --x
		x := g.convertUnaryExpr(d.UnaryExpression)
		return g.NewCPrefixExpr(x, true)
	case cc.UnaryExpressionSizeofExpr: // sizeof x
		return g.NewCUnaryExprT(
			UnarySizeof,
			g.convertUnaryExpr(d.UnaryExpression),
			g.convertTypeOper(d.Operand, d.Position()),
		)
	case cc.UnaryExpressionSizeofType: // sizeof tp
		return g.SizeofT(
			g.convertType(IdentConfig{}, d.TypeName.Type(), d.Position()),
			nil,
		)
	case cc.UnaryExpressionAlignofType: // alignof tp
		return g.AlignofT(
			g.convertType(IdentConfig{}, d.TypeName.Type(), d.Position()),
			nil,
		)
	}
	var op UnaryOp
	switch d.Case {
	case cc.UnaryExpressionAddrof: // &x
		x := g.convertCastExpr(d.CastExpression)
		return g.cAddr(x)
	case cc.UnaryExpressionDeref: // *x
		x := g.convertCastExpr(d.CastExpression)
		typ := g.convertTypeOper(d.Operand, d.Position())
		return g.cDerefT(x, typ)
	case cc.UnaryExpressionPlus: // +x
		op = UnaryPlus
	case cc.UnaryExpressionMinus: // -x
		op = UnaryMinus
	case cc.UnaryExpressionCpl: // ~x
		op = UnaryXor
	case cc.UnaryExpressionNot: // !x
		x := g.convertCastExpr(d.CastExpression)
		return g.cNot(x)
	default:
		panic(d.Case.String())
	}
	x := g.convertCastExpr(d.CastExpression)
	if d.Operand == nil {
		return g.NewCUnaryExpr(
			op, x,
		)
	}
	return g.NewCUnaryExprT(
		op, x,
		g.convertTypeOper(d.Operand, d.Position()),
	)
}

func (g *translator) convertConstExpr(d *cc.ConstantExpression) Expr {
	return g.convertCondExpr(d.ConditionalExpression)
}

func (g *translator) convertAssignExpr(d *cc.AssignmentExpression) Expr {
	switch d.Case {
	case cc.AssignmentExpressionCond:
		return g.convertCondExpr(d.ConditionalExpression)
	}
	x := g.convertUnaryExpr(d.UnaryExpression)
	y := g.convertAssignExpr(d.AssignmentExpression)
	var op BinaryOp
	switch d.Case {
	case cc.AssignmentExpressionAssign:
		op = ""
	case cc.AssignmentExpressionMul:
		op = BinOpMult
	case cc.AssignmentExpressionDiv:
		op = BinOpDiv
	case cc.AssignmentExpressionMod:
		op = BinOpMod
	case cc.AssignmentExpressionAdd:
		op = BinOpAdd
	case cc.AssignmentExpressionSub:
		op = BinOpSub
	case cc.AssignmentExpressionLsh:
		op = BinOpLsh
	case cc.AssignmentExpressionRsh:
		op = BinOpRsh
	case cc.AssignmentExpressionAnd:
		op = BinOpBitAnd
	case cc.AssignmentExpressionXor:
		op = BinOpBitXor
	case cc.AssignmentExpressionOr:
		op = BinOpBitOr
	default:
		panic(d.Case.String())
	}
	return g.NewCAssignExpr(
		x, op, y,
	)
}

func (g *translator) convertLabelStmt(st *cc.LabeledStatement) []CStmt {
	switch st.Case {
	case cc.LabeledStatementLabel: // label:
		stmts := g.convertStmt(st.Statement)
		return append([]CStmt{
			&CLabelStmt{Label: st.Token.Value.String()},
		}, stmts...)
	case cc.LabeledStatementCaseLabel: // case xxx:
		return []CStmt{
			g.NewCaseStmt(
				g.convertConstExpr(st.ConstantExpression),
				g.convertStmt(st.Statement)...,
			),
		}
	case cc.LabeledStatementDefault: // default:
		return []CStmt{
			g.NewCaseStmt(
				nil,
				g.convertStmt(st.Statement)...,
			),
		}
	default:
		panic(st.Case.String())
	}
}

func (g *translator) convertExprStmt(st *cc.ExpressionStatement) []CStmt {
	var exprs []*cc.AssignmentExpression
	for e := st.Expression; e != nil; e = e.Expression {
		exprs = append(exprs, e.AssignmentExpression)
	}
	var stmts []CStmt
	for i := len(exprs) - 1; i >= 0; i-- {
		stmts = append(stmts, NewCExprStmt(g.convertAssignExpr(exprs[i]))...)
	}
	return stmts
}

func (g *translator) convertSelStmt(st *cc.SelectionStatement) []CStmt {
	switch st.Case {
	case cc.SelectionStatementIf: // if (x)
		cond := g.convertExpr(st.Expression)
		return []CStmt{
			g.NewCIfStmt(
				g.ToBool(cond),
				[]CStmt{g.convertBlockStmt(st.Statement)},
				nil,
			),
		}
	case cc.SelectionStatementIfElse: // if (x) else
		cond := g.convertExpr(st.Expression)
		return []CStmt{
			g.NewCIfStmt(
				g.ToBool(cond),
				[]CStmt{g.convertBlockStmt(st.Statement)},
				g.toElseStmt(g.convertOneStmt(st.Statement2)),
			),
		}
	case cc.SelectionStatementSwitch: // switch (x)
		return []CStmt{g.NewCSwitchStmt(
			g.convertExpr(st.Expression),
			[]CStmt{g.convertBlockStmt(st.Statement)},
		)}
	default:
		panic(st.Case.String())
	}
}

func (g *translator) convertIterStmt(st *cc.IterationStatement) []CStmt {
	switch st.Case {
	case cc.IterationStatementWhile:
		x := g.convertExprOpt(st.Expression)
		var cond BoolExpr
		if x != nil {
			cond = g.ToBool(x)
		}
		return []CStmt{
			g.NewCForStmt(
				nil,
				cond,
				nil,
				[]CStmt{g.convertBlockStmt(st.Statement)},
			),
		}
	case cc.IterationStatementDo:
		return []CStmt{
			g.NewCDoWhileStmt(
				g.convertExprOpt(st.Expression),
				[]CStmt{g.convertBlockStmt(st.Statement)},
			),
		}
	case cc.IterationStatementFor:
		x := g.convertExprOpt(st.Expression2)
		var cond BoolExpr
		if x != nil {
			cond = g.ToBool(x)
		}
		return []CStmt{
			g.NewCForStmt(
				g.convertExprOpt(st.Expression),
				cond,
				g.convertExprOpt(st.Expression3),
				[]CStmt{g.convertBlockStmt(st.Statement)},
			),
		}
	case cc.IterationStatementForDecl:
		var cur *CVarDecl
		for _, d := range g.convertDecl(st.Declaration) {
			d := d.(*CVarDecl)
			if cur == nil {
				cur = d
				continue
			}
			if !types.Same(cur.Type, d.Type) {
				panic(fmt.Errorf("different types in a declaration: %v vs %v (%s)", cur.Type, d.Type, st.Position()))
			}
			cur.Single = true
			n1, n2 := len(cur.Names), len(d.Names)
			cur.Names = append(cur.Names, d.Names...)
			if len(cur.Inits) == 0 && len(d.Inits) == 0 {
				continue
			}
			if len(cur.Inits) == 0 {
				cur.Inits = make([]Expr, n1, n1+n2)
			}
			if len(d.Inits) == 0 {
				cur.Inits = append(cur.Inits, make([]Expr, n2)...)
			} else {
				cur.Inits = append(cur.Inits, d.Inits...)
			}
		}
		x := g.convertExprOpt(st.Expression)
		var cond BoolExpr
		if x != nil {
			cond = g.ToBool(x)
		}
		return []CStmt{
			g.NewCForDeclStmt(
				cur,
				cond,
				g.convertExprOpt(st.Expression2),
				[]CStmt{g.convertBlockStmt(st.Statement)},
			),
		}
	default:
		panic(st.Case.String() + " " + st.Position().String())
	}
}

func (g *translator) convertJumpStmt(st *cc.JumpStatement) []CStmt {
	switch st.Case {
	case cc.JumpStatementGoto: // goto x
		return []CStmt{
			&CGotoStmt{Label: st.Token2.Value.String()},
		}
	case cc.JumpStatementContinue: // continue
		return []CStmt{
			&CContinueStmt{},
		}
	case cc.JumpStatementBreak: // break
		return []CStmt{
			&CBreakStmt{},
		}
	case cc.JumpStatementReturn: // return
		return g.NewReturnStmt(
			g.convertExprOpt(st.Expression),
			nil,
		)
	default:
		panic(st.Case.String())
	}
}

func (g *translator) convertAsmStmt(d *cc.AsmStatement) []CStmt {
	// TODO
	return NewCExprStmt(&CAsmExpr{e: g.env.Env, typ: types.UnkT(1)})
}

func (g *translator) convertStmt(d *cc.Statement) []CStmt {
	switch d.Case {
	case cc.StatementLabeled:
		return g.convertLabelStmt(d.LabeledStatement)
	case cc.StatementCompound:
		return g.convertCompStmt(d.CompoundStatement)
	case cc.StatementExpr:
		return g.convertExprStmt(d.ExpressionStatement)
	case cc.StatementSelection:
		return g.convertSelStmt(d.SelectionStatement)
	case cc.StatementIteration:
		return g.convertIterStmt(d.IterationStatement)
	case cc.StatementJump:
		return g.convertJumpStmt(d.JumpStatement)
	case cc.StatementAsm:
		return g.convertAsmStmt(d.AsmStatement)
	default:
		panic(d.Case.String())
	}
}

func (g *translator) convertOneStmt(d *cc.Statement) CStmt {
	stmts := g.convertStmt(d)
	if len(stmts) == 1 {
		return stmts[0]
	}
	return g.NewCBlock(stmts...)
}

func (g *translator) convertBlockStmt(d *cc.Statement) *BlockStmt {
	stmts := g.convertStmt(d)
	if len(stmts) == 1 {
		if b, ok := stmts[0].(*BlockStmt); ok {
			return b
		}
	}
	return g.NewCBlock(stmts...)
}
