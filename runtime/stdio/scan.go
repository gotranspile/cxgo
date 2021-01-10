package stdio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dennwc/cxgo/runtime/libc"
)

func FscanfGo(r io.Reader, format string, args ...interface{}) (int, error) {
	left := args
	var (
		goformat bytes.Buffer
		goargs   []interface{}
		post     []func() error
	)
	popArg := func() (interface{}, error) {
		if len(left) == 0 {
			return nil, errors.New("not enough arguments to scan")
		}
		v := left[0]
		left = left[1:]
		return v, nil
	}
	processArg := func(typ string) error {
		switch typ {
		case "%%":
			goformat.WriteString("%%")
			return nil
		case "%s":
			p, err := popArg()
			if err != nil {
				return err
			}
			goformat.WriteString("%s")
			var v string
			goargs = append(goargs, &v)

			rv, err := libc.AsPtr(p)
			if err != nil {
				panic(err)
			}
			pv := (*byte)(rv)
			post = append(post, func() error {
				libc.StrCpyGoZero(pv, []byte(v))
				return nil
			})
			return nil
		case "%d":
			p, err := popArg()
			if err != nil {
				return err
			}
			goformat.WriteString("%d")
			var v int
			goargs = append(goargs, &v)
			rv, err := libc.AsPtr(p)
			if err != nil {
				panic(err)
			}
			pv := (*int)(rv)
			post = append(post, func() error {
				*pv = v
				return nil
			})
			return nil
		case "%f":
			p, err := popArg()
			if err != nil {
				return err
			}
			goformat.WriteString("%f")
			var v float32
			goargs = append(goargs, &v)
			rv, err := libc.AsPtr(p)
			if err != nil {
				panic(err)
			}
			pv := (*float32)(rv)
			post = append(post, func() error {
				*pv = v
				return nil
			})
			return nil
		case "%*s":
			goformat.WriteString("%s")
			var v string
			goargs = append(goargs, &v)
			// ignore
			return nil
		default:
			panic(typ)
		}
	}
	for _, w := range parseFormat(format) {
		if !w.Verb {
			goformat.WriteString(w.Str)
			continue
		}
		if err := processArg(w.Str); err != nil {
			log.Printf("fscanf(%q, %+v): process: %v", format, args, err)
			return 0, err
		}
	}
	r = &spaceReader{r: r}
	_, err := fmt.Fscanf(r, goformat.String(), goargs...)
	if err != nil {
		log.Printf("fscanf(%q, %+v): scan: %v", format, args, err)
		return 0, err
	}
	for _, fnc := range post {
		if err := fnc(); err != nil {
			log.Printf("fscanf(%q, %+v): post: %v", format, args, err)
			return 0, err
		}
	}
	return len(args), nil
}

// spaceReader replaces new line characters with spaces.
// It's needed because scanf in C interprets new lines as spaces, while Go's fmt.Scanf requires new lines to be included
// in the format string.
type spaceReader struct {
	r io.Reader
}

func (r *spaceReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	for i := range p[:n] {
		switch p[i] {
		case '\n', '\r':
			p[i] = ' '
		}
	}
	return n, err
}

func Scanf(format string, args ...interface{}) int {
	n, err := FscanfGo(os.Stdin, format, args...)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return n
}

func Vscanf(format string, args libc.ArgList) int {
	return Scanf(format, args.Args()...)
}

func Sscanf(buf *byte, format string, args ...interface{}) int {
	str := libc.GoBytes(buf)
	n, err := FscanfGo(bytes.NewReader(str), format, args...)
	if err != nil {
		libc.SetErr(err)
		return -1
	}
	return n
}

func Vsscanf(buf *byte, format string, args libc.ArgList) int {
	return Sscanf(buf, format, args.Args()...)
}
