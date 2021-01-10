package libs

import (
	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysIoctlH = "sys/ioctl.h"
)

func init() {
	RegisterLibrary(sysIoctlH, func(c *Env) *Library {
		uintptrT := c.Go().Uintptr()
		return &Library{
			// TODO
			Imports: map[string]string{
				"csys": RuntimePrefix + "csys",
			},
			Idents: map[string]*types.Ident{
				"FIONREAD": c.NewIdent("FIONREAD", "csys.FIONREAD", csys.FIONREAD, uintptrT),
				"ioctl":    c.NewIdent("ioctl", "csys.Ioctl", csys.Ioctl, c.VarFuncTT(types.IntT(4), uintptrT, uintptrT)),
			},
			Header: `
#include <` + sysTypesH + `>

const _cxgo_go_uintptr FIONREAD = 1;
_cxgo_sint32 ioctl(_cxgo_go_uintptr fildes, _cxgo_go_uintptr request, ... /* arg */); 
`,
		}
	})
}
