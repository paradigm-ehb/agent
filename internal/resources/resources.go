package resources

import (
	"log"
	wr "paradigm-ehb/agent/pkg/cgowrap"
)

// SystemResources aggregates all system resource information.
type SystemResources struct {
	CPU    wr.Cpu
	Memory wr.Ram
	Device wr.Device
	Disks  []wr.Disk
	Procs  []wr.Process
}

func GetCpu(arena *wr.Arena) {

	cpu, err := wr.CpuCreate(arena)
	if err != nil {

		log.Fatal("Failed to allocate an arena for a complete snapshot")
	}

	wr.CpuRead(cpu)

}

func GetRam(arena *wr.Arena) {

	ram, err := wr.RamCreate(arena)
	if err != nil {

		log.Fatal("Failed to allocate an arena for a complete snapshot")
	}

	wr.RamRead(ram)

}

func GetDisk(arena *wr.Arena) {

	disk, err := wr.DiskCreate(arena)
	if err != nil {

		log.Fatal("Failed to allocate an arena for a complete snapshot")
	}

	wr.DiskRead(disk)

}

func GetDevice(arena *wr.Arena) {

	device, err := wr.DeviceCreate(arena)
	if err != nil {

		log.Fatal("Failed to allocate an arena for a complete snapshot")
	}

	wr.DeviceRead(device, arena)

}

func GetCompleteSystemResources() (*SystemResources, error) {

	var sys SystemResources

	arena, err := wr.AllocateArena(wr.MiB(2))
	if err != nil {
		log.Fatal("Failed to allocate an arena for a complete snapshot")
	}

	defer wr.DestroyArena(arena)

	return &sys, nil
}

func UpdateCpuFrequency() (string, error) {

	return "", nil
}

func UpdateMemoryUsage() (string, error) {

	return "", nil
}
