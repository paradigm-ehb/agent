package grpc_handler

import (
	"context"
	"fmt"

	v2 "paradigm-ehb/agent/gen/services/v2"
	manager "paradigm-ehb/agent/internal/dbusservices"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	servicecontrol "paradigm-ehb/agent/internal/dbusservices/systemd"

)

type HandlerServiceV2 struct {
	v2.UnimplementedHandlerServiceServer
}


func (s *HandlerServiceV2) PerformUnitAction(
	_ context.Context,
	in *v2.UnitActionRequest,
) (*v2.UnitActionReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v2.UnitActionReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	var action servicecontrol.UnitAction
	var actionName string

	switch in.Action {
	case v2.UnitActionRequest_UNIT_ACTION_START:
		action = servicecontrol.UnitActionStart
		actionName = "start"
	case v2.UnitActionRequest_UNIT_ACTION_STOP:
		action = servicecontrol.UnitActionStop
		actionName = "stop"
	case v2.UnitActionRequest_UNIT_ACTION_RESTART:
		action = servicecontrol.UnitActionRestart
		actionName = "restart"
	default:
		return &v2.UnitActionReply{
			Success:      false,
			ErrorMessage: "unspecified unit action",
		}, nil
	}

	err = manager.RunAction(conn, action, in.UnitName)
	if err != nil {
		return &v2.UnitActionReply{
			Status:       []byte(fmt.Sprintf("failed to %s unit", actionName)),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.UnitActionReply{
		Status:  []byte("success"),
		Success: true,
	}, nil
}


func (s *HandlerServiceV2) PerformUnitFileAction(
	_ context.Context,
	in *v2.UnitFileActionRequest,
) (*v2.UnitFileActionReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v2.UnitFileActionReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	var action servicecontrol.UnitFileAction

	switch in.Action {
	case v2.UnitFileActionRequest_UNIT_FILE_ACTION_ENABLE:
		action = servicecontrol.UnitFileActionEnable
	case v2.UnitFileActionRequest_UNIT_FILE_ACTION_DISABLE:
		action = servicecontrol.UnitFileActionDisable
	default:
		return &v2.UnitFileActionReply{
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
		return &v2.UnitFileActionReply{
			Status:       []byte("unit file action failed"),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.UnitFileActionReply{
		Status:  []byte("success"),
		Success: true,
	}, nil
}


func (s *HandlerServiceV2) GetAllUnits(
	_ context.Context,
	_ *v2.GetUnitsRequest,
) (*v2.GetUnitsReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	/**
		
		TODO(nasr): fix the design issue and return the correct mapped types
		from the correct namespace

	*/
	_, err = manager.RunRetrievalDeprecated(conn, true)
	if err != nil {
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitsReply{
		Units:   nil,
		Success: true,
	}, nil
}

func (s *HandlerServiceV2) GetLoadedUnits(
	_ context.Context,
	_ *v2.GetUnitsRequest,
) (*v2.GetUnitsReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	/**
		
		TODO(nasr): fix the design issue and return the correct mapped types
		from the correct namespace

	*/
	_, err = manager.RunRetrievalDeprecated(conn, false)
	if err != nil {
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitsReply{
		Units:   nil,
		Success: true,
	}, nil
}


func (s *HandlerServiceV2) GetUnitStatus(
	_ context.Context,
	in *v2.GetUnitStatusRequest,
) (*v2.GetUnitStatusReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return &v2.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	obj, err := dh.CreateSystemdObject(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create system object")
	}

	state, err := manager.UnitStatus(obj, in.UnitName)
	if err != nil {
		return &v2.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitStatusReply{
		State:   state,
		Success: true,
	}, nil
}
