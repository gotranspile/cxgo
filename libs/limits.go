package libs

import (
	"fmt"
	"math"
	"strings"
)

const (
	limitsH = "limits.h"
)

func init() {
	RegisterLibrary(limitsH, func(c *Env) *Library {
		var buf strings.Builder
		buf.WriteString("#define CHAR_BIT 8")

		intMinMax(&buf, "SCHAR", math.MinInt8, math.MaxInt8)
		uintMax(&buf, "UCHAR", math.MaxUint8)
		intMinMax(&buf, "CHAR", math.MinInt8, math.MaxInt8)

		intMinMax(&buf, "SHRT", math.MinInt16, math.MaxInt16)
		uintMax(&buf, "USHRT", math.MaxUint16)

		switch c.IntSize() {
		case 4:
			intMinMax(&buf, "INT", math.MinInt32, math.MaxInt32)
			uintMax(&buf, "UINT", math.MaxUint32)
			intMinMax(&buf, "LONG", math.MinInt32, math.MaxInt32)
			uintMax(&buf, "ULONG", math.MaxUint32)
		case 8:
			intMinMax(&buf, "INT", math.MinInt64, math.MaxInt64)
			uintMax(&buf, "UINT", math.MaxUint64)
			intMinMax(&buf, "LONG", math.MinInt64, math.MaxInt64)
			uintMax(&buf, "ULONG", math.MaxUint64)
		}

		intMinMax(&buf, "LLONG", math.MinInt64, math.MaxInt64)
		uintMax(&buf, "ULLONG", math.MaxUint64)

		return &Library{
			Header: buf.String(),
		}
	})
}

func uintMax(buf *strings.Builder, pref string, max uint64) {
	_, _ = fmt.Fprintf(buf, "#define %s_MAX %d\n", pref, max)
}

func intMinMax(buf *strings.Builder, pref string, min, max int64) {
	_, _ = fmt.Fprintf(buf, "#define %s_MIN %d\n", pref, min)
	uintMax(buf, pref, uint64(max))
}
