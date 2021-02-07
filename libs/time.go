package libs

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

const (
	timeH = "time.h"
)

func init() {
	RegisterLibrary(timeH, func(c *Env) *Library {
		gintT := c.Go().Int()
		intT := types.IntT(4)
		longT := types.IntT(8)
		strT := c.C().String()
		timeT := types.NamedTGo("time_t", "libc.Time", types.IntT(4))
		timerT := types.NamedTGo("timer_t", "libc.Timer", types.UintT(4))
		clockT := types.NamedTGo("clock_t", "libc.Clock", types.IntT(4))
		clockIDT := types.NamedTGo("clockid_t", "libc.ClockID", types.UintT(4))
		tvT := types.NamedTGo("timeval", "libc.TimeVal", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("tv_sec", "Sec", timeT)},
			{Name: types.NewIdentGo("tv_usec", "USec", longT)},
		}))
		tsT := types.NamedTGo("timespec", "libc.TimeSpec", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("tv_sec", "Sec", timeT)},
			{Name: types.NewIdentGo("tv_nsec", "NSec", longT)},
		}))
		tmT := types.NamedTGo("tm", "libc.TimeInfo", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("tm_sec", "Sec", intT)},
			{Name: types.NewIdentGo("tm_min", "Min", intT)},
			{Name: types.NewIdentGo("tm_hour", "Hour", intT)},
			{Name: types.NewIdentGo("tm_mday", "Day", intT)},
			{Name: types.NewIdentGo("tm_mon", "Month", intT)},
			{Name: types.NewIdentGo("tm_year", "Year", intT)},
			{Name: types.NewIdentGo("tm_wday", "WeekDay", intT)},
			{Name: types.NewIdentGo("tm_yday", "YearDay", intT)},
			{Name: types.NewIdentGo("tm_isdst", "IsDst", intT)},
			{Name: types.NewIdentGo("tm_zone", "Timezone", c.PtrT(nil))},
			{Name: types.NewIdentGo("tm_gmtoff", "GMTOffs", intT)},
		}))
		return &Library{
			Imports: map[string]string{
				"libc": RuntimeLibc,
			},
			Types: map[string]types.Type{
				"time_t":    timeT,
				"timer_t":   timerT,
				"clock_t":   clockT,
				"clockid_t": clockIDT,
				"tm":        tmT,
				"timeval":   tvT,
				"timespec":  tsT,
			},
			Idents: map[string]*types.Ident{
				"time":           c.NewIdent("time", "libc.GetTime", libc.GetTime, c.FuncTT(timeT, c.PtrT(timeT))),
				"mktime":         c.NewIdent("mktime", "libc.MakeTime", libc.MakeTime, c.FuncTT(timeT, c.PtrT(tmT))),
				"localtime":      c.NewIdent("localtime", "libc.LocalTime", libc.LocalTime, c.FuncTT(c.PtrT(tmT), c.PtrT(timeT))),
				"clock":          c.NewIdent("clock", "libc.ClockTicks", libc.ClockTicks, c.FuncTT(clockT)),
				"clock_getres":   c.NewIdent("clock_getres", "libc.ClockGetRes", libc.ClockGetRes, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"clock_settime":  c.NewIdent("clock_settime", "libc.ClockSetTime", libc.ClockSetTime, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"clock_gettime":  c.NewIdent("clock_gettime", "libc.ClockGetTime", libc.ClockGetTime, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"asctime":        c.NewIdent("asctime", "libc.AscTime", libc.AscTime, c.FuncTT(strT, c.PtrT(tmT))),
				"CLOCK_REALTIME": c.NewIdent("CLOCK_REALTIME", "libc.CLOCK_REALTIME", libc.CLOCK_REALTIME, gintT),
				"CLOCKS_PER_SEC": c.NewIdent("CLOCKS_PER_SEC", "libc.CLOCKS_PER_SEC", libc.CLOCKS_PER_SEC, gintT),
			},
			// TODO
			Header: `
#include <` + BuiltinH + `>
#include <` + sysTypesH + `>

const _cxgo_go_int CLOCK_REALTIME = 1;
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


`,
		}
	})
}
