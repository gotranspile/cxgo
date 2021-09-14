package libs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

type Library struct {
	created bool

	Name string
	// Header is a content of a C header file for the library.
	// It will be protected by a ifndef guard automatically.
	Header string
	// Types overrides type definitions parsed from the library's header with a custom ones.
	Types map[string]types.Type

	Idents map[string]*types.Ident

	Imports map[string]string

	ForceMacros map[string]bool
}

func (l *Library) GetType(name string) types.Type {
	t, ok := l.Types[name]
	if !ok {
		panic(errors.New("cannot find type: " + name))
	}
	return t
}

// Declare adds a new ident to the library. It also writes a corresponding C definition to the header.
func (l *Library) Declare(ids ...*types.Ident) {
	for _, id := range ids {
		l.declare(id)
	}
}

func (l *Library) declare(id *types.Ident) {
	if l.Idents == nil {
		l.Idents = make(map[string]*types.Ident)
	}
	l.Idents[id.Name] = id
	switch tp := id.CType(nil).(type) {
	case *types.FuncType:
		l.Header += fmt.Sprintf("%s %s(", cTypeStr(tp.Return()), id.Name)
		for i, a := range tp.Args() {
			if i != 0 {
				l.Header += ", "
			}
			l.Header += cTypeStr(a.Type())
		}
		if tp.Variadic() {
			if tp.ArgN() > 0 {
				l.Header += ", "
			}
			l.Header += "..."
		}
		l.Header += ");\n"
	default:
		l.Header += fmt.Sprintf("%s %s;\n", cTypeStr(tp), id.Name)
	}
}

func cTypeStr(t types.Type) string {
	switch t := t.(type) {
	case nil:
		return "void"
	case types.Named:
		return t.Name().Name
	case types.PtrType:
		if t.Elem() == nil {
			return "void*"
		} else if el, ok := t.Elem().(types.Named); ok && el.Name().GoName == "byte" {
			return "char*"
		}
		return cTypeStr(t.Elem()) + "*"
	case types.IntType:
		s := ""
		if t.Signed() {
			s = "signed "
		} else {
			s = "unsigned "
		}
		s += fmt.Sprintf("__int%d", t.Sizeof()*8)
		return s
	case types.FloatType:
		switch t.Sizeof() {
		case 4:
			return "float"
		case 8:
			return "double"
		default:
			return fmt.Sprintf("_cxgo_float%d", t.Sizeof()*8)
		}
	case types.BoolType:
		return "_Bool"
	default:
		panic(fmt.Errorf("TODO: %T", t))
	}
}

var libs = make(map[string]LibraryFunc)

type LibraryFunc func(c *Env) *Library

// RegisterLibrary registers an override for a C library.
func RegisterLibrary(name string, fnc LibraryFunc) {
	if name == "" {
		panic("empty name")
	}
	if fnc == nil {
		panic("no constructor")
	}
	if _, ok := libs[name]; ok {
		panic("already registered")
	}
	libs[name] = fnc
}

const IncludePath = "/_cxgo_overrides"

var defPathReplacer = strings.NewReplacer(
	"/", "_",
	".", "_",
)

// GetLibrary finds or initializes the library, given a C include filename.
func (c *Env) GetLibrary(name string) (*Library, bool) {
	if l, ok := c.libs[name]; ok {
		return l, true
	}
	fnc, ok := libs[name]
	if !ok {
		return nil, false
	}
	l := fnc(c)
	l.created = true
	l.Name = name
	//for name, typ := range l.Types {
	//	named, ok := typ.(types.Named)
	//	if !ok {
	//		continue
	//	}
	//	if _, ok := l.Idents[name]; !ok {
	//		if l.Idents == nil {
	//			l.Idents = make(map[string]*types.Ident)
	//		}
	//		l.Idents[name] = named.Name()
	//	}
	//}
	c.libs[name] = l
	for k, v := range l.Imports {
		c.imports[k] = v
	}
	for k, v := range l.ForceMacros {
		c.macros[k] = v
	}

	ifdef := "_cxgo_" + strings.ToUpper(defPathReplacer.Replace(name))
	l.Header = fmt.Sprintf(`
#ifndef %s
#define %s

%s

#endif // %s
`,
		ifdef,
		ifdef,
		l.Header,
		ifdef,
	)
	return l, true
}

// GetLibraryType is a helper for GetLibrary followed by GetType.
func (c *Env) GetLibraryType(lib, typ string) types.Type {
	l, ok := c.GetLibrary(lib)
	if !ok {
		panic("cannot find library: " + lib)
	}
	return l.GetType(typ)
}

func (c *Env) NewLibrary(path string) (*Library, bool) {
	if !strings.HasPrefix(path, IncludePath+"/") {
		return nil, false // only override ones in our fake lookup path
	}
	name := strings.TrimPrefix(path, IncludePath+"/")
	return c.GetLibrary(name)
}

func (c *Env) TypeByName(name string) (types.Type, bool) {
	for _, l := range c.libs {
		if t, ok := l.Types[name]; ok {
			return t, true
		}
	}
	return nil, false
}

func (c *Env) LibIdentByName(name string) (*Library, *types.Ident, bool) {
	for _, l := range c.libs {
		if id, ok := l.Idents[name]; ok {
			return l, id, true
		}
	}
	return nil, nil, false
}

func (c *Env) IdentByName(name string) (*types.Ident, bool) {
	_, id, ok := c.LibIdentByName(name)
	return id, ok
}

func (c *Env) ForceMacro(name string) bool {
	ok := c.macros[name]
	return ok
}
