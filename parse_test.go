package cxgo

import (
	"bytes"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"modernc.org/cc/v3"

	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/types"
)

const testPkg = "lib"

type configFunc func(c *Config)
type envFunc func(c *types.Config)

type parseCase struct {
	name        string
	inc         string
	src         string
	exp         string
	skipExp     string
	skip        bool
	builtins    bool
	configFuncs []configFunc
	envFuncs    []envFunc
}

func (c parseCase) shouldSkip() bool {
	return c.skip || c.skipExp != ""
}

func withIdent(ic IdentConfig) configFunc {
	return func(c *Config) {
		c.Idents = append(c.Idents, ic)
	}
}

func withIdentField(name string, f IdentConfig) configFunc {
	return withIdent(IdentConfig{Name: name, Fields: []IdentConfig{f}})
}

func withAlias(name string) configFunc {
	return func(c *Config) {
		c.Idents = append(c.Idents, IdentConfig{Name: name, Alias: true})
	}
}

func withRename(from, to string) configFunc {
	return func(c *Config) {
		c.Idents = append(c.Idents, IdentConfig{Name: from, Rename: to})
	}
}

func withDoNotEdit(val bool) configFunc {
	return func(c *Config) {
		c.DoNotEdit = val
	}
}

var casesTranslate = []parseCase{
	{
		name: "empty",
		src:  " ",
	},
	{
		name: "push defines",
		src: `
#define BLAH 10

int v1 = BLAH;

#pragma push_macro("BLAH")
#undef BLAH
#define BLAH 5

int v2 = BLAH;

#pragma pop_macro("BLAH")

int v3 = BLAH;

#pragma push_macro("NON_EXISTENT")
#pragma pop_macro("NON_EXISTENT")
`,
		exp: `
const BLAH = 10
const BLAH = 5

var v1 int32 = BLAH
var v2 int32 = BLAH
var v3 int32 = BLAH
`,
		skipExp: `
const BLAH = 10

var v1 int32 = BLAH
var v2 int32 = BLAH
var v3 int32 = BLAH
`,
	},
	{
		name: "switch",
		src: `
void foo(int a) {
	switch (a) {
	case 1:
		foo(1);
		break;
	case 2:
		foo(2);
	default:
		foo(0);
	case 3:
		foo(3);
		break;
	case 4:
	case 5:
		foo(5);
		return;
	case 6:
		foo(6);
	}
}
`,
		exp: `
func foo(a int32) {
	switch a {
	case 1:
		foo(1)
	case 2:
		foo(2)
		fallthrough
	default:
		foo(0)
		fallthrough
	case 3:
		foo(3)
	case 4:
		fallthrough
	case 5:
		foo(5)
		return
	case 6:
		foo(6)
	}
}
`,
	},
	{
		name: "switch unreachable",
		src: `
void foo(int a) {
	switch (a) {
		a++;
	case 1:
		foo(1);
		break;
		a++;
	case 2:
		foo(2);
	default:
		foo(0);
	case 3:
		foo(3);
		break;
	case 4:
	case 5:
		foo(5);
		return;
		a++;
	case 6:
		foo(6);
	}
}
`,
		exp: `
func foo(a int32) {
	switch a {
	case 1:
		foo(1)
		break
		a++
		fallthrough
	case 2:
		foo(2)
		fallthrough
	default:
		foo(0)
		fallthrough
	case 3:
		foo(3)
	case 4:
		fallthrough
	case 5:
		foo(5)
		return
		a++
		fallthrough
	case 6:
		foo(6)
	}
}
`,
	},
	{
		skip: true,
		name: "switch cases everywhere",
		src: `
void foo(int p, char s) {
    switch (p) {
    case 0:
        if (s == 'a') {
    case 1:
			s = 'c';
        }
        break;
    }
}
`,
		exp: `
func foo(p int32, s int8) {
	switch p {
	case 1:
		goto case_1
	case 0:
		if s == 'a' {
		case_1:
			s = 'c'
		}
	}
}
`,
	},
	{
		skip: true,
		name: "switch cases everywhere 2",
		src: `
void foo(int p, char s) {
    switch (p) {
	if (s) {
    case 0:
        if (s == 'a') {
    case 1:
			s = 'c';
        }
        break;
	}
	case 2:
		s = 'd';
        break;
    }
}
`,
		exp: `
func foo(p int32, s int8) {
	switch {
	case s != 0 && p == 1:
		goto case_1
	case s != 0 && p == 0:
		if s == 'a' {
		case_1:
			s = 'c'
		}
	case p == 2
		s = 'd'
	}
}
`,
	},
	{
		name: "revert last if",
		skip: !optimizeStatements,
		src: `
int foo(int* a, int* b) {
	if (a) {
		foo(a);
		foo(a);
		foo(a);
	}
	foo(b);
	return 1;
}
`,
		exp: `
func foo(a *int32, b *int32) int32 {
	if a == nil {
		foo(b)
		return 1
	}
	foo(a)
	foo(a)
	foo(a)
	foo(b)
	return 1
}
`,
	},
	{
		name: "revert last if goto",
		skip: !optimizeStatements,
		src: `
int foo(int* a, int* b) {
	if (a) {
		foo(a);
		foo(a);
		foo(a);
	}
LABEL_X:
	foo(b);
	return 1;
}
`,
		exp: `
func foo(a *int32, b *int32) int32 {
	if a == nil {
		foo(b)
		return 1
	}
	foo(a)
	foo(a)
	foo(a)
	foo(b)
	return 1
}
`,
	},
	{
		name: "revert last if void",
		skip: !optimizeStatements,
		src: `
void foo(int* a) {
	if (a) {
		foo(a);
		foo(a);
		foo(a);
	}
}
`,
		exp: `
func foo(a *int32) {
	if a == nil {
		return
	}
	foo(a)
	foo(a)
	foo(a)
}
`,
	},
	{
		name: "move return to if",
		skip: !optimizeStatements,
		src: `
int foo(int* a) {
	if (a) {
		foo(a);
	} else {
		foo(a);
	}
	return 1;
}
`,
		exp: `
func foo(a *int32) int32 {
	if a != nil {
		foo(a)
		return 1
	}
	foo(a)
	return 1
}
`,
	},
	{
		name: "move return to if nested",
		skip: !optimizeStatements,
		src: `
int foo(int a) {
	if (a == 1) {
		foo(2);
		if (a == 3) {
			foo(4);
		} else if (a == 5) {
			foo(6);
		}
	}
	return 1;
}
`,
		exp: `
func foo(a int32) int32 {
	if a != 1 {
		return 1
	}
	foo(2)
	if a == 3 {
		foo(4)
		return 1
	} else if a == 5 {
		foo(6)
		return 1
	}
	return 1
}
`,
	},
	{
		name: "revert last if cost",
		skip: !optimizeStatements,
		src: `
int foo(int* a) {
	if (a) {
		foo(a);
		foo(a);
		if (a) {
			foo(a);
		}
		foo(a);
		return 1;
	}
	foo(a);
	foo(a);
	foo(a);
	return 0;
}
`,
		exp: `
func foo(a *int32) int32 {
	if a == nil {
		foo(a)
		foo(a)
		foo(a)
		return 0
	}
	foo(a)
	foo(a)
	if a != nil {
		foo(a)
	}
	foo(a)
	return 1
}
`,
	},
	{
		name: "sub_40BC10",
		inc:  `typedef unsigned int _DWORD;`,
		src: `
unsigned char blob[10];
char* foo(int a) {
	return (char*)(*(_DWORD*)& blob[3] + 160 * a);
}
`,
		exp: `
var blob [10]uint8

func foo(a int32) *byte {
	return (*byte)(unsafe.Pointer(uintptr(*(*_DWORD)(unsafe.Pointer(&blob[3])) + _DWORD(a*160))))
}
`,
	},
	{
		name: "inline gotos",
		skip: !optimizeStatements,
		src: `
int foo(int* a) {
	if (a) {
		foo(a);
LABEL_2:
		foo(a);
		if (a) {
			foo(a);
			goto LABEL_1;
		}
LABEL_1:
		foo(a);
		return 1;
	}
	foo(a);
	foo(a);
	goto LABEL_2;
}
`,
		exp: `
func foo(a *int32) int32 {
	if a == nil {
		foo(a)
		foo(a)
		goto LABEL_2
	}
	foo(a)
LABEL_2:
	foo(a)
	if a == nil {
		foo(a)
		return 1
	}
	foo(a)
	foo(a)
	return 1
}
`,
	},
	{
		name: "inline gotos chain",
		skip: !optimizeStatements,
		src: `
int foo(int* a) {
	if (a) {
		foo(a);
		foo(a);
		if (a) {
LABEL_2:
			foo(a);
			goto LABEL_1;
		}
LABEL_1:
		foo(a);
		return 1;
	}
	foo(a);
	foo(a);
	goto LABEL_2;
}
`,
		exp: `
func foo(a *int32) int32 {
	if a == nil {
		foo(a)
		foo(a)
		foo(a)
		foo(a)
		return 1
	}
	foo(a)
	foo(a)
	if a == nil {
		foo(a)
		return 1
	}
	foo(a)
	foo(a)
	return 1
}
`,
	},
	{
		name: "blocks with vars",
		src: `
#define set(s) \
{ int t = s;\
}

void main() {
  int s;
  if (0) {
    set(s)
    set(s)
    set(s)
  }
}
`,
		exp: `
func main() {
	var s int32
	if false {
		{
			var t int32 = s
			_ = t
		}
		{
			var t int32 = s
			_ = t
		}
		{
			var t int32 = s
			_ = t
		}
	}
}
`,
	},
	{
		name: "blocks no vars",
		src: `
#define set(s) \
{ t = s;\
}

void main() {
  int s, t;
  if (0) {
    set(s)
    set(s)
    set(s)
  }
}
`,
		exp: `
func main() {
	var (
		s int32
		t int32
	)
	_ = t
	if false {
		t = s
		t = s
		t = s
	}
}
`,
	},
	{
		name:     "sizeof only",
		builtins: true,
		src: `
int a = sizeof(int);
`,
		exp: `
var a int32 = int32(unsafe.Sizeof(int32(0)))
`,
	},
	{
		name:     "sizeof",
		builtins: true,
		src: `
void foo() {
	size_t a;
	int b[10];
	unsigned char b2[10];
	void* c;
	int* d;
	int (*e)(void);
	a = sizeof(b);
	a = sizeof(b2);
	a = sizeof(c);
	a = sizeof(d);
	a = sizeof(e);
}
`,
		exp: `
func foo() {
	var a size_t
	_ = a
	var b [10]int32
	_ = b
	var b2 [10]uint8
	_ = b2
	var c unsafe.Pointer
	_ = c
	var d *int32
	_ = d
	var e func() int32
	_ = e
	a = size_t(unsafe.Sizeof([10]int32{}))
	a = size_t(10)
	a = size_t(unsafe.Sizeof(unsafe.Pointer(nil)))
	a = size_t(unsafe.Sizeof((*int32)(nil)))
	a = size_t(unsafe.Sizeof(uintptr(0)))
}
`,
	},
	{
		name: "assign expr",
		src: `
void foo() {
	int a;
	int b;
	b = a = 1;
}
`,
		exp: `
func foo() {
	var (
		a int32
		b int32
	)
	_ = b
	b = func() int32 {
		a = 1
		return a
	}()
}
`,
	},
	{
		name: "stdint override",
		src: `
#include <stdint.h>

void foo() {
	int16_t a1;
	uint32_t a2;
}
`,
		exp: `
func foo() {
	var a1 int16
	_ = a1
	var a2 uint32
	_ = a2
}
`,
	},
	{
		name: "void cast",
		src: `
void foo(int a, int* b) {
	int c;
	(void)a;
	(void)b;
	(void)c;
}
`,
		exp: `
func foo(a int32, b *int32) {
	var c int32
	_ = a
	_ = b
	_ = c
}
`,
	},
	{
		name: "assign ternary",
		src: `
void foo(int a) {
	a = a ? 1 : 0;
	a = -(a ? 1 : 0);
	int b = a ? 1 : 0;
}
`,
		exp: `
func foo(a int32) {
	if a != 0 {
		a = 1
	} else {
		a = 0
	}
	if a != 0 {
		a = -1
	} else {
		a = 0
	}
	var b int32
	_ = b
	if a != 0 {
		b = 1
	} else {
		b = 0
	}
}
`,
	},
	{
		name: "return ternary",
		src: `
int foo(int a) {
	return a ? 1 : 0;
}
`,
		exp: `
func foo(a int32) int32 {
	if a != 0 {
		return 1
	}
	return 0
}
`,
	},
	{
		skip: true,
		name: "multiple assignments",
		src: `
void foo(int a) {
	int b, c;
	c = b = a;
	c = b = a++;
	c = b = ++a;
}
`,
		exp: `
func foo(a int32) {
	var b, c int32
	b, c = a, a
	b, c = a, a
	a++
	a++
	b, c = a, a
}
`,
	},
	{
		skip: true,
		name: "if init 1",
		src: `
void foo(int a) {
	if (a = 1, a) {
		a = 0;
	}
}
`,
		exp: `
func foo(a int32) {
	if a = 1; a != 0 {
		a = 0
	}
}
`,
	},
	{
		skip: true,
		name: "if init 2",
		src: `
void foo(int a) {
	if (a = 1) {
		a = 0;
	}
}
`,
		exp: `
func foo(a int32) {
	if a = 1; a != 0 {
		a = 0
	}
}
`,
	},
	{
		skip: true,
		name: "if init 3",
		src: `
void foo(int a) {
	if ((a = 1) != 0) {
		a = 0;
	}
}
`,
		exp: `
func foo(a int32) {
	if a = 1; a != 0 {
		a = 0
	}
}
`,
	},
	{
		name: "multiple stmt splits",
		src: `
typedef struct list_t list_t;
typedef struct list_t {
	list_t* next;
} list_t;

void foo() {
	list_t *elt, *list;
	for (elt=list; elt ? (list=elt->next, elt->next=0), 1 : 0; elt=list) {}
}
`,
		exp: `
type list_t struct {
	Next *list_t
}

func foo() {
	var (
		elt  *list_t
		list *list_t
	)
	for elt = list; func() int {
		if elt != nil {
			return func() int {
				list = elt.Next
				elt.Next = nil
				return 1
			}()
		}
		return 0
	}() != 0; elt = list {
	}
}
`,
		envFuncs: []envFunc{func(c *types.Config) {
			c.UseGoInt = true
		}},
	},
	{
		name: "do not edit",
		src: `

void foo() {}
`,
		exp: `
// Code generated by cxgo. DO NOT EDIT.

package lib

func foo() {
}
`,
		configFuncs: []configFunc{withDoNotEdit(true)},
	},
}

