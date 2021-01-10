package libc

import "math/rand"

func Rand() int32 {
	return int32(rand.Int())
}

func SeedRand(seed uint32) {
	rand.Seed(int64(seed))
}
