package cxgo

import (
	"fmt"
)

type Block interface {
	AddPrevBlock(b2 Block)
	PrevBlocks() []Block
	NextBlocks() []Block
	ReplaceNext(old, rep Block)
}

func replacePrev(v Block, p, to Block) {
	if v == nil {
		return
	}
	prev := v.PrevBlocks()
	found := false
	for i, p2 := range prev {
		if p == p2 {
			prev[i] = to
			found = true
		}
	}
	if !found {
		panic("not found")
	}
}

func replaceBlock(a, b Block) {
	if a == b {
		return
	}
	for _, p := range a.PrevBlocks() {
		if p == a {
			b.AddPrevBlock(b)
		} else {
			b.AddPrevBlock(p)
			p.ReplaceNext(a, b)
		}
	}
}

func (g *translator) NewControlFlow(stmts []CStmt) *ControlFlow {
	cf := &ControlFlow{
		g:      g,
		labels: make(map[string]Block),
	}
	cf.Start, _ = cf.process(stmts, nil)
	cf.buildDoms()
	return cf
}

type ControlFlow struct {
	g      *translator
	Start  Block
	labels map[string]Block
	breaks []Block
	conts  []Block
	doms   map[Block]BlockSet
}

func (cf *ControlFlow) eachBlock(fnc func(b Block)) {
	cf.eachBlockSub(cf.Start, fnc, make(BlockSet))
}

func (cf *ControlFlow) allBlocks() BlockSet {
	m := make(BlockSet)
	cf.eachBlockSub(cf.Start, nil, m)
	return m
}

func (cf *ControlFlow) eachBlockSub(b Block, fnc func(b Block), seen BlockSet) {
	if b == nil {
		return
	}
	if _, ok := seen[b]; ok {
		return
	}
	seen[b] = struct{}{}
	if fnc != nil {
		fnc(b)
	}
	for _, b2 := range b.NextBlocks() {
		cf.eachBlockSub(b2, fnc, seen)
	}
}

