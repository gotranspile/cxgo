package cnet

import (
	"net"
	"testing"

	"github.com/dennwc/cxgo/runtime/libc"
	"github.com/stretchr/testify/require"
)

func TestNtoa(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	b := Ntoa(Address{Addr: ipToInt4(ip)})
	s := libc.GoString(b)
	require.Equal(t, "127.0.0.1", s)
}
