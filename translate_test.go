package cxgo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testDataDir = "./.testdata"

func failOrSkip(t testing.TB, err error, skip bool) {
	t.Helper()
	if err == nil {
		return
	}
	if skip {
		t.Skip(err)
		return
	}
	require.NoError(t, err)
}

func goBuild(t testing.TB, wd, bin string, files []string, c runConfig) {
	buf := bytes.NewBuffer(nil)
	args := []string{"build", "-o", bin}
	args = append(args, files...)
	t.Logf("go %q", args)
	cmd := exec.Command("go", args...)
	cmd.Dir = wd
	cmd.Env = os.Environ()
	if c.Arch32 {
		cmd.Env = append(cmd.Env, "GOARCH=386")
	}
	cmd.Stderr = buf
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("compilation failed: %w; output:\n%s\n", err, buf)
	}
	failOrSkip(t, err, c.Skip)
}

func cmdRun(t testing.TB, wd, bin string, skip bool, args ...string) []byte {
	buf := bytes.NewBuffer(nil)
	ebuf := bytes.NewBuffer(nil)
	cmd := exec.Command(bin, args...)
	cmd.Dir = wd
	cmd.Stderr = io.MultiWriter(ebuf, buf)
	cmd.Stdout = buf
	timeout := time.NewTimer(time.Minute / 2)
	defer timeout.Stop()
	err := cmd.Start()
	require.NoError(t, err)
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()
	select {
	case <-timeout.C:
		_ = cmd.Process.Kill()
		require.Fail(t, "timeout")
	case err = <-errc:
		if err != nil && skip {
			t.Skipf("program failed; output:\n%s\n", buf)
		}
		require.NoError(t, err, "program failed; output:\n%s\n", buf)
	}
	return buf.Bytes()
}

type runConfig struct {
	Arch32    bool
	Skip      bool
	BuildOnly bool
}

func goRun(t testing.TB, wd string, files []string, c runConfig, args ...string) []byte {
	f, err := ioutil.TempFile("", "cxgo_bin_")
	require.NoError(t, err)
	_ = f.Close()
	bin := f.Name()
	defer os.Remove(bin)
	goBuild(t, wd, bin, files, c)
	if c.BuildOnly {
		return nil
	}
	return cmdRun(t, wd, bin, c.Skip, args...)
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return d.Close()
}

func TestTranslateSnippets(t *testing.T) {
	runTestTranslate(t, casesTranslateSnippets)
}

var casesTranslateSnippets = []parseCase{
	{
		name: "files",
		src: `
#include <stdio.h>

void foo() {
	char* mode = "r";
	FILE* f = fopen("file.dat", mode);
	if (f == 0) {
		return;
	}
	char b[10];
	fread(b, 10, 1, f);
	fclose(f);
}
`,
		exp: `
func foo() {
	var (
		mode *byte       = libc.CString("r")
		f    *stdio.File = stdio.FOpen("file.dat", libc.GoString(mode))
	)
	if f == nil {
		return
	}
	var b [10]byte
	f.ReadN(&b[0], 10, 1)
	f.Close()
}
`,
	},
}

const (
	kindUint = 1
	kindInt  = 2
)

type intVal struct {
	Size int
	Kind int
	ValI int64
	ValU uint64
}

func (v intVal) CType() string {
	u := ""
	if v.Kind == kindUint {
		u = "u"
	}
	return fmt.Sprintf("%sint%d_t", u, v.Size*8)
}

func (v intVal) CValue() string {
	switch v.Kind {
	case kindUint:
		return strconv.FormatUint(v.ValU, 10) + "ul"
	case kindInt:
		return strconv.FormatInt(v.ValI, 10) + "l"
	default:
		panic("should not happen")
	}
}

type binopTest struct {
	To intVal
	X  intVal
	Op BinaryOp
	Y  intVal
}

func (b binopTest) CSrc() string {
	verb := "d"
	if b.To.Kind == kindUint {
		verb = "u"
	}
	if b.To.Size == 8 {
		verb = "ll" + verb
	}
	return fmt.Sprintf(`
#include <stdint.h>
#include <stdio.h>

int main() {
	%s a = %s;
	%s b = %s;
	%s c = 0;
	c = a %s b;
	printf("%%%s\n", c);
	return 0;
}
`,
		b.X.CType(), b.X.CValue(),
		b.Y.CType(), b.Y.CValue(),
		b.To.CType(),
		string(b.Op),
		verb,
	)
}

