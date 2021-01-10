package cxgo

import (
	"go/ast"
	"go/token"
	"math"
	"strconv"
	"strings"

	"github.com/gotranspile/cxgo/types"
)

type Number interface {
	Expr
	IsZero() bool
	IsOne() bool
	Negate() Number
}

func parseCIntLit(s string) (IntLit, error) {
	s = strings.ToLower(s)
	s = strings.TrimRight(s, "ul")
	base := 10
	if strings.HasPrefix(s, "0x") {
		base = 16
		s = s[2:]
	} else if strings.HasPrefix(s, "0b") {
		base = 2
		s = s[2:]
	} else if len(s) > 1 && strings.HasPrefix(s, "0") {
		base = 8
		s = s[1:]
	}
	uv, err := strconv.ParseUint(s, base, 64)
	if err == nil {
		return cUintLit(uv), nil
	}
	iv, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return IntLit{}, err
	}
	return cIntLit(iv), nil
}

func cIntLit(v int64) IntLit {
	if v >= 0 {
		return cUintLit(uint64(v))
	}
	l := IntLit{val: uint64(-v), neg: true}
	if v >= math.MinInt8 && v <= math.MaxInt8 {
		l.typ = types.IntT(1)
	} else if v >= math.MinInt16 && v <= math.MaxInt16 {
		l.typ = types.IntT(2)
	} else if v >= math.MinInt32 && v <= math.MaxInt32 {
		l.typ = types.IntT(4)
	} else {
		l.typ = types.IntT(8)
	}
	return l
}

func cUintLit(v uint64) IntLit {
	l := IntLit{val: v}
	if v <= math.MaxUint8 {
		l.typ = types.AsUntypedIntT(types.UintT(1))
	} else if v <= math.MaxUint16 {
		l.typ = types.AsUntypedIntT(types.UintT(2))
	} else if v <= math.MaxUint32 {
		l.typ = types.AsUntypedIntT(types.UintT(4))
	} else {
		l.typ = types.AsUntypedIntT(types.UintT(8))
	}
	return l
}

func litCanStore(t types.IntType, v IntLit) bool {
	if v.neg {
		if !t.Signed() {
			return false
		}
		switch t.Sizeof() {
		case 1:
			return v.val <= uint64(-math.MinInt8)
		case 2:
			return v.val <= uint64(-math.MinInt16)
		case 4:
			return v.val <= uint64(-math.MinInt32)
		case 8:
			return v.val <= uint64(-math.MinInt64)
		}
		return true
	}
	if t.Signed() {
		switch t.Sizeof() {
		case 1:
			return v.val <= math.MaxInt8
		case 2:
			return v.val <= math.MaxInt16
		case 4:
			return v.val <= math.MaxInt32
		case 8:
			return v.val <= math.MaxInt64
		}
		return true
	}
	switch t.Sizeof() {
	case 1:
		return v.val <= math.MaxUint8
	case 2:
		return v.val <= math.MaxUint16
	case 4:
		return v.val <= math.MaxUint32
	case 8:
		return v.val <= math.MaxUint64
	}
	return true
}

var _ Number = IntLit{}

type IntLit struct {
	typ types.IntType
	val uint64
	neg bool
}

func (IntLit) Visit(v Visitor) {}

func (l IntLit) String() string {
	v := strconv.FormatUint(l.val, 10)
	if l.neg {
		return "-" + v
	}
	return v
}

func (l IntLit) CType(exp types.Type) types.Type {
	if t, ok := types.Unwrap(exp).(types.IntType); ok {
		if !t.Signed() && !l.neg && l.typ.Sizeof() <= t.Sizeof() {
			return exp
		} else if t.Signed() && l.typ.Sizeof() <= t.Sizeof() && litCanStore(t, l) {
			return exp
		}
	}
	return l.typ
}

func (IntLit) IsConst() bool {
	return true
}

func (IntLit) HasSideEffects() bool {
	return false
}

func (l IntLit) IsZero() bool {
	return l.val == 0
}

func (l IntLit) IsOne() bool {
	return l.val == 1
}

func (l IntLit) Negate() Number {
	return l.NegateLit()
}

func (l IntLit) NegateLit() IntLit {
	if l.neg {
		return cUintLit(l.val)
	}
	if l.val > math.MaxInt64 {
		panic("cannot negate")
	}
	return cIntLit(-int64(l.val))
}

func (l IntLit) MulLit(v int64) IntLit {
	if v < 0 {
		l.neg = !l.neg
		v = -v
	}
	l.val *= uint64(v)
	return l
}

func (l IntLit) IsUint() bool {
	return !l.neg || l.val > math.MaxInt64
}

func (l IntLit) IsNeg() bool {
	return l.neg
}

func (l IntLit) Int() int64 {
	if l.val > math.MaxInt64 {
		panic("value is too big!")
	}
	v := int64(l.val)
	if l.neg {
		return -v
	}
	return v
}

func (l IntLit) Uint() uint64 {
	if l.neg {
		panic("value is negative!")
	}
	return l.val
}

func (l IntLit) OverflowInt(sz int) IntLit {
	switch sz {
	case 8:
		v := int64(l.Uint())
		return cIntLit(v)
	case 4:
		v := int32(uint32(l.Uint()))
		return cIntLit(int64(v))
	case 2:
		v := int16(uint16(l.Uint()))
		return cIntLit(int64(v))
	case 1:
		v := int8(uint8(l.Uint()))
		return cIntLit(int64(v))
	}
	return l
}

