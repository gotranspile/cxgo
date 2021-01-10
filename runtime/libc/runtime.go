package libc

import (
	"unsafe"
	_ "unsafe" // required by go:linkname
)

//go:linkname findnull runtime.findnull
func findnull(str *byte) int

//go:linkname findnullw runtime.findnullw
func findnullw(str *uint16) int

//go:linkname gobytes runtime.gobytes
func gobytes(p *byte, n int) []byte

//go:linkname gostring runtime.gostring
func gostring(p *byte) string

//go:linkname gostringnocopy runtime.gostringnocopy
func gostringnocopy(p *byte) string

//go:linkname gostringn runtime.gostringn
func gostringn(p *byte, l int) string

//go:linkname gostringw runtime.gostringw
func gostringw(strw *uint16) string

type rtype = unsafe.Pointer

type emptyInterface struct {
	typ  rtype
	word unsafe.Pointer
}

func typeof(i interface{}) rtype {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return eface.typ
}

func sizeof(t rtype) uintptr {
	if t == nil {
		return 0
	}
	return *(*uintptr)(t)
}

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(typ rtype) unsafe.Pointer

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(typ rtype, size int) unsafe.Pointer
