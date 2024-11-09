package types

import (
	"os"
	"sort"
	"unsafe"
)

// Default returns a default config.
func Default() Config {
	var c Config
	c.setDefaults()
	return c
}

// Config32 returns a default types config for 32 bit systems.
func Config32() Config {
	c := Config{PtrSize: 4, IntSize: 4}
	c.setDefaults()
	return c
}

// Config64 returns a default types config for 64 bit systems.
func Config64() Config {
	c := Config{PtrSize: 8, IntSize: 8}
	c.setDefaults()
	return c
}

// Config stores configuration for base types.
type Config struct {
	PtrSize     int  // size of pointers in bytes
	IntSize     int  // default int size in bytes
	WCharSize   int  // wchar_t size
	WCharSigned bool // is wchar_t signed?
	UseGoInt    bool // use Go int for C int and long
}

func (c *Config) setDefaults() {
	if c.WCharSize == 0 {
		c.WCharSize = 2
	}
	if c.PtrSize == 0 {
		switch os.Getenv("GOARCH") {
		case "386":
			c.PtrSize = 4
		case "amd64":
			c.PtrSize = 8
		default:
			c.PtrSize = int(unsafe.Sizeof((*int)(nil)))
		}
	}
	if c.IntSize == 0 {
		switch os.Getenv("GOARCH") {
		case "386":
			c.IntSize = 4
		case "amd64":
			c.IntSize = 8
		default:
			c.IntSize = int(unsafe.Sizeof(int(0)))
		}
	}
}

func NewEnv(c Config) *Env {
	e := &Env{
		conf: c,
		pkgs: make(map[string]*Package),
	}
	e.conf.setDefaults()
	// Go must be first, because C definitions may depend on it
	e.initGo()
	e.initC()

	// conversion functions
	e.stringGo2C = NewIdent("libc.CString", e.FuncTT(
		e.C().String(),
		e.Go().String(),
	))
	e.wstringGo2C = NewIdent("libc.CWString", e.FuncTT(
		e.C().WString(),
		e.Go().String(),
	))
	e.stringC2Go = NewIdent("libc.GoString", e.FuncTT(
		e.Go().String(),
		e.C().String(),
	))
	e.wstringC2Go = NewIdent("libc.GoWString", e.FuncTT(
		e.Go().String(),
		e.C().WString(),
	))
	return e
}

type Env struct {
	conf Config

	c           C
	g           *Go
	pkgs        map[string]*Package
	stringGo2C  *Ident
	wstringGo2C *Ident
	stringC2Go  *Ident
	wstringC2Go *Ident
}

// PtrSize returns size of the pointer.
func (e *Env) PtrSize() int {
	return e.conf.PtrSize
}

// IntSize returns default size of the integer.
func (e *Env) IntSize() int {
	return e.conf.IntSize
}

// PtrT returns a pointer type with a specified element.
func (e *Env) PtrT(t Type) PtrType {
	return PtrT(e.conf.PtrSize, t)
}

// IntPtrT returns a signed int type that can hold a pointer diff.
func (e *Env) IntPtrT() IntType {
	return IntT(e.conf.PtrSize)
}

// UintPtrT returns a unsigned int type that can hold a pointer.
// It is different from Go().Uintptr(), since it returns uint32/uint64 type directly.
func (e *Env) UintPtrT() IntType {
	return UintT(e.conf.PtrSize)
}

// DefIntT returns a default signed int type.
// It is different from Go().Int(), since it returns int32/int64 type directly.
func (e *Env) DefIntT() Type {
	if e.conf.UseGoInt {
		return e.g.Int()
	}
	return IntT(e.conf.IntSize)
}

// DefUintT returns a default unsigned int type.
// It is different from Go().Uint(), since it returns uint32/uint64 type directly.
func (e *Env) DefUintT() Type {
	if e.conf.UseGoInt {
		return e.g.Uint()
	}
	return UintT(e.conf.IntSize)
}

// DefFloatT returns a default float type.
func (e *Env) DefFloatT() Type {
	return FloatT(e.conf.IntSize)
}

// FuncT returns a function type with a given return type and named arguments.
// It's mostly useful for function declarations. See FuncTT for simplified version.
func (e *Env) FuncT(ret Type, args ...*Field) *FuncType {
	return FuncT(e.conf.PtrSize, ret, args...)
}

