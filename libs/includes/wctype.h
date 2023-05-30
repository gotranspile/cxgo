#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#define wint_t _cxgo_go_rune
const wint_t WEOF = -1;

_Bool      iswalnum(wint_t);
_Bool      iswalpha(wint_t);
_Bool      iswblank(wint_t);
_Bool      iswcntrl(wint_t);
//_Bool      iswctype(wint_t, wctype_t);
_Bool      iswdigit(wint_t);
_Bool      iswgraph(wint_t);
_Bool      iswlower(wint_t);
_Bool      iswprint(wint_t);
_Bool      iswpunct(wint_t);
_Bool      iswspace(wint_t);
_Bool      iswupper(wint_t);
_Bool      iswxdigit(wint_t);
//wint_t    towctrans(wint_t, wctrans_t);
wint_t    towlower(wint_t);
wint_t    towupper(wint_t);
//wctrans_t wctrans(const char *);
//wctype_t  wctype(const char *);
/*
_Bool      iswalnum_l(wint_t, locale_t);
_Bool      iswalpha_l(wint_t, locale_t);
_Bool      iswblank_l(wint_t, locale_t);
_Bool      iswcntrl_l(wint_t, locale_t);
_Bool      iswctype_l(wint_t, wctype_t, locale_t);
_Bool      iswdigit_l(wint_t, locale_t);
_Bool      iswgraph_l(wint_t, locale_t);
_Bool      iswlower_l(wint_t, locale_t);
_Bool      iswprint_l(wint_t, locale_t);
_Bool      iswpunct_l(wint_t, locale_t);
_Bool      iswspace_l(wint_t, locale_t);
_Bool      iswupper_l(wint_t, locale_t);
_Bool      iswxdigit_l(wint_t, locale_t);
wint_t    towctrans_l(wint_t, wctrans_t, locale_t);
wint_t    towlower_l(wint_t, locale_t);
wint_t    towupper_l(wint_t, locale_t);
wctrans_t wctrans_l(const char *, locale_t);
wctype_t  wctype_l(const char *, locale_t);
*/
