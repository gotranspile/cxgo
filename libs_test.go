package cxgo

import "testing"

var casesTranslateLibs = []parseCase{
	{
		name: "panic",
		src: `
#include <stdlib.h>

void foo() {
	abort();
	__builtin_abort();
	__builtin_trap();
	__builtin_unreachable();
}
`,
		exp: `
func foo() {
	panic("abort")
	panic("abort")
	panic("trap")
	panic("unreachable")
}
`,
	},
	{
		name: "assert",
		src: `
#include <assert.h>

void foo(int a) {
	assert(a);
	assert(a != 5);
	assert(0);
	assert(!"fail");
}
`,
		exp: `
func foo(a int32) {
	if a == 0 {
		panic("assert failed")
	}
	if a == 5 {
		panic("assert failed")
	}
	panic(0)
	panic("fail")
}
`,
	},
	{
		name: "varargs",
		src: `
#include <stdarg.h>

void foo(int a, ...) {
	va_list va;
	va_start(va, a);
	int b = va_arg(va, int);
}
`,
		exp: `
func foo(a int32, _rest ...interface{}) {
	var va libc.ArgList
	va.Start(a, _rest)
	var b int32 = va.Arg().(int32)
	_ = b
}
`,
	},
	{
		name: "inet",
		src: `
#include <arpa/inet.h>

void foo() {
	in_addr_t a = inet_addr("1.2.3.4");
}
`,
		exp: `
func foo() {
	var a cnet.Addr = cnet.ParseAddr("1.2.3.4")
	_ = a
}
`,
	},
	{
		name: "sys socket",
		src: `
#include <sys/socket.h>

void foo() {
	in_addr_t a = inet_addr("1.2.3.4");
}
`,
		exp: `
func foo() {
	var a cnet.Addr = cnet.ParseAddr("1.2.3.4")
	_ = a
}
`,
	},
	{
		name: "math",
		src: `
#include <math.h>

void foo(float x, double y) {
	x = M_PI; y = M_PI;
	x = modff(x, &x); y = modf(y, &y);
	x = sinf(x);      y = sin(y);
	x = coshf(x);     y = cosh(y);
	x = atanf(x);     y = atan(y);
	x = roundf(x);    y = round(y);
	x = fabsf(x);     y = fabs(y);
	x = powf(x, x);   y = pow(y, y);
	y = fmod(y, y);
}
`,
		exp: `
func foo(x float32, y float64) {
	x = float32(math.Pi)
	y = math.Pi
	x = cmath.Modff(x, &x)
	y = cmath.Modf(y, &y)
	x = math32.Sin(x)
	y = math.Sin(y)
	x = math32.Cosh(x)
	y = math.Cosh(y)
	x = math32.Atan(x)
	y = math.Atan(y)
	x = float32(math.Round(float64(x)))
	y = math.Round(y)
	x = math32.Abs(x)
	y = math.Abs(y)
	x = math32.Pow(x, x)
	y = math.Pow(y, y)
	y = math.Mod(y, y)
}
`,
	},
	{
		name: "abs",
		src: `
#include <stdlib.h>
#include <math.h>

void foo(int abs) {}
`,
		exp: `
func foo(abs int32) {
}
`,
	},
	{
		name: "main no args",
		src: `
#include <stdlib.h>

void main() {
	exit(0);
}
`,
		exp: `
func main() {
	os.Exit(0)
}
`,
	},
	{
		name: "main with args",
		src: `
int main(int argc, char **argv) {
	if (argc == 1) {
		return argv != 0;
	}
	return 1;
}
`,
		exp: `
func main() {
	var (
		argc int32  = int32(len(os.Args))
		argv **byte = libc.CStringSlice(os.Args)
	)
	if argc == 1 {
		os.Exit(int(libc.BoolToInt(argv != nil)))
	}
	os.Exit(1)
}
`,
	},
	{
		name: "undef malloc",
		src: `
#include <stdlib.h>

#undef malloc
void foo() {
	void* p = malloc(10);
}
`,
		exp: `
func foo() {
	var p unsafe.Pointer = libc.Malloc(10)
	_ = p
}
`,
	},
}

var casesRunLibs = []struct {
	name string
	src  string
}{
	{
		name: "strtok",
		src: `
#include <string.h>
#include <stdio.h>

void printStr(const char* name, const char* s) {
	printf("%s, %d: '%s'\n", name, strlen(s), s);
}

int main() {
	const char* s = "_ a some __string here";
	char* buf = malloc(strlen(s)+1);
	strcpy(buf, s);
	printStr("buf 0", buf);

	char* tok = strtok(buf, " _");
	printStr("tok 1", tok);
	printStr("buf 1", buf);
	printf("diff 1: %d\n", tok-buf);
	for (int i = 0; i < 5; i++) {
		tok = strtok(0, " _");
		if (!tok) break;
		printStr("tok N", tok);
		printf("diff N: %d\n", tok-buf);
	}
	printStr("buf end", buf);
	free(buf);
	return 0;
}
`,
	},
	{
		name: "char",
		src: `
#include <stdio.h>

int main() {
	char a = 0;
	char b = a-1;
	printf("%d\n", (int)b);
	return 0;
}
`,
	},
	{
		name: "uchar",
		src: `
#include <stdio.h>

int main() {
	unsigned char a = 0;
	unsigned char b = a-1;
	printf("%d\n", (int)b);
	return 0;
}
`,
	},
	{
		name: "schar",
		src: `
#include <stdio.h>

int main() {
	signed char a = 0;
	signed char b = a-1;
	printf("%d\n", (int)b);
	return 0;
}
`,
	},
	{
		name: "char mult overflow",
		src: `
#include <stdio.h>

int main() {
	char a = 100;
	char b = 200;
	char x = a * b;
	int y = a * b;
	printf("%d, %d\n", (int)x, y);
}
`,
	},
	{
		name: "int mult overflow",
		src: `
#include <stdio.h>

int main() {
	int a = 0xf000000;
	int b = 0x1000000;
	int x = a * b;
	long long y = a * b;
	long long z = (long long)a * b;
	printf("%d, %lld, %lld\n", x, y, z);
}
`,
	},
	{
		name: "func name",
		src: `
#include <stdio.h>

void foo() {
	printf("%s, %s, %s\n", __func__, __FUNCTION__, __PRETTY_FUNCTION__);
}

int main() {
	foo();
}
`,
	},
}

func TestTranslateLibs(t *testing.T) {
	runTestTranslate(t, casesTranslateLibs)
}

func TestRunLibs(t *testing.T) {
	for _, c := range casesRunLibs {
		c := c
		t.Run(c.name, func(t *testing.T) {
			testTranspileOut(t, c.src)
		})
	}
}
