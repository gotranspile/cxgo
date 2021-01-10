package libs

import (
	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/types"
)

const (
	netdbH = "netdb.h"
)

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
			Header: `
#include <` + stdintH + `>

struct hostent {
	char  *h_name;      // Official name of the host. 
	char **h_aliases;   // A pointer to an array of pointers to alternative host names, terminated by a null pointer. 
	int32_t h_addrtype;  // Address type. 
	int32_t h_length;    // The length, in bytes, of the address. 
	char **h_addr_list; // A pointer to an array of pointers to network addresses (in network byte order) for the host, terminated by a null pointer.
};

struct netent {
	char     *n_name;     // Official, fully-qualified (including the domain) name of the host. 
	char    **n_aliases;  // A pointer to an array of pointers to alternative network names, terminated by a null pointer. 
	int       n_addrtype; // The address type of the network. 
	uint32_t  n_net;      // The network number, in host byte order.
};

struct protoent {
	char   *p_name;     // Official name of the protocol. 
	char  **p_aliases;  // A pointer to an array of pointers to alternative protocol names, terminated by a null pointer. 
	int     p_proto;    // The protocol number.
};

struct servent {
	char   *s_name;     // Official name of the service. 
	char  **s_aliases;  // A pointer to an array of pointers to alternative service names, terminated by a null pointer. 
	int     s_port;     // The port number at which the service resides, in network byte order. 
	char   *s_proto;    // The name of the protocol to use when contacting the service.
};

void              endhostent(void);
void              endnetent(void);
void              endprotoent(void);
void              endservent(void);
void              freeaddrinfo(struct addrinfo *);
const char       *gai_strerror(int);
int               getaddrinfo(const char *restrict, const char *restrict,
                      const struct addrinfo *restrict,
                      struct addrinfo **restrict);
struct hostent   *gethostbyaddr(const void *, socklen_t, int);
struct hostent   *gethostbyname(const char *);
struct hostent   *gethostent(void);
int               getnameinfo(const struct sockaddr *restrict, socklen_t,
                      char *restrict, socklen_t, char *restrict,
                      socklen_t, int);
struct netent    *getnetbyaddr(uint32_t, int);
struct netent    *getnetbyname(const char *);
struct netent    *getnetent(void);
struct protoent  *getprotobyname(const char *);
struct protoent  *getprotobynumber(int);
struct protoent  *getprotoent(void);
struct servent   *getservbyname(const char *, const char *);
struct servent   *getservbyport(int, const char *);
struct servent   *getservent(void);
void              sethostent(int);
void              setnetent(int);
void              setprotoent(int);
void              setservent(int);
`,
		}
	})
}
