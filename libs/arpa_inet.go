package libs

import (
	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/types"
)

// https://pubs.opengroup.org/onlinepubs/000095399/basedefs/arpa/inet.h.html

const (
	arpaInetH = "arpa/inet.h"
)

func init() {
	RegisterLibrary(arpaInetH, func(c *Env) *Library {
		inAddrT := types.NamedTGo("in_addr_t", "cnet.Addr", types.UintT(8))
		inPortT := types.NamedTGo("in_port_t", "cnet.Port", types.UintT(4))
		inAddr := types.NamedTGo("in_addr", "cnet.Address", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("s_addr", "Addr", types.UintT(4))},
		}))
		sockLenT := types.UintT(4)
		strT := c.C().String()
		gstrT := c.Go().String()
		return &Library{
			Imports: map[string]string{
				"cnet": RuntimePrefix + "cnet",
			},
			Types: map[string]types.Type{
				"in_addr":   inAddr,
				"in_addr_t": inAddrT,
				"in_port_t": inPortT,
			},
			Idents: map[string]*types.Ident{
				"inet_addr": c.NewIdent("inet_addr", "cnet.ParseAddr", cnet.ParseAddr, c.FuncTT(inAddrT, gstrT)),
				"htonl":     c.NewIdent("htonl", "cnet.Htonl", cnet.Htonl, c.FuncTT(types.UintT(4), types.UintT(4))),
				"htons":     c.NewIdent("htons", "cnet.Htons", cnet.Htons, c.FuncTT(types.UintT(2), types.UintT(2))),
				"ntohl":     c.NewIdent("ntohl", "cnet.Ntohl", cnet.Ntohl, c.FuncTT(types.UintT(4), types.UintT(4))),
				"ntohs":     c.NewIdent("ntohs", "cnet.Ntohs", cnet.Ntohs, c.FuncTT(types.UintT(2), types.UintT(2))),
				"inet_ntoa": c.NewIdent("inet_ntoa", "cnet.Ntoa", cnet.Ntoa, c.FuncTT(strT, inAddr)),
				"inet_ntop": c.NewIdent("inet_ntop", "cnet.Ntop", cnet.Ntop, c.FuncTT(strT, types.IntT(4), c.PtrT(nil), strT, sockLenT)),
				"inet_pton": c.NewIdent("inet_pton", "cnet.Pton", cnet.Pton, c.FuncTT(strT, types.IntT(4), strT, c.PtrT(nil))),
			},
			// TODO: split the file as per the spec
			Header: `
#include <` + stdintH + `>

uint32_t htonl(uint32_t);
uint16_t htons(uint16_t);
uint32_t ntohl(uint32_t);
uint16_t ntohs(uint16_t);

struct in_addr {
    uint32_t s_addr;
};

typedef uint64_t in_addr_t;
typedef uint32_t in_port_t;
#define socklen_t uint32_t

in_addr_t    inet_addr(const char *);
char        *inet_ntoa(struct in_addr);
const char  *inet_ntop(int32_t, const void *restrict, char *restrict, socklen_t);
int32_t      inet_pton(int32_t, const char *restrict, void *restrict);
`,
		}
	})
}
