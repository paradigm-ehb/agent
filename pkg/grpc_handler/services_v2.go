package grpc_handler

import (
	"context"
	"log"
	manager "paradigm-ehb/agent/internal/dbusservices"
	"paradigm-ehb/agent/gen/services/v2"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	"paradigm-ehb/agent/internal/dbusservices/systemd"
	"github.com/godbus/dbus"
)

type HandlerServiceV2 struct {
	v2.UnimplementedHandlerServiceServer
}

func (s *HandlerService) UnitActionV2(_ context.Context, in *v2.ServiceActionRequest) (*v2.ServiceActionReply, error) {
	var out []byte
	conn, err := dh.CreateSystemBus()
	if err != nil {
		log.Println("failed to create systembus")
		out = []byte("failed system bus")
	}
	defer func(conn *dbus.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	if in.GetUnitFileAction() == v2.ServiceActionRequest_UNIT_FILE_ACTION_ENABLE {
		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionEnable, true, true, []string{in.ServiceName})
		if err != nil {
			log.Println("failed to enable service")
			out = []byte("failed unit file action")
		}
	} else if in.GetUnitFileAction() == v2.ServiceActionRequest_UNIT_FILE_ACTION_DISABLE {
		err = manager.RunSymlinkAction(conn, servicecontrol.UnitFileActionDisable, true, true, []string{in.ServiceName})
		if err != nil {
			log.Println("failed to disable service")
			out = []byte("failed unit file action")
		}
	} else {
		log.Println("failed to do everything")
		out = []byte("failed alot")
	}

	if in.GetUnitAction() == v2.ServiceActionRequest_UNIT_ACTION_START {
		err = manager.RunAction(conn, servicecontrol.UnitActionStart, in.ServiceName)
		if err != nil {
			log.Println("failed to run action")
			out = []byte("failed unit action")
		}
	} else if in.GetUnitAction() == v2.ServiceActionRequest_UNIT_ACTION_STOP {
		err = manager.RunAction(conn, servicecontrol.UnitActionStop, in.ServiceName)
		if err != nil {
			log.Println("failed to run action")
			out = []byte("failed unit action")
		}
	} else if in.GetUnitAction() == v2.ServiceActionRequest_UNIT_ACTION_RESTART {
		err = manager.RunAction(conn, servicecontrol.UnitActionRestart, in.ServiceName)
		if err != nil {
			log.Println("failed to run action")
			out = []byte("failed unit action")
		}
	} else {
		log.Println("failed to to do everything")
		out = []byte("failed a lot")
	}

	err = manager.RunRetrieval(conn, true)
	if err != nil {
		log.Println("failed to do everything")
		out = []byte("failed even more")
	}

	return &v2.ServiceActionReply{Status: out}, nil
}
