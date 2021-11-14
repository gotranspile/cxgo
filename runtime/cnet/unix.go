package cnet

import (
	"log"
	"os"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

func GetHostname(buf *byte, sz int) int {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("gethostname: %v", err)
		libc.SetErr(err)
		return -1
	}
	b := unsafe.Slice(buf, sz)
	n := copy(b, name)
	if n+1 < len(b) {
		b[n] = 0
	}
	return 0
}
