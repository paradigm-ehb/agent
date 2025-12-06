// Package that handles dbus connections
package dbushandler

import (
	"fmt"
	"github.com/godbus/dbus"
)

// Creates a global session bus connection
// Session bus is bound to the user session
// @return *dbus.Conn, error
// @param none
func CreateSessionBus() (*dbus.Conn, error) {

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make session bus")
	}

	return conn, nil
}

// Creates a global session bus connection
// Session bus is bound to the user session
// @return *dbus.Conn, error
// @param none
func CreateSystemBus() (*dbus.Conn, error) {

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to create System bus")
	}

	return conn, nil

}
