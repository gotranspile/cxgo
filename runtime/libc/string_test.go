package libc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrLen(t *testing.T) {
	n := int(StrLen(nil))
	require.Equal(t, 0, n)

	b := make([]byte, 10)

	n = int(StrLen(&b[0]))
	require.Equal(t, 0, n)

	b[0] = 'a'
	n = int(StrLen(&b[0]))
	require.Equal(t, 1, n)

	b[1] = 'b'
	b[2] = 'c'
	n = int(StrLen(&b[0]))
	require.Equal(t, 3, n)

	b[0] = 0
	n = int(StrLen(&b[0]))
	require.Equal(t, 0, n)
	n = int(StrLen(&b[1]))
	require.Equal(t, 2, n)

	n = int(StrLen(CString("\x00")))
	require.Equal(t, 0, n)

	n = int(StrLen(CString("abc")))
	require.Equal(t, 3, n)
}

func TestStrChr(t *testing.T) {
	b := []byte("abcded\x00")

	p := StrChr(&b[0], 'd')
	require.Equal(t, &b[3], p)

	p = StrChr(&b[0], 'f')
	require.Equal(t, (*byte)(nil), p)

	b[3] = 0
	p = StrChr(&b[0], 'e')
	require.Equal(t, (*byte)(nil), p)

	p = StrChr(nil, 'e')
	require.Equal(t, (*byte)(nil), p)
}

func TestStrCpy(t *testing.T) {
	a := []byte("0000000000\x00")
	b := []byte("abcded\x00")

	// Copies the C string pointed by source into the array pointed by destination,
	// including the terminating null character (and stopping at that point).
	p := StrCpy(&a[0], &b[0])
	require.Equal(t, &a[0], p)
	require.Equal(t, "abcded\x00000\x00", string(a))
}

func TestStrNCpy(t *testing.T) {
	a := []byte("0000000000\x00")
	b := []byte("abcded\x00")

	// No null-character is implicitly appended at the end of destination if source is longer than num.
	p := StrNCpy(&a[0], &b[0], 3)
	require.Equal(t, &a[0], p)
	require.Equal(t, "abc0000000\x00", string(a))

	// If the end of the source C string (which is signaled by a null-character) is found before num characters have been copied,
	// destination is padded with zeros until a total of num characters have been written to it.
	p = StrNCpy(&a[0], &b[0], len(b)+2)
	require.Equal(t, &a[0], p)
	require.Equal(t, "abcded\x00\x00\x000\x00", string(a))
}

func TestStrCat(t *testing.T) {
	a := []byte("1\x00000000000\x00")
	b := []byte("abcded\x00")

	p := StrCat(&a[0], &b[0])
	require.Equal(t, &a[0], p)
	require.Equal(t, "1abcded\x00000\x00", string(a))
}

func TestStrNCat(t *testing.T) {
	a := []byte("0\x00000000000\x00")
	b := []byte("abcded\x00")

	p := StrNCat(&a[0], &b[0], 3)
	require.Equal(t, &a[0], p)
	require.Equal(t, "0abc\x00000000\x00", string(a))

	a = []byte("0\x00000000000\x00")
	p = StrNCat(&a[0], &b[0], 10)
	require.Equal(t, &a[0], p)
	require.Equal(t, "0abcded\x00000\x00", string(a))
}

func TestStrTok(t *testing.T) {
	a := []byte("- This, a sample string.\x00")
	toks := []byte(" ,.-\x00")

	var (
		lines []string
		ptrs  []*byte
	)
	for p := StrTok(&a[0], &toks[0]); p != nil; p = StrTok(nil, &toks[0]) {
		lines = append(lines, GoString(p))
		ptrs = append(ptrs, p)
	}
	// should split words without delimiters
	require.Equal(t, []string{
		"This",
		"a",
		"sample",
		"string",
	}, lines)
	// should point to the same string
	require.Equal(t, []*byte{
		&a[2],
		&a[8],
		&a[10],
		&a[17],
	}, ptrs)
	// ...which means it should insert zero bytes into it
	require.Equal(t, "- This\x00 a\x00sample\x00string\x00\x00", string(a))
}

func TestStrSpn(t *testing.T) {
	a := []byte("123abc\x00")
	toks := []byte("1234567890\x00")
	i := StrSpn(&a[0], &toks[0])
	require.Equal(t, 3, int(i))
}

func TestStrStr(t *testing.T) {
	a := []byte("123abc\x00")
	b := []byte("3a\x00")
	p := StrStr(&a[0], &b[0])
	require.Equal(t, &a[2], p)
}

func TestStrDup(t *testing.T) {
	a := []byte("123abc\x00")
	p := StrDup(&a[0])
	require.True(t, &a[0] != p)
	b := BytesN(p, len(a))
	require.Equal(t, string(a), string(b))
}
