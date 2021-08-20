package libs

import (
	"github.com/gotranspile/cxgo/types"
)

const (
	glfw3H = "GLFW/glfw3.h"
)

func init() {
	RegisterLibrary(glfw3H, func(env *Env) *Library {
		hintT := types.NamedTGo("Hint", "glfw.Hint", env.Go().Int())
		actionT := types.NamedTGo("Action", "glfw.Action", env.Go().Int())
		joystickHatStateT := types.NamedTGo("JoystickHatState", "glfw.JoystickHatState", env.Go().Int())
		keyT := types.NamedTGo("Key", "glfw.Key", env.Go().Int())
		mouseButtonT := types.NamedTGo("MouseButton", "glfw.MouseButton", env.Go().Int())
		modKeyT := types.NamedTGo("ModifierKey", "glfw.ModifierKey", env.Go().Int())
		inputModeT := types.NamedTGo("InputMode", "glfw.InputMode", env.Go().Int())
		perEventT := types.NamedTGo("PeripheralEvent", "glfw.PeripheralEvent", env.Go().Int())
		windowPtrT := env.PtrT(nil)
		keyCbT := types.NamedTGo("GLFWkeyfun", "glfw.KeyCallback", env.FuncTT(nil, windowPtrT, keyT, env.Go().Int(), actionT, modKeyT))
		charCbT := types.NamedTGo("GLFWcharfun", "glfw.CharCallback", env.FuncTT(nil, windowPtrT, env.Go().Rune()))
		frameBufCb := types.NamedTGo("GLFWframebuffersizefun", "glfw.FramebufferSizeCallback", env.FuncTT(nil, windowPtrT, env.Go().Rune()))
		videoModeT := types.NamedTGo("GLFWvidmode", "glfw.VidMode", types.StructT([]*types.Field{
			{Name: types.NewIdentGo("width", "Width", env.Go().Int())},
			{Name: types.NewIdentGo("height", "Height", env.Go().Int())},
			{Name: types.NewIdentGo("redBits", "RedBits", env.Go().Int())},
			{Name: types.NewIdentGo("greenBits", "GreenBits", env.Go().Int())},
			{Name: types.NewIdentGo("blueBits", "BlueBits", env.Go().Int())},
			{Name: types.NewIdentGo("refreshRate", "RefreshRate", env.Go().Int())},
		}))
		monitorT := types.NamedTGo("GLFWmonitor", "glfw.Monitor", env.MethStructT(map[string]*types.FuncType{
			"GetVideoMode":    env.FuncTT(env.PtrT(videoModeT), nil),
			"GetPos":          env.FuncTT(nil, env.PtrT(env.Go().Int()), env.PtrT(env.Go().Int())), //TODO: fix sig
			"GetPhysicalSize": env.FuncTT(nil, env.PtrT(env.Go().Int()), env.PtrT(env.Go().Int())), //TODO: fix sig
			"GetName":         env.FuncTT(env.Go().String(), nil),
		}))
		joystickT := types.NamedTGo("_GLFWJoystick", "glfw.Joystick", env.MethStructT(map[string]*types.FuncType{
			"GetAxes":        env.FuncTT(types.SliceT(types.FloatT(4)), env.PtrT(env.Go().Int())),   // TODO: fix sig
			"GetButtons":     env.FuncTT(types.SliceT(actionT), env.PtrT(env.Go().Int())),           // TODO: fix sig
			"GetHats":        env.FuncTT(types.SliceT(joystickHatStateT), env.PtrT(env.Go().Int())), // TODO: fix sig
			"GetName":        env.FuncTT(env.Go().String(), nil),
			"GetGamepadName": env.FuncTT(env.Go().String(), nil),
			"IsGamepad":      env.FuncTT(env.Go().Bool(), nil),
			"GetGUID":        env.FuncTT(env.Go().String(), nil),
		}))
		windowT := types.NamedTGo("GLFWwindow", "glfw.Window", env.MethStructT(map[string]*types.FuncType{
			"MakeContextCurrent": env.FuncTT(nil, nil),
			"ShouldClose":        env.FuncTT(env.Go().Bool(), nil),
			"SwapBuffers":        env.FuncTT(nil, nil),
			"GetKey":             env.FuncTT(actionT, keyT),
			"GetInputMode":       env.FuncTT(env.Go().Int(), env.Go().Int()),
			"SetInputMode":       env.FuncTT(nil, env.Go().Int(), env.Go().Int()),
			"SetShouldClose":     env.FuncTT(nil, env.Go().Bool()),
			"GetFramebufferSize": env.FuncTT(nil, env.Go().Int(), env.Go().Int()), //FIXME: incorrect signature
			"Destroy":            env.FuncTT(nil, nil),
			"Focus":              env.FuncTT(nil, nil),
			"Maximize":           env.FuncTT(nil, nil),
			"Show":               env.FuncTT(nil, nil),
			"Hide":               env.FuncTT(nil, nil),
			"Iconify":            env.FuncTT(nil, nil),
			"Restore":            env.FuncTT(nil, nil),
			"SetTitle":           env.FuncTT(nil, env.C().String()),
			"SetSize":            env.FuncTT(nil, env.Go().Int(), env.Go().Int()),
			"SetPos":             env.FuncTT(nil, env.Go().Int(), env.Go().Int()),
			// callbacks
			"SetKeyCallback":             env.FuncTT(keyCbT, keyCbT),
			"SetCharCallback":            env.FuncTT(charCbT, charCbT),
			"SetFramebufferSizeCallback": env.FuncTT(frameBufCb, frameBufCb),
			// no real implementation
			"GetWindowUserPointer": env.FuncTT(env.PtrT(nil), nil),
		}))
		windowPtrT.SetElem(windowT)
		l := &Library{
			Imports: map[string]string{
				"glfw": "github.com/go-gl/glfw/v3.3/glfw",
			},
			Types: map[string]types.Type{
				"GLFWwindow":    windowT,
				"GLFWmonitor":   monitorT,
				"GLFWvidmode":   videoModeT,
				"_GLFWjoystick": joystickT,
			},
			Idents: map[string]*types.Ident{
				// functions
				"glfwWindowHint": types.NewIdentGo("glfwWindowHint", "glfw.WindowHint", env.FuncTT(nil, hintT, env.Go().Int())),
				"glfwGetKeyName": types.NewIdentGo("glfwGetKeyName", "glfw.GetKeyName", env.FuncTT(env.Go().String(), keyT, env.Go().Int())),
				// constants
				"GLFW_OPENGL_CORE_PROFILE": types.NewIdentGo("GLFW_OPENGL_CORE_PROFILE", "glfw.OpenGLCoreProfile", env.Go().Int()),
				"GLFW_TRUE":                types.NewIdentGo("GLFW_TRUE", "glfw.True", env.Go().Int()),
				// key and button actions
				"GLFW_PRESS":   types.NewIdentGo("GLFW_PRESS", "glfw.Press", actionT),
				"GLFW_RELEASE": types.NewIdentGo("GLFW_RELEASE", "glfw.Release", actionT),
				"GLFW_REPEAT":  types.NewIdentGo("GLFW_REPEAT", "glfw.Repeat", actionT),
				/* The unknown key */
				"GLFW_KEY_UNKNOWN": types.NewIdentGo("GLFW_KEY_UNKNOWN", "glfw.KeyUnknown", keyT),
				/* Printable keys */
				"GLFW_KEY_SPACE":         types.NewIdentGo("GLFW_KEY_SPACE", "glfw.KeySpace", keyT),
				"GLFW_KEY_APOSTROPHE":    types.NewIdentGo("GLFW_KEY_APOSTROPHE", "glfw.KeyApostrophe", keyT),
				"GLFW_KEY_COMMA":         types.NewIdentGo("GLFW_KEY_COMMA", "glfw.KeyComma", keyT),
				"GLFW_KEY_MINUS":         types.NewIdentGo("GLFW_KEY_MINUS", "glfw.KeyMinus", keyT),
				"GLFW_KEY_PERIOD":        types.NewIdentGo("GLFW_KEY_PERIOD", "glfw.KeyPeriod", keyT),
				"GLFW_KEY_SLASH":         types.NewIdentGo("GLFW_KEY_SLASH", "glfw.KeySlash", keyT),
				"GLFW_KEY_0":             types.NewIdentGo("GLFW_KEY_0", "glfw.Key0", keyT),
				"GLFW_KEY_1":             types.NewIdentGo("GLFW_KEY_1", "glfw.Key1", keyT),
				"GLFW_KEY_2":             types.NewIdentGo("GLFW_KEY_2", "glfw.Key2", keyT),
				"GLFW_KEY_3":             types.NewIdentGo("GLFW_KEY_3", "glfw.Key3", keyT),
				"GLFW_KEY_4":             types.NewIdentGo("GLFW_KEY_4", "glfw.Key4", keyT),
				"GLFW_KEY_5":             types.NewIdentGo("GLFW_KEY_5", "glfw.Key5", keyT),
				"GLFW_KEY_6":             types.NewIdentGo("GLFW_KEY_6", "glfw.Key6", keyT),
				"GLFW_KEY_7":             types.NewIdentGo("GLFW_KEY_7", "glfw.Key7", keyT),
				"GLFW_KEY_8":             types.NewIdentGo("GLFW_KEY_8", "glfw.Key8", keyT),
				"GLFW_KEY_9":             types.NewIdentGo("GLFW_KEY_9", "glfw.Key9", keyT),
				"GLFW_KEY_SEMICOLON":     types.NewIdentGo("GLFW_KEY_SEMICOLON", "glfw.KeySemicolon", keyT),
				"GLFW_KEY_EQUAL":         types.NewIdentGo("GLFW_KEY_EQUAL", "glfw.KeyEqual", keyT),
				"GLFW_KEY_A":             types.NewIdentGo("GLFW_KEY_A", "glfw.KeyA", keyT),
				"GLFW_KEY_B":             types.NewIdentGo("GLFW_KEY_B", "glfw.KeyB", keyT),
				"GLFW_KEY_C":             types.NewIdentGo("GLFW_KEY_C", "glfw.KeyC", keyT),
				"GLFW_KEY_D":             types.NewIdentGo("GLFW_KEY_D", "glfw.KeyD", keyT),
				"GLFW_KEY_E":             types.NewIdentGo("GLFW_KEY_E", "glfw.KeyE", keyT),
				"GLFW_KEY_F":             types.NewIdentGo("GLFW_KEY_F", "glfw.KeyF", keyT),
				"GLFW_KEY_G":             types.NewIdentGo("GLFW_KEY_G", "glfw.KeyG", keyT),
				"GLFW_KEY_H":             types.NewIdentGo("GLFW_KEY_H", "glfw.KeyH", keyT),
				"GLFW_KEY_I":             types.NewIdentGo("GLFW_KEY_I", "glfw.KeyI", keyT),
				"GLFW_KEY_J":             types.NewIdentGo("GLFW_KEY_J", "glfw.KeyJ", keyT),
				"GLFW_KEY_K":             types.NewIdentGo("GLFW_KEY_K", "glfw.KeyK", keyT),
				"GLFW_KEY_L":             types.NewIdentGo("GLFW_KEY_L", "glfw.KeyL", keyT),
				"GLFW_KEY_M":             types.NewIdentGo("GLFW_KEY_M", "glfw.KeyM", keyT),
				"GLFW_KEY_N":             types.NewIdentGo("GLFW_KEY_N", "glfw.KeyN", keyT),
				"GLFW_KEY_O":             types.NewIdentGo("GLFW_KEY_O", "glfw.KeyO", keyT),
				"GLFW_KEY_P":             types.NewIdentGo("GLFW_KEY_P", "glfw.KeyP", keyT),
				"GLFW_KEY_Q":             types.NewIdentGo("GLFW_KEY_Q", "glfw.KeyQ", keyT),
				"GLFW_KEY_R":             types.NewIdentGo("GLFW_KEY_R", "glfw.KeyR", keyT),
				"GLFW_KEY_S":             types.NewIdentGo("GLFW_KEY_S", "glfw.KeyS", keyT),
				"GLFW_KEY_T":             types.NewIdentGo("GLFW_KEY_T", "glfw.KeyT", keyT),
				"GLFW_KEY_U":             types.NewIdentGo("GLFW_KEY_U", "glfw.KeyU", keyT),
				"GLFW_KEY_V":             types.NewIdentGo("GLFW_KEY_V", "glfw.KeyV", keyT),
				"GLFW_KEY_W":             types.NewIdentGo("GLFW_KEY_W", "glfw.KeyW", keyT),
				"GLFW_KEY_X":             types.NewIdentGo("GLFW_KEY_X", "glfw.KeyX", keyT),
				"GLFW_KEY_Y":             types.NewIdentGo("GLFW_KEY_Y", "glfw.KeyY", keyT),
				"GLFW_KEY_Z":             types.NewIdentGo("GLFW_KEY_Z", "glfw.KeyZ", keyT),
				"GLFW_KEY_LEFT_BRACKET":  types.NewIdentGo("GLFW_KEY_LEFT_BRACKET", "glfw.KeyLeftBracket", keyT),
				"GLFW_KEY_BACKSLASH":     types.NewIdentGo("GLFW_KEY_BACKSLASH", "glfw.KeyBackslash", keyT),
				"GLFW_KEY_RIGHT_BRACKET": types.NewIdentGo("GLFW_KEY_RIGHT_BRACKET", "glfw.KeyRightBracket", keyT),
				"GLFW_KEY_GRAVE_ACCENT":  types.NewIdentGo("GLFW_KEY_GRAVE_ACCENT", "glfw.KeyGraveAccent", keyT),
				"GLFW_KEY_WORLD_1":       types.NewIdentGo("GLFW_KEY_WORLD_1", "glfw.KeyWorld1", keyT),
				"GLFW_KEY_WORLD_2":       types.NewIdentGo("GLFW_KEY_WORLD_2", "glfw.KeyWorld2", keyT),
				// function key constants
				"GLFW_KEY_ESCAPE":        types.NewIdentGo("GLFW_KEY_ESCAPE", "glfw.KeyEscape", keyT),
				"GLFW_KEY_ENTER":         types.NewIdentGo("GLFW_KEY_ENTER", "glfw.KeyEnter", keyT),
				"GLFW_KEY_TAB":           types.NewIdentGo("GLFW_KEY_TAB", "glfw.KeyTab", keyT),
				"GLFW_KEY_BACKSPACE":     types.NewIdentGo("GLFW_KEY_BACKSPACE", "glfw.KeyBackspace", keyT),
				"GLFW_KEY_INSERT":        types.NewIdentGo("GLFW_KEY_INSERT", "glfw.KeyInsert", keyT),
				"GLFW_KEY_DELETE":        types.NewIdentGo("GLFW_KEY_DELETE", "glfw.KeyDelete", keyT),
				"GLFW_KEY_RIGHT":         types.NewIdentGo("GLFW_KEY_RIGHT", "glfw.KeyRight", keyT),
				"GLFW_KEY_LEFT":          types.NewIdentGo("GLFW_KEY_LEFT", "glfw.KeyLeft", keyT),
				"GLFW_KEY_DOWN":          types.NewIdentGo("GLFW_KEY_DOWN", "glfw.KeyDown", keyT),
				"GLFW_KEY_UP":            types.NewIdentGo("GLFW_KEY_UP", "glfw.KeyUp", keyT),
				"GLFW_KEY_PAGE_UP":       types.NewIdentGo("GLFW_KEY_PAGE_UP", "glfw.KeyPageUp", keyT),
				"GLFW_KEY_PAGE_DOWN":     types.NewIdentGo("GLFW_KEY_PAGE_DOWN", "glfw.KeyPageDown", keyT),
				"GLFW_KEY_HOME":          types.NewIdentGo("GLFW_KEY_HOME", "glfw.KeyHome", keyT),
				"GLFW_KEY_END":           types.NewIdentGo("GLFW_KEY_END", "glfw.KeyEnd", keyT),
				"GLFW_KEY_CAPS_LOCK":     types.NewIdentGo("GLFW_KEY_CAPS_LOCK", "glfw.KeyCapsLock", keyT),
				"GLFW_KEY_SCROLL_LOCK":   types.NewIdentGo("GLFW_KEY_SCROLL_LOCK", "glfw.KeyScrollLock", keyT),
				"GLFW_KEY_NUM_LOCK":      types.NewIdentGo("GLFW_KEY_NUM_LOCK", "glfw.KeyNumLock", keyT),
				"GLFW_KEY_PRINT_SCREEN":  types.NewIdentGo("GLFW_KEY_PRINT_SCREEN", "glfw.KeyPrintScreen", keyT),
				"GLFW_KEY_PAUSE":         types.NewIdentGo("GLFW_KEY_PAUSE", "glfw.KeyPause", keyT),
				"GLFW_KEY_F1":            types.NewIdentGo("GLFW_KEY_F1", "glfw.KeyF1", keyT),
				"GLFW_KEY_F2":            types.NewIdentGo("GLFW_KEY_F2", "glfw.KeyF2", keyT),
				"GLFW_KEY_F3":            types.NewIdentGo("GLFW_KEY_F3", "glfw.KeyF3", keyT),
				"GLFW_KEY_F4":            types.NewIdentGo("GLFW_KEY_F4", "glfw.KeyF4", keyT),
				"GLFW_KEY_F5":            types.NewIdentGo("GLFW_KEY_F5", "glfw.KeyF5", keyT),
				"GLFW_KEY_F6":            types.NewIdentGo("GLFW_KEY_F6", "glfw.KeyF6", keyT),
				"GLFW_KEY_F7":            types.NewIdentGo("GLFW_KEY_F7", "glfw.KeyF7", keyT),
				"GLFW_KEY_F8":            types.NewIdentGo("GLFW_KEY_F8", "glfw.KeyF8", keyT),
				"GLFW_KEY_F9":            types.NewIdentGo("GLFW_KEY_F9", "glfw.KeyF9", keyT),
				"GLFW_KEY_F10":           types.NewIdentGo("GLFW_KEY_F10", "glfw.KeyF10", keyT),
				"GLFW_KEY_F11":           types.NewIdentGo("GLFW_KEY_F11", "glfw.KeyF11", keyT),
				"GLFW_KEY_F12":           types.NewIdentGo("GLFW_KEY_F12", "glfw.KeyF12", keyT),
				"GLFW_KEY_F13":           types.NewIdentGo("GLFW_KEY_F13", "glfw.KeyF13", keyT),
				"GLFW_KEY_F14":           types.NewIdentGo("GLFW_KEY_F14", "glfw.KeyF14", keyT),
				"GLFW_KEY_F15":           types.NewIdentGo("GLFW_KEY_F15", "glfw.KeyF15", keyT),
				"GLFW_KEY_F16":           types.NewIdentGo("GLFW_KEY_F16", "glfw.KeyF16", keyT),
				"GLFW_KEY_F17":           types.NewIdentGo("GLFW_KEY_F17", "glfw.KeyF17", keyT),
				"GLFW_KEY_F18":           types.NewIdentGo("GLFW_KEY_F18", "glfw.KeyF18", keyT),
				"GLFW_KEY_F19":           types.NewIdentGo("GLFW_KEY_F19", "glfw.KeyF19", keyT),
				"GLFW_KEY_F20":           types.NewIdentGo("GLFW_KEY_F20", "glfw.KeyF20", keyT),
				"GLFW_KEY_F21":           types.NewIdentGo("GLFW_KEY_F21", "glfw.KeyF21", keyT),
				"GLFW_KEY_F22":           types.NewIdentGo("GLFW_KEY_F22", "glfw.KeyF22", keyT),
				"GLFW_KEY_F23":           types.NewIdentGo("GLFW_KEY_F23", "glfw.KeyF23", keyT),
				"GLFW_KEY_F24":           types.NewIdentGo("GLFW_KEY_F24", "glfw.KeyF24", keyT),
				"GLFW_KEY_F25":           types.NewIdentGo("GLFW_KEY_F25", "glfw.KeyF25", keyT),
				"GLFW_KEY_KP_0":          types.NewIdentGo("GLFW_KEY_KP_0", "glfw.KeyKp0", keyT),
				"GLFW_KEY_KP_1":          types.NewIdentGo("GLFW_KEY_KP_1", "glfw.KeyKp1", keyT),
				"GLFW_KEY_KP_2":          types.NewIdentGo("GLFW_KEY_KP_2", "glfw.KeyKp2", keyT),
				"GLFW_KEY_KP_3":          types.NewIdentGo("GLFW_KEY_KP_3", "glfw.KeyKp3", keyT),
				"GLFW_KEY_KP_4":          types.NewIdentGo("GLFW_KEY_KP_4", "glfw.KeyKp4", keyT),
				"GLFW_KEY_KP_5":          types.NewIdentGo("GLFW_KEY_KP_5", "glfw.KeyKp5", keyT),
				"GLFW_KEY_KP_6":          types.NewIdentGo("GLFW_KEY_KP_6", "glfw.KeyKp6", keyT),
				"GLFW_KEY_KP_7":          types.NewIdentGo("GLFW_KEY_KP_7", "glfw.KeyKp7", keyT),
				"GLFW_KEY_KP_8":          types.NewIdentGo("GLFW_KEY_KP_8", "glfw.KeyKp8", keyT),
				"GLFW_KEY_KP_9":          types.NewIdentGo("GLFW_KEY_KP_9", "glfw.KeyKp9", keyT),
				"GLFW_KEY_KP_DECIMAL":    types.NewIdentGo("GLFW_KEY_KP_DECIMAL", "glfw.KeyKpDecimal", keyT),
				"GLFW_KEY_KP_DIVIDE":     types.NewIdentGo("GLFW_KEY_KP_DIVIDE", "glfw.KeyKpDivide", keyT),
				"GLFW_KEY_KP_MULTIPLY":   types.NewIdentGo("GLFW_KEY_KP_MULTIPLY", "glfw.KeyKpMultiply", keyT),
				"GLFW_KEY_KP_SUBTRACT":   types.NewIdentGo("GLFW_KEY_KP_SUBTRACT", "glfw.KeyKpSubtract", keyT),
				"GLFW_KEY_KP_ADD":        types.NewIdentGo("GLFW_KEY_KP_ADD", "glfw.KeyKpAdd", keyT),
				"GLFW_KEY_KP_ENTER":      types.NewIdentGo("GLFW_KEY_KP_ENTER", "glfw.KeyKpEnter", keyT),
				"GLFW_KEY_KP_EQUAL":      types.NewIdentGo("GLFW_KEY_KP_EQUAL", "glfw.KeyKpEqual", keyT),
				"GLFW_KEY_LEFT_SHIFT":    types.NewIdentGo("GLFW_KEY_LEFT_SHIFT", "glfw.KeyLeftShift", keyT),
				"GLFW_KEY_LEFT_CONTROL":  types.NewIdentGo("GLFW_KEY_LEFT_CONTROL", "glfw.KeyLeftControl", keyT),
				"GLFW_KEY_LEFT_ALT":      types.NewIdentGo("GLFW_KEY_LEFT_ALT", "glfw.KeyLeftAlt", keyT),
				"GLFW_KEY_LEFT_SUPER":    types.NewIdentGo("GLFW_KEY_LEFT_SUPER", "glfw.KeyLeftSuper", keyT),
				"GLFW_KEY_RIGHT_SHIFT":   types.NewIdentGo("GLFW_KEY_RIGHT_SHIFT", "glfw.KeyRightShift", keyT),
				"GLFW_KEY_RIGHT_CONTROL": types.NewIdentGo("GLFW_KEY_RIGHT_CONTROL", "glfw.KeyRightControl", keyT),
				"GLFW_KEY_RIGHT_ALT":     types.NewIdentGo("GLFW_KEY_RIGHT_ALT", "glfw.KeyRightAlt", keyT),
				"GLFW_KEY_RIGHT_SUPER":   types.NewIdentGo("GLFW_KEY_RIGHT_SUPER", "glfw.KeyRightSuper", keyT),
				"GLFW_KEY_MENU":          types.NewIdentGo("GLFW_KEY_MENU", "glfw.KeyMenu", keyT),
				"GLFW_KEY_LAST":          types.NewIdentGo("GLFW_KEY_LAST", "glfw.KeyLast", keyT),
				// mouse buttons
				"GLFW_MOUSE_BUTTON_1":      types.NewIdentGo("GLFW_MOUSE_BUTTON_1", "glfw.MouseButton1", mouseButtonT),
				"GLFW_MOUSE_BUTTON_2":      types.NewIdentGo("GLFW_MOUSE_BUTTON_2", "glfw.MouseButton2", mouseButtonT),
				"GLFW_MOUSE_BUTTON_3":      types.NewIdentGo("GLFW_MOUSE_BUTTON_3", "glfw.MouseButton3", mouseButtonT),
				"GLFW_MOUSE_BUTTON_4":      types.NewIdentGo("GLFW_MOUSE_BUTTON_4", "glfw.MouseButton4", mouseButtonT),
				"GLFW_MOUSE_BUTTON_5":      types.NewIdentGo("GLFW_MOUSE_BUTTON_5", "glfw.MouseButton5", mouseButtonT),
				"GLFW_MOUSE_BUTTON_6":      types.NewIdentGo("GLFW_MOUSE_BUTTON_6", "glfw.MouseButton6", mouseButtonT),
				"GLFW_MOUSE_BUTTON_7":      types.NewIdentGo("GLFW_MOUSE_BUTTON_7", "glfw.MouseButton7", mouseButtonT),
				"GLFW_MOUSE_BUTTON_8":      types.NewIdentGo("GLFW_MOUSE_BUTTON_8", "glfw.MouseButton8", mouseButtonT),
				"GLFW_MOUSE_BUTTON_LAST":   types.NewIdentGo("GLFW_MOUSE_BUTTON_LAST", "glfw.MouseButtonLast", mouseButtonT),
				"GLFW_MOUSE_BUTTON_LEFT":   types.NewIdentGo("GLFW_MOUSE_BUTTON_LEFT", "glfw.MouseButtonLeft", mouseButtonT),
				"GLFW_MOUSE_BUTTON_RIGHT":  types.NewIdentGo("GLFW_MOUSE_BUTTON_RIGHT", "glfw.MouseButtonRight", mouseButtonT),
				"GLFW_MOUSE_BUTTON_MIDDLE": types.NewIdentGo("GLFW_MOUSE_BUTTON_MIDDLE", "glfw.MouseButtonMiddle", mouseButtonT),
				// modifier keys
				"GLFW_MOD_SHIFT":     types.NewIdentGo("GLFW_MOD_SHIFT", "glfw.ModShift", modKeyT),
				"GLFW_MOD_CONTROL":   types.NewIdentGo("GLFW_MOD_CONTROL", "glfw.ModControl", modKeyT),
				"GLFW_MOD_ALT":       types.NewIdentGo("GLFW_MOD_ALT", "glfw.ModAlt", modKeyT),
				"GLFW_MOD_SUPER":     types.NewIdentGo("GLFW_MOD_SUPER", "glfw.ModSuper", modKeyT),
				"GLFW_MOD_CAPS_LOCK": types.NewIdentGo("GLFW_MOD_CAPS_LOCK", "glfw.ModCapsLock", modKeyT),
				"GLFW_MOD_NUM_LOCK":  types.NewIdentGo("GLFW_MOD_NUM_LOCK", "glfw.ModNumLock", modKeyT),
				// input mode
				"GLFW_CURSOR":               types.NewIdentGo("GLFW_CURSOR", "glfw.Cursor", inputModeT),
				"GLFW_STICKY_KEYS":          types.NewIdentGo("GLFW_STICKY_KEYS", "glfw.StickyKeys", inputModeT),
				"GLFW_STICKY_MOUSE_BUTTONS": types.NewIdentGo("GLFW_STICKY_MOUSE_BUTTONS", "glfw.StickyMouseButtons", inputModeT),
				"GLFW_LOCK_KEY_MODS":        types.NewIdentGo("GLFW_LOCK_KEY_MODS", "glfw.LockKeyMods", inputModeT),
				"GLFW_RAW_MOUSE_MOTION":     types.NewIdentGo("GLFW_RAW_MOUSE_MOTION", "glfw.RawMouseMotion", inputModeT),
				// peripheral events
				"GLFW_CONNECTED":    types.NewIdentGo("GLFW_CONNECTED", "glfw.Connected", perEventT),
				"GLFW_DISCONNECTED": types.NewIdentGo("GLFW_DISCONNECTED", "glfw.Disconnected", perEventT),
				// hints
				"GLFW_FOCUSED":                  types.NewIdentGo("GLFW_FOCUSED", "glfw.Focused", hintT),
				"GLFW_ICONIFIED":                types.NewIdentGo("GLFW_ICONIFIED", "glfw.Iconified", hintT),
				"GLFW_RESIZABLE":                types.NewIdentGo("GLFW_RESIZABLE", "glfw.Resizable", hintT),
				"GLFW_VISIBLE":                  types.NewIdentGo("GLFW_VISIBLE", "glfw.Visible", hintT),
				"GLFW_DECORATED":                types.NewIdentGo("GLFW_DECORATED", "glfw.Decorated", hintT),
				"GLFW_AUTO_ICONIFY":             types.NewIdentGo("GLFW_AUTO_ICONIFY", "glfw.AutoIconify", hintT),
				"GLFW_FLOATING":                 types.NewIdentGo("GLFW_FLOATING", "glfw.Floating", hintT),
				"GLFW_MAXIMIZED":                types.NewIdentGo("GLFW_MAXIMIZED", "glfw.Maximized", hintT),
				"GLFW_CENTER_CURSOR":            types.NewIdentGo("GLFW_CENTER_CURSOR", "glfw.CenterCursor", hintT),
				"GLFW_TRANSPARENT_FRAMEBUFFER":  types.NewIdentGo("GLFW_TRANSPARENT_FRAMEBUFFER", "glfw.TransparentFramebuffer", hintT),
				"GLFW_HOVERED":                  types.NewIdentGo("GLFW_HOVERED", "glfw.Hovered", hintT),
				"GLFW_FOCUS_ON_SHOW":            types.NewIdentGo("GLFW_FOCUS_ON_SHOW", "glfw.FocusOnShow", hintT),
				"GLFW_MOUSE_PASSTHROUGH":        types.NewIdentGo("GLFW_MOUSE_PASSTHROUGH", "glfw.MousePassthrough", hintT),
				"GLFW_RED_BITS":                 types.NewIdentGo("GLFW_RED_BITS", "glfw.RedBits", hintT),
				"GLFW_GREEN_BITS":               types.NewIdentGo("GLFW_GREEN_BITS", "glfw.GreenBits", hintT),
				"GLFW_BLUE_BITS":                types.NewIdentGo("GLFW_BLUE_BITS", "glfw.BlueBits", hintT),
				"GLFW_ALPHA_BITS":               types.NewIdentGo("GLFW_ALPHA_BITS", "glfw.AlphaBits", hintT),
				"GLFW_DEPTH_BITS":               types.NewIdentGo("GLFW_DEPTH_BITS", "glfw.DepthBits", hintT),
				"GLFW_STENCIL_BITS":             types.NewIdentGo("GLFW_STENCIL_BITS", "glfw.StencilBits", hintT),
				"GLFW_ACCUM_RED_BITS":           types.NewIdentGo("GLFW_ACCUM_RED_BITS", "glfw.AccumRedBits", hintT),
				"GLFW_ACCUM_GREEN_BITS":         types.NewIdentGo("GLFW_ACCUM_GREEN_BITS", "glfw.AccumGreenBits", hintT),
				"GLFW_ACCUM_BLUE_BITS":          types.NewIdentGo("GLFW_ACCUM_BLUE_BITS", "glfw.AccumBlueBits", hintT),
				"GLFW_ACCUM_ALPHA_BITS":         types.NewIdentGo("GLFW_ACCUM_ALPHA_BITS", "glfw.AccumAlphaBits", hintT),
				"GLFW_AUX_BUFFERS":              types.NewIdentGo("GLFW_AUX_BUFFERS", "glfw.AuxBuffers", hintT),
				"GLFW_STEREO":                   types.NewIdentGo("GLFW_STEREO", "glfw.Stereo", hintT),
				"GLFW_SAMPLES":                  types.NewIdentGo("GLFW_SAMPLES", "glfw.Samples", hintT),
				"GLFW_SRGB_CAPABLE":             types.NewIdentGo("GLFW_SRGB_CAPABLE", "glfw.SrgbCapable", hintT),
				"GLFW_REFRESH_RATE":             types.NewIdentGo("GLFW_REFRESH_RATE", "glfw.RefreshRate", hintT),
				"GLFW_DOUBLEBUFFER":             types.NewIdentGo("GLFW_DOUBLEBUFFER", "glfw.Doublebuffer", hintT),
				"GLFW_CLIENT_API":               types.NewIdentGo("GLFW_CLIENT_API", "glfw.ClientApi", hintT),
				"GLFW_CONTEXT_VERSION_MAJOR":    types.NewIdentGo("GLFW_CONTEXT_VERSION_MAJOR", "glfw.ContextVersionMajor", hintT),
				"GLFW_CONTEXT_VERSION_MINOR":    types.NewIdentGo("GLFW_CONTEXT_VERSION_MINOR", "glfw.ContextVersionMinor", hintT),
				"GLFW_CONTEXT_REVISION":         types.NewIdentGo("GLFW_CONTEXT_REVISION", "glfw.ContextRevision", hintT),
				"GLFW_CONTEXT_ROBUSTNESS":       types.NewIdentGo("GLFW_CONTEXT_ROBUSTNESS", "glfw.ContextRobustness", hintT),
				"GLFW_OPENGL_FORWARD_COMPAT":    types.NewIdentGo("GLFW_OPENGL_FORWARD_COMPAT", "glfw.OpenglForwardCompat", hintT),
				"GLFW_CONTEXT_DEBUG":            types.NewIdentGo("GLFW_CONTEXT_DEBUG", "glfw.ContextDebug", hintT),
				"GLFW_OPENGL_DEBUG_CONTEXT":     types.NewIdentGo("GLFW_OPENGL_DEBUG_CONTEXT", "glfw.OpenglDebugContext", hintT),
				"GLFW_OPENGL_PROFILE":           types.NewIdentGo("GLFW_OPENGL_PROFILE", "glfw.OpenglProfile", hintT),
				"GLFW_CONTEXT_RELEASE_BEHAVIOR": types.NewIdentGo("GLFW_CONTEXT_RELEASE_BEHAVIOR", "glfw.ContextReleaseBehavior", hintT),
				"GLFW_CONTEXT_NO_ERROR":         types.NewIdentGo("GLFW_CONTEXT_NO_ERROR", "glfw.ContextNoError", hintT),
				"GLFW_CONTEXT_CREATION_API":     types.NewIdentGo("GLFW_CONTEXT_CREATION_API", "glfw.ContextCreationApi", hintT),
				"GLFW_SCALE_TO_MONITOR":         types.NewIdentGo("GLFW_SCALE_TO_MONITOR", "glfw.ScaleToMonitor", hintT),
				"GLFW_COCOA_RETINA_FRAMEBUFFER": types.NewIdentGo("GLFW_COCOA_RETINA_FRAMEBUFFER", "glfw.CocoaRetinaFramebuffer", hintT),
				"GLFW_COCOA_FRAME_NAME":         types.NewIdentGo("GLFW_COCOA_FRAME_NAME", "glfw.CocoaFrameName", hintT),
				"GLFW_COCOA_GRAPHICS_SWITCHING": types.NewIdentGo("GLFW_COCOA_GRAPHICS_SWITCHING", "glfw.CocoaGraphicsSwitching", hintT),
				"GLFW_X11_CLASS_NAME":           types.NewIdentGo("GLFW_X11_CLASS_NAME", "glfw.X11ClassName", hintT),
				"GLFW_X11_INSTANCE_NAME":        types.NewIdentGo("GLFW_X11_INSTANCE_NAME", "glfw.X11InstanceName", hintT),
				"GLFW_WIN32_KEYBOARD_MENU":      types.NewIdentGo("GLFW_WIN32_KEYBOARD_MENU", "glfw.Win32KeyboardMenu", hintT),
			},
			Header: `
#include <` + BuiltinH + `>
#define GLFW_TRUE                   1
#define GLFW_FALSE                  0

// key and button actions
#define GLFW_RELEASE                0
#define GLFW_PRESS                  1
#define GLFW_REPEAT                 2

// hints
#define GLFW_FOCUSED                0x00020001
#define GLFW_ICONIFIED              0x00020002
#define GLFW_RESIZABLE              0x00020003
#define GLFW_VISIBLE                0x00020004
#define GLFW_DECORATED              0x00020005
#define GLFW_AUTO_ICONIFY           0x00020006
#define GLFW_FLOATING               0x00020007
#define GLFW_MAXIMIZED              0x00020008
#define GLFW_CENTER_CURSOR          0x00020009
#define GLFW_TRANSPARENT_FRAMEBUFFER 0x0002000A
#define GLFW_HOVERED                0x0002000B
#define GLFW_FOCUS_ON_SHOW          0x0002000C
#define GLFW_MOUSE_PASSTHROUGH      0x0002000D
#define GLFW_RED_BITS               0x00021001
#define GLFW_GREEN_BITS             0x00021002
#define GLFW_BLUE_BITS              0x00021003
#define GLFW_ALPHA_BITS             0x00021004
#define GLFW_DEPTH_BITS             0x00021005
#define GLFW_STENCIL_BITS           0x00021006
#define GLFW_ACCUM_RED_BITS         0x00021007
#define GLFW_ACCUM_GREEN_BITS       0x00021008
#define GLFW_ACCUM_BLUE_BITS        0x00021009
#define GLFW_ACCUM_ALPHA_BITS       0x0002100A
#define GLFW_AUX_BUFFERS            0x0002100B
#define GLFW_STEREO                 0x0002100C
#define GLFW_SAMPLES                0x0002100D
#define GLFW_SRGB_CAPABLE           0x0002100E
#define GLFW_REFRESH_RATE           0x0002100F
#define GLFW_DOUBLEBUFFER           0x00021010
#define GLFW_CLIENT_API             0x00022001
#define GLFW_CONTEXT_VERSION_MAJOR  0x00022002
#define GLFW_CONTEXT_VERSION_MINOR  0x00022003
#define GLFW_CONTEXT_REVISION       0x00022004
#define GLFW_CONTEXT_ROBUSTNESS     0x00022005
#define GLFW_OPENGL_FORWARD_COMPAT  0x00022006
#define GLFW_CONTEXT_DEBUG          0x00022007
#define GLFW_OPENGL_DEBUG_CONTEXT   GLFW_CONTEXT_DEBUG
#define GLFW_OPENGL_PROFILE         0x00022008
#define GLFW_CONTEXT_RELEASE_BEHAVIOR 0x00022009
#define GLFW_CONTEXT_NO_ERROR       0x0002200A
#define GLFW_CONTEXT_CREATION_API   0x0002200B
#define GLFW_SCALE_TO_MONITOR       0x0002200C
#define GLFW_COCOA_RETINA_FRAMEBUFFER 0x00023001
#define GLFW_COCOA_FRAME_NAME         0x00023002
#define GLFW_COCOA_GRAPHICS_SWITCHING 0x00023003
#define GLFW_X11_CLASS_NAME         0x00024001
#define GLFW_X11_INSTANCE_NAME      0x00024002
#define GLFW_WIN32_KEYBOARD_MENU    0x00025001

/* The unknown key */
#define GLFW_KEY_UNKNOWN            -1

/* Printable keys */
#define GLFW_KEY_SPACE              32
#define GLFW_KEY_APOSTROPHE         39  /* ' */
#define GLFW_KEY_COMMA              44  /* , */
#define GLFW_KEY_MINUS              45  /* - */
#define GLFW_KEY_PERIOD             46  /* . */
#define GLFW_KEY_SLASH              47  /* / */
#define GLFW_KEY_0                  48
#define GLFW_KEY_1                  49
#define GLFW_KEY_2                  50
#define GLFW_KEY_3                  51
#define GLFW_KEY_4                  52
#define GLFW_KEY_5                  53
#define GLFW_KEY_6                  54
#define GLFW_KEY_7                  55
#define GLFW_KEY_8                  56
#define GLFW_KEY_9                  57
#define GLFW_KEY_SEMICOLON          59  /* ; */
#define GLFW_KEY_EQUAL              61  /* = */
#define GLFW_KEY_A                  65
#define GLFW_KEY_B                  66
#define GLFW_KEY_C                  67
#define GLFW_KEY_D                  68
#define GLFW_KEY_E                  69
#define GLFW_KEY_F                  70
#define GLFW_KEY_G                  71
#define GLFW_KEY_H                  72
#define GLFW_KEY_I                  73
#define GLFW_KEY_J                  74
#define GLFW_KEY_K                  75
#define GLFW_KEY_L                  76
#define GLFW_KEY_M                  77
#define GLFW_KEY_N                  78
#define GLFW_KEY_O                  79
#define GLFW_KEY_P                  80
#define GLFW_KEY_Q                  81
#define GLFW_KEY_R                  82
#define GLFW_KEY_S                  83
#define GLFW_KEY_T                  84
#define GLFW_KEY_U                  85
#define GLFW_KEY_V                  86
#define GLFW_KEY_W                  87
#define GLFW_KEY_X                  88
#define GLFW_KEY_Y                  89
#define GLFW_KEY_Z                  90
#define GLFW_KEY_LEFT_BRACKET       91  /* [ */
#define GLFW_KEY_BACKSLASH          92  /* \ */
#define GLFW_KEY_RIGHT_BRACKET      93  /* ] */
#define GLFW_KEY_GRAVE_ACCENT       96  /*  */
#define GLFW_KEY_WORLD_1            161 /* non-US #1 */
#define GLFW_KEY_WORLD_2            162 /* non-US #2 */

/* Function keys */
#define GLFW_KEY_ESCAPE             256
#define GLFW_KEY_ENTER              257
#define GLFW_KEY_TAB                258
#define GLFW_KEY_BACKSPACE          259
#define GLFW_KEY_INSERT             260
#define GLFW_KEY_DELETE             261
#define GLFW_KEY_RIGHT              262
#define GLFW_KEY_LEFT               263
#define GLFW_KEY_DOWN               264
#define GLFW_KEY_UP                 265
#define GLFW_KEY_PAGE_UP            266
#define GLFW_KEY_PAGE_DOWN          267
#define GLFW_KEY_HOME               268
#define GLFW_KEY_END                269
#define GLFW_KEY_CAPS_LOCK          280
#define GLFW_KEY_SCROLL_LOCK        281
#define GLFW_KEY_NUM_LOCK           282
#define GLFW_KEY_PRINT_SCREEN       283
#define GLFW_KEY_PAUSE              284
#define GLFW_KEY_F1                 290
#define GLFW_KEY_F2                 291
#define GLFW_KEY_F3                 292
#define GLFW_KEY_F4                 293
#define GLFW_KEY_F5                 294
#define GLFW_KEY_F6                 295
#define GLFW_KEY_F7                 296
#define GLFW_KEY_F8                 297
#define GLFW_KEY_F9                 298
#define GLFW_KEY_F10                299
#define GLFW_KEY_F11                300
#define GLFW_KEY_F12                301
#define GLFW_KEY_F13                302
#define GLFW_KEY_F14                303
#define GLFW_KEY_F15                304
#define GLFW_KEY_F16                305
#define GLFW_KEY_F17                306
#define GLFW_KEY_F18                307
#define GLFW_KEY_F19                308
#define GLFW_KEY_F20                309
#define GLFW_KEY_F21                310
#define GLFW_KEY_F22                311
#define GLFW_KEY_F23                312
#define GLFW_KEY_F24                313
#define GLFW_KEY_F25                314
#define GLFW_KEY_KP_0               320
#define GLFW_KEY_KP_1               321
#define GLFW_KEY_KP_2               322
#define GLFW_KEY_KP_3               323
#define GLFW_KEY_KP_4               324
#define GLFW_KEY_KP_5               325
#define GLFW_KEY_KP_6               326
#define GLFW_KEY_KP_7               327
#define GLFW_KEY_KP_8               328
#define GLFW_KEY_KP_9               329
#define GLFW_KEY_KP_DECIMAL         330
#define GLFW_KEY_KP_DIVIDE          331
#define GLFW_KEY_KP_MULTIPLY        332
#define GLFW_KEY_KP_SUBTRACT        333
#define GLFW_KEY_KP_ADD             334
#define GLFW_KEY_KP_ENTER           335
#define GLFW_KEY_KP_EQUAL           336
#define GLFW_KEY_LEFT_SHIFT         340
#define GLFW_KEY_LEFT_CONTROL       341
#define GLFW_KEY_LEFT_ALT           342
#define GLFW_KEY_LEFT_SUPER         343
#define GLFW_KEY_RIGHT_SHIFT        344
#define GLFW_KEY_RIGHT_CONTROL      345
#define GLFW_KEY_RIGHT_ALT          346
#define GLFW_KEY_RIGHT_SUPER        347
#define GLFW_KEY_MENU               348

#define GLFW_KEY_LAST               GLFW_KEY_MENU

// mouse buttons
#define GLFW_MOUSE_BUTTON_1         0
#define GLFW_MOUSE_BUTTON_2         1
#define GLFW_MOUSE_BUTTON_3         2
#define GLFW_MOUSE_BUTTON_4         3
#define GLFW_MOUSE_BUTTON_5         4
#define GLFW_MOUSE_BUTTON_6         5
#define GLFW_MOUSE_BUTTON_7         6
#define GLFW_MOUSE_BUTTON_8         7
#define GLFW_MOUSE_BUTTON_LAST      GLFW_MOUSE_BUTTON_8
#define GLFW_MOUSE_BUTTON_LEFT      GLFW_MOUSE_BUTTON_1
#define GLFW_MOUSE_BUTTON_RIGHT     GLFW_MOUSE_BUTTON_2
#define GLFW_MOUSE_BUTTON_MIDDLE    GLFW_MOUSE_BUTTON_3

// shift keys
#define GLFW_MOD_SHIFT           0x0001
#define GLFW_MOD_CONTROL         0x0002
#define GLFW_MOD_ALT             0x0004
#define GLFW_MOD_SUPER           0x0008
#define GLFW_MOD_CAPS_LOCK       0x0010
#define GLFW_MOD_NUM_LOCK        0x0020

// input modes
#define GLFW_CURSOR                 0x00033001
#define GLFW_STICKY_KEYS            0x00033002
#define GLFW_STICKY_MOUSE_BUTTONS   0x00033003
#define GLFW_LOCK_KEY_MODS          0x00033004
#define GLFW_RAW_MOUSE_MOTION       0x00033005

// peripheral events
#define GLFW_CONNECTED              0x00040001
#define GLFW_DISCONNECTED           0x00040002

typedef struct GLFWwindow GLFWwindow;
typedef void (* GLFWerrorfun)(int,const char*);
typedef void (* GLFWkeyfun)(GLFWwindow*,int,int,int,int);
typedef void (* GLFWcharfun)(GLFWwindow*,unsigned int);
typedef void (* GLFWframebuffersizefun)(GLFWwindow*,int,int);

struct GLFWwindow {
	void (*MakeContextCurrent)(void);
	_Bool (*ShouldClose)(void);
	void (*SwapBuffers)(void);
	int (*GetKey)(int);
	void (*SetShouldClose)(_Bool);
	void (*GetFramebufferSize)(int* width, int* height);
	void (*Destroy)(void);
	void (*Focus)(void);
	void (*Maximize)(void);
	void (*Show)(void);
	void (*Hide)(void);
	void (*Iconify)(void);
	void (*Restore)(void);
	void (*SetTitle)(const char* title);
	void (*SetSize)(int width, int height);
	void (*SetPos)(int x, int y);
	int (*GetInputMode)(int);
	void (*SetInputMode)(int, int);

	void* (*GetWindowUserPointer)(void);

	// callbacks
	GLFWkeyfun (*SetKeyCallback)(GLFWkeyfun);
	GLFWcharfun (*SetCharCallback)(GLFWcharfun);
	GLFWframebuffersizefun (*SetFramebufferSizeCallback)(GLFWframebuffersizefun);
};
#define glfwMakeContextCurrent(win) ((GLFWwindow*)win)->MakeContextCurrent()
#define glfwWindowShouldClose(win) ((GLFWwindow*)win)->ShouldClose()
#define glfwSwapBuffers(win) ((GLFWwindow*)win)->SwapBuffers()
#define glfwGetKey(win, k) ((GLFWwindow*)win)->GetKey(k)
#define glfwSetWindowShouldClose(win, b) ((GLFWwindow*)win)->SetShouldClose(b)
#define glfwSetKeyCallback(win, cb) ((GLFWwindow*)win)->SetKeyCallback(cb)
#define glfwSetCharCallback(win, cb) ((GLFWwindow*)win)->SetCharCallback(cb)
#define glfwSetFramebufferSizeCallback(win, cb) ((GLFWwindow*)win)->SetFramebufferSizeCallback(cb)
#define glfwGetFramebufferSize(win, w, h) ((GLFWwindow*)win)->GetFramebufferSize(w, h)
#define glfwDestroyWindow(win) ((GLFWwindow*)win)->Destroy()
#define glfwFocusWindow(win) ((GLFWwindow*)win)->Focus()
#define glfwMaximizeWindow(win) ((GLFWwindow*)win)->Maximize()
#define glfwShowWindow(win) ((GLFWwindow*)win)->Show()
#define glfwHideWindow(win) ((GLFWwindow*)win)->Hide()
#define glfwIconifyWindow(win) ((GLFWwindow*)win)->Iconify()
#define glfwRestoreWindow(win) ((GLFWwindow*)win)->Restore()
#define glfwSetWindowTitle(win, title) ((GLFWwindow*)win)->SetTitle(title)
#define glfwSetWindowSize(win, w, h) ((GLFWwindow*)win)->SetSize(w, h)
#define glfwSetWindowPos(win, x, y) ((GLFWwindow*)win)->SetPos(x, y)
#define glfwGetWindowUserPointer(win) ((GLFWwindow*)win)->GetWindowUserPointer()
#define glfwGetInputMode(win, mode) ((GLFWwindow*)win)->GetInputMode(mode)
#define glfwSetInputMode(win, mode, v) ((GLFWwindow*)win)->SetInputMode(mode, v)

typedef struct GLFWvidmode {
    int width;
    int height;
    int redBits;
    int greenBits;
    int blueBits;
    int refreshRate;
} GLFWvidmode;

typedef struct GLFWmonitor {
	GLFWvidmode* (*GetVideoMode)(void);
	void (*GetPos)(int*, int*);
	void (*GetPhysicalSize)(int*, int*);
	char* (*GetName)(void);
} GLFWmonitor;
#define glfwGetVideoMode(mon) ((GLFWmonitor*)mon)->GetVideoMode()
#define glfwGetMonitorPos(mon, x, y) ((GLFWmonitor*)mon)->GetPos(x, y)
#define glfwGetMonitorPhysicalSize(mon, x, y) ((GLFWmonitor*)mon)->GetPhysicalSize(x, y)
#define glfwGetMonitorName(mon) ((GLFWmonitor*)mon)->GetName()

typedef struct _GLFWjoystick {
	const float* (*GetAxes)(int* count);
	const unsigned char* (*GetButtons)(int* count);
	const unsigned char* (*GetHats)(int* count);
	const char* (*GetName)(void);
	_Bool (*IsGamepad)(void);
	const char* (*GetGUID)(void);
	const char* GetGamepadName(void);
} _GLFWjoystick;
#define glfwGetJoystickAxes(j, a) ((_GLFWjoystick*)j)->GetAxes(a)
#define glfwGetJoystickButtons(j, a) ((_GLFWjoystick*)j)->GetButtons(a)
#define glfwGetJoystickHats(j, a) ((_GLFWjoystick*)j)->GetHats(a)
#define glfwGetJoystickName(j) ((_GLFWjoystick*)j)->GetName()
#define glfwJoystickIsGamepad(j) ((_GLFWjoystick*)j)->IsGamepad()
#define glfwGetJoystickGUID(j) ((_GLFWjoystick*)j)->GetGUID()
#define glfwGetGamepadName(j) ((_GLFWjoystick*)j)->GetGamepadName()

void glfwWindowHint(int, int);
const char* glfwGetKeyName(int key, int scancode);
GLFWerrorfun glfwSetErrorCallback(GLFWerrorfun callback); // no go equivalent
`,
		}
		l.Declare(
			// functions
			types.NewIdentGo("glfwInit", "glfw.Init", env.FuncTT(env.C().Int(), nil)), // returns an error instead of an int
			types.NewIdentGo("glfwTerminate", "glfw.Terminate", env.FuncTT(nil, nil)),
			// createWindow returns an error along with the window
			types.NewIdentGo("glfwCreateWindow", "glfw.CreateWindow", env.FuncTT(env.PtrT(windowT), env.Go().Int(), env.Go().Int(), env.Go().String(), env.PtrT(monitorT), env.PtrT(windowT))),
			types.NewIdentGo("glfwGetProcAddress", "glfw.GetProcAddress", env.FuncTT(env.PtrT(nil), env.Go().String())),
			types.NewIdentGo("glfwPollEvents", "glfw.PollEvents", env.FuncTT(nil, nil)),
			types.NewIdentGo("glfwSwapInterval", "glfw.SwapInterval", env.FuncTT(nil, env.Go().Int())),
			types.NewIdentGo("glfwGetTime", "glfw.GetTime", env.FuncTT(env.C().Float(), nil)),
			types.NewIdentGo("glfwGetCurrentContext", "glfw.GetCurrentContext", env.FuncTT(env.PtrT(windowT), nil)),
			types.NewIdentGo("glfwSetClipboardString", "glfw.SetClipboardString", env.FuncTT(nil, env.Go().String())),
			types.NewIdentGo("glfwWaitEvents", "glfw.WaitEvents", env.FuncTT(nil, nil)),
		)
		return l
	})
}
