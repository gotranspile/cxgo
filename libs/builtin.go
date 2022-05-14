package libs

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/gotranspile/cxgo/types"
)

// TODO: spec URL

const (
	BuiltinH = `cxgo_builtin.h`

	typeFixedIntPref = "_cxgo_"

	RuntimeOrg         = "github.com/gotranspile"
	RuntimePackage     = RuntimeOrg + "/cxgo"
	RuntimePackageVers = "main"
	RuntimePrefix      = RuntimePackage + "/runtime/"
	RuntimeLibc        = RuntimePrefix + "libc"
)

func buildinFixedIntName(sz int, unsigned bool) string {
	return typeFixedIntPref + goIntName(sz, unsigned)
}

func goIntName(sz int, unsigned bool) string {
	name := "sint"
	if unsigned {
		name = "uint"
	}
	return name + strconv.Itoa(sz)
}

func init() {
	RegisterLibrary(BuiltinH, func(c *Env) *Library {
		int32T := types.IntT(4)
		charP := c.C().String()
		lsz := c.C().Long().Sizeof()
		psz := c.PtrT(nil).Sizeof()

		// builtin variable arguments list
		ifaceT := c.Go().Iface()
		sliceT := c.Go().IfaceSlice()
		valistPtr := c.PtrT(nil)
		valistT := types.NamedTGo("__builtin_va_list", "libc.ArgList", c.MethStructT(map[string]*types.FuncType{
			"Start": c.FuncTT(nil, ifaceT, sliceT),
			"Arg":   c.FuncTT(ifaceT),
			"End":   c.FuncTT(nil),
		}))
		valistPtr.SetElem(valistT)

		// memory model
		pre := ""
		switch {
		case lsz == 4 && psz == 8:
			pre += `
#define __LLP64__ 1
#define _LLP64_ 1
`
		case lsz == 8 && psz == 8:
			pre += `
#define __LP64__ 1
#define _LP64_ 1
`
		case psz == 4:
			pre += `
#define __ILP32__ 1
#define _ILP32_ 1
`
		}
		var post strings.Builder
		maxIntTypeDefs(&post, "ptr", c.PtrSize()*8)
		post.WriteString(`
#define NULL 0
typedef intptr_t ptrdiff_t;
typedef uintptr_t size_t;
`)

		l := &Library{
			Header: pre + `
#define __ORDER_BIG_ENDIAN__ 4321
#define __ORDER_LITTLE_ENDIAN__ 1234
#define __BYTE_ORDER__ __ORDER_LITTLE_ENDIAN__

#define __CXGO__
#define __linux__
#define _CXGO_WINAPI

#define __SIZEOF_INT8__ 1
#define __SIZEOF_INT16__ 2
#define __SIZEOF_INT32__ 4
#define __SIZEOF_INT64__ 8

#define _cxgo_int8  __int8
#define _cxgo_int16 __int16
#define _cxgo_int32 __int32
#define _cxgo_int64 __int64

#define _cxgo_sint8  signed __int8
#define _cxgo_sint16 signed __int16
#define _cxgo_sint32 signed __int32
#define _cxgo_sint64 signed __int64

#define _cxgo_uint8  unsigned __int8
#define _cxgo_uint16 unsigned __int16
#define _cxgo_uint32 unsigned __int32
#define _cxgo_uint64 unsigned __int64

#define _cxgo_float32 float
#define _cxgo_float64 double

typedef _Bool _cxgo_go_bool;
typedef _cxgo_uint8 _cxgo_go_byte;
typedef _cxgo_int32 _cxgo_go_rune;
typedef signed int _cxgo_go_int;
typedef unsigned int _cxgo_go_uint;
typedef unsigned int _cxgo_go_uintptr;
typedef void* _cxgo_go_unsafeptr;

typedef struct {
	_cxgo_go_unsafeptr ptr; 
	_cxgo_go_uintptr len_;
} _cxgo_go_string;

typedef struct {
	_cxgo_go_unsafeptr ptr; 
	_cxgo_go_uintptr len_;
	_cxgo_go_uintptr cap_;
} _cxgo_go_slice;

typedef struct{
	_cxgo_go_uintptr typ; 
	_cxgo_go_unsafeptr ptr; 
} _cxgo_go_iface;

typedef _cxgo_go_slice _cxgo_go_iface_slice;

typedef struct __builtin_va_list __builtin_va_list;
typedef struct __builtin_va_list {
	void (*Start) (_cxgo_go_uint n, _cxgo_go_iface_slice rest);
	_cxgo_go_iface (*Arg) ();
	void (*End) ();
} __builtin_va_list;
#define __gnuc_va_list __builtin_va_list

// a hack for C parser to resolve macro references to _rest arg
_cxgo_go_iface_slice _rest;

void* malloc(_cxgo_go_int);

int __predefined_declarator;

` + post.String() + includeHacks,
			Idents: map[string]*types.Ident{
				"malloc": c.C().MallocFunc(),
			},
			Imports: map[string]string{
				"libc":   RuntimeLibc,
				"stdio":  RuntimePrefix + "stdio", // for printf
				"atomic": "sync/atomic",
			},
			Types: map[string]types.Type{
				"__builtin_va_list": valistT,
				"intptr_t":          c.IntPtrT(),
				"uintptr_t":         c.UintPtrT(),
			},
			ForceMacros: map[string]bool{
				"NULL": true,
			},
		}
		for _, t := range c.Go().Types() {
			if nt, ok := t.(types.Named); ok {
				l.Types[nt.Name().Name] = nt
			}
		}
		l.Declare(
			c.Go().LenFunc(),
			c.Go().PanicFunc(),
			c.C().MallocFunc(),
			c.C().MemmoveFunc(),
			c.C().MemcpyFunc(),
			c.C().MemsetFunc(),
			c.NewIdent("_cxgo_func_name", "libc.FuncName", libc.FuncName, c.FuncTT(c.Go().String())),
			c.NewIdent("__builtin_strcpy", "libc.StrCpy", libc.StrCpy, c.FuncTT(charP, charP, charP)),
			c.NewIdent("__sync_lock_test_and_set", "atomic.SwapInt32", atomic.SwapInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_add", "libc.LoadAddInt32", libc.LoadAddInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_sub", "libc.LoadSubInt32", libc.LoadSubInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_or", "libc.LoadOrInt32", libc.LoadOrInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_and", "libc.LoadAndInt32", libc.LoadAndInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_xor", "libc.LoadXorInt32", libc.LoadXorInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__sync_fetch_and_nand", "libc.LoadNandInt32", libc.LoadNandInt32, c.FuncTT(int32T, c.PtrT(int32T), int32T)),
			c.NewIdent("__builtin_printf", "stdio.Printf", stdio.Printf, c.VarFuncTT(c.Go().Int(), c.Go().String())),
			c.NewIdent("_cxgo_va_copy", "libc.ArgCopy", libc.ArgCopy, c.FuncTT(nil, valistPtr, valistPtr)),
			c.NewIdent("printf", "stdio.Printf", stdio.Printf, c.VarFuncTT(c.Go().Int(), c.Go().String())),
		)
		l.Header += `
_cxgo_go_int _cxgo_offsetof(_cxgo_go_iface, _cxgo_go_string);
#define __builtin_offsetof(type, member) _cxgo_offsetof((type)0, "#member")
#define offsetof(type, member) _cxgo_offsetof((type)0, "#member")

#define __builtin_abort() _cxgo_go_panic("abort")
#define __builtin_trap() _cxgo_go_panic("trap")
#define __builtin_unreachable() _cxgo_go_panic("unreachable")
#define memcpy __builtin_memcpy
#define memset __builtin_memset
#define malloc __builtin_malloc

#define __sync_fetch_and_sub(p, v) __sync_fetch_and_add(p, -(v))

#define __builtin_va_start(va, t) va.Start(t, _rest)
#define __builtin_va_arg(va, typ) ((typ)(va.Arg()))
#define __builtin_va_end(va) va.End()
#define __builtin_va_copy(dst, src) _cxgo_va_copy(&dst, &src)
void __builtin_va_arg_pack();

#define _Static_assert(x, y) /* x, y */
#define __builtin_expect(x, y) ((x) == (y))

#define NULL 0
`
		l.Header += fmt.Sprintf("#define __SIZE_TYPE__ _cxgo_uint%d\n", c.PtrSize()*8)
		l.Header += fmt.Sprintf("#define __PTRDIFF_TYPE__ _cxgo_int%d\n", c.PtrSize()*8)
		l.Header += fmt.Sprintf("#define __INTPTR_TYPE__ _cxgo_int%d\n", c.PtrSize()*8)
		return l
	})
}

const includeHacks = `
#define __func__ _cxgo_func_name()
#define __FUNCTION__ _cxgo_func_name()
#define __PRETTY_FUNCTION__ _cxgo_func_name()

#define __aligned(x) /* x */
#define __packed
#define _asm int

// asm volatile("xxx")
#define volatile(x) (x)
`
