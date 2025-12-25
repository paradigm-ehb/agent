package resources

/*
#cgo CFLAGS:  -I${SRCDIR}/agent-resources
#cgo LDFLAGS: -L${SRCDIR}/agent-resources/build -lagent_resources
#include "agent_resources.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Cpu struct {
	Vendor    string
	Model     string
	Frequency string
	MaxCore   uint32
}

type Memory struct {
	Total string
	Free  string
}

type DiskPartition struct {
	Name   string
	Major  uint32
	Minor  uint32
	Blocks uint64
}

type Disk struct {
	Partitions []DiskPartition
}

type Device struct {
	OsVersion string
	Uptime    string
}

type ProcesState rune

const (
	Running        ProcesState = 'R'
	Sleeping       ProcesState = 'S'
	DiskSleeping   ProcesState = 'D'
	Stopped        ProcesState = 'T'
	TracingStopped ProcesState = 't'
	Zombie         ProcesState = 'Z'
	Dead           ProcesState = 'X'
)

type Process struct {
	PID        uint32
	Name       string
	State      ProcesState
	UTime      uint64
	NumThreads uint32
}

type SystemResources struct {
	CPU    Cpu
	Memory Memory
	Device Device
	Disks  []Disk
	Procs  []Process
}

func Make() (*SystemResources, error) {
	arena := C.arena_create(C.MiB(8))
	if arena == nil {
		return nil, fmt.Errorf("arena_create failed")
	}
	defer C.arena_destroy(arena)

	ram := C.ram_create(arena)
	if ram == nil || C.ram_read(ram) != C.OK {
		return nil, fmt.Errorf("ram_read failed")
	}

	memGo := Memory{
		Total: C.GoString(&ram.total[0]),
		Free:  C.GoString(&ram.free[0]),
	}

	cpu := C.cpu_create(arena)
	if cpu == nil || C.cpu_read(cpu) != C.OK {
		return nil, fmt.Errorf("cpu_read failed")
	}

	cpuGo := Cpu{
		Vendor:    C.GoString(&cpu.vendor[0]),
		Model:     C.GoString(&cpu.model[0]),
		Frequency: C.GoString(&cpu.frequency[0]),
		MaxCore:   uint32(cpu.cores),
	}

	device := C.device_create(arena)
	if device == nil || C.device_read(device) != C.OK {
		return nil, fmt.Errorf("device_read failed")
	}

	deviceGo := Device{
		OsVersion: C.GoString(&device.os_version[0]),
		Uptime:    C.GoString(&device.uptime[0]),
	}

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
			State:      ProcesState(p.state),
			UTime:      uint64(p.utime),
			NumThreads: uint32(p.num_threads),
		})
	}

	disk := C.disk_create(arena)
	if disk == nil || C.disk_read(disk, arena) != C.OK {
		return nil, fmt.Errorf("disk_read failed")
	}

	diskGo := Disk{
		Partitions: make([]DiskPartition, 0, disk.count),
	}

	for i := C.size_t(0); i < disk.count; i++ {
		part := (*C.Partition)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(disk.partitions)) +
					uintptr(i)*unsafe.Sizeof(*disk.partitions),
			),
		)

		diskGo.Partitions = append(diskGo.Partitions, DiskPartition{
			Name:   C.GoString(&part.name[0]),
			Major:  uint32(part.major),
			Minor:  uint32(part.minor),
			Blocks: uint64(part.blocks),
		})
	}

	return &SystemResources{
		CPU:    cpuGo,
		Memory: memGo,
		Device: deviceGo,
		Disks:  []Disk{diskGo},
		Procs:  procs,
	}, nil
}
