package types

import "fmt"

const GoPrefix = "_cxgo_go_"

var (
	goArch4 = newGo(4)
	goArch8 = newGo(8)
)

func GoArch(size int) *Go {
	switch size {
	case 4:
		return goArch4
	case 8:
		return goArch8
	default:
		return newGo(size)
	}
}

func newGo(size int) *Go {
	pkg := newPackage("", "")
	// Those are native Go types that will always be mapped to themselves when transpiling
	// All other types might be mapped differently from C to Go
	g := &Go{
		size: size, pkg: pkg,
		// register basic Go types
		boolT: pkg.NewAlias(GoPrefix+"bool", "bool", BoolT()),
		byteT: pkg.NewTypeGo(GoPrefix+"byte", "byte", UintT(1)),
		runeT: pkg.NewTypeGo(GoPrefix+"rune", "rune", IntT(4)),

		uintptrT: pkg.NewTypeGo(GoPrefix+"uintptr", "uintptr", UintT(size)),
		intT:     pkg.NewTypeGo(GoPrefix+"int", "int", IntT(size)),
		uintT:    pkg.NewTypeGo(GoPrefix+"uint", "uint", UintT(size)),
		stringT:  pkg.NewTypeGo(GoPrefix+"string", "string", UnkT(size*3)),
		anyT:     pkg.NewTypeGo(GoPrefix+"any", "any", UnkT(size*2)),

		// register well-know slice types
		bytesT:    pkg.NewTypeGo(GoPrefix+"bytes", "[]byte", UnkT(size*3)),
		sliceAnyT: pkg.NewTypeGo(GoPrefix+"slice_any", "[]any", UnkT(size*3)),
	}

	// register fixed-size builtin Go types
	for _, sz := range []int{
		1, 2, 4, 8,
	} {
		name := fmt.Sprintf("int%d", sz*8)
		g.pkg.NewAlias(GoPrefix+name, name, IntT(sz))          // intN
		g.pkg.NewAlias(GoPrefix+"u"+name, "u"+name, UintT(sz)) // uintN
		if sz >= 4 {
			name = fmt.Sprintf("float%d", sz*8)
			g.pkg.NewAlias(GoPrefix+name, name, FloatT(sz)) // floatN
		}
	}

	// identifiers
	g.iot = NewIdentGo(GoPrefix+"iota", "iota", UntypedIntT(g.size))
	g.lenF = NewIdentGo(GoPrefix+"len", "len", FuncTT(g.size, g.intT, g.anyT))
	g.copyF = NewIdentGo(GoPrefix+"copy", "copy", FuncTT(g.size, g.intT, g.anyT, g.anyT))
	g.panicF = NewIdentGo(GoPrefix+"panic", "panic", FuncTT(g.size, nil, g.stringT))

	// stdlib
	g.osExitF = NewIdentGo("_Exit", "os.Exit", FuncTT(g.size, nil, g.intT))
	return g
}

type Go struct {
	size int // size of (u)int and pointers
	pkg  *Package

	// don't forget to update g.Types() when adding new types here

	boolT     Type
	byteT     Type
	runeT     Type
	uintptrT  Type
	intT      Type
	uintT     Type
	anyT      Type
	sliceAnyT Type
	stringT   Type
	bytesT    Type

	iot     *Ident
	lenF    *Ident
	copyF   *Ident
	panicF  *Ident
	osExitF *Ident
}

// Go returns a package containing builtin Go types.
func (e *Env) Go() *Go {
	return e.g
}

func (e *Env) initGo() {
	// TODO: we are assuming Go arch = C arch here
	e.g = GoArch(e.conf.PtrSize)
}

func (g *Go) Types() []Type {
	return []Type{
		g.boolT,
		g.byteT,
		g.runeT,
		g.uintptrT,
		g.intT,
		g.uintT,
		g.anyT,
		g.sliceAnyT,
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
	return PtrT(g.size, nil)
}

// Int returns Go int type.
func (g *Go) Int() Type {
	return g.intT
}

// Uint returns Go uint type.
func (g *Go) Uint() Type {
	return g.uintT
}

// Any returns Go any type.
func (g *Go) Any() Type {
	return g.anyT
}

// SliceOfAny returns Go []any type.
func (g *Go) SliceOfAny() Type {
	return g.sliceAnyT
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
