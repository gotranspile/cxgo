package stdio

import (
	"log"
	"math"
	"os"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/runtime/libc"
)

func Create(path *byte, mode csys.Mode) uintptr {
	panic("TODO")
}

func Open(path *byte, flags int32, ctls ...interface{}) uintptr {
	spath := libc.GoString(path)
	if len(ctls) != 0 {
		log.Printf("FIXME: open(%q, %x): ignoring controls: %#v", spath, flags, ctls)
	}
	if flags != 0 {
		log.Printf("FIXME: open(%q, %x): ignoring flags: 0x%x", spath, flags, flags)
	}
	// FIXME: handle flags
	f := defaultFS.Open(spath, os.O_RDONLY)
	log.Printf("open(%q, %x): %v", spath, flags, f)
	if f == nil {
		return 0
	}
	return f.FileNo()
}

func FDControl(fd uintptr, flags int32, ctls ...interface{}) int32 {
	panic("TODO")
}

func Chdir(path *byte) int32 {
	spath := libc.GoString(path)
	err := FS().Chdir(spath)
	log.Printf("chdir(%q): %v", spath, err)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return 0
}

func Rmdir(path *byte) int32 {
	name := libc.GoString(path)
	err := FS().Rmdir(name)
	log.Printf("rmdir(%q): %v", name, err)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return 0
}

func Unlink(path *byte) int32 {
	name := libc.GoString(path)
	err := FS().Unlink(name)
	log.Printf("unlink(%q): %v", name, err)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return 0
}

func Access(path *byte, flags int32) int32 {
	name := libc.GoString(path)
	_, err := FS().Stat(name)
	log.Printf("access(%q, %x): %v", name, flags, err)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	// TODO: check access bits
	return 0
}

func Lseek(fd uintptr, offs uint64, whence int32) uint64 {
	f := ByFD(fd)
	if f == nil {
		return math.MaxUint64
	}
	off := f.Seek(int64(offs), whence)
	if off < 0 {
		return math.MaxUint64
	}
	return uint64(off)
}

func GetCwd(p *byte, sz int) *byte {
	dir, err := FS().Getwd()
	if err != nil {
		libc.SetErr(err)
		return nil
	}
	dst := unsafe.Slice(p, sz)
	copy(dst, dir)
	return p
}
