package libs

import (
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/assert.h.html
// https://pubs.opengroup.org/onlinepubs/9699919799/functions/assert.html

const (
	assertH = "assert.h"
)

func init() {
	RegisterLibrary(assertH, func(c *Env) *Library {
		strT := c.C().String()
		l := &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Header: `
#include <` + stdboolH + `>
#define static_assert(x, y) /* x, y */
`,
			Idents: map[string]*types.Ident{
				"_Static_assert": c.NewIdent("_Static_assert", "libc.StaticAssert", staticAssert, c.FuncTT(nil, types.BoolT(), strT)),
			},
		}
		l.Declare(c.C().AssertFunc())
		return l
	})
}

func staticAssert(_ bool, _ *byte) {}
