package libc

import "sync/atomic"

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
