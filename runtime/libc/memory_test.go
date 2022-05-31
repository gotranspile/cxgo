package libc

import "testing"

func TestRealloc(t *testing.T) {
	p := Malloc(1)
	p = Realloc(p, 32*1024*1024)
	Free(p)
}
