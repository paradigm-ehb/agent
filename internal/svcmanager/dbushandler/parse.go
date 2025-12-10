// Package that handles the parsing of dbus types
package dbushandler

import (
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// TODO: fix the return type to something explicit
// ListUnits(out a(ssssssouso) units);

// Method
// parses
// @return [][]string
func Parse(in chan svctypes.Ass, out chan svctypes.Ass) {

	input := <-in
	// for i := range input {
	// 	for j := range input {
	//
	// 		input[i][j] += "\n"
	// 	}
	// }

	out <- input
}
