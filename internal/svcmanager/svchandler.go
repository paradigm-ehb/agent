package svcmanager 

import (
	"fmt"
	"github.com/godbus/dbus"
	"log"
)


type process struct {

	// TODO: Get all running processes

}

type System struct {

	os string
	processes []process
	services  []service
}

type service struct {

	name string 
	procId uint32
	owner string

}


func Init() error {

	var sys System


	conn, err := CreateDbusSession()
	if err != nil {
		log.Panic("failed to create a dbus session")
	}

	defer conn.Close()

	err = sys.GetProcesses(conn)
	if err != nil {
		fmt.Println("list object error")
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

func CreateDbusSession() (*dbus.Conn, error) {

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make DBus connection for health check")
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


func (sys *System) GetProcesses(conn *dbus.Conn) (error) {

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


/**
* Returns the running dislay manager
*
* @return
* []string, error
*
* returns an empty string if no Display Managers are available
* */

func GetDPObject(conn *dbus.Conn) ([]string, error)  {

	var result []string

	obj := conn.Object("org.freedesktop.DisplayManager", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)


	return result, nil
}


func StartService(proc process) uint32 {


	// TODO: starting a service using StartServiceByName

	return 33 
}

