package service

import (
	"context"
	"fmt"

	resourcespb "paradigm-ehb/agent/gen/resources/v1"
	"paradigm-ehb/agent/internal/resmanager"
)

type ResourcesService struct {
	resourcespb.UnimplementedResourcesServiceServer
}

func (s *ResourcesService) GetSystemResources(
	ctx context.Context,
	req *resourcespb.GetSystemResourcesRequest,
) (*resourcespb.GetSystemResourcesResponse, error) {

	sys, err := resources.Make()
	if err != nil {
		return nil, fmt.Errorf("resource snapshot failed: %w", err)
	}

	return &resourcespb.GetSystemResourcesResponse{
		Resources: &resourcespb.SystemResources{
			Cpu: &resourcespb.Cpu{
				Vendor:    sys.CPU.Vendor,
				Model:     sys.CPU.Model,
				Frequency: sys.CPU.Frequency,
				MaxCore:   sys.CPU.MaxCore,
			},
			Memory: &resourcespb.Memory{
				Total: sys.Memory.Total,
				Free:  sys.Memory.Free,
			},
			Device: &resourcespb.Device{
				OsVersion: sys.Device.OsVersion,
				Uptime:    sys.Device.Uptime,
			},
			Disks:     mapDisks(sys.Disks),
			Processes: mapProcesses(sys.Procs),
		},
	}, nil
}

func mapDisks(disks []resources.Disk) []*resourcespb.Disk {
	out := make([]*resourcespb.Disk, 0, len(disks))
	for _, d := range disks {
		pbDisk := &resourcespb.Disk{
			Partitions: make([]*resourcespb.DiskPartition, 0, len(d.Partitions)),
		}

		for _, p := range d.Partitions {
			pbDisk.Partitions = append(pbDisk.Partitions, &resourcespb.DiskPartition{
				Name:   p.Name,
				Major:  p.Major,
				Minor:  p.Minor,
				Blocks: p.Blocks,
			})
		}

		out = append(out, pbDisk)
	}
	return out
}

func mapProcesses(ps []resources.Process) []*resourcespb.Process {
	out := make([]*resourcespb.Process, 0, len(ps))
	for _, p := range ps {
		out = append(out, &resourcespb.Process{
			Pid:        p.PID,
			Name:       p.Name,
			State:      mapProcessState(p.State),
			Utime:      p.UTime,
			NumThreads: p.NumThreads,
		})
	}
	return out
}

func mapProcessState(s resources.ProcesState) resourcespb.ProcessState {
	switch s {
	case 'R':
		return resourcespb.ProcessState_PROCESS_STATE_RUNNING
	case 'S':
		return resourcespb.ProcessState_PROCESS_STATE_SLEEPING
	case 'D':
		return resourcespb.ProcessState_PROCESS_STATE_DISK_SLEEPING
	case 'T':
		return resourcespb.ProcessState_PROCESS_STATE_STOPPED
	case 't':
		return resourcespb.ProcessState_PROCESS_STATE_TRACING_STOPPED
	case 'Z':
		return resourcespb.ProcessState_PROCESS_STATE_ZOMBIE
	case 'X':
		return resourcespb.ProcessState_PROCESS_STATE_DEAD
	default:
		return resourcespb.ProcessState_PROCESS_STATE_UNSPECIFIED
	}
}
