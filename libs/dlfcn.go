package libs

import (
	"github.com/gotranspile/cxgo/runtime/dlopen"
	"github.com/gotranspile/cxgo/types"
)

const (
	dlfcnH = "dlfcn.h"
)

func init() {
	RegisterLibrary(dlfcnH, func(c *Env) *Library {
		gint := c.Go().Int()
		libT := types.NamedTGo("_cxgo_dllib", "dlopen.Library", c.MethStructT(map[string]*types.FuncType{
			"Sym":   c.FuncTT(c.Go().UnsafePtr(), c.Go().String()),
			"Close": c.FuncTT(gint),
		}))
		l := &Library{
			Idents: map[string]*types.Ident{
				"RTLD_LAZY":   types.NewIdentGo("RTLD_LAZY", "dlopen.RTLD_LAZY", gint),
				"RTLD_NOW":    types.NewIdentGo("RTLD_NOW", "dlopen.RTLD_NOW", gint),
				"RTLD_GLOBAL": types.NewIdentGo("RTLD_GLOBAL", "dlopen.RTLD_GLOBAL", gint),
				"RTLD_LOCAL":  types.NewIdentGo("RTLD_LOCAL", "dlopen.RTLD_LOCAL", gint),
			},
			Types: map[string]types.Type{
				"_cxgo_dllib": libT,
			},
			Imports: map[string]string{
				"libc":   RuntimeLibc,
				"dlopen": RuntimePrefix + "dlopen",
			},
			Header: `
const int RTLD_LAZY = 1;
const int RTLD_NOW = 2;
const int RTLD_GLOBAL = 4;
const int RTLD_LOCAL = 8;

typedef struct _cxgo_dllib {
	void* (*Sym) (_cxgo_go_string);
	_cxgo_go_int (*Close) ();
} _cxgo_dllib;

#define dlsym(l, s) ((_cxgo_dllib*)l)->Sym(s)
#define dlclose(l) ((_cxgo_dllib*)b)->Close(v)
`,
		}
		l.Declare(
			c.NewIdent("dlopen", "dlopen.Open", dlopen.Open, c.FuncTT(c.PtrT(libT), c.Go().String(), gint)),
			c.NewIdent("dlerror", "dlopen.Error", dlopen.Error, c.FuncTT(c.C().String())),
		)
		return l
	})
}
