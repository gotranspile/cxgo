// TODO: split the file as per the spec
#include <stdint.h>

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
