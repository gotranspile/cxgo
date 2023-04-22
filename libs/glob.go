package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/gotranspile/cxgo/types"
)

const (
	globH = "glob.h"
)

//go:embed glob.h
var hglob string

func init() {
	RegisterLibrary(globH, func(c *Env) *Library {
		gint := c.Go().Int()
		intT := types.IntT(4)
		sizeT := intT
		strT := c.C().String()
		globT := types.NamedT("stdio.Glob", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("gl_pathc", "Num", sizeT)},
			{Name: types.NewIdentGo("gl_pathv", "Paths", c.PtrT(strT))},
			{Name: types.NewIdentGo("gl_offs", "Reserve", sizeT)},
			{Name: types.NewIdentGo("Glob", "Glob", c.FuncTT(intT, strT, intT, c.FuncTT(intT, strT, intT)))},
			{Name: types.NewIdentGo("Free", "Free", c.FuncTT(nil))},
		}))
		return &Library{
			Imports: map[string]string{
				"libc":  RuntimeLibc,
				"stdio": RuntimePrefix + "stdio",
			},
			Types: map[string]types.Type{
				"glob_t": globT,
			},
			Idents: map[string]*types.Ident{
				"GLOB_NOESCAPE": c.NewIdent("GLOB_NOESCAPE", "stdio.GlobNoEscape", stdio.GlobNoEscape, gint),
			},
			// TODO
			Header: hglob,
		}
	})
}
