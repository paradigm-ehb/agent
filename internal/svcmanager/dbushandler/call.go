package dbushandler

import (
	"github.com/godbus/dbus"
)

/**
*
* Creating a global dbus connection that we
* can pass a to the receiver
* to other methods
* @return
* a pointer to a dbus.Conn
*
* */

func CreateSystemdObject(conn *dbus.Conn) dbus.BusObject {

	return conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
}
