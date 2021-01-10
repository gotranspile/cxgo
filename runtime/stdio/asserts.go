package stdio

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/dennwc/cxgo/runtime/libc"
)

func asUint(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case uint:
		return uint64(v), nil
	case int:
		return uint64(v), nil
	case uintptr:
		return uint64(v), nil
	case uint64:
		return v, nil
	case int64:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case unsafe.Pointer:
		return uint64(uintptr(v)), nil
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.UnsafePointer, reflect.Slice:
			return uint64(rv.Pointer()), nil
		case reflect.Uint, reflect.Uintptr,
			reflect.Uint64, reflect.Uint32,
			reflect.Uint16, reflect.Uint8:
			return rv.Uint(), nil
		case reflect.Int,
			reflect.Int64, reflect.Int32,
			reflect.Int16, reflect.Int8:
			return uint64(rv.Int()), nil
		}
		return 0, fmt.Errorf("cannot cast to uint: %T", v)
	}
}

func asFloat(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Float64, reflect.Float32:
			return rv.Float(), nil
		}
		return 0, fmt.Errorf("cannot cast to float: %T", v)
	}
}

func asString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case *byte:
		return libc.GoString(v), nil
	case *libc.WChar:
		return libc.GoWString(v), nil
	case unsafe.Pointer:
		return libc.GoString((*byte)(v)), nil
	case uintptr:
		return libc.GoString((*byte)(unsafe.Pointer(v))), nil
	default:
		if v := reflect.ValueOf(v); v.Kind() == reflect.Array && v.Type().Elem().Kind() == reflect.Uint8 {
			b := make([]byte, v.Len())
			reflect.Copy(reflect.ValueOf(b), v)
			return libc.GoStringS(b), nil
		}
		p, err := libc.AsPtr(v)
		if err != nil {
			return "", err
		}
		return libc.GoString((*byte)(p)), nil
	}
}
