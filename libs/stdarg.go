package libs

const (
	StdargH = "stdarg.h"
)

func init() {
	RegisterLibrary(StdargH, func(c *Env) *Library {
		// all types and method are define as builtins, so only macros here
		l := &Library{
			Header: `
#include <` + BuiltinH + `>
`,
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
		}
		l.Header += `
#define va_list __builtin_va_list
#define va_start(va, t) va.Start(t, _rest)
#define va_arg(va, typ) (typ)(va.Arg())
#define va_end(va) va.End()
#define va_copy(dst, src) __builtin_va_copy(dst, src)
`
		return l
	})
}
