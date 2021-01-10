package cxgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/dennwc/cxgo/internal/git"
	"github.com/dennwc/cxgo/libs"
	"github.com/dennwc/cxgo/types"
	"github.com/stretchr/testify/require"
)

func downloadGCC(t testing.TB, dst string) {
	err := os.MkdirAll(dst, 0755)
	require.NoError(t, err)

	const (
		repo   = "https://github.com/gcc-mirror/gcc.git"
		sub    = "gcc/testsuite/gcc.c-torture/compile"
		branch = "releases/gcc-10.2.0"
	)

	dir := filepath.Join(os.TempDir(), "cxgo_gcc_git")
	_ = os.RemoveAll(dir)

	t.Log("cloning", repo, "to", dir)
	err = git.Clone(repo, branch, dir)
	if err != nil {
		os.RemoveAll(dst)
	}
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	files, err := filepath.Glob(filepath.Join(dir, sub, "*.c"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, path := range files {
		base := filepath.Base(path)
		err = copyFile(path, filepath.Join(dst, base))
		require.NoError(t, err)
	}
}

func TestGCCExecute(t *testing.T) {
	if testing.Short() || os.Getenv("CXGO_RUN_TESTS_GCC") != "true" {
		t.SkipNow()
	}
	dir := filepath.Join(testDataDir, "gcc")

	ignoreTests := map[string]string{
		"limits-caselabels": "OOM",
	}

	blacklist := map[string]struct{}{}

	isLib := map[string]struct{}{}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		downloadGCC(t, dir)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.c"))
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	out := filepath.Join(os.TempDir(), "cxgo_test_gcc")
	err = os.MkdirAll(out, 0755)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(out, "go.mod"), []byte(fmt.Sprintf(`module main
go 1.13
require (
	github.com/dennwc/cxgo v0.0.0
)
replace github.com/dennwc/cxgo => %s`, wd)), 0644)
	require.NoError(t, err)

	for _, path := range files {
		path := path
		tname := strings.TrimSuffix(filepath.Base(path), ".c")
		_, skip := blacklist[tname]
		t.Run(tname, func(t *testing.T) {
			if reason, ignore := ignoreTests[tname]; ignore {
				t.Skip(reason)
			}
			//t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					if skip {
						defer debug.PrintStack()
						t.Skipf("panic: %v", r)
					} else {
						require.Nil(t, r)
					}
				}
				if !t.Failed() && !t.Skipped() && skip {
					t.Error("blacklisted test pass")
				}
			}()
			oname := filepath.Base(path) + ".go"
			env := libs.NewEnv(types.Config32())
			_, lib := isLib[tname]
			if data, err := ioutil.ReadFile(path); err == nil && !bytes.Contains(data, []byte("main")) {
				t.Log("testing as a library (no main found)")
				lib = true
			}
			pkg := "main"
			if lib {
				pkg = "lib"
			}
			err = Translate(filepath.Dir(path), path, out, env, Config{
				Predef: `
#include <stdlib.h>
#include <string.h>
`,
				Package:            pkg,
				GoFile:             oname,
				MaxDecls:           -1,
				FixImplicitReturns: true,
			})
			failOrSkip(t, err, skip)

			t.Log(path)
			t.Log(filepath.Join(out, oname))
			goRun(t, out, []string{oname}, runConfig{Arch32: false, Skip: skip, BuildOnly: lib})
		})
	}
}
