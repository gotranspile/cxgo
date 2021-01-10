package libc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestFunc2Ptr(t *testing.T) {
	var (
		v1 bool
		v2 bool
	)
	f1 := func() {
		v1 = true
	}
	f2 := func() int {
		v2 = true
		return 1
	}
	p1 := FuncAddr(f1)
	p2 := FuncAddr(f2)

	f1 = AddrAsFunc(p1, (*func())(nil)).(func())
	f1()
	require.True(t, v1)

	f2 = AddrAsFunc(p2, (*func() int)(nil)).(func() int)
	r := f2()
	require.Equal(t, 1, r)
	require.True(t, v2)
}

func TestConvFuncStruct(t *testing.T) {
	type T1 struct {
		fnc func()
	}
	type T2 struct {
		fnc uint
	}
	var v bool
	var b [unsafe.Sizeof(uintptr(0))]byte

	p1 := (*T2)(unsafe.Pointer(&b[0]))
	p1.fnc = uint(FuncAddr(func() {
		v = true
	}))

	p2 := (*T1)(unsafe.Pointer(&b[0]))
	p2.fnc()
	require.True(t, v)
}

func TestConvFuncStruct2(t *testing.T) {
	type T1 struct {
		fnc func()
	}
	type T2 struct {
		fnc uint
	}
	var v bool
	var b [unsafe.Sizeof(uintptr(0))]byte

	p1 := (*T1)(unsafe.Pointer(&b[0]))
	p1.fnc = func() {
		v = true
	}

	p2 := (*T2)(unsafe.Pointer(&b[0]))
	AddrAsFunc(uintptr(p2.fnc), (*func())(nil)).(func())()
	require.True(t, v)
}
