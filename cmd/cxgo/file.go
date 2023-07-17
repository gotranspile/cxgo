package main

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gotranspile/cxgo"
	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/types"
)

func init() {
	cmdFile := &cobra.Command{
		Use:   "file",
		Short: "transpile a single C file to Go",
	}
	Root.AddCommand(cmdFile)

	fOut := cmdFile.Flags().StringP("out", "o", "", "output file to write to")
	fPkg := cmdFile.Flags().StringP("pkg", "p", "main", "package name for a Go file")
	cmdFile.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly one file must be specified")
		}
		in := args[0]
		out := *fOut
		if out == "" {
			out = strings.TrimSuffix(in, filepath.Ext(in)) + ".go"
		}
		env := libs.NewEnv(types.Config{
			UseGoInt: true,
		})
		fc := cxgo.Config{
			Package:  *fPkg,
			GoFile:   filepath.Base(out),
			MaxDecls: -1,
		}
		return cxgo.Translate("", in, filepath.Dir(out), env, fc)
	}
}
