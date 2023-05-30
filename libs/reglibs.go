package libs

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed includes
var efs embed.FS

func init() {
	const root = "includes"
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
		RegisterLibrarySrc(fname, string(data))
		return nil
	})
	if err != nil {
		panic(err)
	}
}
