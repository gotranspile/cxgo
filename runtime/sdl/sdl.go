package sdl

import (
	"log"
	"strconv"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	noFullscreen = true
	noMouseGrab  = true
)

const (
	INIT_VIDEO = sdl.INIT_VIDEO
	INIT_TIMER = sdl.INIT_TIMER

	WINDOWPOS_UNDEFINED = sdl.WINDOWPOS_UNDEFINED
	WINDOWPOS_CENTERED  = sdl.WINDOWPOS_CENTERED

	WINDOW_RESIZABLE          = sdl.WINDOW_RESIZABLE
	WINDOW_OPENGL             = sdl.WINDOW_OPENGL
	WINDOW_FULLSCREEN_DESKTOP = sdl.WINDOW_FULLSCREEN_DESKTOP

	PIXELFORMAT_RGBA5551 = sdl.PIXELFORMAT_RGBA5551
	PIXELFORMAT_RGB565   = sdl.PIXELFORMAT_RGB565
	PIXELFORMAT_RGB555   = sdl.PIXELFORMAT_RGB555

	RENDERER_ACCELERATED  = sdl.RENDERER_ACCELERATED
	RENDERER_PRESENTVSYNC = sdl.RENDERER_PRESENTVSYNC
)

func Itoa(val int32, str *byte, radix int32) *byte {
	s := strconv.FormatInt(int64(val), int(radix))
	return libc.StrCpy(str, libc.CString(s))
}

func Uitoa(val uint32, str *byte, radix int32) *byte {
	s := strconv.FormatUint(uint64(val), int(radix))
	return libc.StrCpy(str, libc.CString(s))
}

func Init(flags uint32) int32 {
	err := sdl.Init(flags)
	log.Printf("SDL_Init(%x): %v", flags, err)
	if err == nil {
		return 0
	}
	return -1 // TODO
}

func GetError() *byte {
	err := sdl.GetError()
	if err == nil {
		return nil
	}
	return libc.CString(err.Error())
}

type Mutex struct {
	m *sdl.Mutex
}

func CreateMutex() *Mutex {
	m, err := sdl.CreateMutex()
	if err != nil {
		log.Println("sdl.CreateMutex:", err)
		return nil
	}
	return &Mutex{m: m}
}

func (m *Mutex) Lock() {
	m.m.Lock()
}

func (m *Mutex) Unlock() {
	m.m.Unlock()
}

func (m *Mutex) Destroy() {
	m.m.Destroy()
}

const (
	BLENDMODE_NONE = BlendMode(sdl.BLENDMODE_NONE)
)

type BlendMode int32

type Surface struct {
	W, H     int32
	Pitch    int32
	Pixels   unsafe.Pointer
	ClipRect Rect
	s        *sdl.Surface
}

func CreateRGBSurfaceWithFormat(flags uint32, w, h, depth int32, format uint32) *Surface {
	s, err := sdl.CreateRGBSurfaceWithFormat(flags, w, h, depth, format)
	if err != nil {
		log.Println("sdl.CreateRGBSurfaceWithFormat:", err)
		return nil
	}
	var r Rect
	r.setFrom(&s.ClipRect)
	return &Surface{
		W:        s.W,
		H:        s.H,
		Pitch:    s.Pitch,
		Pixels:   s.Data(),
		ClipRect: r,
		s:        s,
	}
}

func BlitSurface(src *Surface, srcRect *Rect, dst *Surface, dstRect *Rect) int32 {
	var (
		srcr sdl.Rect
		dstr sdl.Rect
	)
	srcRect.setTo(&srcr)
	dstRect.setTo(&dstr)
	err := src.s.Blit(&srcr, dst.s, &dstr)
	if err != nil {
		log.Printf("sdl.BlitSurface: %v", err)
		return 1
	}
	return 0
}

func BlitScaled(src *Surface, srcRect *Rect, dst *Surface, dstRect *Rect) int32 {
	var (
		srcr sdl.Rect
		dstr sdl.Rect
	)
	srcRect.setTo(&srcr)
	dstRect.setTo(&dstr)
	err := src.s.BlitScaled(&srcr, dst.s, &dstr)
	if err != nil {
		log.Printf("sdl.BlitScaled: %v", err)
		return 1
	}
	return 0
}

func (s *Surface) Lock() int32 {
	err := s.s.Lock()
	if err != nil {
		log.Printf("sdl.Surface.Lock: %v", err)
		return 1
	}
	return 0
}

func (s *Surface) Unlock() {
	s.s.Unlock()
}

func (s *Surface) Free() {
	s.s.Free()
}

func (s *Surface) SetColorKey(a1 int32, a2 uint32) int32 {
	err := s.s.SetColorKey(a1 != 0, a2)
	if err != nil {
		log.Printf("sdl.Surface.Lock: %v", err)
		return 1
	}
	return 0
}

func (s *Surface) SetBlendMode(m BlendMode) int32 {
	err := s.s.SetBlendMode(sdl.BlendMode(m))
	if err != nil {
		log.Printf("sdl.Surface.SetBlendMode: %v", err)
		return 1
	}
	return 0
}

