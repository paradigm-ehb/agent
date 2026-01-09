package wrapper

/*
#cgo CFLAGS:  -I${SRCDIR}/../agent-resources
#cgo LDFLAGS: -L${SRCDIR}/../agent-resources/build -lagent_resources

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "resources.h"

int
process_read2(i32 pid, Process *out)
{
  char path[PATH_MAX_LEN];
  snprintf(path, sizeof(path), "/proc/%d/status", pid);

  FILE *fp = fopen(path, "r");
  if (!fp)
  {
    return ERR_IO;
  }

  char buf[BUFFER_SIZE_LARGE];

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

    printf("lenght: %lu", len);

    if (!strncmp(buf, "Name:", 5))
    {

      memcpy(out->name, val, len);
      printf("\nout name %s\n", out->name);
    }
    printf("name: %s", out->name);
    if (!strncmp(buf, "State:", 6))
    {
      char state_char = 0;
      for (char *p = val; *p; ++p)
      {
        if (*p >= 'A' && *p <= 'Z' || *p == 't')
        {
          state_char = *p;
          break;
        }
      }


      printf("\nstate char: %d\n", out->pid);
      printf("\nstate char: %c\n", state_char);
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

  printf("sizeof(Process) = %zu\n", sizeof(Process));

  int error = fclose(fp);
  if (error != 0)
  {
    return ERR_IO;
  }

  printf("\n\n\n\nout state: %d\n", out->state);
  return OK;
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ProcessState represents the state of a process.

type ProcessState C.int32_t

const (
	ProcessUndefined      ProcessState = ProcessState(C.PROCESS_UNDEFINED)
	ProcessRunning        ProcessState = ProcessState(C.PROCESS_RUNNING)
	ProcessSleeping       ProcessState = ProcessState(C.PROCESS_SLEEPING)
	ProcessDiskSleep      ProcessState = ProcessState(C.PROCESS_DISK_SLEEP)
	ProcessStopped        ProcessState = ProcessState(C.PROCESS_STOPPED)
	ProcessTracingStopped ProcessState = ProcessState(C.PROCESS_TRACING_STOPPED)
	ProcessZombie         ProcessState = ProcessState(C.PROCESS_ZOMBIE)
	ProcessDead           ProcessState = ProcessState(C.PROCESS_DEAD)
	ProcessIdle           ProcessState = ProcessState(C.PROCESS_IDLE)
)

// Process represents a single process with its attributes.
type Process struct {
	PID        int32
	State      ProcessState
	UTime      uint64
	STime      uint64
	NumThreads uint32
	Name       string
}

func (s ProcessState) String() string {

	switch s {
	case ProcessRunning:
		return "Running"
	case ProcessSleeping:
		return "Sleeping"
	case ProcessDiskSleep:
		return "Disk Sleep"
	case ProcessStopped:
		return "Stopped"
	case ProcessTracingStopped:
		return "Tracing Stopped"
	case ProcessZombie:
		return "Zombie"
	case ProcessDead:
		return "Dead"
	case ProcessIdle:
		return "Idle"
	default:
		return "Undefined"
	}
}

// TODO(nasr): research this, interesting, alias vs true aliasing
type Arena = C.mem_arena

/*
*
KiB converts n to kibibytes (n * 1024).
*/
func KiB(n uint64) uint64 {
	return n << 10
}

/*
*
MiB converts n to mebibytes (n * 1024 * 1024).
*/
func MiB(n uint64) uint64 {
	return n << 20
}

/*
*
GiB converts n to gibibytes (n * 1024 * 1024 * 1024).
*/
func GiB(n uint64) uint64 {
	return n << 30
}

/*
*
AllocateArena creates a new memory arena with 8 MiB capacity.
The arena uses mmap for memory allocation and must be destroyed with DestroyArena.

Returns:
  - *C.mem_arena: Pointer to the allocated arena
  - error: Error if allocation fails
*/
func AllocateArena(size uint64) (*C.mem_arena, error) {

	arena := C.arena_create(C.ulong(size))

	if arena == nil {
		return nil, fmt.Errorf("failed to allocate the arena")
	}
	return arena, nil
}

func PushArena(arena *C.mem_arena, size uint64) {

	if arena != nil {
		C.arena_push(arena, C.ulong(size), 1)
	}

}

/*
*
DestroyArena unmaps and destroys the memory arena.
This should be called to clean up resources when the arena is no longer needed.

Parameters:
  - arena: Pointer to the arena to destroy
*/
func DestroyArena(arena *C.mem_arena) {
	if arena != nil {
		C.arena_destroy(arena)
	}
}

