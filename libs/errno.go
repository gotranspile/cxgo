package libs

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/

func init() {
	RegisterLibrary("errno.h", func(c *Env) *Library {
		gint := c.Go().Int()
		return &Library{
			Header: `
#include <` + BuiltinH + `>

_cxgo_go_int errno = 0;

char* strerror (_cxgo_go_int errnum);

// Argument list too long.
#define E2BIG 2
// Permission denied.
#define EACCES 3
// Address in use.
#define EADDRINUSE 4
// Address not available.
#define EADDRNOTAVAIL 5
// Address family not supported.
#define EAFNOSUPPORT 6
// Resource unavailable, try again (may be the same value as EWOULDBLOCK).
#define EAGAIN 7
// Connection already in progress.
#define EALREADY 8
// Bad file descriptor.
#define EBADF 9
// Bad message.
#define EBADMSG 10
// Device or resource busy.
#define EBUSY 11
// Operation canceled.
#define ECANCELED 12
// No child processes.
#define ECHILD 13
// Connection aborted.
#define ECONNABORTED 14
// Connection refused.
#define ECONNREFUSED 15
// Connection reset.
#define ECONNRESET 16
// Resource deadlock would occur.
#define EDEADLK 17
// Destination address required.
#define EDESTADDRREQ 18
// Mathematics argument out of domain of function.
#define EDOM 19
// Reserved.
#define EDQUOT 20
// File exists.
#define EEXIST 21
// Bad address.
#define EFAULT 22
// File too large.
#define EFBIG 23
// Host is unreachable.
#define EHOSTUNREACH 24
// Identifier removed.
#define EIDRM 25
// Illegal byte sequence.
#define EILSEQ 26
// Operation in progress.
#define EINPROGRESS 27
// Interrupted function.
#define EINTR 28
// Invalid argument.
#define EINVAL 29
// I/O error.
#define EIO 30
// Socket is connected.
#define EISCONN 31
// Is a directory.
#define EISDIR 32
// Too many levels of symbolic links.
#define ELOOP 33
// File descriptor value too large.
#define EMFILE 34
// Too many links.
#define EMLINK 35
// Message too large.
#define EMSGSIZE 36
// Reserved.
#define EMULTIHOP 37
// Filename too long.
#define ENAMETOOLONG 38
// Network is down.
#define ENETDOWN 39
// Connection aborted by network.
#define ENETRESET 40
// Network unreachable.
#define ENETUNREACH 41
// Too many files open in system.
#define ENFILE 42
// No buffer space available.
#define ENOBUFS 43
// No message is available on the STREAM head read queue.
#define ENODATA 44
// No such device.
#define ENODEV 45
// No such file or directory.
#define ENOENT 46
// Executable file format error.
#define ENOEXEC 47
// No locks available.
#define ENOLCK 48
// Reserved.
#define ENOLINK 49
// Not enough space.
#define ENOMEM 50
// No message of the desired type.
#define ENOMSG 51
// Protocol not available.
#define ENOPROTOOPT 52
// No space left on device.
#define ENOSPC 53
// No STREAM resources.
#define ENOSR 54
// Not a STREAM.
#define ENOSTR 55
// Functionality not supported.
#define ENOSYS 56
// The socket is not connected.
#define ENOTCONN 57
// Not a directory or a symbolic link to a directory.
#define ENOTDIR 58
// Directory not empty.
#define ENOTEMPTY 59
// Env not recoverable.
#define ENOTRECOVERABLE 60
// Not a socket.
#define ENOTSOCK 61
// Not supported (may be the same value as EOPNOTSUPP).
#define ENOTSUP 62
// Inappropriate I/O control operation.
#define ENOTTY 63
// No such device or address.
#define ENXIO 64
// Operation not supported on socket (may be the same value as ENOTSUP).
#define EOPNOTSUPP 65
// Value too large to be stored in data type.
#define EOVERFLOW 66
// Previous owner died.
#define EOWNERDEAD 67
// Operation not permitted.
#define EPERM 68
// Broken pipe.
#define EPIPE 69
// Protocol error.
#define EPROTO 79
// Protocol not supported.
#define EPROTONOSUPPORT 80
// Protocol wrong type for socket.
#define EPROTOTYPE 81
// Result too large.
#define ERANGE 82
// Read-only file system.
#define EROFS 83
// Invalid seek.
#define ESPIPE 84
// No such process.
#define ESRCH 85
// Reserved.
#define ESTALE 86
// Stream ioctl() timeout.
#define ETIME 87
// Connection timed out.
#define ETIMEDOUT 88
// Text file busy.
#define ETXTBSY 89
// Operation would block (may be the same value as EAGAIN).
#define EWOULDBLOCK 90
// Cross-device link. 
#define EXDEV 91
`,
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Idents: map[string]*types.Ident{
				"errno":    c.NewIdent("errno", "libc.Errno", libc.Errno, gint),
				"strerror": c.NewIdent("strerror", "libc.StrError", libc.StrError, c.FuncTT(c.C().String(), gint)),
			},
		}
	})
}
