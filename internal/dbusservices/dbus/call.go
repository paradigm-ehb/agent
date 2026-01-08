package dbushandler

import (
	"fmt"
	"github.com/godbus/dbus"
)

/*
* CreateSystemdObject
* used to handle services
* @return BusObjcet
 */
func CreateSystemdObject(conn *dbus.Conn) (dbus.BusObject, error) {

	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	if obj == nil {

		return nil, fmt.Errorf("failed to create a systemd object")
	}

	return obj, nil
}

/**

* CreateLoginObject	
* used to handle services
* @return BusObject 

*/
func CreateLoginObject(conn *dbus.Conn) (dbus.BusObject, error) {
	obj := conn.Object(
		"org.freedesktop.login1",
		"/org/freedesktop/login1",
	)

	if obj == nil {

		return nil, fmt.Errorf("failed to create a login object")
	}

	return obj, nil

}