func (cf *ControlFlow) process(stmts []CStmt, after Block) (Block, bool) {
	if len(stmts) == 0 {
		return after, false
	}
	switch s := stmts[0].(type) {
	case *BlockStmt:
		if len(stmts) == 1 {
			return cf.process(s.Stmts, after)
		}
		next, _ := cf.process(stmts[1:], after)
		return cf.process(s.Stmts, next)
	case *CReturnStmt:
		_, _ = cf.process(stmts[1:], after)
		return &ReturnBlock{
			CReturnStmt: s,
		}, false
	case *CLabelStmt:
		name := s.Label
		next, _ := cf.process(stmts[1:], after)
		if next == nil {
			panic("must not be nil")
		}
		if tmp := cf.labels[name]; tmp != nil {
			replaceBlock(tmp, next)
		}
		cf.labels[name] = next
		return next, false
	case *CGotoStmt:
		_, _ = cf.process(stmts[1:], after)
		name := s.Label
		b, ok := cf.labels[name]
		if !ok {
			b = &CodeBlock{}
			cf.labels[name] = b
		}
		return b, false
	case *CContinueStmt:
		_, _ = cf.process(stmts[1:], after)
		return cf.conts[len(cf.conts)-1], false
	case *CBreakStmt:
		_, _ = cf.process(stmts[1:], after)
		return cf.breaks[len(cf.breaks)-1], false
	//case *CFallthroughStmt:
	//	_, _ = cf.process(stmts[1:], after)
	//	return cf.falls[len(cf.falls)-1], false
	case *CIfStmt:
		next, _ := cf.process(stmts[1:], after)

		then, _ := cf.process(s.Then.Stmts, next)
		els := next
		if s.Else != nil {
			els, _ = cf.process([]CStmt{s.Else}, next)
		}
		return NewCondBlock(s.Cond, then, els), false
	case *CForStmt:
		// continue target, temporary
		cont := &CodeBlock{}
		// break target
		brk, _ := cf.process(stmts[1:], after)

		// put break/continue blocks to the stack
		bi := len(cf.breaks)
		ci := len(cf.conts)
		cf.breaks = append(cf.breaks, brk)
		cf.conts = append(cf.conts, cont)
		// process the body assuming those break/continue blocks
		body, _ := cf.process(s.Body.Stmts, cont)
		// restore the break/continue stack
		cf.breaks = cf.breaks[:bi]
		cf.conts = cf.conts[:ci]
		// if loop is empty - set an empty body
		if body == nil {
			body = &CodeBlock{Next: cont}
			cont.AddPrevBlock(body)
		}

		// loop with no condition should return to the beginning of the body
		loop := body
		if s.Cond != nil {
			// loop with the condition returns to it instead of the body
			// we invert the expression and break in the positive if branch
			loop = NewCondBlock(cf.g.cNot(s.Cond), brk, body)
		}

		// loop without an iter should continue to the beginning of the loop body
		cont.Next = loop
		if s.Iter != nil {
			// inline iter into the temporary continue block
			// do not remove the temporary block - it's now permanent
			cont.Stmts = append(cont.Stmts, s.Iter)
			loop.AddPrevBlock(cont)
		} else {
			if loop == cont {
				// fix for an empty loop body
				loop.AddPrevBlock(loop)
			}
			replaceBlock(cont, loop)
		}

		if s.Init == nil {
			// no init - start from the loop cond/body
			return loop, false
		}
		// start from the init, then continue to the cond/body
		in := &CodeBlock{Stmts: []CStmt{s.Init}}
		in.Next = loop
		loop.AddPrevBlock(in)
		return in, false
	case *CCaseStmt:
		panic("must not contain cases")
	case *CSwitchStmt:
		// TODO: add fallthrough statement and rewrite this
		next, _ := cf.process(stmts[1:], after)

		b := &SwitchBlock{
			Expr:   s.Cond,
			Cases:  make([]Expr, len(s.Cases)),
			Blocks: make([]Block, len(s.Cases)),
		}
		bi := len(cf.breaks)
		cf.breaks = append(cf.breaks, next)
		// process backward because we need to handle falltrough
		hasDef := false
		fall := next
		for i := len(s.Cases) - 1; i >= 0; i-- {
			c := s.Cases[i]
			if c.Expr == nil {
				hasDef = true
			}
			b.Cases[i] = c.Expr
			cb, _ := cf.process(c.Stmts, fall)
			b.Blocks[i] = cb
			cb.AddPrevBlock(b)
			fall = cb
		}
		cf.breaks = cf.breaks[:bi]
		if !hasDef {
			b.Cases = append(b.Cases, nil)
			b.Blocks = append(b.Blocks, next)
			next.AddPrevBlock(b)
		}
		return b, false
	}
	if after == nil {
		after = &ReturnBlock{CReturnStmt: &CReturnStmt{}}
	}
	b := &CodeBlock{}
	b.Stmts = append(b.Stmts, stmts[0])
	b2, merge := cf.process(stmts[1:], after)
	if c2, ok := b2.(*CodeBlock); ok && merge {
		b.Stmts = append(b.Stmts, c2.Stmts...)
		b.Next = c2.Next
		if b.Next != nil {
			replacePrev(b.Next, c2, b)
		} else {
			b.Next = after
			after.AddPrevBlock(b)
		}
		return b, true
	}
	if b2 == nil {
		b2 = after
	}
	b.Next = b2
	b.Next.AddPrevBlock(b)
	return b, true
}

type BaseBlock struct {
	prev []Block
}

func (b *BaseBlock) AddPrevBlock(b2 Block) {
	for _, b := range b.prev {
		if b == b2 {
			return
		}
	}
	b.prev = append(b.prev, b2)
}
func (b *BaseBlock) PrevBlocks() []Block {
	return b.prev
}

type CodeBlock struct {
	BaseBlock
	Stmts []CStmt
	Next  Block
}

func (b *CodeBlock) NextBlocks() []Block {
	if b.Next == nil {
		// FIXME
		//panic("no following block")
		return nil
	}
	return []Block{b.Next}
}
func (b *CodeBlock) ReplaceNext(old, rep Block) {
	if b.Next == old {
		b.Next = rep
	}
}

func NewCondBlock(expr Expr, then, els Block) *CondBlock {
	if then == nil || els == nil {
		panic("both branches must be set")
	}
	b := &CondBlock{
		Expr: expr,
		Then: then,
		Else: els,
	}
	b.Then.AddPrevBlock(b)
	b.Else.AddPrevBlock(b)
	return b
}

type CondBlock struct {
	BaseBlock
	Expr Expr
	Then Block
	Else Block
}

