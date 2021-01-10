package libc

import (
	"bytes"
	"strings"
	"unsafe"
)

// CString makes a new zero-terminated byte array containing a given string.
func CString(s string) *byte {
	p := makePad(len(s)+1, 0)
	copy(p, s)
	return &p[0]
}

// CBytes makes a new zero-terminated byte array containing a given byte slice.
func CBytes(b []byte) *byte {
	p := makePad(len(b)+1, 0)
	copy(p, b)
	return &p[0]
}

// GoBytes makes a Go byte slice from a pointer to a zero-terminated byte array.
// The slice will point to the same memory as ptr.
func GoBytes(ptr *byte) []byte {
	n := findnull(ptr)
	if n == 0 {
		return nil
	}
	return BytesN(ptr, int(n))
}

// GoString makes a Go string from a pointer to a zero-terminated byte array.
func GoString(s *byte) string {
	return gostring(s)
}

// GoBytesS is a Go-friendly analog of GoBytes.
func GoBytesS(s []byte) []byte {
	n := StrLenS(s)
	if n == 0 {
		return nil
	}
	return s[:n]
}

// GoStringS is a Go-friendly analog of GoString.
func GoStringS(s []byte) string {
	n := StrLenS(s)
	if n == 0 {
		return ""
	}
	return string(s[:n])
}

// CStringSlice convers a Go string slice to a zero-terminated array of C string pointers.
func CStringSlice(arr []string) **byte {
	out := make([]*byte, len(arr)+1)
	for i, s := range arr {
		out[i] = CString(s)
	}
	return &out[0]
}

// MemCmp compares the first count characters of the objects pointed to by lhs and rhs. The comparison is done lexicographically.
//
// The sign of the result is the sign of the difference between the values of the first pair of bytes (both interpreted
// as byte) that differ in the objects being compared.
//
// The behavior is undefined if access occurs beyond the end of either object pointed to by lhs and rhs. The behavior is
// undefined if either lhs or rhs is a null pointer.
func MemCmp(lhs, rhs unsafe.Pointer, sz int) int {
	b1, b2 := UnsafeBytesN(lhs, sz), UnsafeBytesN(rhs, sz)
	return bytes.Compare(b1, b2)
}

// MemSet copies the value ch (after conversion to byte as if by byte(ch)) into each of the first count characters of
// the object pointed to by dest.
//
// The behavior is undefined if access occurs beyond the end of the dest array. The behavior is undefined if dest is a
// null pointer.
func MemSet(p unsafe.Pointer, ch byte, sz int) unsafe.Pointer {
	b := UnsafeBytesN(p, sz)
	if ch == 0 {
		copy(b, make([]byte, len(b)))
	} else {
		copy(b, bytes.Repeat([]byte{ch}, len(b)))
	}
	return p
}

// MemMove copies count characters from the object pointed to by src to the object pointed to by dest. Both objects are
// interpreted as arrays of byte. The objects may overlap: copying takes place as if the characters were copied to a
// temporary character array and then the characters were copied from the array to dest.
//
// The behavior is undefined if access occurs beyond the end of the dest array. The behavior is undefined if either dest
// or src is a null pointer.
func MemMove(dst, src unsafe.Pointer, sz int) unsafe.Pointer {
	if sz == 0 {
		return dst
	}
	return MemCpy(dst, src, sz)
}

// MemCpy copies count characters from the object pointed to by src to the object pointed to by dest. Both objects are
// interpreted as arrays of byte.
//
// The behavior is undefined if access occurs beyond the end of the dest array. If the objects overlap (which is a
// violation of the restrict contract), the behavior is undefined. The behavior is undefined if either dest or src is a
// null pointer.
func MemCpy(dst, src unsafe.Pointer, sz int) unsafe.Pointer {
	if dst == nil {
		panic("nil destination")
	}
	if sz == 0 || src == nil {
		return dst
	}
	bdst := UnsafeBytesN(dst, sz)
	bsrc := UnsafeBytesN(src, sz)
	copy(bdst, bsrc)
	return dst
}

// MemChr finds the first occurrence of ch (after conversion to byte as if by byte(ch)) in the initial count characters
// (each interpreted as byte) of the object pointed to by ptr.
//
// The behavior is undefined if access occurs beyond the end of the array searched. The behavior is undefined if ptr is
// a null pointer.
func MemChr(ptr *byte, ch byte, sz int) *byte {
	if ptr == nil || sz == 0 {
		return nil
	}
	b := BytesN(ptr, sz)
	i := bytes.IndexByte(b, ch)
	if i < 0 {
		return nil
	}
	return &b[i]
}

