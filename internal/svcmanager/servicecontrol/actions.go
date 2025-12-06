// Package that handles dbus connection actions
package servicecontrol

type Action string

const (
	Start   Action = "org.freedesktop.systemd1.Manager.StartUnit"   // start unit
	Stop    Action = "org.freedesktop.systemd1.Manager.StopUnit"    // stop unit
	Restart Action = "org.freedesktop.systemd1.Manager.RestartUnit" // restart unit
)

type SymlinkAction string

const (
	Enable  SymlinkAction = "org.freedesktop.systemd1.Manager.EnableUnitFiles"  // enable unit(s)
	Disable SymlinkAction = "org.freedesktop.systemd1.Manager.DisableUnitFiles" // disable unit(s)

)
