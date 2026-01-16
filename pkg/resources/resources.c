/*
 * name: Abdellah El Morabit
 * organization: Paradigm-Ehb
 * year: 2025-2026
 * description: resources gathering library
 *
 */

#include "base.h"
#define _POSIX_C_SOURCE 200809L

#include "resources.h"
#include "arena.h"

#include <assert.h>
#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>
#include <sys/statfs.h>

/*
 * disk_push_partition - Add a partition to the disk structure
 * @d: Pointer to the Disk structure
 * @p: Partition to add
 *
 * Dynamically grows the partition array if needed. Doubles capacity when
 * full. If realloc fails, the partition is not added and function returns
 * silently.
 */

void
disk_push_partition(Disk *d, Partition p, mem_arena *arena)
{
  if (d->part_count == d->part_capacity)
  {
    size_t new_cap = d->part_capacity ? d->part_capacity * 2 : 8;

    Partition *np = PUSH_ARRAY_NZ(arena, Partition, new_cap);
    if (!np)
      return;

    if (d->partitions && d->part_count > 0)
    {
      memcpy(np, d->partitions, d->part_count * sizeof(Partition));
    }

    d->partitions = np;
    d->part_capacity = new_cap;
  }

  d->partitions[d->part_count++] = p;
}

/*
 * is_numeric - Check if a string contains only digits
 * @s: String to check
 *
 * Return: 1 if string contains only numeric characters, 0 otherwise
 */
int
is_numeric(const char *s)
{
  for (; *s; ++s)
  {
    if (*s < '0' || *s > '9')
    {
      return 0;
    }
  }
  return 1;
}

/*
 * cpu_create - Allocate and initialize a new Cpu structure
 *
 * Return: Pointer to newly allocated Cpu, or NULL on allocation failure
 */
Cpu *
cpu_create(mem_arena *m)
{
  return arena_push(m, sizeof(Cpu), 1);
}

/*
 * cpu_read - Read CPU information from /proc/cpuinfo
 * @out: Pointer to Cpu structure to populate
 *
 * Reads vendor_id, model name, cpu MHz, and cpu cores from /proc/cpuinfo.
 * The function reads the first occurrence of each field.
 *
 * Return: OK on success, AGENT_ERR_INVALID if out is NULL,
 *         ERR_IO if /proc/cpuinfo cannot be opened
 */

/**
 * get a read of the enabled cpu cores to get a view of the folder structure
 * before starting to read every single of of them
 *
 * call
 *
 * */
int
cpu_read_enabled_core_cpu_frequency(Cpu *out, int enabled_cpu_count)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  char path[PATH_MAX_LEN];
  for (i8 i = 0; i <= enabled_cpu_count; i++)
  {
    snprintf(path, sizeof(path), "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_cur_freq", i);
  }
  FILE *fp = fopen(path, "r");
  // TODO(nasr): do a **cores  see each frequency individually together with the utilization
  if (!fp)
  {
    assert(fp);
    return ERR_IO;
  }

  u64 freq = 0;
  if (fscanf(fp, "%lu", &freq) != 1)
  {
    fclose(fp);
    assert(0);
    return ERR_PARSE;
  }

  fclose(fp);

  snprintf(out->frequency, sizeof(out->frequency), "%lu", freq);
  return OK;
}

int
cpu_read_cpu_model_name_arm64(Cpu *out)
{
  FILE *of = fopen("/proc/device-tree/model", "rb");
  if (!of)
  {
    assert(0);
    return ERR_IO;
  }

  u8 buffer[BUFFER_SIZE_DEFAULT];

  /**
	*
	* note to self
  *	fread returns the amount of bytes read
	* fwrite returns the amount of bytes written
	*
  */

  size_t n = fread(buffer, 1, sizeof(buffer), of);
  if (n == 0)
  {
    assert(0);
    return ERR_IO;
  }

  size_t len = 0;
  while (len < n && buffer[len] != 0)
  {
    ++len;
  }

  memcpy(out->model, buffer, len);
  out->model[len] = '\0';

  fclose(of);

  return OK;
}

int
cpu_get_cores_enabled_arm(Cpu *out)
{
  assert(out);

  FILE *fp = fopen("/sys/devices/system/cpu/enabled", "r");
  assert(fp);

  char buf[BUFFER_SIZE_DEFAULT];
  assert(fgets(buf, sizeof(buf), fp));
  fclose(fp);

  int max_cpu = 0;
  char *p = buf;

  while (*p)
  {
    if (*p >= '0' && *p <= '9')
    {
      int v = 0;
      while (*p >= '0' && *p <= '9')
      {
        v = v * 10 + (*p++ - '0');
      }

      if (v > max_cpu)
      {
        max_cpu = v;
      }
    }
    else
    {
      ++p;
    }
  }

  out->cores = (u32)max_cpu;
  return OK;
}

