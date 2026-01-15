#ifndef ARENA_H
#define ARENA_H

#include "base.h"

/**
 * Arena Helper macro's
 * */

#define MIN(a, b) (((a) < (b)) ? (a) : (b))
#define MAX(a, b) (((a) > (b)) ? (a) : (b))
#define ALIGN_UP_POW2(n, p) (((u64)(n) + ((u64)(p) - 1)) & (~((u64)(p) - 1)))

#define PUSH_STRUCT(arena, T) (T *)arena_push((arena), sizeof(T), 0)
#define PUSH_STRUCT_NZ(arena, T) (T *)arena_push((arena), sizeof(T), 1)
#define PUSH_ARRAY(arena, T, n) (T *)arena_push((arena), sizeof(T) * (n), 0)
#define PUSH_ARRAY_NZ(arena, T, n) (T *)arena_push((arena), sizeof(T) * (n), 1)




/*
 * Represents a disk partition with major/minor device numbers and block count.
 */

/**
 * replacing malloc/free with arena allocaters
 *
 * */

#define ARENA_BASE_POS (sizeof(mem_arena))
// void * for the size of a pointer on the machine, 64/32bit comp
#define ARENA_ALIGN (sizeof(void *))


static inline u64 KiB(u64 n) { return n << 10; }
static inline u64 MiB(u64 n) { return n << 20; }
static inline u64 GiB(u64 n) { return n << 30; }

typedef struct mem_arena mem_arena;


struct mem_arena 
{
  u64 capacity;
  u64 pos;
} ;

// arena prototypes
mem_arena *
arena_create(u64 capacity);
// make it a void pointer to allow implicit conversion
void
arena_destroy(mem_arena *arena);

void *
arena_push(mem_arena *arena, u64 size, b32 non_zero);

void
arena_pop(mem_arena *arena, u64 size);

void
arena_pop_to(mem_arena *arena, u64 pos);

void
arena_clear(mem_arena *arena);


#define PUSH_STRUCT(arena, T) (T *)arena_push((arena), sizeof(T), 0)
#define PUSH_STRUCT_NZ(arena, T) (T *)arena_push((arena), sizeof(T), 1)
#define PUSH_ARRAY(arena, T, n) (T *)arena_push((arena), sizeof(T) * (n), 0)
#define PUSH_ARRAY_NZ(arena, T, n) (T *)arena_push((arena), sizeof(T) * (n), 1)


#endif
