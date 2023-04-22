package libs

import _ "embed"

const (
	StdargH = "stdarg.h"
)

//go:embed stdarg.h
var hstdarg string

func init() {
	RegisterLibrary(StdargH, func(c *Env) *Library {
		// all types and method are define as builtins, so only macros here
		l := &Library{
			Header: hstdarg,
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
		}
		return l
	})
}
