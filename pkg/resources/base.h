#ifndef BASE_H
#define BASE_H

#include <stddef.h>
#include <stdint.h>

#define OK 0
#define ERR_IO 1
#define ERR_PARSE 2
#define ERR_PERM 3
#define ERR_INVALID 4

#define internal static
#define local_persist static
#define global_variable static

enum
{
  BUFFER_SIZE_SMALL = 128,
  BUFFER_SIZE_DEFAULT = 256,
  BUFFER_SIZE_LARGE = 512,
  PATH_MAX_LEN = 4096
};

typedef uint64_t u64;
typedef uint32_t u32;
typedef uint16_t u16;
typedef uint8_t u8;

typedef int8_t i8;
typedef int16_t i16;
typedef int32_t i32;
typedef int64_t i64;

typedef float f32;
typedef double f64;

typedef i32 b32;
typedef i16 b16;
typedef u8 b8;

#define TRUE 1
#define FALSE 0

u64
parse_u64(char *buf, size_t len);

#endif
