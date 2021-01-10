package stdio

import "os"

func NewLocalFS() Filesystem {
	return localFS{}
}

type localFS struct{}

func (localFS) Stdout() FileI {
	return os.Stdout
}

func (localFS) Stderr() FileI {
	return os.Stderr
}

func (localFS) Stdin() FileI {
	return os.Stdin
}

func (localFS) Getwd() (string, error) {
	return os.Getwd()
}

func (localFS) Chdir(path string) error {
	return os.Chdir(path)
}

func (localFS) Rmdir(path string) error {
	return os.RemoveAll(path)
}

func (localFS) Unlink(path string) error {
	return os.Remove(path)
}

func (localFS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (localFS) Open(path string, flag int, mode os.FileMode) (FileI, error) {
	return os.OpenFile(path, flag, mode)
}
