package cxgo

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/dennwc/cxgo/libs"
	"modernc.org/cc/v3"
)

var timeRun = time.Now()

func newIncludeFS(c *libs.Env) cc.Filesystem {
	return includeFS{c: c}
}

type includeFS struct {
	c *libs.Env
}

func (fs includeFS) content(path string, sys bool) (string, error) {
	if !sys {
		return "", os.ErrNotExist
	}
	l, ok := fs.c.NewLibrary(path)
	if !ok {
		return "", os.ErrNotExist
	}
	return l.Header, nil
}

func (fs includeFS) Stat(path string, sys bool) (os.FileInfo, error) {
	data, err := fs.content(path, sys)
	if err != nil {
		return nil, err
	}
	return includeFI{name: path, data: data}, nil
}

func (fs includeFS) Open(path string, sys bool) (io.ReadCloser, error) {
	data, err := fs.content(path, sys)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(strings.NewReader(data)), nil
}

type includeFI struct {
	name string
	data string
}

func (fi includeFI) Name() string {
	return fi.name
}

func (fi includeFI) Size() int64 {
	return int64(len(fi.data))
}

func (fi includeFI) Mode() os.FileMode {
	return 0
}

func (fi includeFI) ModTime() time.Time {
	return timeRun
}

func (fi includeFI) IsDir() bool {
	return false
}

func (fi includeFI) Sys() interface{} {
	return fi
}
