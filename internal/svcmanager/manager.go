package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	// svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func Run() error {

	sysConn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer sysConn.Close()

	// TODO: introduce control flow to enable disable start restart stop
	return nil
}
