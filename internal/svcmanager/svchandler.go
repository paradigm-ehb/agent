package svcmanager

import (
	"fmt"

	"github.com/godbus/dbus"
)

// Values to send if the a call returns an error
type ServerState int

const (
	Healthy ServerState = iota
	Unhealthy
	AttentionNeeded
)

type Action string

const (
	Start   Action = "org.freedesktop.systemd1.Manager.StartUnit"   // start unit
	Stop    Action = "org.freedesktop.systemd1.Manager.StopUnit"    // stop unit
	Restart Action = "org.freedesktop.systemd1.Manager.RestartUnit" // restart unit
)

type SymlinkAction string

const (
	Enable  Action = "org.freedesktop.systemd1.Manager.EnableUnitFiles"  // enable unit(s)
	Disable Action = "org.freedesktop.systemd1.Manager.DisableUnitFiles" // disable unit(s)

)

type System struct {
	os        string    // os version
	processes []process // list of processes on the server
	services  []service // list of services on the server
}

type process struct {
	name string // process name
	id   uint32 // PID
}

type service struct {
	name  string // unit file name
	id    uint32 // service PID
	owner string // unit file owner
}

func Init() error {

	var sys System

	sysConn, err := createSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer sysConn.Close()

	obj := createSystemdObject(sysConn)

	arr, err := getUnits(obj)
	if err != nil {
		fmt.Println("failed to get array of units")
	} else {
		fmt.Println("array of units", arr)
	}

	/**
	*
	* Stop mariadb
	* when passing a service you must add the extension
	* .service, ...
	* */

	name := "mariadb.service"
	namesList := []string{"mariadb.service"}

	enableForRunTime := true
	replaceExistingSynmlink := true

	fmt.Println("\n\n\n\n\nStopping mariadb")
	err = handleActionOnUnit(obj, name, Action(Start))
	if err != nil {
		fmt.Println("failed upper layer starting the service, will try to enable the service before starting it", err)
		// try to enable or disable the service!!
		//
		err = handleSymlinkCreationAction(obj, namesList, SymlinkAction(Enable), enableForRunTime, replaceExistingSynmlink)
		if err != nil {
			fmt.Println("failed upper layer enabling the service, will try to enable the service before starting it", err)
		} else {
			fmt.Println("\n\nlist of unit files: ", arr)
		}
	}

	// DBG: debug sys
	// TODO: add a proper way of printing stuff out
	fmt.Println(sys)

	return nil
}

/**
*
* Creating a global dbus connection that we
* can pass a to the receiver
* to other methods
* @return
* a pointer to a dbus.Conn
*
* */

func createSessionBus() (*dbus.Conn, error) {

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make session bus")
	}

	return conn, nil
}

func createSystemBus() (*dbus.Conn, error) {

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to create System bus")
	}

	return conn, nil

}

func createSystemdObject(conn *dbus.Conn) dbus.BusObject {

	return conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1/unit")
}

/**
*
* Previousl
*
* Saves a list of all running processes on a server
* to the System reference
* By getting all running DBus names
* @return
*   error
*
* */

func (sys *System) getServicesSessionBus(conn *dbus.Conn) error {

	var output []string

	obj := conn.Object("org.freedesktop.DBus", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&output)

	for i := range len(output) {

		var ser service
		ser.name = output[i]
		sys.services = append(sys.services, ser)

	}

	return nil
}

func getStatus(obj dbus.BusObject, name string) {

	obj.Call("org.freedesktop.systemd1.Manager.GetUnitFileState", dbus.Flags(dbus.NameFlagReplaceExisting), name)

}

func getUnits(obj dbus.BusObject) ([]string, error) {

	var result []string
	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnitFiles", dbus.FlagAllowInteractiveAuthorization, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("failed to list unit files %v", call.Err)
	}
	call.Store(&result)

	return result, nil
}

func handleActionOnUnit(obj dbus.BusObject, name string, action Action) error {

	switch action {

	case Action(Start):

		call := obj.Call("org.freedesktop.systemd1.Manager.StartUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
			fmt.Println(call.Body)
			return fmt.Errorf("failed to start %s, %v", name, call.Err)
		}

	case Action(Stop):

		call := obj.Call("org.freedesktop.systemd1.Manager.StopUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
			fmt.Println(call.Body)
			return fmt.Errorf("failed to stop %s, %v", name, call.Err)
		}

	case Action(Restart):

		call := obj.Call("org.freedesktop.systemd1.Manager.RestartUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
		if call.Err != nil {
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
func handleSymlinkCreationAction(obj dbus.BusObject, name []string, action SymlinkAction, enableForRunTime bool, replaceExistingSynmlink bool) error {

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
		call := obj.Call("org.freedesktop.systemd1.Manager.EnableUnitFiles", dbus.FlagAllowInteractiveAuthorization, name, enableForRunTime, replaceExistingSynmlink)
		if call.Err != nil {
			fmt.Println(call.Body)
			return fmt.Errorf("failed to enable a unit file %v", call.Err)
		}

	case SymlinkAction(Disable):
		call := obj.Call("org.freedesktop.systemd1.Manager.EnableUnitFiles", dbus.FlagAllowInteractiveAuthorization, name)
		if call.Err != nil {
			fmt.Println(call.Body)
			return fmt.Errorf("failed to enable a unit file %v", call.Err)
		}

	}

	return nil
}
