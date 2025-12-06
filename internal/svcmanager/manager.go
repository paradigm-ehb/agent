package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func Run(action svc.Action, symLinkAction svc.SymlinkAction, name string) error {

	sysConn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer sysConn.Close()

	// TODO: introduce control flow to enable disable start restart stop
	return nil
}
