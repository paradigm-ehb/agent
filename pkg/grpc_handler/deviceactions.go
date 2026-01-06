package grpc_handler

import (
	"context"
	"log"

	"paradigm-ehb/agent/gen/device_actions/v1"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	da "paradigm-ehb/agent/internal/deviceactions"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceActionsService struct {
	device_actions.UnimplementedDeviceActionsServiceServer
}

func (s *DeviceActionsService) Action(ctx context.Context, req *device_actions.DeviceActionRequest) (*device_actions.DeviceActionReply, error) {

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled")
	default:
	}

	var out string
	bus, err := dh.CreateSystemBus()
	if err != nil {
		log.Println("failed to create systembus")
		out = "failed system bus"
	}

	switch req.DeviceActionRequest {
	case device_actions.DeviceActionRequest_DEVICE_ACTION_REQUEST_SHUTDOWN:
		out = string(da.PerformDeviceAction(bus.BusObject(), da.DeviceActionShutdown))
	case device_actions.DeviceActionRequest_DEVICE_ACTION_REQUEST_REBOOT:
		out = string(da.PerformDeviceAction(bus.BusObject(), da.DeviceActionReboot))
	case device_actions.DeviceActionRequest_DEVICE_ACTION_REQUEST_SUSPEND:
		out = string(da.PerformDeviceAction(bus.BusObject(), da.DeviceActionSuspend))
	case device_actions.DeviceActionRequest_DEVICE_ACTION_REQUEST_HIBERNATE:
		out = string(da.PerformDeviceAction(bus.BusObject(), da.DeviceActionHibernate))
	default:
		out = "unknown action"
	}
	return &device_actions.DeviceActionReply{
		Result: out,
	}, nil
}
