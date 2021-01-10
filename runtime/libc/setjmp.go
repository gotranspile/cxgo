package libc

type JumpBuf struct {
	val int // for non-zero size
}

func (b *JumpBuf) SetJump() int {
	// returns 0 on the first call (correct behavior)
	return b.val
}

func (b *JumpBuf) LongJump(val int) {
	b.val = val
	panic("longjmp")
}
