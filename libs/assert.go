package libs

import (
	"github.com/dennwc/cxgo/runtime/libc"
	"github.com/dennwc/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/assert.h.html
// https://pubs.opengroup.org/onlinepubs/9699919799/functions/assert.html

const (
	assertH = "assert.h"
)

func init() {
	RegisterLibrary(assertH, func(c *Env) *Library {
		strT := c.C().String()
		return &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Header: `
#include <` + stdboolH + `>
void _cxgo_assert(bool);
#define assert _cxgo_assert
#define static_assert(x, y) /* x, y */
`,
			Idents: map[string]*types.Ident{
				"_cxgo_assert":   c.NewIdent("_cxgo_assert", "libc.Assert", libc.Assert, c.FuncTT(nil, types.BoolT())),
				"_Static_assert": c.NewIdent("_Static_assert", "libc.StaticAssert", staticAssert, c.FuncTT(nil, types.BoolT(), strT)),
			},
		}
	})
}

func staticAssert(_ bool, _ *byte) {}
