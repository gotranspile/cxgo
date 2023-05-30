package libs

import (
	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/gotranspile/cxgo/types"
)

const (
	unistdH = "unistd.h"
)

func init() {
	RegisterLibrary(unistdH, func(c *Env) *Library {
		modeT := c.GetLibraryType(sysStatH, "mode_t")
		uintptrT := c.Go().Uintptr()
		fdT := uintptrT
		intT := types.IntT(4)
		gintT := c.Go().Int()
		ulongT := types.UintT(8)
		strT := c.C().String()
		return &Library{
			Imports: map[string]string{
				"stdio": RuntimePrefix + "stdio",
				"csys":  RuntimePrefix + "csys",
				"cnet":  RuntimePrefix + "cnet",
			},
			Idents: map[string]*types.Ident{
				"creat":       c.NewIdent("creat", "stdio.Create", stdio.Create, c.FuncTT(fdT, strT, modeT)),
				"open":        c.NewIdent("open", "stdio.Open", stdio.Open, c.VarFuncTT(fdT, strT, intT)),
				"fcntl":       c.NewIdent("fcntl", "stdio.FDControl", stdio.FDControl, c.VarFuncTT(intT, fdT, intT)),
				"chdir":       c.NewIdent("chdir", "stdio.Chdir", stdio.Chdir, c.FuncTT(intT, strT)),
				"rmdir":       c.NewIdent("rmdir", "stdio.Rmdir", stdio.Rmdir, c.FuncTT(intT, strT)),
				"unlink":      c.NewIdent("unlink", "stdio.Unlink", stdio.Unlink, c.FuncTT(intT, strT)),
				"access":      c.NewIdent("access", "stdio.Access", stdio.Access, c.FuncTT(intT, strT, intT)),
				"lseek":       c.NewIdent("lseek", "stdio.Lseek", stdio.Lseek, c.FuncTT(ulongT, fdT, ulongT, intT)),
				"getcwd":      c.NewIdent("getcwd", "stdio.GetCwd", stdio.GetCwd, c.FuncTT(strT, strT, gintT)),
				"gethostname": c.NewIdent("gethostname", "cnet.GetHostname", cnet.GetHostname, c.FuncTT(gintT, strT, gintT)),
			},
		}
	})
}
