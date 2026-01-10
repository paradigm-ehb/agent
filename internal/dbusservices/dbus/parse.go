// Package that handles the parsing of dbus types
package dbushandler

import (
	types "paradigm-ehb/agent/internal/dbusservices/types"
	"strings"
)

// TODO(nasr): implement interfaces maybe

// ParseUnits filters units to only include services
// @param input []types.Unit
// @return []types.Unit, error
func ParseUnits(input []types.Unit) ([]types.Unit, error) {
	/**
	filter the units on services and remove devices etc
	*/
	buffer := make([]types.Unit, 0, len(input))
	for _, value := range input {
		if strings.HasSuffix(value.Name, ".service") {
			buffer = append(buffer, value)
		}
	}
	return buffer, nil
}

// ParseLoadedUnits filters loaded units to only include services
// @param input []types.LoadedUnit
// @return []types.LoadedUnit, error
func ParseLoadedUnits(input []types.LoadedUnit) ([]types.LoadedUnit, error) {
	/**
	filter the units on services and remove devices etc
	*/
	buffer := make([]types.LoadedUnit, 0, len(input))
	for _, value := range input {
		if strings.HasSuffix(value.Name, ".service") {
			buffer = append(buffer, value)
		}
	}
	return buffer, nil
}
