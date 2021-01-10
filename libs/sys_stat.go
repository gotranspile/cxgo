package libs

import (
	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysStatH = "sys/stat.h"
)

func init() {
	RegisterLibrary(sysStatH, func(c *Env) *Library {
		intT := types.IntT(4)
		strT := c.C().String()
		timeLib := c.GetLib(timeH)
		timevalT := timeLib.GetType("timeval")
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
			{Name: types.NewIdentGo("st_atim", "ATime", timevalT)},
			{Name: types.NewIdentGo("st_mtim", "MTime", timevalT)},
			{Name: types.NewIdentGo("st_ctim", "CTime", timevalT)},
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
			Header: `
#include <` + sysTypesH + `>
#include <` + sysTimeH + `>

typedef _cxgo_sint32 mode_t;

struct stat {
    _cxgo_sint32  st_dev;     /* ID of device containing file */
    _cxgo_sint32  st_ino;     /* inode number */
    mode_t    st_mode;    /* protection */
    _cxgo_sint32     st_nlink;   /* number of hard links */
    _cxgo_sint32       st_uid;     /* user ID of owner */
    _cxgo_sint32       st_gid;     /* group ID of owner */
    _cxgo_sint32       st_rdev;    /* device ID (if special file) */
    off_t       st_size;    /* total size, in bytes */
    struct timeval      st_atim;   /* time of last access */
    struct timeval      st_mtim;   /* time of last modification */
    struct timeval      st_ctim;   /* time of last status change */
    _cxgo_sint32   st_blksize; /* blocksize for filesystem I/O */
    _cxgo_sint32    st_blocks;  /* number of blocks allocated */
};

_cxgo_sint32  chmod(const char *, mode_t);
int    fchmod(int, mode_t);
int    fstat(int, struct stat *);
int    lstat(const char *restrict, struct stat *restrict);
_cxgo_sint32  mkdir(const char *, mode_t);
int    mkfifo(const char *, mode_t);
_cxgo_sint32    stat(const char *restrict, struct stat *restrict);
mode_t umask(mode_t);

_cxgo_sint32 S_ISDIR(mode_t m);
`,
		}
	})
}
