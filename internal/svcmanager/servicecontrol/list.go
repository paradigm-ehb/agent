package servicecontrol

import (
	"fmt"
	"github.com/godbus/dbus"
)

// ListUnitFiles() returns an array of unit names plus their enablement status.
// Note that ListUnit() returns a list of units currently loaded into memory, while ListUnitFiles()
// returns a list of unit files that could be found on disk. Note that while most units are read directly from a
// unit file with the same name some units are not backed by files, and some
// files (templates) cannot directly be loaded as units but need to be instantiated.
// ---------------------------------------------------------------------------------------
// Method returns an array of all currently loaded units,

func GetLoadedUnits(obj dbus.BusObject) any {

	// TODO: replace any with unit interface
	var result any
	// takes no in
	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0)
	if call.Err != nil {
		fmt.Printf("failed to list unit files that are loaded in memory %v", call.Err)
		return nil
	}

	call.Store(&result)

	return result
}

func GetAllUnits(obj dbus.BusObject, ch chan [][]string) <-chan [][]string {

	// ListUnitFiles(out a(ss) files);
	// an array of struct string string
	// i think

	var result [][]string

	go func(in [][]string) {

		call := obj.Call("org.freedesktop.systemd1.Manager.ListUnitFiles", 0)

	}()

	if call.Err != nil {

		fmt.Println("failed to retrieve unit files")
		return
	}

	call.Store(&result)

	ch <- result

	return ch
}
