package cxgo

import "testing"

var casesTranslateFuncs = []parseCase{
	{
		name: "func ptr",
		src: `
int foo() {
	int (*a)();
	a = &foo;
	return a != 0;
}
`,
		exp: `
func foo() int32 {
	var a func() int32
	a = foo
	return libc.BoolToInt(a != nil)
}
`,
	},
	{
		name: "wrong func ptr",
		src: `
void foo() {
	int (*a)();
	a = &foo;
	return 1;
}
`,
		exp: `
func foo() {
	var a func() int32
	_ = a
	a = foo
	return 1
}
`,
	},
	{
		name: "call func from array",
		src: `
void foo(void) {}

static void (*functions[1])(void) = {
    foo,
};

void foo2() {
    functions[0]();
}
`,
		exp: `
func foo() {
}

var functions [1]*func() = [1]*func(){(*func())(unsafe.Pointer(libc.FuncAddr(foo)))}

func foo2() {
	libc.AsFunc(functions[0], (*func())(nil)).(func())()
}
`,
	},
}

func TestFunctions(t *testing.T) {
	runTestTranslate(t, casesTranslateFuncs)
}
