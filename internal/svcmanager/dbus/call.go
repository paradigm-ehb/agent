package dbushandler

import (
	"github.com/godbus/dbus"
)

/**
* Creating a systemd Object
* used to handle services
* @return BusObject
* */
func CreateSystemdObject(conn *dbus.Conn) dbus.BusObject {

	return conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
}
