package csys

import "os"

const (
	O_RDONLY = int32(os.O_RDONLY)
	O_WRONLY = int32(os.O_WRONLY)
	O_RDWR   = int32(os.O_RDWR)
	O_CREAT  = int32(os.O_CREATE)
	O_EXCL   = int32(os.O_EXCL)
	O_TRUNC  = int32(os.O_TRUNC)
)
