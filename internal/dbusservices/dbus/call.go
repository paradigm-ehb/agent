package dbushandler

import (
	"github.com/godbus/dbus"
)

/*
* CreateSystemdObject
* used to handle services
* @return BusObjcet
 */
func CreateSystemdObject(conn *dbus.Conn) dbus.BusObject {

	return conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
}
