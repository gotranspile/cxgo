#include <cxgo_builtin.h>
#include <stddef.h>
#include <stdlib.h>

#define memcpy __builtin_memcpy
#define memmove __builtin_memmove
#define memset __builtin_memset

void    *memccpy(void *restrict, const void *restrict, int, _cxgo_go_int);
void    *memchr(const void *, _cxgo_go_byte, _cxgo_go_int);
_cxgo_go_int      memcmp(const void *, const void *, _cxgo_go_int);
char    *stpcpy(char *restrict, const char *restrict);
char    *stpncpy(char *restrict, const char *restrict, size_t);
char    *strcat(char *restrict, const char *restrict);
char    *strchr(const char *, _cxgo_go_byte);
_cxgo_go_int      strcmp(const char *, const char *);
_cxgo_go_int      strcoll(const char *, const char *);
//int      strcoll_l(const char *, const char *, locale_t);
char    *strcpy(char *restrict, const char *restrict);
_cxgo_go_int   strcspn(const char *, const char *);
#define strdup __builtin_strdup
#define strndup __builtin_strndup
char    *strerror(int);
//char    *strerror_l(int, locale_t);
int      strerror_r(int, char *, size_t);
_cxgo_go_int   strlen(const char *);
char    *strncat(char *restrict, const char *restrict, _cxgo_go_int);
_cxgo_go_int      strncmp(const char *, const char *, _cxgo_go_int);
char    *strncpy(char *restrict, const char *restrict, _cxgo_go_int);
_cxgo_go_int   strnlen(const char *, _cxgo_go_int);
char    *strpbrk(const char *, const char *);
char    *strrchr(const char *, _cxgo_go_byte);
char    *strsignal(int);
_cxgo_go_int   strspn(const char *, const char *);
char    *strstr(const char *, const char *);
char    *strtok(char *restrict, const char *restrict);
char    *strtok_r(char *restrict, const char *restrict, char **restrict);
size_t   strxfrm(char *restrict, const char *restrict, size_t);
//size_t   strxfrm_l(char *restrict, const char *restrict, size_t, locale_t);
_cxgo_go_int strcasecmp(const char *s1, const char *s2);
_cxgo_go_int strncasecmp(const char *s1, const char *s2, _cxgo_go_int n);
