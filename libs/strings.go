package libs

func init() {
	RegisterLibrary("strings.h", func(c *Env) *Library {
		// TODO: support strings.h
		return &Library{}
	})
}
