package libc

import (
	"math"
	"math/rand"
)

const RandMax = math.MaxInt32-1

func Rand() int32 {
	return int32(rand.Intn(RandMax+1))
}

func SeedRand(seed uint32) {
	rand.Seed(int64(seed))
}
