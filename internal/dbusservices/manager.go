package dbus_services

import (
	"fmt"

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

	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		return fmt.Errorf("failed to create systemd object %v", err)
	}

	if !obj.Path().IsValid() {
		return fmt.Errorf("invalid systemd path")
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

func UnitStatus(
	obj dbus.BusObject,
	name string) (string, error) {

	out, err := systemd.GetStatusCall(obj, name)
	if err != nil {
		return "Failed", fmt.Errorf("failed to execute status call %v", err)
	}

	return out, nil

}

/**
 * Map to the correct types
 */
func MapLoadedUnits(conn *dbus.Conn) ([]*types.LoadedUnit, error) {
	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create systemd object: %w", err)
	}
	
	loaded, err := systemd.GetLoadedUnits(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to get loaded units: %w", err)
	}
	
	parsed, err := dbushelper.ParseLoadedUnits(loaded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse loaded units: %w", err)
	}
	
	units := make([]*types.LoadedUnit, 0, len(parsed))
	for _, u := range parsed {
		units = append(units, &types.LoadedUnit{
			Name:        u.Name,
			Description: u.Description,
			LoadState:   u.LoadState,
			SubState:    u.SubState,
			ActiveState: u.ActiveState,
			DepUnit:     u.DepUnit,
			ObjectPath:  u.ObjectPath,
			QueudJob:    u.QueudJob, // Keep the typo for consistency
			JobType:     u.JobType,
			JobPath:     u.JobPath,
		})
	}
	return units, nil
}

func MapFilteredUnits(conn *dbus.Conn, filters []string) ([]*types.LoadedUnit, error) {
	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create systemd object: %w", err)
	}
	
	entries, err := systemd.GetUnitsFiltered(obj, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered units: %w", err)
	}
	
	parsed, err := dbushelper.ParseLoadedUnits(entries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filtered units: %w", err)
	}
	
	units := make([]*types.LoadedUnit, 0, len(parsed))
	for _, e := range parsed {
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

/**
 * Fill in proper unit types
 */
func MapUnits(conn *dbus.Conn) ([]*types.LoadedUnit, error) {
	obj, err := dbushelper.CreateSystemdObject(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create systemd object: %w", err)
	}
	
	result, err := systemd.GetUnits(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to get units: %w", err)
	}
	
	parsedUnits, err := dbushelper.ParseUnits(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse units: %w", err)
	}
	
	units := make([]*types.LoadedUnit, 0, len(parsedUnits))
	for _, e := range parsedUnits {
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
	return units, nil
}

/**
 * @param requestAllUnitsOnDisk: true for all units on disk, false for loaded units
 * A loaded unit is a unit that has been activated before and is available in memory
 * for the server to start up
 * @return []*types.LoadedUnit, error
 */
func RunRetrieval(requestAllUnitsOnDisk bool) ([]*types.LoadedUnit, error) {
	conn, err := dbushelper.CreateSystemBus()

	if err != nil {
		return nil, fmt.Errorf("failed to create system bus connection: %w", err)
	}
	
	if requestAllUnitsOnDisk {
		return MapUnits(conn)
	}
	
	return MapLoadedUnits(conn)
}

func RunRetrievalDeprecated(conn *dbus.Conn, requestAllUnitsOnDisk bool) ([]*types.LoadedUnit, error) {
	
	if requestAllUnitsOnDisk {
		return MapUnits(conn)
	}
	
	return MapLoadedUnits(conn)
}
