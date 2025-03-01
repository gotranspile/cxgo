package cxgo

import (
	"testing"

	"github.com/gotranspile/cxgo/types"
)

var casesTranslateDecls = []parseCase{
	{
		name: "typedef primitive",
		src:  `typedef int int_t;`,
		exp:  `type int_t int32`,
	},
	{
		name: "function decl",
		src: `
void foo() {}
`,
		exp: `
func foo() {
}
`,
	},
	{
		name: "function forward decl",
		src: `
void foo1();
int foo2();
int foo3(int*);
int (*a)();
int (*b)(int*);
`,
		exp: `
func foo1()
func foo2() int32
func foo3(*int32) int32

var a func() int32
var b func(*int32) int32
`,
	},
	{
		name: "function forward decl 2",
		src: `
void foo1();
int foo2();
int foo3(int);
void foo1() {}
int foo2() {}
int foo3(int) {}
`,
		exp: `
func foo1() {
}
func foo2() int32 {
}
func foo3(int32) int32 {
}
`,
	},
	{
		name: "var",
		src: `
int foo;
char byte;
`,
		exp: `
var foo int32
var byte_ int8
`,
	},
	{
		name: "var init",
		src: `
int foo = 1;
`,
		exp: `
var foo int32 = 1
`,
	},
	{
		name: "multiple vars",
		src: `
int foo = 1, *bar;
`,
		exp: `
var foo int32 = 1
var bar *int32
`,
	},
	{
		name: "multiple vars 2",
		src: `
int a = 1, *p, f(void), (*pf)(double);
`,
		exp: `
var a int32 = 1
var p *int32

func f() int32

var pf func(float64) int32
`,
	},
	{
		name: "complex var",
		src: `
int (*(*foo)(double))[3] = 0;
`,
		exp: `
var foo func(float64) *[3]int32 = nil
`,
	},
	{
		name: "function var",
		src: `
void (*foo)(void);
`,
		exp: `
var foo func()
`,
	},
	{
		name: "struct forward decl",
		src: `
struct foo;
struct foo2;

struct foo {
	int a;
};
`,
		exp: `
type foo struct {
	A int32
}
type foo2 struct {
}
`,
	},
	{
		name: "typedef struct",
		inc: `
typedef struct bar {
	int d;
} BAR;
`,
		src: `
typedef struct { int a; } foo;
typedef struct foo2 { int b; } foo2;
typedef struct foo3 { int c; } FOO3;
struct bar bar1;
BAR bar2;
`,
		exp: `
type foo struct {
	A int32
}
type foo2 struct {
	B int32
}
type foo3 struct {
	C int32
}
type FOO3 foo3

var bar1 bar
var bar2 BAR
`,
	},
	{
		name: "typedef struct 3",
		src: `
struct baz { int d; } bar3;
`,
		exp: `
type baz struct {
	D int32
}

var bar3 baz
`,
	},
	{
		name: "typedef alias",
		src: `
typedef struct { int a; } A;
typedef A B;
typedef B C;

struct T1 {
  C (*f1)[3];
  C *f2;
};
typedef struct T1 T2;

void foo(C* c) {}
`,
		exp: `
type A struct {
	A int32
}
type T1 struct {
	F1 *[3]A
	F2 *A
}
type T2 T1

func foo(c *A) {
}
`,
		configFuncs: []configFunc{
			withAlias("B"),
			withAlias("C"),
		},
	},
	{
		name: "typedef alias 2",
		inc: `
typedef struct { int a; } A;
typedef A B;
typedef B C;
`,
		src: `
struct T1 {
  C (*f1)[3];
  C *f2;
};
typedef struct T1 T2;

void foo(C* c) {}
`,
		exp: `
type T1 struct {
	F1 *[3]A
	F2 *A
}
type T2 T1

func foo(c *A) {
}
`,
		configFuncs: []configFunc{
			withAlias("B"),
			withAlias("C"),
		},
	},
	{
		name: "recursive struct",
		inc: `
typedef struct _A* HA;
typedef struct _B* HB;
`,
		src: `
struct _A {
	HB b;
};
struct _B {
	HA a;
};
`,
		exp: `
type _A struct {
	B HB
}
type _B struct {
	A HA
}
`,
	},
	{
		name: "named enum",
		src: `
	enum Enum
	{
	   VALUE_1,
	   VALUE_2,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)
	`,
	},
	{
		name: "forward enum",
		src: `
	enum Enum;
	enum Enum
	{
	   VALUE_1,
	   VALUE_2,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)
	`,
	},
	{
		name: "return enum",
		src: `
	enum Enum
	{
	   VALUE_1,
	   VALUE_2,
	};
	extern enum Enum foo();
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)

func foo() Enum
	`,
	},
	{
		name: "unnamed enum",
		src: `
enum
{
    VALUE_1,
    VALUE_2,
};
`,
		exp: `
const (
	VALUE_1 = iota
	VALUE_2
)
`,
	},
	{
		name: "typedef enum",
		src: `
typedef enum {
    VALUE_1,
    VALUE_2,
} Enum;
`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)
`,
	},
	{
		name: "enum zero",
		src: `
	enum Enum
	{
	   VALUE_1 = 0,
	   VALUE_2,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)
	`,
	},
	{
		name: "enum neg",
		src: `
	enum {
	   VALUE_1 = -1,
	   VALUE_2 = 0,
	   VALUE_3,
	};
	`,
		exp: `

const (
	VALUE_1 = -1
	VALUE_2 = 0
	VALUE_3 = 1
)
	`,
	},
	{
		name: "enum start",
		src: `
	enum Enum
	{
	   VALUE_1 = 1,
	   VALUE_2,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota + 1)
	VALUE_2
)
	`,
	},
	{
		name: "enum fixed",
		src: `
	enum Enum
	{
	   VALUE_1 = 1,
	   VALUE_2 = 2,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 Enum = 1
	VALUE_2 Enum = 2
)
	`,
	},
	{
		name: "enum no zero",
		src: `
	enum Enum
	{
	   VALUE_1,
	   VALUE_2 = 1,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 Enum = 0
	VALUE_2 Enum = 1
)
	`,
	},
	{
		name: "enum no zero 2",
		src: `
	enum Enum
	{
	   VALUE_1,
	   VALUE_2 = 42,
	   VALUE_3,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 Enum = 0
	VALUE_2 Enum = 42
	VALUE_3 Enum = 43
)
	`,
	},
	{
		name: "enum negative",
		src: `
	enum Enum
	{
	   VALUE_1 = -3,
	   VALUE_2,
	   VALUE_3 = 1,
	};
	`,
		exp: `
