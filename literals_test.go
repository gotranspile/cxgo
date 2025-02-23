package cxgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntLitString(t *testing.T) {
	cases := []struct {
		name string
		v    IntLit
		exp  string
	}{
		{"int(1)", cIntLit(1, 10), "1"},
		{"int(-1)", cIntLit(-1, 10), "-1"},
		{"-int(-1)", cIntLit(-1, 10).NegateLit(), "1"},
		{"uint(1)", cUintLit(1, 10), "1"},
		{"-uint(1)", cUintLit(1, 10).NegateLit(), "-1"},
		{"int(0x10)", cIntLit(0x10, 16), "0x10"},
		{"uint(0x10)", cUintLit(0x10, 16), "0x10"},
		{"-int(0x10)", cIntLit(0x10, 16).NegateLit(), "-16"},
		{"-uint(0x10)", cUintLit(0x10, 16).NegateLit(), "-16"},
		{"int8(0x80)", cUintLit(0x80, 16).OverflowInt(1), "-128"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.exp, c.v.String())
		})
	}
}

var casesTranslateLiterals = []parseCase{
	{
		name: "literal statement",
		src: `
void foo() {
	0;
}
`,
		exp: `
func foo() {
}
`,
	},
	{
		name: "string -> byte ptr",
		src: `
void foo() {
	char* a;
	a = "abc";
}
`,
		exp: `
func foo() {
	var a *byte
	_ = a
	a = libc.CString("abc")
}
`,
	},
	{
		name:     "string len",
		builtins: true,
		src: `
void foo() {
	int a;
	a = sizeof("abc");
}
`,
		exp: `
func foo() {
	var a int32
	_ = a
	a = int32(len("abc") + 1)
}
`,
	},
	{
		name: "rune -> int",
		src: `
void foo() {
	int a = 'a';
}
`,
		exp: `
func foo() {
	var a int32 = 'a'
	_ = a
}
`,
	},
	{
		name: "wstring -> wchar ptr",
		src: `
#include <stddef.h>

void foo() {
	wchar_t* a;
	a = L"abc";
}
`,
		exp: `
func foo() {
	var a *libc.WChar
	_ = a
	a = libc.CWString("abc")
}
`,
	},
	{
		name: "negative char",
		src: `
void foo() {
	char a;
	a = -1;
	a = -2;
}
`,
		exp: `
func foo() {
	var a int8
	_ = a
	a = -1
	a = -2
}
`,
	},
	{
		name: "negative uchar",
		src: `
void foo() {
	unsigned char a;
	a = -1;
	a = -2;
}
`,
		exp: `
func foo() {
	var a uint8
	_ = a
	a = 255
	a = 254
}
`,
	},
	{
		name: "negative ushort",
		src: `
void foo() {
	unsigned short a;
	a = -1;
	a = -2;
}
`,
		exp: `
func foo() {
	var a uint16
	_ = a
	a = 65535
	a = 65534
}
`,
	},
	{
		name: "negative uint",
		src: `
void foo() {
	unsigned int a;
	a = -1;
	a = -2;
}
`,
		exp: `
func foo() {
	var a uint32
	_ = a
	a = 4294967295
	a = 4294967294
}
`,
	},
	{
		name: "int overflow",
		src: `
void foo(int a) {
	if (a & 0xFFFF0000) {
		return;
	}
	a = 0x80000000;
	a = 2415929931;
}
`,
		exp: `
func foo(a int32) {
	if uint32(a)&0xFFFF0000 != 0 {
		return
	}
	a = -2147483648
	a = -1879037365
}
`,
	},
	{
		name: "float const -> int",
		src: `
void foo() {
	int a = 1.1 + 0.5;
}
`,
		exp: `
func foo() {
	var a int32 = int32(math.Floor(1.1 + 0.5))
	_ = a
}
`,
	},
	{
		name: "float const + int",
		src: `
void foo(int a) {
	int b = a + (3.0 / 2);
}
`,
		exp: `
func foo(a int32) {
	var b int32 = int32(float64(a) + 3.0/2)
	_ = b
}
`,
	},
	{
		// TODO: seems like cc doesn't set macro token on unary expr
		name: "stdint const override",
		src: `
#include <stdint.h>

void foo() {
	int16_t a1 = INT16_MAX;
	a1 = INT16_MIN;
}
`,
		exp: `
func foo() {
	var a1 int16 = math.MaxInt16
	_ = a1
	a1 = -32768
}
`,
	},
	{
		name: "var init sum",
		src: `
int foo = 1+2;
`,
		exp: `
var foo int32 = 1 + 2
`,
	},
	{
		name: "char var init",
		src: `
char foo = '"';
`,
		exp: `
var foo int8 = '"'
`,
	},
	{
		name: "comp lit zero init",
		src: `
typedef struct A {
	int x;
} A;
typedef struct B {
	int x;
	int y;
} B;
typedef struct C {
	B x;
	int y;
} C;
A v1 = {0};
B v2 = {0};
C v3 = {0};
`,
		exp: `
type A struct {
	X int32
}
type B struct {
	X int32
	Y int32
}
type C struct {
	X B
	Y int32
}

var v1 A = A{}
var v2 B = B{}
var v3 C = C{}
`,
	},
	{
		name: "comp lit zero compare and assign",
		src: `
typedef struct A {
	int x;
} A;

void foo(void) {
	A v1;
	if (v1 == 0) {
		v1 = 0;
	}
}
`,
		exp: `
type A struct {
	X int32
}

func foo() {
	var v1 A
	if v1 == (A{}) {
		v1 = A{}
	}
}
`,
	},
	{
		name: "array lit",
		src: `
int arr1[10] = { 12, 34, 56, 78, 90, 123, 456, 789, 8642, 9753 };
int arr2[10] = { 12, 34, 56, 78, 90, 123, 456, 789, 8642, 9753, };
`,
		exp: `
var arr1 [10]int32 = [10]int32{12, 34, 56, 78, 90, 123, 456, 789, 8642, 9753}
var arr2 [10]int32 = [10]int32{12, 34, 56, 78, 90, 123, 456, 789, 8642, 9753}
`,
	},
	{
		name: "array lit indexes",
		src: `
int arr1[] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
    3,
  [4] = 4,
    5,
    6,
};
int arr2[7] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
    3,
  [4] = 4,
    5,
};
`,
		exp: `
var arr1 [7]int32 = [7]int32{0, 1, 2, 3, 4, 5, 6}
var arr2 [7]int32 = [7]int32{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
`,
	},
	{
		name: "nested struct fields init",
		src: `
struct inner {
   int f;
   int g;
   int h;
};
struct outer {
   int A;
   struct inner B;
   int C;
};
struct outer x = {
   .C = 100,
   .B.g = 200,
   .A = 300,
   .B.f = 400,
};
	`,
		exp: `
type inner struct {
	F int32
	G int32
	H int32
}
type outer struct {
	A int32
	B inner
	C int32
}

var x outer = outer{C: 100, B: inner{G: 200, F: 400}, A: 300}
	`,
	},
	{
		name: "primitive type as struct",
		src: `
float* a = &(float){0.5};
float b = (float){0.5};
	`,
		exp: `
var a *float32 = func() *float32 {
	var tmp float32 = 0.5
	return &tmp
}()
var b float32 = func() float32 {
	var tmp float32 = 0.5
	return tmp
}()
	`,
	},
	{
		name: "string literal ternary",
		src: `
int a;
char* b = a ? "1" : "2";
`,
		exp: `
var a int32
var b *byte = libc.CString(func() string {
	if a != 0 {
		return "1"
	}
	return "2"
}())
`,
	},
	{
		name: "init byte string",
		src: `
char b[] = "===";
`,
		exp: `
var b [4]byte = func() [4]byte {
	var t [4]byte
	copy(t[:], []byte("==="))
	return t
}()
`,
	},
	{
		name: "init named uint8_t string",
		src: `
#include <stdint.h>
typedef uint8_t  MYubyte;
MYubyte vendor[] = "something here";
`,
		exp: `
type MYubyte uint8

var vendor [15]MYubyte = func() [15]MYubyte {
	var t [15]MYubyte
	copy(t[:], []MYubyte("something here"))
	return t
}()
`,
	},
	{
		name: "double negate",
		src: `
#define MONE (-1)
int x;
int y = x - MONE;
`,
		exp: `
const MONE = -1

var x int32
var y int32 = x - (-1)
`,
	},
	{
		name: "float div literal",
		src: `
void foo() {
	float x = 4/3.0;
}
`,
		exp: `
func foo() {
	var x float32 = 4 / 3.0
	_ = x
}
`,
	},
}

func TestLiterals(t *testing.T) {
	runTestTranslate(t, casesTranslateLiterals)
}
