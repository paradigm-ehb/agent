package svcmanager

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

// TODO: decide on using log or fmt for once :)
//
type System struct {

	os string
	processes []process
	services  []service
	displayManager string
}

type process struct {

	// TODO: Get all running processes

}


type service struct {

	name string 
	procId uint32
	owner string

}


func Init() error {

	var sys System


	/**
	* Creating a Session Dbus Connection
	* */
	sesssionConn, err := createSessionBus()
	if err != nil {
		log.Panic("failed to create a dbus session")
	}

	defer sesssionConn.Close()

	// testing this
 	sesssionConn.Hello()

	err = sys.getServicesSessionBus(sesssionConn)
	if err != nil {
		fmt.Println("list object error")
	} 

	err = sys.getDP(sesssionConn)
	if err != nil {
		fmt.Println("No display manager found")
	}


	/**
	* Creating a System Dbus Connection
	* */

	sysConn ,err := createSystemBus()
	if err != nil {
	}
	
	defer sysConn.Close()

	err = sys.getServicesSessionBus(sesssionConn)
	if err != nil {
		fmt.Println("list object error")
	} 

	err = sys.getDP(sesssionConn)
	if err != nil {
		fmt.Println("No display manager found")
	}


	// TODO: add error handling in some way to this

	// DBG: testing this
	sysConn.Hello()
		
	// DBG: debug sys
	fmt.Println(sys)	


	/**
	*  
	* Stop mariadb
	* when passing a service you must add the extension
	* .service, ...
	* */


	fmt.Println("\n\n\n\n\nStopping mariadb")
	err = stopUnit(sysConn, "mariadb.service")
	if err != nil {
		fmt.Println(err)
	}

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



/**
*
* Saves a list of all running processes on a server
* to the System reference 
* By getting all running DBus names
* @return 
*   error
*
* */


func (sys *System) getServicesSessionBus(conn *dbus.Conn) (error) {

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

func (sys *System) getServicesSystemBus(conn *dbus.Conn) (error) {

	var output []string

	obj := conn.Object("org.freedesktop.DisplayManager", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&output)

	for i := range len(output) {

		var ser service
		ser.name = output[i]
		sys.services = append(sys.services, ser)
	}

	return nil
}


/**
* Returns the running dislay manager
*
* @return
* []string, error
*
* returns an empty string if no Display Managers are available
* */

func (sys *System) getDP(conn *dbus.Conn) error  {

	var result []string

	obj := conn.Object("org.freedesktop.DisplayManager", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)

	if len(result) > 0 {

		sys.displayManager = result[0]
	} else {
		return fmt.Errorf("no error found")
	}

	return nil
}


func getStatus(conn *dbus.Conn, name string) error {

	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1/unit")
	obj.Call("org.freedesktop.systemd1.Manager.GetUnitFileStatus", dbus.Flags(dbus.NameFlagReplaceExisting), name)

	// DEBUG: placeholders
	fmt.Println(obj)
	fmt.Println(name)


	// TODO:  retrieve all units in system

	return nil

}

/*

      GetUnitFileState(in  s file,
                       out s state);
      EnableUnitFiles(in  as files,
                      in  b runtime,
                      in  b force,
                      out b carries_install_info,
                      out a(sss) changes);
      DisableUnitFiles(in  as files,
                       in  b runtime,
                       out a(sss) changes);

					   */

func stopUnit(conn *dbus.Conn, name string) error {
	
	// bugfix -> 
	// failed to disable mariadb, Unknown method StopUnit or interface org.freedesktop.systemd1.Manager.  
	//
	// Failed to disable mariadb.service, Access denied as the requested operation requires interactive authentication. Howev│
	// er, interactive authentication has not been enabled by the calling program. 
	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	call := obj.Call("org.freedesktop.systemd1.Manager.StopUnit", dbus.FlagAllowInteractiveAuthorization, name, "replace")
	if call.Err != nil {
		return fmt.Errorf("failed to disable %s, %v",name, call.Err)
	}


	fmt.Println(call.Body)

	return nil


}

/**
* sudo dbus-send --system --print-reply --dest=org.freedesktop.systemd1 /org/freedesktop/systemd1 org.freedesktop.systemd1.Manager.StopUnit
* string:"nginx.service" string:"replace"
*/
