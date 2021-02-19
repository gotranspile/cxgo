package libs

import (
	"fmt"
	"math"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

const (
	limitsH = "limits.h"
)

func init() {
	RegisterLibrary(limitsH, func(c *Env) *Library {
		var buf strings.Builder
		buf.WriteString("#define CHAR_BIT 8")

		idents := make(map[string]*types.Ident)
		intMinMax(&buf, idents, "SCHAR", "Int", math.MinInt8, math.MaxInt8, 8)
		uintMax(&buf, idents, "UCHAR", "Uint", math.MaxUint8, 8)
		intMinMax(&buf, idents, "CHAR", "Int", math.MinInt8, math.MaxInt8, 8)

		intMinMax(&buf, idents, "SHRT", "Int", math.MinInt16, math.MaxInt16, 16)
		uintMax(&buf, idents, "USHRT", "Uint", math.MaxUint16, 16)

		switch c.IntSize() {
		case 4:
			intMinMax(&buf, idents, "INT", "Int", math.MinInt32, math.MaxInt32, 32)
			uintMax(&buf, idents, "UINT", "Uint", math.MaxUint32, 32)
			intMinMax(&buf, idents, "LONG", "Int", math.MinInt32, math.MaxInt32, 32)
			uintMax(&buf, idents, "ULONG", "Uint", math.MaxUint32, 32)
		case 8:
			intMinMax(&buf, idents, "INT", "Int", math.MinInt64, math.MaxInt64, 64)
			uintMax(&buf, idents, "UINT", "Uint", math.MaxUint64, 64)
			intMinMax(&buf, idents, "LONG", "Int", math.MinInt64, math.MaxInt64, 64)
			uintMax(&buf, idents, "ULONG", "Uint", math.MaxUint64, 64)
		}

		intMinMax(&buf, idents, "LLONG", "Int", math.MinInt64, math.MaxInt64, 64)
		uintMax(&buf, idents, "ULLONG", "Uint", math.MaxUint64, 64)

		return &Library{
			Idents: idents,
			Header: buf.String(),
		}
	})
}

func uintMax(buf *strings.Builder, m map[string]*types.Ident, cPref, goPref string, max uint64, size int) {
	cName := cPref + "_MAX"
	if m != nil && goPref != "" {
		m[cName] = types.NewIdentGo(cName, fmt.Sprintf("math.Max%s%d", goPref, size), types.UntypedIntT(size/8))
	}
	_, _ = fmt.Fprintf(buf, "#define %s %du\n", cName, max)
}

func intMinMax(buf *strings.Builder, m map[string]*types.Ident, cPref, goPref string, min, max int64, size int) {
	cName := cPref + "_MIN"
	if m != nil && goPref != "" {
		m[cName] = types.NewIdentGo(cName, fmt.Sprintf("math.Min%s%d", goPref, size), types.UntypedIntT(size/8))
	}
	_, _ = fmt.Fprintf(buf, "#define %s %d\n", cName, min)
	uintMax(buf, m, cPref, goPref, uint64(max), size)
}
