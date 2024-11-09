package cxgo

import (
	"fmt"
	"testing"
)

var casesTranslateNumbers = []parseCase{
	{
		name: "binary operators",
		src: `
void foo(int a, int b) {
	int c;
	c = a + b;
	c = a - b;
	c = a * b;
	c = a / b;
	c = a + 1;
	c = a - 1;
	c = a * 1;
	c = a / 1;
	c = 1 + b;
	c = 1 - b;
	c = 1 * b;
	c = 1 / b;
}
`,
		exp: `
func foo(a int32, b int32) {
	var c int32
	_ = c
	c = a + b
	c = a - b
	c = a * b
	c = a / b
	c = a + 1
	c = a - 1
	c = a * 1
	c = a / 1
	c = b + 1
	c = 1 - b
	c = b * 1
	c = 1 / b
}
`,
	},
	{
		name: "binary conversions",
		src: `
void foo(short a, int b) {
	int c;
	short d;
	c = a + b;
	c = a - b;
	c = a * b;
	c = a / b;
	c = b + a;
	c = b - a;
	c = b * a;
	c = b / a;
	d = a + b;
	d = a - b;
	d = a * b;
	d = a / b;
	d = b + a;
	d = b - a;
	d = b * a;
	d = b / a;
}
`,
		exp: `
func foo(a int16, b int32) {
	var c int32
	_ = c
	var d int16
	_ = d
	c = int32(a) + b
	c = int32(a) - b
	c = int32(a) * b
	c = int32(a) / b
	c = b + int32(a)
	c = b - int32(a)
	c = b * int32(a)
	c = b / int32(a)
	d = int16(int32(a) + b)
	d = int16(int32(a) - b)
	d = int16(int32(a) * b)
	d = int16(int32(a) / b)
	d = int16(b + int32(a))
	d = int16(b - int32(a))
	d = int16(b * int32(a))
	d = int16(b / int32(a))
}
`,
	},
	{
		name: "binary conversions sign",
		src: `
void foo(int a, unsigned int b) {
	int c;
	unsigned int d;
	c = a + b;
	c = a - b;
	c = a * b;
	c = a / b;
	c = b + a;
	c = b - a;
	c = b * a;
	c = b / a;
	d = a + b;
	d = a - b;
	d = a * b;
	d = a / b;
	d = b + a;
	d = b - a;
	d = b * a;
	d = b / a;
}
`,
		exp: `
func foo(a int32, b uint32) {
	var c int32
	_ = c
	var d uint32
	_ = d
	c = int32(uint32(a) + b)
	c = int32(uint32(a) - b)
	c = int32(uint32(a) * b)
	c = int32(uint32(a) / b)
	c = int32(b + uint32(a))
	c = int32(b - uint32(a))
	c = int32(b * uint32(a))
	c = int32(b / uint32(a))
	d = uint32(a) + b
	d = uint32(a) - b
	d = uint32(a) * b
	d = uint32(a) / b
	d = b + uint32(a)
	d = b - uint32(a)
	d = b * uint32(a)
	d = b / uint32(a)
}
`,
	},
	{
		name: "int -> named int",
		inc:  `typedef int NINT;`,
		src: `
void foo() {
	int a;
	NINT b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a int32
		b NINT
	)
	_ = b
	b = NINT(a)
}
`,
	},
	{
		name:     "int -> go int",
		builtins: true,
		src: `
void foo() {
	int a;
	_cxgo_go_int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a int32
		b int
	)
	_ = b
	b = int(a)
}
`,
	},
	{
		name: "int16 - uint16",
		src: `
void foo() {
	short a;
	unsigned short b;
	b = a;
	a = b;
	if (a < b) {
		return;
	}
}
`,
		exp: `
func foo() {
	var (
		a int16
		b uint16
	)
	b = uint16(a)
	a = int16(b)
	if int32(a) < int32(b) {
		return
	}
}
`,
	},
	{
		skip: true, // TODO: verify
		name: "int16 and uint16",
		src: `
void foo() {
	short a;
	unsigned short b;
	b = 8 - a;
	if (b < 8 - a) {
		return;
	}
}
`,
		exp: `
func foo() {
	var (
		a int16
		b uint16
	)
	b = uint16(8 - a)
	if int16(b) < 8 - a {
		return
	}
}
`,
	},
	{
		skip: true, // TODO: verify
		name: "types",
		src: `
void foo() {
	__int8  i8;
	__int16 i16;
	__int32 i32;
	__int64 i64;
	unsigned __int8  u8;
	unsigned __int16 u16;
	unsigned __int32 u32;
	unsigned __int64 u64;

	i16 = i8 * i8;
	i32 = i8 * i8;
	i32 = i16 * i16;
	i64 = i8 * i8;
	i64 = i16 * i16;
	i64 = i32 * i32;

	i16 = i8 * u8;
	i32 = i8 * u8;
	i32 = i16 * u16;
	i64 = i8 * u8;
	i64 = i16 * u16;
	i64 = i32 * u32;

	u16 = u8 * u8;
	u32 = u8 * u8;
	u32 = u16 * u16;
	u64 = u8 * u8;
	u64 = u16 * u16;
	u64 = u32 * u32;

	u16 = i8 * u8;
	u32 = i8 * u8;
	u32 = i16 * u16;
	u64 = i8 * u8;
	u64 = i16 * u16;
	u64 = i32 * u32;
}
`,
		exp: `
func foo() {
	var (
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
	)
	i16 = int16(int32(i8) * int32(i8))
	i32 = int32(i8) * int32(i8)
	i32 = int32(i16) * int32(i16)
	i64 = int64(int32(i8) * int32(i8))
	i64 = int64(int32(i16) * int32(i16))
	i64 = int64(i32) * int64(i32)
	i16 = int16(i8) * int16(u8)
	i32 = int32(i8) * int32(u8)
	i32 = int32(i16) * int32(u16)
	i64 = int64(i8) * int64(u8)
	i64 = int64(i16) * int64(u16)
	i64 = int64(i32) * int64(u32)
	u16 = uint16(u8) * uint16(u8)
	u32 = uint32(u8) * uint32(u8)
	u32 = uint32(u16) * uint32(u16)
	u64 = uint64(u8) * uint64(u8)
	u64 = uint64(u16) * uint64(u16)
	u64 = uint64(u32) * uint64(u32)
	u16 = uint16(i8) * uint16(u8)
	u32 = uint32(i8) * uint32(u8)
	u32 = uint32(i16) * uint32(u16)
	u64 = uint64(i8) * uint64(u8)
	u64 = uint64(i16) * uint64(u16)
	u64 = uint64(i32) * uint64(u32)
}
`,
	},
	{
		name: "named int -> int",
		inc:  `typedef int NINT;`,
		src: `
void foo() {
	NINT a;
	int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a NINT
		b int32
	)
	_ = b
	b = int32(a)
}
`,
	},
	{
		name: "named int eq int",
		inc:  `typedef int NINT;`,
		src: `
void foo() {
	NINT a;
	int b;
	if (b == a) return;
}
`,
		exp: `
func foo() {
	var (
		a NINT
		b int32
	)
	if b == int32(a) {
		return
	}
}
`,
	},
	{
		name: "diff ints",
		src: `
void foo() {
	unsigned short a;
	int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a uint16
		b int32
	)
	_ = b
	b = int32(a)
}
`,
	},
	{
		name: "sub unsigned const",
		src: `
void foo(unsigned int a) {
	unsigned int b;
	b = a - 0x1;
}
`,
		exp: `
func foo(a uint32) {
	var b uint32
	_ = b
	b = a - 0x1
}
`,
	},
	{
		name: "named int arithm",
		inc:  `typedef int my_size_t;`,
		src: `
void foo() {
	my_size_t a;
	my_size_t b;
	a = 1 + b;
	a = b + 1;
}
`,
		exp: `
func foo() {
	var a my_size_t
	_ = a
	var b my_size_t
	a = b + 1
	a = b + 1
}
`,
	},
	{
		name: "overflow mult char",
		src: `
void foo() {
	char a;
	int b;
	b = 300 * a;
}
`,
		exp: `
func foo() {
	var (
		a int8
		b int32
	)
	_ = b
	b = int32(a) * 300
}
`,
	},
	{
		name: "overflow mult schar",
		src: `
void foo() {
	signed char a;
	int b;
	b = 300 * a;
}
`,
		exp: `
func foo() {
	var (
		a int8
		b int32
	)
	_ = b
	b = int32(a) * 300
}
`,
	},
	{
		name: "overflow mult uchar",
		src: `
void foo() {
	unsigned char a;
	int b;
	b = 300 * a;
}
`,
		exp: `
func foo() {
	var (
		a uint8
		b int32
	)
	_ = b
	b = int32(a) * 300
}
`,
	},
	{
		name: "diff int equality",
		src: `
void foo() {
	short a;
	unsigned int b;
	if (b != (a & 0x17F0)) {
		return;
	}
}
`,
		exp: `
func foo() {
	var (
		a int16
		b uint32
	)
	if b != uint32(int32(a)&0x17F0) {
		return
	}
}
`,
	},
	{
		name: "const overflow",
		src: `
void foo(int a) {
	if (a != 0xdeadface) {
		return;
	}
}
`,
		exp: `
func foo(a int32) {
	if uint32(a) != 0xDEADFACE {
		return
	}
}
`,
	},
	{
		name: "const overflow 2",
		src: `
void foo(unsigned int a) {
	a = -1;
	a += -1;
	a = ~0;
	a ^= 0;
	a ^= ~0;
}
`,
		exp: `
func foo(a uint32) {
	a = 4294967295
	a += 4294967295
	a = uint32(^int32(0))
	a ^= 0
	a ^= uint32(^int32(0))
}
`,
	},
	{
		name: "const and named int",
		inc:  `typedef unsigned int mDWORD;`,
		src: `
void foo(int a, mDWORD b) {
	a = 1 - b;
}
`,
		exp: `
func foo(a int32, b mDWORD) {
	a = int32(1 - b)
}
`,
	},
	{
		name: "ternary",
		src: `
void foo(int a) {
	foo(a != 0 ? 0 : 99999);
}
`,
		exp: `
func foo(a int32) {
	foo(int32(func() uint32 {
		if a != 0 {
			return 0
		}
		return 99999
	}()))
}
`,
	},
	{
		name: "ternary overflow",
		src: `
void foo(int a) {
	foo(a != 0 ? 128 : -128);
}
`,
		exp: `
func foo(a int32) {
	foo(func() int32 {
		if a != 0 {
			return 128
		}
		return -128
	}())
}
`,
	},
	{
		name: "seq",
		src: `
void foo(int a) {
	a = (a = 1, 1);
}
`,
		exp: `
func foo(a int32) {
	a = func() int32 {
		a = 1
		return 1
	}()
}
`,
	},
	{
		name: "float and int",
		src: `
void foo(int a, int b, float c) {
	c = (a - b * c) / 2;
}
`,
		exp: `
func foo(a int32, b int32, c float32) {
	c = (float32(a) - float32(b)*c) / 2
}
`,
	},
	{
		name: "implicit and explicit bool",
		src: `
void foo(float a, float b, int c, float d) {
	d = (a < b) + c;
	d = (float)(a < b) + c;
}
`,
		exp: `
func foo(a float32, b float32, c int32, d float32) {
	d = float32(libc.BoolToInt(a < b) + c)
	d = float32(libc.BoolToInt(a < b)) + float32(c)
}
`,
	},
	{
		name: "sizeof and shift",
		src: `
#include <stdlib.h>

typedef unsigned int word;

void foo() {
	word a = (((word)1)<<((8*((int)sizeof(word)))-1));
}
`,
		exp: `
type word uint32

func foo() {
	var a word = word(1 << ((8 * (int32(unsafe.Sizeof(word(0))))) - 1))
	_ = a
}
`,
	},
}