func (s *Surface) GetClipRect(r *Rect) {
	var sr sdl.Rect
	s.s.GetClipRect(&sr)
	r.setFrom(&sr)
}

type Window struct {
	w *sdl.Window
}

func CreateWindow(title *byte, x, y, w, h int32, flags uint32) *Window {
	stitle := libc.GoString(title)
	// force disable fullscreen
	if noFullscreen {
		flags = flags &^ sdl.WINDOW_FULLSCREEN
	}
	err := sdlSetGLVersion(3, 0)
	if err != nil {
		log.Printf("SetGLCoreVersion(): %v", err)
		return nil
	}
	win, err := sdl.CreateWindow(stitle, x, y, w, h, flags)
	log.Printf("sdl.CreateWindow(%q, %d, %d, %d, %d, %x): %v", stitle, x, y, w, h, flags, err)
	if err != nil {
		return nil
	}
	return &Window{w: win}
}

func (win *Window) SDLWindow() *sdl.Window {
	if win == nil {
		return nil
	}
	return win.w
}

func (win *Window) SetFullscreen(flags uint32) int32 {
	err := win.SDLWindow().SetFullscreen(flags)
	if err != nil {
		log.Printf("SDL_SetWindowFullscreen(%p, %x): %v", win, flags, err)
		return -1
	}
	return 0
}

func (win *Window) GetSize(w, h *int32) {
	cw, ch := win.SDLWindow().GetSize()
	if w != nil {
		*w = cw
	}
	if h != nil {
		*h = ch
	}
}

func (win *Window) GetPosition(x, y *int32) {
	cx, cy := win.SDLWindow().GetPosition()
	if x != nil {
		*x = cx
	}
	if y != nil {
		*y = cy
	}
}

func (win *Window) SetSize(w, h int32) {
	win.SDLWindow().SetSize(w, h)
}

func (win *Window) SetPosition(x, y int32) {
	win.SDLWindow().SetPosition(x, y)
}

func (win *Window) SetResizable(val bool) {
	win.SDLWindow().SetResizable(val)
}
func (win *Window) SetBordered(val bool) {
	win.SDLWindow().SetBordered(val)
}
func (win *Window) SetGrab(val bool) {
	if noMouseGrab {
		return
	}
	win.SDLWindow().SetGrab(val)
}
func (win *Window) Minimize() {
	win.SDLWindow().Minimize()
}
func (win *Window) Restore() {
	win.SDLWindow().Restore()
}
func (win *Window) SetTitle(title *byte) {
	win.SDLWindow().SetTitle(libc.GoString(title))
}

func (win *Window) GetDisplayIndex() int32 {
	ind, err := win.SDLWindow().GetDisplayIndex()
	if err != nil {
		log.Println("SDL_GetWindowDisplayIndex:", err)
		return -1
	}
	return int32(ind)
}

func sdlSetGLVersion(major, minor int) error {
	err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, major)
	if err != nil {
		return err
	}
	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, minor)
	if err != nil {
		return err
	}
	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	if err != nil {
		return err
	}
	return nil
}

type Rect struct {
	X, Y int16
	W, H uint16
}

func (r *Rect) setFrom(r2 *sdl.Rect) {
	r.X = int16(r2.X)
	r.Y = int16(r2.Y)
	r.W = uint16(r2.W)
	r.H = uint16(r2.H)
}

func (r *Rect) setTo(r2 *sdl.Rect) {
	r2.X = int32(r.X)
	r2.Y = int32(r.Y)
	r2.W = int32(r.W)
	r2.H = int32(r.H)
}

func GetDisplayBounds(ind int32, rect *Rect) int32 {
	out, err := sdl.GetDisplayBounds(int(ind))
	if err != nil {
		log.Println("sdl.GetDisplayBounds:", err)
		return -1
	}
	rect.setFrom(&out)
	return 0
}

type Renderer struct {
	r *sdl.Renderer
}

func CreateRenderer(win *Window, index int32, flags uint32) *Renderer {
	r, err := sdl.CreateRenderer(win.w, int(index), flags)
	if err != nil {
		log.Println("sdl.CreateRenderer:", err)
		return nil
	}
	return &Renderer{r: r}
}

func (r *Renderer) Present() {
	r.r.Present()
}

func (r *Renderer) Destroy() {
	r.r.Destroy()
}

type Texture struct {
	t *sdl.Texture
}

func CreateTextureFromSurface(r *Renderer, s *Surface) *Texture {
	t, err := r.r.CreateTextureFromSurface(s.s)
	if err != nil {
		log.Println("sdl.CreateTextureFromSurface:", err)
		return nil
	}
	return &Texture{t: t}
}

func (t *Texture) Destroy() {
	t.t.Destroy()
}

func RenderCopy(r *Renderer, t *Texture, src *Rect, dst *Rect) int32 {
	var sr, dr sdl.Rect
	src.setTo(&sr)
	dst.setTo(&dr)
	err := r.r.Copy(t.t, &sr, &dr)
	if err != nil {
		log.Println("sdl.RenderCopy:", err)
		return -1
	}
	return 0
}
