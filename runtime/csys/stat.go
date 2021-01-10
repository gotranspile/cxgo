package csys

import (
	"log"
	"os"

	"github.com/dennwc/cxgo/runtime/libc"
)

const modeDir = 1

type Mode int32

func IsDir(mode Mode) int32 {
	if mode&modeDir != 0 {
		return 1
	}
	return 0
}

type StatRes struct {
	Dev       int32
	Inode     int32
	Mode      Mode
	Links     int32
	UID       int32
	GID       int32
	RDev      int32
	Size      uint64
	ATime     libc.TimeVal
	MTime     libc.TimeVal
	CTime     libc.TimeVal
	BlockSize int32
	Blocks    int32
}

func Stat(path *byte, dst *StatRes) int32 {
	name := libc.GoString(path)
	if name == "" {
		panic("empty name")
	}
	st, err := os.Stat(name)
	if err != nil {
		log.Printf("stat(%q, %p): %v", name, dst, err)
		libc.SetErr(err)
		return -1
	}
	dst.Mode = 0
	if st.Mode().IsDir() {
		dst.Mode = modeDir
	}
	dst.Size = uint64(st.Size())
	// TODO: other fields
	return 0
}

func Chmod(path *byte, mode Mode) int32 {
	spath := libc.GoString(path)
	log.Printf("TODO: chmod(%q, 0%o)", spath, mode) // TODO
	return 0
}

func Mkdir(path *byte, mode Mode) int32 {
	spath := libc.GoString(path)
	err := os.Mkdir(spath, 0755)
	log.Printf("mkdir(%q, %x): %v", spath, mode, err)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return 0
}