int
cpu_read_arm64(Cpu *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  cpu_get_cores_enabled_arm(out);
  cpu_read_cpu_model_name_arm64(out);
  cpu_read_enabled_core_cpu_frequency(out, (int)out->cores);

  return OK;
}

int
cpu_read_amd64(Cpu *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *f = fopen("/proc/cpuinfo", "r");
  if (!f)
  {
    assert(0);
    return ERR_IO;
  }

  char buf[BUFFER_SIZE_LARGE];
  while (fgets(buf, sizeof(buf), f))
  {
    char *colon = strchr(buf, ':');
    if (!colon)
      continue;

    char *val = colon + 1;
    while (*val == ' ')
      val++;

    size_t len = strcspn(val, "\n");

    if (!strncmp(buf, "vendor_id", 9))
    {
      memcpy(out->vendor, val, len);
    }
    if (!strncmp(buf, "model name", 10))
    {
      memcpy(out->model, val, len);
    }
    if (!strncmp(buf, "cpu MHz", 7))
    {
      memcpy(out->frequency, val, len);
    }
    if (!strncmp(buf, "cpu cores", 9))
    {
      out->cores = (u32)atoi(buf);
    }
  }

  fclose(f);
  return OK;
}

int
cpu_read(Cpu *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

#if defined(__arm__) || defined(__aarch64__)
  if (cpu_read_arm64(out) != OK)
  {
    /**
		 * Debugging!!
		 *
		 * */
    assert(0);
  }

#elif defined(__i386__) || defined(__x86_64__)
  if (cpu_read_amd64(out) != OK)
  {
    /**
		 * Debugging!!
		 *
		 * */
    assert(0);
  }

#else
#error "Unsupported architecture"
#endif
  return OK;
}

int
cpu_read_usage(Cpu *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *f = fopen("/proc/stat", "r");
  if (!f)
  {
    assert(0);
    return ERR_IO;
  }

  char buf[BUFFER_SIZE_LARGE];
  while (fgets(buf, sizeof(buf), f))
  {
    if (strncmp(buf, "cpu ", 4) == 0)
    {
      unsigned long user, nice, system, idle, iowait, irq, softirq, steal;
      if (sscanf(
            buf + 4,
            "%lu %lu %lu %lu %lu %lu %lu %lu",
            &user,
            &nice,
            &system,
            &idle,
            &iowait,
            &irq,
            &softirq,
            &steal) == 8)
      {
        u64 idleAll = idle + iowait;
        u64 total = user + nice + system + idle + iowait + irq + softirq + steal;
        out->idle_time = idleAll;
        out->total_time = total;
        break;
      }
    }
  }

  fclose(f);
  return OK;
}

/*
 * ram_create - Allocate and initialize a new Ram structure
 *
 * Return: Pointer to newly allocated Ram, or NULL on allocation failure
 */
Ram *
ram_create(mem_arena *m)
{
  return arena_push(m, sizeof(Ram), 1);
}

/*
 * ram_read - Read RAM information from /proc/meminfo
 * @out: Pointer to Ram structure to populate
 *
 * Reads MemTotal and MemFree from /proc/meminfo in kilobytes.
 *
 * Return: OK on success, AGENT_ERR_INVALID if out is NULL,
 *         ERR_IO if /proc/meminfo cannot be opened
 */
int
ram_read(Ram *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *f = fopen("/proc/meminfo", "r");
  if (!f)
  {
    assert(0);
    return ERR_IO;
  }

  mem_arena *temp_arena = arena_create(KiB(8));

  size_t total_len;
  size_t free_len;

  char *total_buffer;
  char *free_buffer;

  char buf[BUFFER_SIZE_SMALL];
  while (fgets(buf, sizeof(buf), f))
  {
    char *colon = strchr(buf, ':');
    if (!colon)
    {
      continue;
    }

    char *val = colon + 1;
    while (*val == ' ')
    {
      val++;
    }

    if (!strncmp(buf, "MemTotal", 8))
    {
      total_len = strcspn(val, "k\n");
      total_buffer = PUSH_ARRAY(temp_arena, char, total_len);

      memcpy(total_buffer, val, total_len);
    }

    if (!strncmp(buf, "MemFree", 7))
    {
      free_len = strcspn(val, "k\n");
      free_buffer = PUSH_ARRAY(temp_arena, char, free_len);

      memcpy(free_buffer, val, free_len);
    }
  }

  out->total = parse_u64(total_buffer, total_len);
  out->free = parse_u64(free_buffer, free_len);

  fclose(f);
  return OK;
}

