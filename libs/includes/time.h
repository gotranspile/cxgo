#include <sys/types.h>

const _cxgo_go_int CLOCK_REALTIME = 0;
const _cxgo_go_int CLOCK_MONOTONIC = 1;
const _cxgo_go_int CLOCKS_PER_SEC = 1000000;

typedef _cxgo_int32 time_t;
typedef _cxgo_uint32 timer_t;
typedef _cxgo_int32 clock_t;
typedef _cxgo_uint32 clockid_t;
#define suseconds_t _cxgo_int64

struct tm {
    _cxgo_int32    tm_sec;
    _cxgo_int32    tm_min;
    _cxgo_int32    tm_hour;
    _cxgo_int32    tm_mday;
    _cxgo_int32    tm_mon;
    _cxgo_int32    tm_year;
    _cxgo_int32    tm_wday;
    _cxgo_int32    tm_yday;
    _cxgo_int32    tm_isdst;
    void*  tm_zone;
    _cxgo_int32    tm_gmtoff;
};
struct timeval {
    time_t         tv_sec;
    _cxgo_int64    tv_usec;
};
struct timespec {
    time_t  tv_sec;
    _cxgo_int64    tv_nsec;
};

char      *asctime(const struct tm *);
char      *asctime_r(const struct tm *, char *);
clock_t    clock(void);
_cxgo_int32        clock_getres(clockid_t, struct timespec *);
_cxgo_int32        clock_gettime(clockid_t, struct timespec *);
_cxgo_int32        clock_settime(clockid_t, const struct timespec *);
char      *ctime(const time_t *);
char      *ctime_r(const time_t *, char *);
double     difftime(time_t, time_t);
struct tm *getdate(const char *);
struct tm *gmtime(const time_t *);
struct tm *gmtime_r(const time_t *, struct tm *);
struct tm *localtime(const time_t *);
struct tm *localtime_r(const time_t *, struct tm *);
time_t     mktime(struct tm *);
int        nanosleep(const struct timespec *, struct timespec *);
size_t     strftime(char *, size_t, const char *, const struct tm *);
char      *strptime(const char *, const char *, struct tm *);
time_t     time(time_t *);
int        timer_create(clockid_t, struct sigevent *, timer_t *);
int        timer_delete(timer_t);
int        timer_gettime(timer_t, struct itimerspec *);
int        timer_getoverrun(timer_t);
int        timer_settime(timer_t, int, const struct itimerspec *, struct itimerspec *);
void       tzset(void);
