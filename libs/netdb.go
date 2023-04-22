package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/types"
)

const (
	netdbH = "netdb.h"
)

//go:embed netdb.h
var hnetdb string

func init() {
	RegisterLibrary(netdbH, func(c *Env) *Library {
		strT := c.C().String()
		intT := types.IntT(4)
		hostentT := types.NamedTGo("hostent", "cnet.HostEnt", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("h_name", "Name", strT)},
			{Name: types.NewIdentGo("h_aliases", "Aliases", c.PtrT(strT))},
			{Name: types.NewIdentGo("h_addrtype", "AddrType", intT)},
			{Name: types.NewIdentGo("h_length", "Length", intT)},
			{Name: types.NewIdentGo("h_addr_list", "AddrList", c.PtrT(strT))},
		}))
		return &Library{
			Imports: map[string]string{
				"cnet": RuntimePrefix + "cnet",
			},
			Types: map[string]types.Type{
				"hostent": hostentT,
			},
			Idents: map[string]*types.Ident{
				"gethostbyname": c.NewIdent("gethostbyname", "cnet.GetHostByName", cnet.GetHostByName, c.FuncTT(c.PtrT(hostentT), strT)),
			},
			// TODO
			Header: hnetdb,
		}
	})
}
