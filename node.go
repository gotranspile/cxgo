package cxgo

type Visitor func(n Node)

type Node interface {
	// Visit calls the provided interface for all child nodes.
	// It's the visitor's responsibility to recurse by calling n.Visit(v).
	Visit(v Visitor)
}
