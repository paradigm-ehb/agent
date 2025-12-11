// Package that handles dbus connection actions
package servicecontrol

type UnitAction string

const (
	UnitActionStart   UnitAction = "org.freedesktop.systemd1.Manager.StartUnit"   // start unit
	UnitActionStop    UnitAction = "org.freedesktop.systemd1.Manager.StopUnit"    // stop unit
	UnitActionRestart UnitAction = "org.freedesktop.systemd1.Manager.RestartUnit" // restart unit
)

type UnitFileAction string

const (
	UnitFileActionEnable  UnitFileAction = "org.freedesktop.systemd1.Manager.EnableUnitFiles"  // enable unit(s)
	UnitFileActionDisable UnitFileAction = "org.freedesktop.systemd1.Manager.DisableUnitFiles" // disable unit(s)

)