// FuncTT returns a function type with a given return type and arguments.
// To name arguments, use FuncT.
func (e *Env) FuncTT(ret Type, args ...Type) *FuncType {
	return FuncTT(e.conf.PtrSize, ret, args...)
}

// VarFuncT returns a variadic function type with a given return type and named arguments.
// It's mostly useful for function declarations. See VarFuncTT for simplified version.
func (e *Env) VarFuncT(ret Type, args ...*Field) *FuncType {
	return VarFuncT(e.conf.PtrSize, ret, args...)
}

// VarFuncTT returns a variadic function type with a given return type and arguments.
// To name arguments, use VarFuncT.
func (e *Env) VarFuncTT(ret Type, args ...Type) *FuncType {
	return VarFuncTT(e.conf.PtrSize, ret, args...)
}

func (e *Env) MethStructT(meth map[string]*FuncType) *StructType {
	fields := make([]*Field, 0, len(meth))
	for name, ft := range meth {
		id := NewIdent(name, ft)
		fields = append(fields, &Field{
			Name: id,
		})
	}
	return StructT(fields)
}

// StringGo2C is a builtin function that converts Go string to C string.
func (e *Env) StringGo2C() *Ident {
	return e.stringGo2C
}

// WStringGo2C is a builtin function that converts Go string to C wchar_t string.
func (e *Env) WStringGo2C() *Ident {
	return e.wstringGo2C
}

// StringC2Go is a builtin function that converts C string to Go string.
func (e *Env) StringC2Go() *Ident {
	return e.stringC2Go
}

// WStringC2Go is a builtin function that converts C wchar_t string to Go string.
func (e *Env) WStringC2Go() *Ident {
	return e.wstringC2Go
}

func newPackage(name, path string) *Package {
	return &Package{name: name, path: path}
}

func (e *Env) NewPackage(name, path string) *Package {
	if path == "" {
		path = name
	}
	p := newPackage(name, path)
	e.pkgs[path] = p
	return p
}

func (e *Env) PackageByPath(path string) *Package {
	return e.pkgs[path]
}

func (e *Env) Packages() []*Package {
	out := make([]*Package, 0, len(e.pkgs))
	for _, p := range e.pkgs {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].path < out[j].path
	})
	return out
}

type Package struct {
	name string
	path string

	idents  map[*Ident]struct{}
	cnames  map[string]Type
	gonames map[string]Type
}

func (p *Package) checkNames(cname, goname string) string {
	if cname != "" {
		if _, ok := p.cnames[cname]; ok {
			panic("type already exists: " + cname)
		}
	}
	if goname == "" {
		goname = goName(cname)
	}
	if goname != "" {
		if _, ok := p.gonames[goname]; ok {
			panic("type already exists: " + goname)
		}
	}
	return goname
}

func (p *Package) setNames(cname, goname string, t Type) {
	if cname != "" {
		if p.cnames == nil {
			p.cnames = make(map[string]Type)
		}
		p.cnames[cname] = t
	}
	if goname != "" {
		if p.gonames == nil {
			p.gonames = make(map[string]Type)
		}
		p.gonames[goname] = t
	}
}

func (p *Package) NewAlias(cname, goname string, t Type) Type {
	_ = p.checkNames(cname, goname)
	p.setNames(cname, goname, t)
	return t
}

func (p *Package) NewTypeC(cname string, t Type) Named {
	return p.NewTypeGo(cname, "", t)
}

func (p *Package) NewTypeGo(cname, goname string, t Type) Named {
	goname = p.checkNames(cname, goname)

	nt := NamedTGo(cname, goname, t)
	if p.idents == nil {
		p.idents = make(map[*Ident]struct{})
	}
	p.idents[nt.Name()] = struct{}{}
	if cname != "" {
		if p.cnames == nil {
			p.cnames = make(map[string]Type)
		}
		p.cnames[cname] = nt
	}
	if goname != "" {
		if p.gonames == nil {
			p.gonames = make(map[string]Type)
		}
		p.gonames[goname] = nt
	}
	return nt
}

func (p *Package) CType(name string) Type {
	return p.cnames[name]
}

func (p *Package) GoType(name string) Type {
	return p.gonames[name]
}

func (p *Package) CNamedType(name string) Named {
	t := p.cnames[name]
	if t == nil {
		return nil
	}
	nt, _ := t.(Named)
	return nt
}

func (p *Package) GoNamedType(name string) Named {
	t := p.gonames[name]
	if t == nil {
		return nil
	}
	nt, _ := t.(Named)
	return nt
}
