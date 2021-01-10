package sdl

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func GLSetSwapInterval(interval int32) int32 {
	err := sdl.GLSetSwapInterval(int(interval))
	if err != nil {
		log.Println("sdl.GLSetSwapInterval:", err)
		return -1
	}
	return 0
}

type GLContext struct {
	c sdl.GLContext
}

func GLDeleteContext(c GLContext) {
	sdl.GLDeleteContext(c.c)
}

func (win *Window) GLGetDrawableSize(w, h *int32) {
	cw, ch := win.w.GLGetDrawableSize()
	if w != nil {
		*w = cw
	}
	if h != nil {
		*h = ch
	}
}

func (win *Window) GLSwap() {
	win.w.GLSwap()
}

func (win *Window) GLCreateContext() GLContext {
	c, err := win.w.GLCreateContext()
	if err != nil {
		log.Printf("sdl.Window.GLCreateContext: %v", err)
		return GLContext{}
	}
	glInit()
	return GLContext{c: c}
}

func glInit() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	log.Println("OpenGL version:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("GLSL version:", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(func(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, _ unsafe.Pointer) {
		log.Printf("GL: %d, %d, %d: %s", gltype, id, severity, message)
	}, nil)
}
