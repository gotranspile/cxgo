package libs

import (
	"github.com/dennwc/cxgo/runtime/stdio"
	"github.com/dennwc/cxgo/types"
)

const (
	globH = "glob.h"
)

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
			Header: `
#include <` + stddefH + `>

const _cxgo_go_int GLOB_NOESCAPE = 1;

typedef struct {
    size_t   gl_pathc;    /* Count of paths matched so far  */
    char   **gl_pathv;    /* List of matched pathnames.  */
    size_t   gl_offs;     /* Slots to reserve in gl_pathv.  */
	_cxgo_sint32 (*Glob)(const char *pattern, _cxgo_sint32 flags,
                _cxgo_sint32 (*errfunc) (const char *epath, _cxgo_sint32 eerrno));
	void (*Free)(void);
} glob_t;
#define glob(pattern, flags, errfunc, g) ((glob_t*)g)->Glob(pattern, flags, errfunc)
#define globfree(g) ((glob_t*)g)->Free()
`,
		}
	})
}
