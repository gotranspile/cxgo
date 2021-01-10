package libs

const (
	fcntlH = "fcntl.h"
)

func init() {
	RegisterLibrary(fcntlH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
`,
		}
	})
}
