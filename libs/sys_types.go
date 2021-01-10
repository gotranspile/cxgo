package libs

import (
	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysTypesH = "sys/types.h"
)

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
			Header: `
#include <` + BuiltinH + `>
#include <` + stddefH + `>
#include <` + timeH + `>

#define off_t _cxgo_int64
#define ssize_t _cxgo_int64
#define off_t _cxgo_uint64
#define pid_t _cxgo_uint64
#define gid_t _cxgo_uint32
#define uid_t _cxgo_uint32

#define u_short unsigned short
#define u_long unsigned long


// TODO: should be in fcntl.h
const _cxgo_int32 O_RDONLY = 1;
const _cxgo_int32 O_WRONLY = 2;
const _cxgo_int32 O_RDWR = 3;
const _cxgo_int32 O_CREAT = 4;
const _cxgo_int32 O_EXCL = 5;
const _cxgo_int32 O_TRUNC = 6;
`,
		}
	})
}