/*
 * disk_create - Allocate and initialize a new Disk structure
 *
 * Return: Pointer to newly allocated Disk, or NULL on allocation failure
 */
Disk *
disk_create(mem_arena *m)
{
  return arena_push(m, sizeof(Disk), 1);
}

/*
 * disk_read - Read disk partition information from /proc/partitions
 * @out: Pointer to Disk structure to populate
 *
 * Reads all partitions from /proc/partitions, storing major/minor device
 * numbers, block count, and device name for each partition. Skips the header
 * line and any malformed entries.
 *
 * Return: OK on success, AGENT_ERR_INVALID if out is NULL,
 *         ERR_IO if /proc/partitions cannot be opened
 */

int
disk_read(Disk *out, mem_arena *arena)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *f = fopen("/proc/partitions", "r");
  if (!f)
  {
    assert(0);
    return ERR_IO;
  }

  char buf[BUFFER_SIZE_DEFAULT];

  while (fgets(buf, sizeof(buf), f))
  {
    Partition p = { 0 };
    char name[BUFFER_SIZE_DEFAULT];

    if (sscanf(buf,
          "%lu %lu %lu %255s",
          &p.major,
          &p.minor,
          &p.blocks,
          name) != 4)
    {
      continue;
    }

    size_t len = strlen(name);
    if (len >= sizeof(p.name))
      len = sizeof(p.name) - 1;

    memcpy(p.name, name, len);
    p.name[len] = 0;

    disk_push_partition(out, p, arena);
  }

  fclose(f);
  return OK;
}

int
fs_usage(char *path, Disk *disk)
{
  struct statfs s;
  if (statfs(path, &s) != 0)
  {
    return ERR_IO;
  }

  i64 block_size = s.f_bsize;

  u64 blocks = (u64)s.f_blocks;
  u64 bfree = (u64)s.f_bfree;
  u64 bavail = (u64)s.f_bavail;
  u64 bsize = (u64)block_size;

  disk->disk_usage.total = blocks * bsize;
  disk->disk_usage.free = bfree * bsize;
  disk->disk_usage.available = bavail * bsize;
  disk->disk_usage.used = (blocks - bfree) * bsize;

  return OK;
}

/*
 * device_create - Allocate and initialize a new Device structure
 *
 * Return: Pointer to newly allocated Device, or NULL on allocation failure
 */
Device *
device_create(mem_arena *m)
{
  return arena_push(m, sizeof(Device), 1);
}

/*
 * collect_processes - Collect all running process IDs from /proc
 * @dev: Pointer to Device structure to populate with process IDs
 *
 * Scans /proc directory for numeric entries (process IDs) and stores them
 * as strings in the Device structure. Dynamically grows the process array
 * as needed.
 */

int
process_list_collect(Process_List *list, mem_arena *arena)
{
  DIR *d = opendir("/proc");
  if (!d)
  {
    assert(0);
    return ERR_IO;
  }

  struct dirent *e = 0;

  if (!list->items)
  {
    list->capacity = 8;
    list->count = 0;
    list->items = PUSH_ARRAY_NZ(arena, Process, list->capacity);
    if (!list->items)
    {
      closedir(d);
      assert(0);
      return ERR_IO;
    }
  }

  while ((e = readdir(d)))
  {
    if (!is_numeric(e->d_name))
      continue;

    if (list->count == list->capacity)
    {
      size_t new_cap = list->capacity * 2;
      Process *np = PUSH_ARRAY_NZ(arena, Process, new_cap);
      if (!np)
        break;

      memcpy(np, list->items, sizeof(Process) * list->capacity);
      list->items = np;
      list->capacity = new_cap;
    }

    Process *p = &list->items[list->count++];

    p->pid = atoi(e->d_name);
    p->state = PROCESS_UNDEFINED;
    p->utime = 0;
    p->stime = 0;
    p->num_threads = 0;
    p->name[0] = 0;
  }

  closedir(d);
  return OK;
}

/**
 *
struct Proces {
  char *pid;
        char *name;
  char *state;
  char *utime;
  char *num_threads;
};

*/

