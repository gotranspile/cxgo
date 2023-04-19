package libs

import (
	"github.com/gotranspile/cxgo/types"
	"github.com/stretchr/testify/require"
	"testing"
)

autoRegLibs := []string {

	"AL/alc.h",
	"AL/al.h",
	
	"dirent.h",
	"fcntl.h",
	"float.h",
	"getopt.h",
	 
	"GL/gl.h",
	
	"inttypes.h",
	"libgen.h",
	"sched.h",
	"semaphore.h",
	"signal.h",
	"strings.h",
	 
	"sys/mkdev.h",
	"windows.h",
}

func TestReglibs(t *testing.T) {
	c := NewEnv(types.Config32())

	for _, lib := range autoRegLibs {
		_, ok := c.GetLibrary(lib)
		require.True(t, ok)
	}
}
