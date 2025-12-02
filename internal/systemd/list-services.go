package controller 

import (

	"fmt"
	"github.com/godbus/dbus"
	"io"
	"os"
)

var object dbus.BusObject

type message struct {
	
	string message

}

type service struct {
	
	servicePath string
	objectPath 	string

}

/**
* define the object we want to connect with
* in this case it's going to be systemd1
* */

func createProxy(conn *dbus.Conn) error {
		
	object.SetProperty("")
	object := dbus.Object("/usr/bin/systemd", , , )
	return nil
}


/**
* Establish a dbus connection
**/

func makeConnection() (*dbus.Conn,  error) { 
	
	var a io.ReadWriteCloser
	var s []string
	

	conn, err := dbus.NewConn(a)	
	if err != nil {
		return nil, fmt.Errorf("failed to make dbus connection")
		os.Exit(1)
	}
	defer conn.Close()

	return conn, nil
} 


func accessBusObject(conn *dbus.Conn, object dbus.Object ) error {

	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return fmt.Errorf("Failed to get list of owned names:")
		os.Exit(1)
	}

	fmt.Println("Currently owned names on the session bus:")
	for _, v := range s {
		fmt.Println(v)
	}

	


}

