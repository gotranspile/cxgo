package libs

import "github.com/gotranspile/cxgo/types"

const (
	stdboolH = "stdbool.h"
)

func init() {
	RegisterLibrary(stdboolH, func(c *Env) *Library {
		return &Library{
			Types: map[string]types.Type{
				"bool": types.BoolT(),
			},
			Idents: map[string]*types.Ident{
				"true":  types.NewIdentGo("true", "true", types.BoolT()),
				"false": types.NewIdentGo("false", "false", types.BoolT()),
			},
			Header: `
#define bool _Bool
#define false 0
#define true 1
#define __bool_true_false_are_defined
`,
		}
	})
}
