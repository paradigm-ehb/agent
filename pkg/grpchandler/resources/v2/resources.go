package grpc_handler

import (
	"context"

	proto "paradigm-ehb/agent/gen/resources/v2"
	"paradigm-ehb/agent/internal/resources"
	cgo "paradigm-ehb/agent/pkg/cgowrap"

	"golang.org/x/sys/unix"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * ResourcesServiceV2 implements the gRPC ResourcesServiceServer.
 *
 * It exposes endpoints that provide a snapshot of system-level resources
 * such as CPU, memory, disks, and running processes. All data collection
 * is delegated to the internal resources package and then mapped to
 * protobuf-defined response types.
 */
type ResourcesServiceV2 struct {
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
func (s *ResourcesServiceV2) GetSystemResources(
	ctx context.Context,
	req *proto.GetSystemResourcesRequest,
) (*proto.GetSystemResourcesResponse, error) {

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled")
	default:
	}

	snap, err := resources.GetCompleteSystemResources()
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
* TODO(nasr): update individual resources
 */

func (s *ResourcesServiceV2) ProcessAction(
	ctx context.Context,
	req *proto.ProcessActionRequest,
) (*proto.ProcessActionReply, error) {

	err := cgo.ProcessAction(int(req.Pid), unix.Signal(req.Signal))
	if err != nil {

		return &proto.ProcessActionReply{
			Succes: false,
		}, nil
	}

	return &proto.ProcessActionReply{
		Succes: true,
	}, nil
}

/**
 * mapSystemResources converts an internal SystemResources snapshot
 * into its protobuf representation.
 *
 * This acts as the top-level aggregation mapper, delegating
 * to more specific mapping functions per subsystem.
 */
func mapSystemResources(s *resources.SystemResources) *proto.SystemResources {

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
func mapCPU(c cgo.Cpu) *proto.Cpu {

	return &proto.Cpu{
		Vendor:    c.Vendor,
		Model:     c.Model,
		Frequency: c.Frequency,
		MaxCore:   c.MaxCore,
		TotalTime: c.TotalTime,
		IdleTime:  c.IdleTime,
	}

}

/**
 * mapMemory maps RAM usage information into the protobuf Memory message.
 *
 * Values are expected to be raw byte counts as reported by the system.
 */
func mapMemory(m cgo.Ram) *proto.Memory {

	return &proto.Memory{
		Total: m.Total,
		Free:  m.Free,
	}

}

/**
 * mapDevice maps general device and OS-level metadata
 * into the protobuf Device message.
 */
func mapDevice(d cgo.Device) *proto.Device {

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
func mapDisks(disks []cgo.Disk) []*proto.Disk {

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
 * mapProcessState converts an internal ProcessState enum
 * into the corresponding protobuf ProcessState value.
 *
 * Unknown or unmapped states are converted to PROCESS_STATE_UNSPECIFIED
 * to preserve forward compatibility.
 */
func mapProcessState(s cgo.ProcessState) proto.ProcessState {

	switch s {

	case cgo.ProcessRunning:
		{

			return proto.ProcessState_PROCESS_STATE_RUNNING
		}
	case cgo.ProcessSleeping:
		{

			return proto.ProcessState_PROCESS_STATE_SLEEPING
		}
	case cgo.ProcessDiskSleep:
		{

			return proto.ProcessState_PROCESS_STATE_DISK_SLEEPING
		}
	case cgo.ProcessStopped:
		{

			return proto.ProcessState_PROCESS_STATE_STOPPED
		}
	case cgo.ProcessTracingStopped:
		{

			return proto.ProcessState_PROCESS_STATE_TRACING_STOPPED
		}
	case cgo.ProcessZombie:
		{

			return proto.ProcessState_PROCESS_STATE_ZOMBIE
		}
	case cgo.ProcessDead:
		{

			return proto.ProcessState_PROCESS_STATE_DEAD
		}
	case cgo.ProcessIdle:
		{
			return proto.ProcessState_PROCESS_STATE_IDLE
		}

	case cgo.ProcessUndefined:
		{
			return proto.ProcessState_PROCESS_STATE_UNSPECIFIED

		}
	default:

		return proto.ProcessState_PROCESS_STATE_UNSPECIFIED
	}
}

/**
 * mapProcesses maps a slice of process descriptors into protobuf Process messages.
 *
 * Each process includes basic scheduling and accounting information
 * such as PID, state, CPU time, and thread count.
 */
func mapProcesses(ps []cgo.Process) []*proto.Process {

	out := make([]*proto.Process, 0, len(ps))

	for _, p := range ps {
		out = append(out, &proto.Process{
			Pid:        int32(p.PID),
			Name:       p.Name,
			State:      mapProcessState(p.State),
			Utime:      p.UTime,
			NumThreads: p.NumThreads,
		})
	}

	return out
}
