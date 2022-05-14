package cxgo

const gccPredefine = `
#define __asm asm
#define __attribute(x)
#define __attribute__(x)
#define __builtin___memcpy_chk(x, y, z, t) __BUILTIN___MEMCPY_CHK()
#define __builtin___memset_chk(x, y, z, ...) __BUILTIN___MEMSET_CHK()
#define __builtin_alloca(x) __BUILTIN_ALLOCA()
#define __builtin_classify_type(x) __BUILTIN_CLASSIFY_TYPE()
#define __builtin_constant_p(exp) __BUILTIN_CONSTANT_P()
#define __builtin_isgreater(x, y) __BUILTIN_ISGREATER()
#define __builtin_isless(x, y) __BUILTIN_ISLESS()
#define __builtin_isunordered(x, y) __BUILTIN_ISUNORDERED()
#define __builtin_longjmp(x, y) __BUILTIN_LONGJMP()
#define __builtin_mempcpy(x, y, z) __BUILTIN_MEMPCPY()
#define __builtin_mul_overflow(a, b, c) __BUILTIN_MUL_OVERFLOW()
#define __builtin_signbit(x) __BUILTIN_SIGNBIT()
#define __complex _Complex
#define __complex__ _Complex
#define __const
#define __extension__
#define __imag__
#define __inline inline
#define __real(x) __REAL()
#define __real__
#define __restrict
#define __sync_val_compare_and_swap(x, y, z, ...) __SYNC_VAL_COMPARE_AND_SWAP()
#define __typeof typeof
#define __volatile volatile
%[1]v __builtin_object_size (void*, int);
%[1]v __builtin_strlen(char*);
%[1]v __builtin_strspn(char*, char*);
_Bool __BUILTIN_MUL_OVERFLOW();
char* __builtin___stpcpy_chk(char*, char*, %[1]v);
char* __builtin_stpcpy(char*, char*);
char* __builtin_strchr(char*, int);
char* __builtin_strdup(char*);
char* __builtin_strncpy(char*, char*, %[1]v);
//double _Complex __builtin_cpow(double _Complex, _Complex double);
double __REAL();
double __builtin_copysign(double, double);
double __builtin_copysignl(long double, long double);
double __builtin_inff();
double __builtin_modf(double, double*);
double __builtin_modfl(long double, long double*);
double __builtin_nanf(char *);
//float _Complex __builtin_conjf(float _Complex);
float __builtin_ceilf(float);
float __builtin_copysignf(float, float);
float __builtin_modff(float, float*);
int __BUILTIN_CLASSIFY_TYPE();
int __BUILTIN_CONSTANT_P();
int __BUILTIN_ISGREATER();
int __BUILTIN_ISLESS();
int __BUILTIN_ISUNORDERED();
int __BUILTIN_SIGNBIT();
int __builtin___snprintf_chk (char*, %[1]v, int, %[1]v, char*, ...);
int __builtin___sprintf_chk (char*, int, %[1]v, char*, ...);
int __builtin___vsnprintf_chk (char*, %[1]v, int, %[1]v, char*, void*);
int __builtin___vsprintf_chk (char*, int, %[1]v, char*, void*);
int __builtin_abs(int);
int __builtin_clrsb(int);
int __builtin_clrsbl(long);
int __builtin_clrsbll(long long);
int __builtin_clz(unsigned int);
int __builtin_clzl(unsigned long);
int __builtin_clzll(unsigned long long);
int __builtin_constant_p (exp);
int __builtin_ctz(unsigned int x);
int __builtin_ctzl(unsigned long);
int __builtin_ctzll(unsigned long long);
int __builtin_ffs(int);
int __builtin_ffsl(long);
int __builtin_ffsll(long long);
int __builtin_isinf(double);
int __builtin_isinff(float);
int __builtin_isinfl(long double);
int __builtin_memcmp(void*, void*, %[1]v);
int __builtin_parity (unsigned);
int __builtin_parityl(unsigned long);
int __builtin_parityll (unsigned long long);
int __builtin_popcount (unsigned int x);
int __builtin_popcountl (unsigned long);
int __builtin_popcountll (unsigned long long);
int __builtin_puts(char*);
int __builtin_setjmp(void*);
int __builtin_strcmp(char*, char*);
int __builtin_strncmp(char*, char*, %[1]v);
long long strlen (char*);
unsigned __builtin_bswap32 (unsigned x);
unsigned long long __builtin_bswap64 (unsigned long long x);
unsigned short __builtin_bswap16 (unsigned short x);
void __BUILTIN_LONGJMP();
void __SYNC_VAL_COMPARE_AND_SWAP();
void __builtin_bcopy(void*, void*, %[1]v);
void __builtin_bzero(void*, %[1]v);
void __builtin_prefetch (void*, ...);
void __builtin_stack_restore(void*);
void __builtin_unwind_init();
void* __BUILTIN_ALLOCA();
void* __BUILTIN_MEMPCPY();
void* __BUILTIN___MEMCPY_CHK();
void* __BUILTIN___MEMSET_CHK();
void* __builtin_alloca(int);
void* __builtin_apply (void (*)(), void*, %[1]v);
void* __builtin_apply_args();
void* __builtin_extract_return_addr(void *);
void* __builtin_frame_address(unsigned int);
void* __builtin_return_address (unsigned int);
void* __builtin_stack_save();
`
