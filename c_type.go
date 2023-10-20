package cxgo

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"modernc.org/cc/v3"
	"modernc.org/token"

	"github.com/gotranspile/cxgo/types"
)

func (g *translator) convertTypeOper(p cc.Operand, where token.Position) types.Type {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case error:
			panic(ErrorWithPos(r, where))
		default:
			panic(ErrorWithPos(fmt.Errorf("%v", r), where))
		}
	}()
	if d := p.Declarator(); d != nil {
		where = d.Position()
	}
	var conf IdentConfig
	if d := p.Declarator(); d != nil {
		conf = g.idents[d.Name().String()]
	}
	return g.convertTypeRoot(conf, p.Type(), where)
}

// convertType is similar to newTypeCC but it will first consult the type cache.
func (g *translator) convertType(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("type conversion failed at %v: %v", where, r))
		}
	}()
	// custom type overrides coming from the config
	// note that we don't save them since they might depend
	// not only on the input type, but also on a field name
	switch conf.Type {
	case HintBool:
		return g.env.Go().Bool()
	case HintIface:
		return g.env.Go().Any()
	case HintString:
		return g.env.Go().String()
	case HintSlice:
		ct := g.newTypeCC(IdentConfig{}, t, where)
		var elem types.Type
		switch ct := ct.(type) {
		case types.PtrType:
			elem = ct.Elem()
		case types.ArrayType:
			elem = ct.Elem()
		default:
			panic(fmt.Errorf("expected an array or a pointer, got: %v, %#v; defined at: %v", ct, ct, where))
		}
		if elem == types.UintT(1) {
			elem = g.env.Go().Byte()
		}
		return types.SliceT(elem)
	}
	// allow invalid types, they might still be useful
	// since one may define them in a separate Go file
	// and make the code valid
	if t.Kind() == cc.Invalid {
		return types.UnkT(g.env.PtrSize())
	}
	if ct, ok := g.ctypes[t]; ok {
		return ct
	}
	ct := g.newTypeCC(conf, t, where)
	g.ctypes[t] = ct
	return ct
}

// convertTypeRoot is the same as convertType, but it applies a workaround for
// C function pointers.
func (g *translator) convertTypeRoot(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	ft := g.convertType(conf, t, where)
	if p, ok := ft.(types.PtrType); ok && p.ElemKind().IsFunc() {
		ft = p.Elem()
	}
	if p, ok := ft.(types.ArrayType); ok && p.Len() == 0 {
		ft = types.SliceT(p.Elem())
	}
	return ft
}

// convertTypeOpt is similar to convertType, but it also allows void type by returning nil.
func (g *translator) convertTypeOpt(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	if t == nil || t.Kind() == cc.Void || t.Kind() == cc.Invalid {
		return nil
	}
	return g.convertType(conf, t, where)
}

// convertTypeRootOpt is similar to convertTypeRoot, but it also allows void type by returning nil.
func (g *translator) convertTypeRootOpt(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	if t == nil || t.Kind() == cc.Void || t.Kind() == cc.Invalid {
		return nil
	}
	return g.convertTypeRoot(conf, t, where)
}

// replaceType checks if the type needs to be replaced. It usually happens for builtin types.
func (g *translator) replaceType(name string) (types.Type, bool) {
	if t, ok := g.env.TypeByName(name); ok {
		return t, true
	}
	if t := g.env.C().Type(name); t != nil {
		return t, true
	}
	return nil, false
}

// newNamedTypeAt finds or creates a named type defined by specified CC types and tokens.
func (g *translator) newNamedTypeAt(name string, typ, elem cc.Type, where token.Position) types.Type {
	if typ == elem {
		switch typ.Kind() {
		case cc.Struct, cc.Union:
		default:
			panic(fmt.Errorf("name: %s, elem: (%T) %v", name, elem, elem))
		}
	}
	conf := g.idents[name]
	if typ, ok := g.replaceType(name); ok {
		return typ
	}
	if c, ok := g.idents[name]; ok && c.Alias {
		sub := g.convertTypeRoot(conf, elem, where)
		g.ctypes[typ] = sub
		g.aliases[name] = sub
		return sub
	}
	return g.newOrFindNamedType(name, func() types.Type {
		return g.convertTypeRoot(conf, elem, where)
	})
}

