package libs

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/

func init() {
	// TODO: should include <locale.h>
	RegisterLibrary("string.h", func(c *Env) *Library {
		gintT := c.Go().Int()
		voidP := c.PtrT(nil)
		cstrT := c.C().String()
		return &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Idents: map[string]*types.Ident{
				"memcmp":      c.NewIdent("memcmp", "libc.MemCmp", libc.MemCmp, c.FuncTT(gintT, voidP, voidP, gintT)),
				"memchr":      c.NewIdent("memchr", "libc.MemChr", libc.MemChr, c.FuncTT(cstrT, cstrT, c.Go().Byte(), gintT)),
				"strlen":      c.NewIdent("strlen", "libc.StrLen", libc.StrLen, c.FuncTT(gintT, cstrT)),
				"strchr":      c.NewIdent("strchr", "libc.StrChr", libc.StrChr, c.FuncTT(cstrT, cstrT, c.Go().Byte())),
				"strrchr":     c.NewIdent("strrchr", "libc.StrRChr", libc.StrRChr, c.FuncTT(cstrT, cstrT, c.Go().Byte())),
				"strstr":      c.NewIdent("strstr", "libc.StrStr", libc.StrStr, c.FuncTT(cstrT, cstrT, cstrT)),
				"strcmp":      c.NewIdent("strcmp", "libc.StrCmp", libc.StrCmp, c.FuncTT(gintT, cstrT, cstrT)),
				"strncmp":     c.NewIdent("strncmp", "libc.StrNCmp", libc.StrNCmp, c.FuncTT(gintT, cstrT, cstrT, gintT)),
				"strcasecmp":  c.NewIdent("strcasecmp", "libc.StrCaseCmp", libc.StrCaseCmp, c.FuncTT(gintT, cstrT, cstrT)),
				"strncasecmp": c.NewIdent("strncasecmp", "libc.StrNCaseCmp", libc.StrNCaseCmp, c.FuncTT(gintT, cstrT, cstrT, gintT)),
				"strcpy":      c.NewIdent("strcpy", "libc.StrCpy", libc.StrCpy, c.FuncTT(cstrT, cstrT, cstrT)),
				"strncpy":     c.NewIdent("strncpy", "libc.StrNCpy", libc.StrNCpy, c.FuncTT(cstrT, cstrT, cstrT, gintT)),
				"strcat":      c.NewIdent("strcat", "libc.StrCat", libc.StrCat, c.FuncTT(cstrT, cstrT, cstrT)),
				"strncat":     c.NewIdent("strncat", "libc.StrNCat", libc.StrNCat, c.FuncTT(cstrT, cstrT, cstrT, gintT)),
				"strtok":      c.NewIdent("strtok", "libc.StrTok", libc.StrTok, c.FuncTT(cstrT, cstrT, cstrT)),
				"strspn":      c.NewIdent("strspn", "libc.StrSpn", libc.StrSpn, c.FuncTT(gintT, cstrT, cstrT)),
				"strcspn":     c.NewIdent("strcspn", "libc.StrCSpn", libc.StrCSpn, c.FuncTT(gintT, cstrT, cstrT)),
				"strdup":      c.NewIdent("strdup", "libc.StrDup", libc.StrDup, c.FuncTT(cstrT, cstrT)),
			},
			Header: `
#include <` + BuiltinH + `>
#include <` + stddefH + `>
#include <` + StdlibH + `>

#define memcpy __builtin_memcpy
#define memmove __builtin_memmove
#define memset __builtin_memset

void    *memccpy(void *restrict, const void *restrict, int, _cxgo_go_int);
void    *memchr(const void *, _cxgo_go_byte, _cxgo_go_int);
_cxgo_go_int      memcmp(const void *, const void *, _cxgo_go_int);
char    *stpcpy(char *restrict, const char *restrict);
char    *stpncpy(char *restrict, const char *restrict, size_t);
char    *strcat(char *restrict, const char *restrict);
char    *strchr(const char *, _cxgo_go_byte);
_cxgo_go_int      strcmp(const char *, const char *);
_cxgo_go_int      strcoll(const char *, const char *);
//int      strcoll_l(const char *, const char *, locale_t);
char    *strcpy(char *restrict, const char *restrict);
_cxgo_go_int   strcspn(const char *, const char *);
char    *strdup(const char *);
char    *strerror(int);
//char    *strerror_l(int, locale_t);
int      strerror_r(int, char *, size_t);
_cxgo_go_int   strlen(const char *);
char    *strncat(char *restrict, const char *restrict, _cxgo_go_int);
_cxgo_go_int      strncmp(const char *, const char *, _cxgo_go_int);
char    *strncpy(char *restrict, const char *restrict, _cxgo_go_int);
char    *strndup(const char *, _cxgo_go_int);
_cxgo_go_int   strnlen(const char *, _cxgo_go_int);
char    *strpbrk(const char *, const char *);
char    *strrchr(const char *, _cxgo_go_byte);
char    *strsignal(int);
_cxgo_go_int   strspn(const char *, const char *);
char    *strstr(const char *, const char *);
char    *strtok(char *restrict, const char *restrict);
char    *strtok_r(char *restrict, const char *restrict, char **restrict);
size_t   strxfrm(char *restrict, const char *restrict, size_t);
//size_t   strxfrm_l(char *restrict, const char *restrict, size_t, locale_t);
_cxgo_go_int strcasecmp(const char *s1, const char *s2);
_cxgo_go_int strncasecmp(const char *s1, const char *s2, _cxgo_go_int n);
`,
		}
	})
}
