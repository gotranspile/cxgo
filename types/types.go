package types

import (
	"bytes"
	"fmt"
	"sync"
)

type Node interface {
	Uses() []Usage
}

type Type interface {
	Sizeof() int
	Kind() Kind
	GoType() GoType
}

func UnkT(size int) Type {
	if size <= 0 {
		panic("size must be specified; set to ptr size, for example")
	}
	return &unkType{size: size}
}

type unkType struct {
	isStruct bool
	size     int
}

func (t *unkType) Kind() Kind {
	return Unknown
}

func (t *unkType) Sizeof() int {
	return t.size
}

var (
	uint8Type  = IntType{size: 1, signed: false}
	uint16Type = IntType{size: 2, signed: false}
	uint32Type = IntType{size: 4, signed: false}
	uint64Type = IntType{size: 8, signed: false}
	int8Type   = IntType{size: 1, signed: true}
	int16Type  = IntType{size: 2, signed: true}
	int32Type  = IntType{size: 4, signed: true}
	int64Type  = IntType{size: 8, signed: true}
)

func UntypedIntT(minSize int) IntType {
	if minSize <= 0 {
		panic("size must be specified")
	}
	return IntType{size: minSize, untyped: true}
}

func AsUntypedIntT(t IntType) IntType {
	t.untyped = true
	return t
}

func AsTypedIntT(t IntType) IntType {
	t.untyped = false
	return t
}

func ToPtrType(exp Type) PtrType {
	if t, ok := Unwrap(exp).(PtrType); ok {
		return t
	}
	return nil
}

func IntT(size int) IntType {
	if size <= 0 {
		panic("size must be specified")
	}
	switch size {
	case 1:
		return int8Type
	case 2:
		return int16Type
	case 4:
		return int32Type
	case 8:
		return int64Type
	}
	return IntType{size: size, signed: true}
}

func UintT(size int) IntType {
	if size <= 0 {
		panic("size must be specified")
	}
	switch size {
	case 1:
		return uint8Type
	case 2:
		return uint16Type
	case 4:
		return uint32Type
	case 8:
		return uint64Type
	}
	return IntType{size: size, signed: false}
}

var (
	float32Type = FloatType{size: 4}
	float64Type = FloatType{size: 8}
)

func FloatT(size int) FloatType {
	switch size {
	case 4:
		return float32Type
	case 8:
		return float64Type
	}
	return FloatType{size: size}
}

func AsUntypedFloatT(t FloatType) FloatType {
	t.untyped = true
	return t
}

func AsTypedFloatT(t FloatType) FloatType {
	t.untyped = false
	return t
}

func NilT(size int) PtrType {
	return &ptrType{size: size}
}

type PtrType interface {
	Type
	Elem() Type
	SetElem(e Type)
	ElemKind() Kind
	ElemSizeof() int
}

func PtrT(size int, elem Type) PtrType {
	if size <= 0 {
		panic("size must be set")
	} else if size < 4 {
		panic("unlikely")
	}
	return &ptrType{elem: elem, size: size}
}

var _ PtrType = &ptrType{}

type ptrType struct {
	size int
	zero bool
	elem Type
}

func (t *ptrType) Sizeof() int {
	return t.size
}

func (t *ptrType) Kind() Kind {
	if t.zero {
		return Nil
	} else if t.elem == nil {
		return UnsafePtr
	}
	return Ptr
}

func (t *ptrType) Elem() Type {
	return t.elem
}

func (t *ptrType) SetElem(e Type) {
	t.elem = e
}

func (t *ptrType) ElemKind() Kind {
	e := t.elem
	for e != nil {
		p, ok := e.(*ptrType)
		if !ok {
			break
		}
		e = p.elem
	}
	if e == nil {
		return UnsafePtr
	}
	return e.Kind()
}

func (t *ptrType) ElemSizeof() int {
	if t.elem == nil {
		return 1
	}
	return t.elem.Sizeof()
}

var (
	_ PtrType = namedPtr{}
	_ Named   = namedPtr{}
)

type namedPtr struct {
	name *Ident
	*ptrType
}

func (t namedPtr) Name() *Ident {
	return t.name
}

func (t namedPtr) Underlying() Type {
	return t.ptrType
}

func (t namedPtr) SetUnderlying(typ Type) Named {
	panic("trying to change the named type")
}

func (t namedPtr) Incomplete() bool {
	return false
}

type Field struct {
	Name *Ident
}

func (f *Field) Type() Type {
	return f.Name.CType(nil)
}

func FuncT(ptrSize int, ret Type, args ...*Field) *FuncType {
	return funcT(ptrSize, ret, args, false)
}

func FuncTT(ptrSize int, ret Type, args ...Type) *FuncType {
	fields := make([]*Field, 0, len(args))
	for _, t := range args {
		fields = append(fields, &Field{Name: NewUnnamed(t)})
	}
	return FuncT(ptrSize, ret, fields...)
}

func VarFuncT(ptrSize int, ret Type, args ...*Field) *FuncType {
	return funcT(ptrSize, ret, args, true)
}

func VarFuncTT(ptrSize int, ret Type, args ...Type) *FuncType {
	fields := make([]*Field, 0, len(args))
	for _, t := range args {
		fields = append(fields, &Field{Name: NewUnnamed(t)})
	}
	return VarFuncT(ptrSize, ret, fields...)
}

func checkFields(fields []*Field) {
	for _, f := range fields {
		if f.Name == nil {
			panic("nil argument name")
		}
	}
}

