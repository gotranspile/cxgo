package cxgo

import "testing"

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
	a = math.MaxUint8
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
	a = math.MaxUint16
	a = 0xFFFE
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
	a = math.MaxUint32
	a = 0xFFFFFFFE
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
	a = math.MinInt32
	a = -1879037365
}
`,
	},
	{
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
	a1 = math.MinInt16
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
		name: "array lit indexes",
		src: `
int arr[] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
    3,
  [4] = 4,
    5,
    6,
};
`,
		exp: `
var arr [7]int32 = [7]int32{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
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
	copy(t[:], "===")
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
var y int32 = x - int32(-1)
`,
		// TODO: x - (-1)
	}, {
		name: "float div literal",
		src: `
void foo() {
	float x = 4/3.0;
}
`,
		exp: `
func foo() {
	var x float32 = float32(4 / 3.0)
	_ = x
}
`,
	},
}

func TestLiterals(t *testing.T) {
	runTestTranslate(t, casesTranslateLiterals)
}
