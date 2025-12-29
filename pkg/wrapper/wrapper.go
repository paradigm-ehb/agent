package wrapper

/*
#cgo CFLAGS:  -I${SRCDIR}/../agent-resources
#cgo LDFLAGS: -L${SRCDIR}/../agent-resources/build -lagent_resources
#include "agent_resources.h"
*/
import "C"
import (
	"unsafe"

	"fmt"
	"golang.org/x/sys/unix"
)

func KiB(n uint64) uint64 {

	return n << 10
}

func MiB(n uint64) uint64 {

	return n << 20
}

func GiB(n uint64) uint64 {

	return n << 30
}

func AllocateArena() (*C.mem_arena, error) {

	arena := C.arena_create(C.MiB(8))
	if arena == nil {
		return nil, fmt.Errorf("failed to allocate the arena")
	}

	return arena, nil
}

func DestroyArena(arena *C.mem_arena) {

	C.arena_destroy(arena)

}

func ClearArena(arena *C.mem_arena) {

	C.arena_clear(arena)

}

func CpuCreate(arena *C.mem_arena) (*C.Cpu, error) {

	c := C.cpu_create(arena)
	if c == nil {

		return c, fmt.Errorf("failed to push a Cpu object onto the arena stack")
	}

	return c, nil
}

func RamCreate(arena *C.mem_arena) (*C.Ram, error) {

	r := C.ram_create(arena)
	if r == nil || C.ram_read(r) != C.OK {
		return r, fmt.Errorf("failed to create ram object")
	}

	return r, nil
}

func DiskCreate(arena *C.mem_arena) (*C.Disk, error) {

	di := C.disk_create(arena)

	if di == nil {
		return di, fmt.Errorf("failed to  create disk object")
	}

	return di, nil
}

func DeviceCreate(arena *C.mem_arena) (*C.Device, error) {

	de := C.device_create(arena)
	if de == nil {
		return de, fmt.Errorf("failed to create device")
	}
	return de, nil

}

func CpuRead(c *C.Cpu) (Cpu, error) {

	var cpu Cpu

	C.cpu_read(c)

	cpu = Cpu{
		Vendor:    C.GoString(&c.vendor[0]),
		Model:     C.GoString(&c.model[0]),
		Frequency: C.GoString(&c.frequency[0]),
		MaxCore:   uint32(c.cores),
	}

	return cpu, nil

}

func RamRead(ram *C.Ram) (Ram, error) {

	r := Ram{
		Total: C.GoString(&ram.total[0]),
		Free:  C.GoString(&ram.free[0]),
	}
	return r, nil

}

func DiskRead(disk *C.Disk) (Disk, error) {

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

func DeviceRead(device *C.Device, arena *C.mem_arena) (Device, error) {

	var de Device
	err := C.device_read(device)
	if err != C.OK {

		return de, fmt.Errorf("failed to read device")
	}

	de = Device{
		OsVersion: C.GoString(&device.os_version[0]),
		Uptime:    C.GoString(&device.uptime[0]),
	}

	return de, nil

}

func FetchProcesses(device *C.Device, arena *C.mem_arena) {

	C.process_list_collect(&device.processes, arena)

	procs := make([]Process, 0, device.processes.count)
	items := device.processes.items

	for i := C.size_t(0); i < device.processes.count; i++ {
		p := (*C.Process)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(items)) +
					uintptr(i)*unsafe.Sizeof(*items),
			),
		)

		if C.process_read(p.pid, p) != C.OK {
			continue
		}

		procs = append(procs, Process{
			PID:        uint32(p.pid),
			Name:       C.GoString(&p.name[0]),
			State:      ProcessState(p.state),
			UTime:      uint64(p.utime),
			NumThreads: uint32(p.num_threads),
		})
	}

}

func ReadProcesses(pid int32) {

}

func KillProcess(pid int) error {

	if C.process_kill(C.int(pid), C.int(unix.SIGTERM)) != C.OK {
		return nil
	}
	return nil
}

func Make() {

}
