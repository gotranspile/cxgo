package libs

import (
	"github.com/dennwc/cxgo/runtime/pthread"
	"github.com/dennwc/cxgo/types"
)

const (
	pthreadH = "pthread.h"
)

func init() {
	RegisterLibrary(pthreadH, func(c *Env) *Library {
		gintT := c.Go().Int()
		intT := types.IntT(4)
		argT := c.PtrT(nil)
		retT := c.PtrT(nil)
		timeLib := c.GetLib(timeH)
		timespecT := timeLib.GetType("timespec")
		threadT := types.NamedTGo("pthread_t", "pthread.Thread", c.MethStructT(map[string]*types.FuncType{
			"Join":        c.FuncTT(intT, c.PtrT(retT)),
			"TimedJoinNP": c.FuncTT(intT, c.PtrT(retT), c.PtrT(timespecT)),
		}))
		threadAttrT := types.NamedTGo("pthread_attr_t", "pthread.Attr", c.MethStructT(map[string]*types.FuncType{}))
		mutexAttrT := types.NamedTGo("pthread_mutexattr_t", "pthread.MutexAttr", c.MethStructT(map[string]*types.FuncType{
			"Init":    c.FuncTT(intT),
			"SetType": c.FuncTT(intT, intT),
			"Destroy": c.FuncTT(intT),
		}))
		mutexT := types.NamedTGo("pthread_mutex_t", "pthread.Mutex", c.MethStructT(map[string]*types.FuncType{
			"Init":      c.FuncTT(intT, c.PtrT(mutexAttrT)),
			"Destroy":   c.FuncTT(intT),
			"Lock":      c.FuncTT(intT),
			"TryLock":   c.FuncTT(intT),
			"TimedLock": c.FuncTT(intT, c.PtrT(timespecT)),
			"Unlock":    c.FuncTT(intT),
		}))
		return &Library{
			Imports: map[string]string{
				"pthread": RuntimePrefix + "pthread",
			},
			Types: map[string]types.Type{
				"pthread_t_":          threadT,
				"pthread_t":           c.PtrT(threadT),
				"pthread_attr_t":      threadAttrT,
				"pthread_mutex_t":     mutexT,
				"pthread_mutexattr_t": mutexAttrT,
			},
			Idents: map[string]*types.Ident{
				"PTHREAD_MUTEX_RECURSIVE": c.NewIdent("PTHREAD_MUTEX_RECURSIVE", "pthread.MUTEX_RECURSIVE", pthread.MUTEX_RECURSIVE, gintT),
				"pthread_create":          c.NewIdent("pthread_create", "pthread.Create", pthread.Create, c.FuncTT(intT, c.PtrT(c.PtrT(threadT)), c.PtrT(threadAttrT), c.FuncTT(retT, argT), argT)),
			},
			// TODO
			Header: `
#include <` + BuiltinH + `>
#include <` + timeH + `>

const _cxgo_go_int PTHREAD_MUTEX_RECURSIVE = 1;

typedef struct pthread_attr_t {} pthread_attr_t;

typedef struct{
	_cxgo_sint32 (*Join)(void **retval);
	_cxgo_sint32 (*TimedJoinNP)(void **retval, const struct timespec *abstime);
} pthread_t_;
#define pthread_t pthread_t_*

_cxgo_sint32 pthread_create(pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg);

typedef struct pthread_mutexattr_t {
	_cxgo_sint32 (*Init)(void);
	_cxgo_sint32 (*SetType)(_cxgo_sint32 type);
	_cxgo_sint32 (*Destroy)(void);
} pthread_mutexattr_t;
#define pthread_mutexattr_init(attr) ((pthread_mutexattr_t*)attr)->Init()
#define pthread_mutexattr_settype(attr, type) ((pthread_mutexattr_t*)attr)->SetType(type)
#define pthread_mutexattr_destroy(attr) ((pthread_mutexattr_t*)attr)->Destroy()

typedef struct pthread_mutex_t {
	_cxgo_sint32 (*Destroy)(void);
	_cxgo_sint32 (*Init)(const pthread_mutexattr_t *restrict attr);
	_cxgo_sint32 (*Lock)(void);
	_cxgo_sint32 (*TryLock)(void);
	_cxgo_sint32 (*Unlock)(void);
	_cxgo_sint32 (*TimedLock)(const struct timespec *restrict abstime);
} pthread_mutex_t;
#define pthread_mutex_destroy(mutex) ((pthread_mutex_t*)mutex)->Destroy()
#define pthread_mutex_init(mutex, attr) ((pthread_mutex_t*)mutex)->Init(attr)
#define pthread_mutex_lock(mutex) ((pthread_mutex_t*)mutex)->Lock()
#define pthread_mutex_trylock(mutex) ((pthread_mutex_t*)mutex)->TryLock()
#define pthread_mutex_unlock(mutex) ((pthread_mutex_t*)mutex)->Unlock()
#define pthread_mutex_timedlock(mutex, abstime) ((pthread_mutex_t*)mutex)->TimedLock(abstime)

#define pthread_join(thread, retval) ((pthread_t_*)thread)->Join(retval)
#define pthread_timedjoin_np(thread, retval, abstime) ((pthread_t_*)thread)->TimedJoinNP(retval, abstime)
`,
		}
	})
}
