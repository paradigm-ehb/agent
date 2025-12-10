package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func Run(a svc.Action, s svc.SymlinkAction, service string) error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	ch := make(chan []svctypes.Ass)
	parse := make(chan []svctypes.Ass)

	go svc.GetAllUnits(obj, ch)
	// svc.GetAllUnits(obj, nil)
	go dh.Parse(ch, parse)

	// fmt.Println(<-parse)
	<-parse

	return nil
}
