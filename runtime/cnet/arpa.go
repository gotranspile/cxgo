package cnet

import (
	"encoding/binary"
	"net"
	"unsafe"
)

var (
	netOrder = binary.BigEndian
)

func putHost16(b []byte, v uint16) {
	*(*uint16)(unsafe.Pointer(&b[0])) = v
}

func putHost32(b []byte, v uint32) {
	*(*uint32)(unsafe.Pointer(&b[0])) = v
}

func getHost16(b []byte) uint16 {
	return *(*uint16)(unsafe.Pointer(&b[0]))
}

func getHost32(b []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&b[0]))
}

type Addr uint64
type Port uint32
type Address struct {
	Addr uint32
}

func ParseAddr(s string) Addr {
	panic("TODO")
}

func Htonl(v uint32) uint32 {
	var b [4]byte
	netOrder.PutUint32(b[:], v)
	return getHost32(b[:])
}

func Htons(v uint16) uint16 {
	var b [2]byte
	netOrder.PutUint16(b[:], v)
	return getHost16(b[:])
}

func Ntohl(v uint32) uint32 {
	var b [4]byte
	putHost32(b[:], v)
	return netOrder.Uint32(b[:])
}

func Ntohs(v uint16) uint16 {
	var b [2]byte
	putHost16(b[:], v)
	return netOrder.Uint16(b[:])
}

func ipToInt4(ip net.IP) uint32 {
	ip4 := ip.To4()
	v := *(*uint32)(unsafe.Pointer(&ip4[0]))
	return v
}

func intToIp4(v uint32) net.IP {
	p := *(*[4]byte)(unsafe.Pointer(&v))
	return net.IP(p[:])
}

var ntoaBuf [16]byte

func Ntoa(a Address) *byte {
	ip := intToIp4(a.Addr)
	copy(ntoaBuf[:], ip.String())
	return &ntoaBuf[0]
}

func Ntop(a1 int32, a2 unsafe.Pointer, a3 *byte, a4 uint32) *byte {
	panic("TODO")
}

func Pton(a1 int32, a2 *byte, a3 unsafe.Pointer) *byte {
	panic("TODO")
}
