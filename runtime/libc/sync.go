package libc

import (
	"sync"
	"sync/atomic"
)

type syncMap[K comparable, V any] struct {
	m sync.Map
}

func (m *syncMap[K, V]) Store(k K, v V) {
	m.m.Store(k, v)
}

func (m *syncMap[K, V]) Load(k K) (V, bool) {
	vi, ok := m.m.Load(k)
	if !ok {
		var zero V
		return zero, false
	}
	return vi.(V), true
}

func (m *syncMap[K, V]) Delete(k K) {
	m.m.Delete(k)
}

// old = *p; *p (op)= val; return old;

func LoadAddInt32(p *int32, v int32) int32 {
	nv := atomic.AddInt32(p, v)
	return nv - v
}

func LoadSubInt32(p *int32, v int32) int32 {
	nv := atomic.AddInt32(p, -v)
	return nv + v
}

func LoadOrInt32(p *int32, v int32) int32 {
	for {
		old := atomic.LoadInt32(p)
		nv := old | v
		if atomic.CompareAndSwapInt32(p, old, nv) {
			return old
		}
	}
}

func LoadAndInt32(p *int32, v int32) int32 {
	for {
		old := atomic.LoadInt32(p)
		nv := old & v
		if atomic.CompareAndSwapInt32(p, old, nv) {
			return old
		}
	}
}

func LoadXorInt32(p *int32, v int32) int32 {
	for {
		old := atomic.LoadInt32(p)
		nv := old ^ v
		if atomic.CompareAndSwapInt32(p, old, nv) {
			return old
		}
	}
}

func LoadNandInt32(p *int32, v int32) int32 {
	for {
		old := atomic.LoadInt32(p)
		nv := ^old & v
		if atomic.CompareAndSwapInt32(p, old, nv) {
			return old
		}
	}
}
