package libs

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

const (
	mathH = "math.h"
)

func init() {
	const cpkg = "cmath"
	RegisterLibrary(mathH, func(c *Env) *Library {
		doubleT := types.FloatT(8)
		floatT := types.FloatT(4)
		var buf bytes.Buffer
		buf.WriteString("const double M_PI = 3.1415;\n")
		lib := &Library{
			Imports: map[string]string{
				cpkg:     RuntimePrefix + cpkg,
				"math":   "math",
				"math32": "github.com/chewxy/math32",
			},
			Idents: map[string]*types.Ident{
				"atan2": types.NewIdent("math.Atan2", c.FuncTT(doubleT, doubleT, doubleT)),
				"modf":  types.NewIdent(cpkg+".Modf", c.FuncTT(doubleT, doubleT, c.PtrT(doubleT))),
				"modff": types.NewIdent(cpkg+".Modff", c.FuncTT(floatT, floatT, c.PtrT(floatT))),
				"ldexp": types.NewIdent("math.Ldexp", c.FuncTT(doubleT, doubleT, c.Go().Int())),
				"fmod":  types.NewIdent("math.Mod", c.FuncTT(doubleT, doubleT, doubleT)),
				"M_PI":  types.NewIdent("math.Pi", doubleT),
			},
		}
		func2arg := func(pkg, name, cname string, arg types.Type, argc string) {
			fname := strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
			lib.Idents[cname] = types.NewIdentGo(cname, pkg+"."+fname, c.FuncTT(arg, arg))
			fmt.Fprintf(&buf, "%s %s(%s);\n", argc, cname, argc)
		}
		func3arg := func(pkg, name, cname string, arg types.Type, argc string) {
			fname := strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
			lib.Idents[cname] = types.NewIdentGo(cname, pkg+"."+fname, c.FuncTT(arg, arg, arg))
			fmt.Fprintf(&buf, "%s %s(%s, %s);\n", argc, cname, argc, argc)
		}
		func2dfc := func(cname, name string) {
			func2arg("math", name, cname, doubleT, "double")
			// TODO: add Round to maze.io/x/math32
			if cname == "round" {
				fmt.Fprintf(&buf, "#define %sf(x) %s(x)\n", name, name)
			} else {
				func2arg("math32", name, cname+"f", floatT, "float")
			}
		}
		func2df := func(name string) {
			func2dfc(name, name)
		}
		func3df := func(name string) {
			func3arg("math", name, name, doubleT, "double")
			func3arg("math32", name, name+"f", floatT, "float")
		}
		for _, name := range []string{
			"sin", "cos", "tan",
		} {
			for _, h := range []bool{false, true} {
				for _, a := range []bool{false, true} {
					ap := ""
					if a {
						ap = "a"
					}
					hs := ""
					if h {
						hs = "h"
					}
					func2df(ap + name + hs)
				}
			}
		}
		func2df("round")
		func2df("ceil")
		func2df("floor")
		func2dfc("fabs", "abs")
		func3df("pow")
		func2df("sqrt")
		func2df("exp")
		func2df("exp2")
		func2df("log")
		func2df("log10")
		func2df("log2")
		buf.WriteString("double atan2(double y, double x);\n")
		buf.WriteString("double modf(double x, double *iptr);\n")
		buf.WriteString("float modff(float value, float *iptr);\n")
		buf.WriteString("double ldexp(double x, _cxgo_go_int exp);\n")
		buf.WriteString("double fmod(double x, double exp);\n")
		buf.WriteString("int isnan(double x);\n")
		buf.WriteString("double frexp(double x, int* exp);\n")
		lib.Header = buf.String()
		return lib
	})
}
