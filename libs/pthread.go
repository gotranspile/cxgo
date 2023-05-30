package libs

import (
	"github.com/gotranspile/cxgo/runtime/pthread"
	"github.com/gotranspile/cxgo/types"
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
		timespecT := c.GetLibraryType(timeH, "timespec")
		onceT := types.NamedTGo("pthread_once_t", "sync.Once", c.MethStructT(map[string]*types.FuncType{
			"Do": c.FuncTT(nil, c.FuncTT(nil)),
		}))
		mutexAttrT := types.NamedTGo("pthread_mutexattr_t", "pthread.MutexAttr", c.MethStructT(map[string]*types.FuncType{
			"Init":    c.FuncTT(intT),
			"SetType": c.FuncTT(intT, intT),
			"Destroy": c.FuncTT(intT),
		}))
		mutexT := types.NamedTGo("pthread_mutex_t", "pthread.Mutex", c.MethStructT(map[string]*types.FuncType{
			"Init":      c.FuncTT(intT, c.PtrT(mutexAttrT)),
			"Destroy":   c.FuncTT(intT),
			"CLock":     c.FuncTT(intT),
			"TryLock":   c.FuncTT(intT),
			"TimedLock": c.FuncTT(intT, c.PtrT(timespecT)),
			"CUnlock":   c.FuncTT(intT),
		}))
		condAttrT := types.NamedTGo("pthread_condattr_t", "pthread.CondAttr", types.StructT(nil))
		condT := types.NamedTGo("pthread_cond_t", "sync.Cond", types.StructT([]*types.Field{
			{Name: types.NewIdent("L", c.PtrT(mutexT))},
			{Name: types.NewIdent("Wait", c.FuncTT(nil))},
			{Name: types.NewIdent("Signal", c.FuncTT(nil))},
			{Name: types.NewIdent("Broadcast", c.FuncTT(nil))},
		}))
		threadT := types.NamedTGo("pthread_t", "pthread.Thread", c.MethStructT(map[string]*types.FuncType{
			"Join":        c.FuncTT(intT, c.PtrT(retT)),
			"TimedJoinNP": c.FuncTT(intT, c.PtrT(retT), c.PtrT(timespecT)),
		}))
		threadAttrT := types.NamedTGo("pthread_attr_t", "pthread.Attr", c.MethStructT(map[string]*types.FuncType{}))
		return &Library{
			Imports: map[string]string{
				"sync":    "sync",
				"pthread": RuntimePrefix + "pthread",
			},
			Types: map[string]types.Type{
				"pthread_t_":          threadT,
				"pthread_t":           c.PtrT(threadT),
				"pthread_once_t":      onceT,
				"pthread_cond_t":      condT,
				"pthread_condattr_t":  condAttrT,
				"pthread_attr_t":      threadAttrT,
				"pthread_mutex_t":     mutexT,
				"pthread_mutexattr_t": mutexAttrT,
			},
			Idents: map[string]*types.Ident{
				"PTHREAD_MUTEX_RECURSIVE": c.NewIdent("PTHREAD_MUTEX_RECURSIVE", "pthread.MUTEX_RECURSIVE", pthread.MUTEX_RECURSIVE, gintT),
				"pthread_create":          c.NewIdent("pthread_create", "pthread.Create", pthread.Create, c.FuncTT(intT, c.PtrT(c.PtrT(threadT)), c.PtrT(threadAttrT), c.FuncTT(retT, argT), argT)),
				"pthread_cond_init":       c.NewIdent("pthread_cond_init", "pthread.CondInit", pthread.CondInit, c.FuncTT(intT, c.PtrT(condT), c.PtrT(condAttrT))),
				"pthread_cond_destroy":    c.NewIdent("pthread_cond_destroy", "pthread.CondFree", pthread.CondFree, c.FuncTT(intT, c.PtrT(condT))),
			},
		}
	})
}
