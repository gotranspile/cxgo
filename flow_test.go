package cxgo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dennwc/cxgo/libs"
	"github.com/dennwc/cxgo/types"
	"github.com/stretchr/testify/require"
)

func numStmt(n int) CStmt {
	return NewCExprStmt1(&CallExpr{
		Fun:  FuncIdent{types.NewIdent("foo", types.UnkT(1))},
		Args: []Expr{cIntLit(int64(n))},
	})
}

func varDecl(n int) CStmt {
	return &CDeclStmt{&CVarDecl{
		CVarSpec: CVarSpec{
			Type: types.NamedT("bar", types.UnkT(1)),
			Names: []*types.Ident{
				types.NewIdent(fmt.Sprintf("foo%d", n), types.UnkT(1)),
			},
			Inits: []Expr{
				cIntLit(int64(n)),
			},
		},
	}}
}

func numCond(n int) BoolExpr {
	return BoolIdent{types.NewIdent(fmt.Sprintf("%d", n), types.BoolT())}
}

func ret(n int) CStmt {
	return &CReturnStmt{Expr: cIntLit(int64(n))}
}

func newBlock(stmts ...CStmt) *BlockStmt {
	return &BlockStmt{Stmts: stmts}
}

var casesControlFlow = []struct {
	name string
	tree []CStmt
	exp  string
	dom  string
	flat string
}{
	{
		name: "return",
		tree: []CStmt{
			ret(1),
		},
		exp: `
	n2[label="return 1",shape="box"];
	n1->n2;
`,
		dom: `
	n2[label="return 1",shape="box"];
	n1->n2;
`,
		flat: `
return 1
`,
	},
	{
		name: "no return",
		tree: []CStmt{
			numStmt(1),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="return",shape="box"];
	n1->n2;
	n2->n3;
	n3->n2[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="return",shape="box"];
	n1->n2;
	n2->n3;
`,
		flat: `
foo(1)
goto L_1
L_1:
return
`,
	},
	{
		name: "code and return",
		tree: []CStmt{
			numStmt(1),
			ret(2),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="return 2",shape="box"];
	n1->n2;
	n2->n3;
	n3->n2[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="return 2",shape="box"];
	n1->n2;
	n2->n3;
`,
		flat: `
foo(1)
goto L_1
L_1:
return 2
`,
	},
	{
		name: "code and return 2",
		tree: []CStmt{
			numStmt(1),
			numStmt(2),
			ret(3),
		},
		exp: `
	n2[label="foo(1)\nfoo(2)",shape="box"];
	n3[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n3->n2[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)\nfoo(2)",shape="box"];
	n3[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
`,
		flat: `
foo(1)
foo(2)
goto L_1
L_1:
return 3
`,
	},
	{
		name: "code with decls",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(varDecl(1))},
			varDecl(2),
			ret(1),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="var foo1 bar = 1",shape="box"];
	n4[label="var foo2 bar = 2",shape="box"];
	n5[label="return 1",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n3[color="#99999955"];
	n4->n2[color="#99999955"];
	n4->n5;
	n5->n4[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="var foo1 bar = 1",shape="box"];
	n4[label="var foo2 bar = 2",shape="box"];
	n5[label="return 1",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
`,
		flat: `
var foo1 bar
var foo2 bar
if 1 {
goto L_1
} else {
goto L_2
}
L_1:
foo1 = 1
goto L_2
L_2:
foo2 = 2
goto L_3
L_3:
return 1
`,
	},
	{
		name: "if",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(numStmt(2))},
			numStmt(3),
			numStmt(4),
			ret(5),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="foo(3)\nfoo(4)",shape="box"];
	n5[label="return 5",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n3[color="#99999955"];
	n4->n2[color="#99999955"];
	n4->n5;
	n5->n4[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="foo(3)\nfoo(4)",shape="box"];
	n5[label="return 5",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_2
}
L_1:
foo(2)
goto L_2
L_2:
foo(3)
foo(4)
goto L_3
L_3:
return 5
`,
	},
	{
		name: "if else",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(numStmt(2)), Else: newBlock(numStmt(3))},
			ret(4),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="return 4",shape="box"];
	n5[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n5[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n3[color="#99999955"];
	n4->n5[color="#99999955"];
	n5->n2[color="#99999955"];
	n5->n4;
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="return 4",shape="box"];
	n5[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n2->n5;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_3
}
L_1:
foo(2)
goto L_2
L_2:
return 4
L_3:
foo(3)
goto L_2
`,
	},
	{
		name: "if return else",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(ret(2)), Else: newBlock(numStmt(3))},
			ret(4),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="foo(3)",shape="box"];
	n5[label="return 4",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n4->n2[color="#99999955"];
	n4->n5;
	n5->n4[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="foo(3)",shape="box"];
	n5[label="return 4",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_2
}
L_1:
return 2
L_2:
foo(3)
goto L_3
L_3:
return 4
`,
	},
	{
		name: "if else return",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(numStmt(2)), Else: newBlock(ret(3))},
			ret(4),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="return 4",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n5[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n3[color="#99999955"];
	n5->n2[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="foo(2)",shape="box"];
	n4[label="return 4",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n2->n5;
	n3->n4;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_3
}
L_1:
foo(2)
goto L_2
L_2:
return 4
L_3:
return 3
`,
	},
	{
		name: "if all return",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(ret(2)), Else: newBlock(ret(3))},
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="return 3",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n4->n2[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_2
}
L_1:
return 2
L_2:
return 3
`,
	},
	{
		name: "if all return unreachable",
		tree: []CStmt{
			&CIfStmt{Cond: numCond(1), Then: newBlock(ret(2)), Else: newBlock(ret(3))},
			ret(4),
		},
		exp: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="return 3",shape="box"];
	n1->n2;
	n2->n3[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n4->n2[color="#99999955"];
`,
		dom: `
	n2[label="if 1",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
`,
		flat: `
if 1 {
goto L_1
} else {
goto L_2
}
L_1:
return 2
L_2:
return 3
`,
	},
	{
		name: "switch",
		tree: []CStmt{
			&CSwitchStmt{
				Cond: cIntLit(1),
				Cases: []*CCaseStmt{
					{Expr: cIntLit(1), Stmts: []CStmt{numStmt(1)}},
					{Expr: cIntLit(2), Stmts: []CStmt{numStmt(2)}},
					{Expr: cIntLit(3), Stmts: []CStmt{numStmt(3), &CBreakStmt{}}},
					{Stmts: []CStmt{numStmt(4)}},
				},
			},
			ret(1),
		},
		exp: `
	n2[label="switch 1",shape="trapezium"];
	n3[label="foo(1)",shape="box"];
	n4[label="foo(2)",shape="box"];
	n5[label="foo(3)",shape="box"];
	n6[label="return 1",shape="box"];
	n7[label="foo(4)",shape="box"];
	n1->n2;
	n2->n3[label="1"];
	n2->n4[label="2"];
	n2->n5[label="3"];
	n2->n7[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n2[color="#99999955"];
	n4->n3[color="#99999955"];
	n4->n5;
	n5->n2[color="#99999955"];
	n5->n4[color="#99999955"];
	n5->n6;
	n6->n7[color="#99999955"];
	n6->n5[color="#99999955"];
	n7->n2[color="#99999955"];
	n7->n6;
`,
		dom: `
	n2[label="switch 1",shape="trapezium"];
	n3[label="foo(1)",shape="box"];
	n4[label="foo(2)",shape="box"];
	n5[label="foo(3)",shape="box"];
	n6[label="return 1",shape="box"];
	n7[label="foo(4)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n2->n5;
	n2->n6;
	n2->n7;
`,
		flat: `
switch 1 {
case 1:
goto L_1
case 2:
goto L_2
case 3:
goto L_3
default:
goto L_5
}
L_1:
foo(1)
goto L_2
L_2:
foo(2)
goto L_3
L_3:
foo(3)
goto L_4
L_4:
return 1
L_5:
foo(4)
goto L_4
`,
	},
	{
		name: "switch no default",
		tree: []CStmt{
			&CSwitchStmt{
				Cond: cIntLit(1),
				Cases: []*CCaseStmt{
					{Expr: cIntLit(1), Stmts: []CStmt{numStmt(1)}},
					{Expr: cIntLit(2), Stmts: []CStmt{&CBreakStmt{}}},
					{Expr: cIntLit(3), Stmts: []CStmt{numStmt(3), &CBreakStmt{}}},
				},
			},
			ret(1),
		},
		exp: `
	n2[label="switch 1",shape="trapezium"];
	n3[label="foo(1)",shape="box"];
	n4[label="return 1",shape="box"];
	n5[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3[label="1"];
	n2->n4[label="2"];
	n2->n5[label="3"];
	n2->n4[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n4;
	n4->n5[color="#99999955"];
	n4->n2[color="#99999955"];
	n4->n3[color="#99999955"];
	n5->n2[color="#99999955"];
	n5->n4;
`,
		dom: `
	n2[label="switch 1",shape="trapezium"];
	n3[label="foo(1)",shape="box"];
	n4[label="return 1",shape="box"];
	n5[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n2->n5;
`,
		flat: `
switch 1 {
case 1:
goto L_1
case 2:
goto L_2
case 3:
goto L_3
default:
goto L_2
}
L_1:
foo(1)
goto L_2
L_2:
return 1
L_3:
foo(3)
goto L_2
`,
	},
	{
		name: "goto",
		tree: []CStmt{
			numStmt(1),
			&CIfStmt{
				Cond: numCond(1),
				Then: newBlock(
					numStmt(2),
					&CGotoStmt{Label: "L1"},
				),
			},
			&CLabelStmt{Label: "L1"},
			ret(3),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="if 1",shape="hexagon"];
	n4[label="foo(2)",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n3->n2[color="#99999955"];
	n3->n4[color="#00aa00"];
	n3->n5[color="#aa0000"];
	n4->n3[color="#99999955"];
	n4->n5;
	n5->n4[color="#99999955"];
	n5->n3[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="if 1",shape="hexagon"];
	n4[label="foo(2)",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4;
	n3->n5;
`,
		flat: `
foo(1)
goto L_1
L_1:
if 1 {
goto L_2
} else {
goto L_3
}
L_2:
foo(2)
goto L_3
L_3:
return 3
`,
	},
	{
		name: "goto loop",
		tree: []CStmt{
			numStmt(0),
			&CLabelStmt{Label: "L1"},
			&CIfStmt{
				Cond: numCond(1),
				Then: newBlock(
					numStmt(1),
					&CGotoStmt{Label: "L1"},
				),
			},
			ret(3),
		},
		exp: `
	n2[label="foo(0)",shape="box"];
	n3[label="if 1",shape="hexagon"];
	n4[label="foo(1)",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4[color="#99999955"];
	n3->n2[color="#99999955"];
	n3->n4[color="#00aa00"];
	n3->n5[color="#aa0000"];
	n4->n3[color="#99999955"];
	n4->n3;
	n5->n3[color="#99999955"];
`,
		dom: `
	n2[label="foo(0)",shape="box"];
	n3[label="if 1",shape="hexagon"];
	n4[label="foo(1)",shape="box"];
	n5[label="return 3",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4;
	n3->n5;
`,
		flat: `
foo(0)
goto L_1
L_1:
if 1 {
goto L_2
} else {
goto L_3
}
L_2:
foo(1)
goto L_1
L_3:
return 3
`,
	},
	{
		name: "for infinite empty",
		tree: []CStmt{
			&CForStmt{},
			ret(2),
		},
		exp: `
	n2[label="",shape="box"];
	n1->n2;
	n2->n2[color="#99999955"];
	n2->n2;
`,
		dom: `
	n2[label="",shape="box"];
	n1->n2;
`,
		flat: `
L_1:
goto L_1
`,
	},
	{
		name: "for infinite",
		tree: []CStmt{
			&CForStmt{Body: *newBlock(numStmt(1))},
			ret(2),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n1->n2;
	n2->n2[color="#99999955"];
	n2->n2;
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n1->n2;
`,
		flat: `
L_1:
foo(1)
goto L_1
`,
	},
	{
		name: "for infinite break",
		tree: []CStmt{
			&CForStmt{Body: *newBlock(
				numStmt(1),
				&CBreakStmt{},
			)},
			ret(2),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="return 2",shape="box"];
	n1->n2;
	n2->n3;
	n3->n2[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="return 2",shape="box"];
	n1->n2;
	n2->n3;
`,
		flat: `
foo(1)
goto L_1
L_1:
return 2
`,
	},
	{
		name: "for infinite continue",
		tree: []CStmt{
			&CForStmt{Body: *newBlock(
				numStmt(1),
				&CContinueStmt{},
			)},
			ret(2),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n1->n2;
	n2->n2[color="#99999955"];
	n2->n2;
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n1->n2;
`,
		flat: `
L_1:
foo(1)
goto L_1
`,
	},
	{
		name: "for cond",
		tree: []CStmt{
			&CForStmt{Cond: cIntLit(1), Body: *newBlock(
				numStmt(1),
			)},
			ret(2),
		},
		exp: `
	n2[label="if false",shape="hexagon"];
	n3[label="foo(1)",shape="box"];
	n4[label="return 2",shape="box"];
	n1->n2;
	n2->n3[color="#99999955"];
	n2->n4[color="#00aa00"];
	n2->n3[color="#aa0000"];
	n3->n2[color="#99999955"];
	n3->n2;
	n4->n2[color="#99999955"];
`,
		dom: `
	n2[label="if false",shape="hexagon"];
	n3[label="return 2",shape="box"];
	n4[label="foo(1)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
`,
		flat: `
L_1:
if false {
goto L_2
} else {
goto L_3
}
L_2:
return 2
L_3:
foo(1)
goto L_1
`,
	},
	{
		name: "for cond break",
		tree: []CStmt{
			&CForStmt{
				Cond: cIntLit(1),
				Body: *newBlock(
					&CIfStmt{
						Cond: numCond(2),
						Then: newBlock(
							numStmt(3),
							&CBreakStmt{},
						),
					},
					numStmt(4),
				),
			},
			ret(5),
		},
		exp: `
	n2[label="if false",shape="hexagon"];
	n3[label="foo(4)",shape="box"];
	n4[label="if 2",shape="hexagon"];
	n5[label="foo(3)",shape="box"];
	n6[label="return 5",shape="box"];
	n1->n2;
	n2->n3[color="#99999955"];
	n2->n6[color="#00aa00"];
	n2->n4[color="#aa0000"];
	n3->n4[color="#99999955"];
	n3->n2;
	n4->n2[color="#99999955"];
	n4->n5[color="#00aa00"];
	n4->n3[color="#aa0000"];
	n5->n4[color="#99999955"];
	n5->n6;
	n6->n5[color="#99999955"];
	n6->n2[color="#99999955"];
`,
		dom: `
	n2[label="if false",shape="hexagon"];
	n3[label="return 5",shape="box"];
	n4[label="if 2",shape="hexagon"];
	n5[label="foo(3)",shape="box"];
	n6[label="foo(4)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
	n4->n6;
`,
		flat: `
L_1:
if false {
goto L_2
} else {
goto L_3
}
L_2:
return 5
L_3:
if 2 {
goto L_4
} else {
goto L_5
}
L_4:
foo(3)
goto L_2
L_5:
foo(4)
goto L_1
`,
	},
	{
		name: "for cond nested",
		tree: []CStmt{
			&CForStmt{Cond: cIntLit(2), Body: *newBlock(
				numStmt(2),
				&CForStmt{Cond: cIntLit(3), Body: *newBlock(
					numStmt(3),
				)},
			)},
			ret(1),
		},
		exp: `
	n2[label="if 2 == 0",shape="hexagon"];
	n3[label="if 3 == 0",shape="hexagon"];
	n4[label="foo(3)",shape="box"];
	n5[label="foo(2)",shape="box"];
	n6[label="return 1",shape="box"];
	n1->n2;
	n2->n3[color="#99999955"];
	n2->n6[color="#00aa00"];
	n2->n5[color="#aa0000"];
	n3->n4[color="#99999955"];
	n3->n5[color="#99999955"];
	n3->n2[color="#00aa00"];
	n3->n4[color="#aa0000"];
	n4->n3[color="#99999955"];
	n4->n3;
	n5->n2[color="#99999955"];
	n5->n3;
	n6->n2[color="#99999955"];
`,
		dom: `
	n2[label="if 2 == 0",shape="hexagon"];
	n3[label="return 1",shape="box"];
	n4[label="foo(2)",shape="box"];
	n5[label="if 3 == 0",shape="hexagon"];
	n6[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
	n5->n6;
`,
		flat: `
L_1:
if 2 == 0 {
goto L_2
} else {
goto L_3
}
L_2:
return 1
L_3:
foo(2)
goto L_4
L_4:
if 3 == 0 {
goto L_1
} else {
goto L_5
}
L_5:
foo(3)
goto L_4
`,
	},
	{
		name: "for cond nested 2",
		tree: []CStmt{
			&CForStmt{Cond: cIntLit(2), Body: *newBlock(
				numStmt(2),
				&CForStmt{Cond: cIntLit(3), Body: *newBlock(
					numStmt(3),
				)},
				numStmt(4),
			)},
			ret(1),
		},
		exp: `
	n2[label="if 2 == 0",shape="hexagon"];
	n3[label="foo(4)",shape="box"];
	n4[label="if 3 == 0",shape="hexagon"];
	n5[label="foo(3)",shape="box"];
	n6[label="foo(2)",shape="box"];
	n7[label="return 1",shape="box"];
	n1->n2;
	n2->n3[color="#99999955"];
	n2->n7[color="#00aa00"];
	n2->n6[color="#aa0000"];
	n3->n4[color="#99999955"];
	n3->n2;
	n4->n5[color="#99999955"];
	n4->n6[color="#99999955"];
	n4->n3[color="#00aa00"];
	n4->n5[color="#aa0000"];
	n5->n4[color="#99999955"];
	n5->n4;
	n6->n2[color="#99999955"];
	n6->n4;
	n7->n2[color="#99999955"];
`,
		dom: `
	n2[label="if 2 == 0",shape="hexagon"];
	n3[label="return 1",shape="box"];
	n4[label="foo(2)",shape="box"];
	n5[label="if 3 == 0",shape="hexagon"];
	n6[label="foo(4)",shape="box"];
	n7[label="foo(3)",shape="box"];
	n1->n2;
	n2->n3;
	n2->n4;
	n4->n5;
	n5->n6;
	n5->n7;
`,
		flat: `
L_1:
if 2 == 0 {
goto L_2
} else {
goto L_3
}
L_2:
return 1
L_3:
foo(2)
goto L_4
L_4:
if 3 == 0 {
goto L_5
} else {
goto L_6
}
L_5:
foo(4)
goto L_1
L_6:
foo(3)
goto L_4
`,
	},
	{
		name: "for full",
		tree: []CStmt{
			&CForStmt{
				Init: numStmt(1), Cond: cIntLit(1), Iter: numStmt(2),
				Body: *newBlock(
					numStmt(3),
				),
			},
			ret(4),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="if false",shape="hexagon"];
	n4[label="foo(2)",shape="box"];
	n5[label="foo(3)",shape="box"];
	n6[label="return 4",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4[color="#99999955"];
	n3->n2[color="#99999955"];
	n3->n6[color="#00aa00"];
	n3->n5[color="#aa0000"];
	n4->n5[color="#99999955"];
	n4->n3;
	n5->n3[color="#99999955"];
	n5->n4;
	n6->n3[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="if false",shape="hexagon"];
	n4[label="return 4",shape="box"];
	n5[label="foo(3)",shape="box"];
	n6[label="foo(2)",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4;
	n3->n5;
	n5->n6;
`,
		flat: `
foo(1)
goto L_1
L_1:
if false {
goto L_2
} else {
goto L_3
}
L_2:
return 4
L_3:
foo(3)
goto L_4
L_4:
foo(2)
goto L_1
`,
	},
	{
		name: "for full break",
		tree: []CStmt{
			&CForStmt{
				Init: numStmt(1), Cond: cIntLit(1), Iter: numStmt(2),
				Body: *newBlock(
					&CIfStmt{
						Cond: numCond(2),
						Then: newBlock(
							numStmt(3),
							&CBreakStmt{},
						),
					},
					numStmt(4),
				),
			},
			ret(5),
		},
		exp: `
	n2[label="foo(1)",shape="box"];
	n3[label="if false",shape="hexagon"];
	n4[label="foo(2)",shape="box"];
	n5[label="foo(4)",shape="box"];
	n6[label="if 2",shape="hexagon"];
	n7[label="foo(3)",shape="box"];
	n8[label="return 5",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4[color="#99999955"];
	n3->n2[color="#99999955"];
	n3->n8[color="#00aa00"];
	n3->n6[color="#aa0000"];
	n4->n5[color="#99999955"];
	n4->n3;
	n5->n6[color="#99999955"];
	n5->n4;
	n6->n3[color="#99999955"];
	n6->n7[color="#00aa00"];
	n6->n5[color="#aa0000"];
	n7->n6[color="#99999955"];
	n7->n8;
	n8->n7[color="#99999955"];
	n8->n3[color="#99999955"];
`,
		dom: `
	n2[label="foo(1)",shape="box"];
	n3[label="if false",shape="hexagon"];
	n4[label="return 5",shape="box"];
	n5[label="if 2",shape="hexagon"];
	n6[label="foo(3)",shape="box"];
	n7[label="foo(4)",shape="box"];
	n8[label="foo(2)",shape="box"];
	n1->n2;
	n2->n3;
	n3->n4;
	n3->n5;
	n5->n6;
	n5->n7;
	n7->n8;
`,
		flat: `
foo(1)
goto L_1
L_1:
if false {
goto L_2
} else {
goto L_3
}
L_2:
return 5
L_3:
if 2 {
goto L_4
} else {
goto L_5
}
L_4:
foo(3)
goto L_2
L_5:
foo(4)
goto L_6
L_6:
foo(2)
goto L_1
`,
	},
}

func writeDotFile(name string, data []byte) {
	name = strings.ReplaceAll(name, " ", "_")
	fname := name + ".dot"
	_ = ioutil.WriteFile(fname, data, 0644)
	sdata, _ := exec.Command("dot", "-Tsvg", fname).Output()
	_ = ioutil.WriteFile(name+".svg", sdata, 0644)
	_ = os.Remove(fname)
}

func cleanDot(s string) string {
	s = strings.TrimPrefix(s, `digraph  {`)
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, `n1[label="begin"];`)
	s = strings.TrimSuffix(s, `}`)
	s = strings.TrimSpace(s)
	return s
}

func TestControlFlow(t *testing.T) {
	const dir = "testout"
	_ = os.MkdirAll(dir, 0755)
	for _, c := range casesControlFlow {
		t.Run(c.name, func(t *testing.T) {
			tr := newTranslator(libs.NewEnv(types.Config32()), Config{})
			var fixer Visitor
			fixer = func(n Node) {
				switch n := n.(type) {
				case nil:
					return
				case *CVarDecl:
					n.g = tr
				case *CVarSpec:
					n.g = tr
				}
				n.Visit(fixer)
			}

			for _, st := range c.tree {
				st.Visit(fixer)
			}

			cf := tr.NewControlFlow(c.tree)

			got := cf.dumpDot()
			writeDotFile(filepath.Join(dir, c.name), []byte(got))
			got = cleanDot(got)
			require.Equal(t, strings.TrimSpace(c.exp), got)

			got = cf.dumpDomDot()
			writeDotFile(filepath.Join(dir, c.name+"_dom"), []byte(got))
			got = cleanDot(got)
			require.Equal(t, strings.TrimSpace(c.dom), got)

			stmts := cf.Flatten()
			got = printStmts(stmts)
			got = strings.ReplaceAll(got, "\t", "")
			require.Equal(t, strings.TrimSpace(c.flat), got)
		})
	}
}
