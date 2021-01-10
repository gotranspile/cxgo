package sdl

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NUM_SCANCODES = sdl.NUM_SCANCODES

	PRESSED = sdl.PRESSED

	SCANCODE_SPACE = sdl.SCANCODE_SPACE

	KMOD_LSHIFT = sdl.KMOD_LSHIFT
	KMOD_RSHIFT = sdl.KMOD_RSHIFT
	KMOD_RALT   = sdl.KMOD_RALT
)

type FingerID int32

type Keymod int32

type Scancode int32

func PollEvent(e *Event) int32 {
	ev := sdl.PollEvent()
	if ev == nil {
		return 0
	}
	//log.Printf("sdl.PollEvent: %T", ev)
	switch ev := ev.(type) {
	case *sdl.WindowEvent:
		*e = Event{
			Type: int32(ev.Type),
			Window: WindowEvent{
				Event: ev.Event,
			},
		}
	case *sdl.TextEditingEvent:
		*e = Event{
			Type: int32(ev.Type),
		}
		copy(e.Edit.Text[:], ev.Text[:])
	case *sdl.TextInputEvent:
		*e = Event{
			Type: int32(ev.Type),
		}
		copy(e.Text.Text[:], ev.Text[:])
	case *sdl.KeyboardEvent:
		*e = Event{
			Type: int32(ev.Type),
			Key: KeyboardEvent{
				Keysym: Keysym{
					Scancode: Scancode(ev.Keysym.Scancode),
				},
				State: ev.State,
			},
		}
	case *sdl.MouseButtonEvent:
		*e = Event{
			Type: int32(ev.Type),
			Button: MouseButtonEvent{
				Button: ev.Button,
				State:  ev.State,
				X:      ev.X,
				Y:      ev.Y,
			},
		}
	case *sdl.MouseMotionEvent:
		*e = Event{
			Type: int32(ev.Type),
			Motion: MouseMotionEvent{
				X:    ev.X,
				Y:    ev.Y,
				Xrel: ev.XRel,
				Yrel: ev.YRel,
			},
		}
	case *sdl.MouseWheelEvent:
		*e = Event{
			Type: int32(ev.Type),
			Wheel: MouseWheelEvent{
				X: ev.X,
				Y: ev.Y,
			},
		}
	default:
		return 0 // TODO
	}
	return 1
}

type Event struct {
	Type   int32
	Edit   TextEditingEvent
	Text   TextInputEvent
	Key    KeyboardEvent
	Button MouseButtonEvent
	Motion MouseMotionEvent
	Wheel  MouseWheelEvent
	Window WindowEvent
}

type WindowEvent struct {
	Event uint8
}

type TextInputEvent struct {
	Text [32]byte
}

type TextEditingEvent struct {
	Text [32]byte
}

type MouseWheelEvent struct {
	X int32
	Y int32
}

type MouseMotionEvent struct {
	X    int32
	Y    int32
	Xrel int32
	Yrel int32
}

type MouseButtonEvent struct {
	Button uint8
	State  uint8
	X      int32
	Y      int32
}

type KeyboardEvent struct {
	Keysym Keysym
	State  uint8
}

type Keysym struct {
	Scancode Scancode
}

func StartTextInput() {
	sdl.StartTextInput()
}

func StopTextInput() {
	sdl.StopTextInput()
}

func GetModState() Keymod {
	return Keymod(sdl.GetModState())
}

func SetRelativeMouseMode(enabled bool) int32 {
	if noMouseGrab {
		return 0
	}
	return int32(sdl.SetRelativeMouseMode(enabled))
}

func GetEventState(ev uint32) uint8 {
	return sdl.GetEventState(ev)
}

func GetGlobalMouseState(x, y *int32) uint32 {
	cx, cy, state := sdl.GetGlobalMouseState()
	if x != nil {
		*x = cx
	}
	if y != nil {
		*y = cy
	}
	return state
}
