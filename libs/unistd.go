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
			// TODO
			Header: `
#include <` + BuiltinH + `>
#include <` + StdlibH + `>
#include <` + sysTypesH + `>
#include <` + sysStatH + `>
#include <` + StdioH + `>

_cxgo_go_uintptr  creat(const char *, mode_t);
_cxgo_go_uintptr  open(const char *, _cxgo_sint32, ...);
_cxgo_sint32  fcntl(_cxgo_go_uintptr, _cxgo_sint32, ...);

_cxgo_sint32 access(const char *, _cxgo_sint32);
unsigned     alarm(unsigned);
_cxgo_sint32 chdir(const char *);
_cxgo_sint32 fchdir(int fd);
int          chown(const char *, uid_t, gid_t);
#define close(fd) _cxgo_fileByFD((_cxgo_go_uintptr)fd)->Close()
size_t       confstr(int, char *, size_t);
int          dup(int);
int          dup2(int, int);
int          execl(const char *, const char *, ...);
int          execle(const char *, const char *, ...);
int          execlp(const char *, const char *, ...);
int          execv(const char *, char *const []);
int          execve(const char *, char *const [], char *const []);
int          execvp(const char *, char *const []);
#define _exit(v) _Exit(v)
int          fchown(int, uid_t, gid_t);
pid_t        fork(void);
long         fpathconf(int, int);
int          ftruncate(int, off_t);
char        *getcwd(char *, _cxgo_go_int);
gid_t        getegid(void);
uid_t        geteuid(void);
gid_t        getgid(void);
int          getgroups(int, gid_t []);
_cxgo_go_int gethostname(char *, _cxgo_go_int);
char        *getlogin(void);
int          getlogin_r(char *, size_t);
int          getopt(int, char * const [], const char *);
pid_t        getpgrp(void);
pid_t        getpid(void);
pid_t        getppid(void);
uid_t        getuid(void);
int          isatty(int);
int          link(const char *, const char *);
_cxgo_uint64 lseek(_cxgo_go_uintptr, _cxgo_uint64, _cxgo_sint32);
long         pathconf(const char *, int);
int          pause(void);
int          pipe(int [2]);
#define read(fd, p, sz) _cxgo_fileByFD((_cxgo_go_uintptr)fd)->Read(p, sz)
#define write(fd, p, sz) _cxgo_fileByFD((_cxgo_go_uintptr)fd)->Write(p, sz)
ssize_t      readlink(const char *restrict, char *restrict, size_t);
_cxgo_sint32 rmdir(const char *);
int          setegid(gid_t);
int          seteuid(uid_t);
int          setgid(gid_t);
int          setpgid(pid_t, pid_t);
pid_t        setsid(void);
int          setuid(uid_t);
unsigned     sleep(unsigned);
int          symlink(const char *, const char *);
long         sysconf(int);
pid_t        tcgetpgrp(int);
int          tcsetpgrp(int, pid_t);
char        *ttyname(int);
int          ttyname_r(int, char *, size_t);
_cxgo_sint32 unlink(const char *);
`,
		}
	})
}
