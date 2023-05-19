package libs

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

const StdlibH = "stdlib.h"

//go:embed includes/stdlib.h
var hstdlib string

// https://pubs.opengroup.org/onlinepubs/9699919799/

func init() {
	RegisterLibrary(StdlibH, func(c *Env) *Library {
		voidPtr := c.PtrT(nil)
		gstrT := c.Go().String()
		cstrT := c.C().String()
		wstrT := c.C().WString()
		intT := types.IntT(4)
		gintT := c.Go().Int()
		longT := types.IntT(8)
		uintT := types.UintT(4)
		l := &Library{
			Imports: map[string]string{
				"libc":  RuntimeLibc,
				"cmath": RuntimePrefix + "cmath",
			},
			Idents: map[string]*types.Ident{
				"abs":      c.NewIdent("abs", "cmath.Abs", cmath.Abs, c.FuncTT(longT, longT)),
				"_Exit":    c.Go().OsExitFunc(),
				"malloc":   c.C().MallocFunc(),
				"calloc":   c.C().CallocFunc(),
				"realloc":  c.NewIdent("realloc", "libc.Realloc", libc.Realloc, c.FuncTT(voidPtr, voidPtr, gintT)),
				"free":     c.C().FreeFunc(),
				"atoi":     c.NewIdent("atoi", "libc.Atoi", libc.Atoi, c.FuncTT(gintT, gstrT)),
				"atol":     c.NewIdent("atol", "libc.Atoi", libc.Atoi, c.FuncTT(gintT, gstrT)),
				"atof":     c.NewIdent("atof", "libc.Atof", libc.Atof, c.FuncTT(types.FloatT(8), gstrT)),
				"rand":     c.NewIdent("rand", "libc.Rand", libc.Rand, c.FuncTT(intT)),
				"srand":    c.NewIdent("srand", "libc.SeedRand", libc.SeedRand, c.FuncTT(nil, uintT)),
				"qsort":    c.NewIdent("qsort", "libc.Sort", libc.Sort, c.FuncTT(nil, voidPtr, uintT, uintT, c.FuncTT(intT, voidPtr, voidPtr))),
				"bsearch":  c.NewIdent("bsearch", "libc.Search", libc.Search, c.FuncTT(voidPtr, voidPtr, voidPtr, uintT, uintT, c.FuncTT(intT, voidPtr, voidPtr))),
				"mbstowcs": c.NewIdent("mbstowcs", "libc.Mbstowcs", libc.Mbstowcs, c.FuncTT(uintT, wstrT, cstrT, uintT)),
			},
			Header: hstdlib + fmt.Sprintf("#define RAND_MAX %d\n", libc.RandMax),
		}

		l.Declare(
			c.NewIdent("getenv", "os.Getenv", os.Getenv, c.FuncTT(c.Go().String(), c.Go().String())),
		)
		return l
	})
}
