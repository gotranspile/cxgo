package cxgo

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/laher/mergefs"
	"modernc.org/cc/v4"

	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/libs/libcc"
)

type SourceConfig struct {
	Predef           string
	Define           []Define
	Include          []string
	SysInclude       []string
	IgnoreIncludeDir bool
}

func Parse(c *libs.Env, root, fname string, sconf SourceConfig) (*cc.AST, error) {
	path := filepath.Dir(fname)
	if root == "" {
		root = path
	}
	srcs := []cc.Source{{Name: fname}}
	if sconf.Predef != "" {
		srcs = []cc.Source{
			{Name: "predef.h", Value: sconf.Predef}, // FIXME: this should preappend to the file content instead
			{Name: fname},
		}
	}
	var (
		inc []string
		sys []string
	)
	inc = append(inc, sconf.Include...)
	sys = append(sys, sconf.SysInclude...)
	if !sconf.IgnoreIncludeDir {
		inc = append(inc,
			filepath.Join(root, "includes"),
			filepath.Join(root, "include"),
		)
		sys = append(sys, []string{
			filepath.Join(root, "includes"),
			filepath.Join(root, "include"),
		}...)
	}
	inc = append(inc,
		path,
		"@",
	)
	return ParseSource(c, ParseConfig{
		Sources:     srcs,
		WorkDir:     path,
		Includes:    inc,
		SysIncludes: sys,
		Predefines:  true,
		Define:      sconf.Define,
	})
}

func addIncludeOverridePath(inc []string) []string {
	return append(inc[:len(inc):len(inc)], libs.IncludePath)
}

const (
	tokLBrace = "("
	tokRBrace = ")"
)

type Define struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type ParseConfig struct {
	WorkDir     string
	Includes    []string
	SysIncludes []string
	Predefines  bool
	Define      []Define
	Sources     []cc.Source
}

func ParseSource(env *libs.Env, c ParseConfig) (*cc.AST, error) {
	var srcs []cc.Source
	if len(c.Define) != 0 {
		var buf bytes.Buffer
		for _, d := range c.Define {
			buf.WriteString("#define ")
			buf.WriteString(strings.TrimSpace(d.Name))
			if d.Value != "" {
				buf.WriteByte(' ')
				buf.WriteString(strings.TrimSpace(d.Value))
			}
			buf.WriteByte('\n')
		}
		srcs = append(srcs, cc.Source{Name: "<config-defines>", Value: buf.String()})
	}
	srcs = append(srcs, cc.Source{
		Name: "<cc-builtin>", Value: `
#define __UINT16_TYPE__ unsigned short
#define __UINT32_TYPE__ unsigned int
#define __UINT64_TYPE__ unsigned long long
#define __SIZE_TYPE__ unsigned long long
` + cc.Builtin,
	})
	if c.Predefines {
		srcs = append(srcs, cc.Source{Name: "<cxgo-builtin>", Value: "#include <" + libs.BuiltinH + ">\n"})
		//srcs = append(srcs, cc.Source{Name: "<cxgo-predef>", Value: fmt.Sprintf(gccPredefine, "int")})
	}
	srcs = append(srcs, c.Sources...)
	includes := addIncludeOverridePath(c.Includes)
	sysIncludes := addIncludeOverridePath(c.SysIncludes)
	return cc.Translate(&cc.Config{
		FS:              mergefs.Merge(newIncludeFS(env), os.DirFS("/")),
		ABI:             libcc.NewABI(env.Env),
		IncludePaths:    includes,
		SysIncludePaths: sysIncludes,
		PragmaHandler: func(toks []cc.Token) error {
			if len(toks) == 0 {
				return nil
			}
			name := toks[0].SrcStr()
			toks = toks[1:]
			switch name {
			case "push_macro":
				if len(toks) != 3 {
					return nil
				} else if toks[0].SrcStr() != tokLBrace || toks[2].SrcStr() != tokRBrace {
					return nil
				}
				def := toks[1].SrcStr()
				def, err := strconv.Unquote(def)
				if err != nil {
					return nil
				}
				// FIXME: push macro
				//p.PushMacro(def)
			case "pop_macro":
				if len(toks) != 3 {
					return nil
				} else if toks[0].SrcStr() != tokLBrace || toks[2].SrcStr() != tokRBrace {
					return nil
				}
				def := toks[1].SrcStr()
				def, err := strconv.Unquote(def)
				if err != nil {
					return nil
				}
				// FIXME: pop macro
				//p.PopMacro(def)
			}
			return nil
		},
	}, srcs)
}
