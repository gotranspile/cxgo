package libs

const (
	stdboolH = "stdbool.h"
)

func init() {
	RegisterLibrary(stdboolH, func(c *Env) *Library {
		return &Library{
			Header: `
#define bool _Bool
#define false 0
#define true 1
#define __bool_true_false_are_defined
`,
		}
	})
}
