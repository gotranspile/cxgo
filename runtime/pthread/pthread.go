package pthread

import (
	"sync"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

const MUTEX_RECURSIVE = 1

type Attr struct {
	_ int
}

func Create(th **Thread, attr *Attr, start func(unsafe.Pointer) unsafe.Pointer, arg unsafe.Pointer) int32 {
	panic("TODO")
}

type Thread struct {
	_ int
}

func (th *Thread) Join(ret unsafe.Pointer) int32 {
	panic("TODO")
}

func (th *Thread) TimedJoinNP(ret unsafe.Pointer, abs *libc.TimeSpec) int32 {
	panic("TODO")
}

type MutexAttr struct {
	typ int32
}

func (m *MutexAttr) Init() int32 {
	m.typ = 0
	return 0
}

func (m *MutexAttr) SetType(typ int32) int32 {
	m.typ = typ
	return 0
}

func (m *MutexAttr) Destroy() int32 {
	m.typ = 0
	return 0
}

type Mutex struct {
	mu *sync.Mutex
}

func (m *Mutex) Init(attr *MutexAttr) int32 {
	if attr.typ == MUTEX_RECURSIVE {
		// FIXME
	}
	m.mu = new(sync.Mutex)
	return 0
}

func (m *Mutex) Destroy() int32 {
	m.mu = nil
	return 0
}

func (m *Mutex) Lock() int32 {
	m.mu.Lock()
	return 0
}

func (m *Mutex) TryLock() int32 {
	panic("TODO")
}

func (m *Mutex) TimedLock(t *libc.TimeSpec) int32 {
	panic("TODO")
}

func (m *Mutex) Unlock() int32 {
	m.mu.Unlock()
	return 0
}
