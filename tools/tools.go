// Package tools provides small helper utilities used across the project.
package tools

import (

	"fmt"
	"runtime"
)

// Check if the user is running an supported OS
func CheckOSUser() error {

	if os := runtime.GOOS; os != "linux" {
		return fmt.Errorf("Not supported operating system")
	}
	return nil
}



