package libs

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed includes/embed
var efs embed.FS

func init() {
	const root = "includes/embed"
	err := fs.WalkDir(efs, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if d.IsDir() {
			return nil
		}
		data, err := efs.ReadFile(path)
		if err != nil {
			return err
		}
		fname := strings.TrimPrefix(path, root+"/")
		RegisterLibrary(fname, func(env *Env) *Library {
			return &Library{
				Header: fmt.Sprintf("#include <%s>\n%s", BuiltinH, string(data)),
			}
		})
		return nil
	})
	if err != nil {
		panic(err)
	}
}
