package libs

const (
	StdargH = "stdarg.h"
)

func init() {
	RegisterLibrary(StdargH, func(c *Env) *Library {
		// all types and method are define as builtins, so only macros here
		l := &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
		}
		return l
	})
}
