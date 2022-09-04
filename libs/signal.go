package libs

const (
	sysSignalH = "signal.h"
)

func init() {
	RegisterLibrary(sysSignalH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
`,
		}
	})
}
