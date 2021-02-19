package cxgo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/gotranspile/cxgo/internal/git"
	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/types"
	"github.com/stretchr/testify/require"
)

func downloadTCC(t testing.TB, dst string) {
	err := os.MkdirAll(dst, 0755)
	require.NoError(t, err)

	const (
		repo   = "https://repo.or.cz/tinycc.git"
		sub    = "tests/tests2"
		branch = "release_0_9_27"
	)

	dir := filepath.Join(os.TempDir(), "cxgo_tcc_git")
	_ = os.RemoveAll(dir)

	t.Log("cloning", repo, "to", dir)
	err = git.Clone(repo, branch, dir)
	if err != nil {
		os.RemoveAll(dst)
	}
	require.NoError(t, err)
	//defer os.RemoveAll(dir)

	files, err := filepath.Glob(filepath.Join(dir, sub, "*.c"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, path := range files {
		base := filepath.Base(path)
		err = copyFile(path, filepath.Join(dst, base))
		require.NoError(t, err)
		bexp := strings.TrimSuffix(base, ".c") + ".expect"
		bhdr := strings.TrimSuffix(base, ".c") + ".h"
		err = copyFile(filepath.Join(filepath.Dir(path), bexp), filepath.Join(dst, bexp))
		require.NoError(t, err)
		_ = copyFile(filepath.Join(filepath.Dir(path), bhdr), filepath.Join(dst, bhdr))
	}
}

func TestTCCExecute(t *testing.T) {
	if testing.Short() || os.Getenv("CXGO_RUN_TESTS_TCC") != "true" {
		t.SkipNow()
	}
	dir := filepath.Join(testDataDir, "tcc")

	ignoreTests := map[string]string{
		"81_types":                "incomplete types, invalid calls",
		"85_asm-outside-function": "uses assembly",
		"98_al_ax_extend":         "uses assembly",
		"99_fastcall":             "uses assembly",
	}
	blacklist := map[string]struct{}{
		"34_array_assignment":     {}, // CC parsing failure
		"46_grep":                 {},
		"51_static":               {},
		"54_goto":                 {},
		"55_lshift_type":          {},
		"60_errors_and_warnings":  {},
		"73_arm64":                {},
		"75_array_in_struct_init": {}, // CC type checker failure
		"77_push_pop_macro":       {}, // FIXME: cannot detect if a macro was redefined
		"78_vla_label":            {},
		"79_vla_continue":         {},
		"87_dead_code":            {},
		"88_codeopt":              {},
		"89_nocode_wanted":        {},
		"90_struct-init":          {},
		"92_enum_bitfield":        {},
		"93_integer_promotion":    {}, // TODO: some issues with char promotion
		"94_generic":              {},
		"95_bitfields":            {},
		"95_bitfields_ms":         {},
		"96_nodata_wanted":        {},
		"97_utf8_string_literal":  {},
	}
	overrideExpect := map[string]string{
		// TODO: printf("%f", 12.34 + 56.78); TCC: 69.120003, Go: 69.119995
		"22_floating_point": "69.120003\n69.120000\n-44.440000\n700.665200\n0.217330\n1 1 0 0 0 1\n0 1 1 1 0 0\n0 0 0 1 1 1\n69.119995\n-44.439999\n700.665222\n0.217330\n12.340000\n-12.340000\n2.000000\n0.909297\n",
		// TODO: we leave an additional whitespace after each line; probably tcc test trims them when checking output?
		"38_multiple_array_index": "x=0: 1 2 3 4 \nx=1: 5 6 7 8 \nx=2: 9 10 11 12 \nx=3: 13 14 15 16 \n",
		// TODO: we don't implement __TINYC__ and related extensions
		"70_floating_point_literals": "0.123000\n122999996416.000000\n0.000000\n122999996416.000000\n\n123.123001\n123122997002240.000000\n0.000000\n123122997002240.000000\n\n123.000000\n123000003231744.000000\n0.000000\n123000003231744.000000\n\n123000003231744.000000\n0.000000\n123000003231744.000000\n\n\n428.000000\n0.000026\n428.000000\n\n1756112.000000\n0.104672\n1756592.000000\n\n1753088.000000\n0.104492\n1753088.000000\n\n1753088.000000\n0.104492\n1753088.000000\n\n\n",
		// TODO: additional expected line break; probably trimmed by tcc test?
		"71_macro_empty_arg": "17",
		// TODO: additional expected line break; probably trimmed by tcc test?
		"76_dollars_in_identifiers": "fred=10\njoe=20\nhenry=30\nfred2=10\njoe2=20\nhenry2=30\nfred10=100\njoe_10=2\nlocal=10\na100$=100\na$$=1000\na$c$b=2121\n$100=10000\n$$$=money",
		// TODO: we don't support bitfields yet
		"93_integer_promotion": " unsigned : s.ub\n unsigned : s.u\n unsigned : s.ullb\n unsigned : s.ull\n   signed : s.c\n\n unsigned : (1 ? s.ub : 1)\n unsigned : (1 ? s.u : 1)\n unsigned : (1 ? s.ullb : 1)\n unsigned : (1 ? s.ull : 1)\n   signed : (1 ? s.c : 1)\n\n unsigned : s.ub << 1\n unsigned : s.u << 1\n unsigned : s.ullb << 1\n unsigned : s.ull << 1\n   signed : s.c << 1\n\n   signed : +s.ub\n unsigned : +s.u\n unsigned : +s.ullb\n unsigned : +s.ull\n   signed : +s.c\n\n   signed : -s.ub\n unsigned : -s.u\n unsigned : -s.ullb\n unsigned : -s.ull\n   signed : -s.c\n\n   signed : ~s.ub\n unsigned : ~s.u\n unsigned : ~s.ullb\n unsigned : ~s.ull\n   signed : ~s.c\n\n   signed : !s.ub\n   signed : !s.u\n   signed : !s.ullb\n   signed : !s.ull\n   signed : !s.c\n\n unsigned : +(unsigned)s.ub\n unsigned : -(unsigned)s.ub\n unsigned : ~(unsigned)s.ub\n   signed : !(unsigned)s.ub\n",
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		downloadTCC(t, dir)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.c"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	wd, err := os.Getwd()
	require.NoError(t, err)

	out := filepath.Join(os.TempDir(), "cxgo_test_tcc")
	err = os.MkdirAll(out, 0755)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(out, "go.mod"), []byte(fmt.Sprintf(`module main
go 1.13
require (
	github.com/gotranspile/cxgo v0.0.0
)
replace github.com/gotranspile/cxgo => %s`, wd)), 0644)
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
			err = Translate(filepath.Dir(path), path, out, env, Config{
				Package:     "main",
				GoFile:      oname,
				MaxDecls:    -1,
				ForwardDecl: true,
			})
			failOrSkip(t, err, skip)

			t.Log(path)
			t.Log(filepath.Join(out, oname))

			args := []string{"arg1", "arg2", "arg3", "arg4", "arg5"}
			got := goRun(t, out, []string{"./" + oname}, runConfig{Arch32: true, Skip: skip}, args...)
			epath := strings.TrimSuffix(path, ".c") + ".expect"
			var exp []byte
			if s, ok := overrideExpect[tname]; ok {
				exp = []byte(s)
			} else {
				exp, err = ioutil.ReadFile(epath)
				require.NoError(t, err)
			}
			if string(exp) != string(got) && skip {
				t.Skipf("unexpected output:\nexpected: %q\nactual  : %q", exp, got)
			}
			require.Equal(t, string(exp), string(got))
		})
	}
}
