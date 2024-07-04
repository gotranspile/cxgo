package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSameInt(t *testing.T) {
	cases := []struct {
		name string
		x    IntType
		y    IntType
		exp  bool
	}{
		{"same i8", IntT(1), IntT(1), true},
		{"same i64", IntT(8), IntT(8), true},
		{"same u8", UintT(1), UintT(1), true},
		{"same u64", UintT(8), UintT(8), true},
		{"i8 vs u8", IntT(1), UintT(1), false},
		{"i8 vs i16", IntT(1), IntT(2), false},
		{"u8 vs u16", UintT(1), UintT(2), false},
		{"untyped u8 vs u8", AsUntypedIntT(UintT(1)), UintT(1), true},
		{"untyped u64 vs u64", AsUntypedIntT(UintT(8)), UintT(8), true},
		{"untyped i8 vs i8", AsUntypedIntT(IntT(1)), IntT(1), true},
		{"untyped i64 vs i64", AsUntypedIntT(IntT(8)), IntT(8), true},
		{"untyped u8 vs i8", AsUntypedIntT(UintT(1)), IntT(1), false},
		{"untyped i8 vs u8", AsUntypedIntT(IntT(1)), UintT(1), false},
		{"untyped i8 vs i16", AsUntypedIntT(IntT(1)), IntT(2), true},
		{"untyped u8 vs u16", AsUntypedIntT(UintT(1)), UintT(2), true},
		{"untyped i16 vs i8", AsUntypedIntT(IntT(2)), IntT(1), false},
		{"untyped u16 vs u8", AsUntypedIntT(UintT(2)), UintT(1), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.exp, sameInt(c.x, c.y))
			require.Equal(t, c.exp, sameInt(c.y, c.x))
		})
	}
}
