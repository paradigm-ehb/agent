package grpc_handler

import (
	"context"
	"fmt"

	v3 "paradigm-ehb/agent/gen/services/v3"

	manager "paradigm-ehb/agent/internal/dbusservices"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	servicecontrol "paradigm-ehb/agent/internal/dbusservices/systemd"
	types "paradigm-ehb/agent/internal/dbusservices/types"
)

type HandlerServicev3 struct {
	v3.UnimplementedHandlerServiceServer
}


func mapLoadedUnit(u *types.LoadedUnit) *v3.LoadedUnit {
	if u == nil {
		return nil
	}

	return &v3.LoadedUnit{
		Name:         u.Name,
		Description:  u.Description,
		LoadState:    u.LoadState,
		SubState:     u.SubState,
		ActiveState:  u.ActiveState,
		DepUnit:      u.DepUnit,
		ObjectPath:   string(u.ObjectPath),
		QueuedJob:    u.QueudJob,
		JobType:      u.JobType,
		JobPath:      string(u.JobPath),
	}
}

func mapLoadedUnits(units []*types.LoadedUnit) []*v3.LoadedUnit {
	out := make([]*v3.LoadedUnit, 0, len(units))
	for _, u := range units {
		if pu := mapLoadedUnit(u); pu != nil {
			out = append(out, pu)
		}
	}
	return out
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

	if err := manager.RunAction(conn, action, in.UnitName); err != nil {
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

	if err := manager.RunSymlinkAction(
		conn,
		action,
		in.Runtime,
		in.Force,
		[]string{in.UnitName},
	); err != nil {
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

	units, err := manager.RunRetrieval(conn, true)
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.GetUnitsReply{
		Units:   mapLoadedUnits(units),
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

	units, err := manager.RunRetrieval(conn, false)
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.GetUnitsReply{
		Units:   mapLoadedUnits(units),
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

	var filters []string

	switch in.State {
	case v3.GetUnitsFilteredRequest_LOADED:
		filters = []string{"loaded"}
	case v3.GetUnitsFilteredRequest_NOT_FOUND:
		filters = []string{"not-found"}
	case v3.GetUnitsFilteredRequest_BAD_SETTING:
		filters = []string{"bad-setting"}
	case v3.GetUnitsFilteredRequest_ERROR:
		filters = []string{"error"}
	case v3.GetUnitsFilteredRequest_MASKED:
		filters = []string{"masked"}
	default:
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: "unspecified filter state",
		}, nil
	}

	units, err := manager.MapFilteredUnits(conn, filters)
	if err != nil {
		return &v3.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v3.GetUnitsReply{
		Units:   mapLoadedUnits(units),
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

	obj, err := dh.CreateSystemdObject(conn)
	if err != nil {
		return &v3.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, fmt.Errorf("failed to create systemd object")
	}

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
```
