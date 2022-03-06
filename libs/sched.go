package libs

const (
	schedH = "sched.h"
)

func init() {
	RegisterLibrary(schedH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + BuiltinH + `>
#include <` + sysTypesH + `>

typedef struct sched_param {} sched_param;

int    sched_get_priority_max(int);
int    sched_get_priority_min(int);
int    sched_getparam(pid_t, struct sched_param *);
int    sched_getscheduler(pid_t);
int    sched_rr_get_interval(pid_t, struct timespec *);
int    sched_setparam(pid_t, const struct sched_param *);
int    sched_setscheduler(pid_t, int, const struct sched_param *);
int    sched_yield(void);
`,
		}
	})
}
