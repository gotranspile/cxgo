package types

func Unwrap(t Type) Type {
	for {
		named, ok := t.(Named)
		if !ok {
			return t
		}
		t = named.Underlying()
	}
}

func UnwrapPtr(t Type) PtrType {
	for {
		switch tt := t.(type) {
		case PtrType:
			return tt
		case Named:
			t = tt.Underlying()
		default:
			return nil
		}
	}
}

func IsFuncPtr(t Type) (*FuncType, bool) {
	t = Unwrap(t)
	ptr, ok := t.(PtrType)
	if !ok {
		return nil, false
	}
	fnc, ok := ptr.Elem().(*FuncType)
	return fnc, ok
}

func IsUnsafePtr(t Type) bool {
	return t.Kind().IsUnsafePtr()
}

func IsUnsafePtr2(t Type) bool {
	if p, ok := t.(PtrType); ok && IsUnsafePtr(p.Elem()) {
		return true
	}
	return false
}

func IsPtr(t Type) bool {
	return t.Kind().IsPtr()
}

func IsInt(t Type) bool {
	return t.Kind().IsInt()
}

func sameInt(x, y IntType) bool {
	if x.untyped && y.untyped {
		return x.signed == y.signed
	} else if x.untyped {
		if x.size > y.size {
			return false
		}
		return x.signed == y.signed
	} else if y.untyped {
		if y.size > x.size {
			return false
		}
		return x.signed == y.signed
	}
	return x.Sizeof() == y.Sizeof() && x.Signed() == y.Signed()
}

func Same(x, y Type) bool {
	if (x == nil && y != nil) || (x != nil && y == nil) {
		return false
	} else if _, ok := x.(*unkType); ok {
		return false
	} else if _, ok := y.(*unkType); ok {
		return false
	} else if x == y {
		return true
	}
	switch x := x.(type) {
	case IntType:
		y, ok := y.(IntType)
		if !ok {
			return false
		}
		return sameInt(x, y)
	case PtrType:
		y, ok := y.(PtrType)
		return ok && Same(x.Elem(), y.Elem())
	case ArrayType:
		y, ok := y.(ArrayType)
		return ok && x.size == y.size && x.slice == y.slice && Same(x.elem, y.elem)
	case *namedType:
		y, ok := y.(*namedType)
		return ok && x.name == y.name && Same(x.typ, y.typ)
	case *FuncType:
		y, ok := y.(*FuncType)
		if !ok || !Same(x.Return(), y.Return()) {
			return false
		}
		xargs, yargs := x.Args(), y.Args()
		if len(xargs) != len(yargs) {
			return false
		}
		for i := range xargs {
			if !Same(xargs[i].Type(), yargs[i].Type()) {
				return false
			}
		}
		return true
	}
	return x == y
}
