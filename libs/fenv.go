package libs

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/types"
)

const (
	fenvH = "fenv.h"
)

func init() {
	RegisterLibrary(fenvH, func(c *Env) *Library {
		gint := c.Go().Int()
		intT := types.IntT(4)
		return &Library{
			Imports: map[string]string{
				"cmath": RuntimePrefix + "cmath",
			},
			Idents: map[string]*types.Ident{
				"FE_TOWARDZERO": c.NewIdent("FE_TOWARDZERO", "cmath.TowardZero", cmath.TowardZero, gint),
				"fesetround":    c.NewIdent("fesetround", "cmath.FSetRound", cmath.FSetRound, c.FuncTT(intT, intT)),
			},
			// TODO
			Header: `
#include <` + BuiltinH + `>
const _cxgo_go_int FE_TOWARDZERO = 1;
_cxgo_sint32 fesetround(_cxgo_sint32 round);
`,
		}
	})
}
