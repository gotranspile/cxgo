package libc

import (
	"fmt"
	"reflect"
)

type ArgList struct {
	cur  int
	typ  uint
	args []interface{}
}

func (va *ArgList) Start(typ uint, rest []interface{}) {
	va.cur = 0
	va.typ = typ
	va.args = rest
}

func (va *ArgList) Args() []interface{} {
	return va.args
}

func (va *ArgList) Arg() interface{} {
	if va.cur >= len(va.args) {
		return nil
	}
	cur := va.args[va.cur]
	va.cur++
	return cur
}

func (va *ArgList) ArgUintptr() uintptr {
	v := va.Arg()
	switch v := v.(type) {
	case nil:
		return 0
	case uintptr:
		return v
	case int:
		return uintptr(v)
	case uint:
		return uintptr(v)
	case int32:
		return uintptr(v)
	case uint32:
		return uintptr(v)
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.UnsafePointer:
		return rv.Pointer()
	}
	panic(fmt.Errorf("unsupported type: %T", v))
}

func (va *ArgList) End() {
	va.cur = 0
	va.typ = 0
	va.args = nil
}

func ArgCopy(dst, src *ArgList) {
	dst.cur = src.cur
	dst.typ = src.typ
	dst.args = append([]interface{}{}, src.args...)
}
