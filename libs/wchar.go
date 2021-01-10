package libs

import (
	"github.com/dennwc/cxgo/runtime/libc"
	"github.com/dennwc/cxgo/types"
)

const (
	wcharH = "wchar.h"
)

func init() {
	RegisterLibrary(wcharH, func(c *Env) *Library {
		intT := types.IntT(8)
		return &Library{
			Header: `
#include <` + stddefH + `>
#include <` + wctypeH + `>
int wctob (wint_t wc);
`,
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Idents: map[string]*types.Ident{
				"wctob": c.NewIdent("wctob", "libc.Wctob", libc.Wctob, c.FuncTT(intT, c.C().WChar())),
			},
		}
	})
}
