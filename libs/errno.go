package libs

import (
	"fmt"
	"strings"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/9699919799/

func init() {
	RegisterLibrary("errno.h", func(c *Env) *Library {
		return &Library{
			Header: errnoHeader(),
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Idents: errnoIdents(c),
		}
	})
}

type errnoInf struct {
	name  string
	value int
}

func errnoHeader() string {

	lines := []string{
		"#include <" + BuiltinH + ">",
		"_cxgo_go_int errno = 0;",
		"char* strerror (_cxgo_go_int errnum);",
	}

	for _, v := range errnoInfs {
		lines = append(lines, fmt.Sprintf("const int %s = %d;", v.name, v.value))
	}
	res := strings.Join(lines, "\n")
	return res
}

func errnoIdents(c *Env) map[string]*types.Ident {
	gint := c.Go().Int()

	res := map[string]*types.Ident{
		"errno":    c.NewIdent("errno", "libc.Errno", libc.Errno, gint),
		"strerror": c.NewIdent("strerror", "libc.StrError", libc.StrError, c.FuncTT(c.C().String(), gint)),
	}

	for _, v := range errnoInfs {
		res[v.name] = c.NewIdent(v.name, "libc."+v.name, libc.Errno, gint)
	}

	return res
}

var errnoInfs []errnoInf = []errnoInf{
	// Argument list too long.
	{"E2BIG", 2},
	// Permission denied.
	{"EACCES", 3},
	// Address in use.
	{"EADDRINUSE", 4},
	// Address not available.
	{"EADDRNOTAVAIL", 5},
	// Address family not supported.
	{"EAFNOSUPPORT", 6},
	// Resource unavailable, try again (may be the same value as EWOULDBLOCK).
	{"EAGAIN", 7},
	// Connection already in progress.
	{"EALREADY", 8},
	// Bad file descriptor.
	{"EBADF", 9},
	// Bad message.
	{"EBADMSG", 10},
	// Device or resource busy.
	{"EBUSY", 11},
	// Operation canceled.
	{"ECANCELED", 12},
	// No child processes.
	{"ECHILD", 13},
	// Connection aborted.
	{"ECONNABORTED", 14},
	// Connection refused.
	{"ECONNREFUSED", 15},
	// Connection reset.
	{"ECONNRESET", 16},
	// Resource deadlock would occur.
	{"EDEADLK", 17},
	// Destination address required.
	{"EDESTADDRREQ", 18},
	// Mathematics argument out of domain of function.
	{"EDOM", 19},
	// Reserved.
	{"EDQUOT", 20},
	// File exists.
	{"EEXIST", 21},
	// Bad address.
	{"EFAULT", 22},
	// File too large.
	{"EFBIG", 23},
	// Host is unreachable.
	{"EHOSTUNREACH", 24},
	// Identifier removed.
	{"EIDRM", 25},
	// Illegal byte sequence.
	{"EILSEQ", 26},
	// Operation in progress.
	{"EINPROGRESS", 27},
	// Interrupted function.
	{"EINTR", 28},
	// Invalid argument.
	{"EINVAL", 29},
	// I/O error.
	{"EIO", 30},
	// Socket is connected.
	{"EISCONN", 31},
	// Is a directory.
	{"EISDIR", 32},
	// Too many levels of symbolic links.
	{"ELOOP", 33},
	// File descriptor value too large.
	{"EMFILE", 34},
	// Too many links.
	{"EMLINK", 35},
	// Message too large.
	{"EMSGSIZE", 36},
	// Reserved.
	{"EMULTIHOP", 37},
	// Filename too long.
	{"ENAMETOOLONG", 38},
	// Network is down.
	{"ENETDOWN", 39},
	// Connection aborted by network.
	{"ENETRESET", 40},
	// Network unreachable.
	{"ENETUNREACH", 41},
	// Too many files open in system.
	{"ENFILE", 42},
	// No buffer space available.
	{"ENOBUFS", 43},
	// No message is available on the STREAM head read queue.
	{"ENODATA", 44},
	// No such device.
	{"ENODEV", 45},
	// No such file or directory.
	{"ENOENT", 46},
	// Executable file format error.
	{"ENOEXEC", 47},
	// No locks available.
	{"ENOLCK", 48},
	// Reserved.
	{"ENOLINK", 49},
	// Not enough space.
	{"ENOMEM", 50},
	// No message of the desired type.
	{"ENOMSG", 51},
	// Protocol not available.
	{"ENOPROTOOPT", 52},
	// No space left on device.
	{"ENOSPC", 53},
	// No STREAM resources.
	{"ENOSR", 54},
	// Not a STREAM.
	{"ENOSTR", 55},
	// Functionality not supported.
	{"ENOSYS", 56},
	// The socket is not connected.
	{"ENOTCONN", 57},
	// Not a directory or a symbolic link to a directory.
	{"ENOTDIR", 58},
	// Directory not empty.
	{"ENOTEMPTY", 59},
	// Env not recoverable.
	{"ENOTRECOVERABLE", 60},
	// Not a socket.
	{"ENOTSOCK", 61},
	// Not supported (may be the same value as EOPNOTSUPP).
	{"ENOTSUP", 62},
	// Inappropriate I/O control operation.
	{"ENOTTY", 63},
	// No such device or address.
	{"ENXIO", 64},
	// Operation not supported on socket (may be the same value as ENOTSUP).
	{"EOPNOTSUPP", 65},
	// Value too large to be stored in data type.
	{"EOVERFLOW", 66},
	// Previous owner died.
	{"EOWNERDEAD", 67},
	// Operation not permitted.
	{"EPERM", 68},
	// Broken pipe.
	{"EPIPE", 69},
	// Protocol error.
	{"EPROTO", 79},
	// Protocol not supported.
	{"EPROTONOSUPPORT", 80},
	// Protocol wrong type for socket.
	{"EPROTOTYPE", 81},
	// Result too large.
	{"ERANGE", 82},
	// Read-only file system.
	{"EROFS", 83},
	// Invalid seek.
	{"ESPIPE", 84},
	// No such process.
	{"ESRCH", 85},
	// Reserved.
	{"ESTALE", 86},
	// Stream ioctl() timeout.
	{"ETIME", 87},
	// Connection timed out.
	{"ETIMEDOUT", 88},
	// Text file busy.
	{"ETXTBSY", 89},
	// Operation would block (may be the same value as EAGAIN).
	{"EWOULDBLOCK", 90},
	// Cross-device link. },
	{"EXDEV", 91},
}
