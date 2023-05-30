#include <sys/types.h>
#include <arpa/inet.h>

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