int
process_read(i32 pid, Process *out)
{
  char path[PATH_MAX_LEN];
  snprintf(path, sizeof(path), "/proc/%d/status", pid);

  FILE *fp = fopen(path, "r");
  if (!fp)
  {
    return ERR_IO;
  }

  char buf[BUFFER_SIZE_LARGE];

  /* initialize */
  out->pid = pid;

  while (fgets(
    buf,
    sizeof(buf),
    fp))
  {
    char *colon = strchr(buf, ':');
    if (!colon)
    {
      continue;
    }

    char *val = colon + 1;
    while (*val == ' ' || *val == '\t')
    {
      ++val;
    }

    size_t len = strcspn(val, "\n");

    if (!strncmp(buf, "Name:", 5))
    {
      memcpy(out->name, val, len);
      out->name[len] = 0;
    }
    if (!strncmp(buf, "State:", 6))
    {
      char state_char = 0;
      for (char *p = val; *p; ++p)
      {
        if ((*p >= 'A' && *p <= 'Z') || *p == 't')
        {
          state_char = *p;
          break;
        }
      }

      switch (state_char)
      {
        case 'R':
        {
          out->state = PROCESS_RUNNING;
          break;
        }
        case 'S':
        {
          out->state = PROCESS_SLEEPING;
          break;
        }
        case 'D':
        {
          out->state = PROCESS_DISK_SLEEP;
          break;
        }
        case 'T':
        {
          out->state = PROCESS_STOPPED;
          break;
        }
        case 't':
        {
          out->state = PROCESS_TRACING_STOPPED;
          break;
        }
        case 'Z':
        {
          out->state = PROCESS_ZOMBIE;
          break;
        }
        case 'X':
        {
          out->state = PROCESS_DEAD;
          break;
        }
        case 'I':
        {
          out->state = PROCESS_IDLE;
          break;
        }
        default:
        {
          out->state = PROCESS_UNDEFINED;
          break;
        }
      }
    }

    if (!strncmp(buf, "Threads:", 8))
    {
      out->num_threads = (u32)strtoul(val, 0, 10);
    }
  }

  int error = fclose(fp);
  if (error != 0)
  {
    return ERR_IO;
  }

  return OK;
}

int
device_up_time(Device *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *f = fopen("/proc/uptime", "r");
  if (!f)
  {
    assert(0);
    return ERR_IO;
  }

  i64 s;
  int result = fscanf(f, "%ld", &s);
  fclose(f);

  if (result != 1)
  {
    assert(0);
    return ERR_IO;
  }

  i64 day = s / 86400;
  i64 hour = s % 86400 / 3600;
  i64 min = s % 3600 / 60;
  sprintf(out->uptime, "%ldd %ldh %ldm", day, hour, min);
  return OK;
}

/*
 * device_read - Read device information including OS version, uptime, and
 * processes
 * @out: Pointer to Device structure to populate
 *
 * Reads OS version from /proc/version, uptime from /proc/uptime, and collects
 * all running process IDs from /proc directory.
 *
 * Return: OK on success, AGENT_ERR_INVALID if out is NULL,
 *         ERR_IO if required files cannot be opened
 */
int
device_read(Device *out)
{
  if (!out)
  {
    assert(0);
    return ERR_INVALID;
  }

  FILE *version = fopen("/etc/os-release", "r");
  if (!version)
  {
    if (version)
    {
      fclose(version);
    }
    return ERR_IO;
  }

  char buffer[BUFFER_SIZE_DEFAULT];
  while (fgets(buffer, sizeof(out->os_version), version))
  {
    if (!strncmp(buffer, "NAME=", 5))
    {
      char *start = strchr(buffer, '"');
      if (start)
      {
        start++;
        char *end = strchr(start, '"');
        if (end)
        {
          i64 len = end - start;
          if ((u64)len < sizeof(out->os_version))
          {
            memcpy(out->os_version, start, (u64)len);
            out->os_version[len] = '\0';
          }
        }
      }
    }
  }

  fclose(version);
  device_up_time(out);

  return OK;
}

/*
 * process_kill - Send a signal to a process
 * @pid: Process ID to signal
 * @signal: Signal number to send (e.g., SIGTERM, SIGKILL)
 *
 * Wrapper around kill(2) system call with error handling.
 *
 * Return: OK on success,
 *         ERR_INVALID if pid is invalid or process not found,
 *         ERR_PERM if permission denied,
 *         ERR_IO for other errors
 */
int
process_kill(pid_t pid, int signal)
{
  if (pid <= 0)
  {
    assert(0);
    return ERR_INVALID;
  }

  if (kill(pid, signal) == -1)
  {
    if (errno == EPERM)
    {
      assert(0);
      return ERR_PERM;
    }
    if (errno == ESRCH)
    {
      assert(0);
      return ERR_INVALID;
    }

    assert(0);
    return ERR_IO;
  }
  return OK;
}
