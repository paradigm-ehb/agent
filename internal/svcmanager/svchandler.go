package svcmanager 

import (

	"fmt"
	"github.com/godbus/dbus"
	// "github.com/coreos/go-systemd"
)

type Process struct {

	name []byte 
	procId uint32
	owner string

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

	defer conn.Close()

	return conn, nil
} 

/**
*
* Returns a list of all running DBus objects on a server
* @return 
* []string, error
*
* */


func GetDbusObjectList(conn *dbus.Conn) ([]string, error) {

	var result []string


	obj := conn.Object("org.freedesktop.DBus", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)

	return result, nil 
}

func (*Process) formatDBusObjectList(list []string)  {

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


func StartService(proc Process) uint32 {


	// TODO: starting a service using StartServiceByName


	return 33 
}

