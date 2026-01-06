package grpc_handler

import (
	"context"
	"fmt"
	"log"

	"github.com/godbus/dbus"
	v2 "paradigm-ehb/agent/gen/services/v2"
	manager "paradigm-ehb/agent/internal/dbusservices"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	servicecontrol "paradigm-ehb/agent/internal/dbusservices/systemd"
)

type HandlerServiceV2 struct {
	v2.UnimplementedHandlerServiceServer
}

var SystemBus *dbus.Conn
var SharedBus *dbus.Conn

func ( *HandlerServiceV2) PerformAction(_ context.Context, in *v2.ServiceActionRequest) (*v2.ServiceActionReply, error) {
	conn, err := dh.CreateSystemBus()
	if err != nil {
		log.Printf("failed to create systembus: %v", err)
		return &v2.ServiceActionReply{
			Status:       []byte("failed to create system bus"),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	defer func(conn *dbus.Conn) {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn)

	var operations []string

	if in.GetUnitFileAction() != v2.ServiceActionRequest_UNIT_FILE_ACTION_UNSPECIFIED {
		runtime := true
		force := true
		if in.Runtime != nil {
			runtime = *in.Runtime
		}
		if in.Force != nil {
			force = *in.Force
		}

		var action servicecontrol.UnitFileAction
		switch in.GetUnitFileAction() {
		case v2.ServiceActionRequest_UNIT_FILE_ACTION_ENABLE:
			action = servicecontrol.UnitFileActionEnable
			operations = append(operations, "enable")
		case v2.ServiceActionRequest_UNIT_FILE_ACTION_DISABLE:
			action = servicecontrol.UnitFileActionDisable
			operations = append(operations, "disable")
		}

		err = manager.RunSymlinkAction(conn, action, runtime, force, []string{in.ServiceName})
		if err != nil {
			log.Printf("failed to %s service: %v", operations[len(operations)-1], err)
			return &v2.ServiceActionReply{
				Status:       []byte(fmt.Sprintf("failed to %s service", operations[len(operations)-1])),
				Success:      false,
				ErrorMessage: err.Error(),
			}, nil
		}
	}

	if in.GetUnitAction() != v2.ServiceActionRequest_UNIT_ACTION_UNSPECIFIED {
		var action servicecontrol.UnitAction
		var actionName string

		switch in.GetUnitAction() {
		case v2.ServiceActionRequest_UNIT_ACTION_START:
			action = servicecontrol.UnitActionStart
			actionName = "start"
		case v2.ServiceActionRequest_UNIT_ACTION_STOP:
			action = servicecontrol.UnitActionStop
			actionName = "stop"
		case v2.ServiceActionRequest_UNIT_ACTION_RESTART:
			action = servicecontrol.UnitActionRestart
			actionName = "restart"
		}

		operations = append(operations, actionName)
		err = manager.RunAction(conn, action, in.ServiceName)
		if err != nil {
			log.Printf("failed to %s service: %v", actionName, err)
			return &v2.ServiceActionReply{
				Status:       []byte(fmt.Sprintf("failed to %s service", actionName)),
				Success:      false,
				ErrorMessage: err.Error(),
			}, nil
		}
	}

	var statusMsg string
	if len(operations) > 0 {
		statusMsg = fmt.Sprintf("successfully executed: %v on %s", operations, in.ServiceName)
	} else {
		statusMsg = "no operations specified"
	}

	return &v2.ServiceActionReply{
		Status:       []byte(statusMsg),
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (s *HandlerServiceV2) GetAllUnits(_ context.Context, _ *v2.GetUnitsRequest) (*v2.GetUnitsReply, error) {
	conn, err := dh.CreateSessionBus()
	if err != nil {
		log.Printf("failed to create systembus: %v", err)
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	defer conn.Close()

	err = manager.RunRetrieval(conn, true)
	if err != nil {
		log.Printf("failed to retrieve all units: %v", err)
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitsReply{
		UnitsData:    []byte("all units retrieved successfully"),
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (s *HandlerServiceV2) GetLoadedUnits(_ context.Context, _ *v2.GetUnitsRequest) (*v2.GetUnitsReply, error) {
	conn, err := dh.CreateSystemBus()
	if err != nil {
		log.Printf("failed to create systembus: %v", err)
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	defer conn.Close()

	err = manager.RunRetrieval(conn, false)
	if err != nil {
		log.Printf("failed to retrieve loaded units: %v", err)
		return &v2.GetUnitsReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitsReply{
		UnitsData:    []byte("loaded units retrieved successfully"),
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (s *HandlerServiceV2) GetUnitStatus(
	_ context.Context,
	in *v2.GetUnitStatusRequest,) (*v2.GetUnitStatusReply, error) {

	conn, err := dh.CreateSessionBus()

	if err != nil {
		log.Printf("failed to create system bus: %v", err)
		return &v2.GetUnitStatusReply{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	status, err := manager.GetStatus(obj, in.UnitName)
	if err != nil {
		return &v2.GetUnitStatusReply{
			State:        "",
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &v2.GetUnitStatusReply{
		State:        status,
		Success:      true,
		ErrorMessage: "",
	}, nil
}


