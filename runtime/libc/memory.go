package libc

import (
	"unsafe"
)

var allocs syncMap[unsafe.Pointer, []byte]

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

func malloc(sz int) []byte {
	b := makePad(sz, 0)
	p := unsafe.Pointer(&b[0])
	allocs.Store(p, b)
	return b
}

// Malloc allocates a region of memory.
func Malloc(sz int) unsafe.Pointer {
	b := malloc(sz)
	return unsafe.Pointer(&b[0])
}

// Calloc allocates a region of memory for num elements of size sz.
func Calloc(num, sz int) unsafe.Pointer {
	if num == 0 {
		return nil
	}
	if sz <= 0 {
		panic("size should be > 0")
	}
	b := makePad(num*sz, sz)
	p := unsafe.Pointer(&b[0])
	allocs.Store(p, b)
	return p
}

func withSize(p unsafe.Pointer) ([]byte, bool) {
	return allocs.Load(p)
}

func Realloc(buf unsafe.Pointer, sz int) unsafe.Pointer {
	if buf == nil {
		return Malloc(sz)
	}
	p := malloc(sz)
	src, ok := withSize(buf)
	if !ok {
		panic("realloc of a pointer not managed by cxgo")
	}
	copy(p, src)
	Free(buf)
	return unsafe.Pointer(&p[0])
}

// Free marks the memory as freed. May be a nop in Go.
func Free(p unsafe.Pointer) {
	allocs.Delete(p)
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
//
// Deprecated: use unsafe.Add
func IndexUnsafePtr(p unsafe.Pointer, i int) unsafe.Pointer {
	return unsafe.Add(p, i)
}

// IndexBytePtr unsafely moves a byte pointer by i bytes. An offset may be negative.
//
// Deprecated: use unsafe.Add
func IndexBytePtr(p *byte, i int) *byte {
	return (*byte)(unsafe.Add(unsafe.Pointer(p), i))
}

// UnsafeBytesN makes a slice of a given size starting at ptr.
//
// Deprecated: use unsafe.Slice
func UnsafeBytesN(ptr unsafe.Pointer, sz int) []byte {
	return unsafe.Slice((*byte)(ptr), sz)
}

// BytesN makes a slice of a given size starting at ptr.
// It accepts a *byte instead of unsafe pointer as UnsafeBytesN does, which allows to avoid unsafe import.
//
// Deprecated: use unsafe.Slice
func BytesN(p *byte, sz int) []byte {
	return unsafe.Slice(p, sz)
}

// UnsafeUint16N makes a uint16 slice of a given size starting at ptr.
//
// Deprecated: use unsafe.Slice
func UnsafeUint16N(ptr unsafe.Pointer, sz int) []uint16 {
	return unsafe.Slice((*WChar)(ptr), sz)
}

// Uint16N makes a uint16 slice of a given size starting at ptr.
// It accepts a *uint16 instead of unsafe pointer as UnsafeUint16N does, which allows to avoid unsafe import.
//
// Deprecated: use unsafe.Slice
func Uint16N(p *uint16, sz int) []uint16 {
	return unsafe.Slice(p, sz)
}

// UnsafeUint32N makes a uint32 slice of a given size starting at ptr.
//
// Deprecated: use unsafe.Slice
func UnsafeUint32N(ptr unsafe.Pointer, sz int) []uint32 {
	return unsafe.Slice((*uint32)(ptr), sz)
}

// Uint32N makes a uint32 slice of a given size starting at ptr.
// It accepts a *uint32 instead of unsafe pointer as UnsafeUint32N does, which allows to avoid unsafe import.
//
// Deprecated: use unsafe.Slice
func Uint32N(p *uint32, sz int) []uint32 {
	return unsafe.Slice(p, sz)
}
