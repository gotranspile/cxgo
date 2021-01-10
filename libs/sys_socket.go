package libs

import (
	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/types"
)

const (
	sysSocketH = "sys/socket.h"
)

func init() {
	RegisterLibrary(sysSocketH, func(c *Env) *Library {
		intT := types.IntT(4)
		int16T := types.IntT(2)
		uint16T := types.UintT(2)
		gintT := c.Go().Int()
		uintT := types.UintT(4)
		bytesT := c.C().String()
		inAddrT := c.GetLib(arpaInetH).GetType("in_addr")
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
			Header: `
#include <` + sysTypesH + `>
#include <` + arpaInetH + `>

const _cxgo_go_int AF_INET = 2;	/* Internet IP Protocol 	*/

const _cxgo_go_int SOCK_STREAM = 1;		/* stream (connection) socket	*/
const _cxgo_go_int SOCK_DGRAM = 2;		/* datagram (conn.less) socket	*/

const _cxgo_uint16 SOL_SOCKET = 0xffff;

const _cxgo_go_int SO_BROADCAST = 0x20;

#define socklen_t _cxgo_go_int

struct sockaddr {
    _cxgo_int16 sa_family;
    char sa_data[14];
};

// TODO: should be elsewhere
struct in_addr {
    _cxgo_uint32 s_addr;
};

struct sockaddr_in {
    _cxgo_int16    sin_family;   // e.g. AF_INET
    _cxgo_uint16   sin_port;     // e.g. htons(3490)
    struct in_addr   sin_addr;
    char             sin_zero[8];
};

_cxgo_go_int     accept(_cxgo_go_int, struct sockaddr *restrict, socklen_t *restrict);
_cxgo_go_int     bind(_cxgo_go_int, const struct sockaddr *, socklen_t);
int     connect(int, const struct sockaddr *, socklen_t);
int     getpeername(int, struct sockaddr *restrict, socklen_t *restrict);
int     getsockname(int, struct sockaddr *restrict, socklen_t *restrict);
int     getsockopt(int, int, int, void *restrict, socklen_t *restrict);
_cxgo_int32     listen(_cxgo_int32, _cxgo_int32);
_cxgo_int32 recv(_cxgo_int32, void *, _cxgo_uint32, _cxgo_int32);
_cxgo_int32 recvfrom(_cxgo_int32, void *restrict, _cxgo_uint32, _cxgo_int32, struct sockaddr *restrict, socklen_t *restrict);
ssize_t recvmsg(int, struct msghdr *, int);
_cxgo_int32 send(_cxgo_int32, const void *, _cxgo_int32, _cxgo_int32);
ssize_t sendmsg(int, const struct msghdr *, int);
_cxgo_int32 sendto(_cxgo_int32, const void *, _cxgo_uint32, _cxgo_int32, const struct sockaddr *, socklen_t);
_cxgo_int32     setsockopt(_cxgo_int32, _cxgo_int32, _cxgo_int32, const void *, socklen_t);
_cxgo_int32     shutdown(_cxgo_int32, _cxgo_int32);
_cxgo_go_int socket(_cxgo_go_int, _cxgo_go_int, _cxgo_go_int);
int     sockatmark(int);
int     socketpair(int, int, int, int[2]);
`,
		}
	})
}
