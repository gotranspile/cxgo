package libs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotranspile/cxgo/types"
)

func TestStdInt(t *testing.T) {
	c := NewEnv(types.Config32())
	l, ok := c.GetLibrary(stdintH)
	require.True(t, ok)
	require.Equal(t, strings.TrimSpace(`#ifndef _cxgo_STDINT_H
#define _cxgo_STDINT_H
#include <cxgo_builtin.h>

typedef _cxgo_sint8 int8_t;
typedef _cxgo_sint16 int16_t;
typedef _cxgo_sint32 int32_t;
typedef _cxgo_sint64 int64_t;

typedef _cxgo_uint8 uint8_t;
typedef _cxgo_uint16 uint16_t;
typedef _cxgo_uint32 uint32_t;
typedef _cxgo_uint64 uint64_t;

typedef _cxgo_sint8 int_least8_t;
typedef _cxgo_sint16 int_least16_t;
typedef _cxgo_sint32 int_least32_t;
typedef _cxgo_sint64 int_least64_t;

typedef _cxgo_uint8 uint_least8_t;
typedef _cxgo_uint16 uint_least16_t;
typedef _cxgo_uint32 uint_least32_t;
typedef _cxgo_uint64 uint_least64_t;

typedef _cxgo_sint8 int_fast8_t;
typedef _cxgo_sint16 int_fast16_t;
typedef _cxgo_sint32 int_fast32_t;
typedef _cxgo_sint64 int_fast64_t;

typedef _cxgo_uint8 uint_fast8_t;
typedef _cxgo_uint16 uint_fast16_t;
typedef _cxgo_uint32 uint_fast32_t;
typedef _cxgo_uint64 uint_fast64_t;

typedef _cxgo_sint32 intptr_t;
typedef _cxgo_uint32 uintptr_t;

typedef _cxgo_sint64 intmax_t;
typedef _cxgo_uint64 uintmax_t;

#define INT8_MIN -128
#define INT8_MAX 127u
#define UINT8_MAX 255u
#define INT16_MIN -32768
#define INT16_MAX 32767u
#define UINT16_MAX 65535u
#define INT32_MIN -2147483648
#define INT32_MAX 2147483647u
#define UINT32_MAX 4294967295u
#define INT64_MIN -9223372036854775808
#define INT64_MAX 9223372036854775807u
#define UINT64_MAX 18446744073709551615u

#define INT_LEAST8_MIN -128
#define INT_LEAST8_MAX 127u
#define UINT_LEAST8_MAX 255u
#define INT_LEAST16_MIN -32768
#define INT_LEAST16_MAX 32767u
#define UINT_LEAST16_MAX 65535u
#define INT_LEAST32_MIN -2147483648
#define INT_LEAST32_MAX 2147483647u
#define UINT_LEAST32_MAX 4294967295u
#define INT_LEAST64_MIN -9223372036854775808
#define INT_LEAST64_MAX 9223372036854775807u
#define UINT_LEAST64_MAX 18446744073709551615u

#define INT_FAST8_MIN -128
#define INT_FAST8_MAX 127u
#define UINT_FAST8_MAX 255u
#define INT_FAST16_MIN -32768
#define INT_FAST16_MAX 32767u
#define UINT_FAST16_MAX 65535u
#define INT_FAST32_MIN -2147483648
#define INT_FAST32_MAX 2147483647u
#define UINT_FAST32_MAX 4294967295u
#define INT_FAST64_MIN -9223372036854775808
#define INT_FAST64_MAX 9223372036854775807u
#define UINT_FAST64_MAX 18446744073709551615u

#define INTPTR_MIN -2147483648
#define INTPTR_MAX 2147483647u
#define UINTPTR_MAX 4294967295u

#define INTMAX_MIN -9223372036854775808
#define INTMAX_MAX 9223372036854775807u
#define UINTMAX_MAX 18446744073709551615u

#define PTRDIFF_MIN -2147483648
#define PTRDIFF_MAX 2147483647u

#define SIZE_MAX 4294967295u

#define WCHAR_MIN 0
#define WCHAR_MAX 65535u

#define WINT_MIN 0
#define WINT_MAX 4294967295u



#endif // _cxgo_STDINT_H
`), strings.TrimSpace(l.Header))
}
