package libc

import (
	"sort"
	"unsafe"
)

type ptrSort struct {
	base unsafe.Pointer
	tmp  []byte
	num  int
	size int
	cmp  func(a, b unsafe.Pointer) int32
}

func (p *ptrSort) Len() int {
	return p.num
}

func (p *ptrSort) elem(i int) unsafe.Pointer {
	if i < 0 || i >= p.num {
		panic("index out of bounds")
	}
	return unsafe.Pointer(uintptr(p.base) + uintptr(i*p.size))
}

func (p *ptrSort) elems(i int) []byte {
	return unsafe.Slice((*byte)(p.elem(i)), p.size)
}

func (p *ptrSort) Less(i, j int) bool {
	a, b := p.elem(i), p.elem(j)
	return p.cmp(a, b) < 0
}

func (p *ptrSort) Swap(i, j int) {
	a, b := p.elems(i), p.elems(j)
	copy(p.tmp, a)
	copy(a, b)
	copy(b, p.tmp)
}

func Sort(base unsafe.Pointer, num, size uint32, compar func(a, b unsafe.Pointer) int32) {
	sort.Sort(&ptrSort{
		tmp:  make([]byte, size),
		base: base, num: int(num), size: int(size), cmp: compar,
	})
}

func Search(key, base unsafe.Pointer, num, size uint32, compar func(a, b unsafe.Pointer) int32) unsafe.Pointer {
	i := sort.Search(int(num), func(i int) bool {
		cur := unsafe.Pointer(uintptr(base) + uintptr(i)*uintptr(size))
		return compar(key, cur) < 0
	})
	i--
	if i < 0 {
		return nil
	}
	cur := unsafe.Pointer(uintptr(base) + uintptr(i)*uintptr(size))
	if compar(key, cur) == 0 {
		return cur
	}
	return nil
}
