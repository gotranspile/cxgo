
#include <sys/types.h>

typedef struct {
	int pad;
} DIR;

typedef struct dirent {
	ino_t  d_ino;
	char   d_name[];
} dirent;

int            closedir(DIR *);
DIR           *opendir(const char *);
struct dirent *readdir(DIR *);
int            readdir_r(DIR *, struct dirent *, struct dirent **);
void           rewinddir(DIR *);
void           seekdir(DIR *, long int);
long int       telldir(DIR *);
