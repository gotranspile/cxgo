package csys

import (
	"syscall"

	"github.com/dennwc/cxgo/runtime/libc"
)

const FIONREAD uintptr = 21531

func Ioctl(fd uintptr, req uintptr, args ...interface{}) int32 {
	var err syscall.Errno
	switch req {
	case FIONREAD:
		if len(args) != 1 {
			panic("invalid number of args")
		}
		p, err := libc.AsPtr(args[0])
		if err != nil {
			panic(err)
		}
		_, _, err = syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(p))
	}
	if err != 0 {
		return -1
	}
	return 0
}
