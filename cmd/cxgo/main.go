package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/gotranspile/cxgo"
	"github.com/gotranspile/cxgo/internal/git"
	"github.com/gotranspile/cxgo/libs"
	"github.com/gotranspile/cxgo/types"
)

var Root = &cobra.Command{
	Use:   "cxgo",
	Short: "transpile a C project to Go",
	RunE:  run,
}

var (
	version = "dev"
	commit  = ""
	date    = ""
)

var configPath = "cxgo.yml"

func printVersion() {
	vers := version
	if s := commit; s != "" {
		vers = fmt.Sprintf("%s (%s)", vers, s[:8])
	}
	fmt.Printf("version: %s\n", vers)
	if date != "" {
		fmt.Printf("built: %s\n", date)
	}
}

func init() {
	Root.Flags().StringVarP(&configPath, "config", "c", configPath, "config file path")
	Root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print cxgo version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	})
}

func main() {
	if err := Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Replacement struct {
	Old string `yaml:"old"`
	Re  string `yaml:"regexp"`
	New string `yaml:"new"`
}

func (r Replacement) Build() (*cxgo.Replacer, error) {
	if r.Re == "" && r.Old == "" {
		return nil, errors.New("either 'regexp' or 'old' must be set")
	}
	var re *regexp.Regexp
	if r.Re != "" {
		reg, err := regexp.Compile(r.Re)
		if err != nil {
			return nil, err
		}
		re = reg
	}
	return &cxgo.Replacer{
		Old: r.Old,
		Re:  re,
		New: r.New,
	}, nil
}

type SrcFile struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
	Perm    int    `yaml:"perm"`
}

type File struct {
	Disabled    bool               `yaml:"disabled"`
	Name        string             `yaml:"name"`
	Content     string             `yaml:"content"`
	Predef      string             `yaml:"predef"`
	GoFile      string             `yaml:"go"`
	FlattenAll  *bool              `yaml:"flatten_all"`
	ForwardDecl *bool              `yaml:"forward_decl"`
	MaxDecls    int                `yaml:"max_decl"`
	Skip        []string           `yaml:"skip"`
	Idents      []cxgo.IdentConfig `yaml:"idents"`
	Replace     []Replacement      `yaml:"replace"`
}

type Config struct {
	VCS        string            `yaml:"vcs"`
	Branch     string            `yaml:"branch"`
	Root       string            `yaml:"root"`
	Out        string            `yaml:"out"`
	Package    string            `yaml:"package"`
	Include    []string          `yaml:"include"`
	SysInclude []string          `yaml:"sys_include"`
	IncludeMap map[string]string `yaml:"include_map"`
	Hooks      bool              `yaml:"hooks"`
	Define     []cxgo.Define     `yaml:"define"`
	Predef     string            `yaml:"predef"`
	SubPackage bool              `yaml:"subpackage"`

	IntSize   int  `yaml:"int_size"`
	PtrSize   int  `yaml:"ptr_size"`
	WcharSize int  `yaml:"wchar_size"`
	UseGoInt  bool `yaml:"use_go_int"`

	ForwardDecl      bool               `yaml:"forward_decl"`
	FlattenAll       bool               `yaml:"flatten_all"`
	FlattenFunc      []string           `yaml:"flatten"`
	Skip             []string           `yaml:"skip"`
	Replace          []Replacement      `yaml:"replace"`
	Idents           []cxgo.IdentConfig `yaml:"idents"`
	ImplicitReturns  bool               `yaml:"implicit_returns"`
	IgnoreIncludeDir bool               `yaml:"ignore_include_dir"`
	UnexportedFields bool               `yaml:"unexported_fields"`
	IntReformat      bool               `yaml:"int_reformat"`
	KeepFree         bool               `yaml:"keep_free"`
	NoLibs           bool               `yaml:"no_libs"`
	DoNotEdit        bool               `yaml:"do_not_edit"`

	SrcFiles []*SrcFile `yaml:"src_files"`
	FilePref string     `yaml:"file_pref"`
	Files    []*File    `yaml:"files"`

	ExecBefore []string `yaml:"exec_before"`
	ExecAfter  []string `yaml:"exec_after"`
}

func mergeBool(val *bool, def bool) bool {
	if val == nil {
		return def
	}
	return *val
}

func run(cmd *cobra.Command, args []string) error {
	defer cxgo.CallFinals()
	conf, _ := cmd.Flags().GetString("config")
	data, err := os.ReadFile(conf)
	if err != nil {
		return err
	}
	var c Config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return err
	}
	return Run(filepath.Dir(conf), &c)
}