// StrLen returns the length of the given null-terminated byte string, that is, the number of characters in a character
// array whose first element is pointed to by str up to and not including the first null character.
//
// The behavior is undefined if str is not a pointer to a null-terminated byte string.
func StrLen(str *byte) int {
	return findnull(str)
}

// StrLenS is a Go-friendly analog of StrLen.
func StrLenS(s []byte) int {
	if len(s) == 0 {
		return 0
	}
	i := bytes.IndexByte(s, 0)
	if i < 0 {
		return len(s)
	}
	return i
}

// StrChr finds the first occurrence of ch (after conversion to byte as if by byte(ch)) in the null-terminated byte
// string pointed to by str (each character interpreted as unsigned char). The terminating null character is considered
// to be a part of the string and can be found when searching for '\x00'.
//
// The behavior is undefined if str is not a pointer to a null-terminated byte string.
//
// The return value is a pointer to the found character in str, or null pointer if no such character is found.
func StrChr(str *byte, ch byte) *byte {
	if str == nil {
		return nil
	}
	b := GoBytes(str)
	i := bytes.IndexByte(b, ch)
	if i < 0 {
		return nil
	}
	return &b[i]
}

// StrRChr finds the last occurrence of ch (after conversion to byte as if by byte(ch)) in the null-terminated byte
// string pointed to by str (each character interpreted as unsigned char). The terminating null character is considered
// to be a part of the string and can be found when searching for '\x00'.
//
// The behavior is undefined if str is not a pointer to a null-terminated byte string.
//
// The return value is a pointer to the found character in str, or null pointer if no such character is found.
func StrRChr(str *byte, ch byte) *byte {
	if str == nil {
		return nil
	}
	b := GoBytes(str)
	i := bytes.LastIndexByte(b, ch)
	if i < 0 {
		return nil
	}
	return &b[i]
}

// StrStr finds the first occurrence of the null-terminated byte string pointed to by substr in the null-terminated byte
// string pointed to by str. The terminating null characters are not compared.
//
// The behavior is undefined if either str or substr is not a pointer to a null-terminated byte string.
//
// The return value is a pointer to the first character of the found substring in str, or NULL if such substring is not
// found. If substr points to an empty string, str is returned.
func StrStr(str, substr *byte) *byte {
	if str == nil {
		return nil
	} else if substr == nil {
		return str
	}
	sub := GoBytes(substr)
	if len(sub) == 0 {
		return str
	}
	b := GoBytes(str)
	if len(b) == 0 {
		return nil
	}
	i := bytes.Index(b, sub)
	if i < 0 {
		return nil
	}
	return &b[i]
}

func StrCmp(a, b *byte) int {
	s1 := GoString(a)
	s2 := GoString(b)
	return strings.Compare(s1, s2)
}

func StrNCmp(a, b *byte, sz int) int {
	s1 := GoString(a)
	s2 := GoString(b)
	if len(s1) > sz {
		s1 = s1[:sz]
	}
	if len(s2) > sz {
		s2 = s2[:sz]
	}
	return strings.Compare(s1, s2)
}

func StrCaseCmp(a, b *byte) int {
	s1 := strings.ToLower(GoString(a))
	s2 := strings.ToLower(GoString(b))
	return strings.Compare(s1, s2)
}

func StrNCaseCmp(a, b *byte, sz int) int {
	s1 := strings.ToLower(GoString(a))
	s2 := strings.ToLower(GoString(b))
	if len(s1) > sz {
		s1 = s1[:sz]
	}
	if len(s2) > sz {
		s2 = s2[:sz]
	}
	return strings.Compare(s1, s2)
}

// StrCpyGo copies a Go slice into a C string pointed by dst. It won't add the null terminator.
func StrCpyGo(dst *byte, src []byte) {
	d := BytesN(dst, len(src))
	copy(d, src)
}

// StrCpyGoZero is the same as StrCpyGo, but adds a null terminator.
func StrCpyGoZero(dst *byte, src []byte) {
	d := BytesN(dst, len(src)+1)
	n := copy(d, src)
	d[n] = 0
}

