package libc

import (
	"reflect"
	"unsafe"
)

// makePad creates a slice with a given size, but adds padding before and after the slice.
// It is required to make some unsafe C code work, e.g. indexing elements after the slice end.
func makePad(sz int, pad int) []byte {
	if sz <= 0 {
		panic("size should be > 0")
	}
	if pad == 0 {
		pad = int(unsafe.Sizeof(uintptr(0)))
	}
	p := make([]byte, sz+pad*2)
	p = p[pad:]
	p = p[:sz:sz]
	return p
}

// Malloc allocates a region of memory.
func Malloc(sz int) unsafe.Pointer {
	p := makePad(sz, 0)
	return unsafe.Pointer(&p[0])
}

// Calloc allocates a region of memory for num elements of size sz.
func Calloc(num, sz int) unsafe.Pointer {
	if num == 0 {
		return nil
	}
	if sz <= 0 {
		panic("size should be > 0")
	}
	p := makePad(num*sz, sz)
	return unsafe.Pointer(&p[0])
}

func Realloc(buf unsafe.Pointer, sz int) unsafe.Pointer {
	if buf == nil {
		return Malloc(sz)
	}
	p := Malloc(sz)
	MemCpy(p, buf, sz)
	return p
}

// Free marks the memory as freed. May be a nop in Go.
func Free(p unsafe.Pointer) {
	// nop
}

// ToPointer converts a uintptr to unsafe.Pointer.
func ToPointer(p uintptr) unsafe.Pointer {
	return unsafe.Pointer(p)
}

// ToUintptr converts a unsafe.Pointer to uintptr.
func ToUintptr(p unsafe.Pointer) uintptr {
	return uintptr(p)
}

// PointerDiff calculates (a - b).
func PointerDiff(a, b unsafe.Pointer) int {
	return int(uintptr(a) - uintptr(b))
}

// IndexUnsafePtr unsafely moves a pointer by i bytes. An offset may be negative.
func IndexUnsafePtr(p unsafe.Pointer, i int) unsafe.Pointer {
	if i == 0 {
		return p
	}
	if i > 0 {
		return unsafe.Pointer(uintptr(p) + uintptr(i))
	}
	return unsafe.Pointer(uintptr(p) - uintptr(-i))
}

// IndexBytePtr unsafely moves a byte pointer by i bytes. An offset may be negative.
func IndexBytePtr(p *byte, i int) *byte {
	if i == 0 {
		return p
	}
	return (*byte)(IndexUnsafePtr(unsafe.Pointer(p), i))
}

// UnsafeBytesN makes a slice of a given size starting at ptr.
func UnsafeBytesN(ptr unsafe.Pointer, sz int) []byte {
	if sz < 0 {
		panic("negative size")
	}
	if ptr == nil {
		if sz == 0 {
			return nil
		}
		panic("nil pointer")
	}
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(ptr)
	h.Len = sz
	h.Cap = sz
	return b
}

// BytesN makes a slice of a given size starting at ptr.
// It accepts a *byte instead of unsafe pointer as UnsafeBytesN does, which allows to avoid unsafe import.
func BytesN(p *byte, sz int) []byte {
	return UnsafeBytesN(unsafe.Pointer(p), sz)
}

// UnsafeUint16N makes a uint16 slice of a given size starting at ptr.
func UnsafeUint16N(ptr unsafe.Pointer, sz int) []uint16 {
	if sz < 0 {
		panic("negative size")
	}
	if ptr == nil {
		if sz == 0 {
			return nil
		}
		panic("nil pointer")
	}
	var b []uint16
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(ptr)
	h.Len = sz
	h.Cap = sz
	return b
}

// Uint16N makes a uint16 slice of a given size starting at ptr.
// It accepts a *uint16 instead of unsafe pointer as UnsafeUint16N does, which allows to avoid unsafe import.
func Uint16N(p *uint16, sz int) []uint16 {
	return UnsafeUint16N(unsafe.Pointer(p), sz)
}

// UnsafeUint32N makes a uint32 slice of a given size starting at ptr.
func UnsafeUint32N(ptr unsafe.Pointer, sz int) []uint32 {
	if sz < 0 {
		panic("negative size")
	}
	if ptr == nil {
		if sz == 0 {
			return nil
		}
		panic("nil pointer")
	}
	var b []uint32
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(ptr)
	h.Len = sz
	h.Cap = sz
	return b
}

// Uint32N makes a uint32 slice of a given size starting at ptr.
// It accepts a *uint32 instead of unsafe pointer as UnsafeUint32N does, which allows to avoid unsafe import.
func Uint32N(p *uint32, sz int) []uint32 {
	return UnsafeUint32N(unsafe.Pointer(p), sz)
}
