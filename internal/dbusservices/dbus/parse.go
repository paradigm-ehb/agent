// Package that handles the parsing of dbus types
package dbushandler

import (
	types "paradigm-ehb/agent/internal/dbusservices/types"
	"strings"
)

// TODO(nasr): implement interfaces maybe

// Method
// @param chan a(ss), chan a(ss)
// @param chan UnitFileEntry, chan UnitFileEntry
// @return nil
func ParseUnits(in chan []types.Unit, out chan []types.Unit) {

	input := <-in

	/**
	filter the units on services and remove devices etc
	*/
	buffer := make([]types.Unit, 0, len(input))
	for _, value := range input {

		if strings.HasSuffix(".service", value.Name) {
			buffer = append(buffer, value)
		}
	}

	out <- buffer

	return

}

func ParseLoadedUnits(in chan []types.LoadedUnit, out chan []types.LoadedUnit) {

	input := <-in

	/**
	filter the units on services and remove devices etc
	*/
	buffer := make([]types.LoadedUnit, 0, len(input))
	for _, value := range input {

		if strings.HasSuffix(".service", string(value.Name)) {
			buffer = append(buffer, value)
		}
	}

	out <- buffer

	return

}
