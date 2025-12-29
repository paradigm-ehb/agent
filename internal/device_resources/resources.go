package resources

import (
	wr "paradigm-ehb/agent/pkg/wrapper"
	// "unsafe"
	// "golang.org/x/sys/unix"
)

type SystemResources struct {
	CPU    wr.Cpu
	Memory wr.Ram
	Device wr.Device
	Disks  []wr.Disk
	Procs  []wr.Process
}

func Make() {}
