#include <assert.h>
#include <dirent.h>
#include <fcntl.h>
#include <stdint.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <unistd.h>

#include "arena.h"

mem_arena *
arena_create(u64 capacity)
{
  mem_arena *arena = mmap(0, capacity, PROT_READ | PROT_WRITE | PROT_EXEC,
    MAP_SHARED | MAP_ANONYMOUS, -1, 0);
  if (arena == MAP_FAILED)
  {
    assert(0);
  }

  arena->capacity = capacity;
  arena->pos = ARENA_BASE_POS;

  return arena;
}

// make it a void pointer to allow implicit conversion
void
arena_destroy(mem_arena *arena)
{
  munmap(arena, arena->capacity);
}

void *
arena_push(mem_arena *arena, u64 size, b32 non_zero)
{
  u64 pos_aligned = ALIGN_UP_POW2(arena->pos, ARENA_ALIGN);
  u64 new_pos = pos_aligned + size;

  if (new_pos > arena->capacity)
  {
    assert(0);
    return NULL;
  }

  arena->pos = new_pos;
  // cast to u8 to be able to do pointer arithemtic
  u8 *out = (u8 *)arena + pos_aligned;

  if (!non_zero)
  {
    memset(out, 0, size);
  }
  return out;
}
void
arena_pop(mem_arena *arena, u64 size)
{
  size = MIN(size, arena->pos - ARENA_BASE_POS);
  arena->pos -= size;
}

void
arena_pop_to(mem_arena *arena, u64 pos)
{
  u64 size = pos < arena->pos ? arena->pos - pos : 0;
  arena_pop(arena, size);
}

void
arena_clear(mem_arena *arena)
{
  arena_pop_to(arena, ARENA_BASE_POS);
}
