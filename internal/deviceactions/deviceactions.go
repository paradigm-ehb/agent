package deviceactions

// Package that handles device actions over dbus

import (
	"github.com/godbus/dbus"
)

/*
* PerformDeviceAction(in  s action,
* out o unit)
* Performs the specified device action
*
* @return dbus.ObjectPath
 */
func PerformDeviceAction(obj dbus.BusObject, DeviceActionRequest DeviceAction) dbus.ObjectPath {

	var result dbus.ObjectPath
	err := obj.Call(string(DeviceActionRequest), 0, false).Store(&result)
	if err != nil {
		return ""
	}

	return result
}
