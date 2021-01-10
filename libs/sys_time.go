package libs

import (
	"github.com/dennwc/cxgo/runtime/csys"
	"github.com/dennwc/cxgo/types"
)

const (
	sysTimeH = "sys/time.h"
)

func init() {
	RegisterLibrary(sysTimeH, func(c *Env) *Library {
		intT := types.IntT(4)
		timeLib := c.GetLib(timeH)
		timevalT := timeLib.GetType("timeval")
		timespecT := timeLib.GetType("timespec")
		_ = timespecT
		return &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
				"csys": RuntimePrefix + "csys",
			},
			Idents: map[string]*types.Ident{
				"gettimeofday": c.NewIdent("gettimeofday", "csys.GetTimeOfDay", csys.GetTimeOfDay, c.FuncTT(intT, c.PtrT(timevalT), c.PtrT(nil))),
			},
			// TODO
			Header: `
#include <` + timeH + `>
#include <` + sysTypesH + `>

typedef struct fd_set {
	long fds_bits[];
} fd_set;

_cxgo_int32   getitimer(_cxgo_int32, struct itimerval *);
_cxgo_int32   gettimeofday(struct timeval *restrict, void *restrict);
int   select(int, fd_set *restrict, fd_set *restrict, fd_set *restrict, struct timeval *restrict);
int   setitimer(int, const struct itimerval *restrict, struct itimerval *restrict);
`,
		}
	})
}
