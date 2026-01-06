package grpc_handler

import (
	"context"
	"log"

	manager "paradigm-ehb/agent/internal/dbusservices"

	"paradigm-ehb/agent/gen/services/v1"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	"paradigm-ehb/agent/internal/dbusservices/systemd"

	"github.com/godbus/dbus"
)

type HandlerService struct {
	v1.UnimplementedHandlerServiceServer
}

func (s *HandlerService) UnitAction(_ context.Context, in *v1.ServiceActionRequest) (*v1.ServiceActionReply, error) {

	var out string
	conn, err := dh.CreateSystemBus()
	if err != nil {

		log.Println("failed to create systembus")
		out = "failed system bus"
	}

	defer func(conn *dbus.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	if in.GetUnitFileAction() == v1.ServiceActionRequest_UNIT_FILE_ACTION_ENABLE {

		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionEnable, true, true, []string{in.ServiceName})
		if err != nil {

			log.Println("failed to enable service")
			out = "failed unit file action"
		}

	} else if in.GetUnitFileAction() == v1.ServiceActionRequest_UNIT_FILE_ACTION_DISABLE {

		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionDisable, true, true, []string{in.ServiceName})
		if err != nil {

			log.Println("failed to disable service")
			out = "failed unit file action"
		}

	} else {

		log.Println("Bad input")
		out = "Bad input"

	}

	if in.GetUnitAction() == v1.ServiceActionRequest_UNIT_ACTION_START {

		err = manager.RunAction(conn, servicecontrol.UnitActionStart, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else if in.GetUnitAction() == v1.ServiceActionRequest_UNIT_ACTION_STOP {

		err = manager.RunAction(conn, servicecontrol.UnitActionStop, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else if in.GetUnitAction() == v1.ServiceActionRequest_UNIT_ACTION_RESTART {

		err = manager.RunAction(conn, servicecontrol.UnitActionRestart, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else {

		log.Println("external bad input")
		out = "external bad input"
	}

	err = manager.RunRetrieval(conn, true)
	if err != nil {
		log.Println("failed to do everything")
		out = "failed even more"
	}

	return &v1.ServiceActionReply{Status: out}, nil
}
