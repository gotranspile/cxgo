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

// Argument list too long.
const E2BIG = 2

// Permission denied.
const EACCES = 3

// Address in use.
const EADDRINUSE = 4

// Address not available.
const EADDRNOTAVAIL = 5

// Address family not supported.
const EAFNOSUPPORT = 6

// Resource unavailable, try again (may be the same value as EWOULDBLOCK).
const EAGAIN = 7

// Connection already in progress.
const EALREADY = 8

// Bad file descriptor.
const EBADF = 9

// Bad message.
const EBADMSG = 10

// Device or resource busy.
const EBUSY = 11

// Operation canceled.
const ECANCELED = 12

// No child processes.
const ECHILD = 13

// Connection aborted.
const ECONNABORTED = 14

// Connection refused.
const ECONNREFUSED = 15

// Connection reset.
const ECONNRESET = 16

// Resource deadlock would occur.
const EDEADLK = 17

// Destination address required.
const EDESTADDRREQ = 18

// Mathematics argument out of domain of function.
const EDOM = 19

// Reserved.
const EDQUOT = 20

// File exists.
const EEXIST = 21

// Bad address.
const EFAULT = 22

// File too large.
const EFBIG = 23

// Host is unreachable.
const EHOSTUNREACH = 24

// Identifier removed.
const EIDRM = 25

// Illegal byte sequence.
const EILSEQ = 26

// Operation in progress.
const EINPROGRESS = 27

// Interrupted function.
const EINTR = 28

// Invalid argument.
const EINVAL = 29

// I/O error.
const EIO = 30

// Socket is connected.
const EISCONN = 31

// Is a directory.
const EISDIR = 32

// Too many levels of symbolic links.
const ELOOP = 33

// File descriptor value too large.
const EMFILE = 34

// Too many links.
const EMLINK = 35

// Message too large.
const EMSGSIZE = 36

// Reserved.
const EMULTIHOP = 37

// Filename too long.
const ENAMETOOLONG = 38

// Network is down.
const ENETDOWN = 39

// Connection aborted by network.
const ENETRESET = 40

// Network unreachable.
const ENETUNREACH = 41

// Too many files open in system.
const ENFILE = 42

// No buffer space available.
const ENOBUFS = 43

// No message is available on the STREAM head read queue.
const ENODATA = 44

// No such device.
const ENODEV = 45

// No such file or directory.
const ENOENT = 46

// Executable file format error.
const ENOEXEC = 47

// No locks available.
const ENOLCK = 48

// Reserved.
const ENOLINK = 49

// Not enough space.
const ENOMEM = 50

// No message of the desired type.
const ENOMSG = 51

// Protocol not available.
const ENOPROTOOPT = 52

// No space left on device.
const ENOSPC = 53

// No STREAM resources.
const ENOSR = 54

// Not a STREAM.
const ENOSTR = 55

// Functionality not supported.
const ENOSYS = 56

// The socket is not connected.
const ENOTCONN = 57

// Not a directory or a symbolic link to a directory.
const ENOTDIR = 58

// Directory not empty.
const ENOTEMPTY = 59

// Env not recoverable.
const ENOTRECOVERABLE = 60

// Not a socket.
const ENOTSOCK = 61

// Not supported (may be the same value as EOPNOTSUPP).
const ENOTSUP = 62

// Inappropriate I/O control operation.
const ENOTTY = 63

// No such device or address.
const ENXIO = 64

// Operation not supported on socket (may be the same value as ENOTSUP).
const EOPNOTSUPP = 65

// Value too large to be stored in data type.
const EOVERFLOW = 66

// Previous owner died.
const EOWNERDEAD = 67

// Operation not permitted.
const EPERM = 68

// Broken pipe.
const EPIPE = 69

// Protocol error.
const EPROTO = 79

// Protocol not supported.
const EPROTONOSUPPORT = 80

// Protocol wrong type for socket.
const EPROTOTYPE = 81

// Result too large.
const ERANGE = 82

// Read-only file system.
const EROFS = 83

// Invalid seek.
const ESPIPE = 84

// No such process.
const ESRCH = 85

// Reserved.
const ESTALE = 86

// Stream ioctl() timeout.
const ETIME = 87

// Connection timed out.
const ETIMEDOUT = 88

// Text file busy.
const ETXTBSY = 89

// Operation would block (may be the same value as EAGAIN).
const EWOULDBLOCK = 90

// Cross-device link.
const EXDEV = 91