func randIntVal() intVal {
	v := intVal{
		Kind: kindUint,
		Size: 1 << (rand.Int() % 4),
	}
	switch v.Size {
	case 1, 2, 4, 8:
	default:
		panic(v.Size)
	}
	if rand.Int()%2 == 0 {
		v.Kind = kindInt
	}
	switch v.Kind {
	case kindUint:
		v.ValU = rand.Uint64()
		switch v.Size {
		case 1:
			v.ValU %= math.MaxUint8 + 1
		case 2:
			v.ValU %= math.MaxUint16 + 1
		case 4:
			v.ValU %= math.MaxUint32 + 1
		}
	case kindInt:
		v.ValI = rand.Int63()
		switch v.Size {
		case 1:
			v.ValI %= math.MaxInt8 + 1
		case 2:
			v.ValI %= math.MaxInt16 + 1
		case 4:
			v.ValI %= math.MaxInt32 + 1
		}
		if rand.Int()%2 == 0 {
			v.ValI = -v.ValI
		}
		switch v.Size {
		case 1:
			for v.ValI > math.MaxInt8 {
				v.ValI -= math.MaxInt8
			}
			for v.ValI < math.MinInt8 {
				v.ValI += math.MinInt8
			}
		case 2:
			for v.ValI > math.MaxInt16 {
				v.ValI -= math.MaxInt16
			}
			for v.ValI < math.MinInt16 {
				v.ValI += math.MinInt16
			}
		case 4:
			for v.ValI > math.MaxInt32 {
				v.ValI -= math.MaxInt32
			}
			for v.ValI < math.MinInt32 {
				v.ValI += math.MinInt32
			}
		}
	}
	return v
}

func randBinop(op BinaryOp) binopTest {
	return binopTest{
		To: randIntVal(),
		X:  randIntVal(),
		Op: op,
		Y:  randIntVal(),
	}
}

func TestImplicitCompat(t *testing.T) {
	for _, op := range []BinaryOp{
		BinOpAdd,
		BinOpSub,
		BinOpMult,
		BinOpDiv,
		BinOpMod,
	} {
		op := op
		t.Run(string(op), func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 25; i++ {
				func() {
					b := randBinop(op)
					for (b.Op == BinOpMod || b.Op == BinOpDiv) && b.Y.ValI == 0 && b.Y.ValU == 0 {
						b = randBinop(op)
					}
					testTranspileOut(t, b.CSrc())
				}()
			}
		})
	}
}

func testTranspileOut(t testing.TB, csrc string) {
	dir, err := ioutil.TempDir("", "cxgo_cout")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cdir := filepath.Join(dir, "c")
	err = os.MkdirAll(cdir, 0755)
	require.NoError(t, err)

	cfile := filepath.Join(cdir, "main.c")
	err = ioutil.WriteFile(cfile, []byte(csrc), 0644)
	require.NoError(t, err)

	// this is required for test to wait other goroutines in case it fails earlier on the main one
	var wg sync.WaitGroup
	wg.Add(1)
	cch := make(chan progOut, 1)
	go func() {
		defer wg.Done()
		cch <- gccCompileAndExec(t, cdir, cfile)
	}()
	defer wg.Wait()

	godir := filepath.Join(dir, "golang")
	err = os.MkdirAll(godir, 0755)
	require.NoError(t, err)

	goout := goTranspileAndExec(t, ".", godir, cfile)
	cout := <-cch
	require.Equal(t, cout.Code, goout.Code,
		"\nC code: %d (%v)\nGo code: %s (%v)",
		cout.Code, cout.Err,
		goout.Code, goout.Err,
	)
	require.Equal(t, cout.Err, goout.Err)
	t.Logf("// === Output ===\n%s", cout.Out)
	require.Equal(t, cout.Out, goout.Out, "\n// === C source ===\n%s", csrc)
}
