package types

type AccessKind int

const (
	AccessUnknown = AccessKind(iota)
	AccessDefine
	AccessRead
	AccessWrite
)

type Usage struct {
	*Ident
	Access AccessKind
}

func useNodes(arr []Node, acc AccessKind) []Usage {
	var out []Usage
	for _, n := range arr {
		if n == nil {
			continue
		}
		for _, u := range n.Uses() {
			if u.Access == AccessUnknown {
				u.Access = acc
			}
			out = append(out, u)
		}
	}
	return out
}

func UseRead(n ...Node) []Usage {
	return useNodes(n, AccessRead)
}

func UseWrite(n ...Node) []Usage {
	return useNodes(n, AccessWrite)
}

func NewUnnamed(typ Type) *Ident {
	return &Ident{typ: typ}
}

func NewIdent(name string, typ Type) *Ident {
	if typ == nil {
		panic("must have a type; use UnkT if it's unknown")
	}
	return &Ident{Name: name, typ: typ}
}

func NewIdentGo(cname, goname string, typ Type) *Ident {
	if typ == nil {
		panic("must have a type; use UnkT if it's unknown")
	}
	return &Ident{Name: cname, GoName: goname, typ: typ}
}

type Ident struct {
	typ    Type
	Name   string
	GoName string
}

func (e *Ident) IsUnnamed() bool {
	return e.Name == "" && e.GoName == ""
}

func (e *Ident) String() string {
	if e.GoName != "" {
		return e.GoName
	}
	return e.Name
}

func (e *Ident) CType(exp Type) Type {
	if exp == nil {
		return e.typ
	}
	if tk := e.typ.Kind(); tk.IsUntyped() {
		if tk.IsInt() {
			if exp.Kind().IsInt() {
				return exp
			}
		}
	}
	return e.typ
}
