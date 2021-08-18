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
		keyT := types.NamedTGo("Key", "glfw.Key", env.Go().Int())
		modKeyT := types.NamedTGo("ModifierKey", "glfw.ModifierKey", env.Go().Int())
		keyCbT := types.NamedTGo("GLFWkeyfun", "glfw.KeyCallback", env.FuncTT(nil, keyT, env.Go().Int(), actionT, modKeyT))
		windowT := types.NamedTGo("GLFWwindow", "glfw.Window", env.MethStructT(map[string]*types.FuncType{
			"MakeContextCurrent": env.FuncTT(nil, nil),
			"ShouldClose":        env.FuncTT(env.Go().Bool(), nil),
			"SwapBuffers":        env.FuncTT(nil, nil),
			"GetKey":             env.FuncTT(actionT, keyT),
			"SetShouldClose":     env.FuncTT(nil, env.Go().Bool()),
			"SetKeyCallback":     env.FuncTT(keyCbT, keyCbT),
			"GetFramebufferSize": env.FuncTT(nil, env.Go().Int(), env.Go().Int()), //FIXME: incorrect signature
			"Destroy":            env.FuncTT(nil, nil),
		}))
		monitorT := types.NamedTGo("GLFWmonitor", "glfw.Monitor", types.StructT(nil))
		l := &Library{
			Imports: map[string]string{
				"glfw": "github.com/go-gl/glfw/v3.3/glfw",
			},
			Types: map[string]types.Type{
				"GLFWwindow": windowT,
			},
			Idents: map[string]*types.Ident{
				// functions
				"glfwWindowHint": types.NewIdentGo("glfwWindowHint", "glfw.WindowHint", env.FuncTT(nil, hintT, env.Go().Int())),
				// constants
				"GLFW_CONTEXT_VERSION_MAJOR": types.NewIdentGo("GLFW_CONTEXT_VERSION_MAJOR", "glfw.ContextVersionMajor", hintT),
				"GLFW_CONTEXT_VERSION_MINOR": types.NewIdentGo("GLFW_CONTEXT_VERSION_MINOR", "glfw.ContextVersionMinor", hintT),
				"GLFW_OPENGL_PROFILE":        types.NewIdentGo("GLFW_OPENGL_PROFILE", "glfw.OpenGLProfile", hintT),
				"GLFW_OPENGL_CORE_PROFILE":   types.NewIdentGo("GLFW_OPENGL_CORE_PROFILE", "glfw.OpenGLCoreProfile", env.Go().Int()),
				"GLFW_OPENGL_FORWARD_COMPAT": types.NewIdentGo("GLFW_OPENGL_FORWARD_COMPAT", "glfw.OpenGLForwardCompatible", hintT),
				"GLFW_TRUE":                  types.NewIdentGo("GLFW_TRUE", "glfw.True", env.Go().Int()),
				"GLFW_PRESS":                 types.NewIdentGo("GLFW_PRESS", "glfw.Press", actionT),
				"GLFW_KEY_ESCAPE":            types.NewIdentGo("GLFW_KEY_ESCAPE", "glfw.KeyEscape", env.Go().Int()),
			},
			Header: `
const int GLFW_CONTEXT_VERSION_MAJOR = 0x00022002;
const int GLFW_CONTEXT_VERSION_MINOR = 0x00022003;
const int GLFW_OPENGL_PROFILE = 0x00022008;
const int GLFW_OPENGL_CORE_PROFILE = 0x00032001;
const int GLFW_OPENGL_FORWARD_COMPAT = 0x00022006;
const int GLFW_TRUE = 1;
const int GLFW_PRESS = 1;

// functions keys
const int GLFW_KEY_ESCAPE = 256;

typedef struct GLFWwindow GLFWwindow;
typedef void (* GLFWerrorfun)(int,const char*);
typedef void (* GLFWkeyfun)(GLFWwindow*,int,int,int,int);
typedef void (* GLFWframebuffersizefun)(GLFWwindow*,int,int);

struct GLFWwindow {
	void (*MakeContextCurrent)(void);
	_Bool (*ShouldClose)(void);
	void (*SwapBuffers)(void);
	int (*GetKey)(int);
	void (*SetShouldClose)(_Bool);
	void (*GetFramebufferSize)(int* width, int* height);
	void (*Destroy)();
	// callbacks
	GLFWframebuffersizefun (*SetKeyCallback)(GLFWframebuffersizefun)
};
#define glfwMakeContextCurrent(win) ((GLFWwindow*)win)->MakeContextCurrent()
#define glfwWindowShouldClose(win) ((GLFWwindow*)win)->ShouldClose()
#define glfwSwapBuffers(win) ((GLFWwindow*)win)->SwapBuffers()
#define glfwGetKey(win, k) ((GLFWwindow*)win)->GetKey(k)
#define glfwSetWindowShouldClose(win, b) ((GLFWwindow*)win)->SetShouldClose(b)
#define glfwSetKeyCallback(win, cb) ((GLFWwindow*)win)->SetKeyCallback(cb)
#define glfwGetFramebufferSize(win, w, h) ((GLFWwindow*)win)->GetFramebufferSize(w, h)
#define glfwDestroyWindow(win) ((GLFWwindow*)win)->Destroy()
typedef struct GLFWmonitor GLFWmonitor;

void glfwWindowHint(int, int);
GLFWerrorfun glfwSetErrorCallback(GLFWerrorfun callback); // no go equivalent
`,
		}
		l.Declare(
			// functions
			types.NewIdentGo("glfwInit", "glfw.Init", env.FuncTT(env.C().Int(), nil)),
			types.NewIdentGo("glfwTerminate", "glfw.Terminate", env.FuncTT(nil, nil)),
			types.NewIdentGo("glfwCreateWindow", "glfw.CreateWindow", env.FuncTT(env.PtrT(windowT), env.Go().Int(), env.Go().Int(), env.Go().String(), env.PtrT(monitorT), env.PtrT(windowT))),
			types.NewIdentGo("glfwGetProcAddress", "glfw.GetProcAddress", env.FuncTT(env.PtrT(nil), env.Go().String())),
			types.NewIdentGo("glfwPollEvents", "glfw.PollEvents", env.FuncTT(nil, nil)),
			types.NewIdentGo("glfwSwapInterval", "glfw.SwapInterval", env.FuncTT(nil, env.Go().Int())),
			types.NewIdentGo("glfwGetTime", "glfw.GetTime", env.FuncTT(env.C().Float(), nil)),
		)
		return l
	})
}