func (g *translator) newOrFindNamedTypedef(name string, underlying func() types.Type) types.Named {
	if c, ok := g.idents[name]; ok && c.Alias {
		if _, ok := g.aliases[name]; ok {
			return nil
		}
		// we should register the underlying type with the current name,
		// so all the accesses will use underlying type
		t := underlying()
		g.aliases[name] = t
		// and we suppress the definition of this type
		return nil
	}
	return g.newOrFindNamedType(name, underlying)
}

// newOrFindNamedType finds or creates a new named type with a given underlying type.
// The function is given because types may be recursive.
func (g *translator) newOrFindNamedType(name string, underlying func() types.Type) types.Named {
	if _, ok := g.aliases[name]; ok {
		panic("alias")
	}
	if typ, ok := g.named[name]; ok {
		return typ
	}
	und := underlying()
	if typ, ok := g.named[name]; ok {
		return typ
	}
	return g.newNamedType(name, und)
}

// newNamedType creates a new named type based on the underlying type.
func (g *translator) newNamedType(name string, underlying types.Type) types.Named {
	if _, ok := g.named[name]; ok {
		panic("type with a same name already exists: " + name)
	}
	goname := ""
	if c, ok := g.idents[name]; ok && c.Rename != "" {
		goname = c.Rename
	}
	nt := types.NamedTGo(name, goname, underlying)
	g.named[name] = nt
	return nt
}

// newNamedTypeFrom creates a new named type based on a given CC type.
// It is similar to newNamedType, but accepts a CC type that should be bound to a new type.
func (g *translator) newNamedTypeFrom(name string, underlying types.Type, from cc.Type) types.Named {
	if _, ok := g.ctypes[from]; ok {
		panic("same C type already exists")
	}
	nt := g.newNamedType(name, underlying)
	g.ctypes[from] = nt
	return nt
}

// newOrFindNamedTypeFrom finds or creates a type with a given name, underlying type and source C type.
// It cannot return the NamedType because the type may have an override.
func (g *translator) newOrFindNamedTypeFrom(name string, elem func() types.Type, from cc.Type) types.Type {
	if t, ok := g.ctypes[from]; ok {
		return t
	}
	nt := g.newOrFindNamedType(name, elem)
	g.ctypes[from] = nt
	return nt
}

// newOrFindIncompleteNamedTypeFrom finds or creates a incomplete type with a given name and source C type.
// It cannot return the NamedType because the type may have an override, or the type may have been resolved to something else.
func (g *translator) newOrFindIncompleteNamedTypeFrom(name string, from cc.Type) types.Type {
	return g.newOrFindNamedTypeFrom(name, nil, from)
}

