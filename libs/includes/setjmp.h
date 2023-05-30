typedef struct jmp_buf {
	_cxgo_go_int (*SetJump) ();
	void (*LongJump) (_cxgo_go_int);
} jmp_buf;

#define setjmp(b) ((jmp_buf)b).SetJump()
#define longjmp(b, v) ((jmp_buf)b).LongJump(v)
