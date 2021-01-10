package libc

import (
	"errors"
	"os"
	"syscall"
)

var (
	Errno int
	errno int // used to check if C code have changed Errno after SetErr was called
	goerr error
)

func strError(e int) error {
	if e == 0 {
		return nil
	}
	if goerr != nil && errno == e {
		return goerr
	}
	return syscall.Errno(e)
}

// StrError returns a C string for the current error, if any.
func StrError(e int) *byte {
	return CString(strError(e).Error())
}

// SetErr sets the Errno value to a specified Go error equivalent.
func SetErr(err error) {
	code := ErrCode(err)
	Errno = code
	// preserve  original Go errors as well
	errno = code
	goerr = err
}

// ErrCode returns an error code corresponding to a Go error.
func ErrCode(err error) int {
	if err == nil {
		return 0
	}
	if os.IsPermission(err) {
		return 13
	} else if os.IsNotExist(err) {
		return 2
	} else if os.IsExist(err) {
		return 17
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return int(errno)
	}
	return 1
}

// Error returns a Go error value that corresponds to the current Errno.
func Error() error {
	if Errno == 0 {
		return nil
	}
	return strError(Errno)
}
