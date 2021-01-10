package sdl

import (
	"sync"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func GetTicks() uint32 {
	return sdl.GetTicks()
}

func Delay(ms uint32) {
	sdl.Delay(ms)
}

type TimerID int32

type TimerFunc = func(interval uint32, arg unsafe.Pointer) uint32

type timer struct {
	cancel chan struct{}
	fnc    TimerFunc
	arg    unsafe.Pointer
}

func (t *timer) start(interval uint32) {
	dt := time.Duration(interval) * time.Millisecond
	tm := time.NewTimer(dt)
	go func() {
		defer tm.Stop()
		for {
			select {
			case <-t.cancel:
				return
			case <-tm.C:
				interval = t.fnc(interval, t.arg)
				if interval == 0 {
					return
				}
				dt = time.Duration(interval) * time.Millisecond
				tm.Reset(dt)
			}
		}
	}()
}

var timers struct {
	sync.Mutex
	last TimerID
	byID map[TimerID]*timer
}

func AddTimer(interval uint32, fnc TimerFunc, arg unsafe.Pointer) TimerID {
	t := &timer{
		cancel: make(chan struct{}),
		fnc:    fnc,
		arg:    arg,
	}
	timers.Lock()
	timers.last++
	id := timers.last
	if timers.byID == nil {
		timers.byID = make(map[TimerID]*timer)
	}
	timers.byID[id] = t
	timers.Unlock()
	t.start(interval)
	return id
}

func RemoveTimer(timer TimerID) bool {
	timers.Lock()
	t, ok := timers.byID[timer]
	if ok {
		delete(timers.byID, timer)
		close(t.cancel)
	}
	timers.Unlock()
	return ok
}
