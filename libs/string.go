package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/

const (
	stringH = "string.h"
)

//go:embed string.h
var hstring string

func init() {
	// TODO: should include <locale.h>
	RegisterLibrary(stringH, func(c *Env) *Library {
		gintT := c.Go().Int()
		voidP := c.PtrT(nil)
		cstrT := c.C().String()
		return &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Idents: map[string]*types.Ident{
				"memcmp":      c.NewIdent("memcmp", "libc.MemCmp", libc.MemCmp, c.FuncTT(gintT, voidP, voidP, gintT)),
				"memchr":      c.NewIdent("memchr", "libc.MemChr", libc.MemChr, c.FuncTT(cstrT, cstrT, c.Go().Byte(), gintT)),
				"strlen":      c.NewIdent("strlen", "libc.StrLen", libc.StrLen, c.FuncTT(gintT, cstrT)),
				"strchr":      c.NewIdent("strchr", "libc.StrChr", libc.StrChr, c.FuncTT(cstrT, cstrT, c.Go().Byte())),
				"strrchr":     c.NewIdent("strrchr", "libc.StrRChr", libc.StrRChr, c.FuncTT(cstrT, cstrT, c.Go().Byte())),
				"strstr":      c.NewIdent("strstr", "libc.StrStr", libc.StrStr, c.FuncTT(cstrT, cstrT, cstrT)),
				"strcmp":      c.NewIdent("strcmp", "libc.StrCmp", libc.StrCmp, c.FuncTT(gintT, cstrT, cstrT)),
				"strncmp":     c.NewIdent("strncmp", "libc.StrNCmp", libc.StrNCmp, c.FuncTT(gintT, cstrT, cstrT, gintT)),
				"strcasecmp":  c.NewIdent("strcasecmp", "libc.StrCaseCmp", libc.StrCaseCmp, c.FuncTT(gintT, cstrT, cstrT)),
				"strncasecmp": c.NewIdent("strncasecmp", "libc.StrNCaseCmp", libc.StrNCaseCmp, c.FuncTT(gintT, cstrT, cstrT, gintT)),
				"strcpy":      c.NewIdent("strcpy", "libc.StrCpy", libc.StrCpy, c.FuncTT(cstrT, cstrT, cstrT)),
				"strncpy":     c.NewIdent("strncpy", "libc.StrNCpy", libc.StrNCpy, c.FuncTT(cstrT, cstrT, cstrT, gintT)),
				"strcat":      c.NewIdent("strcat", "libc.StrCat", libc.StrCat, c.FuncTT(cstrT, cstrT, cstrT)),
				"strncat":     c.NewIdent("strncat", "libc.StrNCat", libc.StrNCat, c.FuncTT(cstrT, cstrT, cstrT, gintT)),
				"strtok":      c.NewIdent("strtok", "libc.StrTok", libc.StrTok, c.FuncTT(cstrT, cstrT, cstrT)),
				"strspn":      c.NewIdent("strspn", "libc.StrSpn", libc.StrSpn, c.FuncTT(gintT, cstrT, cstrT)),
				"strcspn":     c.NewIdent("strcspn", "libc.StrCSpn", libc.StrCSpn, c.FuncTT(gintT, cstrT, cstrT)),
			},
			Header: hstring,
		}
	})
}