/*
*
ClearArena resets the arena position to its base, effectively clearing all allocations.
This does not free the memory but allows reusing the same arena space.

Parameters:
  - arena: Pointer to the arena to clear
*/
func ClearArena(arena *C.mem_arena) {
	if arena != nil {
		C.arena_clear(arena)
	}
}

/*
*
CpuCreate allocates a Cpu structure in the arena.

Parameters:
  - arena: Memory arena to allocate from

Returns:
  - *C.Cpu: Pointer to the allocated Cpu structure
  - error: Error if allocation fails
*/
func CpuCreate(arena *C.mem_arena) (*C.Cpu, error) {
	c := C.cpu_create(arena)
	if c == nil {
		return nil, fmt.Errorf("failed to push a Cpu object onto the arena stack")
	}
	return c, nil
}

/*
*
CpuRead reads CPU information from the system and returns a Go Cpu struct.
This reads from /proc/cpuinfo on AMD64 or /proc/device-tree and /sys on ARM64.

Parameters:
  - c: Pointer to the C Cpu structure to read into

Returns:
  - Cpu: Go struct containing CPU information
  - error: Error if reading fails
*/
func CpuRead(c *C.Cpu) (Cpu, error) {
	if c == nil {
		return Cpu{}, fmt.Errorf("nil Cpu pointer")
	}

	if C.cpu_read(c) != C.OK {
		return Cpu{}, fmt.Errorf("failed to read CPU information")
	}

	cpu := Cpu{
		Vendor:    C.GoString(&c.vendor[0]),
		Model:     C.GoString(&c.model[0]),
		Frequency: C.GoString(&c.frequency[0]),
		MaxCore:   uint32(c.cores),
	}
	return cpu, nil
}

/*
*
RamCreate allocates a Ram structure in the arena and reads RAM information.

Parameters:
  - arena: Memory arena to allocate from

Returns:
  - *C.Ram: Pointer to the allocated and populated Ram structure
  - error: Error if allocation or reading fails
*/
func RamCreate(arena *C.mem_arena) (*C.Ram, error) {
	r := C.ram_create(arena)
	if r == nil {
		return nil, fmt.Errorf("failed to create ram object")
	}
	if C.ram_read(r) != C.OK {
		return nil, fmt.Errorf("failed to read ram information")
	}
	return r, nil
}

/*
*
RamRead converts a C Ram structure to a Go Ram struct.
This reads total and free memory from /proc/meminfo.

Parameters:
  - ram: Pointer to the C Ram structure

Returns:
  - Ram: Go struct containing RAM information in kilobytes
  - error: Error if ram pointer is nil
*/
func RamRead(ram *C.Ram) (Ram, error) {
	if ram == nil {
		return Ram{}, fmt.Errorf("nil Ram pointer")
	}

	r := Ram{
		Total: C.GoString(&ram.total[0]),
		Free:  C.GoString(&ram.free[0]),
	}
	return r, nil
}

/**
Disk Functions

DiskCreate allocates a Disk structure in the arena and reads partition information.

Parameters:
  - arena: Memory arena to allocate from

Returns:
  - *C.Disk: Pointer to the allocated Disk structure
  - error: Error if allocation fails
*/

func DiskCreate(arena *C.mem_arena) (*C.Disk, error) {
	di := C.disk_create(arena)
	if di == nil {
		return nil, fmt.Errorf("failed to create disk object")
	}
	if C.disk_read(di, arena) != C.OK {
		return nil, fmt.Errorf("failed to read disk information")
	}
	return di, nil
}

/*
*
DiskRead converts a C Disk structure to a Go Disk struct with all partitions.
This reads partition information from /proc/partitions.

Parameters:
  - disk: Pointer to the C Disk structure

Returns:
  - Disk: Go struct containing all disk partitions
  - error: Error if disk pointer is nil
*/
func DiskRead(disk *C.Disk) (Disk, error) {
	if disk == nil {
		return Disk{}, fmt.Errorf("nil Disk pointer")
	}

	d := Disk{
		Partitions: make([]DiskPartition, 0, disk.count),
	}

	for i := C.size_t(0); i < disk.count; i++ {
		part := (*C.Partition)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(disk.partitions)) +
					uintptr(i)*unsafe.Sizeof(*disk.partitions),
			),
		)
		d.Partitions = append(d.Partitions, DiskPartition{
			Name:   C.GoString(&part.name[0]),
			Major:  uint32(part.major),
			Minor:  uint32(part.minor),
			Blocks: uint64(part.blocks),
		})
	}
	return d, nil
}

