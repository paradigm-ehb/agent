package service

import (
	"context"
	"log"

	manager "paradigm-ehb/agent/internal/dbus_services"

	"paradigm-ehb/agent/gen/services/v1"
	dh "paradigm-ehb/agent/internal/dbus_services/dbus"
	"paradigm-ehb/agent/internal/dbus_services/systemd"

	"github.com/godbus/dbus"
)

type HandlerService struct {
	services.UnimplementedHandlerServiceServer
}

func (s *HandlerService) Action(_ context.Context, in *services.ServiceActionRequest) (*services.ServiceActionReply, error) {

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

	if in.GetUnitFileAction() == services.ServiceActionRequest_UNIT_FILE_ACTION_ENABLE {

		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionEnable, true, true, []string{in.ServiceName})
		if err != nil {

			log.Println("failed to enable service")
			out = "failed unit file action"
		}

	} else if in.GetUnitFileAction() == services.ServiceActionRequest_UNIT_FILE_ACTION_DISABLE {

		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionDisable, true, true, []string{in.ServiceName})
		if err != nil {

			log.Println("failed to disable service")
			out = "failed unit file action"
		}

	} else {

		log.Println("failed to do everything")
		out = "failed alot"

	}

	if in.GetUnitAction() == services.ServiceActionRequest_UNIT_ACTION_START {

		err = manager.RunAction(conn, servicecontrol.UnitActionStart, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else if in.GetUnitAction() == services.ServiceActionRequest_UNIT_ACTION_STOP {

		err = manager.RunAction(conn, servicecontrol.UnitActionStop, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else if in.GetUnitAction() == services.ServiceActionRequest_UNIT_ACTION_RESTART {

		err = manager.RunAction(conn, servicecontrol.UnitActionRestart, in.ServiceName)
		if err != nil {

			log.Println("failed to run action")
			out = "failed unit action"
		}

	} else {

		log.Println("failed to to do everything")
		out = "failed a lot"
	}

	err = manager.RunRetrieval(conn, true)
	if err != nil {
		log.Println("failed to do everything")
		out = "failed even more"
	}

	return &services.ServiceActionReply{Status: out}, nil
}
