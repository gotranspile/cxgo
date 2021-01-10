package libcc

import (
	"encoding/binary"

	"github.com/dennwc/cxgo/types"
	"modernc.org/cc/v3"
)

func NewABI(c *types.Env) cc.ABI {
	intSize := c.IntSize()
	ptrSize := c.PtrSize()
	return cc.ABI{
		ByteOrder: binary.LittleEndian,
		Types: map[cc.Kind]cc.ABIType{
			cc.Bool:      {1, 1, 1},
			cc.Char:      {1, 1, 1},
			cc.Int:       {uintptr(intSize), intSize, intSize},
			cc.Long:      {uintptr(ptrSize), ptrSize, ptrSize},
			cc.LongLong:  {8, 8, intSize},
			cc.SChar:     {1, 1, 1},
			cc.Short:     {2, 2, 2},
			cc.UChar:     {1, 1, 1},
			cc.UInt:      {uintptr(intSize), intSize, intSize},
			cc.ULong:     {uintptr(ptrSize), ptrSize, ptrSize},
			cc.ULongLong: {8, 8, intSize},
			cc.UShort:    {2, 2, 2},

			cc.Int8:   {1, 1, 1},
			cc.UInt8:  {1, 1, 1},
			cc.Int16:  {2, 2, 2},
			cc.UInt16: {2, 2, 2},
			cc.Int32:  {4, 4, 4},
			cc.UInt32: {4, 4, 4},
			cc.Int64:  {8, intSize, intSize},
			cc.UInt64: {8, intSize, intSize},

			cc.Float:      {4, 4, intSize},
			cc.Double:     {8, 8, intSize},
			cc.LongDouble: {8, 8, intSize},

			cc.Void: {1, 1, 1},
			cc.Ptr:  {uintptr(ptrSize), ptrSize, ptrSize},

			cc.Enum: {uintptr(intSize), intSize, intSize},
		},
	}
}
