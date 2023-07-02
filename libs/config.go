package libs

import (
	"errors"

	"github.com/gotranspile/cxgo/types"
)

// NewEnv creates a new environment. It uses GOARCH env to set sensible defaults.
func NewEnv(conf types.Config) *Env {
	return &Env{
		Env:  types.NewEnv(conf),
		libs: make(map[string]*Library),
		imports: map[string]string{
			"unsafe": "unsafe",
			"math":   "math",
			"libc":   RuntimeLibc,
		},
		macros: make(map[string]bool),
	}
}

type Env struct {
	*types.Env
	NoLibs  bool              // completely disable library lookups
	Map     map[string]string // when searching for library name, consult the map first and search that name instead
	libs    map[string]*Library
	imports map[string]string
	macros  map[string]bool
}

func (c *Env) Clone() *Env {
	c2 := &Env{Env: c.Env, NoLibs: c.NoLibs}
	c2.libs = make(map[string]*Library)
	for k, v := range c.libs {
		c2.libs[k] = v
	}
	c2.imports = make(map[string]string)
	for k, v := range c.imports {
		c2.imports[k] = v
	}
	c2.macros = make(map[string]bool)
	for k, v := range c.macros {
		c2.macros[k] = v
	}
	c2.Map = make(map[string]string)
	for k, v := range c.Map {
		c2.Map[k] = v
	}
	return c2
}

func (c *Env) ResolveImport(name string) string {
	path := c.imports[name]
	if path == "" {
		path = name
	}
	return path
}

// LookupLibrary finds an already loaded Library. It is useful to prevent import loops.
//
// Typically, the GetLibrary function should be used instead, because it will load the library automatically, if needed.
func (c *Env) LookupLibrary(name string) *Library {
	if v, ok := c.Map[name]; ok {
		name = v
	}
	if c.NoLibs && name != BuiltinH {
		return nil
	}
	if v, ok := c.Map[name]; ok {
		name = v
	}
	l, ok := c.libs[name]
	if !ok {
		panic(errors.New("cannot find library: " + name))
	}
	return l
}
