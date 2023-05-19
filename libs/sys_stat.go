package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysStatH = "sys/stat.h"
)

//go:embed includes/sys_stat.h
var hsys_stat string

func init() {
	RegisterLibrary(sysStatH, func(c *Env) *Library {
		intT := types.IntT(4)
		strT := c.C().String()
		timevalT := c.GetLibraryType(timeH, "timeval")
		modeT := types.NamedTGo("mode_t", "csys.Mode", intT)
		statT := types.NamedTGo("stat", "csys.StatRes", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("st_dev", "Dev", intT)},
			{Name: types.NewIdentGo("st_ino", "Inode", intT)},
			{Name: types.NewIdentGo("st_mode", "Mode", modeT)},
			{Name: types.NewIdentGo("st_nlink", "Links", intT)},
			{Name: types.NewIdentGo("st_uid", "UID", intT)},
			{Name: types.NewIdentGo("st_gid", "GID", intT)},
			{Name: types.NewIdentGo("st_rdev", "RDev", intT)},
			{Name: types.NewIdentGo("st_size", "Size", types.UintT(8))},
			{Name: types.NewIdentGo("st_atime", "ATime", timevalT)},
			{Name: types.NewIdentGo("st_mtime", "MTime", timevalT)},
			{Name: types.NewIdentGo("st_ctime", "CTime", timevalT)},
			{Name: types.NewIdentGo("st_blksize", "BlockSize", intT)},
			{Name: types.NewIdentGo("st_blocks", "Blocks", intT)},
		}))
		return &Library{
			Imports: map[string]string{
				"csys": RuntimePrefix + "csys",
			},
			Types: map[string]types.Type{
				"mode_t": modeT,
				"stat":   statT,
			},
			Idents: map[string]*types.Ident{
				"stat":    c.NewIdent("stat", "csys.Stat", csys.Stat, c.FuncTT(intT, strT, c.PtrT(statT))),
				"chmod":   c.NewIdent("chmod", "csys.Chmod", csys.Chmod, c.FuncTT(intT, strT, modeT)),
				"mkdir":   c.NewIdent("mkdir", "csys.Mkdir", csys.Mkdir, c.FuncTT(intT, strT, modeT)),
				"S_ISDIR": c.NewIdent("S_ISDIR", "csys.IsDir", csys.IsDir, c.FuncTT(intT, modeT)),
			},
			// TODO
			Header: hsys_stat,
		}
	})
}