func TestNumbers(t *testing.T) {
	t.Run("int types", testNumbersInts)
	t.Run("fixed types", testNumbersFixed)
	t.Run("cxgo types", testNumberscxgo)
	t.Run("stdlib types", testNumbersStdlib)
	t.Run("go types", testNumbersGo)
	runTestTranslate(t, casesTranslateNumbers)
}

func testNumbersInts(t *testing.T) {
	const (
		kindVal = iota
		kindPtr
		kindArr
	)
	const (
		signUndefined = iota
		signSigned
		signUnsigned
	)
	type intType struct {
		cname string
		size  int
	}
	for _, kind := range []int{
		kindVal, kindPtr, kindArr,
	} {
		for _, sign := range []int{
			signUndefined, signSigned, signUnsigned,
		} {
			for _, size := range []intType{
				{"char", 1},
				{"short", 2},
				{"int", 4},
				{"long", 4},
				{"long long", 8},
			} {
				ctype := size.cname
				switch sign {
				case signSigned:
					ctype = "signed " + ctype
				case signUnsigned:
					ctype = "unsigned " + ctype
				}
				csuff := ""
				switch kind {
				case kindPtr:
					ctype += "*"
				case kindArr:
					csuff = "[2]"
				}
				t.Run(ctype+csuff, func(t *testing.T) {
					gotype := fmt.Sprintf("int%d", size.size*8)
					switch sign {
					case signUnsigned:
						gotype = "u" + gotype
					}
					// special case for char pointers/arrays -> turn to byte
					if size.cname == "char" && sign == signUndefined && kind != kindVal {
						gotype = "byte"
					}
					switch kind {
					case kindPtr:
						gotype = "*" + gotype
					case kindArr:
						gotype = "[2]" + gotype
					}
					runTestTranslateCase(t, parseCase{
						src: ctype + " v" + csuff + ";",
						exp: "var v " + gotype,
					})
				})
			}
		}
	}
}

