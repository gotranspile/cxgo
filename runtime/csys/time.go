package csys

import (
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

func GetTimeOfDay(t *libc.TimeVal, p unsafe.Pointer) int32 {
	panic("TODO")
}
