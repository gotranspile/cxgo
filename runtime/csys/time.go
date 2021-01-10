package csys

import (
	"unsafe"

	"github.com/dennwc/cxgo/runtime/libc"
)

func GetTimeOfDay(t *libc.TimeVal, p unsafe.Pointer) int32 {
	panic("TODO")
}
