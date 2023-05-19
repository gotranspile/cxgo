package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysTypesH = "sys/types.h"
)

//go:embed includes/sys_types.h
var hsysTypes string

func init() {
	RegisterLibrary(sysTypesH, func(c *Env) *Library {
		intT := types.IntT(4)
		return &Library{
			Idents: map[string]*types.Ident{
				"O_RDONLY": c.NewIdent("O_RDONLY", "csys.O_RDONLY", csys.O_RDONLY, intT),
				"O_WRONLY": c.NewIdent("O_WRONLY", "csys.O_WRONLY", csys.O_WRONLY, intT),
				"O_RDWR":   c.NewIdent("O_RDWR", "csys.O_RDWR", csys.O_RDWR, intT),
				"O_CREAT":  c.NewIdent("O_CREAT", "csys.O_CREAT", csys.O_CREAT, intT),
				"O_EXCL":   c.NewIdent("O_EXCL", "csys.O_EXCL", csys.O_EXCL, intT),
				"O_TRUNC":  c.NewIdent("O_TRUNC", "csys.O_TRUNC", csys.O_TRUNC, intT),
			},
			// TODO
			Header: hsysTypes,
		}
	})
}
