package svcmanager

import (
	"fmt"

	"github.com/godbus/dbus"
	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func RunAction(ac svc.Action, service string) error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	call := obj.Call(string(ac), 0, service, "replace")

	if call.Err != nil {
		fmt.Println("error on action, ", call.Err)
	}

	return nil
}

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunSymlinkAction(sc svc.SymlinkAction, enableForRunTime bool, enableForce bool, service []string) error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	/* EnableUnitFiles(in  as files,
	*                 in  b runtime,
	*                 in  b force,
	*                 out b carries_install_info,
	*			      out a(sss) changes);
	 */

	/**
	 * DisableUnitFiles(in  as files,
	 *                  in  b runtime,
	 *                  out a(sss) changes);
	 */

	switch sc {

	case svc.Enable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime, enableForce)
		fmt.Println(call.Body)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
	case svc.Disable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime)
		fmt.Println(call.Body)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
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
