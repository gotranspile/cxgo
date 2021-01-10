package cxgo

import "testing"

var casesTranslateBools = []parseCase{
	{
		name: "stdbool override",
		src: `
#include <stdbool.h>

void foo() {
	bool a = true;
	if (a) {}
}
`,
		exp: `
func foo() {
	var a bool = true
	if a {
	}
}
`,
	},
	{
		name: "if int",
		src: `
void foo() {
	int a;
	if (a) return;
}
`,
		exp: `
func foo() {
	var a int32
	if a != 0 {
		return
	}
}
`,
	},
	{
		name: "if not int",
		src: `
void foo() {
	int a;
	if (!a) return;
}
`,
		exp: `
func foo() {
	var a int32
	if a == 0 {
		return
	}
}
`,
	},
	{
		name: "if int eq",
		src: `
void foo() {
	int a;
	if (a == 1) return;
}
`,
		exp: `
func foo() {
	var a int32
	if a == 1 {
		return
	}
}
`,
	},
	{
		name: "if bool eq int 0",
		src: `
#include <stdbool.h>
void foo() {
	bool a;
	if (a == 0) return;
}
`,
		exp: `
func foo() {
	var a bool
	if !a {
		return
	}
}
`,
	},
	{
		name: "if bool neq int 0",
		src: `
#include <stdbool.h>
void foo() {
	bool a;
	if (a != 0) return;
}
`,
		exp: `
func foo() {
	var a bool
	if a {
		return
	}
}
`,
	},
	{
		name: "typedef bool",
		src: `
typedef int bool;

void foo() {
	bool a = 1;
}
`,
		exp: `
type bool int32

func foo() {
	var a bool = 1
	_ = a
}
`,
	},
	{
		name: "bool include",
		src: `
#include <stdbool.h>

void foo() {
	bool a = 1;
}
`,
		exp: `
func foo() {
	var a bool = true
	_ = a
}
`,
	},
	{
		name: "bool -> int",
		inc:  `#include <stdbool.h>`,
		src: `
void foo() {
	bool a;
	int b;
	b = a;
}
`,
		exp: `
func foo() {
	var (
		a bool
		b int32
	)
	_ = b
	b = libc.BoolToInt(a)
}
`,
	},
	{
		name: "add int to bool",
		src: `
int foo(int a) {
	int b;
	b = (a <= 0) + 0x7FFFFFFF;
	return b;
}
`,
		exp: `
func foo(a int32) int32 {
	var b int32
	b = libc.BoolToInt(a <= 0) + math.MaxInt32
	return b
}
`,
	},
	{
		name: "bool arithm",
		src: `
#include <stdbool.h>

void foo() {
	bool a;
	int b;
	b = -a - (a - 1);
}
`,
		exp: `
func foo() {
	var (
		a bool
		b int32
	)
	_ = b
	b = -libc.BoolToInt(a) - (libc.BoolToInt(a) - 1)
}
`,
	},
}

func TestBools(t *testing.T) {
	runTestTranslate(t, casesTranslateBools)
}