func (b *CondBlock) NextBlocks() []Block {
	if b.Then == nil || b.Else == nil {
		panic("no following block")
	}
	return []Block{b.Then, b.Else}
}
func (b *CondBlock) ReplaceNext(old, rep Block) {
	if b.Then == old {
		b.Then = rep
	}
	if b.Else == old {
		b.Else = rep
	}
}

type SwitchBlock struct {
	BaseBlock
	Expr   Expr
	Cases  []Expr
	Blocks []Block
}

func (b *SwitchBlock) NextBlocks() []Block {
	return append([]Block{}, b.Blocks...)
}
func (b *SwitchBlock) ReplaceNext(old, rep Block) {
	for i, p := range b.Blocks {
		if p == old {
			b.Blocks[i] = rep
		}
	}
}

type ReturnBlock struct {
	BaseBlock
	*CReturnStmt
}

func (b *ReturnBlock) NextBlocks() []Block {
	return nil
}
func (b *ReturnBlock) ReplaceNext(old, rep Block) {}

func (cf *ControlFlow) Dom(a, b Block) bool {
	if a == b {
		return true
	}
	_, ok := cf.doms[b][a]
	return ok
}

func (cf *ControlFlow) SDom(a, b Block) bool {
	if a == b {
		// strict dominance - nodes shouldn't be the same
		return false
	}
	return cf.Dom(a, b)
}

func (cf *ControlFlow) IDom(a Block) Block {
	var doms []Block
	for d := range cf.doms[a] {
		if d != a {
			doms = append(doms, d)
		}
	}
	for len(doms) > 1 {
	loop:
		for i := 0; i < len(doms); i++ {
			d1 := doms[i]
			for j := 0; j < len(doms); j++ {
				if i == j {
					continue
				}
				d2 := doms[j]
				if _, ok := cf.doms[d2][d1]; ok {
					doms = append(doms[:i], doms[i+1:]...)
					i--
					continue loop
				}
			}
		}
	}
	if len(doms) == 1 {
		return doms[0]
	}
	return nil
}

type BlockSet map[Block]struct{}

func (b BlockSet) Clone() BlockSet {
	b2 := make(BlockSet, len(b))
	for k := range b {
		b2[k] = struct{}{}
	}
	return b2
}

func (b BlockSet) Union(b2 BlockSet) BlockSet {
	if len(b) < len(b2) {
		b, b2 = b2, b
	}
	all := true
	for k := range b2 {
		if _, ok := b[k]; !ok {
			all = false
			break
		}
	}
	if all {
		return b
	}
	m := make(BlockSet, len(b))
	for k := range b {
		m[k] = struct{}{}
	}
	for k := range b2 {
		m[k] = struct{}{}
	}
	return m
}

func (b BlockSet) Intersect(b2 BlockSet) BlockSet {
	if len(b) > len(b2) {
		b, b2 = b2, b
	}
	m := make(BlockSet, len(b))
	for k := range b {
		if _, ok := b2[k]; ok {
			m[k] = struct{}{}
		}
	}
	return m
}

