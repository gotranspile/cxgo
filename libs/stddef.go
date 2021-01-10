package libs

import (
	"fmt"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

// TODO: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/stddef.h.html

const (
	stddefH = "stddef.h"
)

func init() {
	RegisterLibrary(stddefH, func(c *Env) *Library {
		return &Library{
			Header: incStdDef(c.Env),
			Types:  typesStdDef(c.Env),
		}
	})
}

func incStdDef(e *types.Env) string {
	var buf strings.Builder
	buf.WriteString(`
#include <` + stdintH + `>
#define NULL 0
typedef intptr_t ptrdiff_t;
typedef uintptr_t size_t;
`)
	sz := e.C().WCharSize()
	signed := e.C().WCharSigned()
	_, _ = fmt.Fprintf(&buf, "typedef %s wchar_t;\n",
		buildinFixedIntName(sz*8, !signed),
	)
	return buf.String()
}

func typesStdDef(e *types.Env) map[string]types.Type {
	return map[string]types.Type{
		"wchar_t":   e.C().WChar(),
		"size_t":    e.UintPtrT(),
		"ptrdiff_t": e.IntPtrT(),
	}
}