func asExportedName(s string) string {
	if len(s) == 0 {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// newTypeCC creates a new type based on a specified C type. This function will not consult the cache for a given type.
// It will recursively convert all the underlying and sub-types using convertType.
func (g *translator) newTypeCC(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	sname := t.Name().String()
	if nt, ok := g.replaceType(sname); ok {
		g.ctypes[t] = nt
		return nt
	}
	// it's handled separately because it's the only type that is allowed to be incomplete
	if t != t.Alias() {
		if t.IsIncomplete() {
			if nt, ok := g.named[sname]; ok {
				return nt
			}
			return g.newOrFindNamedType(sname, func() types.Type {
				return types.StructT(nil)
			})
		}
		if t, ok := g.aliases[sname]; ok {
			return t
		}
		conf := g.idents[sname]
		u := t.Alias()
		sub := g.convertType(conf, u, where)
		if sub, ok := sub.(types.Named); ok {
			if name := t.Name(); name == u.Name() {
				g.named[name.String()] = sub
				g.ctypes[t] = sub
				return sub
			}
		}
		return g.newNamedTypeAt(sname, t, u, where)
	}
	switch t.Kind() {
	case cc.Struct, cc.Union:
		return g.convertStructType(conf, t, where)
	}
	if u := t.Alias(); t != u {
		panic(fmt.Errorf("unhandled alias type: %T", t))
	}
	switch kind := t.Kind(); kind {
	case cc.UInt64, cc.UInt32, cc.UInt16, cc.UInt8:
		return types.UintT(int(t.Size()))
	case cc.Int64, cc.Int32, cc.Int16, cc.Int8:
		return types.IntT(int(t.Size()))
	case cc.SChar:
		return g.env.C().SignedChar()
	case cc.UChar:
		return g.env.C().UnsignedChar()
	case cc.Short:
		return g.env.C().Short()
	case cc.UShort:
		return g.env.C().UnsignedShort()
	case cc.Int:
		return g.env.C().Int()
	case cc.UInt:
		return g.env.C().UnsignedInt()
	case cc.Long:
		return g.env.C().Long()
	case cc.ULong:
		return g.env.C().UnsignedLong()
	case cc.LongLong:
		return g.env.C().LongLong()
	case cc.ULongLong:
		return g.env.C().UnsignedLongLong()
	case cc.Float:
		return g.env.C().Float()
	case cc.Double:
		return g.env.C().Double()
	case cc.LongDouble:
		return types.FloatT(int(t.Size()))
	case cc.Char:
		return g.env.C().Char()
	case cc.Bool:
		return g.env.C().Bool()
	case cc.Function:
		return g.convertFuncType(conf, nil, t, where)
	case cc.Ptr:
		if t.Elem().Kind() == cc.Char {
			return g.env.C().String()
		}
		if e := t.Elem(); e.Kind() == cc.Struct && e.NumField() == 1 {
			// Go slices defined via cxgo builtins
			f := e.FieldByIndex([]int{0})
			if f.Name().String() == types.GoPrefix+"slice_data" {
				elem := g.convertType(IdentConfig{}, f.Type(), where)
				return types.SliceT(elem)
			}
		}
		var ptr types.PtrType
		if name := t.Elem().Name(); name != 0 {
			if pt, ok := g.namedPtrs[name.String()]; ok {
				return pt
			}
			ptr = g.env.PtrT(nil) // incomplete
			g.namedPtrs[name.String()] = ptr
		}
		// use Opt because of the void*
		elem := g.convertTypeOpt(IdentConfig{}, t.Elem(), where)
		if ptr != nil {
			ptr.SetElem(elem)
			return ptr
		}
		return g.env.PtrT(elem)
	case cc.Array:
		if t.Elem().Kind() == cc.Char {
			return types.ArrayT(g.env.Go().Byte(), int(t.Len()))
		}
		elem := g.convertType(IdentConfig{}, t.Elem(), where)
		return types.ArrayT(
			elem,
			int(t.Len()),
		)
	case cc.Union:
		if name := t.Name(); name != 0 {
			u := t.Alias()
			if name == u.Name() {
				u = u.Alias()
			}
			return g.newNamedTypeAt(name.String(), t, u, where)
		}
		fconf := make(map[string]IdentConfig)
		for _, f := range conf.Fields {
			fconf[f.Name] = f
		}
		var fields []*types.Field
		for i := 0; i < t.NumField(); i++ {
			f := t.FieldByIndex([]int{i})
			name := f.Name().String()
			fc := fconf[name]
			ft := g.convertTypeRoot(fc, f.Type(), where)
			fields = append(fields, &types.Field{
				Name: g.newIdent(name, ft),
			})
		}
		return types.UnionT(fields)
	case cc.Enum:
		return g.newTypeCC(IdentConfig{}, t.EnumType(), where)
	default:
		panic(fmt.Errorf("%T, %s (%s)", t, kind, t.String()))
	}
}

func (g *translator) convertStructType(conf IdentConfig, t cc.Type, where token.Position) types.Type {
	sname := t.Name().String()
	if c, ok := g.idents[sname]; ok {
		conf = c
	}
	fconf := make(map[string]IdentConfig)
	for _, f := range conf.Fields {
		fconf[f.Name] = f
	}
	buildType := func() types.Type {
		var fields []*types.Field
		for i := 0; i < t.NumField(); i++ {
			f := t.FieldByIndex([]int{i})
			fc := fconf[f.Name().String()]
			ft := g.convertTypeRoot(fc, f.Type(), where)
			if f.Name() == 0 {
				st := types.Unwrap(ft).(*types.StructType)
				fields = append(fields, st.Fields()...)
				continue
			}
			fname := g.newIdent(f.Name().String(), ft)
			if fc.Rename != "" {
				fname.GoName = fc.Rename
			} else if !g.conf.UnexportedFields {
				fname.GoName = asExportedName(fname.Name)
			}
			fields = append(fields, &types.Field{
				Name: fname,
			})
		}
		if !where.IsValid() {
			panic(where)
		}
		var s *types.StructType
		if t.Kind() == cc.Union {
			s = types.UnionT(fields)
		} else {
			s = types.StructT(fields)
		}
		s.Where = where.String()
		if t.Name() == 0 {
			return s
		}
		return s
	}
	if t.Name() == 0 {
		return buildType()
	}
	return g.newOrFindNamedType(sname, buildType)
}

func (g *translator) convertFuncType(conf IdentConfig, d *cc.Declarator, t cc.Type, where token.Position) *types.FuncType {
	if kind := t.Kind(); kind != cc.Function {
		panic(kind)
	}
	if d != nil {
		where = d.Position()
	}
	var rconf IdentConfig
	aconf := make(map[string]IdentConfig)
	iconf := make(map[int]IdentConfig)
	for _, f := range conf.Fields {
		if f.Name != "" {
			if f.Name == "return" {
				rconf = f
			} else {
				aconf[f.Name] = f
			}
		} else {
			iconf[f.Index] = f
		}
	}
	var (
		args  []*types.Field
		named int
	)
	for i, p := range t.Parameters() {
		pt := p.Type()
		if pt.Kind() == cc.Void {
			continue
		}
		var fc IdentConfig
		if ac, ok := aconf[p.Name().String()]; ok {
			fc = ac
		} else if ac, ok = iconf[i]; ok {
			fc = ac
		}
		at := g.convertTypeRoot(fc, pt, where)
		var name *types.Ident
		if d != nil && p.Name() != 0 {
			name = g.convertIdent(d.ParamScope(), p.Declarator().NameTok(), at).Ident
			named++
		} else if p.Name() != 0 {
			name = g.convertIdentWith(p.Declarator().NameTok().String(), at, p.Declarator()).Ident
			named++
		} else {
			name = types.NewUnnamed(at)
		}
		args = append(args, &types.Field{
			Name: name,
		})
	}
	if named != 0 && len(args) != named {
		for i, a := range args {
			if a.Name.Name == "" && a.Name.GoName == "" {
				a.Name.GoName = fmt.Sprintf("a%d", i+1)
			}
		}
	}
	ret := g.convertTypeRootOpt(rconf, t.Result(), where)
	if t.IsVariadic() {
		return g.env.VarFuncT(ret, args...)
	}
	return g.env.FuncT(ret, args...)
}

func propagateConst(t types.Type) bool {
	switch t := t.(type) {
	case types.PtrType:
		if !propagateConst(t.Elem()) {
			//t.Const = true // TODO
		}
		return true
	case types.ArrayType:
		return propagateConst(t.Elem())
	case types.IntType:
		//t.Const = true
		return true
	}
	return false
}

func (g *translator) ZeroValue(t types.Type) Expr {
	if t == nil {
		panic("nil type")
	}
	switch t.Kind().Major() {
	case types.Ptr, types.Func:
		return g.Nil()
	case types.Int, types.Float:
		return cUintLit(0, 10)
	case types.Struct, types.Array:
		return &CCompLitExpr{Type: t}
	default:
		panic(t)
	}
}
