package libs

import (
	"github.com/dennwc/cxgo/types"
)

const (
	sdl2H          = "SDL2/SDL.h"
	sdl2OpenGLH    = "SDL2/SDL_opengl.h"
	sdl2OpenGLExtH = "SDL2/SDL_opengl_glext.h"
	sdl2StdIncH    = "SDL2/SDL_stdinc.h"
)

func init() {
	RegisterLibrary(sdl2H, func(c *Env) *Library {
		boolT := types.BoolT()
		sintT := types.IntT(4)
		uintT := types.UintT(4)
		sint16T := types.IntT(2)
		uint16T := types.UintT(2)
		uint8T := types.UintT(1)
		strT := c.C().String()
		bufT := c.PtrT(nil)
		scanT := types.NamedT("sdl.Scancode", sintT)
		blendT := types.NamedTGo("SDL_BlendMode", "sdl.BlendMode", sintT)
		timerT := types.NamedTGo("SDL_TimerID", "sdl.TimerID", sintT)
		keymodT := types.NamedTGo("SDL_Keymod", "sdl.Keymod", sintT)
		fingerT := types.NamedTGo("SDL_FingerID", "sdl.FingerID", sintT)
		rectT := types.NamedT("sdl.Rect", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("x", "X", sint16T)},
			{Name: types.NewIdentGo("y", "Y", sint16T)},
			{Name: types.NewIdentGo("w", "W", uint16T)},
			{Name: types.NewIdentGo("h", "H", uint16T)},
		}))
		muT := types.NamedT("sdl.Mutex", c.MethStructT(map[string]*types.FuncType{
			"Lock":    c.FuncTT(nil),
			"Unlock":  c.FuncTT(nil),
			"Destroy": c.FuncTT(nil),
		}))
		ksymT := types.NamedT("sdl.Keysym", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("scancode", "Scancode", scanT)},
		}))
		evkeybT := types.NamedT("sdl.KeyboardEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("keysym", "Keysym", ksymT)},
			{Name: types.NewIdentGo("state", "Env", uint8T)},
		}))
		evwinT := types.NamedT("sdl.WindowEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("event", "Event", uint8T)},
		}))
		evmbT := types.NamedT("sdl.MouseButtonEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("button", "Button", uint8T)},
			{Name: types.NewIdentGo("state", "Env", uint8T)},
			{Name: types.NewIdentGo("x", "X", sintT)},
			{Name: types.NewIdentGo("y", "Y", sintT)},
		}))
		evmmT := types.NamedT("sdl.MouseMotionEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("x", "X", sintT)},
			{Name: types.NewIdentGo("y", "Y", sintT)},
			{Name: types.NewIdentGo("xrel", "Xrel", sintT)},
			{Name: types.NewIdentGo("yrel", "Yrel", sintT)},
		}))
		evmwT := types.NamedT("sdl.MouseWheelEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("x", "X", sintT)},
			{Name: types.NewIdentGo("y", "Y", sintT)},
		}))
		textT := c.C().BytesN(32)
		evtxtT := types.NamedT("sdl.TextInputEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("text", "Text", textT)},
		}))
		evtxteT := types.NamedT("sdl.TextEditingEvent", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("text", "Text", textT)},
		}))
		surfaceT := types.NamedTGo("SDL_Surface", "sdl.Surface", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("w", "W", sintT)},
			{Name: types.NewIdentGo("h", "H", sintT)},
			{Name: types.NewIdentGo("pitch", "Pitch", sintT)},
			{Name: types.NewIdentGo("pixels", "Pixels", bufT)},
			{Name: types.NewIdentGo("clip_rect", "ClipRect", rectT)},
			{Name: types.NewIdent("Lock", c.FuncTT(sintT))},
			{Name: types.NewIdent("Unlock", c.FuncTT(nil))},
			{Name: types.NewIdent("Free", c.FuncTT(nil))},
			{Name: types.NewIdent("SetColorKey", c.FuncTT(sintT, sintT, uintT))},
			{Name: types.NewIdent("GetClipRect", c.FuncTT(nil, c.PtrT(rectT)))},
			{Name: types.NewIdent("SetBlendMode", c.FuncTT(sintT, blendT))},
		}))
		textureT := types.NamedTGo("SDL_Texture", "sdl.Texture", c.MethStructT(map[string]*types.FuncType{
			"Destroy": c.FuncTT(nil),
		}))
		evT := types.NamedT("sdl.Event", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("type", "Type", sintT)},
			{Name: types.NewIdentGo("edit", "Edit", evtxteT)},
			{Name: types.NewIdentGo("text", "Text", evtxtT)},
			{Name: types.NewIdentGo("key", "Key", evkeybT)},
			{Name: types.NewIdentGo("button", "Button", evmbT)},
			{Name: types.NewIdentGo("motion", "Motion", evmmT)},
			{Name: types.NewIdentGo("wheel", "Wheel", evmwT)},
			{Name: types.NewIdentGo("window", "Window", evwinT)},
		}))
		glCtxT := types.NamedTGo("SDL_GLContext", "sdl.GLContext", types.StructT(nil))
		winT := types.NamedT("sdl.Window", c.MethStructT(map[string]*types.FuncType{
			"SetFullscreen":     c.FuncTT(sintT, uintT),
			"GetSize":           c.FuncTT(nil, c.PtrT(sintT), c.PtrT(sintT)),
			"GetPosition":       c.FuncTT(nil, c.PtrT(sintT), c.PtrT(sintT)),
			"SetSize":           c.FuncTT(nil, sintT, sintT),
			"SetPosition":       c.FuncTT(nil, sintT, sintT),
			"Minimize":          c.FuncTT(nil),
			"Restore":           c.FuncTT(nil),
			"SetGrab":           c.FuncTT(nil, boolT),
			"SetResizable":      c.FuncTT(nil, boolT),
			"SetBordered":       c.FuncTT(nil, boolT),
			"SetTitle":          c.FuncTT(nil, strT),
			"GetDisplayIndex":   c.FuncTT(sintT, strT),
			"GLGetDrawableSize": c.FuncTT(nil, c.PtrT(sintT), c.PtrT(sintT)),
			"GLSwap":            c.FuncTT(nil),
			"GLCreateContext":   c.FuncTT(glCtxT),
		}))
		renderT := types.NamedTGo("SDL_Renderer", "sdl.Renderer", c.MethStructT(map[string]*types.FuncType{
			"Present": c.FuncTT(nil),
			"Destroy": c.FuncTT(nil),
		}))
		timerFuncT := types.NamedTGo("SDL_TimerCallback", "sdl.TimerFunc", c.FuncTT(uintT, uintT, c.PtrT(nil)))
		return &Library{
			Imports: map[string]string{
				"sdl": RuntimePrefix + "sdl",
			},
			Types: map[string]types.Type{
				"SDL_Window":           winT,
				"SDL_Renderer":         renderT,
				"SDL_Scancode":         scanT,
				"SDL_BlendMode":        blendT,
				"SDL_Keymod":           keymodT,
				"SDL_FingerID":         fingerT,
				"SDL_TimerID":          timerT,
				"SDL_Rect":             rectT,
				"SDL_mutex":            muT,
				"SDL_Keysym":           ksymT,
				"SDL_KeyboardEvent":    evkeybT,
				"SDL_WindowEvent":      evwinT,
				"SDL_MouseButtonEvent": evmbT,
				"SDL_MouseMotionEvent": evmmT,
				"SDL_MouseWheelEvent":  evmwT,
				"SDL_TextInputEvent":   evtxtT,
				"SDL_TextEditingEvent": evtxteT,
				"SDL_Event":            evT,
				"SDL_Surface":          surfaceT,
				"SDL_Texture":          textureT,
				"SDL_GLContext":        glCtxT,
				"SDL_TimerCallback":    timerFuncT,
			},
			Idents: map[string]*types.Ident{
				"SDL_itoa":                       types.NewIdent("sdl.Itoa", c.FuncTT(strT, sintT, strT, sintT)),
				"SDL_uitoa":                      types.NewIdent("sdl.Uitoa", c.FuncTT(strT, uintT, strT, sintT)),
				"SDL_GetError":                   types.NewIdent("sdl.GetError", c.FuncTT(strT)),
				"SDL_GetTicks":                   types.NewIdent("sdl.GetTicks", c.FuncTT(uintT)),
				"SDL_PollEvent":                  types.NewIdent("sdl.PollEvent", c.FuncTT(sintT, c.PtrT(evT))),
				"SDL_Delay":                      types.NewIdent("sdl.Delay", c.FuncTT(nil, uintT)),
				"SDL_Init":                       types.NewIdent("sdl.Init", c.FuncTT(sintT, uintT)),
				"SDL_CreateMutex":                types.NewIdent("sdl.CreateMutex", c.FuncTT(c.PtrT(muT))),
				"SDL_CreateWindow":               types.NewIdent("sdl.CreateWindow", c.FuncTT(c.PtrT(winT), strT, sintT, sintT, sintT, sintT, uintT)),
				"SDL_CreateRGBSurfaceWithFormat": types.NewIdentGo("SDL_CreateRGBSurfaceWithFormat", "sdl.CreateRGBSurfaceWithFormat", c.FuncTT(c.PtrT(surfaceT), uintT, sintT, sintT, sintT, uintT)),
				"SDL_BlitSurface":                types.NewIdentGo("SDL_BlitSurface", "sdl.BlitSurface", c.FuncTT(sintT, c.PtrT(surfaceT), c.PtrT(rectT), c.PtrT(surfaceT), c.PtrT(rectT))),
				"SDL_BlitScaled":                 types.NewIdentGo("SDL_BlitScaled", "sdl.BlitScaled", c.FuncTT(sintT, c.PtrT(surfaceT), c.PtrT(rectT), c.PtrT(surfaceT), c.PtrT(rectT))),
				"SDL_AddTimer":                   types.NewIdentGo("SDL_AddTimer", "sdl.AddTimer", c.FuncTT(timerT, uintT, timerFuncT, c.PtrT(nil))),
				"SDL_RemoveTimer":                types.NewIdentGo("SDL_RemoveTimer", "sdl.RemoveTimer", c.FuncTT(boolT, timerT)),
				"SDL_GL_DeleteContext":           types.NewIdentGo("SDL_GL_DeleteContext", "sdl.GLDeleteContext", c.FuncTT(nil, glCtxT)),
				"SDL_GL_SetSwapInterval":         types.NewIdentGo("SDL_GL_SetSwapInterval", "sdl.GLSetSwapInterval", c.FuncTT(sintT, sintT)),
				"SDL_GetDisplayBounds":           types.NewIdentGo("SDL_GetDisplayBounds", "sdl.GetDisplayBounds", c.FuncTT(sintT, sintT, c.PtrT(rectT))),
				"SDL_GetGlobalMouseState":        types.NewIdentGo("SDL_GetGlobalMouseState", "sdl.GetGlobalMouseState", c.FuncTT(uintT, c.PtrT(sintT), c.PtrT(sintT))),
				"SDL_GetEventState":              types.NewIdentGo("SDL_GetEventState", "sdl.GetEventState", c.FuncTT(uint8T, uintT)),
				"SDL_SetRelativeMouseMode":       types.NewIdentGo("SDL_SetRelativeMouseMode", "sdl.SetRelativeMouseMode", c.FuncTT(sintT, boolT)),
				"SDL_GetModState":                types.NewIdentGo("SDL_GetModState", "sdl.GetModState", c.FuncTT(keymodT)),
				"SDL_StartTextInput":             types.NewIdentGo("SDL_StartTextInput", "sdl.StartTextInput", c.FuncTT(nil)),
				"SDL_StopTextInput":              types.NewIdentGo("SDL_StopTextInput", "sdl.StopTextInput", c.FuncTT(nil)),
				"SDL_CreateRenderer":             types.NewIdentGo("SDL_CreateRenderer", "sdl.CreateRenderer", c.FuncTT(c.PtrT(renderT), c.PtrT(winT), sintT, uintT)),
				"SDL_CreateTextureFromSurface":   types.NewIdentGo("SDL_CreateTextureFromSurface", "sdl.CreateTextureFromSurface", c.FuncTT(c.PtrT(textureT), c.PtrT(renderT), c.PtrT(surfaceT))),
				"SDL_RenderCopy":                 types.NewIdentGo("SDL_RenderCopy", "sdl.RenderCopy", c.FuncTT(sintT, c.PtrT(renderT), c.PtrT(textureT), c.PtrT(rectT), c.PtrT(rectT))),

				"SDL_BLENDMODE_NONE":            types.NewIdentGo("SDL_BLENDMODE_NONE", "sdl.BLENDMODE_NONE", blendT),
				"SDL_PIXELFORMAT_RGBA5551":      types.NewIdentGo("SDL_PIXELFORMAT_RGBA5551", "sdl.PIXELFORMAT_RGBA5551", sintT),
				"SDL_PIXELFORMAT_RGB565":        types.NewIdentGo("SDL_PIXELFORMAT_RGB565", "sdl.PIXELFORMAT_RGB565", sintT),
				"SDL_PIXELFORMAT_RGB555":        types.NewIdentGo("SDL_PIXELFORMAT_RGB555", "sdl.PIXELFORMAT_RGB555", sintT),
				"SDL_RENDERER_ACCELERATED":      types.NewIdentGo("SDL_RENDERER_ACCELERATED", "sdl.RENDERER_ACCELERATED", sintT),
				"SDL_RENDERER_PRESENTVSYNC":     types.NewIdentGo("SDL_RENDERER_PRESENTVSYNC", "sdl.RENDERER_PRESENTVSYNC", sintT),
				"SDL_NUM_SCANCODES":             types.NewIdentGo("SDL_NUM_SCANCODES", "sdl.NUM_SCANCODES", sintT),
				"SDL_PRESSED":                   types.NewIdentGo("SDL_PRESSED", "sdl.PRESSED", sintT),
				"SDL_SCANCODE_SPACE":            types.NewIdentGo("SDL_SCANCODE_SPACE", "sdl.SCANCODE_SPACE", sintT),
				"KMOD_LSHIFT":                   types.NewIdentGo("KMOD_LSHIFT", "sdl.KMOD_LSHIFT", sintT),
				"KMOD_RSHIFT":                   types.NewIdentGo("KMOD_RSHIFT", "sdl.KMOD_RSHIFT", sintT),
				"KMOD_RALT":                     types.NewIdentGo("KMOD_RALT", "sdl.KMOD_RALT", sintT),
				"SDL_INIT_VIDEO":                types.NewIdent("sdl.INIT_VIDEO", uintT),
				"SDL_INIT_TIMER":                types.NewIdent("sdl.INIT_TIMER", uintT),
				"SDL_WINDOWPOS_UNDEFINED":       types.NewIdent("sdl.WINDOWPOS_UNDEFINED", sintT),
				"SDL_WINDOWPOS_CENTERED":        types.NewIdent("sdl.WINDOWPOS_CENTERED", sintT),
				"SDL_WINDOW_RESIZABLE":          types.NewIdent("sdl.WINDOW_RESIZABLE", uintT),
				"SDL_WINDOW_OPENGL":             types.NewIdent("sdl.WINDOW_OPENGL", uintT),
				"SDL_WINDOW_FULLSCREEN_DESKTOP": types.NewIdent("sdl.WINDOW_FULLSCREEN_DESKTOP", uintT),
			},
			// TODO
			Header: `
#include <` + BuiltinH + `>

#define Uint32 _cxgo_uint32
#define Uint16 _cxgo_uint16
#define Uint8 _cxgo_uint8
#define Sint32 _cxgo_sint32
#define Sint16 _cxgo_sint16
#define Sint8 _cxgo_sint8
#define SDL_bool _cxgo_bool
typedef Sint32 SDL_FingerID;
typedef Sint32 SDL_Scancode;
typedef Sint32 SDL_Keymod;
typedef Sint32 SDL_BlendMode;
typedef Sint32 SDL_TimerID;


#define SDL_TRUE 1
#define SDL_FALSE 0
const Sint32 SDL_PIXELFORMAT_RGBA5551 = 1;
const Sint32 SDL_PIXELFORMAT_RGB565 = 2;
const Sint32 SDL_PIXELFORMAT_RGB555 = 3;

const Sint32 SDL_RENDERER_ACCELERATED = 1;
const Sint32 SDL_RENDERER_PRESENTVSYNC = 2;

#define SDL_NUM_SCANCODES 512

typedef struct{
  Sint16 x, y;
  Uint16 w, h;
} SDL_Rect;

#define SDL_TEXTEDITING 770
#define SDL_TEXTINPUT 771
#define SDL_KEYDOWN 768
#define SDL_KEYUP 769
#define SDL_MOUSEBUTTONDOWN 1025
#define SDL_MOUSEBUTTONUP 1026
#define SDL_MOUSEMOTION 1024
#define SDL_MOUSEWHEEL 1027
#define SDL_WINDOWEVENT 512

#define SDL_WINDOWEVENT_FOCUS_LOST 13
#define SDL_WINDOWEVENT_FOCUS_GAINED 12

#define SDL_BUTTON_LEFT 1
#define SDL_BUTTON_RIGHT 3
#define SDL_BUTTON_MIDDLE 2

const Sint32 SDL_PRESSED = 1;

const Sint32 SDL_SCANCODE_SPACE = 1;
const Sint32 KMOD_LSHIFT = 1;
const Sint32 KMOD_RSHIFT = 2;
const Sint32 KMOD_RALT = 3;

const Uint32 SDL_INIT_VIDEO = 1;
const Uint32 SDL_INIT_TIMER = 2;
const Sint32 SDL_WINDOWPOS_UNDEFINED = -1;
const Sint32 SDL_WINDOWPOS_CENTERED = -2;
const Uint32 SDL_WINDOW_RESIZABLE = 1;
const Uint32 SDL_WINDOW_OPENGL = 2;
const SDL_BlendMode SDL_BLENDMODE_NONE = 1;

typedef struct{
    SDL_Scancode scancode;
} SDL_Keysym;

typedef struct{
    SDL_Keysym keysym;
    Uint8 state;
} SDL_KeyboardEvent;

typedef struct{
    Uint8 button;
    Uint8 state;
    Sint32 x,y;
} SDL_MouseButtonEvent;

typedef struct{
    Sint32 x, y;
    Sint32 xrel, yrel;
} SDL_MouseMotionEvent;

typedef struct{
    Sint32 x, y;
} SDL_MouseWheelEvent;

typedef struct{} SDL_TouchFingerEvent;

typedef struct{
    char text[32];
} SDL_TextEditingEvent;

typedef struct{
    char text[32];
} SDL_TextInputEvent;

typedef struct{
    Uint8 event;
} SDL_WindowEvent;

typedef struct{
    int type;
    SDL_TextEditingEvent edit;
    SDL_TextInputEvent text;
    SDL_KeyboardEvent key;
    SDL_MouseButtonEvent button;
    SDL_MouseMotionEvent motion;
    SDL_MouseWheelEvent wheel;
    SDL_WindowEvent window;
} SDL_Event;

const Uint32 SDL_WINDOW_FULLSCREEN_DESKTOP = 0;

typedef struct SDL_GLContext {} SDL_GLContext;

typedef struct SDL_Window {
	Sint32 (*SetFullscreen)(Uint32 flags);
	void (*GetSize)(Sint32 *w, Sint32 *h);
	void (*GetPosition)(Sint32 *x, Sint32 *y);
	void (*SetSize)(Sint32 w, Sint32 h);
	void (*SetPosition)(Sint32 x, Sint32 y);
	void (*Minimize)(void);
	void (*Restore)(void);
	void (*SetGrab)(SDL_bool);
	void (*SetResizable)(SDL_bool);
	void (*SetBordered)(SDL_bool);
	void (*SetTitle)(const char*);
	Sint32 (*GetDisplayIndex)(void);
	void (*GLGetDrawableSize)(Sint32* w, Sint32* h);
	void (*GLSwapWindow)(void);
	SDL_GLContext (*GLCreateContext)(void);
} SDL_Window;

const char* SDL_GetError(void);

Uint32 SDL_GetTicks(void);
SDL_Window* SDL_CreateWindow(const char* title, Sint32 x, Sint32 y, Sint32 w, Sint32 h, Uint32 flags);
#define SDL_SetWindowFullscreen(win, flags) ((SDL_Window*)win)->SetFullscreen(flags)
#define SDL_SetWindowSize(win, w, h) ((SDL_Window*)win)->SetSize(w, h)
#define SDL_SetWindowPosition(win, x, y) ((SDL_Window*)win)->SetPosition(x, y)
#define SDL_SetWindowGrab(win, v) ((SDL_Window*)win)->SetGrab(v)
#define SDL_SetWindowTitle(win, v) ((SDL_Window*)win)->SetTitle(v)
#define SDL_MinimizeWindow(win) ((SDL_Window*)win)->Minimize()
#define SDL_RestoreWindow(win) ((SDL_Window*)win)->Restore()
#define SDL_GetWindowDisplayIndex(win) ((SDL_Window*)win)->GetDisplayIndex()
#define SDL_SetWindowResizable(win, v) ((SDL_Window*)win)->SetResizable(v)
#define SDL_SetWindowBordered(win, v) ((SDL_Window*)win)->SetBordered(v)
#define SDL_GetWindowPosition(win, x, y) ((SDL_Window*)win)->GetPosition(x, y)
#define SDL_GetWindowSize(win, w, h) ((SDL_Window*)win)->GetSize(w, h)

#define SDL_GL_GetDrawableSize(win, w, h) ((SDL_Window*)win)->GLGetDrawableSize(w, h)
#define SDL_GL_SwapWindow(win) ((SDL_Window*)win)->GLSwap()
#define SDL_GL_CreateContext(win) ((SDL_Window*)win)->GLCreateContext()

typedef struct SDL_Renderer {
	void (*Destroy)(void);
	void (*Present)(void);
} SDL_Renderer;

SDL_Renderer* SDL_CreateRenderer(SDL_Window* window, Sint32 index, Uint32 flags);

#define SDL_RenderPresent(r) ((SDL_Renderer*)r)->Present()
#define SDL_DestroyRenderer(r) ((SDL_Renderer*)r)->Destroy()

char *SDL_itoa(Sint32 value, char *str, Sint32 radix);
char *SDL_uitoa(Uint32 value, char *str, Sint32 radix);
void SDL_Delay(Uint32 ms);

typedef struct SDL_mutex {
	void (*Lock)(void);
	void (*Unlock)(void);
	void (*Destroy)(void);
} SDL_mutex;

SDL_mutex* SDL_CreateMutex(void);

#define SDL_DestroyMutex(m) ((SDL_mutex*)m)->Destroy()
#define SDL_LockMutex(m) ((SDL_mutex*)m)->Lock()
#define SDL_UnlockMutex(m) ((SDL_mutex*)m)->Unlock()


typedef struct SDL_Surface {
    Sint32 w, h;
    Sint32 pitch;
    void* pixels;
    SDL_Rect clip_rect;

	Sint32 (*Lock)(void);
	void (*Unlock)(void);
	void (*Free)(void);
	Sint32 (*SetColorKey)(Sint32 flag, Uint32 key);
	void (*GetClipRect)(SDL_Rect* rect);
	Sint32 (*SetBlendMode)(SDL_BlendMode blendMode);
} SDL_Surface;

SDL_Surface* SDL_CreateRGBSurfaceWithFormat(Uint32 flags, Sint32 width, Sint32 height, Sint32 depth, Uint32 format);

#define SDL_LockSurface(s) ((SDL_Surface*)s)->Lock()
#define SDL_UnlockSurface(s) ((SDL_Surface*)s)->Unlock()
#define SDL_FreeSurface(s) ((SDL_Surface*)s)->Free()
#define SDL_SetColorKey(s, f, k) ((SDL_Surface*)s)->SetColorKey(f, k)
#define SDL_GetClipRect(s, r) ((SDL_Surface*)s)->GetClipRect(r)
#define SDL_SetSurfaceBlendMode(s, m) ((SDL_Surface*)s)->SetBlendMode(m)
Sint32 SDL_BlitSurface(SDL_Surface* src, const SDL_Rect* srcrect, SDL_Surface* dst, SDL_Rect* dstrect);
Sint32 SDL_BlitScaled(SDL_Surface* src, const SDL_Rect* srcrect, SDL_Surface* dst, SDL_Rect* dstrect);

typedef struct SDL_Texture {
	void (*Destroy)(void);
} SDL_Texture;
SDL_Texture* SDL_CreateTextureFromSurface(SDL_Renderer* renderer, SDL_Surface*  surface);
#define SDL_DestroyTexture(t) ((SDL_Texture*)t)->Destroy()

Sint32 SDL_RenderCopy(SDL_Renderer* renderer, SDL_Texture* texture, const SDL_Rect* srcrect, const SDL_Rect* dstrect);


Sint32 SDL_GL_SetSwapInterval(Sint32 interval);
void SDL_GL_DeleteContext(SDL_GLContext context);
Sint32 SDL_Init(Uint32 flags);
Sint32 SDL_PollEvent(SDL_Event* event);
void SDL_StartTextInput(void);
void SDL_StopTextInput(void);
SDL_Keymod SDL_GetModState(void);
Sint32 SDL_SetRelativeMouseMode(SDL_bool enabled);
Uint8 SDL_GetEventState(Uint32 type);
Uint32 SDL_GetGlobalMouseState(int* x, int* y);
Sint32 SDL_GetDisplayBounds(Sint32 displayIndex, SDL_Rect* rect);

typedef Uint32 (*SDL_TimerCallback)(Uint32 interval, void* param);
SDL_TimerID SDL_AddTimer(Uint32 interval, SDL_TimerCallback callback, void* param);
SDL_bool SDL_RemoveTimer(SDL_TimerID id);
`,
		}
	})
	RegisterLibrary(sdl2OpenGLH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + BuiltinH + `>
`,
		}
	})
	RegisterLibrary(sdl2OpenGLExtH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + BuiltinH + `>
`,
		}
	})
	RegisterLibrary(sdl2StdIncH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + BuiltinH + `>
`,
		}
	})
}
