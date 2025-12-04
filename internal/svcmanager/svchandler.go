package svcmanager 

import (

	"fmt"
	"github.com/godbus/dbus"
	// "github.com/coreos/go-systemd"
)



// define the dbus typer
type types struct {

	// TODO: create the binding types for dbus
	

}


/**
*
*
* bus -> proces = could be anything a unix proces something 
* connection -> think about it like a dns name
*
* ------------
* what the url, 
* object -> the path
* interface -> the type youre going to call
* member -> the url
*
* tools
*
* part of systemd -> 
* dbus-send
* busctl
*
*
*/


func ListDbusObject() ([]string, error) {

	var result []string
		
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make DBus connection for health check")
	}

	defer conn.Close()


    obj := conn.Object("org.freedesktop.DBus", "/")
    obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)

	return result, nil 


}

func dbusConnectionHealth(*dbus.Conn) (bool, error) {

	var result []string
		
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, fmt.Errorf("failed to make DBus connection for health check")
	}

    obj := conn.Object("org.freedesktop.DBus", "/")
    obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)
	
	// TODO: do an actual check 
	return true, nil


}

/**
*
* Method for testing how the Dbus objects work
* and how the units are handled
*
* */
func GetDisplayManager() ([]string, error)  {

	var result []string

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to make DBus connection for health check")
	}

	defer conn.Close()
	
	obj := conn.Object("org.freedesktop.DisplayManager", "/")
	obj.Call("org.freedesktop.DBus.ListNames", 0).Store(&result)


	return result, nil
}


/*
*
* Apparanent authentication implementation
==== AUTHENTICATING FOR org.freedesktop.systemd1.manage-units ====
Authentication is required to stop 'mariadb.service'.
Authenticating as: nasr
Password:
==== AUTHENTICATION COMPLETE ====
**/


func authDbus() {


}

/*
* UInt32 StartServiceByName(String name,
*                           UInt32 flags)
*/

func startService() {
	
	// TODO:

}

func stopService() {

	// TODO:
}

func restartService() {
	
	// TODO:
}

func statusService() {

	// TODO:
}


