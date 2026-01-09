package grpc_handler

import (
	"context"
	"fmt"

	v3 "paradigm-ehb/agent/gen/services/v3"
	manager "paradigm-ehb/agent/internal/dbusservices"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	servicecontrol "paradigm-ehb/agent/internal/dbusservices/systemd"
)

type HandlerServicev3 struct {
	v3.UnimplementedHandlerServiceServer
}

func (s *HandlerServicev3) PerformUnitAction(
	_ context.Context,
	in *v3.UnitActionRequest,
) (*v3.UnitActionReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v3.UnitActionReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	var action servicecontrol.UnitAction
	var actionName string

	switch in.Action {
	case v3.UnitActionRequest_UNIT_ACTION_START:
		action = servicecontrol.UnitActionStart
		actionName = "start"
	case v3.UnitActionRequest_UNIT_ACTION_STOP:
		action = servicecontrol.UnitActionStop
		actionName = "stop"
	case v3.UnitActionRequest_UNIT_ACTION_RESTART:
		action = servicecontrol.UnitActionRestart
		actionName = "restart"
	default:
		return &v3.UnitActionReply{
			Success:      false,
			ErrorMessage: "unspecified unit action",
		}, nil
	}

	err = manager.RunAction(conn, action, in.UnitName)
	if err != nil {
		return &v3.UnitActionReply{
			Status:       []byte(fmt.Sprintf("failed to %s unit", actionName)),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.UnitActionReply{
		Status:  []byte("success"),
		Success: true,
	}, nil
}

func (s *HandlerServicev3) PerformUnitFileAction(
	_ context.Context,
	in *v3.UnitFileActionRequest,
) (*v3.UnitFileActionReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v3.UnitFileActionReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	var action servicecontrol.UnitFileAction

	switch in.Action {
	case v3.UnitFileActionRequest_UNIT_FILE_ACTION_ENABLE:
		action = servicecontrol.UnitFileActionEnable
	case v3.UnitFileActionRequest_UNIT_FILE_ACTION_DISABLE:
		action = servicecontrol.UnitFileActionDisable
	default:
		return &v3.UnitFileActionReply{
			Success:      false,
			ErrorMessage: "unspecified unit file action",
		}, nil
	}

	err = manager.RunSymlinkAction(
		conn,
		action,
		in.Runtime,
		in.Force,
		[]string{in.UnitName},
	)

	if err != nil {
		return &v3.UnitFileActionReply{
			Status:       []byte("unit file action failed"),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.UnitFileActionReply{
		Status:  []byte("success"),
		Success: true,
	}, nil
}

func (s *HandlerServicev3) GetAllUnits(
	_ context.Context,
	_ *v3.GetUnitsRequest,
) (*v3.GetUnitsReply, error) {

	conn, err := dh.CreateSystemBus()

	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	/**
		
	Handle error hanlding
	*/
	_, err = manager.RunRetrieval(conn, true)
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// mappedUnits := make(*v3.GetUnitsReply, 0, len(units))
	// mappedUnits = &units{}

	return &v3.GetUnitsReply{
		Units:   nil,
		Success: true,
	}, nil
}

func (s *HandlerServicev3) GetLoadedUnits(
	_ context.Context,
	_ *v3.GetUnitsRequest,
) (*v3.GetUnitsReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	/**
		TODO(nasr): handle error
	*/
	_, err = manager.RunRetrieval(conn, false)
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.GetUnitsReply{
		Units:   nil,
		Success: true,
	}, nil
}

func (s *HandlerServicev3) GetUnitStatus(
	_ context.Context,
	in *v3.GetUnitStatusRequest,
) (*v3.GetUnitStatusReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v3.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	/**
		TODO(nasr): todo handle error
	*/
	obj, _ := dh.CreateSystemdObject(conn)

	state, err := manager.GetStatus(obj, in.UnitName)
	if err != nil {
		return &v3.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.GetUnitStatusReply{
		State:   state,
		Success: true,
	}, nil
}

func (s *HandlerServicev3) GetFilteredUnits(
	_ context.Context,
	in *v3.GetUnitsFilteredRequest,
) (*v3.GetUnitsReply, error) {

	conn, err := dh.CreateSystemBus()

	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	/**
		
		TODO(nasr): handle the correct types
		doing now

	*/
	_, err = manager.MapFilteredUnits(conn)
	if err != nil {

		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, fmt.Errorf("failed to map filtered units %v: ", err)
	}

	return &v3.GetUnitsReply{}, nil
}