/**
DeviceCreate allocates a Device structure in the arena.

Parameters:
  - arena: Memory arena to allocate from

Returns:
  - *C.Device: Pointer to the allocated Device structure
  - error: Error if allocation fails

*/

func DeviceCreate(arena *C.mem_arena) (*C.Device, error) {
	de := C.device_create(arena)
	if de == nil {
		return nil, fmt.Errorf("failed to create device")
	}
	return de, nil
}

/*
*
DeviceRead reads device information including OS version and uptime.
This reads from /etc/os-release and /proc/uptime.

Parameters:
  - device: Pointer to the C Device structure
  - arena: Memory arena (not currently used but kept for consistency)

Returns:
  - Device: Go struct containing device information
  - error: Error if reading fails
*/
func DeviceRead(device *C.Device, arena *C.mem_arena) (Device, error) {
	if device == nil {
		return Device{}, fmt.Errorf("nil Device pointer")
	}

	if C.device_read(device) != C.OK {
		return Device{}, fmt.Errorf("failed to read device information")
	}

	de := Device{
		OsVersion: C.GoString(&device.os_version[0]),
		Uptime:    C.GoString(&device.uptime[0]),
	}
	return de, nil
}

/*
*
Process Functions

FetchProcesses collects all running process IDs from /proc directory.
This populates the device's process list but does not read detailed information.

Parameters:
  - device: Pointer to the C Device structure to populate
  - arena: Memory arena to allocate process list from

Returns:
  - error: Error if collection fails
*/
func FetchProcesses(device *C.Device, arena *C.mem_arena) error {
	if device == nil {
		return fmt.Errorf("nil Device pointer")
	}
	if arena == nil {
		return fmt.Errorf("nil arena pointer")
	}

	if C.process_list_collect(&device.processes, arena) != C.OK {
		return fmt.Errorf("failed to collect process list")
	}
	return nil
}

/*
*
ReadProcesses reads detailed information for all processes in the device's process list.
This reads from /proc/[pid]/status for each process.

Parameters:
  - device: Pointer to the C Device structure containing the process list

Returns:
  - []Process: Slice of Process structs with detailed information
  - error: Error if device pointer is nil
*/
func ReadProcesses(device *C.Device) ([]Process, error) {


	if device == nil {
		return nil, fmt.Errorf("dvice null pointer")
	}

	count := int(device.processes.count)
	items := device.processes.items

	// slice := unsafe.Slice(items, count)

	procs := make([]Process, 0, count)

	// p := unsafe.Pointer(&items[i])

	// p := (*C.Process)(unsafe.Pointer(items), uintptr(i)*unsafe.Sizeof(*items))

	for i := 0; i < count; i++ {
		// p := &slice[i]


	p := (*C.Process)(unsafe.Pointer(uintptr(unsafe.Pointer(items)) + uintptr(i)*unsafe.Sizeof(*items)))


		err := C.process_read2(p.pid, p)
		if err != C.OK {
			fmt.Errorf("failed reading processes")
		}

	fmt.Println("Reading processes...", device.processes)


		fmt.Println("pid:",p.pid)
		fmt.Println("state:", p.state)

		// if p.state == C.PROCESS_UNDEFINED {
		//
		// 	C.process_read(p.pid, p)
		// }
		//

		fmt.Println("\t\t\t\tstate: ", C.GoString(&p.name[0]))
		fmt.Println(unsafe.Sizeof(C.Process{}))


		fmt.Println("pid offset  ", unsafe.Offsetof(p.pid))
		fmt.Println("state offset", unsafe.Offsetof(p.state))
		fmt.Println("utime offset", unsafe.Offsetof(p.utime))

		procs = append(procs, Process{
			PID:        int32(p.pid),
			State:      ProcessState(p.state), //BREAKPOINT
			UTime:      uint64(p.utime),
			STime:      uint64(p.stime),
			NumThreads: uint32(p.num_threads),
			Name:       C.GoString(&p.name[0]),
		})

	}

	return procs, nil
}

/*
*
KillProcess sends a SIGTERM signal to the specified process.

Parameters:
  - pid: Process ID to terminate

Returns:
  - error: Error if the process cannot be killed or doesn't exist
*/
func KillProcess(pid int) error {

	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	if C.process_kill(C.int(pid), C.int(unix.SIGTERM)) != C.OK {
		return fmt.Errorf("failed to kill process %d", pid)
	}
	return nil
}
