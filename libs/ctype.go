package libs

import (
	"unicode"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

const (
	ctypeH = "ctype.h"
)

func init() {
	RegisterLibrary(ctypeH, func(c *Env) *Library {
		runeT := c.Go().Rune()
		boolT := c.Go().Bool()
		isT := c.FuncTT(boolT, runeT)
		toT := c.FuncTT(runeT, runeT)
		return &Library{
			Imports: map[string]string{
				"libc":    RuntimeLibc,
				"unicode": "unicode",
			},
			Idents: map[string]*types.Ident{
				"isalnum": c.NewIdent("isalnum", "libc.IsAlnum", libc.IsAlnum, isT),
				"isalpha": c.NewIdent("isalpha", "libc.IsAlpha", libc.IsAlpha, isT),
				//"isascii": c.NewIdent("isascii", "libc.IsASCII", libc.IsASCII, isT),
				"iscntrl": c.NewIdent("iscntrl", "unicode.IsControl", unicode.IsControl, isT),
				"isdigit": c.NewIdent("isdigit", "unicode.IsDigit", unicode.IsDigit, isT),
				"isgraph": c.NewIdent("isgraph", "unicode.IsGraphic", unicode.IsGraphic, isT),
				"islower": c.NewIdent("islower", "unicode.IsLower", unicode.IsLower, isT),
				"isprint": c.NewIdent("isprint", "unicode.IsPrint", unicode.IsPrint, isT),
				"ispunct": c.NewIdent("ispunct", "unicode.IsPunct", unicode.IsPunct, isT),
				"isspace": c.NewIdent("isspace", "unicode.IsSpace", unicode.IsSpace, isT),
				"isupper": c.NewIdent("isupper", "unicode.IsUpper", unicode.IsUpper, isT),
				//"isxdigit": c.NewIdent("isxdigit", "libc.IsXDigit", libc.IsXDigit, isT),
				//"toascii": c.NewIdent("toascii", "libc.ToASCII", libc.ToASCII, toT),
				"tolower": c.NewIdent("tolower", "unicode.ToLower", unicode.ToLower, toT),
				"toupper": c.NewIdent("toupper", "unicode.ToUpper", unicode.ToUpper, toT),
			},
			// TODO
			Header: `
_cxgo_go_bool isalnum(_cxgo_go_rune);
_cxgo_go_bool isalpha(_cxgo_go_rune);
_cxgo_go_bool isascii(_cxgo_go_rune);
_cxgo_go_bool iscntrl(_cxgo_go_rune);
_cxgo_go_bool isdigit(_cxgo_go_rune);
_cxgo_go_bool isgraph(_cxgo_go_rune);
_cxgo_go_bool islower(_cxgo_go_rune);
_cxgo_go_bool isprint(_cxgo_go_rune);
_cxgo_go_bool ispunct(_cxgo_go_rune);
_cxgo_go_bool isspace(_cxgo_go_rune);
_cxgo_go_bool isupper(_cxgo_go_rune);
_cxgo_go_bool isxdigit(_cxgo_go_rune);
_cxgo_go_rune toascii(_cxgo_go_rune);
_cxgo_go_rune tolower(_cxgo_go_rune);
_cxgo_go_rune toupper(_cxgo_go_rune);
`,
		}
	})
}
