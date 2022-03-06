package libs

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/gotranspile/cxgo/types"
)

func (c *Env) NewIdent(cname, goname string, v interface{}, typ types.Type) *types.Ident {
	err := c.checkType(reflect.TypeOf(v), typ)
	if err != nil {
		panic(fmt.Errorf("unexpected type for %s: %w", goname, err))
	}
	return types.NewIdentGo(cname, goname, typ)
}

func (c *Env) checkErr(ok bool, exp reflect.Type, got types.Type) error {
	if ok {
		return nil
	}
	return fmt.Errorf("expected %v (%v), got: %T, %v (%+v)", exp, exp.Kind(), got, got.Kind(), got)
}

var (
	rtUnsafe = reflect.TypeOf(unsafe.Pointer(nil))
)

func (c *Env) checkType(t1 reflect.Type, t2 types.Type) error {
	if t1 != rtUnsafe && strings.Contains(t1.PkgPath(), ".") {
		name1 := t1.Name()
		nt2, ok := t2.(types.Named)
		if !ok {
			return fmt.Errorf("expected named type, got: %T", t2)
		}
		id := nt2.Name()
		name2 := id.GoName
		if name2 == "" {
			name2 = id.Name
		}
		if name2 == "" {
			return fmt.Errorf("expected type with a name, got: %+v", id)
		}
		if i := strings.LastIndex(name2, "."); i > 0 {
			name2 = name2[i+1:]
		}
		if name1 != name2 {
			return fmt.Errorf("expected type named %q, got: %q", name1, name2)
		}
		return nil
	}
	switch t1.Kind() {
	case reflect.Bool:
		return c.checkErr(t2 == types.BoolT(), t1, t2)
	case reflect.Int:
		return c.checkErr(t2 == c.Go().Int() || (t2.Kind().IsInt() && t2.Kind().IsUntyped()), t1, t2)
	case reflect.Uint:
		return c.checkErr(t2 == c.Go().Uint(), t1, t2)
	case reflect.UnsafePointer:
		return c.checkErr(t2 == c.Go().UnsafePtr() || types.Same(t2, c.PtrT(nil)), t1, t2)
	case reflect.Uintptr:
		return c.checkErr(t2 == c.Go().Uintptr(), t1, t2)
	case reflect.String:
		return c.checkErr(t2 == c.Go().String(), t1, t2)
	case reflect.Float32:
		return c.checkErr(t2 == types.FloatT(4), t1, t2)
	case reflect.Float64:
		return c.checkErr(t2 == types.FloatT(8), t1, t2)
	case reflect.Int8:
		return c.checkErr(t2 == types.IntT(1), t1, t2)
	case reflect.Uint8:
		return c.checkErr(t2 == types.UintT(1) || types.Same(t2, c.Go().Byte()), t1, t2)
	case reflect.Int16:
		return c.checkErr(t2 == types.IntT(2), t1, t2)
	case reflect.Uint16:
		if nt, ok := t2.(types.Named); ok && nt.Name().GoName == "libc.WChar" {
			return nil
		}
		return c.checkErr(t2 == types.UintT(2), t1, t2)
	case reflect.Int32:
		return c.checkErr(t2 == types.IntT(4) || types.Same(t2, c.Go().Rune()), t1, t2)
	case reflect.Uint32:
		return c.checkErr(t2 == types.UintT(4), t1, t2)
	case reflect.Int64:
		return c.checkErr(t2 == types.IntT(8), t1, t2)
	case reflect.Uint64:
		return c.checkErr(t2 == types.UintT(8), t1, t2)
	case reflect.Ptr:
		p, ok := t2.(types.PtrType)
		if !ok {
			return fmt.Errorf("expected pointer, got: %T", t2)
		}
		return c.checkType(t1.Elem(), p.Elem())
	case reflect.Func:
		f, ok := t2.(*types.FuncType)
		if !ok {
			return fmt.Errorf("expected func, got: %T", t2)
		}
		if t1.NumOut() == 0 {
			if f.Return() != nil {
				return fmt.Errorf("expected void func, got: %T", f.Return())
			}
		} else if t1.NumOut() == 1 {
			err := c.checkType(t1.Out(0), f.Return())
			if err != nil {
				return fmt.Errorf("unexpected return: %w", err)
			}
		} else {
			return fmt.Errorf("func with multiple returns: %v", t1)
		}
		if t1.IsVariadic() != f.Variadic() {
			return fmt.Errorf("variadic flags are different")
		}
		exp := t1.NumIn()
		if t1.IsVariadic() {
			exp--
		}
		if exp != f.ArgN() {
			return fmt.Errorf("unexpected number of arguments: %d vs %d", exp, f.ArgN())
		}
		args := f.Args()
		for i := 0; i < exp; i++ {
			err := c.checkType(t1.In(i), args[i].Type())
			if err != nil {
				return fmt.Errorf("arg %d: %w", i, err)
			}
		}
		return nil
	case reflect.Struct:
		_, ok := types.Unwrap(t2).(*types.StructType)
		if !ok {
			return fmt.Errorf("expected struct, got: %T", t2)
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %v (%v)", t1, t1.Kind())
	}
}
