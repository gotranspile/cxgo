package cxgo

import (
	"log"

	"modernc.org/cc/v4"
)

func (g *translator) runASTPluginsC(cur string, _ *cc.AST, decl []CDecl) []CDecl {
	if g.conf.Hooks {
		for _, f := range astHooksC {
			if err := f(g.conf, cur, decl); err != nil {
				log.Println("error executing hook:", err)
			}
		}
	}
	return decl
}