func TestTranslate(t *testing.T) {
	runTestTranslate(t, casesTranslate)
}

func runTestTranslateCase(t *testing.T, c parseCase) {
	if c.shouldSkip() {
		defer func() {
			if r := recover(); r != nil {
				defer debug.PrintStack()
				t.Skip(r)
			}
		}()
	}
	fname := strings.ReplaceAll(c.name, " ", "_") + ".c"
	var srcs []cc.Source
	if c.inc != "" {
		srcs = append(srcs, cc.Source{
			Name:  strings.TrimSuffix(fname, ".c") + "_predef.h",
			Value: c.inc,
		})
	}
	srcs = append(srcs, cc.Source{
		Name:  fname,
		Value: c.src,
	})
	econf := types.Config32()
	for _, f := range c.envFuncs {
		f(&econf)
	}
	env := libs.NewEnv(econf)
	ast, err := ParseSource(env, ParseConfig{
		WorkDir:    "",
		Predefines: c.builtins,
		Sources:    srcs,
	})
	if c.skip {
		t.SkipNow()
	} else {
		require.NoError(t, err)
	}

	tconf := Config{ForwardDecl: true}
	for _, f := range c.configFuncs {
		f(&tconf)
	}
	decls, err := TranslateAST(fname, ast, env, tconf)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = PrintGo(buf, testPkg, decls, tconf.DoNotEdit)
	assert.NoError(t, err)

	exp := c.exp
	a, b := strings.TrimSpace(exp), strings.TrimSpace(strings.TrimPrefix(buf.String(), "package lib"))
	if a != b {
		if c.skipExp != "" {
			a = strings.TrimSpace(c.skipExp)
			require.Equal(t, a, b)
			t.SkipNow()
		} else if c.skip {
			t.SkipNow()
		} else {
			require.Equal(t, a, b)
		}
	}
	if c.skip && !t.Failed() {
		require.Fail(t, "skipped test passes")
	}
}

