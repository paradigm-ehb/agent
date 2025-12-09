package svcmanager

import (
	"fmt"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
	svc "paradigm-ehb/agent/internal/svcmanager/servicecontrol"
	"sync"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
// TODO: add parameters and handling
func Run() error {

	conn, err := dh.CreateSystemBus()
	if err != nil {
		return fmt.Errorf("failed to create a systembus %v", err)
	}

	defer conn.Close()

	obj := dh.CreateSystemdObject(conn)

	var wg sync.WaitGroup
	ch := make(chan [][]string, 1)

	wg.Add(2)

	go svc.GetAllUnits(obj, ch)

	go dh.Parse(ch)

	wg.Wait()

	result := <-ch
	fmt.Println(result)

	// TODO: introduce control flow to enable disable start restart stop
	return nil
}
