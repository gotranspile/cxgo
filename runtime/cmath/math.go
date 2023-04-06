package cmath

import (
	"math"

	"github.com/chewxy/math32"
)

const TowardZero = 1

func FSetRound(r int32) int32 {
	// FIXME
	return 0
}

func Abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func Modf(v float64, iptr *float64) float64 {
	intg, frac := math.Modf(v)
	*iptr = intg
	return frac
}

func Modff(v float32, iptr *float32) float32 {
	intg, frac := math32.Modf(v)
	*iptr = intg
	return frac
}