func runTestTranslate(t *testing.T, cases []parseCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runTestTranslateCase(t, c)
		})
	}
}

func TestTypeResolution(t *testing.T) {
	const (
		fname = "resolve.c"
	)
	env := libs.NewEnv(types.Config32())
	ast, err := ParseSource(env, ParseConfig{
		WorkDir: "",
		Sources: []cc.Source{
			//		{
			//			Name: "resolve.h",
			//			Value:`
			//void unused(char*,int);
			//`,
			//		},
			{
				Name: fname,
				Value: `
typedef int int_t;
int_t a, b;
struct bar {
	int_t x;
	int_t y;
};
struct bar *f1(int_t c, struct bar e) {
	int_t b;
	int_t d = a + b + c + e.x + e.y;
}
`,
			},
		},
	})
	require.NoError(t, err)

	decls, err := TranslateCAST(fname, ast, env, Config{})
	require.NoError(t, err)
	require.Len(t, decls, 5)

	// find the first int typedef
	d1, ok := decls[0].(*CTypeDef)
	require.True(t, ok)
	// check the named type that we get
	intT := d1.Named
	assert.Equal(t, "int_t", intT.Name().Name)
	assert.Equal(t, types.IntT(4), intT.Underlying())

	// both variable declarations should have exactly the same type (by pointer)
	a1d, ok := decls[1].(*CVarDecl)
	require.True(t, ok)
	assert.Equal(t, "a", a1d.Names[0].Name)
	assert.True(t, intT == a1d.Type, "invalid type in global")

	b1d, ok := decls[2].(*CVarDecl)
	require.True(t, ok)
	assert.Equal(t, "b", b1d.Names[0].Name)
	assert.True(t, intT == b1d.Type, "invalid type in global 2")

	// test struct definition
	bard, ok := decls[3].(*CTypeDef)
	require.True(t, ok)
	assert.Equal(t, "bar", bard.Name().Name)
	bar, ok := bard.Underlying().(*types.StructType)
	require.True(t, ok)
	require.Len(t, bar.Fields(), 2)
	xf := bar.Fields()[0]
	yf := bar.Fields()[1]
	assert.True(t, intT == xf.Type(), "invalid type in struct field 1")
	assert.True(t, intT == yf.Type(), "invalid type in struct field 2")

	// check that the same type works in function args
	f1d, ok := decls[4].(*CFuncDecl)
	require.True(t, ok)
	assert.Equal(t, "f1", f1d.Name.Name)
	f1t := f1d.Type
	require.Len(t, f1t.Args(), 2)
	assert.True(t, intT == f1t.Args()[0].Type(), "invalid type in func arg")
	assert.True(t, bard.Named == f1t.Args()[1].Type(), "invalid type in func arg")
	p1, ok := f1t.Return().(types.PtrType)
	require.True(t, ok)
	assert.True(t, bard.Named == p1.Elem(), "invalid type in func return")

	require.Len(t, f1d.Body.Stmts, 3)
	// first decl is always __func__
	s1, ok := f1d.Body.Stmts[0].(*CDeclStmt)
	require.True(t, ok)
	fnc_, ok := s1.Decl.(*CVarDecl)
	require.True(t, ok)
	require.Equal(t, "__func__", fnc_.Names[0].Name)

	// check if the type works in function body
	s1, ok = f1d.Body.Stmts[1].(*CDeclStmt)
	require.True(t, ok)
	b2d, ok := s1.Decl.(*CVarDecl)
	require.True(t, ok)
	assert.Equal(t, "b", b2d.Names[0].Name)
	assert.True(t, intT == b2d.Type, "invalid type in local")
	// and that we can distinguish globals from locals
	assert.True(t, b1d != b2d)
	assert.True(t, b1d.Names[0] != b2d.Names[0])

	// check same conditions for a second local decl
	s1, ok = f1d.Body.Stmts[2].(*CDeclStmt)
	require.True(t, ok)
	d1d, ok := s1.Decl.(*CVarDecl)
	require.True(t, ok)
	assert.Equal(t, "d", d1d.Names[0].Name)
	assert.True(t, intT == d1d.Type, "invalid type in local 2")
	// check variables - one is global, second is local, third is an arg
	e1, ok := d1d.Inits[0].(*CBinaryExpr)
	require.True(t, ok, "%T", d1d.Inits[0])
	e2, ok := e1.Right.(*CSelectExpr)
	assert.True(t, ok, "%T", e1.Right)
	assert.True(t, e2.Sel == yf.Name, "invalid field ref 5")                  // arg
	assert.True(t, e2.Expr == IdentExpr{f1t.Args()[1].Name}, "invalid ref 5") // arg

	e1, ok = e1.Left.(*CBinaryExpr)
	require.True(t, ok, "%T", e1.Left)
	e2, ok = e1.Right.(*CSelectExpr)
	assert.True(t, ok, "%T", e1.Right)
	assert.True(t, e2.Sel == xf.Name, "invalid field ref 4")                  // arg
	assert.True(t, e2.Expr == IdentExpr{f1t.Args()[1].Name}, "invalid ref 4") // arg

	e1, ok = e1.Left.(*CBinaryExpr)
	require.True(t, ok, "%T", e1.Left)
	assert.True(t, IdentExpr{f1t.Args()[0].Name} == e1.Right, "invalid ref 3") // arg

	e1, ok = e1.Left.(*CBinaryExpr)
	require.True(t, ok)
	assert.True(t, IdentExpr{a1d.Names[0]} == e1.Left, "invalid ref 1")  // global
	assert.True(t, IdentExpr{b2d.Names[0]} == e1.Right, "invalid ref 2") // local
}
