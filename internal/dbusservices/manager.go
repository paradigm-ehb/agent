package dbus_services

import (
	"fmt"
	"log"

	dbushelper "paradigm-ehb/agent/internal/dbusservices/dbus"
	systemd "paradigm-ehb/agent/internal/dbusservices/systemd"
	types "paradigm-ehb/agent/internal/dbusservices/types"

	"github.com/godbus/dbus"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunAction(
	conn *dbus.Conn,
	ac systemd.UnitAction,
	service string) error {

	obj, _ := dbushelper.CreateSystemdObject(conn)
	if !obj.Path().IsValid() {
		return fmt.Errorf("object path is invalid")
	}

	call := obj.Call(string(ac), 0, service, "replace")

	if call.Err != nil {
		return fmt.Errorf("failed to execute object on in unit action, %v", call.Err)
	}

	return nil
}

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunSymlinkAction(
	conn *dbus.Conn,
	sc systemd.UnitFileAction,
	enableForRunTime bool,
	enableForce bool,
	service []string) error {

	obj, _ := dbushelper.CreateSystemdObject(conn)

	if !obj.Path().IsValid() {
		fmt.Println("invalid systemd path")
	}

	/** EnableUnitFiles(in  as files, in  b runtime, in  b force, out b carries_install_info, out a(sss) changes); */
	/** DisableUnitFiles(in  as files, in  b runtime, out a(sss) changes); */

	switch sc {

	case systemd.UnitFileActionEnable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime, enableForce)
		if call.Err != nil {
			return fmt.Errorf("error %v", call.Err)
		}
	case systemd.UnitFileActionDisable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
	}

	return nil
}

// MapLoadedUnits /*
func MapLoadedUnits(conn *dbus.Conn) []*types.LoadedUnit {

	obj, _ := dbushelper.CreateSystemdObject(conn)
	ch := make(chan []types.LoadedUnit)
	parse := make(chan []types.LoadedUnit)

	go systemd.GetLoadedUnits(obj, ch)
	go dbushelper.ParseLoadedUnits(ch, parse)

	loaded := <-parse

	units := make([]*types.LoadedUnit, 0, len(loaded))

	for _, u := range loaded {
		units = append(units, &types.LoadedUnit{
			Name:        u.Name,
			Description: u.Description,
			LoadState:   u.LoadState,
			SubState:    u.SubState,
			ActiveState: u.ActiveState,
			DepUnit:     u.DepUnit,
			ObjectPath:  u.ObjectPath,
			/*oops typo in queued job :)*/
			QueudJob: u.QueudJob,
			JobType:  u.JobType,
			JobPath:  u.JobPath,
		})
	}

	return units
}

func MapFilteredUnits(conn *dbus.Conn, filters []string) ([]*types.LoadedUnit, error) {

	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		fmt.Errorf("failed to create systemd object")
	}

	in := make(chan []types.LoadedUnit)
	out := make(chan []types.LoadedUnit)

	go systemd.GetUnitsFiltered(obj, in, filters)
	go dbushelper.ParseLoadedUnits(in, out)

	var entries []types.LoadedUnit

	units := make([]*types.LoadedUnit, 0, len(entries))

	entries = <-out

	for _, e := range entries {

		units = append(units, &types.LoadedUnit{
			Name:        e.Name,
			Description: "Not available",
			LoadState:   e.LoadState,
			SubState:    "Not Available",
			ActiveState: "Not Available",
			DepUnit:     "Not Available",
			ObjectPath:  "Not Available",
			QueudJob:    0,
			JobType:     "Not Available",
			JobPath:     "Not Available",
		})
	}

	return units, nil
}

func MapUnits(conn *dbus.Conn) []*types.LoadedUnit {

	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		fmt.Errorf("failed to create systemd object")
	}

	in := make(chan []types.Unit)
	out := make(chan []types.Unit)

	go systemd.GetUnits(obj, in)
	go dbushelper.ParseUnits(in, out)

	var entries []types.Unit

	units := make([]*types.LoadedUnit, 0, len(entries))

	entries = <-out

	for _, e := range entries {

		units = append(units, &types.LoadedUnit{
			Name:        e.Name,
			Description: "Not available",
			LoadState:   e.State,
			SubState:    "Not Available",
			ActiveState: "Not Available",
			DepUnit:     "Not Available",
			ObjectPath:  "Not Available",
			QueudJob:    0,
			JobType:     "Not Available",
			JobPath:     "Not Available",
		})
	}

	return units
}

/*
*
* @param, true for all on disk, false for loaded units
* a loaded unit is a unit that has been activated before
* and is available in memoery for the server to start up
* or something like that
* @return []*types.LoadedUnit, error
 */
func RunRetrieval(
	conn *dbus.Conn,
	requestAllUnitsOnDisk bool,
) ([]*types.LoadedUnit, error) {

	conn, err := dbushelper.CreateSystemBus()

	if err != nil {
		log.Printf("failed to create a system bus connection for retrieving units %v", err)
	}

	if requestAllUnitsOnDisk {
		return MapUnits(conn), nil
	}

	return MapLoadedUnits(conn), nil
}

func UnitStatus(
	obj dbus.BusObject,
	name string) (string, error) {

		out, err := systemd.GetStatusCall(obj, name)
		if err != nil {
			return "Failed", fmt.Errorf("failed to execute status call ", err)
		}

		return out, nil 

}
