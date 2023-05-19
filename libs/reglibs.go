package libs

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed includes/embed/*
var efs embed.FS

func init() {

	onFile := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		return registerInclude(path)
	}

	err := fs.WalkDir(efs, ".", onFile)

	if err != nil {
		panic(err)
	}
}

func registerInclude(fName string) error {

	fBuff, fErr := efs.ReadFile(fName)

	if fErr != nil {
		return fErr
	}

	RegisterLibrary(strings.TrimPrefix(fName, "includes/"), func(env *Env) *Library {
		return &Library{
			Header: fmt.Sprintf("#include <%s>\n%s", BuiltinH, string(fBuff)),
		}
	})

	return nil
}
