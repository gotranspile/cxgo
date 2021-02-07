package libc

import (
	"time"
	"unsafe"
)

const CLOCK_REALTIME = 1

const CLOCKS_PER_SEC = 1000000 // us

type Time int32

func (t Time) GoTime() time.Time {
	return time.Unix(int64(t), 0)
}

type Timer uint32
type Clock int32
type ClockID uint32

var clockStart = time.Now()

func ClockTicks() Clock {
	return Clock(time.Since(clockStart).Microseconds())
}

func GetTime(dst *Time) Time {
	t := Time(time.Now().Unix())
	if dst != nil {
		*dst = t
	}
	return t
}

func MakeTime(src *TimeInfo) Time {
	if src == nil {
		return Time(time.Now().Unix())
	}
	return Time(time.Date(1900+int(src.Year), time.Month(src.Month), int(src.Day), int(src.Hour), int(src.Min), int(src.Sec), 0, time.Local).Unix())
}

func LocalTime(src *Time) *TimeInfo {
	t := src.GoTime()
	return &TimeInfo{
		Sec:     int32(t.Second()),
		Min:     int32(t.Minute()),
		Hour:    int32(t.Hour()),
		Day:     int32(t.Day()),
		Month:   int32(t.Month()),
		Year:    int32(t.Year()) - 1900,
		WeekDay: int32(t.Weekday()),
		YearDay: int32(t.YearDay()),
		// TODO
		IsDst:    0,
		Timezone: nil,
		GMTOffs:  0,
	}
}

func AscTime(tm *TimeInfo) *byte {
	t := tm.GoTime()
	s := t.Format("Mon Jan _2 15:04:05 2006")
	return CString(s)
}

type TimeVal struct {
	Sec  Time
	USec int64
}

type TimeSpec struct {
	Sec  Time
	NSec int64
}

type TimeInfo struct {
	Sec      int32
	Min      int32
	Hour     int32
	Day      int32
	Month    int32
	Year     int32
	WeekDay  int32
	YearDay  int32
	IsDst    int32
	Timezone unsafe.Pointer
	GMTOffs  int32
}

func (tm *TimeInfo) GoTime() time.Time {
	return MakeTime(tm).GoTime()
}

func ClockGetRes(c Clock, ts *TimeSpec) int32 {
	panic("TODO")
}

func ClockSetTime(c Clock, ts *TimeSpec) int32 {
	panic("TODO")
}

func ClockGetTime(c Clock, ts *TimeSpec) int32 {
	panic("TODO")
}
