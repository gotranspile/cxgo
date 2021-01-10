package libc

import (
	"reflect"
	"unsafe"
)

func BoolToInt(v bool) int32 {
	if v {
		return 1
	}
	return 0
}

// interfaceAddr extracts an address of an object stored in the interface value.
func interfaceAddr(v interface{}) uintptr {
	ifc := *(*[2]uintptr)(unsafe.Pointer(&v))
	return ifc[1]
}

// interfaceType extracts a type of an object stored in the interface value.
func interfaceType(v interface{}) uintptr {
	ifc := *(*[2]uintptr)(unsafe.Pointer(&v))
	return ifc[0]
}

// makeInterface creates a new interface value from an object address and a type.
func makeInterface(typ, ptr uintptr) interface{} {
	v := [2]uintptr{typ, ptr}
	return *(*interface{})(unsafe.Pointer(&v))
}

// FuncAddr converts a function value to a uintptr.
func FuncAddr(v interface{}) uintptr {
	if reflect.TypeOf(v).Kind() != reflect.Func {
		panic(v)
	}
	return interfaceAddr(v)
}

// FuncAddrUnsafe converts a function value to a unsafe.Pointer.
func FuncAddrUnsafe(v interface{}) unsafe.Pointer {
	return unsafe.Pointer(FuncAddr(v))
}

// AddrAsFunc converts a function address to a function value.
// The caller must type-assert to an expected function type.
// The function type can be specified as (*func())(nil).
func AddrAsFunc(addr uintptr, typ interface{}) interface{} {
	tp := reflect.TypeOf(typ).Elem()
	if tp.Kind() != reflect.Func {
		panic(typ)
	}
	t := interfaceType(reflect.Zero(tp).Interface())
	return makeInterface(t, addr)
}

var (
	reflUnsafePtr = reflect.TypeOf(unsafe.Pointer(nil))
	reflUintptr   = reflect.TypeOf(uintptr(0))
)

// AsFunc is a less restrictive version of AddrAsFunc.
// It accepts a pointer value that may be one of: uintptr, unsafe.Pointer, *T, int, uint.
func AsFunc(p interface{}, typ interface{}) interface{} {
	switch p := p.(type) {
	case uintptr:
		return AddrAsFunc(p, typ)
	case unsafe.Pointer:
		return AddrAsFunc(uintptr(p), typ)
	case uint:
		return AddrAsFunc(uintptr(p), typ)
	case uint64:
		return AddrAsFunc(uintptr(p), typ)
	case uint32:
		return AddrAsFunc(uintptr(p), typ)
	case int:
		return AddrAsFunc(uintptr(p), typ)
	case int64:
		return AddrAsFunc(uintptr(p), typ)
	case int32:
		return AddrAsFunc(uintptr(p), typ)
	}
	if t := reflect.TypeOf(p); t.Kind() != reflect.Ptr {
		panic("unsupported type: " + t.String())
	}
	pv := reflect.ValueOf(p).Pointer()
	return AddrAsFunc(pv, typ)
}