type Enum int32

const (
	VALUE_1 Enum = -3
	VALUE_2 Enum = -2
	VALUE_3 Enum = 1
)
	`,
	},
	{
		name: "use enum",
		src: `
enum Enum
{
   VALUE_1,
   VALUE_2,
};
enum Enum foo() {
	return VALUE_1;
}
	`,
		exp: `
type Enum int32

const (
	VALUE_1 = Enum(iota)
	VALUE_2
)

func foo() Enum {
	return VALUE_1
}
	`,
	},
	{
		name: "enum in func",
		src: `
void foo() {
	typedef enum { A, B, C } Enum;
	typedef enum { D, E, F } Enum2;
	Enum x;
}
	`,
		exp: `
func foo() {
	type Enum int32
	const (
		A = Enum(iota)
		B
		C
	)
	type Enum2 int32
	const (
		D = Enum2(iota)
		E
		F
	)
	var x Enum
	_ = x
}
	`,
	},
	{
		name: "struct and func", skip: true, // TODO
		src: `
struct foo;
void foo();
`,
		exp: `
type foo struct {
}

func foo_2()
`,
	},
	{
		name: "struct and var",
		inc: `
struct foo {};
`,
		src: `
struct foo foo;
`,
		exp: `
var foo foo
`,
	},
	{
		name: "func arg",
		src: `
void foo(int (*a)(void)) {
}
`,
		exp: `
func foo(a func() int32) {
}
`,
	},
	{
		name: "local func var",
		src: `
void foo() {
	int (*a)();
}
`,
		exp: `
func foo() {
	var a func() int32
	_ = a
}
`,
	},
	{
		name: "for init",
		src: `
void foo() {
	for (int i = 0; i < 5; i++) {}
}
`,
		exp: `
func foo() {
	for i := int32(0); i < 5; i++ {
	}
}
`,
	},
	{
		name: "for init multiple",
		skip: true, // TODO
		src: `
	void foo() {
		for (int i = 0, j = 1; i < 5; i++) {}
	}
	`,
		exp: `
	func foo() {
		for i, j := int32(0), int32(1); i < 5; i++ {
		}
	}
	`,
	},
	{
		name: "for init ternary define",
		src: `
void foo(int x) {
	for (int i = x ? 1 : 0; i < 5; i++) {}
}
	`,
		exp: `
func foo(x int32) {
	for i := int32(func() int32 {
		if x != 0 {
			return 1
		}
		return 0
	}()); i < 5; i++ {
	}
}
	`,
	},
	{
		name: "for init ternary reuse",
		src: `
void foo(int x) {
	int i;
	for (i = x ? 1 : 0; i < 5; i++) {}
}
	`,
		exp: `
func foo(x int32) {
	var i int32
	for i = func() int32 {
		if x != 0 {
			return 1
		}
		return 0
	}(); i < 5; i++ {
	}
}
	`,
	},
	{
		name: "extern var",
		src: `
