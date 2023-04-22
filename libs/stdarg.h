#include <cxgo_builtin.h>

#define va_list __builtin_va_list
#define va_start(va, t) va.Start(t, _rest)
#define va_arg(va, typ) (typ)(va.Arg())
#define va_end(va) va.End()
#define va_copy(dst, src) __builtin_va_copy(dst, src)
