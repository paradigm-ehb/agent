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

	/**
	*
	* Stop mariadb
	* when passing a service you must add the extension
	* .service, ...
	* */

	fmt.Println("\n\n\n\n\nStopping mariadb")
	err = handleActionOnUnit(obj, "mariadb.service", Action(Start))
	if err != nil {
		fmt.Println(err)
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

func getUnits(obj *dbus.BusObject) error {

	// TODO
	// GetUnitProcesses(in  s name,
	//                  out a(sus) processes);
	return nil
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
