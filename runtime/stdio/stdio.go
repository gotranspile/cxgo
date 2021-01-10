package stdio

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/dennwc/cxgo/runtime/libc"
)

const (
	defPermFile = 0644
)

const (
	EOF      = -2
	SEEK_SET = int32(io.SeekStart)
	SEEK_CUR = int32(io.SeekCurrent)
	SEEK_END = int32(io.SeekEnd)
)

func Stdout() *File {
	return defaultFS.OpenFrom(defaultFS.fs.Stdout())
}

func Stderr() *File {
	return defaultFS.OpenFrom(defaultFS.fs.Stderr())
}

func Stdin() *File {
	return defaultFS.OpenFrom(defaultFS.fs.Stdin())
}

func Remove(path string, _ ...interface{}) int {
	panic("TODO")
}

func Rename(path1, path2 string, _ ...interface{}) int {
	panic("TODO")
}

func openFlags(mode string) int {
	mode = strings.ReplaceAll(mode, "b", "")
	switch mode {
	case "r": // open file for reading
		return os.O_RDONLY
	case "w": // truncate to zero length or create file for writing
		return os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	case "a": // append; open or create file for writing at end-of-file
		return os.O_CREATE | os.O_APPEND | os.O_WRONLY
	case "r+": // open file for update (reading and writing)
		return os.O_RDWR
	case "w+": // truncate to zero length or create file for update
		return os.O_CREATE | os.O_TRUNC | os.O_RDWR
	case "a+": // append; open or create file for update, writing at end-of-file
		return os.O_CREATE | os.O_APPEND | os.O_RDWR
	default:
		panic("unknown file mode: " + mode)
	}
}

func FOpen(path, mode string) *File {
	return defaultFS.OpenS(path, mode)
}

func OpenFrom(f FileI) *File {
	return defaultFS.OpenFrom(f)
}

func (fs *filesystem) OpenFrom(f FileI) *File {
	ff := &File{
		fs:   fs,
		file: f,
		fd:   f.Fd(),
	}
	fs.Lock()
	fs.byFD[ff.fd] = ff
	fs.Unlock()
	return ff
}

func (fs *filesystem) Open(path string, flag int) *File {
	f, err := fs.fs.Open(path, flag, defPermFile)
	log.Printf("fopen(%q, %v): %v", path, flag, err)
	if err != nil {
		libc.SetErr(err)
		return nil
	}
	return fs.OpenFrom(f)
}

func (fs *filesystem) OpenS(path, mode string) *File {
	flags := openFlags(mode)
	return fs.Open(path, flags)
}

func FDOpen(fd uintptr, mode string) *File {
	f := ByFD(fd)
	log.Printf("fdopen(%d, %q): %v", fd, mode, f)
	if f == nil {
		return nil
	}
	flags := openFlags(mode)
	_ = flags // FIXME: use flags
	return f
}

func FDOpenS(fd uintptr, mode string) *File {
	f := ByFD(fd)
	log.Printf("fdopen(%d, %q): %v", fd, mode, f)
	if f == nil {
		return nil
	}
	flags := openFlags(mode)
	_ = flags // FIXME: use flags
	return f
}

func FReOpen(path, mode string, f *File) *File {
	panic("TODO")
}

func Fscanf(file *File, format string, args ...interface{}) int {
	n, err := FscanfGo(file.file, format, args...)
	if err != nil {
		file.err = err
		return -1
	}
	return n
}

func Vfscanf(file *File, format string, args libc.ArgList) int {
	return Fscanf(file, format, args.Args()...)
}

type File struct {
	fs   *filesystem
	fd   uintptr
	file FileI
	err  error
	c    *int
}

func (f *File) SetErr(err error) {
	f.err = err
}

func (f *File) IsEOF() int32 {
	if f.err == io.EOF {
		return 1
	}
	return 0
}

func (f *File) Error() int64 {
	if f.err == nil {
		return 0
	}
	return int64(libc.ErrCode(f.err))
}

func (f *File) FileNo() uintptr {
	return f.fd
}

func (f *File) Flush() int32 {
	err := f.file.Sync()
	if err != nil {
		f.err = err
		return -1
	}
	return 0
}

func (f *File) Close() int32 {
	if f == nil {
		return -1
	}
	if err := f.file.Close(); err != nil {
		f.err = err
	}
	f.fs.Lock()
	delete(f.fs.byFD, f.fd)
	f.fs.Unlock()
	return 0 // TODO
}

func (f *File) WriteN(p *byte, size, cnt int) int32 {
	n := f.Write(p, size*cnt)
	if n > 0 {
		n /= int32(size)
	}
	return n
}

func (f *File) Write(p *byte, sz int) int32 {
	if f == nil {
		return -1
	}
	n, err := f.file.Write(libc.BytesN(p, sz))
	if err != nil {
		f.err = err
	}
	return int32(n)
}

func (f *File) ReadN(p *byte, size, cnt int) int32 {
	n := f.Read(p, size*cnt)
	if n > 0 {
		n /= int32(size)
	}
	return n
}

func (f *File) Read(p *byte, sz int) int32 {
	if f == nil {
		return -1
	} else if sz == 0 {
		return 0
	}
	n, err := f.file.Read(libc.BytesN(p, sz))
	if err != nil {
		f.err = err
	}
	return int32(n)
}

func (f *File) GetC() int {
	if f.c != nil {
		c := *f.c
		f.c = nil
		return c
	}
	var b [1]byte
	_, err := f.file.Read(b[:])
	if err != nil {
		//log.Printf("fgetc(): %v", err)
		f.err = err
		if err == io.EOF {
			return EOF
		}
		return -1
	}
	return int(b[0])
}

func (f *File) UnGetC(c int) int {
	f.c = &c
	return 0
}

func (f *File) GetS(buf *byte, sz int32) *byte {
	dst := libc.BytesN(buf, int(sz))
	var b [1]byte
	for len(dst) > 1 {
		_, err := f.file.Read(b[:])
		if err != nil {
			log.Printf("fgets(%q, %d): %v", f.file.Name(), sz, err)
			f.err = err
			return nil
		}
		dst[0] = b[0]
		dst = dst[1:]
		if b[0] == '\n' {
			break
		}
	}
	dst[0] = 0
	return buf
}

func (f *File) PutC(c int) int64 {
	if f == nil {
		return -1
	}
	if c < 0 || c > 0xff {
		panic("TODO")
	}
	n, err := f.file.Write([]byte{byte(c)})
	if err != nil {
		f.err = err
	}
	return int64(n)
}

func (f *File) PutS(s *byte) int64 {
	return int64(f.Write(s, libc.StrLen(s)))
}

func (f *File) Scanf(format string, args ...interface{}) int64 {
	panic("TODO")
}

func (f *File) Tell() int64 {
	cur, err := f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		libc.SetErr(err)
		f.err = err
		return -1
	}
	return cur
}

func (f *File) Seek(off int64, whence int32) int32 {
	_, err := f.file.Seek(off, int(whence))
	if err != nil {
		f.err = err
		return -1
	}
	return 0
}
