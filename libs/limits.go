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
		switch c.IntSize() {
		case 4:
			intMinMax(&buf, "INT", math.MinInt32, math.MaxInt32)
		case 8:
			intMinMax(&buf, "INT", math.MinInt64, math.MaxInt64)
		}
		return &Library{
			Header: buf.String(),
		}
	})
}

func intMinMax(buf *strings.Builder, pref string, min, max int64) {
	_, _ = fmt.Fprintf(buf, "#define %s_MIN %d\n", pref, min)
	_, _ = fmt.Fprintf(buf, "#define %s_MAX %d\n", pref, max)
}
