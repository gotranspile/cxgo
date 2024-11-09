package cxgo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/types"
)

func gccCompile(t testing.TB, out, cfile string) {
	var buf bytes.Buffer
	cmd := exec.Command("gcc", "-O0", "-o", out, cfile)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("gcc compilation failed: %v\n%s", err, buf.String())
	}
}

func gccCompileAndExec(t testing.TB, dir, cfile string) progOut {
	cbin := filepath.Join(dir, "out")
	gccCompile(t, cbin, cfile)
	return progExecute(dir, cbin)
}

func goCompile(t testing.TB, out, dir string) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = dir
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("go compilation failed: %v\n%s", err, buf.String())
	}
}

func goCompileAndExec(t testing.TB, dir string) progOut {
	gobin := filepath.Join(dir, "out")
	goCompile(t, gobin, dir)
	return progExecute(dir, gobin)
}

func goTranspileAndExec(t testing.TB, cxgo, dir string, cfile string) progOut {
	goProject(t, dir, cxgo)
	gofile := filepath.Join(dir, "main.go")
	env := libs.NewEnv(types.Config32())
	err := Translate(filepath.Dir(cfile), cfile, dir, env, Config{
		Package:     "main",
		GoFile:      gofile,
		MaxDecls:    -1,
		ForwardDecl: false,
	})
	require.NoError(t, err)
	goProjectMod(t, dir)
	gosrc, err := os.ReadFile(gofile)
	require.NoError(t, err)
	t.Logf("// === Go source ===\n%s", string(gosrc))
	return goCompileAndExec(t, dir)
}

type progOut struct {
	Code int
	Err  error
	Out  string
}

func execInDir(wd string, bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = wd
	return cmd.Run()
}

func progExecute(wd, bin string) progOut {
	var buf bytes.Buffer
	cmd := exec.Command(bin)
	cmd.Dir = wd
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := progOut{
		Err: err,
		Out: buf.String(),
	}
	if e, ok := err.(*exec.ExitError); ok && e.Exited() {
		out.Code = e.ExitCode()
		out.Err = nil
	}
	return out
}

func goProject(t testing.TB, out, cxgo string) {
	cxgo, err := filepath.Abs(cxgo)
	require.NoError(t, err)

	err = os.MkdirAll(out, 0755)
	require.NoError(t, err)

	gomod := fmt.Sprintf(`module main
go 1.19
require (
	github.com/gotranspile/cxgo v0.0.0-local
)
replace github.com/gotranspile/cxgo v0.0.0-local => %s`, cxgo)

	err = os.WriteFile(filepath.Join(out, "go.mod"), []byte(gomod), 0644)
	require.NoError(t, err)

	// allows running go mod tidy without having other source files and still keep require above
	err = os.WriteFile(filepath.Join(out, "dummy.go"), []byte(`
package main

import _ "github.com/gotranspile/cxgo/runtime/libc"
`), 0644)
	require.NoError(t, err)
}

func goProjectMod(t testing.TB, out string) {
	err := execInDir(out, "go", "mod", "tidy")
	require.NoError(t, err)
}
