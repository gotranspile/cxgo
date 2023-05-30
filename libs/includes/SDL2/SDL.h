#define Uint32 _cxgo_uint32
#define Uint16 _cxgo_uint16
#define Uint8 _cxgo_uint8
#define Sint32 _cxgo_sint32
#define Sint16 _cxgo_sint16
#define Sint8 _cxgo_sint8
#define SDL_bool _cxgo_go_bool
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