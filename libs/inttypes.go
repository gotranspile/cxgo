package libs

// TODO: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/inttypes.h.html

func init() {
	RegisterLibrary("inttypes.h", func(c *Env) *Library {
		return &Library{
			Header: `
#include <` + stdintH + `>
`,
		}
	})
}
