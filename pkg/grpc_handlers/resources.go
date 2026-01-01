package service

import (
	"context"

	proto "paradigm-ehb/agent/gen/resources/v1"
	res "paradigm-ehb/agent/internal/resources"
	wr "paradigm-ehb/agent/pkg/cgowrap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * ResourcesService implements the gRPC ResourcesServiceServer.
 *
 * It exposes endpoints that provide a snapshot of system-level resources
 * such as CPU, memory, disks, and running processes. All data collection
 * is delegated to the internal resources package and then mapped to
 * protobuf-defined response types.
 */
type ResourcesService struct {
	proto.UnimplementedResourcesServiceServer
}

/**
 * GetSystemResources returns a full snapshot of system resources.
 *
 * The method:
 *   - Respects context cancellation to avoid unnecessary work
 *   - Collects system information via the internal resources layer
 *   - Maps internal domain structures to protobuf response messages
 *
 * Errors:
 *   - Returns codes.Canceled if the request context is canceled
 *   - Returns codes.Internal if resource collection fails
 */
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

/**
 * mapSystemResources converts an internal SystemResources snapshot
 * into its protobuf representation.
 *
 * This acts as the top-level aggregation mapper, delegating
 * to more specific mapping functions per subsystem.
 */
func mapSystemResources(s *res.SystemResources) *proto.SystemResources {
	return &proto.SystemResources{
		Cpu:       mapCPU(s.CPU),
		Memory:    mapMemory(s.Memory),
		Device:    mapDevice(s.Device),
		Disks:     mapDisks(s.Disks),
		Processes: mapProcesses(s.Procs),
	}
}

/**
 * mapCPU maps CPU metadata from the cgo wrapper type
 * into the protobuf Cpu message.
 */
func mapCPU(c wr.Cpu) *proto.Cpu {
	return &proto.Cpu{
		Vendor:    c.Vendor,
		Model:     c.Model,
		Frequency: c.Frequency,
		MaxCore:   c.MaxCore,
	}
}

/**
 * mapMemory maps RAM usage information into the protobuf Memory message.
 *
 * Values are expected to be raw byte counts as reported by the system.
 */
func mapMemory(m wr.Ram) *proto.Memory {
	return &proto.Memory{
		Total: m.Total,
		Free:  m.Free,
	}
}

/**
 * mapDevice maps general device and OS-level metadata
 * into the protobuf Device message.
 */
func mapDevice(d wr.Device) *proto.Device {
	return &proto.Device{
		OsVersion: d.OsVersion,
		Uptime:    d.Uptime,
	}
}

/**
 * mapDisks maps a slice of disk descriptors into protobuf Disk messages.
 *
 * Each disk contains a list of partitions, which are also converted
 * field-by-field into their protobuf equivalents.
 */
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

/**
 * mapProcesses maps a slice of process descriptors into protobuf Process messages.
 *
 * Each process includes basic scheduling and accounting information
 * such as PID, state, CPU time, and thread count.
 */
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

/**
 * mapProcessState converts an internal ProcessState enum
 * into the corresponding protobuf ProcessState value.
 *
 * Unknown or unmapped states are converted to PROCESS_STATE_UNSPECIFIED
 * to preserve forward compatibility.
 */
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
