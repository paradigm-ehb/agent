#include "base.h"

/**
 * Helper function to parse strings to int using ascii codes
 * */
u64
parse_u64(char *buf, size_t len)
{
  u64 value = 0;

  for (
    size_t buffer_idx = 0;
    buffer_idx < len;
    ++buffer_idx)
  {
    char c = buf[buffer_idx];
    if (c < '0' || c > '9')
    {
      break;
    }
    u64 digit = (u64)(c - '0');
    value = value * 10 + digit;
  }

  return value;
}
