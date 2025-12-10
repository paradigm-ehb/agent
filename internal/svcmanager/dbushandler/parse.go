// Package that handles the parsing of dbus types
package dbushandler

import (
	"fmt"
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// TODO: fix the return type to something explicit
// ListUnits(out a(ssssssouso) units);

// Method
// @param chan a(ss), chan a(ss)
// @return nil
func Parse(in chan []svctypes.Ass, out chan []svctypes.Ass) {

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
