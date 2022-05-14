package libcc

import (
	"modernc.org/cc/v4"

	"github.com/gotranspile/cxgo/types"
)

func NewABI(c *types.Env) *cc.ABI {
	// TODO: expose some other constructor?
	intSize := c.IntSize()
	ptrSize := c.PtrSize()
	if intSize != ptrSize {
		panic("TODO: int and pointer size differ")
	}
	arch := "amd64"
	if ptrSize == 4 {
		arch = "386"
	}
	abi, err := cc.NewABI("linux", arch)
	if err != nil {
		panic(err)
	}
	return abi
}
