package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	// svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
)

func Run() error {

	sysConn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer sysConn.Close()

	return nil
}
