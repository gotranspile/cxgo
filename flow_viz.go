package cxgo

import (
	"bytes"
	"go/format"
	"go/token"
	"strconv"
	"strings"

	"github.com/emicklei/dot"
)

func (cf *ControlFlow) dumpDot() string {
	g := dot.NewGraph()
	s := g.Node("n1").Label("begin")
	start := cf.dumpBlockDot(g, cf.Start, make(map[Block]dot.Node))
	g.Edge(s, start)
	return g.String()
}

func printExpr(e Expr) string {
	buf := bytes.NewBuffer(nil)
	fs := token.NewFileSet()
	if err := format.Node(buf, fs, e.AsExpr()); err != nil {
		panic(err)
	}
	return buf.String()
}

func dotAddNode(g *dot.Graph, id string, b Block) dot.Node {
	var (
		label string
		shape = "circle"
	)
	switch b := b.(type) {
	case nil:
		panic("must not be nil")
	case *CodeBlock:
		label = printStmts(b.Stmts)
		shape = "box"
	case *CondBlock:
		label = "if " + printExpr(b.Expr)
		shape = "hexagon"
	case *ReturnBlock:
		label = printStmts([]CStmt{b.CReturnStmt})
		shape = "box"
	case *SwitchBlock:
		label = "switch " + printExpr(b.Expr)
		shape = "trapezium"
	default:
		panic(b)
	}
	label = strings.ReplaceAll(label, "\t", "  ")
	return g.Node(id).Label(label).Attr("shape", shape)
}

func (cf *ControlFlow) dumpBlockDot(g *dot.Graph, b Block, seen map[Block]dot.Node) dot.Node {
	if n, ok := seen[b]; ok {
		return n
	}
	id := "n" + strconv.Itoa(len(seen)+2) // n1 is "begin" node
	n := dotAddNode(g, id, b)
	seen[b] = n
	for _, p := range b.PrevBlocks() {
		pn := cf.dumpBlockDot(g, p, seen)
		g.Edge(n, pn).Attr("color", "#99999955")
	}
	switch b := b.(type) {
	case *CondBlock:
		g.Edge(n, cf.dumpBlockDot(g, b.Then, seen)).Attr("color", "#00aa00")
		g.Edge(n, cf.dumpBlockDot(g, b.Else, seen)).Attr("color", "#aa0000")
	case *CodeBlock:
		if b.Next != nil {
			g.Edge(n, cf.dumpBlockDot(g, b.Next, seen))
		}
	case *SwitchBlock:
		for i, c := range b.Blocks {
			e := g.Edge(n, cf.dumpBlockDot(g, c, seen))
			if v := b.Cases[i]; v == nil {
				e.Attr("color", "#aa0000")
			} else {
				e.Label(printExpr(v))
			}
		}
	case *ReturnBlock:
	default:
		panic(b)
	}
	return n
}

func (cf *ControlFlow) dumpDomDot() string {
	g := dot.NewGraph()
	s := g.Node("n1").Label("begin")
	nodes := make(map[Block]dot.Node)
	cf.eachBlock(func(b Block) {
		id := "n" + strconv.Itoa(len(nodes)+2) // n1 is "begin" node
		n := dotAddNode(g, id, b)
		nodes[b] = n
	})
	g.Edge(s, nodes[cf.Start])
	cf.eachBlock(func(b Block) {
		d := cf.IDom(b)
		if d == nil {
			return
		}
		g.Edge(nodes[d], nodes[b])
	})
	return g.String()
}

func printStmts(stmts []CStmt) string {
	var out []GoStmt
	for _, s := range stmts {
		out = append(out, s.AsStmt()...)
	}
	out = fixLabels(out)

	buf := bytes.NewBuffer(nil)
	fs := token.NewFileSet()
	for i, s := range out {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if err := format.Node(buf, fs, s); err != nil {
			panic(err)
		}
	}
	return buf.String()
}
