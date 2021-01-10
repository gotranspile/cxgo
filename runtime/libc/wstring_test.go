package libc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCWString(t *testing.T) {
	const s = "abcабв"
	p := CWString(s)
	v := GoWString(p)
	require.Equal(t, s, v)
}

func TestWStrLen(t *testing.T) {
	n := int(WStrLen(nil))
	require.Equal(t, 0, n)

	b := make([]uint16, 10)

	n = int(WStrLen(&b[0]))
	require.Equal(t, 0, n)

	b[0] = 'a'
	n = int(WStrLen(&b[0]))
	require.Equal(t, 1, n)

	b[1] = 'b'
	b[2] = 'c'
	n = int(WStrLen(&b[0]))
	require.Equal(t, 3, n)

	b[0] = 0
	n = int(WStrLen(&b[0]))
	require.Equal(t, 0, n)
	n = int(WStrLen(&b[1]))
	require.Equal(t, 2, n)

	n = int(WStrLen(CWString("\x00")))
	require.Equal(t, 0, n)

	n = int(WStrLen(CWString("abc")))
	require.Equal(t, 3, n)
}

/*
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
*/

func wrepeat(c WChar, n int) []WChar {
	b := make([]WChar, n)
	for i := range b {
		b[i] = c
	}
	return b
}

func TestWStrNCpy(t *testing.T) {
	a := wrepeat('0', 10)
	b := []WChar{'a', 'b', 'c', 'd', 'e', 'd', 0}

	// No null-character is implicitly appended at the end of destination if source is longer than num.
	p := WStrNCpy(&a[0], &b[0], 3)
	require.Equal(t, &a[0], p)
	require.Equal(t, []WChar{'a', 'b', 'c', '0', '0', '0', '0', '0', '0', '0'}, a)

	// If the end of the source C string (which is signaled by a null-character) is found before num characters have been copied,
	// destination is padded with zeros until a total of num characters have been written to it.
	p = WStrNCpy(&a[0], &b[0], uint32(len(b)+1))
	require.Equal(t, &a[0], p)
	require.Equal(t, []WChar{'a', 'b', 'c', 'd', 'e', 'd', 0, 0, '0', '0'}, a)
}

func TestWStrCat(t *testing.T) {
	a := wrepeat('0', 10)
	a[0] = '1'
	a[1] = 0
	a[9] = 0
	b := []WChar{'a', 'b', 'c', 'd', 'e', 'd', 0}

	p := WStrCat(&a[0], &b[0])
	require.Equal(t, &a[0], p)
	require.Equal(t, []WChar{'1', 'a', 'b', 'c', 'd', 'e', 'd', 0, '0', 0}, a)
}

/*
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
*/