// StrCpy copies the C string pointed by source into the array pointed by destination, including the terminating null character
// (and stopping at that point).
//
// To avoid overflows, the size of the array pointed by destination shall be long enough to contain the same C string as source
// (including the terminating null character), and should not overlap in memory with source.
func StrCpy(dst, src *byte) *byte {
	s := GoBytes(src)
	StrCpyGoZero(dst, s)
	return dst
}

// StrNCpy copies the first num characters of source to destination. If the end of the source C string
// (which is signaled by a null-character) is found before num characters have been copied, destination is padded with zeros
// until a total of num characters have been written to it.
//
// No null-character is implicitly appended at the end of destination if source is longer than num.
// Thus, in this case, destination shall not be considered a null terminated C string (reading it as such would overflow).
//
// Destination and source shall not overlap (see MemMove for a safer alternative when overlapping).
func StrNCpy(dst, src *byte, sz int) *byte {
	d := BytesN(dst, sz)
	if len(d) == 0 {
		return dst
	}
	s := GoBytes(src)
	pad := 0
	if len(s) > sz {
		s = s[:sz]
	} else if len(s) < sz {
		pad = sz - len(s)
	}
	n := copy(d, s)
	for i := 0; i < pad; i++ {
		d[n+i] = 0
	}
	return &d[0]
}

// StrCat appends a copy of the source string to the destination string. The terminating null character in destination
// is overwritten by the first character of source, and a null-character is included at the end of the new string
// formed by the concatenation of both in destination.
//
// Destination and source shall not overlap.
func StrCat(dst, src *byte) *byte {
	s := GoBytes(src)
	i := StrLen(dst)
	n := i + len(s)
	d := BytesN(dst, n+1)
	copy(d[i:], s)
	d[n] = 0
	return &d[0]
}

// StrNCat appends the first num characters of source to destination, plus a terminating null-character.
//
// If the length of the C string in source is less than num, only the content up to the terminating null-character is copied.
func StrNCat(dst, src *byte, sz int) *byte {
	s := GoBytes(src)
	if len(s) > sz {
		s = s[:sz]
	}
	n := StrLen(dst)
	d := BytesN(dst, n+len(s)+1)[n:]
	n = copy(d, s)
	d[n] = 0
	return dst
}

var strtok struct {
	data []byte
	ind  int
}

// StrTok is used in a sequence of calls to split str into tokens, which are sequences of contiguous characters
// separated by any of the characters that are part of delimiters.
//
// On a first call, the function expects a C string as argument for str, whose first character is used as the starting
// location to scan for tokens. In subsequent calls, the function expects a null pointer and uses the position
// right after the end of the last token as the new starting location for scanning.
//
// To determine the beginning and the end of a token, the function first scans from the starting location for the first
// character not contained in delimiters (which becomes the beginning of the token). And then scans starting from this
// beginning of the token for the first character contained in delimiters, which becomes the end of the token.
// The scan also stops if the terminating null character is found.
//
// This end of the token is automatically replaced by a null-character, and the beginning of the token is returned by the function.
//
// Once the terminating null character of str is found in a call to strtok, all subsequent calls to this function
// (with a null pointer as the first argument) return a null pointer.
//
// The point where the last token was found is kept internally by the function to be used on the next call
// (particular library implementations are not required to avoid data races).
func StrTok(src, delim *byte) *byte {
	if src != nil {
		strtok.data = GoBytes(src)
		strtok.ind = 0
	}
	d := GoString(delim)
	for ; strtok.ind < len(strtok.data); strtok.ind++ {
		if strings.IndexByte(d, strtok.data[strtok.ind]) < 0 {
			// start of a new token
			tok := strtok.data[strtok.ind:]
			if i := bytes.IndexAny(tok, d); i >= 0 {
				tok[i] = 0
				strtok.ind += i + 1
			} else {
				strtok.data = nil
				strtok.ind = 0
			}
			return &tok[0]
		}
		// skip delimiters
	}
	strtok.data = nil
	strtok.ind = 0
	return nil
}

// StrSpn returns the length of the initial portion of str1 which consists only of characters that are part of str2.
//
// The search does not include the terminating null-characters of either strings, but ends there.
func StrSpn(str, chars *byte) int {
	s := GoBytes(str)
	c := GoBytes(chars)
	i := 0
	for ; i < len(s); i++ {
		if bytes.IndexByte(c, s[i]) < 0 {
			break
		}
	}
	return i
}

func StrCSpn(a, b *byte) int {
	panic("TODO")
}

// StrDup copies a null-terminated C string.
func StrDup(s *byte) *byte {
	return CString(GoString(s))
}
