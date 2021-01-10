package cnet

import (
	"log"
	"os"

	"github.com/dennwc/cxgo/runtime/libc"
)

func GetHostname(buf *byte, sz int) int {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("gethostname: %v", err)
		libc.SetErr(err)
		return -1
	}
	b := libc.BytesN(buf, int(sz))
	n := copy(b, name)
	if n+1 < len(b) {
		b[n] = 0
	}
	return 0
}
