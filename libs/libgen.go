package libs

const (
	libgenH = "libgen.h"
)

func init() {
	RegisterLibrary(libgenH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
char *dirname(char *path);
char *basename(char *path);
`,
		}
	})
}
