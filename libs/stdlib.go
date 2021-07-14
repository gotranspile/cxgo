package libs

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
	"os"
)

const StdlibH = "stdlib.h"

// https://pubs.opengroup.org/onlinepubs/9699919799/

func init() {
	RegisterLibrary(StdlibH, func(c *Env) *Library {
		voidPtr := c.PtrT(nil)
		gstrT := c.Go().String()
		cstrT := c.C().String()
		wstrT := c.C().WString()
		intT := types.IntT(4)
		gintT := c.Go().Int()
		longT := types.IntT(8)
		uintT := types.UintT(4)
		l := &Library{
			Imports: map[string]string{
				"libc":  RuntimeLibc,
				"cmath": RuntimePrefix + "cmath",
			},
			Idents: map[string]*types.Ident{
				"abs":      c.NewIdent("abs", "cmath.Abs", cmath.Abs, c.FuncTT(longT, longT)),
				"_Exit":    c.Go().OsExitFunc(),
				"malloc":   c.C().MallocFunc(),
				"calloc":   c.C().CallocFunc(),
				"realloc":  c.NewIdent("realloc", "libc.Realloc", libc.Realloc, c.FuncTT(voidPtr, voidPtr, gintT)),
				"free":     c.C().FreeFunc(),
				"atoi":     c.NewIdent("atoi", "libc.Atoi", libc.Atoi, c.FuncTT(gintT, gstrT)),
				"atol":     c.NewIdent("atol", "libc.Atoi", libc.Atoi, c.FuncTT(gintT, gstrT)),
				"atof":     c.NewIdent("atof", "libc.Atof", libc.Atof, c.FuncTT(types.FloatT(8), gstrT)),
				"rand":     c.NewIdent("rand", "libc.Rand", libc.Rand, c.FuncTT(intT)),
				"srand":    c.NewIdent("srand", "libc.SeedRand", libc.SeedRand, c.FuncTT(nil, uintT)),
				"qsort":    c.NewIdent("qsort", "libc.Sort", libc.Sort, c.FuncTT(nil, voidPtr, uintT, uintT, c.FuncTT(intT, voidPtr, voidPtr))),
				"bsearch":  c.NewIdent("bsearch", "libc.Search", libc.Search, c.FuncTT(voidPtr, voidPtr, voidPtr, uintT, uintT, c.FuncTT(intT, voidPtr, voidPtr))),
				"mbstowcs": c.NewIdent("mbstowcs", "libc.Mbstowcs", libc.Mbstowcs, c.FuncTT(uintT, wstrT, cstrT, uintT)),
			},
			Header: `
#include <` + BuiltinH + `>
#include <` + stddefH + `>

#define EXIT_FAILURE 1
#define EXIT_SUCCESS 0

#define malloc __builtin_malloc
#define abort __builtin_abort
void free(void*);

typedef struct {
	int quot;
	int rem;
} div_t;

typedef struct {
	long quot;
	long rem;
} ldiv_t;

typedef struct {
	long long quot;
	long long rem;
} lldiv_t;

void          _Exit(_cxgo_go_int);
long          a64l(const char *);
_cxgo_int64   abs(_cxgo_int64);
int           atexit(void (*)(void));
double        atof(const char *);
_cxgo_go_int  atoi(const char *);
_cxgo_go_int  atol(const char *);
long long     atoll(const char *);
void         *bsearch(const void *, const void *, _cxgo_uint32, _cxgo_uint32, _cxgo_int32 (*)(const void *, const void *));
void         *calloc(_cxgo_go_int, _cxgo_go_int);
div_t         div(int, int);
double        drand48(void);
double        erand48(unsigned short [3]);
#define exit(x) _Exit(x)
int           getsubopt(char **, char *const *, char **);
int           grantpt(int);
char         *initstate(unsigned, char *, size_t);
long          jrand48(unsigned short [3]);
char         *l64a(long);
long          labs(long);
void          lcong48(unsigned short [7]);
ldiv_t        ldiv(long, long);
long long     llabs(long long);
lldiv_t       lldiv(long long, long long);
long          lrand48(void);
void         *malloc(_cxgo_go_int);
int           mblen(const char *, size_t);
_cxgo_uint32  mbstowcs(wchar_t *restrict, const char *restrict, _cxgo_uint32);
int           mbtowc(wchar_t *restrict, const char *restrict, size_t);
char         *mkdtemp(char *);
int           mkstemp(char *);
long          mrand48(void);
long          nrand48(unsigned short [3]);
int           posix_memalign(void **, size_t, size_t);
int           posix_openpt(int);
char         *ptsname(int);
int           putenv(char *);
void          qsort(void *, _cxgo_uint32, _cxgo_uint32, _cxgo_int32 (*)(const void *, const void *));
_cxgo_sint32  rand(void);
int           rand_r(unsigned *);
long          random(void);
void         *realloc(void *, _cxgo_go_int);
char         *realpath(const char *restrict, char *restrict);
unsigned short *seed48(unsigned short [3]);
int           setenv(const char *, const char *, int);
void          setkey(const char *);
char         *setstate(char *);
void          srand(_cxgo_uint32);
void          srand48(long);
void          srandom(unsigned);
double        strtod(const char *restrict, char **restrict);
float         strtof(const char *restrict, char **restrict);
long          strtol(const char *restrict, char **restrict, int);
long double   strtold(const char *restrict, char **restrict);
long long     strtoll(const char *restrict, char **restrict, int);
unsigned long strtoul(const char *restrict, char **restrict, int);
unsigned long long strtoull(const char *restrict, char **restrict, int);
int           system(const char *);
int           unlockpt(int);
int           unsetenv(const char *);
size_t        wcstombs(char *restrict, const wchar_t *restrict, size_t);
int           wctomb(char *, wchar_t);
`,
		}
		l.Header += fmt.Sprintf("#define RAND_MAX %d\n", libc.RandMax)
		l.Declare(
			c.NewIdent("getenv", "os.Getenv", os.Getenv, c.FuncTT(c.Go().String(), c.Go().String())),
		)
		return l
	})
}
