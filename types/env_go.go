package types

import "fmt"

const GoPrefix = "_cxgo_go_"

type Go struct {
	e   *Env
	pkg *Package

	// don't forget to update g.Types() when adding new types here

	boolT    Type
	byteT    Type
	runeT    Type
	uintptrT Type
	intT     Type
	uintT    Type
	ifaceT   Type
	ifaceSlT Type
	stringT  Type
	bytesT   Type

	iot     *Ident
	lenF    *Ident
	copyF   *Ident
	panicF  *Ident
	osExitF *Ident
}

// Go returns a package containing builtin Go types.
func (e *Env) Go() *Go {
	return &e.g
}

func (e *Env) initGo() {
	e.g.e = e
	e.g.pkg = e.newPackage("", "")
	e.g.init()
}

func (g *Go) Types() []Type {
	return []Type{
		g.boolT,
		g.byteT,
		g.runeT,
		g.uintptrT,
		g.intT,
		g.uintT,
		g.ifaceT,
		g.ifaceSlT,
		g.stringT,
		g.bytesT,
	}
}

func (g *Go) IsBuiltinType(t Type) bool {
	for _, t2 := range g.Types() {
		if t == t2 {
			return true
		}
	}
	return false
}

func (g *Go) init() {
	const cpref = GoPrefix

	// Those are native Go types that will always be mapped to themselves when transpiling
	// All other types might be mapped differently from C to Go

	// register basic Go types
	g.boolT = g.pkg.NewAlias(cpref+"bool", "bool", BoolT())
	g.byteT = g.pkg.NewTypeGo(cpref+"byte", "byte", UintT(1))
	g.runeT = g.pkg.NewTypeGo(cpref+"rune", "rune", IntT(4))

	ptrSize := g.e.PtrSize()
	// TODO: we are assuming Go arch = C arch here
	g.uintptrT = g.pkg.NewTypeGo(cpref+"uintptr", "uintptr", UintT(ptrSize))
	g.intT = g.pkg.NewTypeGo(cpref+"int", "int", IntT(g.e.conf.IntSize))
	g.uintT = g.pkg.NewTypeGo(cpref+"uint", "uint", UintT(g.e.conf.IntSize))
	g.stringT = g.pkg.NewTypeGo(cpref+"string", "string", UnkT(ptrSize*3))
	g.ifaceT = g.pkg.NewTypeGo(cpref+"iface", "interface{}", UnkT(ptrSize*2))

	// register well-know slice types
	g.bytesT = g.pkg.NewTypeGo(cpref+"bytes", "[]byte", UnkT(ptrSize*3))
	g.ifaceSlT = g.pkg.NewTypeGo(cpref+"iface_slice", "[]interface{}", UnkT(ptrSize*3))

	// register fixed-size builtin Go types
	for _, sz := range []int{
		1, 2, 4, 8,
	} {
		name := fmt.Sprintf("int%d", sz*8)
		g.pkg.NewAlias(cpref+name, name, IntT(sz))          // intN
		g.pkg.NewAlias(cpref+"u"+name, "u"+name, UintT(sz)) // uintN
		if sz >= 4 {
			name = fmt.Sprintf("float%d", sz*8)
			g.pkg.NewAlias(cpref+name, name, FloatT(sz)) // floatN
		}
	}

	// identifiers
	g.iot = NewIdentGo(cpref+"iota", "iota", g.e.UntypedIntT())
	g.lenF = NewIdentGo(cpref+"len", "len", g.e.FuncTT(g.intT, g.ifaceT))
	g.copyF = NewIdentGo(cpref+"copy", "copy", g.e.FuncTT(g.intT, g.ifaceT, g.ifaceT))
	g.panicF = NewIdentGo(cpref+"panic", "panic", g.e.FuncTT(nil, g.stringT))

	// stdlib
	g.osExitF = NewIdentGo("_Exit", "os.Exit", g.e.FuncTT(nil, g.intT))
}

// Bool returns Go bool type.
func (g *Go) Bool() Type {
	return g.boolT
}

// Byte returns Go byte type.
func (g *Go) Byte() Type {
	return g.byteT
}

// Rune returns Go rune type.
func (g *Go) Rune() Type {
	return g.runeT
}

// Uintptr returns Go uintptr type.
func (g *Go) Uintptr() Type {
	return g.uintptrT
}

// UnsafePtr returns Go unsafe.Pointer type.
func (g *Go) UnsafePtr() Type {
	// TODO: reserve a special type for it?
	return PtrT(g.e.PtrSize(), nil)
}

// Int returns Go int type.
func (g *Go) Int() Type {
	return g.intT
}

// Uint returns Go uint type.
func (g *Go) Uint() Type {
	return g.uintT
}

// Iface returns Go interface{} type.
func (g *Go) Iface() Type {
	return g.ifaceT
}

// IfaceSlice returns Go []interface{} type.
func (g *Go) IfaceSlice() Type {
	return g.ifaceSlT
}

// String returns Go string type.
func (g *Go) String() Type {
	return g.stringT
}

// Bytes returns Go []byte type.
func (g *Go) Bytes() Type {
	return g.bytesT
}

// Iota returns Go iota identifier.
func (g *Go) Iota() *Ident {
	return g.iot
}

// LenFunc returns Go len function identifier.
func (g *Go) LenFunc() *Ident {
	return g.lenF
}

// CopyFunc returns Go copy function identifier.
func (g *Go) CopyFunc() *Ident {
	return g.copyF
}

// PanicFunc returns Go panic function identifier.
func (g *Go) PanicFunc() *Ident {
	return g.panicF
}

// OsExitFunc returns Go os.Exit function identifier.
func (g *Go) OsExitFunc() *Ident {
	return g.osExitF
}
