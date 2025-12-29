// Package that handles the parsing of dbus types
package dbushandler

import (
	"fmt"
	svctypes "paradigm-ehb/agent/internal/dbus_services/types"
)

// TODO: implement interfaces maybe

// Method
// @param chan a(ss), chan a(ss)
// @param chan UnitFileEntry, chan UnitFileEntry
// @return nil
func ParseUnitFileEntries(in chan []svctypes.UnitFileEntry, out chan []svctypes.UnitFileEntry) {

	input := <-in

	for i := range input {

		if input[i].State == "enabled" {
			fmt.Println("==================Enabled======================")
			fmt.Println(input[i].Name)
		} else {

			fmt.Println("==================Disabled======================")
			fmt.Println(input[i].Name)
		}
	}
	out <- input

}

func ParseLoadedUnits(in chan []svctypes.LoadedUnit, out chan []svctypes.LoadedUnit) {

	input := <-in

	for i := range input {

		fmt.Println(input[i])
	}

	out <- input

}
