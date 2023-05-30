package libs

import (
	"unicode"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/wctype.h.html

const (
	wctypeH = "wctype.h"
)

func init() {
	RegisterLibrary(wctypeH, func(c *Env) *Library {
		runeT := c.Go().Rune()
		isT := c.FuncTT(types.BoolT(), runeT)
		toT := c.FuncTT(runeT, runeT)
		return &Library{
			Imports: map[string]string{
				"libc":    RuntimeLibc,
				"unicode": "unicode",
			},
			Idents: map[string]*types.Ident{
				"iswalpha": c.NewIdent("iswalpha", "libc.IsAlpha", libc.IsAlpha, isT),
				"iswalnum": c.NewIdent("iswalnum", "libc.IsAlnum", libc.IsAlnum, isT),
				//"iswblank": c.NewIdent("iswblank", "libc.IsBlank", libc.IsBlank, isT),
				"iswcntrl": c.NewIdent("iswcntrl", "unicode.IsControl", unicode.IsControl, isT),
				"iswdigit": c.NewIdent("iswdigit", "unicode.IsDigit", unicode.IsDigit, isT),
				"iswgraph": c.NewIdent("iswgraph", "unicode.IsGraphic", unicode.IsGraphic, isT),
				"iswlower": c.NewIdent("iswlower", "unicode.IsLower", unicode.IsLower, isT),
				"iswprint": c.NewIdent("iswprint", "unicode.IsPrint", unicode.IsPrint, isT),
				"iswpunct": c.NewIdent("iswpunct", "unicode.IsPunct", unicode.IsPunct, isT),
				"iswspace": c.NewIdent("iswspace", "unicode.IsSpace", unicode.IsSpace, isT),
				"iswupper": c.NewIdent("iswupper", "unicode.IsUpper", unicode.IsUpper, isT),
				"towlower": c.NewIdent("towlower", "unicode.ToLower", unicode.ToLower, toT),
				"towupper": c.NewIdent("towupper", "unicode.ToUpper", unicode.ToUpper, toT),
			},
		}
	})
}
