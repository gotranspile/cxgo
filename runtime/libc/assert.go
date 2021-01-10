package libc

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Assert(cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func AsPtr(v interface{}) (unsafe.Pointer, error) {
	switch v := v.(type) {
	case unsafe.Pointer:
		return v, nil
	case *byte:
		return unsafe.Pointer(v), nil
	case *WChar:
		return unsafe.Pointer(v), nil
	case uintptr:
		return unsafe.Pointer(v), nil
	case uint64:
		return unsafe.Pointer(uintptr(v)), nil
	case int64:
		return unsafe.Pointer(uintptr(v)), nil
	case uint32:
		return unsafe.Pointer(uintptr(v)), nil
	case int32:
		return unsafe.Pointer(uintptr(v)), nil
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.UnsafePointer, reflect.Slice:
			return unsafe.Pointer(rv.Pointer()), nil
		case reflect.Uint, reflect.Uintptr, reflect.Uint64, reflect.Uint32:
			return unsafe.Pointer(uintptr(rv.Uint())), nil
		case reflect.Int, reflect.Int64, reflect.Int32:
			return unsafe.Pointer(uintptr(rv.Uint())), nil
		}
		return nil, fmt.Errorf("cannot cast to pointer: %T", v)
	}
}
