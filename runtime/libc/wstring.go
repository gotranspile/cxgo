package libc

import (
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"
)

const WCharSize = int(unsafe.Sizeof(WChar(0)))

// IndexWCharPtr unsafely moves a wide char pointer by i elements. An offset may be negative.
//
// Deprecated: use unsafe.Add
func IndexWCharPtr(p *WChar, i int) *WChar {
	return (*WChar)(unsafe.Add(unsafe.Pointer(p), i*WCharSize))
}

// UnsafeWCharN makes a slice of a given size starting at ptr.
//
// Deprecated: use unsafe.Slice
func UnsafeWCharN(ptr unsafe.Pointer, sz uint64) []WChar {
	return unsafe.Slice((*WChar)(ptr), sz)
}

// WCharN makes a slice of a given size starting at ptr.
// It accepts a *WChar instead of unsafe pointer as UnsafeWCharN does, which allows to avoid unsafe import.
//
// Deprecated: use unsafe.Slice
func WCharN(p *WChar, sz uint64) []WChar {
	return unsafe.Slice(p, sz)
}

// CWString makes a new zero-terminated wide char array containing a given string.
func CWString(s string) *WChar {
	sz := utf8.RuneCountInString(s)
	p := makePad((sz+2)*int(unsafe.Sizeof(WChar(0))), 0)

	w := unsafe.Slice((*WChar)(unsafe.Pointer(&p[0])), uint64(sz))
	w = w[:0]
	for _, r := range s {
		w = append(w, WChar(r))
	}
	if cap(w) != sz {
		panic("must not happen")
	}
	return &w[0]
}

func WStrLen(str *WChar) uint64 {
	return uint64(findnull(str))
}

func GoWSlice(ptr *WChar) []WChar {
	n := WStrLen(ptr)
	if n == 0 {
		return nil
	}
	return unsafe.Slice(ptr, n)
}

func GoWString(s *WChar) string {
	arr := GoWSlice(s)
	runes := make([]rune, 0, len(arr))
	for _, r := range arr {
		runes = append(runes, rune(r))
	}
	return string(runes)
}

func WStrCpy(dst, src *WChar) *WChar {
	s := GoWSlice(src)
	d := unsafe.Slice(dst, len(s)+1)
	n := copy(d, s)
	d[n] = 0
	return dst
}

func WStrNCpy(dst, src *WChar, sz uint32) *WChar {
	d := unsafe.Slice(dst, sz)
	s := GoWSlice(src)
	pad := 0
	if len(s) > int(sz) {
		s = s[:sz]
	} else if len(s) < int(sz) {
		pad = int(sz) - len(s)
	}
	n := copy(d, s)
	for i := 0; i < pad; i++ {
		d[n+i] = 0
	}
	return dst
}

func WStrChr(str *WChar, ch int64) *WChar {
	if ch < 0 || ch > 0xffff {
		panic(ch)
	}
	if str == nil {
		return nil
	}
	b := GoWSlice(str)
	for i, c := range b {
		if c == WChar(ch) {
			return &b[i]
		}
	}
	return nil
}

func WStrCat(dst, src *WChar) *WChar {
	s := GoWSlice(src)
	i := WStrLen(dst)
	n := int(i) + len(s)
	d := unsafe.Slice(dst, n+1)
	copy(d[i:], s)
	d[n] = 0
	return &d[0]
}

func WStrCmp(a, b *WChar) int {
	s1 := GoWString(a)
	s2 := GoWString(b)
	return strings.Compare(s1, s2)
}

func WStrCaseCmp(a, b *WChar) int {
	s1 := strings.ToLower(GoWString(a))
	s2 := strings.ToLower(GoWString(b))
	return strings.Compare(s1, s2)
}

func WStrtol(s *WChar, end **WChar, base int) int {
	if end != nil {
		panic("TODO")
	}
	str := GoWString(s)
	v, err := strconv.ParseInt(str, base, 32)
	if err != nil {
		return 0
	}
	return int(v)
}

func Mbstowcs(a1 *WChar, a2 *byte, a3 uint32) uint32 {
	panic("TODO")
}

func wIndexAny(s, chars []WChar) int {
	for i, c1 := range s {
		for _, c2 := range chars {
			if c1 == c2 {
				return i
			}
		}
	}
	return -1
}

func wIndex(s []WChar, c WChar) int {
	for i, c1 := range s {
		if c1 == c {
			return i
		}
	}
	return -1
}

var wstrtok struct {
	data []WChar
	ind  int
}

func WStrTok(src, delim *WChar) *WChar {
	if src != nil {
		wstrtok.data = GoWSlice(src)
		wstrtok.ind = 0
	}
	d := GoWSlice(delim)
	for ; wstrtok.ind < len(wstrtok.data); wstrtok.ind++ {
		if wIndex(d, wstrtok.data[wstrtok.ind]) < 0 {
			// start of a new token
			tok := wstrtok.data[wstrtok.ind:]
			if i := wIndexAny(tok, d); i >= 0 {
				tok[i] = 0
				wstrtok.ind += i + 1
			} else {
				wstrtok.data = nil
				wstrtok.ind = 0
			}
			return &tok[0]
		}
		// skip delimiters
	}
	wstrtok.data = nil
	wstrtok.ind = 0
	return nil
}
