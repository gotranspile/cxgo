#include <stddef.h>

const _cxgo_go_int GLOB_NOESCAPE = 1;

typedef struct {
    size_t   gl_pathc;    /* Count of paths matched so far  */
    char   **gl_pathv;    /* List of matched pathnames.  */
    size_t   gl_offs;     /* Slots to reserve in gl_pathv.  */
	_cxgo_sint32 (*Glob)(const char *pattern, _cxgo_sint32 flags,
                _cxgo_sint32 (*errfunc) (const char *epath, _cxgo_sint32 eerrno));
	void (*Free)(void);
} glob_t;
#define glob(pattern, flags, errfunc, g) ((glob_t*)g)->Glob(pattern, flags, errfunc)
#define globfree(g) ((glob_t*)g)->Free()
