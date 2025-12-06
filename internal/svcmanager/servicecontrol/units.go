package servicecontrol

import (
	"github.com/godbus/dbus"
)

/*
* GetUnit(in  s name,
* out o unit);
*
* Retrieves the object path of a unit
*
* @return dbus.ObjectPath
 */
func GetUnit(obj dbus.BusObject, name string) dbus.ObjectPath {

	var result dbus.ObjectPath
	obj.Call("org.freedesktop.systemd1.Manager.GetUnit", 0, name).Store(&result)

	return result
}
