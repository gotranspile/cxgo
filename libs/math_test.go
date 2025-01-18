package libs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotranspile/cxgo/types"
)

var expMathH = `#ifndef _cxgo_MATH_H
#define _cxgo_MATH_H
#include <cxgo_builtin.h>

const double M_PI_val = 3.1415;
#define M_PI M_PI_val
double sin(double);
float sinf(float);
double asin(double);
float asinf(float);
double sinh(double);
float sinhf(float);
double asinh(double);
float asinhf(float);
double cos(double);
float cosf(float);
double acos(double);
float acosf(float);
double cosh(double);
float coshf(float);
double acosh(double);
float acoshf(float);
double tan(double);
float tanf(float);
double atan(double);
float atanf(float);
double tanh(double);
float tanhf(float);
double atanh(double);
float atanhf(float);
double round(double);
#define roundf(x) round(x)
double ceil(double);
float ceilf(float);
double floor(double);
float floorf(float);
double fabs(double);
float fabsf(float);
double pow(double, double);
float powf(float, float);
double sqrt(double);
float sqrtf(float);
double exp(double);
float expf(float);
double exp2(double);
float exp2f(float);
double log(double);
float logf(float);
double log10(double);
float log10f(float);
double log2(double);
float log2f(float);
double atan2(double y, double x);
float atan2f(float y, float x);
double modf(double x, double *iptr);
float modff(float value, float *iptr);
double ldexp(double x, _cxgo_go_int exp);
double fmod(double x, double exp);
int isnan(double x);
double frexp(double x, int* exp);
double hypot(double x, double y);
float hypotf(float x, float y);
double fmax(double x, double y);
float fmaxf(float x, float y);
double fmin(double x, double y);
float fminf(float x, float y);
int isfinite(double x);


#endif // _cxgo_MATH_H
`

func TestMathH(t *testing.T) {
	c := NewEnv(types.Config32())
	l, ok := c.GetLibrary(mathH)
	require.True(t, ok)
	require.Equal(t, expMathH, l.Header, "\n%s", l.Header)
}
