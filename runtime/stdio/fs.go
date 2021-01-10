package stdio

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/dennwc/cxgo/runtime/libc"
)

type FileI interface {
	Fd() uintptr
	Name() string
	Sync() error
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

type Filesystem interface {
	Stdout() FileI
	Stderr() FileI
	Stdin() FileI
	Getwd() (string, error)
	Chdir(path string) error
	Rmdir(path string) error
	Unlink(path string) error
	Open(path string, flag int, mode os.FileMode) (FileI, error)
	Stat(path string) (os.FileInfo, error)
}

var defaultFS *filesystem

func init() {
	SetFS(nil)
}

func SetFS(fs Filesystem) {
	if fs == nil {
		fs = localFS{}
	}
	// TODO: mount stdout, stderr, stdin
	defaultFS = &filesystem{
		fs:   fs,
		byFD: make(map[uintptr]*File),
	}
}

func FS() Filesystem {
	if defaultFS == nil {
		SetFS(nil)
	}
	return defaultFS.fs
}

func ByFD(fd uintptr) *File {
	f, err := defaultFS.fileByFD(fd)
	if err != nil {
		libc.SetErr(err)
		return nil
	}
	return f
}

type filesystem struct {
	fs Filesystem
	sync.RWMutex
	byFD map[uintptr]*File
}

func (fs *filesystem) fileByFD(fd uintptr) (*File, error) {
	fs.RLock()
	f := fs.byFD[fd]
	fs.RUnlock()
	if f == nil {
		return nil, errors.New("invalid fd")
	}
	return f, nil
}
