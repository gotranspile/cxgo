package types

import "strings"

type Kind int

const (
	Unknown = Kind(0)
	Untyped = Kind(1 << iota)
	unsafeKind
	Ptr
	Int
	Float
	Bool
	Struct
	Func
	Array
	Signed
	Unsigned

	UnsafePtr  = unsafeKind | Ptr
	UntypedInt = Untyped | Int
	Nil        = Untyped | Ptr
)

func (k Kind) Is(k2 Kind) bool {
	if k2 == Unknown {
		return k == k2
	}
	return k&k2 == k2
}

func (k Kind) IsUntyped() bool {
	return k.Is(Untyped)
}

func (k Kind) IsPtr() bool {
	return k.Is(Ptr)
}

func (k Kind) IsUnsafePtr() bool {
	return k.Is(UnsafePtr)
}

func (k Kind) IsFunc() bool {
	return k.Is(Func)
}

func (k Kind) IsRef() bool {
	return k.IsPtr() || k.IsFunc()
}

func (k Kind) IsInt() bool {
	return k.Is(Int)
}

func (k Kind) IsSigned() bool {
	return k.Is(Signed)
}

func (k Kind) IsUnsigned() bool {
	return k.Is(Unsigned)
}

func (k Kind) IsFloat() bool {
	return k.Is(Float)
}

func (k Kind) IsBool() bool {
	return k.Is(Bool)
}

func (k Kind) Major() Kind {
	return k & (Func | Ptr | Array | Struct | Int | Float | Bool)
}

var kindNames = []struct {
	Kind Kind
	Name string
}{
	{Nil, "Nil"},
	{UntypedInt, "UntypedInt"},
	{UnsafePtr, "UnsafePtr"},
	{Unsigned, "Unsigned"},
	{Signed, "Signed"},
	{Int, "Int"},
	{Float, "Float"},
	{Bool, "Bool"},
	{Untyped, "Untyped"},
	{Ptr, "Ptr"},
	{Array, "Array"},
	{Func, "Func"},
	{Struct, "Struct"},
}

func (k Kind) String() string {
	if k == 0 {
		return "Unknown"
	}
	var kinds []string
	for _, s := range kindNames {
		if k&s.Kind != 0 {
			kinds = append(kinds, s.Name)
			k &^= s.Kind
		}
	}
	return strings.Join(kinds, "|")
}