func testNumbersFixed(t *testing.T) {
	const (
		kindVal = iota
		kindPtr
		kindArr
	)
	const (
		signUndefined = iota
		signSigned
		signUnsigned
	)
	for _, kind := range []int{
		kindVal, kindPtr, kindArr,
	} {
		for _, sign := range []int{
			signUndefined, signSigned, signUnsigned,
		} {
			for _, size := range []int{
				1, 2, 4, 8,
			} {
				ctype := fmt.Sprintf("__int%d", size*8)
				switch sign {
				case signSigned:
					ctype = "signed " + ctype
				case signUnsigned:
					ctype = "unsigned " + ctype
				}
				csuff := ""
				switch kind {
				case kindPtr:
					ctype += "*"
				case kindArr:
					csuff = "[2]"
				}
				t.Run(ctype+csuff, func(t *testing.T) {
					gotype := fmt.Sprintf("int%d", size*8)
					switch sign {
					case signUnsigned:
						gotype = "u" + gotype
					}
					switch kind {
					case kindPtr:
						gotype = "*" + gotype
					case kindArr:
						gotype = "[2]" + gotype
					}
					runTestTranslateCase(t, parseCase{
						src: ctype + " v" + csuff + ";",
						exp: "var v " + gotype,
					})
				})
			}
		}
	}
}

