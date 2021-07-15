package libc

import (
	"math"
	"math/rand"
)

const RandMax = math.MaxInt32

func Rand() int32 {
	return int32(rand.Intn(RandMax))
}

func SeedRand(seed uint32) {
	rand.Seed(int64(seed))
}
