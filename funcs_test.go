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
}

func TestFunctions(t *testing.T) {
	runTestTranslate(t, casesTranslateFuncs)
}
