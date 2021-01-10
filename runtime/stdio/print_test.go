package stdio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprintf(t *testing.T) {
	a := []byte("0000000000\x00")
	n := Sprintf(&a[0], "hi %lu\n", int32(10))
	require.Equal(t, 6, int(n))
	require.Equal(t, "hi 10\n", string(a[:n]))
	require.Equal(t, "hi 10\n\x00000\x00", string(a))
}
