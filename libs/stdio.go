package libs

import (
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/gotranspile/cxgo/types"
)

// TODO: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/stdio.h.html

const (
	StdioH = "stdio.h"
)

func init() {
	// TODO: must #include <stdarg.h>
	// TODO: must #include <sys/types.h>
	// TODO: off_t and ssize_t must be in <sys/types.h>
	RegisterLibrary(StdioH, func(c *Env) *Library {
		gintT := c.Go().Int()
		gstrT := c.Go().String()
		intT := types.IntT(4)
		longT := types.IntT(8)
		cstrT := c.C().String()
		cbytesT := c.C().String()
		bytesT := c.C().String()
		filePtr := c.PtrT(nil)
		valistT := c.GetLibraryType(BuiltinH, "__builtin_va_list")
		fileT := types.NamedTGo("FILE", "stdio.File", c.MethStructT(map[string]*types.FuncType{
			"PutC":   c.FuncTT(intT, gintT),
			"PutS":   c.FuncTT(intT, cstrT),
			"WriteN": c.FuncTT(intT, cbytesT, gintT, gintT),
			"Write":  c.FuncTT(intT, cbytesT, gintT),
			"ReadN":  c.FuncTT(intT, bytesT, gintT, gintT),
			"Read":   c.FuncTT(intT, bytesT, gintT),
			"Tell":   c.FuncTT(longT),
			"Seek":   c.FuncTT(intT, longT, intT),
			"GetC":   c.FuncTT(gintT),
			"UnGetC": c.FuncTT(intT, gintT),
			"GetS":   c.FuncTT(cstrT, cstrT, intT),
			"IsEOF":  c.FuncTT(intT),
			"FileNo": c.FuncTT(c.Go().Uintptr()),
			"Flush":  c.FuncTT(intT),
			"Close":  c.FuncTT(intT),
		}))
		filePtr.SetElem(fileT)
		l := &Library{
			Imports: map[string]string{
				"stdio": RuntimePrefix + "stdio",
			},
			Types: map[string]types.Type{
				"FILE": fileT,
			},
			Idents: map[string]*types.Ident{
				"SEEK_SET": c.NewIdent("SEEK_SET", "stdio.SEEK_SET", stdio.SEEK_SET, intT),
				"SEEK_CUR": c.NewIdent("SEEK_CUR", "stdio.SEEK_CUR", stdio.SEEK_CUR, intT),
				"SEEK_END": c.NewIdent("SEEK_END", "stdio.SEEK_END", stdio.SEEK_END, intT),

				"_cxgo_EOF": c.NewIdent("_cxgo_EOF", "stdio.EOF", stdio.EOF, types.UntypedIntT(1)),
			},
			Header: `
#include <` + stddefH + `>
#include <` + StdargH + `>
#include <` + StdlibH + `>
#include <` + sysTypesH + `>

typedef struct FILE FILE;

typedef struct FILE {
	_cxgo_sint32 (*PutC)(_cxgo_go_int);
	_cxgo_sint32 (*PutS)(const char*);
	_cxgo_sint32 (*WriteN)(const void*, _cxgo_go_int, _cxgo_go_int);
	_cxgo_sint32 (*Write)(const void*, _cxgo_go_int);
	_cxgo_sint32 (*ReadN)(void*, _cxgo_go_int, _cxgo_go_int);
	_cxgo_sint32 (*Read)(void*, _cxgo_go_int);
	_cxgo_sint64 (*Tell)(void);
	_cxgo_sint32 (*Seek)(_cxgo_sint64, _cxgo_sint32);
	_cxgo_go_int (*GetC)(void);
	_cxgo_sint32 (*UnGetC)(_cxgo_go_int);
	char* (*GetS)(char*, _cxgo_sint32);
	_cxgo_sint32 (*IsEOF)(void);
	_cxgo_go_uintptr (*FileNo)(void);
	_cxgo_sint32 (*Flush)(void);
	_cxgo_sint32 (*Close)(void);
} FILE;

#define fpos_t _cxgo_uint64
#define int _cxgo_int64

#define FILENAME_MAX 255

const _cxgo_sint32 SEEK_SET = 0;
const _cxgo_sint32 SEEK_CUR = 1;
const _cxgo_sint32 SEEK_END = 2;

const int _cxgo_EOF = -1;
#define EOF _cxgo_EOF

`,
		}

		l.Declare(
			c.NewIdent("_cxgo_getStdout", "stdio.Stdout", stdio.Stdout, c.FuncTT(filePtr)),
			c.NewIdent("_cxgo_getStderr", "stdio.Stderr", stdio.Stderr, c.FuncTT(filePtr)),
			c.NewIdent("_cxgo_getStdin", "stdio.Stdin", stdio.Stdin, c.FuncTT(filePtr)),
			c.NewIdent("_cxgo_fileByFD", "stdio.ByFD", stdio.ByFD, c.FuncTT(filePtr, c.Go().Uintptr())),
			c.NewIdent("fopen", "stdio.FOpen", stdio.FOpen, c.FuncTT(filePtr, gstrT, gstrT)),
			c.NewIdent("fdopen", "stdio.FDOpen", stdio.FDOpen, c.FuncTT(filePtr, c.Go().Uintptr(), gstrT)),
			// note: printf itself is considered a builtin
			c.NewIdent("vprintf", "stdio.Vprintf", stdio.Vprintf, c.FuncTT(gintT, gstrT, valistT)),
			c.NewIdent("sprintf", "stdio.Sprintf", stdio.Sprintf, c.VarFuncTT(gintT, cstrT, gstrT)),
			c.NewIdent("vsprintf", "stdio.Vsprintf", stdio.Vsprintf, c.FuncTT(gintT, cstrT, gstrT, valistT)),
			c.NewIdent("snprintf", "stdio.Snprintf", stdio.Snprintf, c.VarFuncTT(gintT, cstrT, gintT, gstrT)),
			c.NewIdent("vsnprintf", "stdio.Vsnprintf", stdio.Vsnprintf, c.FuncTT(gintT, cstrT, gintT, gstrT, valistT)),
			c.NewIdent("fprintf", "stdio.Fprintf", stdio.Fprintf, c.VarFuncTT(gintT, filePtr, gstrT)),
			c.NewIdent("vfprintf", "stdio.Vfprintf", stdio.Vfprintf, c.FuncTT(gintT, filePtr, gstrT, valistT)),
			c.NewIdent("scanf", "stdio.Scanf", stdio.Scanf, c.VarFuncTT(gintT, gstrT)),
			c.NewIdent("vscanf", "stdio.Vscanf", stdio.Vscanf, c.FuncTT(gintT, gstrT, valistT)),
			c.NewIdent("sscanf", "stdio.Sscanf", stdio.Sscanf, c.VarFuncTT(gintT, cstrT, gstrT)),
			c.NewIdent("vsscanf", "stdio.Vsscanf", stdio.Vsscanf, c.FuncTT(gintT, cstrT, gstrT, valistT)),
			c.NewIdent("fscanf", "stdio.Fscanf", stdio.Fscanf, c.VarFuncTT(gintT, filePtr, gstrT)),
			c.NewIdent("vfscanf", "stdio.Vfscanf", stdio.Vfscanf, c.FuncTT(gintT, filePtr, gstrT, valistT)),
			c.NewIdent("remove", "stdio.Remove", stdio.Remove, c.VarFuncTT(gintT, gstrT)),
			c.NewIdent("rename", "stdio.Rename", stdio.Rename, c.VarFuncTT(gintT, gstrT, gstrT)),
		)

		l.Header += `
#define stdout _cxgo_getStdout()
#define stderr _cxgo_getStderr()
#define stdin _cxgo_getStdin()

#define fclose(f) ((FILE*)(f))->Close()
#define feof(f) ((FILE*)(f))->IsEOF()
#define fflush(f) ((FILE*)(f))->Flush()
#define fgetc(f) ((FILE*)(f))->GetC()
#define fgets(buf, sz, f) ((FILE*)(f))->GetS(buf, sz)
#define fileno(f) ((FILE*)(f))->FileNo()
#define fputc(v, f) ((FILE*)(f))->PutC(v)
#define fputs(v, f) ((FILE*)(f))->PutS(v)
#define fread(p, sz, cnt, f) ((FILE*)(f))->ReadN(p, sz, cnt)
#define fwrite(p, sz, cnt, f) ((FILE*)(f))->WriteN(p, sz, cnt)
#define fseek(f, a1, a2) ((FILE*)(f))->Seek(a1, a2)
#define ftell(f) ((FILE*)(f))->Tell()
#define getc(f) ((FILE*)(f))->GetC()
#define ungetc(c, f) ((FILE*)(f))->UnGetC(c)
#define putchar(v) ((FILE*)stdout)->PutC(c)

void     clearerr(FILE *);
char    *ctermid(char *);
_cxgo_go_int      dprintf(_cxgo_go_int, const char *restrict, ...);
int      ferror(FILE *);
int      fgetpos(FILE *restrict, fpos_t *restrict);
void     flockfile(FILE *);
FILE    *fmemopen(void *restrict, size_t, const char *restrict);
FILE    *freopen(const char *restrict, const char *restrict, FILE *restrict);
int      fseeko(FILE *, off_t, int);
int      fsetpos(FILE *, const fpos_t *);
off_t    ftello(FILE *);
int      ftrylockfile(FILE *);
void     funlockfile(FILE *);
int      getchar(void);
int      getc_unlocked(FILE *);
int      getchar_unlocked(void);
ssize_t  getdelim(char **restrict, size_t *restrict, int, FILE *restrict);
ssize_t  getline(char **restrict, size_t *restrict, FILE *restrict);
char    *gets(char *);
FILE    *open_memstream(char **, size_t *);
int      pclose(FILE *);
void     perror(const char *);
FILE    *popen(const char *, const char *);
int      putc(int, FILE *);
int      putc_unlocked(int, FILE *);
int      putchar_unlocked(int);
int      puts(const char *);
int      renameat(int, const char *, int, const char *);
void     rewind(FILE *);
void     setbuf(FILE *restrict, char *restrict);
int      setvbuf(FILE *restrict, char *restrict, int, size_t);
char    *tempnam(const char *, const char *);
FILE    *tmpfile(void);
char    *tmpnam(char *);
int      vdprintf(int, const char *restrict, va_list);

#undef int
`
		return l
	})
}
