// Package that handles dbus connections
package dbushandler

import (
	"fmt"
	"github.com/godbus/dbus"
)

func CreateSessionBus() (*dbus.Conn, error) {

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make session bus")
	}

	return conn, nil
}

func CreateSystemBus() (*dbus.Conn, error) {

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to create System bus")
	}

	return conn, nil

}
