package libs

const (
	getoptH = "getopt.h"
)

func init() {
	RegisterLibrary(getoptH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
struct option {
   const char *name;
   int         has_arg;
   int        *flag;
   int         val;
};

int getopt(int argc, char * const argv[], const char *optstring);
int getopt_long(int argc, char * const argv[], const char *optstring,
                  const struct option *longopts, int *longindex);
int getopt_long_only(int argc, char * const argv[], const char *optstring,
                  const struct option *longopts, int *longindex);
char *optarg;
int optind, opterr, optopt, required_argument, no_argument;
`,
		}
	})
}
