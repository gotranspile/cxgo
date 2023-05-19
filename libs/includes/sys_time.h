#include <time.h>
#include <sys/types.h>

typedef struct fd_set {
	long fds_bits[];
} fd_set;

_cxgo_int32   getitimer(_cxgo_int32, struct itimerval *);
_cxgo_int32   gettimeofday(struct timeval *restrict, void *restrict);
int   select(int, fd_set *restrict, fd_set *restrict, fd_set *restrict, struct timeval *restrict);
int   setitimer(int, const struct itimerval *restrict, struct itimerval *restrict);
