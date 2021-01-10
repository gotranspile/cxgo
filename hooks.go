package cxgo

type ASTHookCFunc func(c Config, fname string, decls []CDecl) error

var (
	astHooksC []ASTHookCFunc
	finals    []func() error
)

func RegisterASTHookC(fnc ASTHookCFunc) {
	astHooksC = append(astHooksC, fnc)
}

func RegisterFinal(fnc func() error) {
	finals = append(finals, fnc)
}

func CallFinals() error {
	for _, f := range finals {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
