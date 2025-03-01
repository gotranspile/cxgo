package cxgo

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"modernc.org/cc/v3"

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

var (
	tokLBrace = cc.String("(")
	tokRBrace = cc.String(")")
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

func newCCConfigs(env *libs.Env, c ParseConfig) (*cc.Config, []string, []string, []cc.Source) {
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
		srcs = append(srcs, cc.Source{Name: "cxgo_config_defines.h", Value: buf.String()})
	}
	if c.Predefines {
		srcs = append(srcs, cc.Source{Name: "cxgo_predef.h", Value: fmt.Sprintf(gccPredefine, "int")})
	}
	srcs = append(srcs, c.Sources...)
	includes := addIncludeOverridePath(c.Includes)
	sysIncludes := addIncludeOverridePath(c.SysIncludes)
	cconf := &cc.Config{
		Config3: cc.Config3{
			WorkingDir: c.WorkDir,
			Filesystem: cc.Overlay(cc.LocalFS(), newIncludeFS(env)),
		},
		ABI: libcc.NewABI(env.Env),
		PragmaHandler: func(p cc.Pragma, toks []cc.Token) {
			if len(toks) == 0 {
				return
			}
			name := toks[0].Value.String()
			toks = toks[1:]
			switch name {
			case "push_macro":
				if len(toks) != 3 {
					return
				} else if toks[0].Value != tokLBrace || toks[2].Value != tokRBrace {
					return
				}
				def := toks[1].Value.String()
				def, err := strconv.Unquote(def)
				if err != nil {
					return
				}
				p.PushMacro(def)
			case "pop_macro":
				if len(toks) != 3 {
					return
				} else if toks[0].Value != tokLBrace || toks[2].Value != tokRBrace {
					return
				}
				def := toks[1].Value.String()
				def, err := strconv.Unquote(def)
				if err != nil {
					return
				}
				p.PopMacro(def)
			}
		},
	}
	return cconf, includes, sysIncludes, srcs
}

func PreprocessSource(w io.Writer, env *libs.Env, c ParseConfig) error {
	cconf, includes, sysIncludes, srcs := newCCConfigs(env, c)
	err := cc.Preprocess(cconf, includes, sysIncludes, srcs, w)
	return err
}

func ParseSource(env *libs.Env, c ParseConfig) (*cc.AST, error) {
	cconf, includes, sysIncludes, srcs := newCCConfigs(env, c)
	return cc.Translate(cconf, includes, sysIncludes, srcs)
}
