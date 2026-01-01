package service

import (
	"context"

	proto "paradigm-ehb/agent/gen/resources/v1"
	res "paradigm-ehb/agent/internal/resources"
	wr "paradigm-ehb/agent/pkg/cgowrap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResourcesService struct {
	proto.UnimplementedResourcesServiceServer
}

func (s *ResourcesService) GetSystemResources(
	ctx context.Context,
	req *proto.GetSystemResourcesRequest,
) (*proto.GetSystemResourcesResponse, error) {

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled")
	default:
	}

	snap, err := res.GetCompleteSystemResources()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to collect system resources: %v",
			err,
		)
	}

	return &proto.GetSystemResourcesResponse{
		Resources: mapSystemResources(snap),
	}, nil
}

func mapSystemResources(s *res.SystemResources) *proto.SystemResources {
	return &proto.SystemResources{
		Cpu:       mapCPU(s.CPU),
		Memory:    mapMemory(s.Memory),
		Device:    mapDevice(s.Device),
		Disks:     mapDisks(s.Disks),
		Processes: mapProcesses(s.Procs),
	}
}

func mapCPU(c wr.Cpu) *proto.Cpu {
	return &proto.Cpu{
		Vendor:    c.Vendor,
		Model:     c.Model,
		Frequency: c.Frequency,
		MaxCore:   c.MaxCore,
	}
}

func mapMemory(m wr.Ram) *proto.Memory {
	return &proto.Memory{
		Total: m.Total,
		Free:  m.Free,
	}
}

func mapDevice(d wr.Device) *proto.Device {
	return &proto.Device{
		OsVersion: d.OsVersion,
		Uptime:    d.Uptime,
	}
}

func mapDisks(disks []wr.Disk) []*proto.Disk {
	out := make([]*proto.Disk, 0, len(disks))

	for _, d := range disks {
		parts := make([]*proto.DiskPartition, 0, len(d.Partitions))
		for _, p := range d.Partitions {
			parts = append(parts, &proto.DiskPartition{
				Name:   p.Name,
				Major:  p.Major,
				Minor:  p.Minor,
				Blocks: p.Blocks,
			})
		}

		out = append(out, &proto.Disk{
			Partitions: parts,
		})
	}

	return out
}

func mapProcesses(ps []wr.Process) []*proto.Process {
	out := make([]*proto.Process, 0, len(ps))

	for _, p := range ps {
		out = append(out, &proto.Process{
			Pid:        p.PID,
			Name:       p.Name,
			State:      mapProcessState(p.State),
			Utime:      p.UTime,
			NumThreads: p.NumThreads,
		})
	}

	return out
}

func mapProcessState(s wr.ProcessState) proto.ProcessState {
	switch s {
	case wr.ProcessRunning:
		return proto.ProcessState_PROCESS_STATE_RUNNING
	case wr.ProcessSleeping:
		return proto.ProcessState_PROCESS_STATE_SLEEPING
	case wr.ProcessDiskSleep:
		return proto.ProcessState_PROCESS_STATE_DISK_SLEEPING
	case wr.ProcessStopped:
		return proto.ProcessState_PROCESS_STATE_STOPPED
	case wr.ProcessTracingStopped:
		return proto.ProcessState_PROCESS_STATE_TRACING_STOPPED
	case wr.ProcessZombie:
		return proto.ProcessState_PROCESS_STATE_ZOMBIE
	case wr.ProcessDead:
		return proto.ProcessState_PROCESS_STATE_DEAD
	default:
		return proto.ProcessState_PROCESS_STATE_UNSPECIFIED
	}
}
