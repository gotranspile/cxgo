package libs

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

func init() {
	RegisterLibraries()
}

//go:embed includes/*
var efs embed.FS

func RegisterLibraries() {

	for {
		fNames, err := fileNames()

		if err != nil {
			panic(err)
		}

		for _, fName := range fNames {

			fBuff, fErr := efs.ReadFile(fName)
			if fErr != nil {
				break
			}

			include := strings.TrimPrefix(fName, "includes/")
			registerInclude(include, string(fBuff))
		}

		break
	}

	return
}

func fileNames() (files []string, err error) {

	wdFunc := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}

	if err := fs.WalkDir(efs, ".", wdFunc); err != nil {
		return nil, err
	}

	return files, nil
}

func registerInclude(include, header string) {

	header = fmt.Sprintf("#include <%s>\n%s", BuiltinH, header)

	RegisterLibrary(include, func(env *Env) *Library {
		return &Library{
			Header: header,
		}
	})
}
