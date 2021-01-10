package types

// CommonType find common type to convert two operands to.
func (e *Env) CommonType(x, y Type) (otyp Type) {
	if x == y && !x.Kind().IsInt() && !y.Kind().IsInt() {
		return x
	}
	if x.Sizeof() < y.Sizeof() {
		return e.CommonType(y, x)
	}
	xk, yk := x.Kind(), y.Kind()
	defer func() {
		if otyp.Kind().IsUntyped() && (xk.IsInt() || yk.IsInt()) && (!xk.IsUntyped() || !yk.IsUntyped()) {
			panic("returning untyped")
		}
	}()

	def := e.DefIntT()
	if xk != yk && ((xk.IsInt() && !yk.IsInt()) || (!xk.IsInt() && yk.IsInt())) {
		if xk.IsInt() && yk.IsFloat() {
			return y
		} else if xk.IsFloat() && yk.IsInt() {
			return x
		}
		if xk.IsInt() && yk.IsUntyped() {
			return x
		}
		if xk.IsInt() && yk.IsInt() {
			xi, yi := Unwrap(x).(IntType), Unwrap(y).(IntType)
			if xi.Sizeof() == yi.Sizeof() && xi.Signed() == yi.Signed() {
				if _, ok := x.(Named); ok {
					return y // prefer unnamed types
				}
				return x
			}
			if xi.Sizeof() < def.Sizeof() && yi.Sizeof() < def.Sizeof() {
				return def // C implicit type conversion to int
			}
			if xk.IsUntyped() && yk.IsInt() {
				if x.Sizeof() == y.Sizeof() {
					return y
				}
				// y is smaller
				var xi Type
				if yk.Is(Signed) {
					xi = e.DefIntT()
				} else {
					xi = e.DefUintT()
				}
				return xi
			}
		}
	}

	switch x := x.(type) {
	case IntType:
		switch y := y.(type) {
		case IntType:
			if x.Kind().IsUntyped() {
				x = AsTypedIntT(x)
			}
			if y.Kind().IsUntyped() {
				y = AsTypedIntT(y)
			}
			if x.Sizeof() < def.Sizeof() && y.Sizeof() < def.Sizeof() {
				return def // C implicit type conversion to int
			}
			if x.Signed() == y.Signed() {
				// same sign -> pick largest
				if x.Sizeof() >= y.Sizeof() {
					return x
				}
				return y
			}
			// make X always correspond to unsigned
			if !y.Signed() {
				x, y = y, x
			}
			if x.Sizeof() >= y.Sizeof() {
				// is unsigned is larger or equal - prefer it
				return x
			}
			return y
		case BoolType:
			// int+bool = int
			if x.Kind().IsUntyped() {
				return e.DefIntT()
			}
			return x
		}
	case BoolType:
		switch y := y.(type) {
		case IntType:
			// bool+int = int
			if y.Kind().IsUntyped() {
				return e.DefIntT()
			}
			return y
		}
	case ArrayType:
		switch y.(type) {
		case IntType:
			e := e.PtrT(x.Elem())
			return e
		}
	}
	if x.Kind().IsUntyped() {
		return y
	} else if y.Kind().IsUntyped() {
		return x
	}
	// TODO
	return x
}
