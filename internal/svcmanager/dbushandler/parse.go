// Package that handles the parsing of dbus types
package dbushandler

import (
	"fmt"
	svctypes "paradigm-ehb/agent/internal/svcmanager/system"
)

// TODO: implement interfaces maybe

// Method
// @param chan a(ss), chan a(ss)
// @return nil
func ParseAllUnits(in chan []svctypes.Ass, out chan []svctypes.Ass) {

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

func ParseLoadedUnits(in chan []svctypes.Assssssouso, out chan []svctypes.Assssssouso) {

	input := <-in

	for i := range input {

		fmt.Println(input[i])
	}

	out <- input

}