func Run(root string, c *Config) error {
	if c.VCS != "" {
		name := strings.TrimSuffix(c.VCS, ".git")
		if i := strings.LastIndex(name, "/"); i > 0 {
			name = name[i:]
		}
		dir := filepath.Join(os.TempDir(), name)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			log.Printf("clonning %s to %s", c.VCS, dir)
			if err := git.Clone(c.VCS, c.Branch, dir); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			log.Printf("already cloned %s to %s", c.Root, dir)
		}
		c.Root = filepath.Join(dir, c.Root)
	}
	if !filepath.IsAbs(c.Root) {
		c.Root = filepath.Join(root, c.Root)
	}
	for _, f := range c.SrcFiles {
		if f.Name == "" {
			return errors.New("src_files entry with no name")
		}
		perm := os.FileMode(f.Perm)
		if perm == 0 {
			perm = 0644
		}
		log.Printf("writing %q (%o)", f.Name, perm)
		if err := os.WriteFile(filepath.Join(c.Root, f.Name), []byte(f.Content), perm); err != nil {
			return err
		}
	}
	if !filepath.IsAbs(c.Out) {
		c.Out = filepath.Join(root, c.Out)
		if abs, err := filepath.Abs(c.Out); err == nil {
			c.Out = abs
		}
	}
	log.Printf("writing to %s", c.Out)
	if err := os.MkdirAll(c.Out, 0755); err != nil {
		return err
	}
	tconf := types.Default()
	if c.UseGoInt {
		tconf.UseGoInt = c.UseGoInt
	}
	if c.IntSize != 0 {
		tconf.IntSize = c.IntSize
	}
	if c.PtrSize != 0 {
		tconf.PtrSize = c.PtrSize
	}
	if c.WcharSize != 0 {
		tconf.WCharSize = c.WcharSize
	}
	for i := range c.Include {
		if filepath.IsAbs(c.Include[i]) {
			continue
		}
		c.Include[i] = filepath.Join(c.Root, c.Include[i])
	}
	for i := range c.SysInclude {
		if filepath.IsAbs(c.SysInclude[i]) {
			continue
		}
		c.SysInclude[i] = filepath.Join(c.Root, c.SysInclude[i])
	}
	seen := make(map[string]struct{})
	processFile := func(f *File) error {
		if _, ok := seen[f.Name]; ok {
			return fmt.Errorf("duplicate entry for file: %q", f.Name)
		}
		seen[f.Name] = struct{}{}
		if base, ok := strings.CutSuffix(f.Name, ".h"); ok {
			if _, ok := seen[base+".c"]; ok {
				log.Println("skipping", f.Name) // Already included declarations from it.
				return nil
			}
		}
		if f.Content != "" {
			data := []byte(f.Content)
			if fdata, err := format.Source(data); err == nil {
				data = fdata
			}
			return os.WriteFile(filepath.Join(c.Out, f.Name), data, 0644)
		}
		idents := make(map[string]cxgo.IdentConfig)
		for _, v := range c.Idents {
			idents[v.Name] = v
		}
		for _, v := range f.Idents {
			idents[v.Name] = v
		}
		ilist := make([]cxgo.IdentConfig, 0, len(idents))
		for _, v := range idents {
			ilist = append(ilist, v)
		}

		env := libs.NewEnv(tconf)
		fc := cxgo.Config{
			Root:               c.Root,
			Package:            c.Package,
			GoFile:             f.GoFile,
			GoFilePref:         c.FilePref,
			FlattenAll:         mergeBool(f.FlattenAll, c.FlattenAll),
			ForwardDecl:        mergeBool(f.ForwardDecl, c.ForwardDecl),
			MaxDecls:           -1,
			Hooks:              c.Hooks,
			Define:             c.Define,
			Predef:             f.Predef,
			Idents:             ilist,
			Include:            c.Include,
			SysInclude:         c.SysInclude,
			IncludeMap:         c.IncludeMap,
			FixImplicitReturns: c.ImplicitReturns,
			IgnoreIncludeDir:   c.IgnoreIncludeDir,
			UnexportedFields:   c.UnexportedFields,
			IntReformat:        c.IntReformat,
			KeepFree:           c.KeepFree,
			DoNotEdit:          c.DoNotEdit,
		}
		env.NoLibs = c.NoLibs
		env.Map = c.IncludeMap
		if f.MaxDecls > 0 {
			fc.MaxDecls = f.MaxDecls
		}
		if fc.Predef == "" {
			fc.Predef = c.Predef
		}
		for _, r := range f.Replace {
			rp, err := r.Build()
			if err != nil {
				return err
			}
			fc.Replace = append(fc.Replace, *rp)
		}
		for _, r := range c.Replace {
			rp, err := r.Build()
			if err != nil {
				return err
			}
			fc.Replace = append(fc.Replace, *rp)
		}
		if len(f.Skip) != 0 {
			fc.SkipDecl = make(map[string]bool)
			for _, s := range f.Skip {
				fc.SkipDecl[s] = true
			}
		}
		log.Println(f.Name)
		if err := cxgo.Translate(c.Root, filepath.Join(c.Root, f.Name), c.Out, env, fc); err != nil {
			return err
		}
		return nil
	}
	if err := runCmd(c.Root, c.ExecBefore); err != nil {
		return err
	}
	for _, f := range c.Files {
		if f.Disabled {
			seen[f.Name] = struct{}{}
			continue
		}
		if strings.Contains(f.Name, "*") {
			paths, err := doublestar.Glob(filepath.Join(c.Root, f.Name))
			if err != nil {
				return err
			}
			for _, path := range paths {
				rel, err := filepath.Rel(c.Root, path)
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				if _, ok := seen[rel]; ok {
					continue
				} else if _, ok = seen["./"+rel]; ok {
					continue
				}
				f2 := *f
				f2.Name = rel
				if err := processFile(&f2); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
			}
		} else {
			if err := processFile(f); err != nil {
				return fmt.Errorf("%s: %w", f.Name, err)
			}
		}
	}
	if !c.SubPackage {
		if _, err := os.Stat(filepath.Join(c.Out, "go.mod")); os.IsNotExist(err) {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, `module %s

go 1.19

require (
	%s %s
)
`, c.Package, libs.RuntimePackage, libs.RuntimePackageVers)
			if err := os.WriteFile(filepath.Join(c.Out, "go.mod"), buf.Bytes(), 0644); err != nil {
				return err
			}
		}
	}
	if err := runCmd(c.Out, c.ExecAfter); err != nil {
		return err
	}
	return nil
}

func runCmd(wd string, args []string) error {
	if len(args) == 0 {
		return nil
	}
	log.Printf("+ %s", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = wd
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
