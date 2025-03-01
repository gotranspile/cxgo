package cxgo

import "testing"

var casesTranslatePtrs = []parseCase{
	{
		name: "if ptr",
		src: `
void foo() {
	int* a;
	if (a) return;
}
`,
		exp: `
func foo() {
	var a *int32
	if a != nil {
		return
	}
}
`,
	},
	{
		name: "if not ptr",
		src: `
void foo() {
	int* a;
	if (!a) return;
}
`,
		exp: `
func foo() {
	var a *int32
	if a == nil {
		return
	}
}
`,
	},
	{
		name: "if unsafe ptr",
		src: `
void foo() {
	void* a;
	if (a) return;
}
`,
		exp: `
func foo() {
	var a unsafe.Pointer
	if a != nil {
		return
	}
}
`,
	},
	{
		name: "if not unsafe ptr",
		src: `
void foo() {
	void* a;
	if (!a) return;
}
`,
		exp: `
func foo() {
	var a unsafe.Pointer
	if a == nil {
		return
	}
}
`,
	},
	{
		name: "if func",
		src: `
void foo() {
	void(*a)(void);
	if (a) return;
}
`,
		exp: `
func foo() {
	var a func()
	if a != nil {
		return
	}
}
`,
	},
	{
		name: "if not func",
		src: `
void foo() {
	void(*a)(void);
	if (!a) return;
}
`,
		exp: `
func foo() {
	var a func()
	if a == nil {
		return
	}
}
`,
	},
	{
		name: "bool -> ptr",
		src: `
#include <stdbool.h>

void foo() {
	bool a;
	int* b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a bool
		b *int32
	)
	_ = b
	b = (*int32)(unsafe.Pointer(uintptr(libc.BoolToInt(a))))
}
`,
	},
	{
		name:     "uintptr",
		builtins: true,
		src: `
void bar(_cxgo_go_uintptr a) {
	bar(a);
}
void foo(unsigned int a, int b) {
	_cxgo_go_uintptr c;
	c = a;
	c = b;
	bar(a);
	bar(b);
	bar(c);
}
`,
		exp: `
func bar(a uintptr) {
	bar(a)
}
func foo(a uint32, b int32) {
	var c uintptr
	c = uintptr(a)
	c = uintptr(b)
	bar(uintptr(a))
	bar(uintptr(b))
	bar(c)
}
`,
	},
	{
		name: "inc ptr",
		src: `
void foo() {
	char* a;
	a++;
	a--;
}
`,
		exp: `
func foo() {
	var a *byte
	a = (*byte)(unsafe.Add(unsafe.Pointer(a), 1))
	a = (*byte)(unsafe.Add(unsafe.Pointer(a), -1))
}
`,
	},
	{
		name: "unsafe ptr add",
		src: `
void foo() {
	int a = 2;
	void* pa;
	pa = pa + a;
	pa = pa + 3;
	pa += a;
	pa += 5;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa unsafe.Pointer
	)
	pa = unsafe.Add(pa, a)
	pa = unsafe.Add(pa, 3)
	pa = unsafe.Add(pa, a)
	pa = unsafe.Add(pa, 5)
}
`,
	},
	{
		name: "named ptr -> named ptr",
		src: `
typedef struct{} A;
typedef struct{} B;
void foo(A* a) {
	B *b = a;
}
`,
		exp: `
type A struct {
}
type B struct {
}

func foo(a *A) {
	var b *B = (*B)(unsafe.Pointer(a))
	_ = b
}
`,
	},
	{
		name: "cast ptr to func ptr",
		src: `
void foo(void* a) {
	(*(void (**)(void))a)();
}
`,
		exp: `
func foo(a unsafe.Pointer) {
	(*(*func())(a))()
}
`,
	},
	{
		name: "cast ptr to func",
		src: `
void foo(int* a) {
	((void (*)(void))a)();
}
`,
		exp: `
func foo(a *int32) {
	(libc.AsFunc(a, (*func())(nil)).(func()))()
}
`,
	},
	{
		name: "cast ptr to int",
		src: `
void foo(int* a) {
	int b = a;
}
`,
		exp: `
func foo(a *int32) {
	var b int32 = int32(uintptr(unsafe.Pointer(a)))
	_ = b
}
`,
	},
	{
		name: "cast struct ptr to int",
		src: `
typedef struct A {} A;

void foo(A* a) {
	int b = a;
}
`,
		exp: `
type A struct {
}

func foo(a *A) {
	var b int32 = int32(uintptr(unsafe.Pointer(a)))
	_ = b
}
`,
	},
	{
		name: "cast func to int",
		src: `
void foo(int a) {
	a = foo;
}
`,
		exp: `
func foo(a int32) {
	a = int32(libc.FuncAddr(foo))
}
`,
	},
	{
		name: "ptr multi cast",
		src: `
typedef int* PT;

void foo(PT a) {
	char b;
	foo((PT)&b);
}
`,
		exp: `
type PT *int32

func foo(a PT) {
	var b int8
	foo(PT(unsafe.Pointer(&b)))
}
`,
	},
	{
		name: "ptr add",
		src: `
void foo() {
	int a = 2;
	int* pa;
	pa = pa + a;
	pa = pa + 3;
	pa += a;
	pa += 5;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *int32
	)
	pa = (*int32)(unsafe.Add(unsafe.Pointer(pa), unsafe.Sizeof(int32(0))*uintptr(a)))
	pa = (*int32)(unsafe.Add(unsafe.Pointer(pa), unsafe.Sizeof(int32(0))*3))
	pa = (*int32)(unsafe.Add(unsafe.Pointer(pa), unsafe.Sizeof(int32(0))*uintptr(a)))
	pa = (*int32)(unsafe.Add(unsafe.Pointer(pa), unsafe.Sizeof(int32(0))*5))
}
`,
	},
	{
		name: "int add ptr",
		src: `
void foo() {
	int a = 2;
	int* pa;
	pa = a + 4;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *int32
	)
	_ = pa
	pa = (*int32)(unsafe.Pointer(uintptr(a + 4)))
}
`,
	},
	{
		name: "int add short ptr",
		src: `
void foo() {
	int a = 2;
	short* pa;
	pa = a + 748;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *int16
	)
	_ = pa
	pa = (*int16)(unsafe.Pointer(uintptr(a + 748)))
}
`,
	},
	{
		name: "int add short ptr unaligned",
		src: `
void foo() {
	int a = 2;
	short* pa;
	pa = a + 3;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *int16
	)
	_ = pa
	pa = (*int16)(unsafe.Pointer(uintptr(a + 3)))
}
`,
	},
	{
		name: "int add byte ptr",
		src: `
void foo() {
	int a = 2;
	char* pa;
	pa = a + 4;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *byte
	)
	_ = pa
	pa = (*byte)(unsafe.Pointer(uintptr(a + 4)))
}
`,
	},
	{
		name: "int add byte ptr unaligned",
		src: `
void foo() {
	int a = 2;
	char* pa;
	pa = a + 3;
}
`,
		exp: `
func foo() {
	var (
		a  int32 = 2
		pa *byte
	)
	_ = pa
	pa = (*byte)(unsafe.Pointer(uintptr(a + 3)))
}
`,
	},
	{
		name: "int -> unsafe ptr",
		src: `
void foo() {
	int a;
	void* b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a int32
		b unsafe.Pointer
	)
	_ = b
	b = unsafe.Pointer(uintptr(a))
}
`,
	},
	{
		name: "unsafe ptr -> int",
		src: `
void foo() {
	void* a;
	int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a unsafe.Pointer
		b int32
	)
	_ = b
	b = int32(uintptr(a))
}
`,
	},
	{
		name: "int -> ptr",
		src: `
void foo() {
	int a;
	int* b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a int32
		b *int32
	)
	_ = b
	b = (*int32)(unsafe.Pointer(uintptr(a)))
}
`,
	},
	{
		name: "ptr -> int",
		src: `
void foo() {
	int* a;
	int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a *int32
		b int32
	)
	_ = b
	b = int32(uintptr(unsafe.Pointer(a)))
}
`,
	},
	{
		name: "array -> ptr",
		src: `
void foo() {
	char a[5];
	int* b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a [5]byte
		b *int32
	)
	_ = b
	b = (*int32)(unsafe.Pointer(&a[0]))
}
`,
	},
	{
		name: "array -> unsafe ptr",
		src: `
void foo() {
	char a[5];
	void* b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a [5]byte
		b unsafe.Pointer
	)
	_ = b
	b = unsafe.Pointer(&a[0])
}
`,
	},
	{
		name: "store func ptr",
		src: `
void foo() {
	char a[8];
	*(void**)& a[4] = &foo;
}
`,
		exp: `
func foo() {
	var a [8]byte
	*(*unsafe.Pointer)(unsafe.Pointer(&a[4])) = unsafe.Pointer(libc.FuncAddr(foo))
}
`,
	},
	{
		name: "get func pointer",
		src: `
void foo() {
	char a[8];
	(*(void(**)(int)) &a)(1);
}
`,
		exp: `
func foo() {
	var a [8]byte
	(*(*func(int32))(unsafe.Pointer(&a[0])))(1)
}
`,
	},
	{
		name: "call func ptr",
		src: `
void foo(void(***a)(void)) {
	void(**b)(void);
	b = *a;
	(*b)();
}
`,
		exp: `
func foo(a **func()) {
	var b *func()
	b = *a
	(*b)()
}
`,
	},
	{
		name: "post inc and set ptr",
		src: `
void foo() {
	int* a;
	int b;
	*a++ = b;
}
`,
		exp: `
func foo() {
	var (
		a *int32
		b int32
	)
	*func() *int32 {
		p_ := &a
		x := *p_
		*p_ = (*int32)(unsafe.Add(unsafe.Pointer(*p_), unsafe.Sizeof(int32(0))*1))
		return x
	}() = b
}
`,
	},
	{
		name: "ptr decr size",
		src: `
void foo() {
  int *ptr = 0;
  ptr -= 3;
}
`,
		exp: `
func foo() {
	var ptr *int32 = nil
	ptr = (*int32)(unsafe.Add(unsafe.Pointer(ptr), -int(unsafe.Sizeof(int32(0))*3)))
}
`,
	},
	{
		name: "negative ptr",
		src: `
void foo() {
	const char* a;
	a = -1;
	a = -2;
	if (a == (const char*)-1) {
		return;
	}
}
`,
		exp: `
func foo() {
	var a *byte
	a = (*byte)(unsafe.Pointer(uintptr(4294967295)))
	a = (*byte)(unsafe.Pointer(uintptr(4294967294)))
	if uintptr(unsafe.Pointer(a)) == uintptr(4294967295) {
		return
	}
}
`,
	},
	{
		name: "negative void ptr",
		src: `
void foo() {
	void* a;
	a = -1;
	a = (void*)-1;
	a = -2;
	if (a == -1) {
		return;
	}
}
`,
		exp: `
func foo() {
	var a unsafe.Pointer
	a = unsafe.Pointer(uintptr(4294967295))
	a = unsafe.Pointer(uintptr(4294967295))
	a = unsafe.Pointer(uintptr(4294967294))
	if uintptr(a) == uintptr(4294967295) {
		return
	}
}
`,
	},
	{
		name: "return zero ptr",
		src: `
int* foo() {
	return 0;
}
`,
		exp: `
func foo() *int32 {
	return nil
}
`,
	},
	{
		name: "return zero ptr named",
		inc:  `typedef int my_size_t;`,
		src: `
my_size_t* foo() {
	return 0;
}
`,
		exp: `
func foo() *my_size_t {
	return nil
}
`,
	},
	{
		name: "compare unsafe pointers",
		src: `
void foo(void* a, void* b) {
	if (a < b) return;
}
`,
		exp: `
func foo(a unsafe.Pointer, b unsafe.Pointer) {
	if uintptr(a) < uintptr(b) {
		return
	}
}
`,
	},
	{
		name: "compare pointers",
		src: `
void foo(int* a, int* b) {
	if (a < b) return;
}
`,
		exp: `
func foo(a *int32, b *int32) {
	if uintptr(unsafe.Pointer(a)) < uintptr(unsafe.Pointer(b)) {
		return
	}
}
`,
	},
	{
		name: "compare diff pointers",
		src: `
void foo(int* a, void* b) {
	if (a < b) return;
}
`,
		exp: `
func foo(a *int32, b unsafe.Pointer) {
	if uintptr(unsafe.Pointer(a)) < uintptr(b) {
		return
	}
}
`,
	},
	{
		name: "equal unsafe pointers",
		src: `
void foo(void* a, void* b) {
	if (a == b) return;
}
`,
		exp: `
func foo(a unsafe.Pointer, b unsafe.Pointer) {
	if a == b {
		return
	}
}
`,
	},
	{
		name: "equal pointers",
		src: `
void foo(int* a, int* b) {
	if (a == b) return;
}
`,
		exp: `
func foo(a *int32, b *int32) {
	if a == b {
		return
	}
}
`,
	},
	{
		name: "equal diff pointers",
		src: `
void foo(int* a, void* b) {
	if (a == b) return;
}
`,
		exp: `
func foo(a *int32, b unsafe.Pointer) {
	if unsafe.Pointer(a) == b {
		return
	}
}
`,
	},
	{
		name: "equal pointer and const",
		src: `
void foo(int* a, void* b) {
	if (a == 1 && b == 1) return;
}
`,
		exp: `
func foo(a *int32, b unsafe.Pointer) {
	if uintptr(unsafe.Pointer(a)) == uintptr(1) && uintptr(b) == uintptr(1) {
		return
	}
}
`,
	},
	{
		name: "diff pointers",
		src: `
#include <stddef.h>
void foo(int* a, int* b) {
	int c = a - b;
}
`,
		exp: `
func foo(a *int32, b *int32) {
	var c int32 = int32(uintptr(unsafe.Pointer(a)) - uintptr(unsafe.Pointer(b)))
	_ = c
}
`,
	},
	{
		name: "array ptr equal",
		src: `
void foo() {
	int a[10];
	if (a == 0x1000) {
		return;
	}
}
`,
		exp: `
func foo() {
	var a [10]int32
	if uintptr(unsafe.Pointer(&a[0])) == uintptr(0x1000) {
		return
	}
}
`,
	},
	{
		name: "equal func",
		src: `
void foo() {
	void(*a)(void);
	if (a == foo) {
		return;
	}
}
`,
		exp: `
func foo() {
	var a func()
	if libc.FuncAddr(a) == libc.FuncAddr(foo) {
		return
	}
}
`,
	},
	{
		name: "array reinterpret",
		inc:  `unsigned char bytes[124];`,
		src: `
void foo() {
	unsigned char* a;
	a = *(unsigned char**)& bytes[4];
}
`,
		exp: `
func foo() {
	var a *uint8
	_ = a
	a = *(**uint8)(unsafe.Pointer(&bytes[4]))
}
`,
	},
	{
		name: "malloc and free",
		src: `
#include <stdlib.h>
void foo() {
	void* a;
	a = malloc(124);
	free(a);
}
`,
		exp: `
func foo() {
	var a unsafe.Pointer
	_ = a
	a = libc.Malloc(124)
	a = nil
}
`,
	},
	{
		name: "malloc sizeof",
		src: `
#include <stdlib.h>
void foo() {
	int* a = malloc(sizeof(int));
}
`,
		exp: `
func foo() {
	var a *int32 = new(int32)
	_ = a
}
`,
	},
	{
		name: "memset sizeof",
		src: `
#include <stdlib.h>
void foo(int* a) {
	memset(a, 0, sizeof(int));
}
`,
		exp: `
func foo(a *int32) {
	*a = 0
}
`,
	},
	{
		skip: true, // TODO
		name: "memset slice",
		src: `
#include <stdlib.h>
void foo(int* a) {
	memset(a, 0, 10*sizeof(int));
}
`,
		exp: `
func foo(a []int32) {
	copy(a[:10], make([]int32, 10))
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "a", Type: HintSlice}),
		},
	},
	{
		name: "memset sizeof struct",
		src: `
#include <stdlib.h>

typedef struct {
	int x;
} A;

void foo(A* a) {
	memset(a, 0, sizeof(A));
}
`,
		exp: `
type A struct {
	X int32
}

func foo(a *A) {
	*a = A{}
}
`,
	},
	{
		name: "array ptr assign",
		src: `
void foo() {
	char bytes[10];
	*(void**)& bytes[1] = &bytes[2];
}
`,
		exp: `
func foo() {
	var bytes [10]byte
	*(*unsafe.Pointer)(unsafe.Pointer(&bytes[1])) = unsafe.Pointer(&bytes[2])
}
`,
	},
	{
		name: "addr of field",
		src: `
typedef struct {
	int field;
} A;

void foo() {
	A a;
	int* b;
	b = &a.field;
}
`,
		exp: `
type A struct {
	Field int32
}

func foo() {
	var (
		a A
		b *int32
	)
	_ = b
	b = &a.Field
}
`,
	},
	{
		name: "implicit 0 field ptr",
		src: `
typedef struct {
	int x;
} A;
typedef struct {
	A a;
	A b;
} B;

void bar(A* arg) {
	int x = arg->x;
}

void foo(B* arg) {
	bar(arg);
}
`,
		exp: `
type A struct {
	X int32
}
type B struct {
	A A
	B A
}

func bar(arg *A) {
	var x int32 = arg.X
	_ = x
}
func foo(arg *B) {
	bar(&arg.A)
}
`,
	},
	{
		name: "implicit array access",
		src: `
struct s {
    int i;
};

struct s ss[] = {
    {0},
};


void foo() {
    int i = ss->i;
}
`,
		exp: `
type s struct {
	I int32
}

var ss [1]s = [1]s{}

func foo() {
	var i int32 = ss[0].I
	_ = i
}
`,
	},
	{
		name: "named ptr call",
		src: `
typedef int H;
typedef struct F {} F;
void foo(F* a);
void bar(void) {
	H a;
	foo(((F*)a));
}
`,
		exp: `
type H int32
type F struct {
}

func foo(a *F)
func bar() {
	var a H
	foo((*F)(unsafe.Pointer(uintptr(a))))
}
`,
	},
	{
		name: "sizeof int conv",
		src: `
#include <stddef.h>

typedef struct {
	int x;
} A;

void foo () {
	int x = sizeof(A);
}
`,
		exp: `
type A struct {
	X int32
}

func foo() {
	var x int32 = int32(unsafe.Sizeof(A{}))
	_ = x
}
`,
	},
	{
		name: "ptr add size mult",
		src: `
#include <stddef.h>

typedef struct {
	int x;
	int y;
} A;

void foo() {
	int* ptr;
	ptr += 20 * (sizeof(A)/2);
}
`,
		exp: `

type A struct {
	X int32
	Y int32
}

func foo() {
	var ptr *int32
	ptr = (*int32)(unsafe.Add(unsafe.Pointer(ptr), unsafe.Sizeof(int32(0))*(20*(unsafe.Sizeof(A{})/2))))
}
`,
	},
	{
		skip: true, // TODO
		name: "ptr compare out of bounds",
		src: `
void foo() {
	void* p;
	int arr[9];
	if (p > &arr[9]) {
		foo();
	}
}
`,
		exp: `
func foo() {
	var (
		p   unsafe.Pointer
		arr [9]int32
	)
	if uintptr(unsafe.Pointer(p)) - uintptr(unsafe.Pointer(&arr)) > unsafe.Sizeof(int32(0))*9 {
		foo()
	}
}
`,
	},
	{
		name: "pointer to first elem",
		src: `
void foo() {
	char b[10];
	char* p;
	p = b;
	p = &b;
	p = &b[0];
}
`,
		exp: `
func foo() {
	var (
		b [10]byte
		p *byte
	)
	_ = p
	p = &b[0]
	p = &b[0]
	p = &b[0]
}
`,
	},
	{
		name: "slice override global",
		src: `
int* foo;
`,
		exp: `
var foo []int32
`,
		configFuncs: []configFunc{
			withIdent(IdentConfig{
				Name: "foo",
				Type: HintSlice,
			}),
		},
	},
	{
		name: "slice override global struct",
		skip: true, // FIXME
		src: `
struct bar {
	int a;
};
struct bar* foo;
`,
		exp: `
type bar struct {
	A int32
}
var foo []bar
`,
		configFuncs: []configFunc{
			withIdent(IdentConfig{
				Name: "foo",
				Type: HintSlice,
			}),
		},
	},
	{
		name: "slice override func arg",
		src: `
void foo(int* a);
void bar(int* b){
}
`,
		exp: `
func foo(a []int32)
func bar(b []int32) {
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "a", Type: HintSlice}),
			withIdentField("bar", IdentConfig{Name: "b", Type: HintSlice}),
		},
	},
	{
		name: "slice override struct field",
		src: `
typedef struct{
	int* a
} foo;

typedef struct bar {
	int* b
} bar;

struct baz {
	int* c
};
`,
		exp: `
type foo struct {
	A []int32
}
type bar struct {
	B []int32
}
type baz struct {
	C []int32
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "a", Type: HintSlice}),
			withIdentField("bar", IdentConfig{Name: "b", Type: HintSlice}),
			withIdentField("baz", IdentConfig{Name: "c", Type: HintSlice}),
		},
	},
	{
		name: "slice and ptr",
		src: `
void foo(int* a, int* b, int c) {
	a = b;
	b = a;
	c = *a;
	c = a[1];
	if (a != 0) {
		a = 0;
	}
}
`,
		exp: `
func foo(a []int32, b *int32, c int32) {
	a = []int32(b)
	b = &a[0]
	c = a[0]
	c = a[1]
	if a != nil {
		a = nil
	}
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "a", Type: HintSlice}),
		},
	},
	{
		name: "slice calloc",
		src: `
#include <stdlib.h>

int n;
int* a = calloc(n, sizeof(int));
int* b = (int*)calloc(n, sizeof(int));
int* c = (int*)calloc(n, sizeof(int));
`,
		exp: `
var n int32
var a []int32 = make([]int32, int(n))
var b []int32 = make([]int32, int(n))
var c *int32 = &make([]int32, int(n))[0]
`,
		configFuncs: []configFunc{
			withIdent(IdentConfig{Name: "a", Type: HintSlice}),
			withIdent(IdentConfig{Name: "b", Type: HintSlice}),
		},
	},
	{
		name: "slice safe calloc",
		src: `
#include <stdlib.h>

void foo(int* a, int n) {
	if ((a = (int*)calloc(n, sizeof(int))) != 0) {
		a = 0;
	}
}
`,
		exp: `
func foo(a []int32, n int32) {
	if (func() []int32 {
		a = make([]int32, int(n))
		return a
	}()) != nil {
		a = nil
	}
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "a", Type: HintSlice}),
		},
	},
	{
		name: "slice index",
		src: `
typedef struct {
	int* ptr
} A;

void foo(A x, int y) {
	y = *x.ptr;
	y = x.ptr[1];
	y = *(x.ptr + 1);
	y = *(x.ptr + 2*y + 1);
	y = *(&(x.ptr + 2*y + 1)[y]);
	y = *(x.ptr + 2 - 1);
}
`,
		exp: `
type A struct {
	Ptr []int32
}

func foo(x A, y int32) {
	y = x.Ptr[0]
	y = x.Ptr[1]
	y = x.Ptr[1]
	y = x.Ptr[y*2+1]
	y = x.Ptr[y*2+1+y]
	y = x.Ptr[2-1]
}
`,
		configFuncs: []configFunc{
			withIdentField("A", IdentConfig{Name: "ptr", Type: HintSlice}),
		},
	},
	{
		name: "slice to bytes",
		src: `
void foo(unsigned char* x) {
}
void foo2(char* x) {
	foo(x);
}
`,
		exp: `
func foo(x []byte) {
}
func foo2(x string) {
	foo([]byte(x))
}
`,
		configFuncs: []configFunc{
			withIdentField("foo", IdentConfig{Name: "x", Type: HintSlice}),
			withIdentField("foo2", IdentConfig{Name: "x", Type: HintString}),
		},
	},
	{
		name:     "slice def index",
		builtins: true,
		src: `
typedef struct {
	_cxgo_go_slice_t(int) ptr
} A;

void foo(A x, int y) {
	y = *x.ptr;
	y = x.ptr[1];
	y = *(x.ptr + 1);
	y = *(x.ptr + 2*y + 1);
	y = *(&(x.ptr + 2*y + 1)[y]);
}
`,
		exp: `
type A struct {
	Ptr []int32
}

func foo(x A, y int32) {
	y = x.Ptr[0]
	y = x.Ptr[1]
	y = x.Ptr[1]
	y = x.Ptr[y*2+1]
	y = x.Ptr[y*2+1+y]
}
`,
	},
	{
		name:     "slice def ops",
		builtins: true,
		src: `
void foo(int y) {
	_cxgo_go_slice_t(int) arr;
	arr = 0;
	arr = _cxgo_go_make(_cxgo_go_slice_t(int), 1);
	arr = _cxgo_go_make(_cxgo_go_slice_t(int), 1, 2);
	arr = _cxgo_go_make_same(arr, 3, 4);
	y = _cxgo_go_len(arr);
	y = _cxgo_go_cap(arr);
	arr = _cxgo_go_slice(arr, -1, -1);
	arr = _cxgo_go_slice(arr, -1, 2);
	arr = _cxgo_go_slice(arr, 1, -1);
	arr = _cxgo_go_slice(arr, 1, 2);
	arr = _cxgo_go_slice(arr, -1, -1, -1);
	arr = _cxgo_go_slice(arr, -1, 2, 3);
	arr = _cxgo_go_slice(arr, 1, 2, 3);
	arr = _cxgo_go_append(arr, 1);
	arr = _cxgo_go_append(arr, 1, 2);
	arr = _cxgo_go_append(arr, arr);
}
`,
		exp: `
func foo(y int32) {
	var arr []int32
	arr = nil
	arr = make([]int32, 1)
	arr = make([]int32, 1, 2)
	arr = make([]int32, 3, 4)
	y = int32(len(arr))
	y = int32(cap(arr))
	arr = arr[:]
	arr = arr[:2]
	arr = arr[1:]
	arr = arr[1:2]
	arr = arr[:]
	arr = arr[:2:3]
	arr = arr[1:2:3]
	arr = append(arr, 1)
	arr = append(arr, 1, 2)
	arr = append(arr, arr...)
}
`,
	},
}

func TestPointers(t *testing.T) {
	runTestTranslate(t, casesTranslatePtrs)
}
