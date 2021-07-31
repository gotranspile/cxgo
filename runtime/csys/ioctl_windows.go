// +build windows

package csys

func Ioctl(fd uintptr, req uintptr, args ...interface{}) int32 {
	panic("TODO")
}
