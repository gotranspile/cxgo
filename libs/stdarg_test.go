package libs

import (
	"github.com/dennwc/cxgo/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStdargH(t *testing.T) {
	c := NewEnv(types.Config32())
	l, ok := c.GetLibrary(StdargH)
	require.True(t, ok)
	require.Equal(t, `
#ifndef _cxgo_STDARG_H
#define _cxgo_STDARG_H


#include <cxgo_builtin.h>

#define va_list __builtin_va_list
#define va_start(va, t) va.Start(t, _rest)
#define va_arg(va, typ) (typ)(va.Arg())
#define va_end(va) va.End()
#define va_copy(dst, src) __builtin_va_copy(dst, src)


#endif // _cxgo_STDARG_H
`, l.Header)
}
