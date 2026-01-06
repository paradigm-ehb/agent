package grpc_handler

import (
	"context"
	"log"

	actions "paradigm-ehb/agent/gen/actions/v1"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	da "paradigm-ehb/agent/internal/deviceactions"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceActionsService struct {
	actions.UnimplementedActionServiceServer
}

func (s *DeviceActionsService) Action(ctx context.Context, req *actions.ActionRequest) (*actions.ActionReply, error) {

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


	obj := bus.Object(
		"org.freedesktop.login1",
		"/org/freedesktop/login1",
	)

	switch req.GetDeviceAction() {
	case actions.DeviceAction_DEVICE_ACTION_SHUTDOWN:
		out = string(da.PerformDeviceAction(obj, da.DeviceActionShutdown))
	case actions.DeviceAction_DEVICE_ACTION_REBOOT:
		out = string(da.PerformDeviceAction(obj, da.DeviceActionReboot))
	case actions.DeviceAction_DEVICE_ACTION_SUSPEND:
		out = string(da.PerformDeviceAction(obj, da.DeviceActionSuspend))
	case actions.DeviceAction_DEVICE_ACTION_HIBERNATE:
		out = string(da.PerformDeviceAction(obj, da.DeviceActionHibernate))
	default:
		out = "unknown action"
	}
	return &actions.ActionReply{
		Status: out,
	}, nil
}
