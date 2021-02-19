package libs

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

const (
	stdintH = "stdint.h"
)

func init() {
	RegisterLibrary(stdintH, func(c *Env) *Library {
		idents := make(map[string]*types.Ident)
		return &Library{
			Idents: idents,
			Header: incStdInt(c.Env, idents),
			Types:  typesStdInt(c.Env),
		}
	})
}

func fixedIntTypeDefs(buf *strings.Builder, part string) {
	for _, unsigned := range []bool{false, true} {
		name := "int"
		if unsigned {
			name = "uint"
		}
		for _, sz := range intSizes {
			_, _ = fmt.Fprintf(buf, "#define %s%s%d_t %s\n",
				name, part, sz,
				buildinFixedIntName(sz, unsigned),
			)
		}
		buf.WriteByte('\n')
	}
}

func fixedIntTypes() string {
	var buf strings.Builder
	fixedIntTypeDefs(&buf, "")
	fixedIntTypeDefs(&buf, "_least")
	fixedIntTypeDefs(&buf, "_fast")
	return buf.String()
}

func maxIntTypeDefs(buf *strings.Builder, part string, sz int) {
	for _, unsigned := range []bool{false, true} {
		name := "int"
		if unsigned {
			name = "uint"
		}
		_, _ = fmt.Fprintf(buf, "typedef %s %s%s_t;\n",
			buildinFixedIntName(sz, unsigned),
			name, part,
		)
	}
	buf.WriteByte('\n')
}

func maxIntTypes(e *types.Env) string {
	var buf strings.Builder
	maxIntTypeDefs(&buf, "ptr", e.PtrSize()*8)
	maxIntTypeDefs(&buf, "max", intSizes[len(intSizes)-1])
	return buf.String()
}

func intLimitsDef(buf *strings.Builder, m map[string]*types.Ident, part string, min, max int64, umax uint64, size int) {
	intMinMax(buf, m, "INT"+part, "Int", min, max, size)
	uintMax(buf, m, "UINT"+part, "Uint", umax, size)
}

func intLimitsDefs(buf *strings.Builder, m map[string]*types.Ident, part string) {
	for i, sz := range intSizes {
		intLimitsDef(buf, m, part+strconv.Itoa(sz), minInts[i], maxInts[i], maxUints[i], sz)
	}
	buf.WriteByte('\n')
}

func intLimits(m map[string]*types.Ident) string {
	var buf strings.Builder
	intLimitsDefs(&buf, m, "")
	intLimitsDefs(&buf, m, "_LEAST")
	intLimitsDefs(&buf, m, "_FAST")
	return buf.String()
}

func intSizeInd(sz int) int {
	switch sz {
	case 0, 1:
		return 0
	case 2:
		return 1
	case 4:
		return 2
	case 8:
		return 3
	default:
		panic(sz)
	}
}

func maxIntLimits(e *types.Env, m map[string]*types.Ident) string {
	var buf strings.Builder
	i := intSizeInd(e.PtrSize())
	intLimitsDef(&buf, m, "PTR", minInts[i], maxInts[i], maxUints[i], e.PtrSize()*8)
	buf.WriteByte('\n')

	i = len(intSizes) - 1
	intLimitsDef(&buf, m, "MAX", minInts[i], maxInts[i], maxUints[i], e.PtrSize()*8)
	buf.WriteByte('\n')
	return buf.String()
}

func otherLimits(e *types.Env, m map[string]*types.Ident) string {
	var buf strings.Builder

	i := intSizeInd(e.PtrSize())
	intMinMax(&buf, m, "PTRDIFF", "Int", minInts[i], maxInts[i], e.PtrSize()*8)
	buf.WriteByte('\n')

	// TODO: SIG_ATOMIC

	i = intSizeInd(e.PtrSize())
	uintMax(&buf, m, "SIZE", "Uint", maxUints[i], e.PtrSize()*8)
	buf.WriteByte('\n')

	i = intSizeInd(e.C().WCharSize())
	if e.C().WCharSigned() {
		intMinMax(&buf, m, "WCHAR", "Int", minInts[i], maxInts[i], e.C().WCharSize()*8)
	} else {
		intMinMax(&buf, m, "WCHAR", "Uint", 0, int64(maxUints[i]), e.C().WCharSize()*8)
	}
	buf.WriteByte('\n')

	i = intSizeInd(e.C().WIntSize())
	if e.C().WCharSigned() {
		intMinMax(&buf, m, "WINT", "Int", minInts[i], maxInts[i], e.C().WIntSize()*8)
	} else {
		intMinMax(&buf, m, "WINT", "Uint", 0, int64(maxUints[i]), e.C().WIntSize()*8)
	}
	buf.WriteByte('\n')
	return buf.String()
}

var (
	intSizes = []int{8, 16, 32, 64}
	minInts  = []int64{math.MinInt8, math.MinInt16, math.MinInt32, math.MinInt64}
	maxInts  = []int64{math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
	maxUints = []uint64{math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}
)

func incStdInt(e *types.Env, m map[string]*types.Ident) string {
	var buf strings.Builder
	buf.WriteString("#include <" + BuiltinH + ">\n")
	buf.WriteString(fixedIntTypes())
	buf.WriteString(maxIntTypes(e))
	buf.WriteString(intLimits(m))
	buf.WriteString(maxIntLimits(e, m))
	buf.WriteString(otherLimits(e, m))
	return buf.String()
}

func typesStdInt(e *types.Env) map[string]types.Type {
	m := make(map[string]types.Type, 2*len(intSizes))
	for _, unsigned := range []bool{false, true} {
		for _, sz := range intSizes {
			name := buildinFixedIntName(sz, unsigned)
			if unsigned {
				m[name] = types.UintT(sz / 8)
			} else {
				m[name] = types.IntT(sz / 8)
			}
		}
	}
	m["intptr_t"] = e.IntPtrT()
	m["uintptr_t"] = e.UintPtrT()
	max := intSizes[len(intSizes)-1] / 8
	m["intmax_t"] = types.IntT(max)
	m["uintmax_t"] = types.UintT(max)
	return m
}
