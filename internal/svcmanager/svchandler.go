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

	// DEBUG
	fmt.Println(sysConn.Names())

	obj := createSystemdObject(sysConn)

	arrOnRam := getCurrentlyLoadedUnits(obj)
	// DEBUG
	fmt.Println("array of units", arrOnRam)

	arrOnDisk, err := getAllUnitsOnDisk(obj)
	if err != nil {
		return nil
	}

	fmt.Println("printing unit files that are on the disk", arrOnDisk)

	name := "mariadb.service"
	namesList := []string{"mariadb.service"}

	enableForRunTime := true
	replaceExistingSynmlink := true

	// DEBUG
	fmt.Println("\n\n\n\n\nStopping mariadb")
	err = handleActionOnUnit(obj, name, Action(Start))
	if err != nil {
		// DEBUG
		fmt.Println("failed to start the unit: ", err)
		err = handleSymlinkCreationAction(obj, namesList, SymlinkAction(Enable), enableForRunTime, replaceExistingSynmlink)
		if err != nil {
			// DEBUG
			fmt.Println("failed to enable the unit: ", err)
		}
	}

	// DEBUG
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

	call := obj.Call("org.freedesktop.systemd1.Manager.GetUnitFileState", dbus.Flags(dbus.NameFlagReplaceExisting), name)
	// DEBUG
	call.Path.IsValid()

}

// ListUnitFiles() returns an array of unit names plus their enablement status.
// Note that ListUnit() returns a list of units currently loaded into memory, while ListUnitFiles()
// returns a list of unit files that could be found on disk. Note that while most units are read directly from a
// unit file with the same name some units are not backed by files, and some
// files (templates) cannot directly be loaded as units but need to be instantiated.
// ---------------------------------------------------------------------------------------
// Method returns an array of all currently loaded units,
func getCurrentlyLoadedUnits(obj dbus.BusObject) any {

	// ListUnits(out a(ssssssouso) units);
	// crazy return type hhhhhhh
	// going to any this for now
	// TODO: fix the return type to something explicit
	// HOW: create a giant struct that handles everything

	var result any
	// takes no in
	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0)
	if call.Err != nil {
		fmt.Printf("failed to list unit files that are loaded in memory %v", call.Err)
		return nil
	}

	call.Store(&result)

	return result
}

func getAllUnitsOnDisk(obj dbus.BusObject) ([][]string, error) {

	// ListUnitFiles(out a(ss) files);
	// an array of struct string string
	// i think

	var result [][]string

	// takes no in either
	call := obj.Call("org.freedesktop.systemd1.Manager.ListUnitsFiles", 0)
	if call.Err != nil {
		return nil, fmt.Errorf("failed to list unit files that on disk %v", call.Err)
	}

	call.Store(&result)

	return result, nil

}

func handleActionOnUnit(obj dbus.BusObject, name string, action Action) error {

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

// TODO: get unit object path using GetUnit() may be used to get the unit object path for a unit name.
// It takes the unit name and returns the object path. If a unit has not been loaded yet by this name this call will fail.
// doing this will allow us to get all units running, get their object path and do handles on them
