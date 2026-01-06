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

func (s *HandlerServiceV2) PerformAction(
	_ context.Context,
	in *v2.ServiceActionRequest) (*v2.ServiceActionReply, error) {
	conn, err := dh.CreateSystemBus()
	if err != nil {

		fmt.Println("failed to create system bus")
	}

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
			{
				action = servicecontrol.UnitFileActionEnable
			}

		case v2.ServiceActionRequest_UNIT_FILE_ACTION_DISABLE:
			{
				action = servicecontrol.UnitFileActionDisable
			}
		}

		err = manager.RunSymlinkAction(conn, action, runtime, force, []string{in.ServiceName})

		if err != nil {
			return &v2.ServiceActionReply{
				Status:       []byte(fmt.Sprintf("NICE IT WORKED (nope joking it failed, tough luck, try again another time)")),
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
			{
				action = servicecontrol.UnitActionStart
			}
		case v2.ServiceActionRequest_UNIT_ACTION_STOP:
			{
				action = servicecontrol.UnitActionStop
			}
		case v2.ServiceActionRequest_UNIT_ACTION_RESTART:
			{
				action = servicecontrol.UnitActionRestart
				actionName = "restart"
			}
		}

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

	return &v2.ServiceActionReply{
		Status:       []byte("succes"),
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

	_, al, err := manager.RunRetrieval(conn, true)

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
		fmt.Println("failed to create system bus")
	}

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
	in *v2.GetUnitStatusRequest) (*v2.GetUnitStatusReply, error) {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		fmt.Println("failed to crete sysembus")
	}
	obj := dh.CreateSystemdObject(conn)

	status, err := manager.GetStatus(obj, in.UnitName)
	if err != nil {
		fmt.Printf("\n %s\n\n", err)
		return &v2.GetUnitStatusReply{
			State:        "failed here?",
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