extern int a;
`,
		exp: `
`,
	},
	{
		name: "use incomplete type",
		src: `
typedef struct MyType MyType;
MyType *new_type(void) {
	return 0;
}
`,
		exp: `
type MyType struct {
}

func new_type() *MyType {
	return nil
}
`,
	},
	{
		name: "unnamed struct var",
		src: `
void foo() {
	struct{
		int field;
	} a = {0};
	struct{
		int field;
		int field2;
	} b = {0};
}
`,
		exp: `
func foo() {
	var a struct {
		Field int32
	} = struct {
		Field int32
	}{}
	_ = a
	var b struct {
		Field  int32
		Field2 int32
	} = struct {
		Field  int32
		Field2 int32
	}{}
	_ = b
}
`,
	},
	{
		name: "empty array decl",
		src: `
void foo() {
	int a[1][0];
}
`,
		exp: `
func foo() {
	var a [1][0]int32
	_ = a
}
`,
	},
	{
		name: "dyn array arg",
		src: `
void foo(int a[]) {
}
`,
		exp: `
func foo(a []int32) {
}
`,
	},
	{
		name: "rename decl func",
		src: `
void foo() {}
`,
		exp: `
func Bar() {
}
`,
		configFuncs: []configFunc{
			withRename("foo", "Bar"),
		},
	},
	{
		name: "rename decl struct",
		src: `
struct foo {};
`,
		exp: `
type Bar struct {
}
`,
		configFuncs: []configFunc{
			withRename("foo", "Bar"),
		},
	},
	{
		name: "args partially named",
		src: `
void (*foo)(void *a, int, const char *);
`,
		exp: `
var foo func(a unsafe.Pointer, a2 int32, a3 *byte)
`,
	},
	{
		name: "go ints",
		src: `
typedef unsigned int word;
`,
		exp: `
type word uint
`,
		envFuncs: []envFunc{
			func(c *types.Config) {
				c.UseGoInt = true
			},
		},
	},
	{
		name: "unused vars",
		src: `
void foo(int x) {
	int a;
}

void bar() {
	int a = 0;
	int b;
	int c = b;
}
`,
		exp: `
func foo(x int32) {
	var a int32
	_ = a
}
func bar() {
	var a int32 = 0
	_ = a
	var b int32
	var c int32 = b
	_ = c
}
`,
	},
	{
		name: "tcc 10",
		src: `
struct z
{
   int a;
} foo;
`,
		exp: `
type z struct {
	A int32
}

var foo z
`,
	},
	{
		name: "stdlib forward decl",
		src: `
int printf(const char*, ...);

void foo() {
	printf("%d\n", 1);
}
`,
		exp: `
func foo() {
	stdio.Printf("%d\n", 1)
}
`,
	},
	{
		name: "macro empty",
		src: `
#define MY_DEF

int a;
`,
		exp: `
var a int32
`,
	},
	{
		name: "macro untyped int",
		src: `
#define MY_CONST 1

int a = MY_CONST;
`,
		exp: `
const MY_CONST = 1

var a int32 = MY_CONST
`,
	},
	{
		// TODO: cc.AST.Eval() doesn't support cast expressions?
		skip: true,
		name: "macro typed int",
		src: `
#define MY_CONST ((int)1)

int a = MY_CONST;
`,
		exp: `
const MY_CONST = int32(1)

var a int32 = MY_CONST
`,
	},
	{
		name: "macro string",
		src: `
#define MY_CONST "abc"

char* a = MY_CONST;
`,
		exp: `
const MY_CONST = "abc"

var a *byte = libc.CString(MY_CONST)
`,
	},
	{
		name: "macro order",
		src: `
#define MY_CONST 1

int a = MY_CONST;

#define MY_CONST_2 2
`,
		exp: `
const MY_CONST = 1

var a int32 = MY_CONST

const MY_CONST_2 = 2
`,
		// TODO: we don't handle order yet
		skipExp: `
const MY_CONST = 1
const MY_CONST_2 = 2

var a int32 = MY_CONST
`,
	},
	{
		name: "tmp var names",
		src: `
typedef struct {
	int x, y;
} vec_t;

void foo(vec_t* p) {
	int x;
	x = p->x = p->y = 0;
}
`,
		exp: `
type vec_t struct {
	X int32
	Y int32
}

func foo(p *vec_t) {
	var x int32
	_ = x
	x = func() int32 {
		p_ := &p.X
		*p_ = func() int32 {
			p_ := &p.Y
			*p_ = 0
			return *p_
		}()
		return *p_
	}()
}
`,
	},
}

func TestDecls(t *testing.T) {
	runTestTranslate(t, casesTranslateDecls)
}
