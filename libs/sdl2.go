package libs

import (
	"github.com/gotranspile/cxgo/types"
)

const (
	sdl2H = "SDL2/SDL.h"
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
				"sdl": RuntimeOrg + "/sdl",
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
		}
	})
}
