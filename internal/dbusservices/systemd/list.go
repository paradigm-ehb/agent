package servicecontrol

import (
	"fmt"
	types "paradigm-ehb/agent/internal/dbusservices/types"

	"github.com/godbus/dbus"
)

// ListUnitFiles() returns an array of unit names plus their enablement status.
// Note that ListUnit() returns a list of units currently loaded into memory, while ListUnitFiles()
// returns a list of unit files that could be found on disk. Note that while most units are read directly from a
// unit file with the same name some units are not backed by files, and some
// files (templates) cannot directly be loaded as units but need to be instantiated.
// ---------------------------------------------------------------------------------------
// Method returns an array of all currently loaded units,
func GetLoadedUnits(
	obj dbus.BusObject,
	out chan []types.LoadedUnit) {

	var result []types.LoadedUnit

	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0)
	if call.Err != nil {
		fmt.Printf("failed to list unit files that are loaded in memory %v", call.Err)
		return
	}

	err := call.Store(&result)

	if err != nil {
		return
	}

	out <- result

}

func GetUnits(
	obj dbus.BusObject,
) ([]types.Unit, error) {

	// ListUnitFiles(out a(ss) files);
	// an array of struct string string

	var result []types.Unit

	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnitFiles", 0)

	if call.Err != nil {

		return nil, fmt.Errorf("failed to call all units loaded on disk")
	}

	err := call.Store(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to store services ")
	}

	/**
	* THIS IS WORKING
	* fmt.Println("channel result: \n", result)
	 */

	/**
	* THIS ISNT
	* out <- result
	 */

	return result, nil
}

func GetUnitsFiltered(
	obj dbus.BusObject,
	out chan []types.LoadedUnit,
	states []string) {

	var result []types.LoadedUnit

	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnitsFiltered", 0, states)

	if call.Err != nil {
		fmt.Println("failed to call filtered list of units")
		return
	}

	err := call.Store(&result)
	if err != nil {
		return
	}

	out <- result

}

func GetStatusCall(obj dbus.BusObject, name string) (string, error) {
	var result string

	call := obj.Call(
		"org.freedesktop.systemd1.Manager.GetUnitFileState",
		0,
		name,
	)

	if call.Err != nil {
		return "call error: ", call.Err
	}

	if err := call.Store(&result); err != nil {
		return "call store: ", err
	}

	return result, nil
}
