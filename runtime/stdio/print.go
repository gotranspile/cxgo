package stdio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dennwc/cxgo/runtime/libc"
)

func FprintfGo(w io.Writer, format string, args ...interface{}) (int, error) {
	words := parseFormat(format)
	var (
		goFormat bytes.Buffer
		goArgs   = make([]interface{}, 0, len(args))
		left     = args
	)
	popArg := func() (interface{}, error) {
		if len(left) == 0 {
			return nil, errors.New("not enough arguments to print")
		}
		v := left[0]
		left = left[1:]
		return v, nil
	}
	for _, w := range words {
		if !w.Verb {
			goFormat.WriteString(w.Str)
			continue
		}
		if strings.HasPrefix(w.Str, "%*") {
			// scanf-specific
			return 0, errors.New("cannot skip args in printf")
		}
		switch w.Str {
		case "%%":
			goFormat.WriteString("%%")
			continue
		case "%c":
			v, err := popArg()
			if err != nil {
				return 0, err
			}
			iv, err := asUint(v)
			if err != nil {
				return 0, err
			}
			goFormat.WriteString("%c")
			goArgs = append(goArgs, rune(iv))
			continue
		}
		last := w.Str[len(w.Str)-1]
		switch last {
		case 's', 'S':
			v, err := popArg()
			if err != nil {
				return 0, err
			}
			s, err := asString(v)
			if err != nil {
				return 0, err
			}
			goFormat.WriteString(strings.ToLower(w.Str))
			goArgs = append(goArgs, s)
			continue
		case 'p':
			v, err := popArg()
			if err != nil {
				return 0, err
			}
			p, err := libc.AsPtr(v)
			if err != nil {
				return 0, err
			}
			goFormat.WriteString(w.Str)
			goArgs = append(goArgs, p)
			continue
		case 'i', 'd', 'u', 'o', 'x':
			if last == 'i' || last == 'u' {
				w.Str = w.Str[:len(w.Str)-1] + "d"
			}
			w.Str = strings.ReplaceAll(w.Str, "l", "")
			v, err := popArg()
			if err != nil {
				return 0, err
			}
			d, err := asUint(v)
			if err != nil {
				return 0, err
			}
			goFormat.WriteString(w.Str)
			switch last {
			case 'i', 'd':
				goArgs = append(goArgs, int64(d))
			default:
				goArgs = append(goArgs, d)
			}
			continue
		case 'f', 'e', 'g', 'F', 'E', 'G':
			v, err := popArg()
			if err != nil {
				return 0, err
			}
			f, err := asFloat(v)
			if err != nil {
				return 0, err
			}
			goFormat.WriteString(w.Str)
			goArgs = append(goArgs, f)
			continue
		default:
			return 0, fmt.Errorf("unsupported verb: %q", w.Str)
		}
	}
	return fmt.Fprintf(w, goFormat.String(), goArgs...)
}

func FprintlnfGo(w io.Writer, format string, args ...interface{}) (int, error) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	return FprintfGo(w, format, args...)
}

func Printf(format string, args ...interface{}) int {
	n, _ := FprintfGo(os.Stdout, format, args...)
	return n
}

func Vprintf(format string, args libc.ArgList) int {
	return Printf(format, args.Args()...)
}

func Dprintf(format string, args ...interface{}) int {
	n, _ := FprintlnfGo(os.Stderr, format, args...)
	return n
}

func Sprintf(buf *byte, format string, args ...interface{}) int {
	var b bytes.Buffer
	n, _ := FprintfGo(&b, format, args...)
	dst := libc.BytesN(buf, b.Len()+1)
	copy(dst, b.Bytes())
	dst[b.Len()] = 0
	return n
}

func Vsprintf(buf *byte, format string, args libc.ArgList) int {
	return Sprintf(buf, format, args.Args()...)
}

func Snprintf(buf *byte, sz int, format string, args ...interface{}) int {
	var b bytes.Buffer
	_, _ = FprintfGo(&b, format, args...)
	dst := libc.BytesN(buf, sz)
	n := copy(dst, b.Bytes())
	if b.Len() < len(dst) {
		dst[b.Len()] = 0
	}
	return n
}

func Vsnprintf(buf *byte, sz int, format string, args libc.ArgList) int {
	return Snprintf(buf, sz, format, args.Args()...)
}

func Fprintf(file *File, format string, args ...interface{}) int {
	n, err := FprintfGo(file.file, format, args...)
	if err != nil {
		file.err = err
		return -1
	}
	return n
}

func Vfprintf(file *File, format string, args libc.ArgList) int {
	return Fprintf(file, format, args.Args()...)
}
