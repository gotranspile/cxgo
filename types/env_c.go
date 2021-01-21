package types

import (
	"strings"
)

type C struct {
	e   *Env
	pkg *Package

	charT    Type
	wcharT   Type
	mallocF  *Ident
	freeF    *Ident
	callocF  *Ident
	memmoveF *Ident
	memcpyF  *Ident
	memsetF  *Ident
}

// C returns a package containing builtin C types.
func (e *Env) C() *C {
	return &e.c
}

func (e *Env) initC() {
	e.c.e = e
	e.c.pkg = e.newPackage("", "")
	e.c.init()
}

func (c *C) init() {
	g := c.e.Go()

	c.pkg.NewAlias("bool", "", c.Bool())
	c.charT = c.pkg.NewAlias("char", "", IntT(1)) // signed
	if c.e.conf.WCharSize == 2 && !c.e.conf.WCharSigned {
		c.wcharT = c.pkg.NewTypeGo("wchar_t", "libc.WChar", UintT(2))
	} else if c.e.conf.WCharSigned {
		c.wcharT = c.pkg.NewTypeC("wchar_t", IntT(c.e.conf.WCharSize))
	} else {
		c.wcharT = c.pkg.NewTypeC("wchar_t", UintT(c.e.conf.WCharSize))
	}

	c.pkg.NewAlias("BOOL", "", c.Bool())
	c.pkg.NewAlias("CHAR", "", c.Char())
	c.pkg.NewAlias("BYTE", "", g.Byte())
	c.pkg.NewAlias("_BYTE", "", g.Byte())
	c.pkg.NewAlias("INT", "", c.Int())
	c.pkg.NewAlias("UINT", "", c.UnsignedInt())
	c.pkg.NewAlias("LONG", "", c.Long())
	c.pkg.NewAlias("ULONG", "", c.UnsignedLong())
	c.pkg.NewAlias("LONGLONG", "", c.LongLong())
	c.pkg.NewAlias("ULONGLONG", "", c.UnsignedLongLong())
	c.pkg.NewAlias("WORD", "", UintT(2))
	c.pkg.NewAlias("_WORD", "", UintT(2))
	c.pkg.NewAlias("DWORD", "", UintT(4))
	c.pkg.NewAlias("_DWORD", "", UintT(4))
	c.pkg.NewAlias("QWORD", "", UintT(8))
	c.pkg.NewAlias("_QWORD", "", UintT(8))

	unsafePtr := g.UnsafePtr()
	c.mallocF = NewIdentGo("__builtin_malloc", "libc.Malloc", c.e.FuncTT(unsafePtr, g.Int()))
	c.freeF = NewIdentGo("free", "libc.Free", c.e.FuncTT(nil, unsafePtr))
	c.callocF = NewIdentGo("calloc", "libc.Calloc", c.e.FuncTT(unsafePtr, g.Int(), g.Int()))
	c.memmoveF = NewIdentGo("__builtin_memmove", "libc.MemMove", c.e.FuncTT(unsafePtr, unsafePtr, unsafePtr, g.Int()))
	c.memcpyF = NewIdentGo("__builtin_memcpy", "libc.MemCpy", c.e.FuncTT(unsafePtr, unsafePtr, unsafePtr, g.Int()))
	c.memsetF = NewIdentGo("__builtin_memset", "libc.MemSet", c.e.FuncTT(unsafePtr, unsafePtr, g.Byte(), g.Int()))
}

func (c *C) WCharSize() int {
	return c.e.conf.WCharSize
}

func (c *C) WIntSize() int {
	return c.WCharSize() * 2
}

func (c *C) WCharSigned() bool {
	return c.e.conf.WCharSigned
}

// Type finds a C builtin type by name.
func (c *C) Type(name string) Type {
	if t := c.pkg.CType(name); t != nil {
		return nil
	}
	// check if the element type of LP* type has an override as well
	if strings.HasPrefix(name, "LP") {
		if elem := c.pkg.CType(name[2:]); elem != nil {
			return c.e.PtrT(elem)
		}
	}
	// same for types with a _PTR suffix
	if strings.HasSuffix(name, "_PTR") {
		if elem := c.pkg.CType(name[:len(name)-4]); elem != nil {
			return c.e.PtrT(elem)
		}
	}
	return nil
}

// Named finds a named C builtin type.
func (c *C) NamedType(name string) Named {
	// to run the same type resolution
	nt, _ := c.Type(name).(Named)
	return nt
}

// Bool returns C bool type.
func (c *C) Bool() Type {
	// TODO: support custom bool types
	return c.e.Go().Bool()
}

// Char returns C char type.
func (c *C) Char() Type {
	return c.charT
}

// WChar returns C wchar_t type.
func (c *C) WChar() Type {
	return c.wcharT
}

// SignedChar returns C signed char type.
func (c *C) SignedChar() Type {
	return IntT(1)
}

// UnsignedChar returns C unsigned char type.
func (c *C) UnsignedChar() Type {
	return UintT(1)
}

// Short returns C short type.
func (c *C) Short() Type {
	return IntT(2)
}

// UnsignedShort returns C unsigned short type.
func (c *C) UnsignedShort() Type {
	return UintT(2)
}

// Int returns C int type.
func (c *C) Int() Type {
	if c.e.conf.UseGoInt {
		return c.e.Go().Int()
	}
	return c.e.DefIntT()
}

// UnsignedInt returns C unsigned int type.
func (c *C) UnsignedInt() Type {
	if c.e.conf.UseGoInt {
		return c.e.Go().Uint()
	}
	return c.e.DefUintT()
}

// Long returns C long type.
func (c *C) Long() Type {
	return c.Int()
}

// UnsignedLong returns C unsigned long type.
func (c *C) UnsignedLong() Type {
	return c.UnsignedInt()
}

// LongLong returns C long long type.
func (c *C) LongLong() Type {
	return IntT(8)
}

// UnsignedLongLong returns C unsigned long long type.
func (c *C) UnsignedLongLong() Type {
	return UintT(8)
}

// Float returns C float type.
func (c *C) Float() Type {
	return FloatT(4)
}

// Double returns C double type.
func (c *C) Double() Type {
	return FloatT(8)
}

// String returns C char* type.
func (c *C) String() PtrType {
	// we represent it differently
	return c.e.PtrT(c.e.Go().Byte())
}

// WString returns C wchar_t* type.
func (c *C) WString() PtrType {
	return c.e.PtrT(c.WChar())
}

// BytesN returns C char[N] type.
func (c *C) BytesN(n int) Type {
	return ArrayT(c.e.Go().Byte(), n)
}

// MallocFunc returns C malloc function ident.
func (c *C) MallocFunc() *Ident {
	return c.mallocF
}

// FreeFunc returns C free function ident.
func (c *C) FreeFunc() *Ident {
	return c.freeF
}

// CallocFunc returns C calloc function ident.
func (c *C) CallocFunc() *Ident {
	return c.callocF
}

// MemmoveFunc returns C memmove function ident.
func (c *C) MemmoveFunc() *Ident {
	return c.memmoveF
}

// MemcpyFunc returns C memcpy function ident.
func (c *C) MemcpyFunc() *Ident {
	return c.memcpyF
}

// MemsetFunc returns C memset function ident.
func (c *C) MemsetFunc() *Ident {
	return c.memsetF
}
