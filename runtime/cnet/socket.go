package cnet

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

const (
	AF_INET      = 2
	SOCK_STREAM  = 1
	SOCK_DGRAM   = 2
	SOL_SOCKET   = uint16(0xffff)
	SO_BROADCAST = 0x20
)

type SockAddr struct {
	Family int16
	Data   [14]byte
}

type SockAddrInet struct {
	Family int16
	Port   uint16
	Addr   Address
	Zero   [8]byte
}

func Accept(a1 int, a2 *SockAddr, a3 *int) int {
	panic("TODO")
}

func Bind(fd int, addr *SockAddr, sz int) int {
	var sa syscall.Sockaddr
	switch addr.Family {
	case AF_INET:
		addr := (*SockAddrInet)(unsafe.Pointer(addr))
		port := Ntohs(addr.Port)
		if addr.Addr.Addr != 0 {
			panic("TODO")
		}
		sa = &syscall.SockaddrInet4{Port: int(port)}
	default:
		err := fmt.Errorf("unsupported address family: %d", addr.Family)
		log.Printf("bind(%d, %p, %d): %v", fd, addr, sz, err)
		libc.SetErr(err)
		return -1
	}
	err := syscall.Bind(Handle(fd), sa)
	if err != nil {
		log.Printf("bind(%d, %p, %d): %v", fd, addr, sz, err)
		libc.SetErr(err)
		return -1
	}
	return 0
}

func Socket(domain, typ, proto int) int {
	fd, err := syscall.Socket(domain, typ, proto)
	if err != nil {
		log.Printf("socket(%d, %d, %d): %v", domain, typ, proto, err)
		libc.SetErr(err)
		return -1
	}
	return int(fd)
}

func Listen(a1, a2 int32) int32 {
	panic("TODO")
}

func Shutdown(a1, a2 int32) int32 {
	panic("TODO")
}

func Send(a1 int32, a2 *byte, a3 int32, a4 int32) int32 {
	panic("TODO")
}

func SendTo(a1 int32, a2 *byte, a3 uint32, a4 int32, a5 *SockAddr, a6 int) int32 {
	panic("TODO")
}

func Recv(a1 int32, a2 *byte, a3 uint32, a4 int32) int32 {
	panic("TODO")
}

func RecvFrom(a1 int32, a2 *byte, a3 uint32, a4 int32, a5 *SockAddr, a6 *int) int32 {
	panic("TODO")
}

func SetSockOpt(fd int, level int, name int, val *byte, sz int) int {
	var err error
	switch name {
	case SO_BROADCAST:
		if sz != 4 {
			panic(sz)
		}
		val := *(*int32)(unsafe.Pointer(val))
		err = syscall.SetsockoptInt(Handle(fd), level, name, int(val))
	default:
		log.Printf("setsockopt(%d, %d, %d, %x, %d)", fd, level, name, val, sz)
		panic("TODO")
	}
	if err != nil {
		log.Printf("setsockopt(%d, %d, %d, %x, %d): %v", fd, level, name, val, sz, err)
		libc.SetErr(err)
		return -1
	}
	return 0
}
