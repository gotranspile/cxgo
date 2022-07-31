package libs

const (
	mkdevH = "sys/mkdev.h"
)

func init() {
	RegisterLibrary(mkdevH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
int major(int);
int minor(int);
`,
		}
	})
}
