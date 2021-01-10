package dlopen

import (
	"fmt"
	"unsafe"

	"github.com/dennwc/cxgo/runtime/libc"
)

// TODO: set correct values
const (
	RTLD_LAZY = 1 << iota
	RTLD_NOW
	RTLD_GLOBAL
	RTLD_LOCAL
)

type Library struct {
	name string
}

var gerr error

func Open(name string, flags int) *Library {
	// FIXME: dlopen
	panic(fmt.Errorf("dlopen(%q, 0x%x)", name, flags))
}

func (l *Library) Sym(name string) unsafe.Pointer {
	// FIXME: dlsym
	panic(fmt.Errorf("dlsym(%q)", name))
}

func (l *Library) Close() int {
	return 0
}

func Error() *byte {
	if gerr == nil {
		return nil
	}
	return libc.CString(gerr.Error())
}