func (l IntLit) OverflowUint(sz int) IntLit {
	switch sz {
	case 8:
		v := uint64(l.Int())
		return cUintLit(v)
	case 4:
		v := uint32(int32(l.Int()))
		return cUintLit(uint64(v))
	case 2:
		v := uint16(int16(l.Int()))
		return cUintLit(uint64(v))
	case 1:
		v := uint8(int8(l.Int()))
		return cUintLit(uint64(v))
	}
	return l
}

func (l IntLit) AsExpr() GoExpr {
	if l.neg {
		val := -int64(l.val)
		switch val {
		case math.MinInt64:
			return ident("math.MinInt64")
		case math.MinInt32:
			return ident("math.MinInt32")
		case math.MinInt16:
			return ident("math.MinInt16")
		case math.MinInt8:
			return ident("math.MinInt8")
		}
		return intLit64(val)
	}
	switch l.val {
	case math.MaxUint64:
		return ident("math.MaxUint64")
	case math.MaxUint32:
		return ident("math.MaxUint32")
	case math.MaxUint16:
		return ident("math.MaxUint16")
	case math.MaxUint8:
		return ident("math.MaxUint8")
	case math.MaxInt64:
		return ident("math.MaxInt64")
	case math.MaxInt32:
		return ident("math.MaxInt32")
	case math.MaxInt16:
		return ident("math.MaxInt16")
	case math.MaxInt8:
		return ident("math.MaxInt8")
	}
	return uintLit64(l.val)
}

func (l IntLit) Uses() []types.Usage {
	return nil
}

var _ Number = FloatLit{}

func parseCFloatLit(s string) (FloatLit, error) {
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, "f")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return FloatLit{}, err
	}
	return FloatLit{val: v}, nil
}

type FloatLit struct {
	val float64
}

func (FloatLit) Visit(v Visitor) {}

func (l FloatLit) CType(exp types.Type) types.Type {
	if t, ok := types.Unwrap(exp).(types.FloatType); ok {
		return t
	}
	return types.FloatT(8)
}

func (l FloatLit) AsExpr() GoExpr {
	s := strconv.FormatFloat(l.val, 'g', -1, 64)
	if float64(int(l.val)) == l.val && !strings.Contains(s, ".") {
		s += ".0"
	}
	return &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: s,
	}
}

func (l FloatLit) IsZero() bool {
	return l.val == 0.0
}

func (l FloatLit) IsOne() bool {
	return l.val == 1.0
}

func (l FloatLit) Negate() Number {
	return FloatLit{val: -l.val}
}

func (l FloatLit) IsConst() bool {
	return true
}

func (l FloatLit) HasSideEffects() bool {
	return false
}

func (l FloatLit) Uses() []types.Usage {
	return nil
}

func (g *translator) parseCStringLit(s string) (StringLit, error) {
	return StringLit{typ: g.env.Go().String(), val: s}, nil
}

func (g *translator) parseCWStringLit(s string) (StringLit, error) {
	v, err := g.parseCStringLit(s)
	if err == nil {
		v.wide = true
	}
	return v, err
}

var _ Expr = StringLit{}

type StringLit struct {
	typ  types.Type
	val  string
	wide bool
}

func (StringLit) Visit(v Visitor) {}

func (l StringLit) String() string {
	return strconv.Quote(l.val)
}

func (l StringLit) Value() string {
	return l.val
}

func (l StringLit) CType(types.Type) types.Type {
	return l.typ
}

func (l StringLit) AsExpr() GoExpr {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(l.val),
	}
}

func (l StringLit) IsWide() bool {
	return l.wide
}

func (l StringLit) IsConst() bool {
	return true
}

func (l StringLit) HasSideEffects() bool {
	return false
}

func (l StringLit) Uses() []types.Usage {
	return nil
}

func isASCII(b byte) bool {
	return b <= '~'
}

type CLitKind int

const (
	CLitChar = CLitKind(iota)
	CLitWChar
)

func cLit(value string, kind CLitKind) Expr {
	return &CLiteral{
		Value: value,
		Kind:  kind,
	}
}

func cLitT(value string, kind CLitKind, typ types.Type) Expr {
	if typ == nil {
		panic("use cLit")
	}
	return cLit(value, kind)
}

type CLiteral struct {
	Value string
	Kind  CLitKind
	Type  types.Type
}

func (*CLiteral) Visit(v Visitor) {}

func (e *CLiteral) CType(types.Type) types.Type {
	if e.Type != nil {
		return e.Type
	}
	switch e.Kind {
	case CLitChar:
		return types.IntT(1)
	case CLitWChar:
		panic("TODO")
	default:
		panic(e.Kind)
	}
}

func (e *CLiteral) IsConst() bool {
	return true
}

func (e *CLiteral) HasSideEffects() bool {
	return false
}

func (e *CLiteral) AsExpr() GoExpr {
	lit := &ast.BasicLit{
		Value: e.Value,
	}
	switch e.Kind {
	case CLitChar, CLitWChar: // FIXME
		r := []rune(lit.Value)
		if len(r) != 1 {
			panic(strconv.Quote(lit.Value))
		}
		lit.Kind = token.CHAR
		if isASCII(byte(r[0])) {
			lit.Value = quoteWith(lit.Value, '\'')
		} else {
			lit.Value = "'\\x" + strconv.FormatUint(uint64(r[0]), 16) + "'"
		}
	default:
		panic(e.Kind)
	}
	return lit
}

func (e *CLiteral) Uses() []types.Usage {
	return nil
}
