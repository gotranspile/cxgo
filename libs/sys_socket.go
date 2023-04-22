package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysSocketH = "sys/socket.h"
)

//go:embed sys_socket.h
var hsys_socket string

func init() {
	RegisterLibrary(sysSocketH, func(c *Env) *Library {
		intT := types.IntT(4)
		int16T := types.IntT(2)
		uint16T := types.UintT(2)
		gintT := c.Go().Int()
		uintT := types.UintT(4)
		bytesT := c.C().String()
		inAddrT := c.GetLibraryType(arpaInetH, "in_addr")
		sockAddrT := types.NamedT("cnet.SockAddr", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("sa_family", "Family", int16T)},
			{Name: types.NewIdentGo("sa_data", "Data", c.C().BytesN(14))},
		}))
		sockAddrInT := types.NamedT("cnet.SockAddrInet", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("sin_family", "Family", int16T)},
			{Name: types.NewIdentGo("sin_port", "Port", uint16T)},
			{Name: types.NewIdentGo("sin_addr", "Addr", inAddrT)},
			{Name: types.NewIdentGo("sin_zero", "Zero", c.C().BytesN(8))},
		}))
		return &Library{
			Imports: map[string]string{
				"csys": RuntimePrefix + "csys",
				"cnet": RuntimePrefix + "cnet",
			},
			Types: map[string]types.Type{
				"sockaddr":    sockAddrT,
				"sockaddr_in": sockAddrInT,
			},
			Idents: map[string]*types.Ident{
				"AF_INET":      c.NewIdent("AF_INET", "cnet.AF_INET", cnet.AF_INET, gintT),
				"SOCK_STREAM":  c.NewIdent("SOCK_STREAM", "cnet.SOCK_STREAM", cnet.SOCK_STREAM, gintT),
				"SOCK_DGRAM":   c.NewIdent("SOCK_DGRAM", "cnet.SOCK_DGRAM", cnet.SOCK_DGRAM, gintT),
				"SOL_SOCKET":   c.NewIdent("SOL_SOCKET", "cnet.SOL_SOCKET", cnet.SOL_SOCKET, types.UintT(2)),
				"SO_BROADCAST": c.NewIdent("SO_BROADCAST", "cnet.SO_BROADCAST", cnet.SO_BROADCAST, gintT),
				"accept":       c.NewIdent("accept", "cnet.Accept", cnet.Accept, c.FuncTT(gintT, gintT, c.PtrT(sockAddrT), c.PtrT(gintT))),
				"bind":         c.NewIdent("bind", "cnet.Bind", cnet.Bind, c.FuncTT(gintT, gintT, c.PtrT(sockAddrT), gintT)),
				"listen":       c.NewIdent("listen", "cnet.Listen", cnet.Listen, c.FuncTT(intT, intT, intT)),
				"shutdown":     c.NewIdent("shutdown", "cnet.Shutdown", cnet.Shutdown, c.FuncTT(intT, intT, intT)),
				"send":         c.NewIdent("send", "cnet.Send", cnet.Send, c.FuncTT(intT, intT, bytesT, intT, intT)),
				"sendto":       c.NewIdent("sendto", "cnet.SendTo", cnet.SendTo, c.FuncTT(intT, intT, bytesT, uintT, intT, c.PtrT(sockAddrT), gintT)),
				"recv":         c.NewIdent("recv", "cnet.Recv", cnet.Recv, c.FuncTT(intT, intT, bytesT, uintT, intT)),
				"recvfrom":     c.NewIdent("recvfrom", "cnet.RecvFrom", cnet.RecvFrom, c.FuncTT(intT, intT, bytesT, uintT, intT, c.PtrT(sockAddrT), c.PtrT(gintT))),
				"socket":       c.NewIdent("socket", "cnet.Socket", cnet.Socket, c.FuncTT(gintT, gintT, gintT, gintT)),
				"setsockopt":   c.NewIdent("setsockopt", "cnet.SetSockOpt", cnet.SetSockOpt, c.FuncTT(gintT, gintT, gintT, gintT, bytesT, gintT)),
			},
			// TODO
			Header: hsys_socket,
		}
	})
}
