package libs

import (
	_ "embed"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/types"
)

const (
	timeH = "time.h"
)

//go:embed time.h
var htime string

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
				"time":            c.NewIdent("time", "libc.GetTime", libc.GetTime, c.FuncTT(timeT, c.PtrT(timeT))),
				"mktime":          c.NewIdent("mktime", "libc.MakeTime", libc.MakeTime, c.FuncTT(timeT, c.PtrT(tmT))),
				"localtime":       c.NewIdent("localtime", "libc.LocalTime", libc.LocalTime, c.FuncTT(c.PtrT(tmT), c.PtrT(timeT))),
				"clock":           c.NewIdent("clock", "libc.ClockTicks", libc.ClockTicks, c.FuncTT(clockT)),
				"clock_getres":    c.NewIdent("clock_getres", "libc.ClockGetRes", libc.ClockGetRes, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"clock_settime":   c.NewIdent("clock_settime", "libc.ClockSetTime", libc.ClockSetTime, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"clock_gettime":   c.NewIdent("clock_gettime", "libc.ClockGetTime", libc.ClockGetTime, c.FuncTT(intT, clockT, c.PtrT(tsT))),
				"asctime":         c.NewIdent("asctime", "libc.AscTime", libc.AscTime, c.FuncTT(strT, c.PtrT(tmT))),
				"CLOCK_REALTIME":  c.NewIdent("CLOCK_REALTIME", "libc.CLOCK_REALTIME", libc.CLOCK_REALTIME, gintT),
				"CLOCK_MONOTONIC": c.NewIdent("CLOCK_MONOTONIC", "libc.CLOCK_MONOTONIC", libc.CLOCK_MONOTONIC, gintT),
				"CLOCKS_PER_SEC":  c.NewIdent("CLOCKS_PER_SEC", "libc.CLOCKS_PER_SEC", libc.CLOCKS_PER_SEC, gintT),
			},
			// TODO
			Header: htime,
		}
	})
}
