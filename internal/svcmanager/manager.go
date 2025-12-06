package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
)

func Init() error {

	var sys System

	sysConn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer sysConn.Close()

	// DEBUG
	fmt.Println(sysConn.Names())

	obj := dh.CreateSystemdObject(sysConn)

	arrOnRam := svc.GetLoadedUnits(obj)
	// DEBUG
	fmt.Println("array of units", arrOnRam)

	arrOnDisk, err := svc.GetAllUnits(obj)
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
	err = svc.HandleActionOnUnit(obj, name, svc.Action(svc.Start))
	if err != nil {
		// DEBUG
		fmt.Println("failed to start the unit: ", err)
		err = svc.HandleSymlinkCreationAction(obj, namesList, svc.SymlinkAction(svc.Enable), enableForRunTime, replaceExistingSynmlink)
		if err != nil {
			// DEBUG
			fmt.Println("failed to enable the unit: ", err)
		}
	}

	// DEBUG
	fmt.Println(sys)

	return nil
}
