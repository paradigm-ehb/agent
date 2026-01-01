package resources

import (
	"fmt"

	wr "paradigm-ehb/agent/pkg/cgowrap"
)

type SystemResources struct {
	CPU    wr.Cpu
	Memory wr.Ram
	Device wr.Device
	Disks  []wr.Disk
	Procs  []wr.Process
}

// TODO(nasr): improve understanding of lifetimes
func withArena(size uint64, fn func(*wr.Arena) error) error {
	arena, err := wr.AllocateArena(size)
	if err != nil {
		return err
	}
	defer wr.DestroyArena(arena)
	return fn(arena)
}

func GetCPU() (wr.Cpu, error) {
	var out wr.Cpu

	err := withArena(wr.MiB(1), func(arena *wr.Arena) error {
		c, err := wr.CpuCreate(arena)
		if err != nil {
			return err
		}

		cpu, err := wr.CpuRead(c)
		if err != nil {
			return err
		}

		out = cpu
		return nil
	})

	return out, err
}

func GetRAM() (wr.Ram, error) {
	var out wr.Ram

	err := withArena(wr.MiB(1), func(arena *wr.Arena) error {
		r, err := wr.RamCreate(arena)
		if err != nil {
			return err
		}

		ram, err := wr.RamRead(r)
		if err != nil {
			return err
		}

		out = ram
		return nil
	})

	return out, err
}

func GetDisks() ([]wr.Disk, error) {
	var out []wr.Disk

	err := withArena(wr.MiB(2), func(arena *wr.Arena) error {
		d, err := wr.DiskCreate(arena)
		if err != nil {
			return err
		}

		disk, err := wr.DiskRead(d)
		if err != nil {
			return err
		}

		out = []wr.Disk{disk}
		return nil
	})

	return out, err
}

func GetDevice() (wr.Device, error) {
	var out wr.Device

	err := withArena(wr.MiB(1), func(arena *wr.Arena) error {
		d, err := wr.DeviceCreate(arena)
		if err != nil {
			return err
		}

		device, err := wr.DeviceRead(d, arena)
		if err != nil {
			return err
		}

		out = device
		return nil
	})

	return out, err
}

func GetProcesses() ([]wr.Process, error) {
	var out []wr.Process

	err := withArena(wr.MiB(4), func(arena *wr.Arena) error {
		d, err := wr.DeviceCreate(arena)
		if err != nil {
			return err
		}

		if err := wr.FetchProcesses(d, arena); err != nil {
			return err
		}

		procs, err := wr.ReadProcesses(d)
		if err != nil {
			return err
		}

		out = procs
		return nil
	})

	return out, err
}

func GetCompleteSystemResources() (*SystemResources, error) {
	var snap SystemResources

	err := withArena(wr.MiB(8), func(arena *wr.Arena) error {

		c, err := wr.CpuCreate(arena)
		if err != nil {
			return fmt.Errorf("cpu create: %w", err)
		}
		snap.CPU, err = wr.CpuRead(c)
		if err != nil {
			return fmt.Errorf("cpu read: %w", err)
		}

		r, err := wr.RamCreate(arena)
		if err != nil {
			return fmt.Errorf("ram create: %w", err)
		}
		snap.Memory, err = wr.RamRead(r)
		if err != nil {
			return fmt.Errorf("ram read: %w", err)
		}

		dv, err := wr.DeviceCreate(arena)
		if err != nil {
			return fmt.Errorf("device create: %w", err)
		}
		snap.Device, err = wr.DeviceRead(dv, arena)
		if err != nil {
			return fmt.Errorf("device read: %w", err)
		}

		dk, err := wr.DiskCreate(arena)
		if err != nil {
			return fmt.Errorf("disk create: %w", err)
		}
		disk, err := wr.DiskRead(dk)
		if err != nil {
			return fmt.Errorf("disk read: %w", err)
		}
		snap.Disks = []wr.Disk{disk}

		if err := wr.FetchProcesses(dv, arena); err != nil {
			return fmt.Errorf("process list: %w", err)
		}
		snap.Procs, err = wr.ReadProcesses(dv)
		if err != nil {
			return fmt.Errorf("process read: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &snap, nil
}