func funcT(ptrSize int, ret Type, args []*Field, vari bool) *FuncType {
	if ptrSize <= 0 {
		panic("size must be set")
	} else if ptrSize < 4 {
		panic("unlikely")
	}
	checkFields(args)
	return &FuncType{
		size: ptrSize,
		args: append([]*Field{}, args...),
		ret:  ret,
		vari: vari,
	}
}

type FuncType struct {
	ptr  bool
	size int

	args []*Field
	ret  Type
	vari bool
}

func (t *FuncType) Kind() Kind {
	return Func
}

func (t *FuncType) Sizeof() int {
	return t.size
}

func (t *FuncType) Return() Type {
	return t.ret
}

func (t *FuncType) Variadic() bool {
	return t.vari
}

func (t *FuncType) ArgN() int {
	return len(t.args)
}

func (t *FuncType) Args() []*Field {
	return append([]*Field{}, t.args...)
}

type IntType struct {
	size    int
	signed  bool
	untyped bool
}

func (t IntType) Kind() Kind {
	s := Unsigned
	if t.signed {
		s = Signed
	}
	if t.untyped {
		return UntypedInt | s
	}
	return Int | s
}

func (t IntType) Sizeof() int {
	return t.size
}

func (t IntType) Signed() bool {
	return t.signed
}

type FloatType struct {
	size    int
	untyped bool
}

func (t FloatType) Kind() Kind {
	if t.untyped {
		return UntypedFloat
	}
	return Float
}

func (t FloatType) Sizeof() int {
	return t.size
}

func BoolT() Type {
	return BoolType{}
}

type BoolType struct{}

func (t BoolType) Kind() Kind {
	return Bool
}

func (t BoolType) Sizeof() int {
	return 1
}

func ArrayT(elem Type, size int) Type {
	if size < 0 {
		panic("negative size")
	}
	return ArrayType{
		elem:  elem,
		size:  size,
		slice: false,
	}
}

func SliceT(elem Type) Type {
	return ArrayType{
		elem:  elem,
		size:  0,
		slice: true,
	}
}

type ArrayType struct {
	elem  Type
	size  int
	slice bool
}

func (t ArrayType) Kind() Kind {
	return Array
}

func (t ArrayType) Elem() Type {
	return t.elem
}

func (t ArrayType) Len() int {
	if t.slice {
		return 0
	}
	return t.size
}

func (t ArrayType) IsSlice() bool {
	return t.slice
}

func (t ArrayType) Sizeof() int {
	sz := t.size
	if sz == 0 {
		sz = 1
	}
	return sz * t.elem.Sizeof()
}

func NamedT(name string, typ Type) Named {
	return NamedTGo(name, "", typ)
}

func NamedTGo(cname, goname string, typ Type) Named {
	if cname == "" {
		panic("name is not set")
	}
	if typ == nil {
		panic("type is not set")
	}
	switch typ := typ.(type) {
	case *ptrType:
		named := &namedPtr{ptrType: typ}
		named.name = NewIdentGo(cname, goname, named)
		return named
	}
	named := &namedType{typ: typ}
	named.name = NewIdentGo(cname, goname, named)
	return named
}

type Named interface {
	Type
	Name() *Ident
	Underlying() Type
}

type namedType struct {
	name *Ident
	typ  Type
}

func (t *namedType) String() string {
	return t.name.String()
}

func (t *namedType) Kind() Kind {
	return t.typ.Kind()
}

func (t *namedType) Name() *Ident {
	return t.name
}

func (t *namedType) Underlying() Type {
	return t.typ
}

func (t *namedType) Sizeof() int {
	return t.typ.Sizeof()
}

var (
	structMu    sync.RWMutex
	structTypes = make(map[string]*StructType)
	unionTypes  = make(map[string]*StructType)
)

func StructT(fields []*Field) *StructType {
	checkFields(fields)
	s := &StructType{
		fields: append([]*Field{}, fields...),
		union:  false,
	}
	h := s.hash()

	structMu.RLock()
	t, ok := structTypes[h]
	structMu.RUnlock()
	if ok {
		return t
	}

	structMu.Lock()
	defer structMu.Unlock()
	if t, ok := structTypes[h]; ok {
		return t
	}
	structTypes[h] = s
	return s
}

func UnionT(fields []*Field) *StructType {
	checkFields(fields)
	s := &StructType{
		fields: append([]*Field{}, fields...),
		union:  true,
	}
	h := s.hash()

	structMu.RLock()
	t, ok := unionTypes[h]
	structMu.RUnlock()
	if ok {
		return t
	}

	structMu.Lock()
	defer structMu.Unlock()
	if t, ok := unionTypes[h]; ok {
		return t
	}
	unionTypes[h] = s
	return s
}

type StructType struct {
	Where  string
	fields []*Field
	union  bool
}

func (t *StructType) hash() string {
	buf := bytes.NewBuffer(nil)
	for _, f := range t.fields {
		buf.WriteString(f.Name.Name)
		buf.WriteByte(0)
		fmt.Fprintf(buf, "%p", f.Type())
		buf.WriteByte(0)
	}
	return buf.String()
}

func (t *StructType) Fields() []*Field {
	return append([]*Field{}, t.fields...)
}

func (t *StructType) Kind() Kind {
	return Struct
}

func (t *StructType) Sizeof() int {
	if t.union {
		max := 0
		for _, f := range t.fields {
			if sz := f.Type().Sizeof(); sz > max {
				max = sz
			}
		}
		return max
	}
	n := 0
	for _, f := range t.fields {
		n += f.Type().Sizeof()
	}
	if n == 0 {
		n = 1
	}
	return n
}
