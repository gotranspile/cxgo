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
#define GLFW_OPENGL_ANY_PROFILE              0
#define GLFW_OPENGL_CORE_PROFILE    0x00032001
#define GLFW_OPENGL_COMPAT_PROFILE  0x00032002
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
typedef void (* GLFWmonitorfun)(GLFWmonitor*,int);
typedef struct GLFWwindow GLFWwindow;
typedef void (* GLFWerrorfun)(int,const char*);
typedef void (* GLFWkeyfun)(GLFWwindow*,int,int,int,int);
typedef void (* GLFWcharfun)(GLFWwindow*,unsigned int);
typedef void (* GLFWframebuffersizefun)(GLFWwindow*,int,int);
typedef void (* GLFWwindowposfun)(GLFWwindow* window, int xpos, int ypos);
typedef void (* GLFWwindowsizefun)(GLFWwindow* window, int width, int height);
typedef void (* GLFWwindowcontentscalefun)(GLFWwindow* window, float xscale, float yscale);
typedef void (* GLFWwindowclosefun)(GLFWwindow* window);
typedef void (* GLFWjoystickfun)(int,int);
typedef void (* GLFWwindowrefreshfun)(GLFWwindow* window);
typedef void (* GLFWwindowfocusfun)(GLFWwindow* window, int focused);
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
	void (*GetSize)(int* width, int* height);
	void (*SetPos)(int x, int y);
	void (*GetPos)(int* x, int* y);
	int (*GetInputMode)(int);
	void (*SetInputMode)(int, int);
	void* (*GetUserPointer)(void);
	void (*SetUserPointer)(void*);
	void (*SetAttrib)(int, int);
	int (*GetAttrib)(int);
	GLFWmonitor* (*GetMonitor)(void);
	void (*SetMonitor)(GLFWmonitor*, int, int, int, int, int);
	// callbacks
	GLFWkeyfun (*SetKeyCallback)(GLFWkeyfun);
	GLFWcharfun (*SetCharCallback)(GLFWcharfun);
	GLFWframebuffersizefun (*SetFramebufferSizeCallback)(GLFWframebuffersizefun);
	GLFWwindowposfun (*SetPosCallback)(GLFWwindowposfun);
	GLFWwindowsizefun (*SetSizeCallback)(GLFWwindowsizefun);
	GLFWwindowcontentscalefun (*SetContentScaleCallback)(GLFWwindowcontentscalefun);
	GLFWwindowclosefun (*SetCloseCallback)(GLFWwindowclosefun);
	GLFWwindowrefreshfun (*SetRefreshCallback)(GLFWwindowrefreshfun);
	GLFWwindowfocusfun (*SetFocusCallback)(GLFWwindowfocusfun);
};
#define glfwGetWindowMonitor(win) ((GLFWwindow*)win)->GetMonitor()
#define glfwSetWindowMonitor(win, mon, x, y, w, h, r) ((GLFWwindow*)win)->SetMonitor(mon, x, y, w, h, r)
#define glfwGetWindowAttrib(win, x, y) ((GLFWwindow*)win)->GetAttrib(x)
#define glfwSetWindowAttrib(win, x, y) ((GLFWwindow*)win)->SetAttrib(x, y)
#define glfwMakeContextCurrent(win) ((GLFWwindow*)win)->MakeContextCurrent()
#define glfwWindowShouldClose(win) ((GLFWwindow*)win)->ShouldClose()
#define glfwSwapBuffers(win) ((GLFWwindow*)win)->SwapBuffers()
#define glfwGetKey(win, k) ((GLFWwindow*)win)->GetKey(k)
#define glfwSetWindowShouldClose(win, b) ((GLFWwindow*)win)->SetShouldClose(b)
#define glfwSetKeyCallback(win, cb) ((GLFWwindow*)win)->SetKeyCallback(cb)
#define glfwSetCharCallback(win, cb) ((GLFWwindow*)win)->SetCharCallback(cb)
#define glfwSetFramebufferSizeCallback(win, cb) ((GLFWwindow*)win)->SetFramebufferSizeCallback(cb)
#define glfwSetWindowPosCallback(win, cb) ((GLFWwindow*)win)->SetPosCallback(cb)
#define glfwSetWindowSizeCallback(win, cb) ((GLFWwindow*)win)->SetSizeCallback(cb)
#define glfwSetWindowContentScaleCallback(win, cb) ((GLFWwindow*)win)->SetContentScaleCallback(cb)
#define glfwSetWindowCloseCallback(win, cb) ((GLFWwindow*)win)->SetCloseCallback(cb)
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
#define glfwGetWindowSize(win, w, h) ((GLFWwindow*)win)->GetSize(w, h)
#define glfwSetWindowPos(win, x, y) ((GLFWwindow*)win)->SetPos(x, y)
#define glfwGetWindowPos(win, x, y) ((GLFWwindow*)win)->GetPos(x, y)
#define glfwGetWindowUserPointer(win) ((GLFWwindow*)win)->GetUserPointer()
#define glfwSetWindowUserPointer(win, ptr) ((GLFWwindow*)win)->SetUserPointer(ptr)
#define glfwGetInputMode(win, mode) ((GLFWwindow*)win)->GetInputMode(mode)
#define glfwSetInputMode(win, mode, v) ((GLFWwindow*)win)->SetInputMode(mode, v)
#define glfwSetWindowRefreshCallback(win, cb) ((GLFWwindow*)win)->SetRefreshCallback(cb)
#define glfwSetWindowFocusCallback(win, cb) ((GLFWwindow*)win)->SetRefreshCallback(cb)
typedef struct _GLFWjoystick {
	const float* (*GetAxes)(int* count);
	const unsigned char* (*GetButtons)(int* count);
	const unsigned char* (*GetHats)(int* count);
	const char* (*GetName)(void);
	_Bool (*IsGamepad)(void);
	const char* (*GetGUID)(void);
	const char* GetGamepadName(void);
} _GLFWjoystick;
#define glfwGetJoystickAxes(j, a) ((_GLFWjoystick)j).GetAxes(a)
#define glfwGetJoystickButtons(j, a) ((_GLFWjoystick)j).GetButtons(a)
#define glfwGetJoystickHats(j, a) ((_GLFWjoystick)j).GetHats(a)
#define glfwGetJoystickName(j) ((_GLFWjoystick)j).GetName()
#define glfwJoystickIsGamepad(j) ((_GLFWjoystick)j).IsGamepad()
#define glfwGetJoystickGUID(j) ((_GLFWjoystick)j).GetGUID()
#define glfwGetGamepadName(j) ((_GLFWjoystick)j).GetGamepadName()
void glfwWindowHint(int, int);
const char* glfwGetKeyName(int key, int scancode);