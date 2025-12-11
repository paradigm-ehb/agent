package svcmanager

import (
	"fmt"

	"github.com/godbus/dbus"
	dh "paradigm-ehb/agent/internal/svcmanager/dbus"
	svc "paradigm-ehb/agent/internal/svcmanager/systemd"
	svctypes "paradigm-ehb/agent/internal/svcmanager/types"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func RunAction(conn *dbus.Conn, ac svc.UnitAction, service string) error {

	obj := dh.CreateSystemdObject(conn)
	if !obj.Path().IsValid() {
		fmt.Println("invalid systemd path")
	}

	call := obj.Call(string(ac), 0, service, "replace")

	if call.Err != nil {
		fmt.Println("error on action, ", call.Err)
	}

	return nil
}

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunSymlinkAction(conn *dbus.Conn, sc svc.UnitFileAction, enableForRunTime bool, enableForce bool, service []string) error {

	obj := dh.CreateSystemdObject(conn)

	if !obj.Path().IsValid() {
		fmt.Println("invalid systemd path")
	}

	/** EnableUnitFiles(in  as files, in  b runtime, in  b force, out b carries_install_info, out a(sss) changes); */
	/** DisableUnitFiles(in  as files, in  b runtime, out a(sss) changes); */

	switch sc {

	case svc.UnitFileActionEnable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime, enableForce)
		fmt.Println(call.Body)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
	case svc.UnitFileActionDisable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime)
		fmt.Println(call.Body)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
	}

	return nil
}

// @param, true for all on disk, false for loaded units
func RunRetrieval(conn *dbus.Conn, all bool) error {

	obj := dh.CreateSystemdObject(conn)

	if all {

		ch := make(chan []svctypes.UnitFileEntry)
		parse := make(chan []svctypes.UnitFileEntry)

		go svc.GetAllUnits(obj, ch)
		go dh.ParseUnitFileEntries(ch, parse)
		<-parse

	} else if !all {

		ch := make(chan []svctypes.LoadedUnit)
		parse := make(chan []svctypes.LoadedUnit)

		go svc.GetLoadedUnits(obj, ch)
		go dh.ParseLoadedUnits(ch, parse)
		<-parse

	} else {
		return fmt.Errorf("failed parameter")
	}

	return nil
}

func GetStatus(obj dbus.BusObject, name string) {

	call := obj.Call("org.freedesktop.systemd1.Manager.GetUnitFileState", dbus.Flags(dbus.NameFlagReplaceExisting), name)
	// DEBUG
	call.Path.IsValid()

}
