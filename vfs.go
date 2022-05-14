package cxgo

import (
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/gotranspile/cxgo/libs"
)

var timeRun = time.Now()

func newIncludeFS(c *libs.Env) fs.StatFS {
	return includeFS{c: c}
}

type includeFS struct {
	c *libs.Env
}

func (ifs includeFS) content(path string) (string, error) {
	l, ok := ifs.c.NewLibrary(path)
	if !ok {
		return "", fs.ErrNotExist
	}
	return l.Header, nil
}

func (ifs includeFS) Stat(path string) (fs.FileInfo, error) {
	data, err := ifs.content(path)
	if err != nil {
		return nil, err
	}
	return includeFile{name: path, size: len(data)}, nil
}

func (ifs includeFS) Open(path string) (fs.File, error) {
	data, err := ifs.content(path)
	if err != nil {
		return nil, err
	}
	return includeFile{name: path, size: len(data), r: strings.NewReader(data)}, nil
}

type includeFile struct {
	name string
	size int
	r    io.Reader
}

func (fi includeFile) Name() string {
	return fi.name
}

func (fi includeFile) Size() int64 {
	return int64(fi.size)
}

func (fi includeFile) Mode() fs.FileMode {
	return 0
}

func (fi includeFile) ModTime() time.Time {
	return timeRun
}

func (fi includeFile) IsDir() bool {
	return false
}

func (fi includeFile) Sys() interface{} {
	return fi
}

func (fi includeFile) Stat() (fs.FileInfo, error) {
	return fi, nil
}

func (fi includeFile) Read(p []byte) (int, error) {
	return fi.r.Read(p)
}

func (fi includeFile) Close() error {
	return nil
}
