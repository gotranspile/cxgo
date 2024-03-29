package libs

import (
	"github.com/gotranspile/cxgo/types"
)

const (
	setjmpH = "setjmp.h"
)

func init() {
	RegisterLibrary(setjmpH, func(c *Env) *Library {
		gint := c.Go().Int()
		bufT := types.NamedTGo("jmp_buf", "libc.JumpBuf", c.MethStructT(map[string]*types.FuncType{
			"SetJump":  c.FuncTT(gint),
			"LongJump": c.FuncTT(nil, gint),
		}))
		return &Library{
			Types: map[string]types.Type{
				"jmp_buf": bufT,
			},
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
		}
	})
}
