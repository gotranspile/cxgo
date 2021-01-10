package stdio

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/dennwc/cxgo/runtime/libc"
)

const GlobNoEscape = 1

type Glob struct {
	Num     int32
	Paths   **byte
	Reserve int32
}

func isHidden(s string) bool {
	for _, p := range filepath.SplitList(s) {
		if strings.HasPrefix(p, ".") {
			return true
		}
	}
	return false
}

func (g *Glob) Glob(pattern *byte, flags int32, errFunc func(epath *byte, errno int32) int32) int32 {
	if errFunc != nil {
		panic("implement me")
	}
	glob := libc.GoString(pattern)
	res, err := filepath.Glob(glob)
	log.Printf("glob(%q, %x, %p, %p): %d, %v", glob, flags, errFunc, g, len(res), err)
	if err != nil {
		libc.SetErr(err)
		return 1 // TODO
	}
	g.Reserve = 0
	paths := make([]*byte, 0, len(res)+1+2)
	paths = append(paths, nil) // padding
	for _, s := range res {
		if isHidden(s) {
			continue
		}
		paths = append(paths, libc.CString(s))
	}
	g.Num = int32(len(paths) - 1) // - padding
	if g.Num == 0 {
		g.Paths = nil
		return 1 // TODO
	}
	paths = append(paths, nil) // null-terminator
	paths = append(paths, nil) // padding
	g.Paths = &paths[1]        // padding
	return 0
}

func (g *Glob) Free() {
	g.Paths = nil
	g.Num = 0
	g.Reserve = 0
}
