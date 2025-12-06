package servicecontrol

import (
	"fmt"
	"github.com/godbus/dbus"
)

func HandleActionOnUnit(obj dbus.BusObject, name string, action Action) error {

	switch action {

	case Action(Start):

		call := obj.Call("org.freedesktop.systemd1.Manager.StartUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
			// DEBUG
			fmt.Println(call.Body)
			return fmt.Errorf("failed to start %s, %v", name, call.Err)
		}

	case Action(Stop):

		call := obj.Call("org.freedesktop.systemd1.Manager.StopUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
			// DEBUG
			fmt.Println(call.Body)
			return fmt.Errorf("failed to stop %s, %v", name, call.Err)
		}

	case Action(Restart):

		call := obj.Call("org.freedesktop.systemd1.Manager.RestartUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
			// DEBUG
			fmt.Println(call.Body)
			return fmt.Errorf("failed to restart %s, %v", name, call.Err)
		}

	}
	return nil
}

/**
*
*  Enable or disable a unit file
*
* */
func HandleSymlinkCreationAction(obj dbus.BusObject, name []string, action SymlinkAction, enableForRunTime bool, replaceExistingSynmlink bool) error {

	/**

	EnableUnitFiles(in  as files,
	in  b runtime,
	in  b force,
	out b carries_install_info,
	out a(sss) changes);

	creates a symllink in /run or somethign

	*/

	switch action {

	case SymlinkAction(Enable):
		call := obj.Call("org.freedesktop.systemd1.Manager.EnableUnitFiles", dbus.Flags(dbus.NameFlagReplaceExisting), name, enableForRunTime, replaceExistingSynmlink)
		if call.Err != nil {
			// DEBUG
			fmt.Println("response body", call.Body)
			return fmt.Errorf("failed to enable a unit file %v", call.Err)
		}

	case SymlinkAction(Disable):
		call := obj.Call("org.freedesktop.systemd1.Manager.DisableUnitFiles", dbus.Flags(dbus.NameFlagReplaceExisting), name, enableForRunTime)
		if call.Err != nil {
			// DEBUG
			fmt.Println("response body", call.Body)
			return fmt.Errorf("failed to disable a unit file %v", call.Err)
		}
	}

	return nil
}

func GetStatus(obj dbus.BusObject, name string) {

	call := obj.Call("org.freedesktop.systemd1.Manager.GetUnitFileState", dbus.Flags(dbus.NameFlagReplaceExisting), name)
	// DEBUG
	call.Path.IsValid()

}
