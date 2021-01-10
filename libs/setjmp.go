package libs

import "github.com/dennwc/cxgo/types"

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
			Header: `
#include <` + BuiltinH + `>

typedef struct jmp_buf {
	_cxgo_go_int (*SetJump) ();
	void (*LongJump) (_cxgo_go_int);
} jmp_buf;

#define setjmp(b) ((jmp_buf)b).SetJump()
#define longjmp(b, v) ((jmp_buf)b).LongJump(v)
`,
		}
	})
}
