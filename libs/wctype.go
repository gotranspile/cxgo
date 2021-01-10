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
			Header: `
#include <` + stdboolH + `>
#include <` + stddefH + `>
#include <` + StdioH + `>

#define wint_t _cxgo_go_rune
const wint_t WEOF = -1;

_Bool      iswalnum(wint_t);
_Bool      iswalpha(wint_t);
_Bool      iswblank(wint_t);
_Bool      iswcntrl(wint_t);
//_Bool      iswctype(wint_t, wctype_t);
_Bool      iswdigit(wint_t);
_Bool      iswgraph(wint_t);
_Bool      iswlower(wint_t);
_Bool      iswprint(wint_t);
_Bool      iswpunct(wint_t);
_Bool      iswspace(wint_t);
_Bool      iswupper(wint_t);
_Bool      iswxdigit(wint_t);
//wint_t    towctrans(wint_t, wctrans_t);
wint_t    towlower(wint_t);
wint_t    towupper(wint_t);
//wctrans_t wctrans(const char *);
//wctype_t  wctype(const char *);
/*
_Bool      iswalnum_l(wint_t, locale_t);
_Bool      iswalpha_l(wint_t, locale_t);
_Bool      iswblank_l(wint_t, locale_t);
_Bool      iswcntrl_l(wint_t, locale_t);
_Bool      iswctype_l(wint_t, wctype_t, locale_t);
_Bool      iswdigit_l(wint_t, locale_t);
_Bool      iswgraph_l(wint_t, locale_t);
_Bool      iswlower_l(wint_t, locale_t);
_Bool      iswprint_l(wint_t, locale_t);
_Bool      iswpunct_l(wint_t, locale_t);
_Bool      iswspace_l(wint_t, locale_t);
_Bool      iswupper_l(wint_t, locale_t);
_Bool      iswxdigit_l(wint_t, locale_t);
wint_t    towctrans_l(wint_t, wctrans_t, locale_t);
wint_t    towlower_l(wint_t, locale_t);
wint_t    towupper_l(wint_t, locale_t);
wctrans_t wctrans_l(const char *, locale_t);
wctype_t  wctype_l(const char *, locale_t);
*/
`,
		}
	})
}
