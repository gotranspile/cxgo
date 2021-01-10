package cnet

import (
	"log"
	"net"

	"github.com/gotranspile/cxgo/runtime/libc"
)

type HostEnt struct {
	Name     *byte
	Aliases  **byte
	AddrType int32
	Length   int32
	AddrList **byte
}

func GetHostByName(name *byte) *HostEnt {
	host := libc.GoString(name)
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Printf("gethostbyname(%q): %v", host, err)
		libc.SetErr(err)
		return nil
	}
	arr := make([]*byte, len(ips)+1)
	for i, ip := range ips {
		ip4 := ip.To4()
		arr[i] = &ip4[0]
	}
	return &HostEnt{
		Name:     name,
		Aliases:  nil,
		AddrType: AF_INET,
		Length:   4,
		AddrList: &arr[0],
	}
}
