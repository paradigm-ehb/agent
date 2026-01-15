#include "base.h"
#include "resources.h"
#include "stdio.h"
#include "arena.h"
#include "arena.c"
#include "resources.c"
#include <assert.h>

int
main(void)
{
  mem_arena *arena = arena_create(MiB(1));

  Ram *ram = ram_create(arena);
  Device *device = device_create(arena);

  ram_read(ram);
  process_list_collect(&device->processes, arena);

  for (size_t proc_idx = 0;
    proc_idx < device->processes.count;
    ++proc_idx)
  {
    Process *proc = (Process *)arena_push(arena,
      sizeof(Process),
      1);

    i32 error = process_read(
      device->processes
        .items[proc_idx]
        .pid,
      proc);

    printf(
      "[total] total=%6s  free=%3s\n",
      ram->total,
      ram->free
      );

    if (error != OK)
    {
      assert(0);
    }
  }

  return 0;
}
