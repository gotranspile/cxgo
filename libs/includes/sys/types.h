#include <stddef.h>
#include <time.h>

#define off_t _cxgo_int64
#define ssize_t _cxgo_int64
#define off_t _cxgo_uint64
#define pid_t _cxgo_uint64
#define gid_t _cxgo_uint32
#define uid_t _cxgo_uint32
#define ino_t _cxgo_uint64

#define u_short unsigned short
#define u_long unsigned long


// TODO: should be in fcntl.h
const _cxgo_int32 O_RDONLY = 1;
const _cxgo_int32 O_WRONLY = 2;
const _cxgo_int32 O_RDWR = 3;
const _cxgo_int32 O_CREAT = 4;
const _cxgo_int32 O_EXCL = 5;
const _cxgo_int32 O_TRUNC = 6;
