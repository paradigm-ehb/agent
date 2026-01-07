package dbus_services

import (
	"fmt"

	v2 "paradigm-ehb/agent/gen/services/v2"
	dh "paradigm-ehb/agent/internal/dbusservices/dbus"
	svc "paradigm-ehb/agent/internal/dbusservices/systemd"
	svctypes "paradigm-ehb/agent/internal/dbusservices/types"

	"github.com/godbus/dbus"
)

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunAction(conn *dbus.Conn, ac svc.UnitAction, service string) error {

	obj := dh.CreateSystemdObject(conn)
	if !obj.Path().IsValid() {
		return fmt.Errorf("object path is invalid")
	}

	call := obj.Call(string(ac), 0, service, "replace")

	if call.Err != nil {
		return fmt.Errorf("failed to execute object on in unit action, %v", call.Err)
	}

	return nil
}

// @param, action [start, stop, restart], symLinkAction [enable, disable], service name format "example.service"
func RunSymlinkAction(conn *dbus.Conn, sc svc.UnitFileAction, enableForRunTime bool, enableForce bool, service []string) error {

	obj := dh.CreateSystemdObject(conn)

	if !obj.Path().IsValid() {
		fmt.Println("invalid systemd path")
	}

	/** EnableUnitFiles(in  as files, in  b runtime, in  b force, out b carries_install_info, out a(sss) changes); */
	/** DisableUnitFiles(in  as files, in  b runtime, out a(sss) changes); */

	switch sc {

	case svc.UnitFileActionEnable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime, enableForce)
		if call.Err != nil {
			return fmt.Errorf("error %v", call.Err)
		}
	case svc.UnitFileActionDisable:
		call := obj.Call(string(sc), dbus.FlagAllowInteractiveAuthorization, service, enableForRunTime)
		if call.Err != nil {
			return fmt.Errorf("something happened here %v", call.Err)
		}
	}

	return nil
}

// @param, true for all on disk, false for loaded units
func RunRetrieval(
	conn *dbus.Conn,
	all bool,
) ([]*v2.LoadedUnit, error) {

	obj := dh.CreateSystemdObject(conn)

	if all {
		ch := make(chan []svctypes.UnitFileEntry)
		parse := make(chan []svctypes.UnitFileEntry)

		go svc.GetAllUnits(obj, ch)
		go dh.ParseUnitFileEntries(ch, parse)

		entries := <-parse

		units := make([]*v2.LoadedUnit, 0, len(entries))
		for _, e := range entries {
			for _, i := range entries {

				/**

				TODO(nasr): replace this with the proper enum
				*/

				ch := make(chan []svctypes.LoadedUnit)
				parse := make(chan []svctypes.LoadedUnit)

				go svc.GetLoadedUnits(obj, ch)
				go dh.ParseLoadedUnits(ch, parse)

				loaded := <-parse

				if i.State == "enabled" {

					svc.GetLoadedUnits(obj, out)

					var luIdx int16
					var alIdx int16

					for luIdx, lu := range <-out {

					}

				}

				units = append(units, &v2.LoadedUnit{
					Name:        e.Name,
					Description: "",
					LoadState:   e.State,
					SubState:    "",
					ActiveState: "",
					DepUnit:     "",
					ObjectPath:  "",
					QueuedJob:   0,
					JobType:     "",
					JobPath:     "",
				})
			}

			loaded := <-parse

			units := make([]*v2.LoadedUnit, 0, len(loaded))
			for _, u := range loaded {
				units = append(units, &v2.LoadedUnit{
					Name:        u.Name,
					Description: u.Description,
					LoadState:   u.LoadState,
					SubState:    u.SubState,
					ActiveState: u.ActiveState,
					DepUnit:     u.DepUnit,
					ObjectPath:  string(u.ObjectPath),
					/*oops typo in queued job :)*/
					QueuedJob: u.QueudJob,
					JobType:   u.JobType,
					JobPath:   string(u.JobPath),
				})
			}

		}

		return units, nil
	}

	ch := make(chan []svctypes.LoadedUnit)
	parse := make(chan []svctypes.LoadedUnit)

	go svc.GetLoadedUnits(obj, ch)
	go dh.ParseLoadedUnits(ch, parse)

	loaded := <-parse

	units := make([]*v2.LoadedUnit, 0, len(loaded))
	for _, u := range loaded {
		units = append(units, &v2.LoadedUnit{
			Name:        u.Name,
			Description: u.Description,
			LoadState:   u.LoadState,
			SubState:    u.SubState,
			ActiveState: u.ActiveState,
			DepUnit:     u.DepUnit,
			ObjectPath:  string(u.ObjectPath),
			/*oops typo in queued job :)*/
			QueuedJob: u.QueudJob,
			JobType:   u.JobType,
			JobPath:   string(u.JobPath),
		})
	}

	return units, nil
}

func GetStatus(obj dbus.BusObject, name string) (string, error) {
	var result string

	call := obj.Call(
		"org.freedesktop.systemd1.Manager.GetUnitFileState",
		0,
		name,
	)

	if call.Err != nil {
		return "call error: ", call.Err
	}

	if err := call.Store(&result); err != nil {
		return "call store: ", err
	}

	fmt.Println("status:", result)

	return result, nil
}
