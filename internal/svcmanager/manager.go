package svcmanager

import (
	"fmt"

	// "github.com/godbus/dbus"
	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func RunAction(ac svc.Action, sc svc.SymlinkAction, service string) error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	// obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, 0)
	call := obj.Call(string(ac), 0, service, "replace")

	if call.Err != nil {
		fmt.Println(call.Err)
	}

	return nil

}

// @param, true for all on disk, false for loaded units
func RunRetrieval(all bool) error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	if all {

		ch := make(chan []svctypes.Ass)
		parse := make(chan []svctypes.Ass)

		go svc.GetAllUnits(obj, ch)
		go dh.ParseAllUnits(ch, parse)
		<-parse

	} else if !all {

		ch := make(chan []svctypes.Assssssouso)
		parse := make(chan []svctypes.Assssssouso)

		go svc.GetLoadedUnits(obj, ch)
		go dh.ParseLoadedUnits(ch, parse)
		<-parse

	} else {
		return fmt.Errorf("failed parameter")
	}

	return nil
}