func testNumberscxgo(t *testing.T) {
	const (
		kindVal = iota
		kindPtr
		kindArr
	)
	const (
		signUndefined = iota
		signSigned
		signUnsigned
	)
	for _, kind := range []int{
		kindVal, kindPtr, kindArr,
	} {
		for _, sign := range []int{
			signUndefined, signSigned, signUnsigned,
		} {
			for _, size := range []int{
				1, 2, 4, 8,
			} {
				csign := ""
				switch sign {
				case signSigned:
					csign = "s"
				case signUnsigned:
					csign = "u"
				}
				ctype := fmt.Sprintf("_cxgo_%sint%d", csign, size*8)
				csuff := ""
				switch kind {
				case kindPtr:
					ctype += "*"
				case kindArr:
					csuff = "[2]"
				}
				t.Run(ctype+csuff, func(t *testing.T) {
					gotype := fmt.Sprintf("int%d", size*8)
					switch sign {
					case signUnsigned:
						gotype = "u" + gotype
					}
					switch kind {
					case kindPtr:
						gotype = "*" + gotype
					case kindArr:
						gotype = "[2]" + gotype
					}
					runTestTranslateCase(t, parseCase{
						builtins: true,
						src:      ctype + " v" + csuff + ";",
						exp:      "var v " + gotype,
					})
				})
			}
		}
	}
}

func testNumbersStdlib(t *testing.T) {
	const (
		kindVal = iota
		kindPtr
		kindArr
	)
	const (
		signSigned = iota
		signUnsigned
	)
	for _, kind := range []int{
		kindVal, kindPtr, kindArr,
	} {
		for _, sign := range []int{
			signSigned, signUnsigned,
		} {
			for _, size := range []int{
				1, 2, 4, 8,
			} {
				ctype := fmt.Sprintf("int%d_t", size*8)
				if sign == signUnsigned {
					ctype = "u" + ctype
				}
				csuff := ""
				switch kind {
				case kindPtr:
					ctype += "*"
				case kindArr:
					csuff = "[2]"
				}
				t.Run(ctype+csuff, func(t *testing.T) {
					gotype := fmt.Sprintf("int%d", size*8)
					switch sign {
					case signUnsigned:
						gotype = "u" + gotype
					}
					switch kind {
					case kindPtr:
						gotype = "*" + gotype
					case kindArr:
						gotype = "[2]" + gotype
					}
					runTestTranslateCase(t, parseCase{
						src: "#include <stdint.h>\n" + ctype + " v" + csuff + ";",
						exp: "var v " + gotype,
					})
				})
			}
		}
	}
}

func testNumbersGo(t *testing.T) {
	for _, gotype := range []string{
		"int", "uint", "uintptr", "byte", "rune",
	} {
		t.Run(gotype, func(t *testing.T) {
			runTestTranslateCase(t, parseCase{
				builtins: true,
				src:      "_cxgo_go_" + gotype + " v;",
				exp:      "var v " + gotype,
			})
		})
	}
}