func (b BlockSet) Contains(b2 BlockSet) bool {
	if len(b) < len(b2) {
		return false
	}
	for k := range b2 {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}

func (cf *ControlFlow) buildDoms() {
	// TODO: quadratic complexity! use a different algorithm
	blocks := cf.allBlocks()
	cf.doms = make(map[Block]BlockSet, len(blocks))
	for b := range blocks {
		if b == cf.Start {
			cf.doms[b] = BlockSet{b: {}}
		} else {
			cf.doms[b] = blocks
		}
	}

	changes := true
	for changes {
		changes = false
		for b := range blocks {
			d := cf.doms[b]
			if len(d) == 1 {
				continue
			}
			var m BlockSet
			for _, b2 := range b.PrevBlocks() {
				if b == b2 {
					continue
				}
				d2 := cf.doms[b2]
				if m == nil {
					m = d2
					continue
				}
				m = m.Intersect(d2)
				if len(m) == 0 {
					break
				}
			}
			m = m.Union(BlockSet{b: {}})
			if len(d) != len(m) {
				changes = true
				cf.doms[b] = m
			}
			if len(d) < len(m) {
				panic(fmt.Errorf("set size increased: %p: %v vs %v", b, d, m))
			}
		}
	}
}

type varDecls struct {
	Decls []*CVarDecl
}

func (cf *ControlFlow) Flatten() []CStmt {
	labels := make(map[Block]int)
	cf.eachBlock(func(b Block) {
		prev := b.PrevBlocks()
		if len(prev) == 0 {
			return
		}
		labels[b] = len(labels) + 1
	})
	var decls varDecls
	stmts := cf.flatten(cf.Start, &decls, nil, labels, make(map[Block]struct{}))
	var out []CStmt
	if len(decls.Decls) != 0 {
		for _, d := range decls.Decls {
			out = append(out, &CDeclStmt{Decl: d})
		}
	}
	out = append(out, stmts...)
	return out
}

func numLabelName(n int) string {
	return fmt.Sprintf("L_%d", n)
}

func numLabel(n int) *CLabelStmt {
	return &CLabelStmt{Label: numLabelName(n)}
}

func numGoto(n int) *CGotoStmt {
	return &CGotoStmt{Label: numLabelName(n)}
}

func (cf *ControlFlow) slitDecls(decl *varDecls, stmts []CStmt) []CStmt {
	out := make([]CStmt, 0, len(stmts))
	for _, st := range stmts {
		ds, ok := st.(*CDeclStmt)
		if !ok {
			out = append(out, st)
			continue
		}
		d, ok := ds.Decl.(*CVarDecl)
		if !ok {
			out = append(out, st)
			continue
		} else if len(d.Inits) == 0 || d.Const {
			out = append(out, st)
			continue
		} else if d.Names[0].Name == "__func__" {
			out = append(out, st)
			continue
		}
		decl.Decls = append(decl.Decls, &CVarDecl{
			Const:  d.Const,
			Single: d.Single,
			CVarSpec: CVarSpec{
				g:     d.g,
				Type:  d.Type,
				Names: d.Names,
			},
		})
		for j, val := range d.Inits {
			if val != nil {
				out = append(out, d.g.NewCAssignStmt(IdentExpr{d.Names[j]}, "", val, false)...)
			}
		}
	}
	return out
}

func (cf *ControlFlow) flatten(b Block, decl *varDecls, stmts []CStmt, labels map[Block]int, seen map[Block]struct{}) []CStmt {
	if b == nil {
		return stmts
	}
	if _, ok := seen[b]; ok {
		l, ok := labels[b]
		if !ok {
			panic(fmt.Errorf("must have a label: %T, %v", b, len(b.PrevBlocks())))
		}
		stmts = append(stmts, numGoto(l))
		return stmts
	}
	seen[b] = struct{}{}
	if id, ok := labels[b]; ok {
		stmts = append(stmts, numLabel(id))
	}
	switch b := b.(type) {
	case *CodeBlock:
		cur := cf.slitDecls(decl, b.Stmts)
		stmts = append(stmts, cur...)
		if l, ok := labels[b.Next]; ok {
			stmts = append(stmts, numGoto(l))
			if _, ok := seen[b.Next]; ok {
				return stmts
			}
		}
		stmts = cf.flatten(b.Next, decl, stmts, labels, seen)
	case *CondBlock:
		then, ok := labels[b.Then]
		if !ok {
			panic("must have a label")
		}
		els, ok := labels[b.Else]
		if !ok {
			panic("must have a label")
		}
		stmts = append(stmts, &CIfStmt{
			Cond: cf.g.ToBool(b.Expr),
			Then: cf.g.NewCBlock(numGoto(then)),
			Else: cf.g.NewCBlock(numGoto(els)),
		})
		if _, ok := seen[b.Then]; !ok {
			stmts = cf.flatten(b.Then, decl, stmts, labels, seen)
		}
		if _, ok := seen[b.Else]; !ok {
			stmts = cf.flatten(b.Else, decl, stmts, labels, seen)
		}
	case *ReturnBlock:
		stmts = append(stmts, b.CReturnStmt)
	case *SwitchBlock:
		s := &CSwitchStmt{
			Cond: b.Expr,
		}
		for i, e := range b.Cases {
			l, ok := labels[b.Blocks[i]]
			if !ok {
				panic("must have a label")
			}
			s.Cases = append(s.Cases, cf.g.NewCaseStmt(
				e, numGoto(l),
			))
		}
		stmts = append(stmts, s)
		for _, b := range b.Blocks {
			if _, ok := seen[b]; !ok {
				stmts = cf.flatten(b, decl, stmts, labels, seen)
			}
		}
	default:
		panic(b)
	}
	return stmts
}
